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
			beego.NSRouter("/cancel", &v1.AccountController{}, "get:Cancel"),
			beego.NSRouter("/wechat/bind", &v1.AccountController{}, "put:BindWeChat"),
			beego.NSRouter("/mobile/change", &v1.AccountController{}, "put:ChangeMobile"),
			beego.NSRouter("/mobile/exist", &v1.AccountController{}, "get:MobileExist"),
			beego.NSRouter("/dy/search/:type", &v1.AccountController{}, "put:DyUserSearchSave"),
			beego.NSRouter("/dy/search/:id", &v1.AccountController{}, "delete:DyUserSearchDel"),
			beego.NSRouter("/dy/search/:type", &v1.AccountController{}, "get:DyUserSearchList"),
			beego.NSNamespace("/collect",
				beego.NSRouter("/list", &v1.AccountController{}, "get:GetCollect"),
				beego.NSRouter("/add", &v1.AccountController{}, "put:AddCollect"),
				beego.NSRouter("/exist", &v1.AccountController{}, "get:IsCollect"),
				beego.NSRouter("/label/get", &v1.AccountController{}, "get:DyCollectLabel"),
				beego.NSRouter("/tag/upd/:id/:tag_id", &v1.AccountController{}, "put:UpdCollectTag"),
				beego.NSRouter("/del/:id", &v1.AccountController{}, "delete:DelCollect"),
				beego.NSNamespace("/dy",
					beego.NSRouter("/remark", &v1.AccountController{}, "put:DyCollectRemark"),
					beego.NSRouter("/tag/list", &v1.AccountController{}, "get:GetDyCollectTags"),
					beego.NSRouter("/tag/add", &v1.AccountController{}, "put:AddDyCollectTag"),
					beego.NSRouter("/tag/upd/:id", &v1.AccountController{}, "put:UpdDyCollectTag"),
					beego.NSRouter("/tag/del/:id", &v1.AccountController{}, "delete:DelDyCollectTag"),
				),
			),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
