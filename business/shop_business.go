package business

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	"dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
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
	list []entity.DyShopProductAnalysis, count []dy.NameValueChart, total int, comError global.CommonError) {
	startKey := ""
	stopKey := "99999999999999999999"
	//if keyword != "" {
	//	startKey, stopKey = es.NewEsProductBusiness().GetSearchRowKey(keyword, category)
	//}
	//if stopKey == "" {
	//	return
	//}
	hbaseList := make([]entity.DyShopProductAnalysis, 0)
	if startKey != stopKey {
		hbaseList, _ = hbase.GetShopProductAnalysisRangDate(shopId, startKey, stopKey, startTime, stopTime)
	}
	hbaseData, err1 := hbase.GetShopProductAnalysisByDate(shopId, stopKey, stopTime.Format("20060102"))
	if err1 == nil {
		hbaseList = append(hbaseList, hbaseData)
	}
	productMap := map[string]entity.DyShopProductAnalysis{}
	cateMap := map[string]int{}
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
			if n, ok := cateMap[v.DcmLevelFirst]; ok {
				cateMap[v.DcmLevelFirst] = n + 1
			} else {
				cateMap[v.DcmLevelFirst] = 1
			}
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
	count = make([]dy.NameValueChart, 0)
	for k, v := range cateMap {
		count = append(count, dy.NameValueChart{
			Name:  k,
			Value: v,
		})
	}
	sort.Slice(count, func(i, j int) bool {
		if count[i].Name == "其他" {
			return false
		}
		if count[j].Name == "其他" {
			return true
		}
		return count[i].Value > count[j].Value
	})
	return
}
