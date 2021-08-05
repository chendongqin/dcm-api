package routers

import (
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	//直播
	ns := beego.NewNamespace("/v1/dy",
		beego.NSNamespace("/live",
			beego.NSRouter("/search", &v1dy.LiveController{}, "get:SearchRoom"),
			beego.NSRouter("/info/:room_id", &v1dy.LiveController{}, "get:LiveInfoData"),
			beego.NSRouter("/promotion/sale/:room_id/:product_id", &v1dy.LiveController{}, "get:LiveProductSaleChart"),
			beego.NSRouter("/promotion/list/:room_id", &v1dy.LiveController{}, "get:LiveProductList"),
			beego.NSRouter("/promotion/cate/:room_id", &v1dy.LiveController{}, "get:LiveProductCateList"),
			beego.NSRouter("/promotion/chart/:room_id", &v1dy.LiveController{}, "get:LivePromotions"),
			beego.NSRouter("/rank/chart/:room_id", &v1dy.LiveController{}, "get:LiveRankTrends"),
			beego.NSRouter("/fans/chart/:room_id", &v1dy.LiveController{}, "get:LiveFansTrends"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
