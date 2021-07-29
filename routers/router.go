package routers

import (
	"dongchamao/controllers/api/v1"
	"github.com/astaxie/beego"
	"github.com/json-iterator/go/extra"
)

func init() {
	// 容忍字符串和数字互转
	extra.RegisterFuzzyDecoders()
	//beego.Get("/v1/ipip", func(ctx *context.Context) {
	//	ip := ctx.Input.IP()
	//	header := fmt.Sprintf("%+v", ctx.Request.Header)
	//	sp := ctx.Input.Header("Server-Protocol")
	//	sh := ctx.Input.Host()
	//	ctx.Output.Body([]byte(ip + sp + sh + header))
	//})

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSRouter("/login", &v1.AccountController{}, "put:Login"),
			beego.NSRouter("/findpwd", &v1.AccountController{}, "put:FindPwd"),
		),
		beego.NSNamespace("/sms",
			beego.NSRouter("/code", &v1.CommonController{}, "post:Sms"),
			beego.NSRouter("/verify/:grant_type/:username/:code", &v1.CommonController{}, "get:CheckSmsCode"),
		),
	)
	// 注册路由组
	beego.AddNamespace(ns)

	beego.Router("/v1/test", &v1.CommonController{}, "get:Test")

}
