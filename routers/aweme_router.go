package routers

import (
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	//视频
	ns := beego.NewNamespace("/v1/dy",
		beego.NSNamespace("/aweme",
			beego.NSRouter("/info/:aweme_id", &v1dy.AwemeController{}, "get:AwemeBaseData"),
			beego.NSRouter("/chart/:aweme_id/:start/:end", &v1dy.AwemeController{}, "get:AwemeChart"),
			beego.NSRouter("/chart/:aweme_id/:start/:end", &v1dy.AwemeController{}, "get:AwemeChart"),
			beego.NSRouter("/hot/words/:aweme_id", &v1dy.AwemeController{}, "get:AwemeCommentHotWords"),
			beego.NSRouter("/comments/:aweme_id", &v1dy.AwemeController{}, "get:AwemeCommentTop"),
			beego.NSRouter("/product/:aweme_id/:start/:end", &v1dy.AwemeController{}, "get:AwemeProductAnalyse"),
			beego.NSRouter("/product/chart/:aweme_id/:start/:end", &v1dy.AwemeController{}, "get:AwemeProductAnalyseChart"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
