package dy

import (
	"dongchamao/business"
	"dongchamao/business/es"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	es2 "dongchamao/models/es"
	dy2 "dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	"fmt"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"math"
	"sort"
	"strings"
	"time"
)

type AuthorController struct {
	controllers.ApiBaseController
}

func (receiver *AuthorController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseLoginMinShowNum, business.DyJewelBaseShowNum)
}

//达人分类
func (receiver *AuthorController) AuthorCate() {
	cateList := business.NewAuthorBusiness().GetCacheAuthorCate(true)
	receiver.SuccReturn(cateList)
	return
}

//达人带货行业
func (receiver *AuthorController) GetCacheAuthorLiveTags() {
	authorBusiness := business.NewAuthorBusiness()
	cateList := authorBusiness.GetCacheAuthorLiveTags(true)
	receiver.SuccReturn(map[string]interface{}{
		"list": cateList,
	})
	return
}

//达人库
func (receiver *AuthorController) BaseSearch() {
	keyword := receiver.GetString("keyword", "")
	category := receiver.GetString("category", "")
	secondCategory := receiver.GetString("second_category", "")
	sellTags := receiver.GetString("sell_tags", "")
	province := receiver.GetString("province", "")
	city := receiver.GetString("city", "")
	fanProvince := receiver.GetString("fan_province", "")
	fanCity := receiver.GetString("fan_city", "")
	sortStr := receiver.GetString("sort", "follower_incre_count")
	orderBy := receiver.GetString("order_by", "desc")
	minFollower, _ := receiver.GetInt64("min_follower", 0)
	maxFollower, _ := receiver.GetInt64("max_follower", 0)
	minWatch, _ := receiver.GetInt64("min_watch", 0)
	maxWatch, _ := receiver.GetInt64("max_watch", 0)
	minDigg, _ := receiver.GetInt64("min_digg", 0)
	maxDigg, _ := receiver.GetInt64("max_digg", 0)
	minGmv, _ := receiver.GetInt64("min_gmv", 0)
	maxGmv, _ := receiver.GetInt64("max_gmv", 0)
	minAge, _ := receiver.GetInt("min_age", 0)
	maxAge, _ := receiver.GetInt("max_age", 0)
	fanAge := receiver.GetString("fan_age", "")
	gender, _ := receiver.GetInt("gender", 0)
	fanGender, _ := receiver.GetInt("fan_gender", 0)
	verification, _ := receiver.GetInt("verification", 0)
	level, _ := receiver.GetInt("level", 0)
	isBrand, _ := receiver.GetInt("is_brand", 0)
	isDelivery, _ := receiver.GetInt("is_delivery", 0)
	superSeller, _ := receiver.GetInt("super_seller", 0)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	pageSize = receiver.CheckPageSize(pageSize)
	receiver.KeywordBan(keyword)
	if !receiver.HasLogin && keyword != "" {
		receiver.FailReturn(global.NewError(4001))
		return
	}
	if !receiver.HasAuth {
		if category != "" || secondCategory != "" || sellTags != "" || province != "" || city != "" || fanProvince != "" || fanCity != "" || sortStr != "follower_incre_count" || orderBy != "desc" ||
			minFollower > 0 || maxFollower > 0 || minWatch > 0 || maxWatch > 0 || minDigg > 0 || maxDigg > 0 || minGmv > 0 || maxGmv > 0 ||
			gender > 0 || minAge > 0 || maxAge > 0 || fanAge != "" || verification > 0 || level > 0 || fanGender > 0 ||
			superSeller == 1 || isDelivery > 0 || isBrand == 1 || page != 1 {
			if !receiver.HasLogin {
				receiver.FailReturn(global.NewError(4001))
				return
			}
			receiver.FailReturn(global.NewError(4004))
			return
		}
		if pageSize > receiver.MaxTotal {
			pageSize = receiver.MaxTotal
		}
	}
	formNum := (page - 1) * pageSize
	if formNum > receiver.MaxTotal {
		receiver.FailReturn(global.NewError(4004))
		return
	}
	authorId := ""
	if strings.Index(keyword, "http://") > 0 || strings.Index(keyword, "https://") > 0 {
		keyword = strings.Replace(keyword, "在抖音，记录美好生活！ ", "", 1)
	}
	if utils.CheckType(keyword, "url") {
		spiderBusiness := business.SpiderBusiness{}
		shortUrl, _ := business.ParseDyShortUrlToSecUid(keyword)
		if shortUrl == "" {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		//authorId = utils.ParseDyAuthorUrl(shortUrl) // 获取authorId  抖音更改版本 获取不到sec_uid
		secUid := utils.ParseDyAuthorSecUrl(shortUrl)
		author := spiderBusiness.GetAuthorBaseInfoV2(secUid)
		if author != nil {
			authorId = author.AuthorId
		} else {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		keyword = ""
	} else {
		keyword = utils.MatchDouyinNewText(keyword)
	}
	EsAuthorBusiness := es.NewEsAuthorBusiness()
	//只带keyword去查
	preTotal := 0
	var comErr global.CommonError
	if page == 1 {
		_, preTotal, comErr = EsAuthorBusiness.BaseSearch(authorId, keyword, "", "", "", "", "", "", "", "",
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, page, pageSize,
			sortStr, orderBy)
		if comErr != nil {
			receiver.FailReturn(comErr)
			return
		}
	}
	list, total, comErr := EsAuthorBusiness.BaseSearch(authorId, keyword, category, secondCategory, sellTags, province, city, fanProvince, fanCity, fanAge,
		minFollower, maxFollower, minWatch, maxWatch, minDigg, maxDigg, minGmv, maxGmv,
		gender, minAge, maxAge, verification, level, isDelivery, isBrand, superSeller, fanGender, page, pageSize,
		sortStr, orderBy)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	if keyword != "" && preTotal == 0 && page == 1 {
		spiderBusiness := business.NewSpiderBusiness()
		authorIncomeRawList, err1 := spiderBusiness.GetAuthorListByKeyword(keyword)
		if err1 != nil {
			receiver.FailReturn(global.NewMsgError(err1.Error()))
			return
		}
		list := make([]es2.DyAuthor, 0)
		for _, v := range authorIncomeRawList {
			var tempAuthor es2.DyAuthor
			tempAuthor.AuthorId = business.IdEncrypt(v.Id)
			tempAuthor.UniqueId = v.UniqueId
			if tempAuthor.UniqueId == "0" || tempAuthor.UniqueId == "" {
				tempAuthor.UniqueId = v.ShortId
			}
			tempAuthor.Avatar = dyimg.Fix(v.Avatar)
			tempAuthor.FollowerCount = v.FollowerCount
			tempAuthor.ShortId = v.ShortId
			tempAuthor.Gender = v.Gender
			tempAuthor.Nickname = v.Nickname
			tempAuthor.Birthday = v.Birthday
			tempAuthor.VerifyName = v.VerifyName
			tempAuthor.VerificationType = v.VerificationType
			tempAuthor.IsCollection = 0
			list = append(list, tempAuthor)
		}
		end := 10
		if len(list) < end {
			end = len(list)
		}
		resList := list[0:end]
		listTotal := len(resList)
		receiver.SuccReturn(map[string]interface{}{
			"list":       resList,
			"total":      listTotal,
			"total_page": 1,
			"max_num":    listTotal,
			"has_auth":   receiver.HasAuth,
			"has_login":  receiver.HasLogin,
			"data_from":  "douyin",
		})
		return
	}
	authorIds := make([]string, 0)
	for _, v := range list {
		authorIds = append(authorIds, v.AuthorId)
	}
	if receiver.HasLogin {
		collectBusiness := business.NewCollectBusiness()
		collect, comErr := collectBusiness.DyListCollect(1, receiver.UserId, authorIds)
		if comErr != nil {
			receiver.FailReturn(comErr)
		}
		for k, v := range list {
			list[k].IsCollect = collect[v.AuthorId]
		}
	}
	authorLiveMap := map[string]string{}
	cacheKey := cache.GetCacheKey(cache.AuthorLiveMap, utils.Md5_encode(strings.Join(authorIds, ",")))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &authorLiveMap)
	} else {
		authorMap, _ := hbase.GetAuthorByIds(authorIds)
		for k, v := range authorMap {
			if v.RoomStatus == 2 {
				authorLiveMap[k] = v.RoomId
			} else {
				authorLiveMap[k] = ""
			}
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(authorLiveMap), 180)
	}
	for k, v := range list {
		list[k].Avatar = dyimg.Fix(v.Avatar)
		if v.UniqueId == "" || v.UniqueId == "0" {
			list[k].UniqueId = v.ShortId
		}
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		if a, ok := authorLiveMap[v.AuthorId]; ok {
			list[k].RoomId = business.IdEncrypt(a)
		} else {
			authorData, _ := hbase.GetAuthor(v.AuthorId)
			list[k].RoomId = business.IdEncrypt(authorData.RoomId)
		}
		list[k].IsCollection = 1
		list[k].DiggFollowerRate = utils.RateMin(list[k].DiggFollowerRate)
		list[k].InteractionRate = utils.RateMin(list[k].InteractionRate)
	}
	totalPage := math.Ceil(float64(total) / float64(pageSize))
	maxPage := math.Ceil(float64(receiver.MaxTotal) / float64(pageSize))
	if totalPage > maxPage {
		totalPage = maxPage
	}
	maxTotal := receiver.MaxTotal
	if maxTotal > total {
		maxTotal = total
	}
	business.NewUserBusiness().KeywordsRecord(keyword)
	receiver.SuccReturn(map[string]interface{}{
		"list":       list,
		"total":      total,
		"total_page": totalPage,
		"max_num":    maxTotal,
		"has_auth":   receiver.HasAuth,
		"has_login":  receiver.HasLogin,
		"data_from":  "local",
	})
	return
}

