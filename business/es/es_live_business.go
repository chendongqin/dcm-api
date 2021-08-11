package es

import (
	"dongchamao/global"
	"dongchamao/global/alias"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/entity"
	"dongchamao/models/es"
	"dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	"dongchamao/services/elasticsearch"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"math"
	"time"
)

type EsLiveBusiness struct {
}

func NewEsLiveBusiness() *EsLiveBusiness {
	return new(EsLiveBusiness)
}

//达人直播间搜索
func (receiver *EsLiveBusiness) SearchAuthorRooms(authorId, keyword, sortStr, orderBy string, page, size int, startDate, endDate time.Time) (list []es.EsAuthorLiveRoom, total int, comErr global.CommonError) {
	if sortStr == "" {
		sortStr = "create_timestamp"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"create_timestamp", "gmv", "sales", "max_user_count"}) {
		comErr = global.NewError(4000)
		return
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		comErr = global.NewError(4000)
		return
	}
	if size > 50 {
		comErr = global.NewError(4000)
		return
	}
	//兼容数据 2021-06-29
	esTable := GetESTableByTime(es.DyAuthorLiveRecordsTable, startDate, endDate)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("author_id", authorId)
	esQuery.SetRange("create_timestamp", map[string]interface{}{
		"gte": startDate.Unix(),
		"lt":  endDate.AddDate(0, 0, 1).Unix(),
	})
	if keyword != "" {
		esQuery.AddCondition(map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match_phrase": map[string]interface{}{
							"title": keyword,
						},
					},
					map[string]interface{}{
						"match_phrase": map[string]interface{}{
							"product_title": keyword,
						},
					},
				},
			},
		})
	}
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*size, size).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	for k, v := range list {
		list[k].Sales = math.Floor(v.Sales)
	}
	total = esMultiQuery.Count
	return
}

//直播间商品统计
func (receiver *EsLiveBusiness) CountRoomProductByRoomId(roomInfo entity.DyLiveInfo) int64 {
	date := time.Unix(roomInfo.DiscoverTime, 0).Format("20060102")
	esTable := fmt.Sprintf(es.DyRoomProductRecordsTable, date)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("room_id", roomInfo.RoomID)
	total, _ := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetMultiQuery().
		FindCount()
	return total
}

