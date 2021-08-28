package business

import (
	"github.com/astaxie/beego"
)

type WxAppBusiness struct {
	appId     string
	appSecret string
}

func NewWxAppBusiness() *WxAppBusiness {
	wxApp := new(WxAppBusiness)
	wxApp.appId = beego.AppConfig.String("wx_app_id")
	wxApp.appSecret = beego.AppConfig.String("wx_app_secret")
	return wxApp
}


