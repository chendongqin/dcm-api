package routers

import (
	controllers "dongchamao/controllers/api"
	v1 "dongchamao/controllers/api/v1"
	"github.com/astaxie/beego"
)

func init() {
	//支付相关
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/pay",
			beego.NSRouter("/order/dy", &v1.PayController{}, "put:CreateDyOrder"),
			beego.NSRouter("/wechat/:channel/:order_id", &v1.PayController{}, "get:WechatPay"),
			beego.NSRouter("/notify/wechat", &controllers.CallbackController{}, "*:WechatNotify"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
