package routers

import (
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	//视频
	ns := beego.NewNamespace("/v1/dy",
		beego.NSNamespace("/rank",
			beego.NSRouter("/author/aweme", &v1dy.RankController{}, "get:DyStartAuthorVideoRank"),
			beego.NSRouter("/author/live", &v1dy.RankController{}, "get:DyStartAuthorLiveRank"),
			beego.NSRouter("/live/hour/:date/:hour", &v1dy.RankController{}, "get:DyLiveHourRank"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
