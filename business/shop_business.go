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
	"fmt"
	jsoniter "github.com/json-iterator/go"
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
	list []entity.DyShopProductAnalysis, total int, comError global.CommonError) {
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
		} else {
			productMap[v.ProductId] = v
		}
	}
	list = []entity.DyShopProductAnalysis{}
	for _, v := range productMap {
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

func (receiver *ShopBusiness) ShopProductAnalysisCount(shopId, keyword string, startTime, stopTime time.Time) (
	count []dy.DyCate, comError global.CommonError) {
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
		_ = global.Cache.Set(cacheKey, utils.SerializeData(hbaseList), 300)
	}
	productMap := map[string]entity.DyShopProductAnalysis{}
	cateMap := map[string]int{}
	cateSonMap := map[string]map[string]int{}
	for _, v := range hbaseList {
		if v.DcmLevelFirst == "" {
			v.DcmLevelFirst = "其他"
		}
		if v.FirstCname == "" {
			v.FirstCname = "其他"
		}
		if keyword != "" {
			if strings.Index(v.Title, keyword) < 0 {
				continue
			}
		}
		if _, exist := productMap[v.ProductId]; !exist {
			if n, ok := cateMap[v.DcmLevelFirst]; ok {
				cateMap[v.DcmLevelFirst] = n + 1
			} else {
				cateMap[v.DcmLevelFirst] = 1
				cateSonMap[v.DcmLevelFirst] = map[string]int{}
			}
			if _, ok := cateSonMap[v.DcmLevelFirst][v.FirstCname]; !ok {
				cateSonMap[v.DcmLevelFirst][v.FirstCname] = 1
			} else {
				cateSonMap[v.DcmLevelFirst][v.FirstCname] += 1
			}
			productMap[v.ProductId] = v
		}
	}
	count = []dy.DyCate{}
	for k, v := range cateMap {
		item := []dy.DyCate{}
		if k != "其他" {
			if c, ok := cateSonMap[k]; ok {
				for ck, cv := range c {
					item = append(item, dy.DyCate{
						Name:    ck,
						Num:     cv,
						SonCate: []dy.DyCate{},
					})
				}
			}
		}
		count = append(count, dy.DyCate{
			Name:    k,
			Num:     v,
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
	esShopBusiness := es.NewEsShopBusiness()
	//直播达人
	startDate := startTime.Format("20060102")
	stopDate := endTime.Format("20060102")
	allLiveList := make([]entity.DyProductAuthorAnalysis, 0)
	cacheKey := cache.GetCacheKey(cache.ShopLiveAuthorAllList, utils.Md5_encode(fmt.Sprintf("%s%s%s%s", shopId, startDate, stopDate, "")))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &allLiveList)
	} else {
		idsList, idTotal, comErr1 := esShopBusiness.GetShopLiveAuthorRowKeys(shopId, "", "", startTime, endTime)
		if comErr1 != nil {
			comErr = comErr1
			return
		}
		if idTotal == 0 {
			return
		}
		for _, v := range idsList {
			startRowKey := v.Key + "_" + startDate + "_"
			stopRowKey := v.Key + "_" + stopDate + "_99999999999999999"
			tmpList, _ := hbase.GetProductAuthorAnalysisRange(startRowKey, stopRowKey)
			allLiveList = append(allLiveList, tmpList...)
		}
		sort.Slice(allLiveList, func(i, j int) bool {
			return allLiveList[i].Date > allLiveList[j].Date
		})
		_ = global.Cache.Set(cacheKey, utils.SerializeData(allLiveList), 300)
	}
	//视频达人
	allAwemeList := make([]entity.DyProductAwemeAuthorAnalysis, 0)
	cacheAwemeKey := cache.GetCacheKey(cache.ShopAwemeAuthorAllList, utils.Md5_encode(fmt.Sprintf("%s%s%s%s", shopId, startDate, stopDate, "")))
	cacheAwemeStr := global.Cache.Get(cacheAwemeKey)
	if cacheAwemeStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &allAwemeList)
	} else {
		idsList, idTotal, comErr1 := esShopBusiness.GetShopVideoAuthorRowKeys(shopId, "", "", startTime, endTime)
		if comErr1 != nil {
			comErr = comErr1
			return
		}
		if idTotal == 0 {
			return
		}
		for _, v := range idsList {
			startRowKey := v.Key + "_" + startDate + "_"
			stopRowKey := v.Key + "_" + stopDate + "_99999999999999999"
			tmpList, _ := hbase.GetProductAwemeAuthorAnalysisRange(startRowKey, stopRowKey)
			allAwemeList = append(allAwemeList, tmpList...)
		}
		sort.Slice(allAwemeList, func(i, j int) bool {
			return allAwemeList[i].CreateSdf > allAwemeList[j].CreateSdf
		})
		_ = global.Cache.Set(cacheAwemeKey, utils.SerializeData(allAwemeList), 300)
	}
	allGmv := map[string]float64{}
	var totalGmv float64 = 0
	for _, v := range allLiveList {
		if v.Gmv == 0 {
			continue
		}
		if _, ok := allGmv[v.AuthorId]; !ok {
			allGmv[v.AuthorId] = v.Gmv
		} else {
			allGmv[v.AuthorId] += v.Gmv
		}
		totalGmv += v.Gmv
	}
	for _, v := range allAwemeList {
		if v.Gmv == 0 {
			continue
		}
		if _, ok := allGmv[v.AuthorId]; !ok {
			allGmv[v.AuthorId] = v.Gmv
		} else {
			allGmv[v.AuthorId] += v.Gmv
		}
		totalGmv += v.Gmv
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
func (receiver *ShopBusiness) ShopLiveAuthorAnalysis(shopId, keyword, tag string, startTime, endTime time.Time, minFollow, maxFollow int64, scoreType, page, pageSize int) (list []entity.DyProductAuthorAnalysis, total int, comErr global.CommonError) {
	list = []entity.DyProductAuthorAnalysis{}
	startDate := startTime.Format("20060102")
	stopDate := endTime.Format("20060102")
	allList := make([]entity.DyProductAuthorAnalysis, 0)
	cacheKey := cache.GetCacheKey(cache.ShopLiveAuthorAllList, utils.Md5_encode(fmt.Sprintf("%s%s%s%s", shopId, startDate, stopDate, keyword)))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &allList)
	} else {
		idsList, idTotal, comErr1 := es.NewEsShopBusiness().GetShopLiveAuthorRowKeys(shopId, "", keyword, startTime, endTime)
		if comErr1 != nil {
			comErr = comErr1
			return
		}
		if idTotal == 0 {
			return
		}
		for _, v := range idsList {
			startRowKey := v.Key + "_" + startDate + "_"
			stopRowKey := v.Key + "_" + stopDate + "_99999999999999999"
			tmpList, _ := hbase.GetProductAuthorAnalysisRange(startRowKey, stopRowKey)
			allList = append(allList, tmpList...)
		}
		sort.Slice(allList, func(i, j int) bool {
			return allList[i].Date > allList[j].Date
		})
		_ = global.Cache.Set(cacheKey, utils.SerializeData(allList), 300)
	}
	authorMap := map[string]entity.DyProductAuthorAnalysis{}
	authorTagMap := map[string]string{}
	authorProductMap := map[string]map[string]entity.DyProductAuthorAnalysis{}
	for _, v := range allList {
		if v.ShortId != shopId {
			continue
		}
		if at, ok := authorTagMap[v.AuthorId]; ok {
			v.ShopTags = at
		} else {
			authorTagMap[v.AuthorId] = v.ShopTags
		}
		if keyword != "" {
			if strings.Index(strings.ToLower(v.Nickname), strings.ToLower(keyword)) < 0 && v.DisplayId != keyword && v.ShortId != keyword {
				continue
			}
		}
		if scoreType != 5 && scoreType != v.Level {
			continue
		}
		if tag == "其他" {
			if v.ShopTags != "" && strings.Index(v.ShopTags, tag) < 0 {
				continue
			}
		} else {
			if tag != "" {
				if strings.Index(v.ShopTags, tag) < 0 {
					continue
				}
			}
		}
		if _, ok := authorProductMap[v.AuthorId]; !ok {
			authorProductMap[v.AuthorId] = map[string]entity.DyProductAuthorAnalysis{}
		}
		if p, ok := authorProductMap[v.AuthorId][v.ProductId]; ok {
			p.Gmv += v.Gmv
			p.Sales += v.Sales
			if p.Date < v.Date {
				p.Date = v.Date
			}
			authorProductMap[v.AuthorId][v.ProductId] = p
		} else {
			authorProductMap[v.AuthorId][v.ProductId] = v
		}
		if d, ok := authorMap[v.AuthorId]; ok {
			d.Gmv += v.Gmv
			d.Sales += v.Sales
			d.RelatedRooms = append(d.RelatedRooms, v.RelatedRooms...)
			authorMap[v.AuthorId] = d
		} else {
			authorMap[v.AuthorId] = v
		}
	}
	for _, v := range authorMap {
		if minFollow > 0 && v.FollowCount < minFollow {
			continue
		}
		if maxFollow > 0 && v.FollowCount >= maxFollow {
			continue
		}
		item := []entity.DyAuthorProductDetail{}
		if p, exist := authorProductMap[v.AuthorId]; exist {
			v.ProductNum = len(p)
			for _, p1 := range p {
				item = append(item, entity.DyAuthorProductDetail{
					Gmv:       p1.Gmv,
					ProductId: p1.ProductId,
					Sales:     p1.Sales,
					Date:      p1.Date,
				})
			}
		}
		v.Products = item
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Gmv > list[j].Gmv
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
	return
}

