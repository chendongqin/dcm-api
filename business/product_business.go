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
	"dongchamao/services"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

type ProductBusiness struct {
}

type ProductLiveTrend struct {
	DateChart    []string
	SaleChart    []int64
	RoomNumChart []int
	PriceChart   []float64
}
type ProductAwemeTrend struct {
	ChartList []dy.ProductSalesTrendChart
	DateChart []int64
}

func NewProductBusiness() *ProductBusiness {
	return new(ProductBusiness)
}

func (receiver *ProductBusiness) GetCacheProductCate(enableCache bool) []dy.DyCate {
	cacheKey := cache.GetCacheKey(cache.LongTimeConfigKeyCache)
	redisService := services.NewRedisService()
	pCate := make([]dy.DyCate, 0)
	if enableCache == true {
		jsonStr := redisService.Hget(cacheKey, "product_cate")
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
		_ = redisService.Hset(cacheKey, "product_cate", jsonData)
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

//达人概览
func (receiver *ProductBusiness) ProductAuthorView(productId string, startTime, endTime time.Time) (
	allTop3 []dy.NameValueInt64PercentChart, liveTop3 []dy.NameValueInt64PercentChart, awemeTop3 []dy.NameValueInt64PercentChart, comErr global.CommonError) {
	allTop3 = []dy.NameValueInt64PercentChart{}
	liveTop3 = []dy.NameValueInt64PercentChart{}
	awemeTop3 = []dy.NameValueInt64PercentChart{}
	allSales := map[string]int64{}
	liveSales := map[string]int64{}
	awemeSales := map[string]int64{}
	var liveTotalSales int64 = 0
	var awemeTotalSales int64 = 0
	//直播达人
	allLiveList, _, comErr := es.NewEsLiveBusiness().SumSearchLiveAuthor(productId, "", startTime, endTime)
	for _, l := range allLiveList {
		if len(l.Data.Hits.Hits) == 0 {
			continue
		}
		v := l.Data.Hits.Hits[0].Source
		sales := utils.ToInt64(math.Floor(l.PredictSales.Value))
		if sales == 0 {
			continue
		}
		liveSales[v.AuthorId] = sales
		allSales[v.AuthorId] = sales
		liveTotalSales += sales
	}
	//视频达人
	allAwemeList, _, comErr := es.NewEsVideoBusiness().SumSearchAwemeAuthor(productId, "", startTime, endTime)
	for _, l := range allAwemeList {
		if len(l.Data.Hits.Hits) == 0 {
			continue
		}
		v := l.Data.Hits.Hits[0].Source
		v.AwemeGmv = l.AwemeGmv.Value
		v.Sales = l.Sales.Value
		if v.Sales == 0 {
			continue
		}
		awemeSales[v.AuthorId] = v.Sales
		allSales[v.AuthorId] = v.Sales
		awemeTotalSales += v.Sales
	}
	totalSales := liveTotalSales + awemeTotalSales
	for k, v := range allSales {
		allTop3 = append(allTop3, dy.NameValueInt64PercentChart{
			Name:  k,
			Value: v,
		})
	}
	for k, v := range liveSales {
		liveTop3 = append(liveTop3, dy.NameValueInt64PercentChart{
			Name:  k,
			Value: v,
		})
	}
	for k, v := range awemeSales {
		awemeTop3 = append(awemeTop3, dy.NameValueInt64PercentChart{
			Name:  k,
			Value: v,
		})
	}
	sort.Slice(allTop3, func(i, j int) bool {
		return allTop3[i].Value > allTop3[j].Value
	})
	sort.Slice(liveTop3, func(i, j int) bool {
		return liveTop3[i].Value > liveTop3[j].Value
	})
	sort.Slice(awemeTop3, func(i, j int) bool {
		return awemeTop3[i].Value > awemeTop3[j].Value
	})
	if len(allTop3) > 3 {
		allTop3 = allTop3[0:3]
	}
	if len(liveTop3) > 3 {
		liveTop3 = liveTop3[0:3]
	}
	if len(awemeTop3) > 3 {
		awemeTop3 = awemeTop3[0:3]
	}
	otherSales := totalSales
	if totalSales > 0 {
		for k, v := range allTop3 {
			author, _ := hbase.GetAuthor(v.Name)
			allTop3[k].Name = author.Data.Nickname
			allTop3[k].Percent = float64(v.Value) / float64(totalSales)
			otherSales -= v.Value
		}
	}
	if otherSales > 0 {
		allTop3 = append(allTop3, dy.NameValueInt64PercentChart{
			Name:    "其他",
			Value:   otherSales,
			Percent: float64(otherSales) / float64(totalSales),
		})
	}
	if liveTotalSales > 0 {
		for k, v := range liveTop3 {
			author, _ := hbase.GetAuthor(v.Name)
			liveTop3[k].Name = author.Data.Nickname
			liveTop3[k].Percent = float64(v.Value) / float64(liveTotalSales)
		}
	}
	if awemeTotalSales > 0 {
		for k, v := range awemeTop3 {
			author, _ := hbase.GetAuthor(v.Name)
			awemeTop3[k].Name = author.Data.Nickname
			awemeTop3[k].Percent = float64(v.Value) / float64(awemeTotalSales)
		}
	}
	return
}

func (receiver *ProductBusiness) ProductAuthorAnalysis(productId, keyword, tag, sortStr, orderBy string, startTime, endTime time.Time, minFollow, maxFollow int64, scoreType, page, pageSize int) (list []entity.DyProductAuthorAnalysis, total int, totalSales int64, totalGmv float64, comErr global.CommonError) {
	list = []entity.DyProductAuthorAnalysis{}
	allList, _, comErr := es.NewEsLiveBusiness().SumSearchLiveAuthor(productId, "", startTime, endTime)
	if comErr != nil {
		return
	}
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
	if total == 0 {
		return
	}
	start := (page - 1) * pageSize
	end := start + pageSize
	if total < end {
		end = total
	}
	if start > total {
		start = total
	}
	if total > 0 {
		list = list[start:end]
	}
	authorIds := []string{}
	for _, v := range list {
		authorIds = append(authorIds, v.AuthorId)
	}
	roomMap, _, _ := es.NewEsLiveBusiness().CountSearchLiveAuthorRoomProductNum(productId, "", authorIds, startTime, endTime)
	for k, v := range list {
		if num, exist := roomMap[v.AuthorId]; exist {
			list[k].RoomNum = num
		}
	}
	return
}

func (receiver *ProductBusiness) ProductAuthorLiveRooms(productId, shopId, authorId string, startTime, endTime time.Time, sortStr, orderBy string, page, pageSize int) (list []entity.DyProductAuthorRelatedRoom, total int) {
	//gmv:销售额,sales：销量，start_ts：开播时间
	if sortStr == "gmv" {
		sortStr = "predict_gmv"
	} else if sortStr == "sales" {
		sortStr = "predict_sales"
	} else if sortStr == "start_ts" {
		sortStr = "live_create_time"
	} else {
		return
	}
	allList, total, _ := es.NewEsLiveBusiness().GetLiveRoomByProductAuthor(productId, shopId, authorId, sortStr, orderBy, startTime, endTime, page, pageSize)
	list = []entity.DyProductAuthorRelatedRoom{}
	for _, v := range allList {
		list = append(list, entity.DyProductAuthorRelatedRoom{
			EndTs:     v.FinishTime,
			Gmv:       v.PredictGmv,
			RoomId:    v.RoomID,
			Sales:     utils.ToInt64(math.Floor(v.PredictSales)),
			StartTs:   v.LiveCreateTime,
			Title:     v.RoomTitle,
			Cover:     v.Cover,
			TotalUser: v.TotalUser,
		})
	}
	return
}

func (receiver *ProductBusiness) ProductAuthorAnalysisCount(productId, keyword string, startTime, endTime time.Time) (countList dy.DyProductLiveCount, comErr global.CommonError) {
	countList = dy.DyProductLiveCount{
		Tags:  []dy.DyCate{},
		Level: []dy.DyIntCate{},
	}
	cKey := cache.GetCacheKey(cache.ProductAuthorCount, productId, startTime.Format("20060102"), endTime.Format("20060102"))
	if keyword == "" {
		countJson := global.Cache.Get(cKey)
		if countJson != "" {
			countJson = utils.DeserializeData(countJson)
			_ = jsoniter.Unmarshal([]byte(countJson), &countList)
			return
		}
	}
	allList, _, comErr := es.NewEsLiveBusiness().SumSearchLiveAuthor(productId, "", startTime, endTime)
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

func (receiver *ProductBusiness) ProductAwemeAuthorAnalysis(productId, shopId, keyword, tag, sortStr, orderBy string, startTime, endTime time.Time, minFollow, maxFollow int64, scoreType, page, pageSize int) (list []entity.DyProductAwemeAuthorAnalysis, total int, totalSales int64, totalGmv float64, comErr global.CommonError) {
	list = []entity.DyProductAwemeAuthorAnalysis{}
	allList, _, comErr := es.NewEsVideoBusiness().SumSearchAwemeAuthor(productId, shopId, startTime, endTime)
	if comErr != nil {
		return
	}
	keyword = strings.ToLower(keyword)
	for _, l := range allList {
		if len(l.Data.Hits.Hits) == 0 {
			continue
		}
		v := l.Data.Hits.Hits[0].Source
		v.AwemeGmv = l.AwemeGmv.Value
		v.Sales = l.Sales.Value
		if minFollow > 0 && v.FollowerCount < minFollow {
			continue
		}
		if maxFollow > 0 && v.FollowerCount >= maxFollow {
			continue
		}
		if tag != "" {
			if strings.Index(v.Tags, tag) < 0 && strings.Index(v.DcmLevelFirst, tag) < 0 {
				continue
			}
		}
		if scoreType != -1 && scoreType != v.Level {
			continue
		}
		if keyword != "" {
			if strings.Index(strings.ToLower(v.Nickname), keyword) < 0 && v.UniqueId != keyword && v.ShortId != keyword {
				continue
			}
		}
		totalSales = totalSales + v.Sales
		totalGmv = totalGmv + v.AwemeGmv
		list = append(list, entity.DyProductAwemeAuthorAnalysis{
			ProductId:   v.ProductId,
			AuthorId:    v.AuthorId,
			Nickname:    v.Nickname,
			CreateSdf:   v.DistDate,
			DisplayId:   v.UniqueId,
			ShortId:     v.ShortId,
			Score:       v.Score,
			Level:       v.Level,
			FirstName:   v.Tags,
			SecondName:  v.TagsLevelTwo,
			Avatar:      v.Avatar,
			FollowCount: v.FollowerCount,
			DiggCount:   v.DiggCount,
			Sales:       v.Sales,
			Gmv:         v.AwemeGmv,
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
	awemeMap, productMap, _ := es.NewEsVideoBusiness().CountSearchAuthorAwemeProductNum(productId, shopId, authorIds, startTime, endTime)
	for k, v := range list {
		if num, exist := awemeMap[v.AuthorId]; exist {
			list[k].AwemesNum = num
		}
		if num, exist := productMap[v.AuthorId]; exist {
			list[k].ProductNum = num
		}
	}
	return
}

func (receiver *ProductBusiness) ProductAuthorAwemes(productId, shopId, authorId string, startTime, endTime time.Time, sortStr, orderBy string, page, pageSize int) (list []entity.DyProductAuthorRelatedAweme, total int) {
	list = []entity.DyProductAuthorRelatedAweme{}
	//esProductBusiness := es.NewEsProductBusiness()
	//allList, _, _ := esProductBusiness.SearchAwemeRangeDateList(productId, shopId, authorId, startTime, endTime, 1, 10000)
	//for _, v := range allList {
	//	rowKey := v.ProductId + "_" + v.CreateSdf + "_" + v.AuthorId
	//	data, err := hbase.GetProductAwemeAuthorAnalysis(rowKey)
	//	if err == nil {
	//		list = append(list, data.RelatedAwemes...)
	//	}
	//}
	tmpSortStr := sortStr
	if !utils.InArrayString(tmpSortStr, []string{"gmv", "sales"}) {
		tmpSortStr = "gmv"
	}
	sumList, total, err := es.NewEsVideoBusiness().AuthorProductAwemeSumList(authorId, productId, shopId, tmpSortStr, orderBy, startTime, endTime, 1, 1000)
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	awemeIds := []string{}
	awemeGmvMap := map[string]float64{}
	awemeSalesMap := map[string]int64{}
	start := (page - 1) * pageSize
	end := start + pageSize
	if total < end {
		end = total
	}
	if start > total {
		start = total
	}
	if total > 0 {
		sumList = sumList[start:end]
	} else {
		sumList = sumList
	}
	for _, v := range sumList {
		awemeIds = append(awemeIds, v.Key)
		awemeGmvMap[v.Key] = v.TotalGmv.Value
		awemeSalesMap[v.Key] = v.TotalSales.Value
	}
	awemes, _ := hbase.GetVideoByIds(awemeIds)
	for _, v := range awemes {
		var gmv float64 = 0
		var sales int64 = 0
		if n, exist := awemeGmvMap[v.AwemeID]; exist {
			gmv = n
		}
		if n, exist := awemeSalesMap[v.AwemeID]; exist {
			sales = n
		}
		list = append(list, entity.DyProductAuthorRelatedAweme{
			CommentCount:    v.Data.CommentCount,
			AwemeTitle:      v.Data.AwemeTitle,
			AwemeId:         v.AwemeID,
			Sales:           sales,
			AwemeGmv:        gmv,
			DiggCount:       v.Data.DiggCount,
			ForwardCount:    v.Data.ForwardCount,
			AwemeCover:      v.Data.AwemeCover,
			AwemeCreateTime: v.Data.AwemeCreateTime,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		switch sortStr {
		case "sales":
			if orderBy == "desc" {
				return list[i].Sales > list[j].Sales
			} else {
				return list[j].Sales > list[i].Sales
			}
		case "aweme_gmv":
			if orderBy == "desc" {
				return list[i].AwemeGmv > list[j].AwemeGmv
			} else {
				return list[j].AwemeGmv > list[i].AwemeGmv
			}
		default:
			if orderBy == "desc" {
				return list[i].AwemeCreateTime > list[j].AwemeCreateTime
			} else {
				return list[j].AwemeCreateTime > list[i].AwemeCreateTime
			}
		}
	})
	//total = len(list)
	//start := (page - 1) * pageSize
	//end := start + pageSize
	//if total < end {
	//	end = total
	//}
	//list = list[start:end]
	return
}

func (receiver *ProductBusiness) NewProductAuthorAwemes(productId, authorId string, startTime, endTime time.Time, sortStr, orderBy string, page, pageSize int) (list []entity.DyProductAuthorRelatedAweme, total int) {
	list = []entity.DyProductAuthorRelatedAweme{}
	if sortStr == "gmv" {
		sortStr = "aweme_gmv"
	}
	awemeList, _, err := es.NewEsVideoBusiness().NewAuthorProductAwemeSumList(authorId, "", "", startTime, endTime, 1, 10000)
	if err != nil {
		return
	}
	for _, v := range awemeList {
		if strings.Index(v.ProductIds, productId) < 0 {
			continue
		}
		list = append(list, entity.DyProductAuthorRelatedAweme{
			CommentCount:    v.CommentCount,
			AwemeTitle:      v.AwemeTitle,
			AwemeId:         v.AwemeId,
			Sales:           v.Sales,
			AwemeGmv:        v.AwemeGmv,
			DiggCount:       v.DiggCount,
			ForwardCount:    v.ShareCount,
			AwemeCover:      v.AwemeCover,
			AwemeCreateTime: v.AwemeCreateTime,
		})
	}
	total = len(list)
	if total == 0 {
		return
	}
	sort.Slice(list, func(i, j int) bool {
		switch sortStr {
		case "sales":
			if orderBy == "desc" {
				return list[i].Sales > list[j].Sales
			} else {
				return list[j].Sales > list[i].Sales
			}
		case "aweme_gmv":
			if orderBy == "desc" {
				return list[i].AwemeGmv > list[j].AwemeGmv
			} else {
				return list[j].AwemeGmv > list[i].AwemeGmv
			}
		default:
			if orderBy == "desc" {
				return list[i].AwemeCreateTime > list[j].AwemeCreateTime
			} else {
				return list[j].AwemeCreateTime > list[i].AwemeCreateTime
			}
		}
	})
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

func (receiver *ProductBusiness) ProductAwemeAuthorAnalysisCount(productId, shopId, keyword string, startTime, endTime time.Time) (countList dy.DyProductAwemeCount, comErr global.CommonError) {
	countList = dy.DyProductAwemeCount{
		Tags:  []dy.DyCate{},
		Level: []dy.DyIntCate{},
	}
	cKey := ""
	if shopId != "" {
		cKey = cache.GetCacheKey(cache.ShopAwemeAuthorCount, shopId, startTime.Format("20060102"), endTime.Format("20060102"))
	} else {
		cKey = cache.GetCacheKey(cache.ProductAwemeAuthorCount, productId, startTime.Format("20060102"), endTime.Format("20060102"))
	}
	allList, _, comErr := es.NewEsVideoBusiness().SumSearchAwemeAuthor(productId, shopId, startTime, endTime)
	tagsMap := map[string]int{}
	levelMap := map[int]int{}
	for _, l := range allList {
		if len(l.Data.Hits.Hits) == 0 {
			continue
		}
		v := l.Data.Hits.Hits[0].Source
		if keyword != "" {
			if strings.Index(v.Nickname, keyword) < 0 && v.UniqueId != keyword && v.ShortId != keyword {
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
		_ = global.Cache.Set(cKey, countJson, 600)
	}
	return
}

func (receiver *ProductBusiness) UrlExplain(anyStr string) (id string) {
	urlInfo, err := url.Parse(anyStr)
	if err != nil {
		return
	}
	switch urlInfo.Host {
	case "v.douyin.com":
		retURL := NewDouyinBusiness().ParseDyShortUrl(anyStr)
		return receiver.UrlExplain(retURL)
	case "u.jd.com": //京东短链匹配
		jdUrl := utils.ReversedJDShortUrl(anyStr)
		return receiver.UrlExplain(jdUrl)
	case "m.tb.cn":
		revertURL := utils.GetLocation(anyStr)
		if revertURL != "" && !strings.Contains(anyStr, "m.tb.cn") {
			return receiver.UrlExplain(revertURL)
		}
		return ""
	case "m-goods.kaola.com", "item.jd.com", "item.m.jd.com", "m.suning.com", "a.m.tmall.com", "a.m.taobao.com":
		pattern := `(\d+)\.[html|htm]`
		re := regexp.MustCompile(pattern)
		s := re.FindStringSubmatch(urlInfo.Path)
		if len(s) > 1 {
			id = strings.TrimLeft(s[1], "0") //苏宁的抹去前导0
		} else {
			id = strings.ReplaceAll(urlInfo.Path, "/", "")
			id = strings.ReplaceAll(id, ".html", "")
		}
		break
	case "":
		//尝试淘口令接口
		id, _ = NewTaoBaoBusiness().TpwdConvert(anyStr)
	default:
		params, err := url.ParseQuery(urlInfo.RawQuery)
		if err != nil {
			return
		}
		idParam := params["id"]
		if len(idParam) > 0 {
			id = idParam[0]
		}
	}
	return
}

func (receiver *ProductBusiness) ExplainTaobaoShortUrl(url string) string {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	response, err := client.Do(request)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	pattern := `(https://item.taobao.com/.*?)\'`
	reg := regexp.MustCompile(pattern)
	matches := reg.FindStringSubmatch(string(content))
	if len(matches) > 1 {
		return matches[1]
	}
	//https://a.m.taobao.com/
	//http://a.m.tmall.com/
	//var url = 'https://a.m.taobao.com/i628122154669.htm?price=68&sourceType=item&sourceType=item&suid=74f0fa1e-b0b7-41b8-8326-25d157e6762b&shareUniqueId=4639705295&ut_sk=1.X4Js%2BFcRiVQDAKjZUjx8nWb6_21646297_1603682881302.Copy.1&un=6a7315ee868246b0ee428784da605ae9&share_crt_v=1&spm=a2159r.13376460.0.0&sp_tk=a0NKSGNSTE9IQXI=&cpp=1&shareurl=true&short_name=h.4159Akg&bxsign=scdV_3t5vYCjx090pisOzYWUTCTueGvuhqk8XdISQZ9jty0vONfkaESSKjThfVZqe6NauFgqcQnCQ7QT2yh0r0nD4cODKIS5p075kAwzVGmlbM';
	pattern = `(http[s*]://a.m.[tmall|taobao]+.com/.*?)\'`
	reg = regexp.MustCompile(pattern)
	matches = reg.FindStringSubmatch(string(content))
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

//获取商品详情直播销量趋势数据
func (receiver *ProductBusiness) GetProductLiveTrend(productId string, startTime, endTime time.Time) (res ProductLiveTrend) {
	infoMap, _ := hbase.GetProductLiveSalesRangDate(productId, startTime, endTime)
	dateChart := make([]string, 0)
	saleChart := make([]int64, 0)
	roomNumChart := make([]int, 0)
	priceChart := make([]float64, 0)
	beginTime := startTime
	for {
		if beginTime.After(endTime) {
			break
		}
		var sale int64 = 0
		roomNum := 0
		var price float64 = 0
		if v, ok := infoMap[beginTime.Format("20060102")]; ok {
			sale = v.Sales
			roomNum = v.RoomNum
			price = v.Price
		}
		dateChart = append(dateChart, beginTime.Format("01/02"))
		saleChart = append(saleChart, sale)
		roomNumChart = append(roomNumChart, roomNum)
		priceChart = append(priceChart, price)
		beginTime = beginTime.AddDate(0, 0, 1)
	}

	res = ProductLiveTrend{
		DateChart:    dateChart,
		SaleChart:    saleChart,
		RoomNumChart: roomNumChart,
		PriceChart:   priceChart,
	}
	return
}

//获取商品详情视频销量趋势数据
func (receiver *ProductBusiness) GetProductAwemeSalesTrend(productId string, startTime, endTime time.Time) (res ProductAwemeTrend, comErr global.CommonError) {
	hbaseDataList, comErr := hbase.GetDyProductAwemeSalesTrendRangeDate(productId, startTime, endTime)
	chartList := make([]dy.ProductSalesTrendChart, 0)
	dateChart := make([]int64, 0)
	for {
		if startTime.After(endTime) {
			break
		}
		timestamp := startTime.Unix()
		v := hbaseDataList[startTime.Format("20060102")]
		dateChart = append(dateChart, timestamp)
		chartList = append(chartList, dy.ProductSalesTrendChart{
			DateTimestamp: timestamp,
			Sales:         v.Sales,
			VideoNum:      v.AwemeNum,
		})
		startTime = startTime.AddDate(0, 0, 1)
	}
	res = ProductAwemeTrend{
		DateChart: dateChart,
		ChartList: chartList,
	}
	return
}
