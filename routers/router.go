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
	beego.Router("/v1/user/login", &v1.AccountController{}, "put:Login")
	beego.Router("/v1/user/findpwd", &v1.AccountController{}, "put:ResetPwd")
	beego.Router("/v1/sms/code", &v1.CommonController{}, "post:Sms")
	beego.Router("/v1/sms/verify/:grant_type/:username", &v1.CommonController{}, "get:CheckSmsCode")

}