//小店直播达人分析统计
func (receiver *ShopBusiness) ShopLiveAuthorAnalysisCount(shopId, keyword string, startTime, endTime time.Time) (countList dy.DyProductLiveCount, comErr global.CommonError) {
	countList = dy.DyProductLiveCount{
		Tags:  []dy.DyCate{},
		Level: []dy.DyIntCate{},
	}
	startDate := startTime.Format("20060102")
	stopDate := endTime.Format("20060102")
	allList := make([]entity.DyProductAuthorAnalysis, 0)
	cacheKey := cache.GetCacheKey(cache.ShopLiveAuthorAllList, utils.Md5_encode(fmt.Sprintf("%s%s%s%s", shopId, startDate, stopDate, keyword)))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &allList)
	} else {
		idsList, idTotal, comErr1 := es.NewEsShopBusiness().GetShopLiveAuthorRowKeys(shopId, "", keyword, startTime, endTime)
		if comErr1 != nil {
			comErr = comErr1
			return
		}
		if idTotal == 0 {
			return
		}
		for _, v := range idsList {
			startRowKey := v.Key + "_" + startDate + "_"
			stopRowKey := v.Key + "_" + stopDate + "_99999999999999999"
			tmpList, _ := hbase.GetProductAuthorAnalysisRange(startRowKey, stopRowKey)
			allList = append(allList, tmpList...)
		}
		sort.Slice(allList, func(i, j int) bool {
			return allList[i].Date > allList[j].Date
		})
		_ = global.Cache.Set(cacheKey, utils.SerializeData(allList), 300)
	}
	tagsMap := map[string]int{}
	levelMap := map[int]int{}
	authorMap := map[string]string{}
	authorTagMap := map[string]string{}
	for _, v := range allList {
		if _, ok := authorMap[v.AuthorId]; ok {
			continue
		}
		if at, ok := authorTagMap[v.AuthorId]; ok {
			v.ShopTags = at
		} else {
			authorTagMap[v.AuthorId] = v.ShopTags
		}
		if keyword != "" {
			if strings.Index(v.Nickname, keyword) < 0 && v.DisplayId != keyword && v.ShortId != keyword {
				continue
			}
		}
		if v.ShopTags == "" || v.ShopTags == "null" {
			v.ShopTags = "其他"
		}
		shopTags := strings.Split(v.ShopTags, "_")
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
		authorMap[v.AuthorId] = v.AuthorId
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
	if otherLevel > 0 {
		countList.Level = append(countList.Level, dy.DyIntCate{
			Name: 0,
			Num:  otherLevel,
		})
	}
	return
}