//直播间商品统计
func (receiver *EsLiveBusiness) CountRoomProductByAuthorId(authorId string, startTime, endTime time.Time) int64 {

	esTable := GetESTableByTime(es.DyRoomProductRecordsTable, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("author_id", authorId)
	esQuery.SetRange("start_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	_ = esMultiQuery.
		SetTable(esTable).
		SetCollapse("product_id.keyword").
		SetFields("product_id").
		AddMust(esQuery.Condition).
		SetMultiQuery().
		Query()
	total := esMultiQuery.Count
	return int64(total)
}

//直播间筛选
func (receiver *EsLiveBusiness) RoomProductByRoomId(roomInfo entity.DyLiveInfo, keyword, sortStr, orderBy, firstLabel, secondLabel, thirdLabel string, page, pageSize int) (list []es.EsAuthorLiveProduct, productCount dy.LiveProductCount, total int, comErr global.CommonError) {
	if sortStr == "" {
		sortStr = "shelf_time"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"shelf_time", "predict_sales", "predict_gmv"}) {
		comErr = global.NewError(4000)
		return
	}
	if pageSize > 50 {
		comErr = global.NewError(4000)
		return
	}
	date := time.Unix(roomInfo.DiscoverTime, 0).Format("20060102")
	esTable := fmt.Sprintf(es.DyRoomProductRecordsTable, date)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("room_id", roomInfo.RoomID)
	if keyword != "" {
		esQuery.SetMatchPhrase("title", keyword)
	}
	if firstLabel != "" {
		if firstLabel == "其他" {
			esQuery.AddCondition(map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{
							"terms": map[string]interface{}{"dcm_level_first.keyword": []string{firstLabel, ""}},
						},
						{
							"bool": map[string]interface{}{
								"must_not": map[string]interface{}{
									"exists": map[string]interface{}{
										"field": "dcm_level_first",
									},
								},
							},
						},
					},
				},
			})
			secondLabel = ""
			thirdLabel = ""
		} else {
			esQuery.SetTerm("dcm_level_first.keyword", firstLabel)
		}
	}
	if secondLabel != "" {
		esQuery.SetTerm("first_cname.keyword", secondLabel)
	}
	if thirdLabel != "" {
		esQuery.SetTerm("second_cname.keyword", thirdLabel)
	}
	orderEs := elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order
	if sortStr == "predict_sales" {
		orderEs = elasticsearch.NewElasticOrder().Add("real_sales").Add(sortStr, orderBy).Order
	} else if sortStr == "predict_gmv" {
		orderEs = elasticsearch.NewElasticOrder().Add("real_gmv").Add(sortStr, orderBy).Order
	}
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit(0, 5000).
		SetOrderBy(orderEs).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	productCount = dy.LiveProductCount{}
	for k, v := range list {
		productCount.ProductNum++
		////todo gmv处理
		//if v.RealGmv > 0 {
		//	var sale float64 = 0
		//	if v.Price > 0 {
		//		sale = math.Floor(v.RealGmv / v.Price)
		//	}
		//	productCount.Sales += sale
		//	productCount.Gmv += v.RealGmv
		//	list[k].PredictGmv = v.RealGmv
		//	list[k].PredictSales = sale
		//} else {
		productCount.Sales += math.Floor(v.PredictSales)
		productCount.Gmv += v.PredictGmv
		//}
		if v.IsReturn == 1 && v.StartTime == v.ShelfTime {
			list[k].IsReturn = 0
		}
	}
	start := (page - 1) * pageSize
	end := start + pageSize
	listLen := len(list)
	if listLen < end {
		end = listLen
	}
	list = list[start:end]
	for k, v := range list {
		list[k].Cover = dyimg.Product(v.Cover)
		////todo 真实gmv存在，按gmv处理
		//if v.RealGmv > 0 {
		//	list[k].PredictGmv = v.RealGmv
		//	if v.Price > 0 {
		//		list[k].PredictSales = math.Floor(v.RealGmv / v.Price)
		//	}
		//} else {
		list[k].PredictSales = math.Floor(v.PredictSales)
		//}
		if v.Pv > 0 {
			list[k].BuyRate = v.PredictSales / float64(v.Pv)
		}
	}
	total = esMultiQuery.Count
	return
}

