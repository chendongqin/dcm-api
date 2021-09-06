package routers

import (
	controllers "dongchamao/controllers/api"
	"github.com/astaxie/beego"
)

func init() {
	//视频
	ns := beego.NewNamespace("/internal",
		beego.NSNamespace("/author",
			beego.NSRouter("/cate/:author_id", &controllers.InternalController{}, "post:ChangeAuthorCate"),
			beego.NSRouter("/product", &controllers.InternalController{}, "get:ProductSearch"),
			beego.NSRouter("/search", &controllers.InternalController{}, "get:AuthorSearch"),
		),
		beego.NSNamespace("/product",
			beego.NSRouter("/cate/:product_id", &controllers.InternalController{}, "post:ChangeProductCate"),
		),
		beego.NSNamespace("/system",
			beego.NSRouter("/cache/clear", &controllers.InternalController{}, "post:ClearCache"),
		),
		beego.NSNamespace("/config",
			beego.NSRouter("/:key_name", &controllers.InternalController{}, "get:GetConfig"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
