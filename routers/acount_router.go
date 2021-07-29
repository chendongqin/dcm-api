package routers

import (
	v1 "dongchamao/controllers/api/v1"
	"github.com/astaxie/beego"
)

func init() {
	//用户相关
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/account",
			beego.NSRouter("/password", &v1.AccountController{}, "put:ResetPwd"),
			beego.NSRouter("/info", &v1.AccountController{}, "get:Info"),
			beego.NSRouter("/logout", &v1.AccountController{}, "get:Logout"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
