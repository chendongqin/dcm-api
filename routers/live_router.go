package routers

import (
	v1 "dongchamao/controllers/api/v1"
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
			beego.NSRouter("/fans/data/:type/:room_id", &v1dy.LiveController{}, "get:LiveFanAnalyse"),
			beego.NSRouter("/fans/product/:room_id", &v1dy.LiveController{}, "get:LiveProductPvAnalyse"),
			beego.NSRouter("/monitor/price", &v1.LiveMonitorController{}, "get:MonitorPrice"),
			beego.NSRouter("/monitor/add", &v1.LiveMonitorController{}, "put:AddLiveMonitor"),
			beego.NSRouter("/monitor/calc/:start/:end", &v1.LiveMonitorController{}, "get:LiveMonitorCalcCount"),
			beego.NSRouter("/monitor/list", &v1.LiveMonitorController{}, "get:LiveMonitorList"),
			beego.NSRouter("/monitor/:monitor_id", &v1.LiveMonitorController{}, "post:CancelLiveMonitor"),
			beego.NSRouter("/monitor/read/:monitor_id", &v1.LiveMonitorController{}, "post:ReadLiveMonitor"),
			beego.NSRouter("/monitor/rooms/:monitor_id", &v1.LiveMonitorController{}, "get:LiveMonitorRooms"),
			beego.NSRouter("/monitor/:monitor_id", &v1.LiveMonitorController{}, "delete:DeleteLiveMonitor"),
			beego.NSRouter("/monitor/number", &v1.LiveMonitorController{}, "get:LiveMonitorNum"),
		),
		beego.NSNamespace("/living",
			beego.NSRouter("/base/:room_id", &v1dy.LiveController{}, "get:LivingBaseData"),
			beego.NSRouter("/sale/:room_id", &v1dy.LiveController{}, "get:LivingSaleData"),
			beego.NSRouter("/watch/chart/:room_id", &v1dy.LiveController{}, "get:LivingWatchChart"),
			beego.NSRouter("/product/:room_id", &v1dy.LiveController{}, "get:LivingProduct"),
			beego.NSRouter("/message/:room_id", &v1dy.LiveController{}, "get:LivingMessage"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
