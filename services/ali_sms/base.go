package AliSms

import (
	"dongchamao/global"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	jsoniter "github.com/json-iterator/go"
)

const AliSign = "维妥科技"
const CodeTemplateCode = "SMS_218595039"

type AliRe struct {
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
	Code      string `json:"Code"`
	BizId     string `json:"BizId,omitempty"`
}

func SmsSend(phoneNum, templateCode, signName string, templateParam map[string]string) (bool, error) {
	aliAccessKey := global.Cfg.String("ali_sms_accessKey")
	aliAccessSecret := global.Cfg.String("ali_sms_secret")
	if aliAccessKey == "" || aliAccessSecret == "" {
		return false, errors.New("配置加载失败")
	}
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", aliAccessKey, aliAccessSecret)
	if err != nil {
		return false, err
	}
	templateParamStr, _ := jsoniter.Marshal(templateParam)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = phoneNum
	request.QueryParams["SignName"] = signName
	request.QueryParams["TemplateCode"] = templateCode              //短信发送模版
	request.QueryParams["TemplateParam"] = string(templateParamStr) //短信发送模版变量

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return false, err
	}
	re := AliRe{}
	err = jsoniter.Unmarshal(response.GetHttpContentBytes(), &re)
	if err != nil {
		return false, err
	}
	if re.Code == "OK" && re.Message == "OK" {
		return true, nil
	}
	return false, errors.New(re.Message)
}


func SmsCode(phoneNum, code string) (bool, error) {
	templateParam := map[string]string{
		"code":      code,
	}
	return SmsSend(phoneNum, CodeTemplateCode, AliSign, templateParam)
}

