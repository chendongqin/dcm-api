package routers

import (
	v1dy "dongchamao/controllers/api/v1/dy"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/v1/dy/product/base/:product_id", &v1dy.ProductController{}, "get:ProductBase")
	beego.Router("/v1/dy/product/:product_id/:start/:end", &v1dy.ProductController{}, "get:ProductBaseAnalysis")
}
