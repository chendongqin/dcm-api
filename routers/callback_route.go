package routers

import (
	v1 "dongchamao/controllers/api/v1"
	"github.com/astaxie/beego"
)

func init() {
	//第三方回调
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/callback",
			beego.NSRouter("/wechat", &v1.WechatController{}, "*:Receive"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
