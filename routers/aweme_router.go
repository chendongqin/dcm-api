package routers

import (
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	//视频
	beego.Router("/v1/dy/aweme/:aweme_id", &v1dy.AwemeController{}, "get:AwemeBaseData")
	beego.Router("/v1/dy/aweme/chart/:aweme_id/:start/:end", &v1dy.AwemeController{}, "get:AwemeChart")
}