//小店视频达人分析
func (receiver *ShopBusiness) ShopAwemeAuthorAnalysis(shopId, keyword, tag string, startTime, endTime time.Time, minFollow, maxFollow int64, scoreType, page, pageSize int) (list []entity.DyProductAwemeAuthorAnalysis, total int, comErr global.CommonError) {
	list = []entity.DyProductAwemeAuthorAnalysis{}
	startDate := startTime.Format("20060102")
	stopDate := endTime.Format("20060102")
	allList := make([]entity.DyProductAwemeAuthorAnalysis, 0)
	cacheKey := cache.GetCacheKey(cache.ShopAwemeAuthorAllList, utils.Md5_encode(fmt.Sprintf("%s%s%s%s", shopId, startDate, stopDate, keyword)))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &allList)
	} else {
		idsList, idTotal, comErr1 := es.NewEsShopBusiness().GetShopVideoAuthorRowKeys(shopId, "", keyword, startTime, endTime)
		if comErr1 != nil {
			comErr = comErr1
			return
		}
		if idTotal == 0 {
			return
		}
		for _, v := range idsList {
			startRowKey := v.Key + "_" + startDate + "_"
			stopRowKey := v.Key + "_" + stopDate + "_99999999999999999"
			tmpList, _ := hbase.GetProductAwemeAuthorAnalysisRange(startRowKey, stopRowKey)
			allList = append(allList, tmpList...)
		}
		sort.Slice(allList, func(i, j int) bool {
			return allList[i].CreateSdf > allList[j].CreateSdf
		})
		_ = global.Cache.Set(cacheKey, utils.SerializeData(allList), 300)
	}
	authorMap := map[string]entity.DyProductAwemeAuthorAnalysis{}
	authorIds := make([]string, 0)
	authorTagMap := map[string]string{}
	authorProductMap := map[string]map[string]entity.DyProductAwemeAuthorAnalysis{}
	for _, v := range allList {
		if scoreType != 5 && scoreType != v.Level {
			continue
		}
		if at, ok := authorTagMap[v.AuthorId]; ok {
			v.FirstName = at
		} else {
			authorTagMap[v.AuthorId] = v.FirstName
		}
		if tag == "其他" {
			if v.FirstName != "" && strings.Index(v.FirstName, tag) < 0 {
				continue
			}
		} else {
			if tag != "" {
				if strings.Index(v.FirstName, tag) < 0 {
					continue
				}
			}
		}
		if keyword != "" {
			if strings.Index(strings.ToLower(v.Nickname), strings.ToLower(keyword)) < 0 && v.DisplayId != keyword && v.ShortId != keyword {
				continue
			}
		}
		if _, ok := authorProductMap[v.AuthorId]; !ok {
			authorProductMap[v.AuthorId] = map[string]entity.DyProductAwemeAuthorAnalysis{}
		}
		if p, ok := authorProductMap[v.AuthorId][v.ProductId]; ok {
			p.Gmv += v.Gmv
			p.Sales += v.Sales
			if p.CreateSdf < v.CreateSdf {
				p.CreateSdf = v.CreateSdf
			}
			authorProductMap[v.AuthorId][v.ProductId] = p
		} else {
			authorProductMap[v.AuthorId][v.ProductId] = v
		}
		if d, ok := authorMap[v.AuthorId]; ok {
			d.Gmv += v.Gmv
			d.Sales += v.Sales
			d.DiggCount += v.DiggCount
			d.RelatedAwemes = append(d.RelatedAwemes, v.RelatedAwemes...)
			authorMap[v.AuthorId] = d
		} else {
			authorMap[v.AuthorId] = v
			authorIds = append(authorIds, v.AuthorId)
		}
	}
	for _, v := range authorMap {
		if minFollow > 0 && v.FollowCount < minFollow {
			continue
		}
		if maxFollow > 0 && v.FollowCount >= maxFollow {
			continue
		}
		item := []entity.DyAuthorProductDetail{}
		if p, exist := authorProductMap[v.AuthorId]; exist {
			v.ProductNum = len(p)
			for _, p1 := range p {
				item = append(item, entity.DyAuthorProductDetail{
					Gmv:       p1.Gmv,
					ProductId: p1.ProductId,
					Sales:     p1.Sales,
					Date:      p1.CreateSdf,
				})
			}
		}
		v.Products = item
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Sales > list[j].Sales
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
	return
}

