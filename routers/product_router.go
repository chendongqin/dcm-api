package routers

import (
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1/dy",
		beego.NSNamespace("/product",
			beego.NSRouter("/cate", &v1dy.ProductController{}, "get:GetCacheProductCate"),
			beego.NSRouter("/search", &v1dy.ProductController{}, "get:Search"),
			beego.NSRouter("/base/:product_id", &v1dy.ProductController{}, "get:ProductBase"),
			beego.NSRouter("/analysis/:product_id/:start/:end", &v1dy.ProductController{}, "get:ProductBaseAnalysis"),
			beego.NSRouter("/live/chart/:product_id/:start/:end", &v1dy.ProductController{}, "get:ProductLiveChart"),
			beego.NSRouter("/live/room/:product_id/:start/:end", &v1dy.ProductController{}, "get:ProductLiveRoomList"),
			beego.NSRouter("/author/:product_id/:start/:end", &v1dy.ProductController{}, "get:ProductLiveAuthorAnalysis"),
			beego.NSRouter("/author/count/:product_id/:start/:end", &v1dy.ProductController{}, "get:ProductLiveAuthorAnalysisCount"),
			beego.NSRouter("/author/room/:product_id/:author_id/:start/:end", &v1dy.ProductController{}, "get:ProductAuthorLiveRooms"),
			beego.NSRouter("/room/:product_id/:start/:end", &v1dy.ProductController{}, "get:ProductRoomsRangeDate"),
			beego.NSRouter("/aweme/sales/chart/:product_id/:start/:end", &v1dy.ProductController{}, "get:ProductAwemeSalesTrend"),
			beego.NSRouter("/aweme/list/:product_id/:start/:end", &v1dy.ProductController{}, "get:ProductAweme"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
