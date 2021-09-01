package routers

import (
	v1 "dongchamao/controllers/api/v1"
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	//抖音达人
	ns := beego.NewNamespace("/v1/dy",
		beego.NSNamespace("/author",
			beego.NSRouter("/red/:type", &v1.CommonController{}, "get:RedAuthorRoom"),
			beego.NSRouter("/top/red", &v1.CommonController{}, "get:RedAuthorLivingRoom"),
			beego.NSRouter("/search", &v1dy.AuthorController{}, "get:BaseSearch"),
			beego.NSRouter("/cate", &v1dy.AuthorController{}, "get:AuthorCate"),
			beego.NSRouter("/live/tags", &v1dy.AuthorController{}, "get:GetCacheAuthorLiveTags"),
			beego.NSRouter("/info/:author_id", &v1dy.AuthorController{}, "get:AuthorBaseData"),
			beego.NSRouter("/view/:author_id", &v1dy.AuthorController{}, "get:AuthorViewData"),
			beego.NSRouter("/reputation/:author_id", &v1dy.AuthorController{}, "get:Reputation"),
			beego.NSRouter("/awemes/:author_id/:start/:end", &v1dy.AuthorController{}, "get:AuthorAwemesByDay"),
			beego.NSRouter("/basic/chart/:author_id/:start/:end", &v1dy.AuthorController{}, "get:AuthorBasicChart"),
			beego.NSRouter("/live/analysis/:author_id/:start/:end", &v1dy.AuthorController{}, "get:CountLiveRoomAnalyse"),
			beego.NSRouter("/live/rooms/:author_id/:start/:end", &v1dy.AuthorController{}, "get:AuthorLiveRooms"),
			beego.NSRouter("/fans/analysis/:author_id", &v1dy.AuthorController{}, "get:AuthorFansAnalyse"),
			beego.NSRouter("/reputation/:author_id", &v1dy.AuthorController{}, "get:AuthorStarSimpleData"),
			beego.NSRouter("/product/analysis/:author_id/:start/:end", &v1dy.AuthorController{}, "get:AuthorProductAnalyse"),
			beego.NSRouter("/product/rooms/:author_id/:product_id/:start/:end", &v1dy.AuthorController{}, "get:AuthorProductRooms"),
			beego.NSRouter("/income", &v1dy.AuthorController{}, "put:AuthorIncome"),
		),
		beego.NSNamespace("/xt/author",
			beego.NSRouter("/index/:author_id", &v1dy.AuthorController{}, "get:AuthorStarSimpleData"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