//达人数据
func (receiver *AuthorController) AuthorBaseData() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	authorBase, comErr := authorBusiness.HbaseGetAuthor(authorId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	reputation, _ := authorBusiness.HbaseGetAuthorReputation(authorId)
	fansClub, _ := hbase.GetAuthorFansClub(authorId)
	//todo 昨天数据取前天
	basicBefore, _ := hbase.GetAuthorBasic(authorId, time.Now().AddDate(0, 0, -1).Format("20060102"))
	yseBasicBefore, _ := hbase.GetAuthorBasic(authorId, time.Now().AddDate(0, 0, -2).Format("20060102"))
	authorBase.Data.ID = business.IdEncrypt(authorBase.Data.ID)
	authorBase.Data.RoomID = business.IdEncrypt(authorBase.Data.RoomID)
	//获取榜单排名
	mapRank := authorBusiness.HbaseGetAuthorRank(authorId)
	mapRank["desc"] = fmt.Sprintf("达人%s%s%s名", mapRank["rank_name"], mapRank["date_type"], mapRank["value"])
	basic := entity.DyAuthorBasic{
		FollowerCount:        basicBefore.FollowerCount,
		FollowerCountBefore:  yseBasicBefore.FollowerCount,
		TotalFansCount:       basicBefore.TotalFansCount,
		TotalFansCountBefore: yseBasicBefore.TotalFansCount,
		TotalFavorited:       basicBefore.TotalFavorited,
		TotalFavoritedBefore: yseBasicBefore.TotalFavorited,
		CommentCount:         basicBefore.CommentCount,
		CommentCountBefore:   yseBasicBefore.CommentCount,
		ForwardCount:         basicBefore.ForwardCount,
		ForwardCountBefore:   yseBasicBefore.ForwardCount,
	}
	authorStore, _ := hbase.GetAuthorStore(authorId)
	returnMap := map[string]interface{}{
		"author_base": authorBase.Data,
		"room_count":  authorBase.LiveCount,
		"reputation": dy2.RepostSimpleReputation{
			Score:         reputation.Score,
			Level:         reputation.Level,
			EncryptShopID: reputation.EncryptShopID,
			ShopName:      reputation.ShopName,
			ShopLogo:      reputation.ShopLogo,
		},
		"fans_club":   fansClub.TotalFansCount,
		"rank":        nil,
		"tags":        authorBase.Tags,
		"second_tags": authorBase.TagsLevelTwo,
		"basic":       basic,
		"shop": dy2.DyAuthorStoreSimple{
			ShopId:   business.IdEncrypt(authorStore.Id),
			ShopName: authorStore.ShopName,
		},
		"rank_infor": mapRank,
	}
	receiver.SuccReturn(returnMap)
	return
}

