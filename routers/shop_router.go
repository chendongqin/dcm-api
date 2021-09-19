package routers

import (
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	//直播
	ns := beego.NewNamespace("/v1/dy",
		beego.NSNamespace("/shop",
			beego.NSRouter("/info/:shop_id", &v1dy.ShopController{}, "get:ShopBase"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
