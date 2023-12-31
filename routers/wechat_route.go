package routers

import (
	v1 "dongchamao/controllers/api/v1"
	"github.com/astaxie/beego"
)

func init() {
	//视频
	ns := beego.NewNamespace("/v1/wechat",
		beego.NSRouter("/qrcode", &v1.WechatController{}, "get:QrCode"),
		beego.NSRouter("/check", &v1.WechatController{}, "get:CheckScan"),
		beego.NSRouter("/phone", &v1.WechatController{}, "put:WechatPhone"),
		beego.NSRouter("/app", &v1.WechatController{}, "get:WechatApp"),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