func (receiver *AuthorController) AuthorViewData() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	esLiveBusiness := es.NewEsLiveBusiness()
	authorBase, comErr := authorBusiness.HbaseGetAuthor(authorId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	todayString := time.Now().Format("20060102")
	monthString := time.Now().Format("200601")
	todayTime, _ := time.ParseInLocation("20060102", todayString, time.Local)
	monthTime, _ := time.ParseInLocation("20060102", monthString+"01", time.Local)
	//lastMonthDay := todayTime.AddDate(0, 0, -29)
	nowWeek := int(time.Now().Weekday())
	if nowWeek == 0 {
		nowWeek = 7
	}
	lastWeekDay := todayTime.AddDate(0, 0, -(nowWeek - 1))
	monthRoom := esLiveBusiness.CountDataByAuthor(authorId, monthTime, todayTime)
	weekRoom := esLiveBusiness.CountDataByAuthor(authorId, lastWeekDay, todayTime)
	productCount := dy2.DyAuthorBaseProductCount{}
	startTime := todayTime.AddDate(0, 0, -31)
	yesterday := todayTime.AddDate(0, 0, -1)
	//room30Count := esLiveBusiness.CountDataByAuthor(authorId, startTime, yesterday)
	cacheKey := cache.GetCacheKey(cache.AuthorViewProductAllList, authorId, startTime.Format("20060102"), yesterday.Format("20060102"))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &productCount)
	} else {
		productMap := map[string]string{}
		liveList, _, _ := es.NewEsLiveBusiness().ScanLiveProductByAuthor(authorId, "", "", "", "", "", "", 0, startTime, yesterday, 1, 10000)
		awemeList, _, _ := es.NewEsVideoBusiness().ScanAwemeProductByAuthor(authorId, "", "", "", "", "", "", 0, startTime, yesterday, 1, 10000)
		brandSaleMap := map[string]float64{}
		brandNumMap := map[string]int{}
		priceNumMap := map[string]int{}
		priceSaleMap := map[string]float64{}
		var totalGmv float64
		var totalSales float64
		for _, v := range liveList {
			if _, exist := productMap[v.ProductID]; !exist {
				productMap[v.ProductID] = v.ProductID
			}
			totalGmv += v.PredictGmv
			totalSales += math.Floor(v.PredictSales)
			category := v.DcmLevelFirst
			if category == "" || category == "null" {
				category = "其他"
			}
			var priceStr string
			if v.Price > 500 {
				priceStr = "500+"
			} else if v.Price > 300 {
				priceStr = "300-500"
			} else if v.Price > 100 {
				priceStr = "100-300"
			} else if v.Price > 50 {
				priceStr = "50-100"
			} else {
				priceStr = "0-50"
			}
			if _, ok := brandSaleMap[category]; !ok {
				brandSaleMap[category] = math.Floor(v.PredictSales)
			} else {
				brandSaleMap[category] += math.Floor(v.PredictSales)
			}
			if _, ok := brandNumMap[category]; !ok {
				brandNumMap[category] = 1
			} else {
				brandNumMap[category] += 1
			}
			if _, ok := priceSaleMap[priceStr]; !ok {
				priceSaleMap[priceStr] = math.Floor(v.PredictSales)
			} else {
				priceSaleMap[priceStr] += math.Floor(v.PredictSales)
			}
			if _, ok := priceNumMap[priceStr]; !ok {
				priceNumMap[priceStr] = 1
			} else {
				priceNumMap[priceStr] += 1
			}
		}
		for _, v := range awemeList {
			if _, exist := productMap[v.ProductId]; !exist {
				productMap[v.ProductId] = v.ProductId
			}
			totalGmv += v.Gmv
			totalSales += float64(v.Sales)
			category := v.DcmLevelFirst
			if category == "" || category == "null" {
				category = "其他"
			}
			var priceStr string
			if v.Price > 500 {
				priceStr = "500+"
			} else if v.Price > 300 {
				priceStr = "300-500"
			} else if v.Price > 100 {
				priceStr = "100-300"
			} else if v.Price > 50 {
				priceStr = "50-100"
			} else {
				priceStr = "0-50"
			}
			if _, ok := brandSaleMap[category]; !ok {
				brandSaleMap[category] = float64(v.Sales)
			} else {
				brandSaleMap[category] += float64(v.Sales)
			}
			if _, ok := brandNumMap[category]; !ok {
				brandNumMap[category] = 1
			} else {
				brandNumMap[category] += 1
			}
			if _, ok := priceSaleMap[priceStr]; !ok {
				priceSaleMap[priceStr] = float64(v.Sales)
			} else {
				priceSaleMap[priceStr] += float64(v.Sales)
			}
			if _, ok := priceNumMap[priceStr]; !ok {
				priceNumMap[priceStr] = 1
			} else {
				priceNumMap[priceStr] += 1
			}
		}
		brandSaleList := make([]dy2.NameValueInt64Chart, 0)
		brandNumList := make([]dy2.NameValueChart, 0)
		for c, v := range brandSaleMap {
			brandSaleList = append(brandSaleList, dy2.NameValueInt64Chart{
				Name:  c,
				Value: utils.ToInt64(v),
			})
		}
		for c, v := range brandNumMap {
			brandNumList = append(brandNumList, dy2.NameValueChart{
				Name:  c,
				Value: v,
			})
		}
		sort.Slice(brandSaleList, func(i, j int) bool {
			return brandSaleList[j].Value < brandSaleList[i].Value
		})
		sort.Slice(brandNumList, func(i, j int) bool {
			return brandNumList[j].Value < brandNumList[i].Value
		})
		listLen := len(brandSaleList)
		topBrandSaleList := make([]dy2.NameValueInt64Chart, 0)
		topBrandNumList := make([]dy2.NameValueChart, 0)
		topSaleCates := make([]string, 0)
		topNumCates := make([]string, 0)
		if listLen > 0 {
			salesStopKey := 0
			numtopKey := 0
			var sale int64
			var num int
			for _, v := range brandSaleList {
				if v.Value == 0 {
					break
				}
				if v.Name == "其他" {
					sale += v.Value
				} else {
					if salesStopKey < 3 {
						topBrandSaleList = append(topBrandSaleList, v)
						topSaleCates = append(topSaleCates, v.Name)
						salesStopKey++
					} else {
						sale += v.Value
					}
				}
			}
			if sale > 0 {
				topBrandSaleList = append(topBrandSaleList, dy2.NameValueInt64Chart{
					Name:  "其他",
					Value: sale,
				})
			}
			for _, v := range brandNumList {
				if v.Value == 0 {
					break
				}
				if v.Name == "其他" {
					num += v.Value
				} else {
					if numtopKey < 3 {
						topBrandNumList = append(topBrandNumList, v)
						topNumCates = append(topNumCates, v.Name)
						numtopKey++
					} else {
						num += v.Value
					}
				}
			}
			if num > 0 {
				topBrandNumList = append(topBrandNumList, dy2.NameValueChart{
					Name:  "其他",
					Value: num,
				})
			}
		}
		productCount = dy2.DyAuthorBaseProductCount{
			Sales30Top3:           topSaleCates,
			ProductNum30Top3:      topNumCates,
			Sales30Top3Chart:      topBrandSaleList,
			ProductNum30Top3Chart: topBrandNumList,
			ProductNum:            len(productMap),
			//Predict30Sales:        totalSales,
			//Predict30Gmv:          totalGmv,
			Sales30Chart: []dy2.DyAuthorBaseProductPriceChart{},
		}
		for p, v := range priceSaleMap {
			num := priceNumMap[p]
			productCount.Sales30Chart = append(productCount.Sales30Chart, dy2.DyAuthorBaseProductPriceChart{
				Name:       p,
				Sales:      utils.ToInt64(v),
				ProductNum: num,
			})
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(productCount), 600)

	}
	//productCount.ProductNum = int(esLiveBusiness.CountRoomProductByAuthorId(authorId, startTime, yesterday))
	videoSumData := es.NewEsVideoBusiness().SumDataByAuthor(authorId, startTime, yesterday)
	liveSumData, room30Count := esLiveBusiness.SumDataByAuthor(authorId, startTime, yesterday)
	dayLiveRoomNum := esLiveBusiness.CountRoomByDayByAuthorId(authorId, 1, startTime, yesterday)
	var avgGmv float64 = 0
	var avgSales float64 = 0
	if dayLiveRoomNum > 0 {
		avgSales = math.Floor(liveSumData.TotalSales.Sum / float64(dayLiveRoomNum))
		avgGmv = math.Floor(liveSumData.TotalGmv.Sum / float64(dayLiveRoomNum))
	}
	productCount.Predict30Gmv = liveSumData.TotalGmv.Sum + videoSumData.Gmv
	productCount.Predict30Sales = utils.ToInt64(math.Floor(liveSumData.TotalSales.Sum)) + videoSumData.Sales
	data := dy2.DyAuthorBaseCount{
		LiveCount: dy2.DyAuthorBaseLiveCount{
			RoomCount:      int64(authorBase.LiveCount),
			Room30Count:    int64(room30Count),
			Predict30Sales: avgSales,
			Predict30Gmv:   utils.FriendlyFloat64(avgGmv),
			AgeDuration:    authorBase.AgeLiveDuration,
			WeekRoomCount:  weekRoom,
			MonthRoomCount: monthRoom,
		},
		VideoCount: dy2.DyAuthorBaseVideoCount{
			VideoCount:       int64(authorBase.AwemeCount),
			Video30Count:     int64(videoSumData.Total),
			DiggFollowerRate: 0,
			Predict30Sales:   float64(videoSumData.Sales),
			Predict30Gmv:     videoSumData.Gmv,
			AgeDuration:      authorBase.Duration / 1000,
		},
		ProductCount: productCount,
	}
	if authorBase.AwemeCount != 0 {
		data.VideoCount.AvgDigg = authorBase.DiggCount / int64(authorBase.AwemeCount)
	}
	// todo 达人近30天视频为0 屏蔽赞粉比
	if data.VideoCount.Video30Count > 0 {
		data.VideoCount.DiggFollowerRate = authorBase.DiggFollowerRate
	}
	data.VideoCount.Avg30Digg = videoSumData.AvgDigg
	firstLiveTimestamp := authorBase.FirstLiveTime - (authorBase.FirstLiveTime % 86400)
	firstVideoTimestamp := authorBase.FirstAwemeTime - (authorBase.FirstAwemeTime % 86400)
	if firstLiveTimestamp > 0 {
		firstLiveTime := time.Unix(firstLiveTimestamp, 0)
		tmpWeek := int(firstLiveTime.Weekday())
		if tmpWeek == 0 {
			tmpWeek = 7
		}
		days := (todayTime.AddDate(0, 0, 7-nowWeek).Unix() - firstLiveTime.AddDate(0, 0, -(tmpWeek-1)).Unix()) / 86400
		days += 1
		weekNum := utils.ToInt64(days / 7)
		if weekNum > 0 {
			data.LiveCount.AvgWeekRoomCount = utils.ToInt64(math.Ceil(float64(authorBase.LiveCount) / float64(weekNum)))
		}
		var month = utils.ToInt64(math.Ceil(float64(time.Now().Unix()-firstLiveTime.Unix()) / (30 * 86400)))
		if month > 0 {
			data.LiveCount.AvgMonthRoomCount = utils.ToInt64(math.Ceil(float64(authorBase.LiveCount) / float64(month)))
		}
	}
	if firstVideoTimestamp > 0 {
		firstVideoTime := time.Unix(firstVideoTimestamp, 0)
		tmpWeek := int(firstVideoTime.Weekday())
		if tmpWeek == 0 {
			tmpWeek = 7
		}
		days := (todayTime.AddDate(0, 0, 7-nowWeek).Unix() - firstVideoTime.AddDate(0, 0, -(tmpWeek-1)).Unix()) / 86400
		days += 1
		weekNum := utils.ToInt64(days / 7)
		if weekNum > 0 {
			data.VideoCount.WeekVideoCount = utils.ToInt64(math.Ceil(float64(authorBase.AwemeCount) / float64(weekNum)))
		}
		var month = utils.ToInt64(math.Ceil(float64(time.Now().Unix()-firstVideoTime.Unix()) / (30 * 86400)))
		if month > 0 {
			data.VideoCount.MonthVideoCount = utils.ToInt64(math.Ceil(float64(authorBase.AwemeCount) / float64(month)))
		}
	}
	receiver.SuccReturn(data)
	return
}

