package es

import (
	"dongchamao/global"
	"dongchamao/global/alias"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/services/elasticsearch"
	"time"
)

type EsProductBusiness struct {
}

func NewEsProductBusiness() *EsProductBusiness {
	return new(EsProductBusiness)
}

func (i *EsProductBusiness) BaseSearch(productId, keyword, category, secondCategory, thirdCategory, platform string,
	minCommissionRate, minPrice, maxPrice float64, commerceType, isCoupon, relateRoom, relateAweme, isStar, notStar, page, pageSize int,
	sortStr, orderBy string) (list []es.DyProduct, total int, comErr global.CommonError) {
	list = []es.DyProduct{}
	if sortStr == "" {
		sortStr = "order_account"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"order_account", "pv", "cvr", "week_order_account", "month_order_account", "commission"}) {
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
	esTable := es.DyProductTable
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	if keyword != "" {
		esQuery.SetMatchPhrase("title", keyword)
	}
	if productId != "" {
		esQuery.SetTerm("product_id", productId)
	}
	if category != "" {
		esQuery.SetTerm("dcm_level_first.keyword", category)
	}
	if secondCategory != "" {
		esQuery.SetTerm("first_cname.keyword", secondCategory)
	}
	if thirdCategory != "" {
		esQuery.SetTerm("second_cname.keyword", thirdCategory)
	}
	if platform != "" {
		esQuery.SetTerm("platform_label.keyword", platform)
	}
	if isStar == 1 {
		esQuery.SetTerm("is_star", 1)
	}
	if notStar == 1 {
		esQuery.SetTerm("is_star", 0)
	}
	if isCoupon == 1 {
		esQuery.SetTerm("is_coupon", 1)
	}
	if relateRoom == 1 {
		esQuery.SetRange("relate_room", map[string]interface{}{"gt": 0})
	}
	if relateAweme == 1 {
		esQuery.SetRange("relate_aweme", map[string]interface{}{"gt": 0})
	}
	if commerceType == 1 {
		esQuery.SetTerms("commerce_type", []int{1, 2})
	} else if commerceType == 2 {
		esQuery.SetTerms("commerce_type", []int{1, 2, 4, 5})
	} else if commerceType == 3 {
		esQuery.SetTerms("commerce_type", []int{3, 4})
	} else if commerceType == 3 {
		esQuery.SetTerms("commerce_type", []int{2, 3, 4, 5})
	}
	if minCommissionRate > 0 {
		esQuery.SetRange("commission_rate", map[string]interface{}{
			"gte": minCommissionRate,
		})
	}
	if minPrice > 0 || maxPrice > 0 {
		rangeMap := map[string]interface{}{}
		if minPrice > 0 {
			rangeMap["gte"] = minPrice
		}
		if maxPrice > 0 {
			rangeMap["lte"] = maxPrice
		}
		esQuery.SetRange("coupon_price", rangeMap)
	}
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

func (i *EsProductBusiness) SearchRangeDateList(productId, authorId string, startTime, endTime time.Time, page, pageSize int) (list []es.DyProductAuthorAnalysis, total int, comErr global.CommonError) {
	esTable := GetESTableByTime(es.DyProductAuthorAnalysisTable, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("productId", productId)
	esQuery.SetRange("createSdf.keyword", map[string]interface{}{
		"gte": startTime.Format("20060102"),
		"lte": endTime.Format("20060102"),
	})
	if authorId != "" {
		esQuery.SetTerm("authorId", authorId)
	}
	results := esMultiQuery.
		SetTable(esTable).
		SetFields("productId", "authorId", "createSdf").
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("_id", "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

func (i *EsProductBusiness) SearchRangeDateRowKey(productId, keyword string, startTime, endTime time.Time) (startRow es.DyProductAuthorAnalysis, stopRow es.DyProductAuthorAnalysis, total int, comErr global.CommonError) {
	esTable := GetESTableByTime(es.DyProductAuthorAnalysisTable, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	_, esMultiQuery2 := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("productId", productId)
	esQuery.SetRange("createSdf.keyword", map[string]interface{}{
		"gte": startTime.Format("20060102"),
		"lte": endTime.Format("20060102"),
	})
	if keyword != "" {
		if utils.HasChinese(keyword) {
			slop := 100
			length := len([]rune(keyword))
			if length <= 3 {
				slop = 2
			}
			esMultiQuery.AddMust(elasticsearch.Query().
				SetMatchPhraseWithParams("nickname", keyword, alias.M{
					"slop": slop,
				}).Condition)
			esMultiQuery2.AddMust(elasticsearch.Query().
				SetMatchPhraseWithParams("nickname", keyword, alias.M{
					"slop": slop,
				}).Condition)
		} else {
			esQuery.SetMultiMatch([]string{"displayId", "shortId"}, keyword)
		}
	}
	result := esMultiQuery.
		SetTable(esTable).
		SetFields("productId", "authorId", "createSdf").
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("_id", "asc").Order).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result, &startRow)
	result2 := esMultiQuery2.
		SetTable(esTable).
		SetFields("productId", "authorId", "createSdf").
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("_id", "desc").Order).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result2, &stopRow)
	total = esMultiQuery.Count
	return
}

func (i *EsProductBusiness) InternalSearch(productId, title, dcmLevelFirst, firstCname, secondCname, thirdCname string, page, pageSize int) (list []es.DyProduct, total int, comErr global.CommonError) {
	esTable := es.DyProductTable
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	if productId != "" {
		esQuery.SetTerm("product_id", productId)
	}
	if title != "" {
		esQuery.SetMatchPhrase("title", title)
	}
	if dcmLevelFirst != "" {
		esQuery.SetMatchPhrase("dcm_level_first", dcmLevelFirst)
	}
	if firstCname != "" {
		esQuery.SetMatchPhrase("first_cname", firstCname)
	}
	if secondCname != "" {
		esQuery.SetMatchPhrase("second_cname", secondCname)
	}
	if thirdCname != "" {
		esQuery.SetMatchPhrase("third_cname", thirdCname)
	}
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("order_account", "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}
