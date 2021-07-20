package routers

import (
	v1 "dongchamao/controllers/api/v1"
	"github.com/astaxie/beego"
)

func init() {
	//抖音达人
	beego.Router("/v1/account/password", &v1.AccountController{}, "put:ResetPwd")
	beego.Router("/v1/account/info", &v1.AccountController{}, "get:Info")
}