//直播间商品分类统计
func (receiver *EsLiveBusiness) AllRoomProductCateByRoomId(roomInfo entity.DyLiveInfo) (productCount dy.LiveProductCateCount) {
	cKey := cache.GetCacheKey(cache.LiveRoomProductCount, roomInfo.RoomID)
	productCountJson := global.Cache.Get(cKey)
	if productCountJson != "" {
		productCountJson = utils.DeserializeData(productCountJson)
		_ = jsoniter.Unmarshal([]byte(productCountJson), &productCount)
		return
	}
	date := time.Unix(roomInfo.DiscoverTime, 0).Format("20060102")
	esTable := fmt.Sprintf(es.DyRoomProductRecordsTable, date)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("room_id", roomInfo.RoomID)
	list := make([]es.EsAuthorLiveProduct, 0)
	results := esMultiQuery.
		SetFields("dcm_level_first", "first_cname", "second_cname").
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit(0, 5000).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("start_time", "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	firstCateCountMap := map[string]int{}
	firstCateMap := map[string]map[string]bool{}
	secondCateMap := map[string]map[string]bool{}
	for _, v := range list {
		if v.DcmLevelFirst == "" {
			v.DcmLevelFirst = "其他"
		}
		if _, ok := firstCateMap[v.DcmLevelFirst]; !ok {
			firstCateMap[v.DcmLevelFirst] = map[string]bool{}
		}
		if _, ok := firstCateCountMap[v.DcmLevelFirst]; !ok {
			firstCateCountMap[v.DcmLevelFirst] = 1
		} else {
			firstCateCountMap[v.DcmLevelFirst] += 1
		}
		if v.FirstCname == "" || v.DcmLevelFirst == "其他" {
			continue
		}
		firstCateMap[v.DcmLevelFirst][v.FirstCname] = true
		if _, ok := secondCateMap[v.FirstCname]; !ok {
			secondCateMap[v.FirstCname] = map[string]bool{}
		}
		if v.SecondCname == "" {
			continue
		}
		secondCateMap[v.FirstCname][v.SecondCname] = true
	}
	productCount.CateList = []dy.DyCate{}
	for k, v := range firstCateMap {
		secondCateList := make([]dy.DyCate, 0)
		for ck, _ := range v {
			if cv, ok := secondCateMap[ck]; ok {
				secondCateItem := dy.DyCate{
					Name: ck,
				}
				for tk, _ := range cv {
					secondCateItem.SonCate = append(secondCateItem.SonCate, dy.DyCate{
						Name:    tk,
						Num:     0,
						SonCate: nil,
					})
				}
				if len(secondCateItem.SonCate) == 0 {
					secondCateItem.SonCate = []dy.DyCate{}
				}
				secondCateList = append(secondCateList, secondCateItem)
			}
		}
		productNumber := 0
		if n, ok := firstCateCountMap[k]; ok {
			productNumber = n
		}
		item := dy.DyCate{
			Name:    k,
			Num:     productNumber,
			SonCate: []dy.DyCate{},
		}
		if len(secondCateList) > 0 {
			item.SonCate = secondCateList
		}
		productCount.CateList = append(productCount.CateList, item)
	}
	var timeout time.Duration = 60
	if roomInfo.FinishTime <= (time.Now().Unix() - 3600) {
		timeout = 1800
	}
	cateListJson := utils.SerializeData(productCount)
	_ = global.Cache.Set(cKey, cateListJson, timeout)
	return
}

//达人直播带货商品直播列表
func (receiver *EsLiveBusiness) GetAuthorProductSearchRoomIds(authorId, productId string, startTime, stopTime time.Time, page, pageSize int) (roomIds []string, total int, comErr global.CommonError) {
	esTable := GetESTableByTime(es.DyRoomProductRecordsTable, startTime, stopTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("author_id", authorId)
	esQuery.SetTerm("product_id", productId)
	esQuery.SetRange("shelf_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  stopTime.AddDate(0, 0, 1).Unix(),
	})
	results := esMultiQuery.
		SetTable(esTable).
		SetFields("room_id").
		AddMust(esQuery.Condition).
		SetMultiQuery().
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("shelf_time", "desc").Order).
		Query()
	list := make([]es.EsAuthorLiveProduct, 0)
	utils.MapToStruct(results, &list)
	for _, v := range list {
		roomIds = append(roomIds, v.RoomID)
	}
	total = esMultiQuery.Count
	return
}

//达人直播间搜索
func (receiver *EsLiveBusiness) SearchProductRooms(productId, keyword, sortStr, orderBy string, page, size int, startTime, endTime time.Time) (list []es.EsAuthorLiveProduct, total int, comErr global.CommonError) {
	if sortStr == "" {
		sortStr = "shelf_time"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"shelf_time", "predict_gmv", "predict_sales"}) {
		comErr = global.NewError(4000)
		return
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		comErr = global.NewError(4000)
		return
	}
	if size > 50 {
		comErr = global.NewError(4000)
		return
	}
	esTable := GetESTableByTime(es.DyRoomProductRecordsTable, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("product_id", productId)
	esQuery.SetRange("live_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if keyword != "" {
		esQuery.AddCondition(map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match_phrase": map[string]interface{}{
							"room_title": keyword,
						},
					},
					map[string]interface{}{
						"match_phrase": map[string]interface{}{
							"nickname": keyword,
						},
					},
				},
			},
		})
	}
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*size, size).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	for k, v := range list {
		list[k].PredictSales = math.Floor(v.PredictSales)
		list[k].Cover = dyimg.Fix(v.Cover)
		list[k].RoomCover = dyimg.Fix(v.RoomCover)
		////todo gmv处理
		//if v.RealGmv > 0 {
		//	var sale float64 = 0
		//	if v.Price > 0 {
		//		sale = math.Floor(v.RealGmv / v.Price)
		//	}
		//	list[k].PredictGmv = v.RealGmv
		//	list[k].PredictSales = sale
		//}
		if v.IsReturn == 1 && v.StartTime == v.ShelfTime {
			list[k].IsReturn = 0
		}
		if v.Pv > 0 {
			list[k].BuyRate = v.PredictSales / float64(v.Pv)
		}
	}
	total = esMultiQuery.Count
	return
}

