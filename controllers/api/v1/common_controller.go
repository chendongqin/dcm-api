package v1

import (
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/services/ali_sms"
)

type CommonController struct {
	controllers.ApiBaseController
}

func (receiver *CommonController) Sms() {
	InputData := receiver.InputFormat()
	grantType := InputData.GetString("grant_type", "")
	mobile := InputData.GetString("mobile", "")
	if !utils.VerifyMobileFormat(mobile) {
		receiver.FailReturn(global.NewError(4205))
		return
	}
	if !utils.InArrayString(grantType, []string{"login", "findpwd"}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	limitIpKey := cache.GetCacheKey(cache.SmsCodeLimitBySome, "ip", receiver.Ip)
	verifyData := global.Cache.Get(limitIpKey)
	if verifyData != "" {
		receiver.FailReturn(global.NewError(4211))
		return
	}
	limitMobileKey := cache.GetCacheKey(cache.SmsCodeLimitBySome, "mobile", mobile)
	verifyData = global.Cache.Get(limitMobileKey)
	if verifyData != "" {
		receiver.FailReturn(global.NewError(4211))
		return
	}
	if utils.InArrayString(grantType, []string{"findpwd"}) {
		user := dcm.DcUser{}
		exist, _ := dcm.GetBy("username", mobile, &user)
		if !exist {
			receiver.FailReturn(global.NewError(4204))
			return
		}
	}
	cacheKey := cache.GetCacheKey(cache.SmsCodeVerify, grantType, mobile)
	code := utils.GetRandomInt(6)
	err := global.Cache.Set(cacheKey, code, 300)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	res, smsErr := aliSms.SmsCode(mobile, code)
	if !res || smsErr != nil {
		receiver.FailReturn(global.NewError(6000))
		return
	}
	global.Cache.Set(limitIpKey, "1", 60)
	global.Cache.Set(limitMobileKey, "1", 60)
	receiver.SuccReturn(nil)
	return
}

//验证码校验
func (receiver *CommonController) CheckSmsCode() {
	mobile := receiver.GetString(":username", "")
	grantType := receiver.GetString(":grant_type", "")
	if !utils.VerifyMobileFormat(mobile) {
		receiver.FailReturn(global.NewError(4205))
		return
	}
	if !utils.InArrayString(grantType, []string{"findpwd"}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	code := receiver.GetString(":code", "")
	codeKey := cache.GetCacheKey(cache.SmsCodeVerify, grantType, mobile)
	verifyCode := global.Cache.Get(codeKey)
	if verifyCode != code {
		receiver.FailReturn(global.NewError(4209))
		return
	}
	receiver.SuccReturn(nil)
	return
}

func (receiver *CommonController) Test() {
	receiver.SuccReturn(nil)
	return
}
