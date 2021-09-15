package routers

import (
	v1 "dongchamao/controllers/api/v1"
	"github.com/astaxie/beego"
)

func init() {
	//用户相关
	ns := beego.NewNamespace("/v1/dy",
		beego.NSNamespace("/team",
			beego.NSRouter("/list", &v1.VipController{}, "get:GetDyTeam"),
			beego.NSRouter("/add", &v1.VipController{}, "put:AddDyTeamSub"),
			beego.NSRouter("/remove", &v1.VipController{}, "put:RemoveDyTeam"),
			beego.NSRouter("/remark", &v1.VipController{}, "put:AddDySubRemark"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