//星图指数数据
func (receiver *AuthorController) AuthorStarSimpleData() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	returnMap := map[string]interface{}{
		"has_star_detail": false,
		"star_detail":     nil,
	}
	authorBusiness := business.NewAuthorBusiness()
	xtDetail, comErr := hbase.GetXtAuthorDetail(authorId)
	if comErr == nil {
		returnMap["has_star_detail"] = true
		returnMap["star_detail"] = authorBusiness.GetDyAuthorScore(xtDetail.LiveScore, xtDetail.Score)
	}
	receiver.SuccReturn(returnMap)
	return
}

//达人口碑
func (receiver *AuthorController) Reputation() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	reputation, _ := authorBusiness.HbaseGetAuthorReputation(authorId)
	reputation.UID = business.IdEncrypt(reputation.UID)
	lenNum := len(reputation.DtScoreList)
	if lenNum > 30 {
		start := lenNum - 30
		reputation.DtScoreList = reputation.DtScoreList[start:]
	}
	receiver.SuccReturn(map[string]interface{}{
		"reputation": reputation,
	})
	return
}

//达人视频概览
func (receiver *AuthorController) AuthorAwemesByDay() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	startDay, endDay, comErr := receiver.GetRangeDate()
	if !receiver.HasAuth && startDay.Sub(endDay).Hours()/24 > 30 {
		receiver.FailReturn(global.NewError(4004))
	}
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	aABusiness := business.NewAuthorAwemeBusiness()
	videoOverview := aABusiness.GetVideoAggRangeDate(authorId, startDay, endDay)
	receiver.SuccReturn(map[string]interface{}{
		"video_overview": videoOverview,
	})
	return
}

