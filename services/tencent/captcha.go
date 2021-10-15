package tencent

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	captcha "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func TencentCaptcha(ticket, randStr, ip string) bool {
	credential := common.NewCredential(global.Cfg.String("tencentSecretId"), global.Cfg.String("tencentSecretKey"))
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "captcha.tencentcloudapi.com"
	client, _ := captcha.NewClient(credential, "", cpf)
	request := captcha.NewDescribeCaptchaResultRequest()
	captchaAppIdTemp := utils.ToInt64(global.Cfg.String("tencentCaptchaAppId"))
	appSecretKey := global.Cfg.String("tencentAppSecretKey")
	var captchaType uint64 = 9
	request.AppSecretKey = &appSecretKey
	captchaAppId := uint64(captchaAppIdTemp)
	request.CaptchaAppId = &captchaAppId
	request.CaptchaType = &captchaType
	request.Ticket = &ticket
	request.Randstr = &randStr
	request.UserIp = &ip
	response, err := client.DescribeCaptchaResult(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return false
	}
	if err != nil {
		return false
	}
	if *response.Response.CaptchaCode == 1 {
		return true
	}
	return false
}
