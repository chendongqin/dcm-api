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
			beego.NSRouter("/search", &controllers.InternalController{}, "get:AuthorSearch"),
		),
		beego.NSNamespace("/product",
			beego.NSRouter("/cate/:product_id", &controllers.InternalController{}, "post:ChangeProductCate"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
