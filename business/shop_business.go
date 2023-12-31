package business

import (
	"dongchamao/business/es"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	"dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	jsoniter "github.com/json-iterator/go"
	"math"
	"sort"
	"strings"
	"time"
)

type ShopBusiness struct {
}

func NewShopBusiness() *ShopBusiness {
	return new(ShopBusiness)
}

//小店商品分析
func (receiver *ShopBusiness) ShopProductAnalysis(shopId, keyword, category, sortStr, orderBy string, startTime, stopTime time.Time, page, pageSize int) (
	list []entity.DyShopProductAnalysis, total int, totalSales int64, totalGmv float64, comError global.CommonError) {
	hbaseList := make([]entity.DyShopProductAnalysis, 0)
	cacheKey := cache.GetCacheKey(cache.ShopProductAnalysisScanList, startTime.Format("20060102"), stopTime.Format("20060102"), shopId)
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &hbaseList)
	} else {
		startKey := ""
		stopKey := "9999999999999999"
		hbaseList, _ = hbase.GetShopProductAnalysisRangDate(shopId, startKey, stopKey, startTime, stopTime)
		_ = global.Cache.Set(cacheKey, utils.SerializeData(hbaseList), 300)
	}
	productMap := map[string]entity.DyShopProductAnalysis{}
	for _, v := range hbaseList {
		if v.DcmLevelFirst == "" {
			v.DcmLevelFirst = "其他"
		}
		if keyword != "" {
			if strings.Index(v.Title, keyword) < 0 {
				continue
			}
		}
		if category != "" {
			if v.DcmLevelFirst != category {
				continue
			}
		}
		if p, exist := productMap[v.ProductId]; exist {
			if p.Price < v.Price {
				p.Price = v.Price
			}
			if p.CommissionRate < v.CommissionRate {
				p.CommissionRate = v.CommissionRate
			}
			p.Gmv += v.Gmv
			p.Sales += v.Sales
			if p.Date < v.Date {
				p.MonthPvCount = v.MonthPvCount
				p.MonthCvr = utils.RateMin(v.MonthCvr)
			}
			productMap[v.ProductId] = p
		} else {
			productMap[v.ProductId] = v
		}
	}
	list = []entity.DyShopProductAnalysis{}
	for _, v := range productMap {
		totalSales = totalSales + v.Sales
		totalGmv = totalGmv + v.Gmv
		list = append(list, v)
	}
	//排序
	sort.Slice(list, func(i, j int) bool {
		var left, right float64
		switch sortStr {
		case "price":
			left = list[i].Price
			right = list[j].Price
		case "gmv":
			left = list[i].Gmv
			right = list[j].Gmv
		case "sales":
			left = utils.ToFloat64(list[i].Sales)
			right = utils.ToFloat64(list[j].Sales)
		case "month_pv_count":
			left = utils.ToFloat64(list[i].MonthPvCount)
			right = utils.ToFloat64(list[j].MonthPvCount)
		case "month_cvr":
			left = utils.ToFloat64(list[i].MonthCvr)
			right = utils.ToFloat64(list[j].MonthCvr)
		case "commission_rate":
			left = utils.ToFloat64(list[i].CommissionRate)
			right = utils.ToFloat64(list[j].CommissionRate)
		}
		if orderBy == "desc" {
			return left > right
		}
		return right > left
	})
	start := (page - 1) * pageSize
	end := start + pageSize
	total = len(list)
	if start > total {
		list = []entity.DyShopProductAnalysis{}
		return
	}
	if total < end {
		end = total
	}
	if total > 0 {
		list = list[start:end]
	}
	productIds := []string{}
	for _, v := range list {
		productIds = append(productIds, v.ProductId)
	}
	productInfoMap, _ := hbase.GetProductByIds(productIds)
	for k, v := range list {
		if productInfo, exist := productInfoMap[v.ProductId]; exist {
			list[k].ProductStatus = productInfo.Status
		}
		list[k].Image = dyimg.Product(v.Image)
		list[k].ProductId = IdEncrypt(v.ProductId)
	}
	return
}

