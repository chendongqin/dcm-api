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
			beego.NSRouter("/alipay/:channel/:order_id", &v1.PayController{}, "get:AliPay"),
			beego.NSRouter("/notify/wechat", &controllers.CallbackController{}, "*:WechatNotify"),
			beego.NSRouter("/notify/alipay", &controllers.CallbackController{}, "*:AlipayNotify"),
			beego.NSRouter("/order/:order_id", &v1.PayController{}, "get:OrderDetail"),
			beego.NSRouter("/order/:order_id", &v1.PayController{}, "delete:OrderDel"),
			beego.NSRouter("/order/list/:platform", &v1.PayController{}, "get:OrderList"),
			beego.NSRouter("/order/surplus", &v1.PayController{}, "get:DySurplusValue"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)
}
