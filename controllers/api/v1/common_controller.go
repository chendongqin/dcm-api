package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/logger"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/services/ali_sms"
	"encoding/json"
	"strings"
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
	if logger.CheckError(err) != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	res, smsErr := aliSms.SmsCode(mobile, code)
	if !res || logger.CheckError(smsErr) != nil {
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

func (receiver *CommonController) IdEncryptDecrypt() {
	id := receiver.Ctx.Input.Param(":id")
	id1 := ""
	if strings.Index(id, "=") < 0 {
		id1 = business.IdEncrypt(id)
	}
	id2 := business.IdDecrypt(id)
	receiver.SuccReturn(map[string]string{
		"id":      id,
		"encrypt": id1,
		"decrypt": id2,
	})
	return
}

func (receiver *CommonController) Test() {
	return
}

func (receiver *CommonController) GetConfig() {
	var configJson dcm.DcConfigJson
	keyName := receiver.GetString(":key_name")
	exist, err := dcm.GetBy("key_name", keyName, &configJson)
	if !exist || err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	if configJson.Auth == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(configJson.Value), &data); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(data)
	return
}