func (receiver *ShopBusiness) ShopProductAnalysisCount(shopId, keyword string, startTime, stopTime time.Time) (count []dy.DyCate, comError global.CommonError) {
	hbaseList := make([]entity.DyShopProductAnalysis, 0)
	cacheKey := cache.GetCacheKey(cache.ShopProductAnalysisCountScanList, startTime.Format("20060102"), stopTime.Format("20060102"), shopId)
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &hbaseList)
	} else {
		startKey := ""
		stopKey := "99999999999999999999"
		hbaseList, _ = hbase.GetShopProductAnalysisRangDate(shopId, startKey, stopKey, startTime, stopTime)
		sort.Slice(hbaseList, func(i, j int) bool {
			return hbaseList[i].Date > hbaseList[j].Date
		})
		_ = global.Cache.Set(cacheKey, utils.SerializeData(hbaseList), 300)
	}
	countMap := map[string]map[string]int{}
	countSonMap := map[string]map[string]map[string]int{}
	productMap := map[string]string{}
	for _, v := range hbaseList {
		if _, exist := productMap[v.ProductId]; exist {
			continue
		}
		if keyword != "" {
			if strings.Index(v.Title, keyword) < 0 {
				continue
			}
		}
		if countMap[v.DcmLevelFirst] == nil {
			countMap[v.DcmLevelFirst] = make(map[string]int)
			countSonMap[v.DcmLevelFirst] = make(map[string]map[string]int)
		}
		if countSonMap[v.DcmLevelFirst][v.FirstCname] == nil {
			countSonMap[v.DcmLevelFirst][v.FirstCname] = make(map[string]int)
		}
		if v.DcmLevelFirst == "" {
			v.DcmLevelFirst = "其他"
		}
		if v.FirstCname == "" {
			v.FirstCname = "其他"
		}
		countMap[v.DcmLevelFirst][v.ProductId]++
		countSonMap[v.DcmLevelFirst][v.FirstCname][v.ProductId]++
	}
	count = []dy.DyCate{}
	for k, v := range countMap {
		item := []dy.DyCate{}
		if k != "其他" {
			if c, ok := countSonMap[k]; ok {
				for ck, cv := range c {
					item = append(item, dy.DyCate{
						Name:    ck,
						Num:     len(cv),
						SonCate: []dy.DyCate{},
					})
				}
			}
		}
		count = append(count, dy.DyCate{
			Name:    k,
			Num:     len(v),
			SonCate: item,
		})
	}
	sort.Slice(count, func(i, j int) bool {
		if count[i].Name == "其他" {
			return false
		}
		if count[j].Name == "其他" {
			return true
		}
		return count[i].Num > count[j].Num
	})
	return
}

//达人概览
func (receiver *ShopBusiness) ShopAuthorView(shopId string, startTime, endTime time.Time) (
	allTop5 []dy.NameValueFloat64PercentChart, comErr global.CommonError) {
	allTop5 = []dy.NameValueFloat64PercentChart{}
	//直播达人
	allLiveList, _, comErr := es.NewEsLiveBusiness().SumSearchLiveAuthor("", shopId, startTime, endTime)
	//视频达人
	allAwemeList, _, comErr := es.NewEsVideoBusiness().SumSearchAwemeAuthor("", shopId, startTime, endTime)
	allGmv := map[string]float64{}
	var totalGmv float64 = 0
	for _, l := range allLiveList {
		if len(l.Data.Hits.Hits) == 0 {
			continue
		}
		v := l.Data.Hits.Hits[0].Source
		if v.PredictGmv == 0 {
			continue
		}
		if _, ok := allGmv[v.AuthorId]; !ok {
			allGmv[v.AuthorId] = v.PredictGmv
		} else {
			allGmv[v.AuthorId] += v.PredictGmv
		}
		totalGmv += v.PredictGmv
	}
	for _, l := range allAwemeList {
		if len(l.Data.Hits.Hits) == 0 {
			continue
		}
		v := l.Data.Hits.Hits[0].Source
		if v.AwemeGmv == 0 {
			continue
		}
		if _, ok := allGmv[v.AuthorId]; !ok {
			allGmv[v.AuthorId] = v.AwemeGmv
		} else {
			allGmv[v.AuthorId] += v.AwemeGmv
		}
		totalGmv += v.AwemeGmv
	}
	for k, v := range allGmv {
		allTop5 = append(allTop5, dy.NameValueFloat64PercentChart{
			Name:  k,
			Value: v,
		})
	}
	sort.Slice(allTop5, func(i, j int) bool {
		return allTop5[i].Value > allTop5[j].Value
	})
	if len(allTop5) > 5 {
		allTop5 = allTop5[0:5]
	}
	otherGmv := totalGmv
	if totalGmv > 0 {
		for k, v := range allTop5 {
			author, _ := hbase.GetAuthor(v.Name)
			allTop5[k].Name = author.Data.Nickname
			allTop5[k].Percent = v.Value / totalGmv
			totalGmv -= v.Value
		}
	}
	if otherGmv > 0 {
		allTop5 = append(allTop5, dy.NameValueFloat64PercentChart{
			Name:    "其他",
			Value:   otherGmv,
			Percent: otherGmv / totalGmv,
		})
	}
	return
}

