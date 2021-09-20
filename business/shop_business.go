package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	"dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
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
func (s *ShopBusiness) ShopProductAnalysis(shopId, keyword, category, sortStr, orderBy string, startTime, stopTime time.Time, page, pageSize int) (
	list []entity.DyShopProductAnalysis, total int, comError global.CommonError) {
	hbaseList := make([]entity.DyShopProductAnalysis, 0)
	cacheKey := cache.GetCacheKey(cache.ShopProductAnalysisScanList, startTime.Format("20060102"), stopTime.Format("20060102"))
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
				p.MonthCvr = v.MonthCvr
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
		}
		if orderBy == "desc" {
			return left > right
		}
		return right > left
	})
	start := (page - 1) * pageSize
	end := start + pageSize
	total = len(list)
	if total < end {
		end = total
	}
	list = list[start:end]
	for k, v := range list {
		list[k].Image = dyimg.Product(v.Image)
		list[k].ProductId = IdEncrypt(v.ProductId)
	}
	return
}

func (s *ShopBusiness) ShopProductAnalysisCount(shopId, keyword string, startTime, stopTime time.Time) (
	count []dy.DyCate, comError global.CommonError) {
	hbaseList := make([]entity.DyShopProductAnalysis, 0)
	cacheKey := cache.GetCacheKey(cache.ShopProductAnalysisScanList, startTime.Format("20060102"), stopTime.Format("20060102"))
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
