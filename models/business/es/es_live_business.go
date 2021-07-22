package es

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/business"
	"dongchamao/models/es"
	"dongchamao/services/dyimg"
	"dongchamao/services/elasticsearch"
	"dongchamao/structinit/repost/dy"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"math"
	"strings"
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
	firstDay, _ := time.ParseInLocation("20060102", "20210701", time.Local)
	if startDate.Before(firstDay) {
		startDate = firstDay
	}
	tableArr := make([]string, 0)
	begin := startDate
	for {
		if begin.After(endDate) {
			break
		}
		tableArr = append(tableArr, fmt.Sprintf(es.DyAuthorLiveRecords, begin.Format("20060102")))
		begin = begin.AddDate(0, 0, 1)
	}
	if len(tableArr) == 0 {
		return
	}
	esTable := strings.Join(tableArr, ",")
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("author_id", authorId)
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

//直播间筛选
func (receiver *EsLiveBusiness) RoomProductByRoomId(roomId, keyword, sortStr, orderBy, firstLabel, secondLabel, thirdLabel string, page, pageSize int) (list []es.EsAuthorLiveProduct, total int, comErr global.CommonError) {
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
	if pageSize > 30 {
		comErr = global.NewError(4000)
		return
	}
	liveBusiness := business.NewLiveBusiness()
	roomInfo, comErr := liveBusiness.HbaseGetLiveInfo(roomId)
	if comErr != nil {
		return
	}
	date := time.Unix(roomInfo.DiscoverTime, 0).Format("20060102")
	esTable := fmt.Sprintf(es.DyRoomProductRecords, date)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("room_id", roomId)
	if keyword != "" {
		esQuery.SetMatchPhrase("title", keyword)
	}
	if firstLabel != "" {
		esQuery.SetTerm("first_cname", firstLabel)
	}
	if secondLabel != "" {
		esQuery.SetTerm("second_cname", secondLabel)
	}
	if thirdLabel != "" {
		esQuery.SetTerm("third_cname", thirdLabel)
	}
	start := (page - 1) * pageSize
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit(start, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	for k, v := range list {
		list[k].Cover = dyimg.Product(v.Cover)
		//真实gmv存在，按gmv处理
		if v.RealGmv > 0 {
			list[k].PredictGmv = v.RealGmv
			if v.Price > 0 {
				list[k].PredictSales = math.Floor(v.RealGmv / v.Price)
			}
		} else {
			list[k].PredictSales = math.Floor(v.PredictSales)
		}
		if v.PredictSales > 0 {
			list[k].BuyRate = float64(v.Pv) / v.PredictSales
		}
	}
	total = esMultiQuery.Count
	return
}

//直播间商品分类统计
func (receiver *EsLiveBusiness) AllRoomProductCateByRoomId(roomId string) (productCount dy.LiveProductCateCount) {
	cKey := cache.GetCacheKey(cache.LiveRoomProductCount, roomId)
	productCountJson := global.Cache.Get(cKey)
	if productCountJson != "" {
		jsoniter.Unmarshal([]byte(productCountJson), &productCount)
		return
	}
	liveBusiness := business.NewLiveBusiness()
	roomInfo, comErr := liveBusiness.HbaseGetLiveInfo(roomId)
	if comErr != nil {
		return
	}
	date := time.Unix(roomInfo.DiscoverTime, 0).Format("20060102")
	esTable := fmt.Sprintf(es.DyRoomProductRecords, date)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("room_id", roomId)
	list := make([]es.EsAuthorLiveProduct, 0)
	results := esMultiQuery.
		SetFields("first_cname", "second_cname", "third_cname").
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
	for _, v := range list {
		if v.RealGmv > 0 {
			gmvNum++
			productCount.Gmv += v.RealGmv
			if v.Price > 0 {
				productCount.Sales += math.Floor(v.RealGmv / v.Price)
			}
		} else {
			productCount.Gmv += v.PredictGmv
			productCount.Sales += v.PredictSales
		}
		productCount.ProductNum++
		if v.FirstCname == "" {
			v.FirstCname = "其他"
		}
		if _, ok := firstCateMap[v.FirstCname]; !ok {
			firstCateMap[v.FirstCname] = map[string]bool{}
		}
		if _, ok := firstCateCountMap[v.FirstCname]; !ok {
			firstCateCountMap[v.FirstCname] = 1
		} else {
			firstCateCountMap[v.FirstCname] += 1
		}
		if v.SecondCname == "" {
			continue
		}
		firstCateMap[v.FirstCname][v.SecondCname] = true
		if _, ok := secondCateMap[v.SecondCname]; !ok {
			secondCateMap[v.SecondCname] = map[string]bool{}
		}
		if v.ThirdCname == "" {
			continue
		}
		secondCateMap[v.SecondCname][v.ThirdCname] = true
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
	if gmvNum == productCount.ProductNum {
		timeout = 1800
	}
	_ = global.Cache.Set(cKey, string(cateListByte), timeout)
	return
}
