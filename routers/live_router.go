package routers

import (
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	//直播
	beego.Router("/v1/dy/live/:room_id", &v1dy.LiveController{}, "get:LiveInfoData")

}
