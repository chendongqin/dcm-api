package v1

import (
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/business"
	"dongchamao/models/dcm"
	"dongchamao/structinit/repost/dy"
	"strings"
	"time"
)

type AccountController struct {
	controllers.ApiBaseController
}

func (receiver *AccountController) Login() {
	InputData := receiver.InputFormat()
	grantType := InputData.GetString("grant_type", "password")
	appId := InputData.GetInt("app_id", 0)
	if !utils.InArrayInt(appId, []int{10000, 10001, 10002, 10003, 10004, 10005}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var user dcm.DcUser
	var comErr global.CommonError
	var authToken string
	var expTime int64
	userBusiness := business.NewUserBusiness()
	if grantType == "password" {
		username := InputData.GetString("username", "")
		password := InputData.GetString("pwd", "")
		user, authToken, expTime, comErr = userBusiness.LoginByPwd(username, password, appId)
	} else if grantType == "sms" {
		username := InputData.GetString("username", "")
		code := InputData.GetString("code", "")
		user, authToken, expTime, comErr = userBusiness.SmsLogin(username, code, appId)
	} else {
		comErr = global.NewError(4000)
	}
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	account := dy.RepostAccountData{
		UserId:   user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}
	vipBusiness := business.NewVipBusiness()
	vipLevelsMap := vipBusiness.GetVipLevels(user.Id)
	for k, v := range vipLevelsMap {
		if k == business.VipPlatformDouYin {
			account.DyLevel.Level = v
		} else if k == business.VipPlatformXiaoHongShu {
			account.XhsLevel.Level = v
		} else if k == business.VipPlatformTaoBao {
			account.TbLevel.Level = v
		}
	}
	account.DyLevel.LevelName = vipBusiness.GetUserLevel(account.DyLevel.Level)
	account.XhsLevel.LevelName = vipBusiness.GetUserLevel(account.XhsLevel.Level)
	account.TbLevel.LevelName = vipBusiness.GetUserLevel(account.TbLevel.Level)
	updateData := map[string]interface{}{
		"login_time": utils.GetNowTimeStamp(),
		"login_ip":   receiver.Ip,
	}
	today := time.Now().Format("20060102")
	yesterdayTime := time.Now().AddDate(0, 0, -1)
	prevTimeDate := user.PrevTime.Format("20060102")
	if today != prevTimeDate {
		updateData["prev_time"] = utils.GetNowTimeStamp()
		if yesterdayTime.Format("20060102") == prevTimeDate || user.Successions == 0 {
			successions := user.Successions + 1
			totalSuccessions := user.TotalSuccessions + 1
			updateData["successions"] = successions
			updateData["total_successions"] = totalSuccessions
			if successions > user.MaxSuccessions {
				updateData["max_successions"] = successions
			}
		}
	}
	dcm.UpdateInfo(nil, user.Id, updateData, new(dcm.DcUser))
	receiver.SuccReturn(map[string]interface{}{
		"account": account,
		"token_info": dy.RepostAccountToken{
			UserId:      user.Id,
			TokenString: authToken,
			ExpTime:     expTime,
		},
	})
	return
}

func (receiver *AccountController) ResetPwd() {
	InputData := receiver.InputFormat()
	mobile := InputData.GetString("username", "")
	if !utils.VerifyMobileFormat(mobile) {
		receiver.FailReturn(global.NewError(4205))
		return
	}
	code := InputData.GetString("code", "")
	newPwd := InputData.GetString("new_pwd", "")
	surePwd := InputData.GetString("sure_pwd", "")
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
	receiver.SuccReturn(nil)
	return
}
