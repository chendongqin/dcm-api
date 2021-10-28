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
	esTable, connection := GetESTable(es.DyProductTable)
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
	} else {
		if keyword == "" {
			esQuery.SetTerm("platform_label.keyword", "小店")
		}
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
	} else if commerceType == 4 {
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
	sortOrder := elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order
	//if utils.InArrayString(sortStr, []string{"order_account", "pv", "cvr"}) {
	//	sortOrder = elasticsearch.NewElasticOrder().Add("is_yesterday", "desc").Add(sortStr, orderBy).Order
	//}
	var cacheTime time.Duration = 120
	var outTime = 10 * time.Second
	esMultiQuery.Timeout = &outTime
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(sortOrder).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

//拼接商品库字段日期后缀
func FixDate(fields string, dateType int) string {
	if !utils.InArrayString(fields, []string{"pv", "cvr", "order_account", "gpm", "is_coupon", "platform_label", "relate_aweme", "relate_room", "relate_author"}) {
		return fields
	}
	var dateTypeMap = map[int]string{0: "", 1: "_7", 2: "_15", 3: "_30"}
	return fields + dateTypeMap[dateType]
}

func (i *EsProductBusiness) BaseSearchNew(productId, keyword, category, secondCategory, thirdCategory, platform string,
	minCommissionRate, minPrice, maxPrice, minGpm, maxGpm float64, commerceType, isCoupon, relateRoom, relateAweme, relateAuthor, isStar, notStar, page, pageSize, dateType int,
	sortStr, orderBy string) (list []es.ProductNew, total int, comErr global.CommonError) {
	list = []es.ProductNew{}
	//"pv","cvr","order_account","gpm","is_coupon","platform_label","relate_aweme","relate_room","relate_author","is_star"
	if sortStr == "" {
		sortStr = "order_account"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"commission", "pv", "order_account", "cvr", "gpm", "relate_aweme", "relate_author", "relate_room"}) {
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
	esTable, connection := GetESTable(es.DyProductTableNew)
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
	} else {
		if keyword == "" {
			esQuery.SetTerm("platform_label.keyword", "小店")
		}
	}
	if isStar == 1 {
		esQuery.SetTerm(FixDate("is_star", dateType), 1)
	}
	if notStar == 1 {
		esQuery.SetTerm(FixDate("is_star", dateType), 0)
	}
	if isCoupon == 1 {
		esQuery.SetTerm(FixDate("is_coupon", dateType), 1)
	}
	if relateRoom == 1 {
		esQuery.SetRange(FixDate("relate_room", dateType), map[string]interface{}{"gt": 0})
	}
	if relateAweme == 1 {
		esQuery.SetRange(FixDate("relate_aweme", dateType), map[string]interface{}{"gt": 0})
	}
	if relateAuthor == 1 {
		esQuery.SetRange(FixDate("relate_author", dateType), map[string]interface{}{"gt": 0})
	}
	var commerceTypeMap = map[int][]int{
		1: {1, 2},
		2: {1, 2, 4, 5},
		3: {3, 4},
		4: {2, 3, 4, 5},
	}
	if commerce, exist := commerceTypeMap[commerceType]; exist {
		esQuery.SetTerms(FixDate("commerce_type", dateType), commerce)
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
	if minGpm > 0 || maxGpm > 0 {
		rangeMap := map[string]interface{}{}
		if minGpm > 0 {
			rangeMap["gte"] = minGpm
		}
		if maxGpm > 0 {
			rangeMap["lte"] = maxGpm
		}
		esQuery.SetRange(FixDate("gpm", dateType), rangeMap)
	}
	sortOrder := elasticsearch.NewElasticOrder().Add(FixDate(sortStr, dateType), orderBy).Order
	var cacheTime time.Duration = 120
	var outTime = 10 * time.Second
	esMultiQuery.Timeout = &outTime
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(sortOrder).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

func (i *EsProductBusiness) SearchRangeDateList(productId, shopId, authorId string, startTime, endTime time.Time, page, pageSize int) (list []es.DyProductAuthorAnalysis, total int, comErr global.CommonError) {
	esTable, connection, err := GetESTableByTime(es.DyProductAuthorAnalysisTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	if productId != "" {
		esQuery.SetTerm("productId", productId)
	}
	esQuery.SetRange("createSdf.keyword", map[string]interface{}{
		"gte": startTime.Format("20060102"),
		"lte": endTime.Format("20060102"),
	})
	if authorId != "" {
		esQuery.SetTerm("authorId", authorId)
	}
	if shopId != "" {
		esQuery.SetTerm("shopId", shopId)
	}
	results := esMultiQuery.
		SetConnection(connection).
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

func (i *EsProductBusiness) SearchAwemeRangeDateList(productId, shopId, authorId string, startTime, endTime time.Time, page, pageSize int) (list []es.DyProductAwemeAuthorAnalysis, total int, comErr global.CommonError) {
	esTable, connection, err := GetESTableByMonthTime(es.DyProductAwemeAuthorAnalysisTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	if productId != "" {
		esQuery.SetTerm("productId", productId)
	}
	if shopId != "" {
		esQuery.SetTerm("shopId", shopId)
	}
	esQuery.SetRange("createSdf.keyword", map[string]interface{}{
		"gte": startTime.Format("20060102"),
		"lte": endTime.Format("20060102"),
	})
	if authorId != "" {
		esQuery.SetTerm("authorId", authorId)
	}
	results := esMultiQuery.
		SetConnection(connection).
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
	esTable, connection, err := GetESTableByTime(es.DyProductAuthorAnalysisTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
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
			esQuery.SetMatchPhraseWithParams("nickname", keyword, alias.M{
				"slop": slop,
			})
		} else {
			esQuery.SetMultiMatch([]string{"nickname", "displayId", "shortId"}, keyword)
		}
	}
	result := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetFields("productId", "authorId", "createSdf").
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("_id", "asc").Order).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result, &startRow)
	result2 := esMultiQuery2.
		SetConnection(connection).
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

func (i *EsProductBusiness) SearchAwemeRangeDateRowKey(productId, keyword string, startTime, endTime time.Time) (startRow es.DyProductAuthorAnalysis, stopRow es.DyProductAuthorAnalysis, total int, comErr global.CommonError) {
	esTable, connection, err := GetESTableByMonthTime(es.DyProductAwemeAuthorAnalysisTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
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
			esQuery.SetMatchPhraseWithParams("nickname", keyword, alias.M{
				"slop": slop,
			})
		} else {
			esQuery.SetMultiMatch([]string{"nickname", "displayId", "shortId"}, keyword)
		}
	}
	result := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetFields("productId", "authorId", "createSdf").
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("_id", "asc").Order).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result, &startRow)
	result2 := esMultiQuery2.
		SetConnection(connection).
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

func (i *EsProductBusiness) SimpleSearch(productId, title, platformLabel, dcmLevelFirst, firstCname, secondCname, thirdCname string, minPrice, maxPrice float64, page, pageSize int) (list []es.DyProduct, total int, comErr global.CommonError) {
	esTable, connection := GetESTable(es.DyProductTable)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	if productId != "" {
		esQuery.SetTerm("product_id", productId)
	}
	if title != "" {
		esQuery.SetMatchPhrase("title", title)
	}
	if platformLabel != "" {
		esQuery.SetTerm("platform_label.keyword", platformLabel)
	}
	if dcmLevelFirst != "" {
		esQuery.SetTerm("dcm_level_first.keyword", dcmLevelFirst)
	}
	if firstCname != "" {
		esQuery.SetTerm("first_cname.keyword", firstCname)
	}
	if secondCname != "" {
		esQuery.SetTerm("second_cname.keyword", secondCname)
	}
	if thirdCname != "" {
		esQuery.SetTerm("third_cname.keyword", thirdCname)
	}
	if minPrice > 0 || maxPrice > 0 {
		if minPrice == maxPrice {
			esQuery.SetTerm("price", minPrice)
		} else {
			rangeMap := map[string]interface{}{}
			if minPrice > 0 {
				rangeMap["gte"] = minPrice
			}
			if maxPrice > 0 {
				rangeMap["lt"] = maxPrice
			}
			esQuery.SetRange("price", rangeMap)
		}
	}
	results := esMultiQuery.
		SetConnection(connection).
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

//后台内部通过productIds获取商品信息
func (i *EsProductBusiness) SimpleSearchByIds(productIds []string, page, pageSize int) (list []es.DyProduct, total int, comErr global.CommonError) {
	if pageSize > 30 {
		comErr = global.NewError(4000)
		return
	}
	esTable, connection := GetESTable(es.DyProductTable)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerms("product_id", productIds)
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

func (i *EsProductBusiness) KeywordSearch(keyword string) (list []es.DyProduct) {
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esTable, connection := GetESTable(es.DyProductTable)
	if keyword != "" {
		esQuery.SetMatchPhrase("title", keyword)
	} else {
		esTable, connection = GetESTable(es.DyProductTableNew)
	}
	var cacheTime time.Duration = 60
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddShould(esQuery.Condition).
		SetLimit(0, 2).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("order_account", "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	return
}

func (i *EsProductBusiness) ProductSalesTopDayRank(day, fCate, sCate, tCate, sortStr, orderBy string,
	page, pageSize int) (list []es.DyProductSalesTopRank, total int, commonError global.CommonError) {
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esTable, connection := GetESTableByDate(es.DyProductSalesTopTable, day)
	if sortStr == "" {
		sortStr = "order_count"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"cos_fee", "order_count", "order_account_count"}) {
		commonError = global.NewError(4000)
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		commonError = global.NewError(4000)
	}
	if fCate != "" {
		esQuery.SetTerm("dcm_cname.keyword", fCate)
	}
	if sCate != "" {
		esQuery.SetTerm("first_cname.keyword", sCate)
	}
	if tCate != "" {
		esQuery.SetTerm("second_cname.keyword", tCate)
	}
	esOrder := elasticsearch.NewElasticOrder().Add(sortStr, orderBy)
	if sortStr != "order_count" {
		esOrder.Add("order_count", "desc")
	}
	var cacheTime time.Duration = 600
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddShould(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(esOrder.Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

func (i *EsProductBusiness) ProductShareTopDayRank(day, fCate, sCate, tCate, sortStr, orderBy string,
	page, pageSize int) (list []es.DyProductShareTopRank, total int, commonError global.CommonError) {
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esTable, connection := GetESTableByDate(es.DyProductShareTopTable, day)
	list = []es.DyProductShareTopRank{}
	if sortStr == "" {
		sortStr = "share_count"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"cos_fee", "share_count"}) {
		commonError = global.NewError(4000)
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		commonError = global.NewError(4000)
	}
	if fCate != "" {
		esQuery.SetTerm("dcm_cname.keyword", fCate)
	}
	if sCate != "" {
		esQuery.SetTerm("first_cname.keyword", sCate)
	}
	if tCate != "" {
		esQuery.SetTerm("second_cname.keyword", tCate)
	}
	esOrder := elasticsearch.NewElasticOrder().Add(sortStr, orderBy)
	if sortStr != "share_count" {
		esOrder.Add("share_count", "desc")
	}
	var cacheTime time.Duration = 300
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddShould(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(esOrder.Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

func (i *EsProductBusiness) LiveProductSalesTopDayRank(day, fCate, sCate, tCate, sortStr, orderBy string,
	page, pageSize int) (list []es.DyLiveProductSaleTopRank, total int, commonError global.CommonError) {
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esTable, connection := GetESTableByDate(es.DyLiveProductSalesTopTable, day)
	if sortStr == "" {
		sortStr = "gmv"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"cos_fee", "sales", "gmv", "live_count", "price"}) {
		commonError = global.NewError(4000)
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		commonError = global.NewError(4000)
	}
	if fCate != "" {
		esQuery.SetTerm("dcm_cname.keyword", fCate)
	}
	if sCate != "" {
		esQuery.SetTerm("first_cname.keyword", sCate)
	}
	if tCate != "" {
		esQuery.SetTerm("second_cname.keyword", tCate)
	}
	esOrder := elasticsearch.NewElasticOrder().Add(sortStr, orderBy)
	if sortStr != "gmv" {
		esOrder.Add("gmv", "desc")
	}
	var cacheTime time.Duration = 600
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddShould(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(esOrder.Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

func (i *EsProductBusiness) SearchProducts(productIds []string) (list []es.DyProduct, total int, comErr global.CommonError) {
	list = []es.DyProduct{}
	if len(productIds) == 0 {
		return
	}
	esTable, connection := GetESTable(es.DyProductTable)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerms("product_id", productIds)
	var cacheTime time.Duration = 300
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetCache(cacheTime).
		SetLimit(0, len(productIds)).
		SetOrderBy(elasticsearch.NewElasticOrder().Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

//获取查询rowkey的productid
func (i *EsProductBusiness) GetSearchRowKey(keyword, category string) (starRowKey string, stopRowKey string) {
	esTable, connection := GetESTable(es.DyProductTable)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	if keyword != "" {
		esQuery.SetMatchPhrase("title", keyword)
	}
	if category != "" {
		esQuery.SetTerm("dcm_level_first.keyword", category)
	}
	sortOrder1 := elasticsearch.NewElasticOrder().Add("_id", "desc").Order
	result1 := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetFields("product_id").
		AddMust(esQuery.Condition).
		SetLimit(0, 1).
		SetOrderBy(sortOrder1).
		SetMultiQuery().
		QueryOne()
	stopRow := es.DyProduct{}
	utils.MapToStruct(result1, &stopRow)
	if esMultiQuery.Count == 0 {
		return
	} else if esMultiQuery.Count == 1 {
		stopRowKey = stopRow.ProductId
		starRowKey = stopRowKey
		return
	}
	_, esMultiQuery2 := elasticsearch.NewElasticQueryGroup()
	sortOrder2 := elasticsearch.NewElasticOrder().Add("_id", "asc").Order
	result2 := esMultiQuery2.
		SetConnection(connection).
		SetTable(esTable).
		SetFields("product_id").
		AddMust(esQuery.Condition).
		SetLimit(0, 1).
		SetOrderBy(sortOrder2).
		SetMultiQuery().
		QueryOne()
	startRow := es.DyProduct{}
	utils.MapToStruct(result2, &startRow)
	stopRowKey = stopRow.ProductId
	starRowKey = startRow.ProductId
	return
}
