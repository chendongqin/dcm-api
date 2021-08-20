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
			beego.NSRouter("/dy/search/:type", &v1.AccountController{}, "put:DyUserSearchSave"),
			beego.NSRouter("/dy/search/:id", &v1.AccountController{}, "delete:DyUserSearchDel"),
			beego.NSRouter("/dy/search/:type", &v1.AccountController{}, "get:DyUserSearchList"),
			beego.NSRouter("/dy/add_collect", &v1.AccountController{}, "get:AddDyCollect"),
			beego.NSRouter("/dy/del_collect/:id", &v1.AccountController{}, "delete:DelDyCollect"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
