package es

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/services/elasticsearch"
)

type EsShopBusiness struct {
}

func NewEsShopBusiness() *EsShopBusiness {
	return new(EsShopBusiness)
}

const (
	DY_SHOP_LIVE           = 1
	DY_SHOP_LIVE_BUT_AWEME = 2
	DY_SHOP_AWEME          = 3
	DY_SHOP_AWEME_BUT_LIVE = 4
	DY_SHOP_EQUALS         = 5
)

//小店库查询
func (receiver *EsShopBusiness) BaseSearch(
	keyword, category, secondCategory, thirdCategory string,
	min30Sales, max30Sales int64, min30Gmv, max30Gmv, min30UnitPrice, max30UnitPrice, minScore, maxScore float64,
	isLive, isBrand, isVideo, page, pageSize int,
	sortStr, orderBy string) (list []es.DyShop, total int, comErr global.CommonError) {
	list = []es.DyShop{}
	if sortStr == "" {
		sortStr = "month_sales"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"month_sales", "month_gmv", "score", "month_single_price"}) {
		comErr = global.NewError(4000)
		return
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		comErr = global.NewError(4000)
		return
	}
	if pageSize > 100 {
		comErr = global.NewError(4000)
		return
	}
	esTable := es.DyShopTable
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	if keyword != "" {
		esQuery.SetMatchPhrase("shop_name", keyword)
	}
	if category != "" {
		esQuery.SetMultiMatch([]string{"dcm_level_first", "dcm_level_first_2", "dcm_level_first_3"}, category)
	}
	if secondCategory != "" {
		esQuery.SetMultiMatch([]string{"first_cname", "first_cname_2", "first_cname_3"}, secondCategory)
	}
	if thirdCategory != "" {
		esQuery.SetMultiMatch([]string{"second_cname", "second_cname_2", "second_cname_3"}, thirdCategory)
	}
	if isBrand == 1 {
		esQuery.SetTerm("is_brand", 1)
	}
	if min30Sales > 0 || max30Sales > 0 {
		rangeMap := map[string]interface{}{}
		if min30Sales > 0 {
			rangeMap["gte"] = min30Sales
		}
		if max30Sales > 0 {
			rangeMap["lt"] = max30Sales
		}
		esQuery.SetRange("month_sales", rangeMap)
	}
	if min30Gmv > 0 || max30Gmv > 0 {
		rangeMap := map[string]interface{}{}
		if min30Gmv > 0 {
			rangeMap["gte"] = min30Gmv
		}
		if max30Gmv > 0 {
			rangeMap["lt"] = max30Gmv
		}
		esQuery.SetRange("month_gmv", rangeMap)
	}
	if min30UnitPrice > 0 || max30UnitPrice > 0 {
		rangeMap := map[string]interface{}{}
		if min30UnitPrice > 0 {
			rangeMap["gte"] = min30UnitPrice
		}
		if max30UnitPrice > 0 {
			rangeMap["lt"] = max30UnitPrice
		}
		esQuery.SetRange("month_single_price", rangeMap)
	}
	if minScore > 0 || maxScore > 0 {
		rangeMap := map[string]interface{}{}
		if minScore > 0 {
			rangeMap["gte"] = minScore
		}
		if maxScore > 0 {
			rangeMap["lt"] = maxScore
		}
		esQuery.SetRange("score", rangeMap)
	}
	if isLive == 1 && isVideo != 1 {
		esQuery.SetTerms("commerce_type", []int{DY_SHOP_LIVE, DY_SHOP_LIVE_BUT_AWEME, DY_SHOP_AWEME_BUT_LIVE, DY_SHOP_EQUALS})
	}
	if isLive != 1 && isVideo == 1 {
		esQuery.SetTerms("commerce_type", []int{DY_SHOP_AWEME, DY_SHOP_LIVE_BUT_AWEME, DY_SHOP_AWEME_BUT_LIVE, DY_SHOP_EQUALS})
	}
	results := esMultiQuery.
		SetTable(esTable).
		SetCache(180).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}
