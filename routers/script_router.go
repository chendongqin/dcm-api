package routers

import (
	controllers "dongchamao/controllers/api"
	"github.com/astaxie/beego"
)

func init() {
	//视频
	ns := beego.NewNamespace("/v1/script",
		beego.NSRouter("/author/tag", &controllers.ScriptController{}, "get:AuthorTag"),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