func (receiver *EsLiveBusiness) SearchLiveRooms(keyword, category, firstName, secondName, thirdName string, minAmount, maxAmount, minAvgUserCount, maxAvgUserCount int64, minUv, maxUv, hasProduct, brand, keywordType int, sortStr, orderBy string, page, size int, startTime, endTime time.Time) (list []es.EsDyLiveInfo, total int, comErr global.CommonError) {
	if sortStr == "" {
		sortStr = "avg_user_count"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"create_time", "predict_gmv", "predict_uv_value", "avg_user_count"}) {
		comErr = global.NewError(4000)
		return
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		comErr = global.NewError(4000)
		return
	}
	if size > 50 {
		comErr = global.NewError(4000)
		return
	}
	esTable := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if minAmount > 0 || maxAmount > 0 {
		rangeMap := map[string]interface{}{}
		if minAmount > 0 {
			rangeMap["gte"] = minAmount
		}
		if maxAmount > 0 {
			rangeMap["lt"] = maxAmount
		}
		esQuery.SetRange("predict_gmv", rangeMap)
	}
	if minUv > 0 || maxUv > 0 {
		rangeMap := map[string]interface{}{}
		if minUv > 0 {
			rangeMap["gte"] = minUv
		}
		if maxUv > 0 {
			rangeMap["lt"] = maxUv
		}
		esQuery.SetRange("predict_uv_value", rangeMap)
	}
	if minAvgUserCount > 0 || maxAvgUserCount > 0 {
		rangeMap := map[string]interface{}{}
		if minAvgUserCount > 0 {
			rangeMap["gte"] = minAvgUserCount
		}
		if maxAvgUserCount > 0 {
			rangeMap["lt"] = maxAvgUserCount
		}
		esQuery.SetRange("avg_user_count", rangeMap)
	}
	if keyword != "" {
		if keywordType == 1 {
			esQuery.SetMatchPhrase("product_title", keyword)
		} else {
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
			} else {
				esQuery.SetMultiMatch([]string{"display_id", "short_id", "nickname"}, keyword)
			}
		}
	}
	if brand == 1 {
		esQuery.SetTerm("brand", 1)
	}
	if hasProduct == 1 {
		esQuery.SetRange("num_promotions", map[string]interface{}{
			"gt": 0,
		})
	}
	if category != "" {
		esQuery.SetMatchPhrase("tags", category)
	}
	if firstName != "" {
		esQuery.SetMatchPhrase("dcm_level_first", firstName)
	}
	if secondName != "" {
		esQuery.SetMatchPhrase("first_cname", secondName)
	}
	if thirdName != "" {
		esQuery.SetMatchPhrase("second_cname", thirdName)
	}
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*size, size).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	for k, v := range list {
		list[k].Cover = dyimg.Fix(v.Cover)
		list[k].Avatar = dyimg.Fix(v.Avatar)
		////todo gmv处理
		//if v.RealGmv > 0 {
		//	list[k].PredictGmv = v.RealGmv
		//}
		//if v.RealUvValue > 0 {
		//	list[k].PredictUvValue = v.RealUvValue
		//}
		list[k].AvgUserCount = math.Floor(v.AvgUserCount)
		if v.DisplayId == "" {
			list[k].DisplayId = v.ShortId
		}
		list[k].TagsArr = v.GetTagsArr()
	}
	total = esMultiQuery.Count
	return
}
