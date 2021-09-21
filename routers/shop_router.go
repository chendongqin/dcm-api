package routers

import (
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	//直播
	ns := beego.NewNamespace("/v1/dy",
		beego.NSNamespace("/shop",
			beego.NSRouter("search", &v1dy.ShopController{}, "get:SearchBase"),
			beego.NSRouter("/info/:shop_id", &v1dy.ShopController{}, "get:ShopBase"),
			beego.NSRouter("/analysis/:shop_id/:start/:end", &v1dy.ShopController{}, "get:ShopBaseAnalysis"),
			beego.NSRouter("/product/analysis/:shop_id/:start/:end", &v1dy.ShopController{}, "get:ShopProductAnalysis"),
			beego.NSRouter("/product/analysis/count/:shop_id/:start/:end", &v1dy.ShopController{}, "get:ShopProductAnalysisCount"),
			beego.NSRouter("/author/gmv/top/:shop_id/:start/:end", &v1dy.ShopController{}, "get:ShopAuthorGmvRate"),
			beego.NSRouter("/author/live/:shop_id/:start/:end", &v1dy.ShopController{}, "get:ShopLiveAuthorAnalysis"),
			beego.NSRouter("/author/live/count/:shop_id/:start/:end", &v1dy.ShopController{}, "get:ShopLiveAuthorAnalysisCount"),
			beego.NSRouter("/author/live/room/:shop_id/:author_id/:start/:end", &v1dy.ShopController{}, "get:ShopLiveAuthorRooms"),
			beego.NSRouter("/author/aweme/:shop_id/:start/:end", &v1dy.ShopController{}, "get:ShopAwemeAuthorAnalysis"),
			beego.NSRouter("/author/aweme/count/:shop_id/:start/:end", &v1dy.ShopController{}, "get:ShopAwemeAuthorAnalysisCount"),
			beego.NSRouter("/author/aweme/list/:shop_id/:author_id/:start/:end", &v1dy.ShopController{}, "get:ShopAuthorAwemes"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
