package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/models/repost"
	"dongchamao/models/repost/dy"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"
)

type AccountController struct {
	controllers.ApiBaseController
}

//登陆
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
func (receiver *AccountController) FindPwd() {
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

//重置密码
func (receiver *AccountController) ResetPwd() {
	InputData := receiver.InputFormat()
	if receiver.UserInfo.SetPassword == 1 {
		oldPwd := InputData.GetString("old_pwd", "")
		if oldPwd == "" {
			receiver.FailReturn(global.NewError(4207))
			return
		}
		oldPwd = utils.Base64Decode(oldPwd)
		if utils.Md5_encode(oldPwd+receiver.UserInfo.Salt) != receiver.UserInfo.Password {
			receiver.FailReturn(global.NewError(4207))
			return
		}
	}
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
	pwd := utils.Md5_encode(newPwd + receiver.UserInfo.Salt)
	userBusiness := business.NewUserBusiness()
	updateData := map[string]interface{}{
		"password":     pwd,
		"set_password": 1,
		"update_time":  utils.GetNowTimeStamp(),
	}
	affect, _ := userBusiness.UpdateUserAndClearCache(nil, receiver.UserId, updateData)
	if affect == 0 {
		receiver.FailReturn(global.NewError(4213))
		return
	}
	receiver.Logout()
	receiver.SuccReturn(nil)
	return
}

func (receiver *AccountController) Info() {
	username := receiver.UserInfo.Username
	account := dy.RepostAccountData{
		UserId:      receiver.UserInfo.Id,
		Username:    username[:3] + "****" + username[7:],
		Nickname:    receiver.UserInfo.Nickname,
		Avatar:      receiver.UserInfo.Avatar,
		PasswordSet: receiver.UserInfo.SetPassword,
	}
	vipBusiness := business.NewVipBusiness()
	vipLevelsMap := vipBusiness.GetVipLevels(receiver.UserInfo.Id)
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
	receiver.SuccReturn(map[string]interface{}{
		"info": account,
	})
	return
}

//登出
func (receiver *AccountController) Logout() {
	cacheKey := cache.GetCacheKey(cache.UserPlatformUniqueToken, receiver.AppId, receiver.UserId)
	_ = global.Cache.Delete(cacheKey)
	//执行登出事件
	receiver.RegisterLogout()
	//uniquetoken更新置为空  旧的token不可用
	userBusiness := business.NewUserBusiness()
	_ = userBusiness.AddOrUpdateUniqueToken(receiver.UserId, receiver.AppId, "")
	userBusiness.DeleteUserInfoCache(receiver.UserInfo.Id)
	receiver.SuccReturn("success")
	return

}

func (receiver *AccountController) DyUserSearchSave() {
	searchType := receiver.GetString(":type")
	data := receiver.ApiDatas
	dataMap, _ := utils.ToMapStringInterface(data)
	searchData := map[string]interface{}{}
	note := ""
	for k, v := range dataMap {
		if k == "note" {
			note = utils.ToString(v)
			continue
		}
		searchData[k] = v
	}
	if note == "" {
		receiver.FailReturn(global.NewMsgError("请输入筛选器昵称"))
		return
	}
	contentByte, _ := jsoniter.Marshal(searchData)
	searchM := dcm.DcUserSearch{
		UserId:     receiver.UserId,
		SearchType: searchType,
		Note:       note,
		Content:    string(contentByte),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	affect, err := dcm.Insert(nil, &searchM)
	if affect == 0 || err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(nil)
	return
}

func (receiver *AccountController) DyUserSearchDel() {
	id := receiver.GetString(":id")
	dbSession := dcm.GetDbSession()
	searchM := dcm.DcUserSearch{}
	affect, err := dbSession.Where("id = ? AND user_id = ?", id, receiver.UserId).Delete(&searchM)
	if affect == 0 || err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(nil)
	return
	receiver.SuccReturn(nil)
	return
}

func (receiver *AccountController) DyUserSearchList() {
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	searchType := receiver.GetString(":type")
	list := make([]dcm.DcUserSearch, 0)
	dbSession := dcm.GetSlaveDbSession()
	start := (page - 1) * pageSize
	total, err := dbSession.
		Where("search_type = ?", searchType).
		Limit(pageSize, start).
		FindAndCount(&list)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	repostList := make([]repost.SearchData, 0)
	for _, v := range list {
		content := map[string]interface{}{}
		_ = jsoniter.Unmarshal([]byte(v.Content), &content)
		repostList = append(repostList, repost.SearchData{
			SearchType: v.SearchType,
			Note:       v.Note,
			Content:    content,
		})
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  repostList,
		"total": total,
	})
	return
}