//小店直播达人分析
func (receiver *ShopBusiness) ShopLiveAuthorAnalysis(shopId, keyword, tag, sortStr, orderBy string, startTime, endTime time.Time, minFollow, maxFollow int64, scoreType, page, pageSize int) (list []entity.DyProductAuthorAnalysis, total int, totalSales int64, totalGmv float64, comErr global.CommonError) {
	list = []entity.DyProductAuthorAnalysis{}
	allList, _, comErr := es.NewEsLiveBusiness().SumSearchLiveAuthor("", shopId, startTime, endTime)
	for _, l := range allList {
		if len(l.Data.Hits.Hits) == 0 {
			continue
		}
		v := l.Data.Hits.Hits[0].Source
		v.PredictGmv = l.PredictGmv.Value
		v.PredictSales = math.Floor(l.PredictSales.Value)
		if keyword != "" {
			if strings.Index(strings.ToLower(v.Nickname), strings.ToLower(keyword)) < 0 && v.DisplayId != keyword && v.ShortId != keyword {
				continue
			}
		}
		if scoreType != -1 && scoreType != v.Level {
			continue
		}
		if tag == "其他" {
			if v.Tags != "" && strings.Index(v.Tags, tag) < 0 {
				continue
			}
		} else {
			if tag != "" {
				if strings.Index(v.Tags, tag) < 0 {
					continue
				}
			}
		}
		if minFollow > 0 && v.FollowerCount < minFollow {
			continue
		}
		if maxFollow > 0 && v.FollowerCount >= maxFollow {
			continue
		}
		totalSales = totalSales + utils.ToInt64(math.Floor(v.PredictSales))
		totalGmv = totalGmv + v.PredictGmv
		list = append(list, entity.DyProductAuthorAnalysis{
			AuthorId:    v.AuthorId,
			DisplayId:   v.DisplayId,
			FollowCount: v.FollowerCount,
			Gmv:         v.PredictGmv,
			Nickname:    v.Nickname,
			Avatar:      v.Avatar,
			Price:       v.Price,
			ProductId:   v.ProductId,
			Sales:       utils.ToInt64(math.Floor(v.PredictSales)),
			Score:       v.Score,
			Level:       v.Level,
			ShopTags:    v.Tags,
			ShortId:     v.ShortId,
			ShopId:      v.ShopId,
			Date:        v.Dt,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		leftMore := true
		switch sortStr {
		case "gmv":
			leftMore = list[i].Gmv > list[j].Gmv
		case "sales":
			leftMore = list[i].Sales > list[j].Sales
		default:
			return true
		}
		if orderBy == "asc" {
			return !leftMore
		}
		return leftMore
	})
	total = len(list)
	start := (page - 1) * pageSize
	end := start + pageSize
	if total < end {
		end = total
	}
	if start > total {
		start = total
	}
	if total == 0 {
		return
	}
	list = list[start:end]
	authorIds := []string{}
	for _, v := range list {
		authorIds = append(authorIds, v.AuthorId)
	}
	roomMap, productMap, _ := es.NewEsLiveBusiness().CountSearchLiveAuthorRoomProductNum("", shopId, authorIds, startTime, endTime)
	for k, v := range list {
		if num, exist := roomMap[v.AuthorId]; exist {
			list[k].RoomNum = num
		}
		if num, exist := productMap[v.AuthorId]; exist {
			list[k].ProductNum = num
		}
	}
	return
}

//小店直播达人分析统计
func (receiver *ShopBusiness) ShopLiveAuthorAnalysisCount(shopId, keyword string, startTime, endTime time.Time) (countList dy.DyProductLiveCount, comErr global.CommonError) {
	countList = dy.DyProductLiveCount{
		Tags:  []dy.DyCate{},
		Level: []dy.DyIntCate{},
	}
	cKey := cache.GetCacheKey(cache.ShopAuthorCount, shopId, startTime.Format("20060102"), endTime.Format("20060102"))
	if keyword == "" {
		countJson := global.Cache.Get(cKey)
		if countJson != "" {
			countJson = utils.DeserializeData(countJson)
			_ = jsoniter.Unmarshal([]byte(countJson), &countList)
			return
		}
	}
	allList, _, comErr := es.NewEsLiveBusiness().SumSearchLiveAuthor("", shopId, startTime, endTime)
	tagsMap := map[string]int{}
	levelMap := map[int]int{}
	for _, l := range allList {
		if len(l.Data.Hits.Hits) == 0 {
			continue
		}
		v := l.Data.Hits.Hits[0].Source
		v.PredictGmv = l.PredictGmv.Value
		v.PredictSales = math.Floor(l.PredictSales.Value)
		if keyword != "" {
			if strings.Index(v.Nickname, keyword) < 0 && v.DisplayId != keyword && v.ShortId != keyword {
				continue
			}
		}
		if v.Tags == "" || v.Tags == "null" {
			v.Tags = "其他"
		}
		shopTags := []string{}
		if strings.Index(v.Tags, "_") >= 0 {
			shopTags = strings.Split(v.Tags, "_")
		} else {
			shopTags = strings.Split(v.Tags, "|")
		}
		for _, s := range shopTags {
			if _, ok := tagsMap[s]; ok {
				tagsMap[s] += 1
			} else {
				tagsMap[s] = 1
			}
		}
		if _, ok := levelMap[v.Level]; ok {
			levelMap[v.Level] += 1
		} else {
			levelMap[v.Level] = 1
		}
	}
	otherTags := 0
	otherLevel := 0
	for k, v := range tagsMap {
		if k == "其他" {
			otherTags = v
			continue
		}
		countList.Tags = append(countList.Tags, dy.DyCate{
			Name: k,
			Num:  v,
		})
	}
	sort.Slice(countList.Tags, func(i, j int) bool {
		return countList.Tags[i].Num > countList.Tags[j].Num
	})
	if otherTags > 0 {
		countList.Tags = append(countList.Tags, dy.DyCate{
			Name: "其他",
			Num:  otherTags,
		})
	}
	for k, v := range levelMap {
		if k == 0 {
			otherLevel = v
			continue
		}
		countList.Level = append(countList.Level, dy.DyIntCate{
			Name: k,
			Num:  v,
		})
	}
	sort.Slice(countList.Level, func(i, j int) bool {
		return countList.Level[i].Num > countList.Level[j].Num
	})
	if otherLevel > 0 {
		countList.Level = append(countList.Level, dy.DyIntCate{
			Name: 0,
			Num:  otherLevel,
		})
	}
	if keyword == "" && (len(countList.Tags) > 0 || len(countList.Level) > 0) {
		countJson := utils.SerializeData(countList)
		_ = global.Cache.Set(cKey, countJson, 300)
	}
	return
}