//达人视频列表
func (receiver *AuthorController) AuthorAwemes() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	hasProduct, _ := receiver.GetInt("has_product", 0)
	keyword := receiver.GetString("keyword", "")
	sortStr := receiver.GetString("sort", "")
	orderBy := receiver.GetString("order_by", "")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 30)
	list, total, totalSales, totalGmv, comErr := es.NewEsVideoBusiness().SearchByAuthor(authorId, keyword, sortStr, orderBy, hasProduct, page, pageSize, startTime, endTime)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		list[k].AwemeId = business.IdEncrypt(v.AwemeId)
		list[k].AwemeCover = dyimg.Fix(v.AwemeCover)
		list[k].Avatar = dyimg.Fix(v.Avatar)
		if v.UniqueId == "" || v.UniqueId == "0" {
			list[k].UniqueId = v.ShortId
		}
		list[k].AwemeUrl = business.AwemeUrl + v.AwemeId
	}
	maxTotal := total
	if total > business.EsMaxShowNum {
		maxTotal = business.EsMaxShowNum
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":           list,
		"total":          total,
		"total_sales":    totalSales,
		"total_gmv":      totalGmv,
		"max_show_total": maxTotal,
	})
	return
}

//基础数据趋势图
func (receiver *AuthorController) AuthorBasicChart() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	startDay := receiver.Ctx.Input.Param(":start")
	endDay := receiver.Ctx.Input.Param(":end")
	if authorId == "" || startDay == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if endDay == "" {
		endDay = time.Now().Format("2006-01-02")
	}
	pslTime := "2006-01-02"
	t1, err := time.ParseInLocation(pslTime, startDay, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	t2, err := time.ParseInLocation(pslTime, endDay, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if t1.After(t2) || t2.After(t1.AddDate(0, 0, 90)) || t2.After(time.Now()) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	data, comErr := authorBusiness.HbaseGetAuthorBasicRangeDate(authorId, t1, t2)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(data)
	return
}

//达人粉丝团
func (receiver *AuthorController) AuthorFansList() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	authorLiveFansClubUser, _ := hbase.GetAuthorLiveFansClubUser(authorId)
	fans := authorLiveFansClubUser.FansInfos
	if len(fans) == 0 {
		fans = []entity.DyLiveFansClubUserInfo{}
	}
	sort.Slice(fans, func(i, j int) bool {
		return fans[i].Intimacy > fans[j].Intimacy
	})
	for k, v := range fans {
		fans[k].Id = business.IdEncrypt(v.Id)
		fans[k].Avatar = dyimg.Fix(v.Avatar)
	}
	receiver.SuccReturn(map[string]interface{}{
		"list": fans,
	})
	return
}

