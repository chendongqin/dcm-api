package routers

import (
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	//抖音达人
	beego.Router("/v1/dy/author/:author_id", &v1dy.AuthorController{}, "get:AuthorBaseData")
	beego.Router("/v1/dy/author/reputation/:author_id", &v1dy.AuthorController{}, "get:Reputation")
	beego.Router("/v1/dy/author/awemes/:author_id/:start/:end", &v1dy.AuthorController{}, "get:AuthorAwemesByDay")
	beego.Router("/v1/dy/author/basic/chart/:author_id/:start/:end", &v1dy.AuthorController{}, "get:AuthorBasicChart")
	beego.Router("/v1/dy/author/fans/analysis/:author_id", &v1dy.AuthorController{}, "get:AuthorFansAnalyse")

	beego.Router("/v1/dy/xt/author/index/:author_id", &v1dy.AuthorController{}, "get:AuthorStarSimpleData")
}
