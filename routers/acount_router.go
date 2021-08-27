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
			beego.NSRouter("/mobile/change", &v1.AccountController{}, "put:ChangeMobile"),
			beego.NSRouter("/mobile/bind", &v1.AccountController{}, "put:BindMobile"),
			beego.NSRouter("/dy/search/:type", &v1.AccountController{}, "put:DyUserSearchSave"),
			beego.NSRouter("/dy/search/:id", &v1.AccountController{}, "delete:DyUserSearchDel"),
			beego.NSRouter("/dy/search/:type", &v1.AccountController{}, "get:DyUserSearchList"),
			beego.NSRouter("/collect/list", &v1.AccountController{}, "get:GetCollect"),
			beego.NSRouter("/collect/add", &v1.AccountController{}, "put:AddCollect"),
			beego.NSRouter("/collect/del/:id", &v1.AccountController{}, "delete:DelCollect"),
			beego.NSRouter("/collect/dy/tag/list", &v1.AccountController{}, "get:GetDyCollectTags"),
			beego.NSRouter("/collect/dy/tag/add", &v1.AccountController{}, "put:AddDyCollectTag"),
			beego.NSRouter("/collect/dy/tag/upd/:id", &v1.AccountController{}, "put:UpdDyCollectTag"),
			beego.NSRouter("/collect/dy/tag/del/:id", &v1.AccountController{}, "delete:DelDyCollectTag"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
