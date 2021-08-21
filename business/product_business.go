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
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

type ProductBusiness struct {
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

func (receiver *ProductBusiness) ProductAuthorAnalysis(productId, keyword, tag string, startTime, endTime time.Time, minFollow, maxFollow int64, scoreType, page, pageSize int) (list []entity.DyProductAuthorAnalysis, total int, comErr global.CommonError) {
	list = []entity.DyProductAuthorAnalysis{}
	esProductBusiness := es.NewEsProductBusiness()
	startRow, stopRow, total, comErr := esProductBusiness.SearchRangeDateRowKey(productId, keyword, startTime, endTime)
	if comErr != nil {
		return
	}
	if startRow.ProductId == "" || stopRow.ProductId == "" {
		return
	}
	startRowKey := startRow.ProductId + "_" + startRow.CreateSdf + "_" + startRow.AuthorId
	stopRowKey := stopRow.ProductId + "_" + stopRow.CreateSdf + "_" + stopRow.AuthorId
	cacheKey := cache.GetCacheKey(cache.ProductAuthorAllList, startRowKey, stopRowKey)
	cacheStr := global.Cache.Get(cacheKey)
	allList := make([]entity.DyProductAuthorAnalysis, 0)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &allList)
	} else {
		allList, _ = hbase.GetProductAuthorAnalysisRange(startRowKey, stopRowKey)
		lastRow, err := hbase.GetProductAuthorAnalysis(stopRowKey)
		if err == nil {
			allList = append(allList, lastRow)
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(allList), 300)
	}
	authorMap := map[string]entity.DyProductAuthorAnalysis{}
	authorIds := make([]string, 0)
	for _, v := range allList {
		if keyword != "" {
			if strings.Index(v.Nickname, keyword) < 0 && v.DisplayId != keyword && v.ShortId != keyword {
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
		if d, ok := authorMap[v.AuthorId]; ok {
			d.Gmv += v.Gmv
			d.Sales += v.Sales
			d.RelatedRooms = append(d.RelatedRooms, v.RelatedRooms...)
			authorMap[v.AuthorId] = d
		} else {
			authorMap[v.AuthorId] = v
			authorIds = append(authorIds, v.AuthorId)
		}
	}
	authorBusiness := NewAuthorBusiness()
	authorDataMap := authorBusiness.GetAuthorByIdsLimitGo(authorIds, 200)
	for _, v := range authorMap {
		if a, ok := authorDataMap[v.AuthorId]; ok {
			v.FollowCount = a.FollowerCount
			if v.DisplayId == "" {
				v.DisplayId = a.UniqueID
				v.ShortId = a.ShortID
			}
		}
		if minFollow > 0 && v.FollowCount < minFollow {
			continue
		}
		if maxFollow > 0 && v.FollowCount >= maxFollow {
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

func (receiver *ProductBusiness) ProductAuthorLiveRooms(productId, authorId string, startTime, endTime time.Time, sortStr, orderBy string, page, pageSize int) (list []entity.DyProductAuthorRelatedRoom, total int) {
	esProductBusiness := es.NewEsProductBusiness()
	allList, _, _ := esProductBusiness.SearchRangeDateList(productId, authorId, startTime, endTime, 1, 1000)
	list = []entity.DyProductAuthorRelatedRoom{}
	for _, v := range allList {
		rowKey := v.ProductId + "_" + v.CreateSdf + "_" + v.AuthorId
		data, err := hbase.GetProductAuthorAnalysis(rowKey)
		if err == nil {
			list = append(list, data.RelatedRooms...)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		switch sortStr {
		case "gmv":
			if orderBy == "desc" {
				return list[i].Gmv > list[j].Gmv
			} else {
				return list[j].Gmv > list[i].Gmv
			}
		case "sale":
			if orderBy == "desc" {
				return list[i].Sales > list[j].Sales
			} else {
				return list[j].Sales > list[i].Sales
			}
		default:
			if orderBy == "desc" {
				return list[i].StartTs > list[j].StartTs
			} else {
				return list[j].StartTs > list[i].StartTs
			}
		}
	})
	total = len(list)
	start := (page - 1) * pageSize
	end := start + pageSize
	if total < end {
		end = total
	}
	list = list[start:end]
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
	esProductBusiness := es.NewEsProductBusiness()
	startRow, stopRow, _, comErr := esProductBusiness.SearchRangeDateRowKey(productId, keyword, startTime, endTime)
	if startRow.ProductId == "" || stopRow.ProductId == "" {
		return
	}
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
	levelMap := map[int]int{}
	authorMap := map[string]string{}
	for _, v := range allList {
		if _, ok := authorMap[v.AuthorId]; ok {
			continue
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
	for k, v := range tagsMap {
		countList.Tags = append(countList.Tags, dy.DyCate{
			Name: k,
			Num:  v,
		})
	}
	for k, v := range levelMap {
		countList.Level = append(countList.Level, dy.DyIntCate{
			Name: k,
			Num:  v,
		})
	}
	if keyword == "" && (len(countList.Tags) > 0 || len(countList.Level) > 0) {
		countJson := utils.SerializeData(countList)
		_ = global.Cache.Set(cKey, countJson, 300)
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