//小店视频达人分析统计
func (receiver *ShopBusiness) ShopAwemeAuthorAnalysisCount(shopId, keyword string, startTime, endTime time.Time) (countList dy.DyProductAwemeCount, comErr global.CommonError) {
	countList = dy.DyProductAwemeCount{
		Tags:  []dy.DyCate{},
		Level: []dy.DyIntCate{},
	}
	startDate := startTime.Format("20060102")
	stopDate := endTime.Format("20060102")
	allList := make([]entity.DyProductAwemeAuthorAnalysis, 0)
	cacheKey := cache.GetCacheKey(cache.ShopAwemeAuthorAllList, utils.Md5_encode(fmt.Sprintf("%s%s%s%s", shopId, startDate, stopDate, keyword)))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &allList)
	} else {
		idsList, idTotal, comErr1 := es.NewEsShopBusiness().GetShopVideoAuthorRowKeys(shopId, "", keyword, startTime, endTime)
		if comErr1 != nil {
			comErr = comErr1
			return
		}
		if idTotal == 0 {
			return
		}
		for _, v := range idsList {
			startRowKey := v.Key + "_" + startDate + "_"
			stopRowKey := v.Key + "_" + stopDate + "_99999999999999999"
			tmpList, _ := hbase.GetProductAwemeAuthorAnalysisRange(startRowKey, stopRowKey)
			allList = append(allList, tmpList...)
		}
		sort.Slice(allList, func(i, j int) bool {
			return allList[i].CreateSdf > allList[j].CreateSdf
		})
		_ = global.Cache.Set(cacheKey, utils.SerializeData(allList), 300)
	}
	tagsMap := map[string]int{}
	levelMap := map[int]int{}
	authorMap := map[string]string{}
	authorTagMap := map[string]string{}
	for _, v := range allList {
		if _, ok := authorMap[v.AuthorId]; ok {
			continue
		}
		if at, ok := authorTagMap[v.AuthorId]; ok {
			v.FirstName = at
		} else {
			authorTagMap[v.AuthorId] = v.FirstName
		}
		if keyword != "" {
			if strings.Index(v.Nickname, keyword) < 0 && v.DisplayId != keyword && v.ShortId != keyword {
				continue
			}
		}
		if v.FirstName == "" || v.FirstName == "null" {
			v.FirstName = "其他"
		}
		shopTags := strings.Split(v.FirstName, "_")
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
		authorMap[v.AuthorId] = v.AuthorId
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
	if otherLevel > 0 {
		countList.Level = append(countList.Level, dy.DyIntCate{
			Name: 0,
			Num:  otherLevel,
		})
	}
	return
}
