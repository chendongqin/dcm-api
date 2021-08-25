package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"github.com/astaxie/beego/context"
	"strings"
	"time"
)

const SignKey = "We6BkysyuB5RNed4"

type AccountAuthBusiness struct {
}

func NewAccountAuthBusiness() *AccountAuthBusiness {
	return new(AccountAuthBusiness)
}

//无参数路由不需登陆白名单
var LoginWitheUri = []string{
	"/v1/id",
	"/v1/user/login",
	"/v1/user/findpwd",
	"/v1/sms/code",
	"/v1/sms/verify",
	"/v1/dy/author/live/tags",
	"/v1/dy/author/cate",
	"/v1/dy/product/cate",
	"/v1/dy/author/search",
	"/v1/dy/live/search",
	"/v1/pay/notify/wechat",
	"/v1/pay/notify/alipay",
	"/v1/callback/wechat",
}

var SignWitheUri = []string{
	"/v1/id",
	"/v1/pay/notify/wechat",
	"/v1/pay/notify/alipay",
	"/v1/callback/wechat",
}

var AuthDyWitheUriMap = []string{}

//登陆白名单校验
func (receiver *AccountAuthBusiness) AuthDyWhiteUri(uri string, level int) bool {
	if level > 0 {
		return true
	}
	if utils.InArrayString(uri, LoginWitheUri) || utils.InArrayString(uri, AuthDyWitheUriMap) {
		return true
	}
	return false
}

//登陆白名单校验
func (receiver *AccountAuthBusiness) AuthLoginWhiteUri(uri string) bool {
	if utils.InArrayString(uri, LoginWitheUri) {
		return true
	}
	return false
}

//签名白名单校验
func (receiver *AccountAuthBusiness) AuthSignWhiteUri(uri string) bool {
	if utils.InArrayString(uri, SignWitheUri) {
		return true
	}
	return false
}

//获取真实uri
func (receiver *AccountAuthBusiness) GetTrueRequestUri(uri string, params map[string]string) string {
	urlData := strings.Split(uri, "?")
	if len(urlData) > 0 {
		uri = urlData[0]
	}
	for _, v := range params {
		uri = strings.Replace(uri, "/"+v, "", 1)
	}
	return uri
}

//校验签名
func (receiver *AccountAuthBusiness) CheckSign(timestamp, random, sign string) global.CommonError {
	tmpStr := timestamp + random + SignKey
	if sign != utils.Md5_encode(tmpStr) {
		return global.NewError(4041)
	}
	nowTime := time.Now().Unix() - 60
	nowTime2 := time.Now().Unix() + 60
	timestampInt64 := utils.ToInt64(timestamp)
	if timestampInt64 < nowTime || timestampInt64 > nowTime2 {
		return global.NewError(4041)
	}
	return nil
}

func (receiver *AccountAuthBusiness) CheckAppIdSign(appId string, ctx *context.Context) global.CommonError {
	timestamp := ctx.Input.Header("TIMESTAMP")
	random := ctx.Input.Header("RANDOM")
	sign := ctx.Input.Header("SIGN")
	secret, _ := receiver.GetAppSecret(appId, true)
	if secret == "" {
		return global.NewError(4041)
	}
	tmpStr := timestamp + random + secret
	if sign != utils.Md5_encode(tmpStr) {
		return global.NewError(4041)
	}
	nowTime := time.Now().Unix() - 60
	nowTime2 := time.Now().Unix() + 60
	timestampInt64 := utils.ToInt64(timestamp)
	if timestampInt64 < nowTime || timestampInt64 > nowTime2 {
		return global.NewError(4041)
	}
	return nil
}

func (receiver *AccountAuthBusiness) GetAppSecret(appId string, enableCache bool) (string, bool) {
	cKey := cache.GetCacheKey(cache.AppIdSecret, appId)

	if enableCache == true {
		secret := global.Cache.Get(cKey)
		if secret != "" {
			return secret, true
		}
	}
	model := &dcm.DcAppid{}
	exist, _ := dcm.GetSlaveDbSession().Where("app_id = ?", appId).
		Get(model)

	if exist {
		_ = global.Cache.Set(cKey, model.Secret, 1800)
	}

	return model.Secret, exist
}
