package es

import (
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/services/dyimg"
	"dongchamao/services/elasticsearch"
	"dongchamao/structinit/repost/dy"
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
	esTable := GetESTableByTime(es.DyAuthorLiveRecords, startDate, endDate)
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
	esTable := fmt.Sprintf(es.DyRoomProductRecords, date)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("room_id", roomInfo.RoomID)
	total, _ := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetMultiQuery().FindCount()
	return total
}

//直播间商品统计
func (receiver *EsLiveBusiness) CountRoomProductByAuthorId(authorId string, startTime, endTime time.Time) int64 {

	esTable := GetESTableByTime(es.DyRoomProductRecords, startTime, endTime)
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
		SetMultiQuery().Query()
	total := esMultiQuery.Count
	return int64(total)
}

//直播间筛选
func (receiver *EsLiveBusiness) RoomProductByRoomId(roomInfo entity.DyLiveInfo, keyword, sortStr, orderBy, firstLabel, secondLabel, thirdLabel string, page, pageSize int) (list []es.EsAuthorLiveProduct, productCount dy.LiveProductCount, total int, comErr global.CommonError) {
	if sortStr == "" {
		sortStr = "start_time"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"start_time", "predict_sales", "predict_gmv"}) {
		comErr = global.NewError(4000)
		return
	}
	if pageSize > 50 {
		comErr = global.NewError(4000)
		return
	}
	date := time.Unix(roomInfo.DiscoverTime, 0).Format("20060102")
	esTable := fmt.Sprintf(es.DyRoomProductRecords, date)
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
							"match_phrase": map[string]interface{}{"dcm_level_first": firstLabel},
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
			esQuery.SetMatchPhrase("dcm_level_first", firstLabel)
		}
	}
	if secondLabel != "" {
		esQuery.SetMatchPhrase("first_cname", secondLabel)
	}
	if thirdLabel != "" {
		esQuery.SetMatchPhrase("second_cname", thirdLabel)
	}
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit(0, 1000).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	productCount = dy.LiveProductCount{}
	for k, v := range list {
		productCount.ProductNum++
		//todo gmv处理
		if v.RealGmv > 0 {
			var sale float64 = 0
			if v.Price > 0 {
				sale = math.Floor(v.RealGmv / v.Price)
			}
			productCount.Sales += sale
			productCount.Gmv += v.RealGmv
			list[k].PredictGmv = v.RealGmv
			list[k].PredictSales = sale
		} else {
			productCount.Sales += math.Floor(v.PredictSales)
			productCount.Gmv += v.PredictGmv
		}
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
		//todo 真实gmv存在，按gmv处理
		if v.RealGmv > 0 {
			list[k].PredictGmv = v.RealGmv
			if v.Price > 0 {
				list[k].PredictSales = math.Floor(v.RealGmv / v.Price)
			}
		} else {
			list[k].PredictSales = math.Floor(v.PredictSales)
		}
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
		jsoniter.Unmarshal([]byte(productCountJson), &productCount)
		return
	}
	date := time.Unix(roomInfo.DiscoverTime, 0).Format("20060102")
	esTable := fmt.Sprintf(es.DyRoomProductRecords, date)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("room_id", roomInfo.RoomID)
	list := make([]es.EsAuthorLiveProduct, 0)
	results := esMultiQuery.
		SetFields("dcm_level_first", "first_cname", "second_cname").
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit(0, 1000).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("start_time", "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	firstCateCountMap := map[string]int{}
	firstCateMap := map[string]map[string]bool{}
	secondCateMap := map[string]map[string]bool{}
	//gmv写入数据数
	gmvNum := 0
	productNum := 0
	for _, v := range list {
		if v.RealGmv > 0 {
			gmvNum++
		}
		productNum++
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
	productCount.CateList = []dy.LiveProductFirstCate{}
	for k, v := range firstCateMap {
		secondCateList := make([]dy.LiveProductSecondCate, 0)
		for ck, _ := range v {
			if cv, ok := secondCateMap[ck]; ok {
				secondCate := dy.LiveProductSecondCate{
					Name: ck,
				}
				for tk, _ := range cv {
					secondCate.Cate = append(secondCate.Cate, tk)
				}
				if len(secondCate.Cate) == 0 {
					secondCate.Cate = []string{}
				}
				secondCateList = append(secondCateList, secondCate)
			}
		}
		productNum := 0
		if n, ok := firstCateCountMap[k]; ok {
			productNum = n
		}
		item := dy.LiveProductFirstCate{
			Name:       k,
			ProductNum: productNum,
			Cate:       []dy.LiveProductSecondCate{},
		}
		if len(secondCateList) > 0 {
			item.Cate = secondCateList
		}
		productCount.CateList = append(productCount.CateList, item)
	}
	cateListByte, _ := jsoniter.Marshal(productCount)
	var timeout time.Duration = 60
	if gmvNum == productNum {
		timeout = 1800
	}
	_ = global.Cache.Set(cKey, string(cateListByte), timeout)
	return
}

//达人直播间搜索
func (receiver *EsLiveBusiness) SearchProductRooms(productId, keyword, sortStr, orderBy string, page, size int, startTime, endTime time.Time) (list []es.EsAuthorLiveProduct, total int, comErr global.CommonError) {
	if sortStr == "" {
		sortStr = "predict_gmv"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"predict_gmv", "predict_sales"}) {
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
	esTable := GetESTableByTime(es.DyRoomProductRecords, startTime, endTime)
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
		//todo gmv处理
		if v.RealGmv > 0 {
			var sale float64 = 0
			if v.Price > 0 {
				sale = math.Floor(v.RealGmv / v.Price)
			}
			list[k].PredictGmv = v.RealGmv
			list[k].PredictSales = sale
		}
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
