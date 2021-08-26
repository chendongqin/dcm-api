package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/models/repost/dy"
	"strings"
)

type LoginController struct {
	controllers.ApiBaseController
}

//登陆
func (receiver *LoginController) Login() {
	InputData := receiver.InputFormat()
	grantType := InputData.GetString("grant_type", "password")
	appId := receiver.AppId
	if !utils.InArrayInt(appId, []int{10000, 10001, 10002, 10003, 10004, 10005}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var user dcm.DcUser
	var comErr global.CommonError
	var authToken string
	var expTime int64
	var isNew int
	setPassword := 0
	userBusiness := business.NewUserBusiness()
	if grantType == "password" {
		username := InputData.GetString("username", "")
		password := InputData.GetString("pwd", "")
		password = utils.Base64Decode(password)
		user, authToken, expTime, comErr = userBusiness.LoginByPwd(username, password, appId)
	} else if grantType == "sms" {
		username := InputData.GetString("username", "")
		code := InputData.GetString("code", "")
		user, authToken, expTime, isNew, comErr = userBusiness.SmsLogin(username, code, appId)
		if isNew == 0 && user.SetPassword == 0 {
			setPassword = 1
		}
	} else if grantType == "wechat" {
		openid := InputData.GetString("openid", "")
		user, authToken, expTime, isNew, comErr = userBusiness.QrLogin(openid, appId)
	} else {
		comErr = global.NewError(4000)
	}
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	updateData := map[string]interface{}{
		"login_time": utils.GetNowTimeStamp(),
		"login_ip":   receiver.Ip,
	}
	if business.AppIdMap[appId] == 2 {
		updateData["is_install_app"] = 1
	}
	_, _ = userBusiness.UpdateUserAndClearCache(nil, user.Id, updateData)
	receiver.RegisterLogin(authToken, expTime)
	receiver.SuccReturn(map[string]interface{}{
		"set_password": setPassword,
		"token_info": dy.RepostAccountToken{
			UserId:      user.Id,
			TokenString: authToken,
			ExpTime:     expTime,
		},
	})
	return
}

//找回密码
func (receiver *LoginController) FindPwd() {
	InputData := receiver.InputFormat()
	mobile := InputData.GetString("username", "")
	if !utils.VerifyMobileFormat(mobile) {
		receiver.FailReturn(global.NewError(4205))
		return
	}
	code := InputData.GetString("code", "")
	newPwd := InputData.GetString("new_pwd", "")
	newPwd = utils.Base64Decode(newPwd)
	surePwd := InputData.GetString("sure_pwd", "")
	surePwd = utils.Base64Decode(surePwd)
	if newPwd == "" || surePwd == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if newPwd != surePwd {
		receiver.FailReturn(global.NewError(4207))
		return
	}
	pwdLen := strings.Count(newPwd, "")
	if pwdLen > 24 || pwdLen < 6 {
		receiver.FailReturn(global.NewError(4210))
		return
	}
	codeKey := cache.GetCacheKey(cache.SmsCodeVerify, "findpwd", mobile)
	verifyCode := global.Cache.Get(codeKey)
	if verifyCode != code {
		receiver.FailReturn(global.NewError(4209))
		return
	}
	user := dcm.DcUser{}
	exist, _ := dcm.GetBy("username", mobile, &user)
	if !exist {
		receiver.FailReturn(global.NewError(4209))
		return
	}
	pwd := utils.Md5_encode(newPwd + user.Salt)
	affect, _ := dcm.UpdateInfo(nil, user.Id, map[string]interface{}{
		"password":     pwd,
		"set_password": 1,
		"update_time":  utils.GetNowTimeStamp(),
	}, new(dcm.DcUser))
	if affect == 0 {
		receiver.FailReturn(global.NewError(4213))
		return
	}
	_ = global.Cache.Delete(codeKey)
	receiver.SuccReturn(nil)
	return
}