//粉丝分布分析
func (receiver *AuthorController) AuthorFansAnalyse() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	authorData, comErr := hbase.GetAuthor(authorId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	detail, comErr := hbase.GetXtAuthorDetail(authorId)
	data := map[string][]entity.XtDistributionsList{}
	var countCity int64 = 0
	var countPro int64 = 0
	if comErr == nil {
		for _, v := range detail.Distributions {
			name := ""
			switch v.Type {
			case entity.XtGenderDistribution:
				name = "gender"
			case entity.XtCityDistribution:
				name = "city"
			case entity.XtAgeDistribution:
				name = "age"
			case entity.XtProvinceDistribution:
				name = "province"
			default:
				continue
			}
			data[name] = v.DistributionList
		}
		for _, v := range data["city"] {
			countCity += v.DistributionValue
		}
		for _, v := range data["province"] {
			countPro += v.DistributionValue
		}

	} else {
		data["gender"] = []entity.XtDistributionsList{}
		data["city"] = []entity.XtDistributionsList{}
		data["age"] = []entity.XtDistributionsList{}
		data["province"] = []entity.XtDistributionsList{}
		//性别处理
		for _, v := range authorData.Gender {
			DistributionKey := ""
			if v.Gender == "男" {
				DistributionKey = "male"
			} else if v.Gender == "女" {
				DistributionKey = "female"
			}
			if DistributionKey == "" {
				continue
			}
			data["gender"] = append(data["gender"], entity.XtDistributionsList{
				DistributionKey:   DistributionKey,
				DistributionValue: utils.ToInt64(v.GenderNum),
			})
		}
		for _, v := range authorData.City {
			value := utils.ToInt64(v.CityNum)
			countCity += value
			data["city"] = append(data["city"], entity.XtDistributionsList{
				DistributionKey:   v.City,
				DistributionValue: value,
			})
		}
		ageMap := map[string]int64{}
		for _, v := range authorData.AgeDistrinbution {
			if _, exist := ageMap[v.AgeDistrinbution]; !exist {
				ageMap[v.AgeDistrinbution] = utils.ToInt64(v.AgeDistrinbutionNum)
			} else {
				ageMap[v.AgeDistrinbution] += utils.ToInt64(v.AgeDistrinbutionNum)
			}
		}
		for k, v := range ageMap {
			data["age"] = append(data["age"], entity.XtDistributionsList{
				DistributionKey:   k,
				DistributionValue: v,
			})
		}
		for _, v := range authorData.Province {
			value := utils.ToInt64(v.ProvinceNum)
			countPro += value
			data["province"] = append(data["province"], entity.XtDistributionsList{
				DistributionKey:   v.Province,
				DistributionValue: value,
			})
		}

	}
	if countCity > 0 {
		for k, v := range data["city"] {
			data["city"][k].DistributionPer = float64(v.DistributionValue) / float64(countCity)
		}
	}
	if countPro > 0 {
		for k, v := range data["province"] {
			data["province"][k].DistributionPer = float64(v.DistributionValue) / float64(countPro)
		}
	}
	var countHour int64 = 0
	var countWeek int64 = 0
	data["active_day"] = []entity.XtDistributionsList{}
	for _, v := range authorData.HourCreateTimeNum {
		value := utils.ToInt64(v.HourCreateTimeNum)
		countHour += value
		data["active_day"] = append(data["active_day"], entity.XtDistributionsList{
			DistributionKey:   v.HourCreateTime,
			DistributionValue: value,
		})
	}
	data["active_week"] = []entity.XtDistributionsList{}
	for _, v := range authorData.WeekCreateTimeNum {
		value := utils.ToInt64(v.WeekCreateTimeNum)
		countWeek += value
		data["active_week"] = append(data["active_week"], entity.XtDistributionsList{
			DistributionKey:   v.WeekCreateTime,
			DistributionValue: utils.ToInt64(v.WeekCreateTimeNum),
		})
	}
	if countHour > 0 {
		for k, v := range data["active_day"] {
			data["active_day"][k].DistributionPer = float64(v.DistributionValue) / float64(countHour)
		}
	}
	if countWeek > 0 {
		for k, v := range data["active_week"] {
			data["active_week"][k].DistributionPer = float64(v.DistributionValue) / float64(countWeek)
		}
	}
	receiver.SuccReturn(data)
	return
}

//达人直播分析
func (receiver *AuthorController) CountLiveRoomAnalyse() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	t1, t2, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	data := authorBusiness.CountLiveRoomAnalyse(authorId, t1, t2)
	receiver.SuccReturn(data)
	return
}

//达人直播间列表
func (receiver *AuthorController) AuthorLiveRooms() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	InputData := receiver.InputFormat()
	keyword := InputData.GetString("keyword", "")
	sortStr := InputData.GetString("sort", "create_time")
	orderBy := InputData.GetString("order_by", "desc")
	page := InputData.GetInt("page", 1)
	size := InputData.GetInt("page_size", 10)
	listType, _ := receiver.GetInt("list_type", 0)
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	t1, t2, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	esLiveBusiness := es.NewEsLiveBusiness()
	list, total, totalSales, totalGmv, comErr := esLiveBusiness.SearchAuthorRooms(authorId, keyword, sortStr, orderBy, page, size, t1, t2)
	if listType == 1 {
		roomIds := []string{}
		for _, v := range list {
			roomIds = append(roomIds, v.RoomId)
		}
		liveMap, _ := hbase.GetLiveInfoByIds(roomIds)
		for k, v := range list {
			if liveInfo, exist := liveMap[v.RoomId]; exist {
				list[k].RoomStatus = liveInfo.RoomStatus
				if liveInfo.RoomStatus == 4 {
					list[k].FinishTime = liveInfo.FinishTime
				} else {
					list[k].FinishTime = time.Now().Unix()
				}
			}
		}
	}
	for k, v := range list {
		list[k].RoomId = business.IdEncrypt(v.RoomId)
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		list[k].PredictSales = math.Floor(v.PredictSales)
		list[k].PredictGmv = math.Floor(v.PredictGmv)
		list[k].Cover = dyimg.Fix(v.Cover)
		list[k].Avatar = dyimg.Fix(v.Avatar)
	}
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":        list,
		"total":       total,
		"total_sales": totalSales,
		"total_gmv":   totalGmv,
	})
	return
}

