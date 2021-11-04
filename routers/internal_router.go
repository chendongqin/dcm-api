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
			beego.NSRouter("/search/ids", &controllers.InternalController{}, "post:AuthorSearchByIds"),
		),
		beego.NSNamespace("/product",
			beego.NSRouter("/search", &controllers.InternalController{}, "get:ProductSearch"),
			beego.NSRouter("/cate/:product_id", &controllers.InternalController{}, "post:ChangeProductCate"),
			beego.NSRouter("/search/ids", &controllers.InternalController{}, "post:ProductSearchByIds"),
		),
		beego.NSNamespace("/spider",
			beego.NSRouter("/red/author/:author_id", &controllers.InternalController{}, "get:SpiderLiveSpeedUp"),
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
		beego.NSNamespace("/decrypt",
			beego.NSRouter("/id/:id", &controllers.InternalController{}, "get:IdEncryptDecrypt"),
			beego.NSRouter("/json", &controllers.InternalController{}, "post:JsonDecrypt"),
		),
		beego.NSNamespace("/log",
			beego.NSRouter("/url", &controllers.InternalController{}, "get:CommonUrlLog"),
			beego.NSRouter("/speed/:type/:days/:end_time/:page", &controllers.InternalController{}, "get:SpeedUp"),
		),
		beego.NSNamespace("/live",
			beego.NSRouter("/room/search", &controllers.InternalController{}, "get:SearchLiveRooms"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
