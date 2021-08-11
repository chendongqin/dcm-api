package business

import (
	"dongchamao/business/es"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	"dongchamao/models/entity"
	"dongchamao/models/repost/dy"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"
)

type ProductBusiness struct {
}

func NewProductBusiness() *ProductBusiness {
	return new(ProductBusiness)
}

func (receiver *ProductBusiness) GetCacheProductCate(enableCache bool) []dy.DyCate {
	memberKey := cache.GetCacheKey(cache.ConfigKeyCache, "product_cate")
	pCate := make([]dy.DyCate, 0)
	if enableCache == true {
		jsonStr := global.Cache.Get(memberKey)
		if jsonStr != "" {
			jsonData := utils.DeserializeData(jsonStr)
			_ = jsoniter.Unmarshal([]byte(jsonData), &pCate)
			return pCate
		}
	}
	allList := make([]dcm.DcProductCate, 0)
	_ = dcm.GetSlaveDbSession().Where("level<4").Desc("weight").Find(&allList)
	firstList := make([]dcm.DcProductCate, 0)
	secondMap := map[int][]dcm.DcProductCate{}
	thirdMap := map[int][]dcm.DcProductCate{}
	for _, v := range allList {
		if v.Level == 3 {
			if _, ok := thirdMap[v.ParentId]; !ok {
				thirdMap[v.ParentId] = []dcm.DcProductCate{}
			}
			thirdMap[v.ParentId] = append(thirdMap[v.ParentId], v)
		} else if v.Level == 2 {
			if _, ok := secondMap[v.ParentId]; !ok {
				secondMap[v.ParentId] = []dcm.DcProductCate{}
			}
			secondMap[v.ParentId] = append(secondMap[v.ParentId], v)
		} else {
			firstList = append(firstList, v)
		}
	}
	for _, v := range firstList {
		item := dy.DyCate{
			Name:    v.Name,
			SonCate: []dy.DyCate{},
		}
		if s, ok := secondMap[v.Id]; ok {
			for _, s1 := range s {
				sonCate := make([]dy.DyCate, 0)
				if t, ok2 := thirdMap[s1.Id]; ok2 {
					for _, t1 := range t {
						sonCate = append(sonCate, dy.DyCate{
							Name:    t1.Name,
							SonCate: []dy.DyCate{},
						})
					}
				}
				item.SonCate = append(item.SonCate, dy.DyCate{
					Name:    s1.Name,
					SonCate: sonCate,
				})
			}
		}
		pCate = append(pCate, item)
	}
	if len(pCate) > 0 {
		jsonData := utils.SerializeData(pCate)
		_ = global.Cache.Set(memberKey, jsonData, 86400)
	}
	return pCate
}

//获取商品url
func (receiver *ProductBusiness) GetProductUrl(platform, productId string) string {
	url := ""
	switch platform {
	case "淘宝":
		url = "https://item.taobao.com/item.htm?id=%s"
	case "京东":
		url = " https://item.m.jd.com/product/%s.html"
	case "天猫":
		url = " https://detail.tmall.com/item.htm?id=%s"
	case "苏宁":
		url = "https://m.suning.com/product/0000000000/0000000%s.html"
	case "小店":
		url = "https://haohuo.jinritemai.com/views/product/item2?id=%s"
	case "唯品会":
		url = "https://m.vip.com/public/go.html?pid=%s"
	case "考拉":
		url = "https://m-goods.kaola.com/product/%s.html"
	}
	if url != "" {
		url = fmt.Sprintf(url, productId)
	}
	return url
}

func (receiver *ProductBusiness) ProductAuthorAnalysis(productId, keyword, tag string, startTime, endTime time.Time, minFollow, maxFollow int64, scoreType, page, pageSize int) (list []entity.DyProductAuthorAnalysis, total int, comErr global.CommonError) {
	esProductBusiness := es.NewEsProductBusiness()
	if tag == "" && minFollow == 0 && maxFollow == 0 && scoreType == 5 {
		searchList, searchTotal, err := esProductBusiness.SearchRangeDateList(productId, keyword, startTime, endTime, page, pageSize)
		if err != nil {
			comErr = err
			return
		}
		total = searchTotal
		for _, v := range searchList {
			rowKey := v.ProductId + "_" + v.CreateSdf + "_" + v.AuthorId
			data, err := hbase.GetProductAuthorAnalysis(rowKey)
			if err == nil {
				list = append(list, data)
			}
		}
	}
	startRow, stopRow, total, comErr := esProductBusiness.SearchRangeDateRowKey(productId, keyword, startTime, endTime)
	if comErr != nil {
		return
	}
	startRowKey := startRow.ProductId + "_" + startRow.CreateSdf + "_" + startRow.AuthorId
	stopRowKey := stopRow.ProductId + "_" + stopRow.CreateSdf + "_" + stopRow.AuthorId
	allList, _ := hbase.GetProductAuthorAnalysisRange(startRowKey, stopRowKey)
	lastRow, err := hbase.GetProductAuthorAnalysis(stopRowKey)
	if err == nil {
		allList = append(allList, lastRow)
	}
	for _, v := range allList {
		if keyword != "" {
			if strings.Index(v.NickName, keyword) < 0 && v.DisplayId != keyword && v.ShortId != keyword {
				continue
			}
		}
		if minFollow > 0 && v.FollowCount < minFollow {
			continue
		}
		if maxFollow > 0 && v.FollowCount >= maxFollow {
			continue
		}
		if scoreType != 5 && scoreType != v.Level {
			continue
		}
		if tag != "" && strings.Index(v.ShopTags, tag) < 0 {
			continue
		}
		list = append(list, v)
	}
	total = len(list)
	start := (page - 1) * pageSize
	end := start + pageSize
	if total < end {
		end = total
	}
	list = list[start:end]
	return
}

func (receiver *ProductBusiness) ProductAuthorAnalysisCount(productId, keyword string, startTime, endTime time.Time) (countList []dy.DyCate, comErr global.CommonError) {
	countList = []dy.DyCate{}
	cKey := cache.GetCacheKey(cache.ProductAuthorCount, startTime.Format("20060102"), endTime.Format("20060102"))
	if keyword == "" {
		countJson := global.Cache.Get(cKey)
		if countJson != "" {
			countJson = utils.DeserializeData(countJson)
			_ = jsoniter.Unmarshal([]byte(countJson), &countList)
			return
		}
	}
	esProductBusiness := es.NewEsProductBusiness()
	startRow, stopRow, _, comErr := esProductBusiness.SearchRangeDateRowKey(productId, keyword, startTime, endTime)
	if comErr != nil {
		return
	}
	startRowKey := startRow.ProductId + "_" + startRow.CreateSdf + "_" + startRow.AuthorId
	stopRowKey := stopRow.ProductId + "_" + stopRow.CreateSdf + "_" + stopRow.AuthorId
	allList, _ := hbase.GetProductAuthorAnalysisRange(startRowKey, stopRowKey)
	lastRow, err := hbase.GetProductAuthorAnalysis(stopRowKey)
	if err == nil {
		allList = append(allList, lastRow)
	}
	tagsMap := map[string]int{}
	for _, v := range allList {
		if keyword != "" {
			if strings.Index(v.NickName, keyword) < 0 && v.DisplayId != keyword && v.ShortId != keyword {
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
	}
	for k, v := range tagsMap {
		countList = append(countList, dy.DyCate{
			Name: k,
			Num:  v,
		})
	}
	if keyword == "" && len(countList) > 0 {
		countJson := utils.SerializeData(countList)
		_ = global.Cache.Set(cKey, countJson, 300)
	}
	return
}