//达人电商分析
func (receiver *AuthorController) AuthorProductAnalyse() {
	authorId := business.IdDecrypt(receiver.GetString(":author_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	firstCate := receiver.GetString("first_cate", "")
	secondCate := receiver.GetString("second_cate", "")
	thirdCate := receiver.GetString("third_cate", "")
	brandName := receiver.GetString("brand_name", "")
	sortStr := receiver.GetString("sort", "")
	orderBy := receiver.GetString("order_by", "")
	shopType, _ := receiver.GetInt("shop_type", 0)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	if brandName == "不限" {
		brandName = ""
	}
	if firstCate == "不限" {
		firstCate = ""
	}
	authorBusiness := business.NewAuthorBusiness()
	list, analysisCount, cateList, brandList, total, comErr := authorBusiness.NewGetAuthorProductAnalyse(authorId, keyword, firstCate, secondCate, thirdCate, brandName, sortStr, orderBy, shopType, startTime, endTime, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":           list,
		"cate_list":      cateList,
		"brand_list":     brandList,
		"analysis_count": analysisCount,
		"total":          total,
	})
	return
}

//达人合作小店
func (receiver *AuthorController) AuthorShopAnalyse() {
	authorId := business.IdDecrypt(receiver.GetString(":author_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	sortStr := receiver.GetString("sort", "")
	orderBy := receiver.GetString("order_by", "")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	authorBusiness := business.NewAuthorBusiness()
	list, total, comErr := authorBusiness.GetAuthorShopAnalyse(authorId, keyword, sortStr, orderBy, startTime, endTime, page, pageSize, receiver.UserId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
	return
}

//达人商品直播间
func (receiver *AuthorController) AuthorProductRooms() {
	authorId := business.IdDecrypt(receiver.GetString(":author_id", ""))
	productId := business.IdDecrypt(receiver.GetString(":product_id", ""))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 5, 10)
	authorBusiness := business.NewAuthorBusiness()
	list, total, comErr := authorBusiness.GetAuthorProductRooms(authorId, productId, startTime, endTime, page, pageSize, "shelf_time", "desc")
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
	return
}

//达人商品直播间统计
func (receiver *AuthorController) SumAuthorProductOfRooms() {
	authorId := business.IdDecrypt(receiver.GetString(":author_id", ""))
	productId := business.IdDecrypt(receiver.GetString(":product_id", ""))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	totalGmv, totalSales, comErr := es.NewEsLiveBusiness().SumAuthorProductOfRoom(authorId, productId, startTime, endTime)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	liveNum, _ := es.NewEsLiveBusiness().SumAuthorProductCountRoom(authorId, productId, startTime, endTime)
	receiver.SuccReturn(map[string]interface{}{
		"gmv":      totalGmv,
		"sales":    totalSales,
		"live_num": liveNum,
	})
	return
}

//达人收录 搜索
func (receiver *AuthorController) AuthorIncomeSearch() {
	var secUid string
	keyword := receiver.GetString("keyword", "")
	if keyword == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	spiderBusiness := business.NewSpiderBusiness()
	if utils.CheckType(keyword, "url") { // 抓换链接
		shortUrl, _ := business.ParseDyShortUrlToSecUid(keyword)
		if shortUrl == "" {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		secUid = utils.ParseDyAuthorSecUrl(shortUrl) // 获取authorId
		// 请求或去 抖音uid
		authorIncome := spiderBusiness.GetAuthorBaseInfoV2(secUid)
		if authorIncome != nil {
			authorIncome.AuthorId = business.IdEncrypt(authorIncome.AuthorId)
			authorIncome.Avatar = dyimg.Fix(authorIncome.Avatar)
			receiver.SuccReturn(authorIncome)
		} else {
			receiver.FailReturn(global.NewError(4000))
		}
		return
	} else {
		// 如果是keyword形式的，先查es，es没有数据就请求爬虫数据接口
		list, total, _ := es.NewEsAuthorBusiness().SimpleSearch(
			"", "", keyword, "", "", 0, 0,
			1, 1)
		if total == 0 {
			authorIncome := spiderBusiness.GetAuthorByKeyword(keyword)
			authorIncome.AuthorId = business.IdEncrypt(authorIncome.AuthorId)
			authorIncome.Avatar = dyimg.Fix(authorIncome.Avatar)
			receiver.SuccReturn(authorIncome)
			return
		} else {
			for _, author := range list {
				authorIncome := dy2.DyAuthorIncome{
					AuthorId:     author.AuthorId,
					Avatar:       author.Avatar,
					Nickname:     author.Nickname,
					UniqueId:     author.UniqueId,
					IsCollection: 1,
				}
				authorIncome.AuthorId = business.IdEncrypt(authorIncome.AuthorId)
				authorIncome.Avatar = dyimg.Fix(authorIncome.Avatar)
				receiver.SuccReturn(authorIncome)
				return
			}
		}
	}
}

//达人收录 调用抖音接口获取10条记录
func (receiver *AuthorController) AuthorListIncomeSearch() {
	var authorIncome = &dy2.DyAuthorIncome{}
	keyword := receiver.GetString("keyword", "")
	isCtnSearch, _ := receiver.GetInt("isCtnSearch", 0)
	if keyword == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	spiderBusiness := business.NewSpiderBusiness()
	authorIncomeList := make([]dy2.DyAuthorIncome, 0)
	if isCtnSearch == 0 {
		if utils.CheckType(keyword, "url") { // 抓换链接
			shortUrl, _ := business.ParseDyShortUrlToSecUid(keyword)
			if shortUrl == "" {
				receiver.FailReturn(global.NewError(4000))
				return
			}
			secUid := utils.ParseDyAuthorSecUrl(shortUrl) // 获取authorId
			authorIncome = spiderBusiness.GetAuthorBaseInfoV2(secUid)
			if authorIncome != nil {
				authorIncome.AuthorId = business.IdEncrypt(authorIncome.AuthorId)
				authorIncome.Avatar = dyimg.Fix(authorIncome.Avatar)
				authorIncomeList = append(authorIncomeList, *authorIncome)
			} else {
				receiver.FailReturn(global.NewError(4000))
				return
			}
		} else {
			// 如果是keyword形式的，先查es，es没有数据就请求爬虫数据接口
			list, total, _ := es.NewEsAuthorBusiness().SimpleSearch(
				"", "", keyword, "", "", 0, 0,
				1, 1)
			if total == 0 {
				authorIncomeRawList, err1 := spiderBusiness.GetAuthorListByKeyword(keyword)
				for _, v := range authorIncomeRawList {
					var tempAuthor dy2.DyAuthorIncome
					tempAuthor.UniqueId = v.UniqueId
					if tempAuthor.UniqueId == "0" || tempAuthor.UniqueId == "" {
						tempAuthor.UniqueId = v.ShortId
					}
					tempAuthor.AuthorId = business.IdEncrypt(v.Id)
					tempAuthor.Avatar = dyimg.Fix(v.Avatar)
					tempAuthor.Nickname = v.Nickname
					tempAuthor.IsCollection = 0
					authorIncomeList = append(authorIncomeList, tempAuthor)
				}
				if err1 != nil {
					receiver.FailReturn(global.NewMsgError(err1.Error()))
					return
				}
			} else {
				for _, author := range list {
					authorIncome := dy2.DyAuthorIncome{
						AuthorId:     author.AuthorId,
						Avatar:       author.Avatar,
						Nickname:     author.Nickname,
						UniqueId:     author.UniqueId,
						IsCollection: 1,
					}
					authorIncome.AuthorId = business.IdEncrypt(authorIncome.AuthorId)
					authorIncome.Avatar = dyimg.Fix(authorIncome.Avatar)
					authorIncomeList = append(authorIncomeList, authorIncome)
				}
			}
		}
	} else { //继续搜索
		authorIncomeRawList, err1 := spiderBusiness.GetAuthorListByKeyword(keyword)
		for _, v := range authorIncomeRawList {
			var tempAuthor dy2.DyAuthorIncome
			tempAuthor.UniqueId = v.UniqueId
			if tempAuthor.UniqueId == "0" || tempAuthor.UniqueId == "" {
				tempAuthor.UniqueId = v.ShortId
			}
			tempAuthor.AuthorId = business.IdEncrypt(v.Id)
			tempAuthor.Avatar = dyimg.Fix(v.Avatar)
			tempAuthor.Nickname = v.Nickname
			tempAuthor.IsCollection = 0
			authorIncomeList = append(authorIncomeList, tempAuthor)
		}
		if err1 != nil {
			receiver.FailReturn(global.NewMsgError(err1.Error()))
			return
		}
	}

	end := 10
	if len(authorIncomeList) < end {
		end = len(authorIncomeList)
	}
	resList := authorIncomeList[0:end]
	listTotal := len(resList)
	receiver.SuccReturn(map[string]interface{}{
		"list":  resList,
		"total": listTotal,
	})
	return
}

//达人收录 确认收入
func (receiver *AuthorController) AuthorIncome() {
	authorId := receiver.InputFormat().GetString("author_id", "")
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	authorIdDec := business.IdDecrypt(authorId)
	spiderBusiness := business.NewSpiderBusiness()
	ret, ok := spiderBusiness.SpiderSpeedUp("author", authorIdDec)
	if ok {
		receiver.SuccReturn([]string{authorIdDec})
	} else {
		receiver.FailReturn(global.NewError(4000))
	}
	logs.Info("[收入达人结果]：", ret)
	return
}

//达人搜索
func (receiver *AuthorController) AuthorSearch() {
	keyword := receiver.GetString("keyword", "")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	authorId := ""
	if strings.Index(keyword, "http://") > 0 || strings.Index(keyword, "https://") > 0 {
		keyword = strings.Replace(keyword, "在抖音，记录美好生活！ ", "", 1)
	}
	if utils.CheckType(keyword, "url") {
		spiderBusiness := business.SpiderBusiness{}
		shortUrl, _ := business.ParseDyShortUrlToSecUid(keyword)
		if shortUrl == "" {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		secUid := utils.ParseDyAuthorSecUrl(shortUrl) // 获取authorId 抖音获取不到author_id 只能获取sec_uid
		author := spiderBusiness.GetAuthorBaseInfoV2(secUid)
		if author != nil {
			authorId = author.AuthorId
		} else {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		keyword = ""
	} else {
		keyword = utils.MatchDouyinNewText(keyword)
	}
	list, total, comErr := es.NewEsAuthorBusiness().SimpleSearch(authorId, "", keyword, "", "", 0, 0, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		list[k].Avatar = dyimg.Fix(v.Avatar)
		if v.UniqueId == "" || v.UniqueId == "0" {
			list[k].UniqueId = v.ShortId
		}
	}
	if total > 10000 {
		total = 10000
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
	return
}

/**爬虫加速**/
func (receiver *AuthorController) SpiderSpeedUp() {
	if !business.UserActionLock(receiver.TrueUri, utils.ToString(receiver.UserId), 5) {
		receiver.FailReturn(global.NewError(6000))
		return
	}

	AuthorId := business.IdDecrypt(receiver.GetString(":author_id", ""))
	if AuthorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	spriderName := "author"
	cacheKey := cache.GetCacheKey(cache.SpiderSpeedUpLimit, spriderName, AuthorId)
	cacheData := global.Cache.Get(cacheKey)
	if cacheData != "" {
		//缓存存在
		receiver.FailReturn(global.NewError(6000))
		return
	}
	//加速
	ret, _ := business.NewSpiderBusiness().SpiderSpeedUp(spriderName, AuthorId)
	global.Cache.Set(cacheKey, "1", 300)

	logs.Info("达人加速，爬虫推送结果：", ret)
	receiver.SuccReturn([]string{})
	return
}
