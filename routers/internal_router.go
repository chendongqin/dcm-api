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
			beego.NSRouter("/info/:author_id", &controllers.InternalController{}, "get:AuthorInfo"),
		),
		beego.NSNamespace("/product",
			beego.NSRouter("/search", &controllers.InternalController{}, "get:ProductSearch"),
			beego.NSRouter("/cate/:product_id", &controllers.InternalController{}, "post:ChangeProductCate"),
		),
		beego.NSNamespace("/system",
			beego.NSRouter("/cache/clear", &controllers.InternalController{}, "post:ClearCache"),
		),
		beego.NSNamespace("/config",
			beego.NSRouter("/:key_name", &controllers.InternalController{}, "get:GetConfig"),
		),
		beego.NSNamespace("/wechat",
			beego.NSRouter("/menu", &controllers.InternalController{}, "get:GetWeChatMenu"),
			beego.NSRouter("/menu/set", &controllers.InternalController{}, "post:SetWeChatMenu"),
			beego.NSRouter("/media/upload", &controllers.InternalController{}, "post:UploadWeChatMedia"),
			beego.NSRouter("/media/list", &controllers.InternalController{}, "get:GetWeChatMediaList"),
			beego.NSRouter("/media/del", &controllers.InternalController{}, "get:DelWeChatMedia"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
