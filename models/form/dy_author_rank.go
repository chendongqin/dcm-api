package form

import (
	"strconv"
	"time"
)

type dyAuthorGoodsRankParams struct {
	Time		int 	`form:"time" valid:"Required;Range(100000, 99999999)"` // 不能为空并且是6到8位数字 6代表月，7代表周，8代表日
	Tags   		string 	`form:"tags" valid:"Match(//)"` // 分类筛选,空代表全部
	Verified    int    	`form:"verified" valid:"Range(0, 1)"` // 1=只看蓝v
	Page		int    	`form:"page" valid:"Range(1, 1000)"` // 页数
	PageSize	int    	`form:"page_size" valid:"Range(1, 1000)"` // 每页数量
	Sort		string 	`form:"sort" valid:"Required; Match(/^(sum_gmv|sum_sales|avg_price)$/)"` // 排序字段 sum_gmv:销售额, sum_sales:销量, avg_price:客单价
	OrderBy		string 	`form:"order_by" valid:"Match(/^(desc|asc)$/)"` // 排序 desc:降序 asc:升序
}

func DefaultDyAuthorGoodsRankParams() dyAuthorGoodsRankParams {
	timeInt, _ := strconv.Atoi(time.Now().Format("20060102"))
	return dyAuthorGoodsRankParams{
		Time: timeInt,
		Tags: "",
		Verified: 0,
		Page: 1,
		PageSize: 50,
		Sort: "sum_gmv",
		OrderBy: "desc",
	}
}