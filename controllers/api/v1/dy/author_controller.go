package dy

import (
	"dongchamao/business"
	"dongchamao/business/es"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	"dongchamao/models/entity"
	dy2 "dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	jsoniter "github.com/json-iterator/go"
	"math"
	"sort"
	"time"
)

type AuthorController struct {
	controllers.ApiBaseController
}

func (receiver *AuthorController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseShowNum)
}

//达人分类
func (receiver *AuthorController) AuthorCate() {
	var cateList []dy2.DyCate
	var cateFirst []dcm.DcAuthorCate
	var cateSecond []dcm.DcAuthorCate
	db := dcm.GetDbSession().Table(dcm.DcAuthorCate{})
	if err := db.Where("level=?", 1).Find(&cateFirst); err != nil {
		panic(err)
		return
	}
	if err := db.Where("level=?", 2).Find(&cateSecond); err != nil {
		panic(err)
		return
	}
	for _, v := range cateFirst {
		var cate = dy2.DyCate{
			Name:    v.Name,
			SonCate: []dy2.DyCate{},
		}
		for _, vv := range cateSecond {
			if vv.ParentId == v.Id {
				cate.SonCate = append(cate.SonCate, dy2.DyCate{
					Name:    vv.Name,
					SonCate: []dy2.DyCate{},
				})
			}
		}
		cateList = append(cateList, cate)
	}
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
	minFanAge, _ := receiver.GetInt("min_fan_age", 0)
	maxFanAge, _ := receiver.GetInt("max_fan_age", 0)
	gender, _ := receiver.GetInt("gender", 0)
	fanGender, _ := receiver.GetInt("fan_gender", 0)
	verification, _ := receiver.GetInt("verification", 0)
	level, _ := receiver.GetInt("level", 0)
	isBrand, _ := receiver.GetInt("is_brand", 0)
	isDelivery, _ := receiver.GetInt("is_delivery", 0)
	superSeller, _ := receiver.GetInt("super_seller", 0)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	receiver.KeywordBan(keyword)
	if !receiver.HasLogin && keyword != "" {
		receiver.FailReturn(global.NewError(4001))
		return
	}
	if !receiver.HasAuth {
		if category != "" || secondCategory != "" || sellTags != "" || province != "" || city != "" || fanProvince != "" || fanCity != "" || sortStr != "follower_incre_count" || orderBy != "desc" ||
			minFollower > 0 || maxFollower > 0 || minWatch > 0 || maxWatch > 0 || minDigg > 0 || maxDigg > 0 || minGmv > 0 || maxGmv > 0 ||
			gender > 0 || minAge > 0 || maxAge > 0 || minFanAge > 0 || maxFanAge > 0 || verification > 0 || level > 0 || fanGender > 0 ||
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
	if utils.CheckType(keyword, "url") {
		shortUrl, _ := business.ParseDyShortUrl(keyword)
		authorId = utils.ParseDyAuthorUrl(shortUrl)
		keyword = ""
	} else {
		keyword = utils.MatchDouyinNewText(keyword)
	}
	EsAuthorBusiness := es.NewEsAuthorBusiness()
	list, total, comErr := EsAuthorBusiness.BaseSearch(authorId, keyword, category, secondCategory, sellTags, province, city, fanProvince, fanCity,
		minFollower, maxFollower, minWatch, maxWatch, minDigg, maxDigg, minGmv, maxGmv,
		gender, minAge, maxAge, minFanAge, maxFanAge, verification, level, isDelivery, isBrand, superSeller, fanGender, page, pageSize,
		sortStr, orderBy)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	authorIds := make([]string, 0)
	for _, v := range list {
		authorIds = append(authorIds, v.AuthorId)
	}
	authorMap := business.NewAuthorBusiness().GetAuthorFormPool(authorIds, 10)
	for k, v := range list {
		list[k].Avatar = dyimg.Fix(v.Avatar)
		if v.UniqueId == "" || v.UniqueId == "0" {
			list[k].UniqueId = v.ShortId
		}
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		if a, ok := authorMap[v.AuthorId]; ok {
			list[k].RoomId = a.RoomId
		} else {
			authorData, _ := hbase.GetAuthor(v.AuthorId)
			list[k].RoomId = business.IdEncrypt(authorData.RoomId)
		}
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
	receiver.SuccReturn(map[string]interface{}{
		"list":       list,
		"total":      total,
		"total_page": totalPage,
		"max_num":    maxTotal,
		"has_auth":   receiver.HasAuth,
		"has_login":  receiver.HasLogin,
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
	basicBefore, _ := hbase.GetAuthorBasic(authorId, time.Now().AddDate(0, 0, -1).Format("20060102"))
	authorBase.Data.ID = business.IdEncrypt(authorBase.Data.ID)
	authorBase.Data.RoomID = business.IdEncrypt(authorBase.Data.RoomID)
	basic := entity.DyAuthorBasic{
		FollowerCount:        authorBase.FollowerCount,
		FollowerCountBefore:  basicBefore.FollowerCount,
		TotalFansCount:       authorBase.TotalFansCount,
		TotalFansCountBefore: basicBefore.TotalFansCount,
		TotalFavorited:       authorBase.TotalFavorited,
		TotalFavoritedBefore: basicBefore.TotalFavorited,
		CommentCount:         authorBase.CommentCount,
		CommentCountBefore:   basicBefore.CommentCount,
		ForwardCount:         authorBase.ForwardCount,
		ForwardCountBefore:   basicBefore.ForwardCount,
	}
	returnMap := map[string]interface{}{
		"author_base": authorBase.Data,
		"room_count":  authorBase.RoomCount,
		"reputation": dy2.RepostSimpleReputation{
			Score:         reputation.Score,
			Level:         reputation.Level,
			EncryptShopID: reputation.EncryptShopID,
			ShopName:      reputation.ShopName,
			ShopLogo:      reputation.ShopLogo,
		},
		"fans_club": fansClub.TotalFansCount,
		"rank":      nil,
		"basic":     basic,
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
	authorBase, comErr := authorBusiness.HbaseGetAuthor(authorId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	todayString := time.Now().Format("20060102")
	monthString := time.Now().Format("200601")
	todayTime, _ := time.ParseInLocation("20060102", todayString, time.Local)
	monthTime, _ := time.ParseInLocation("20060102", monthString+"01", time.Local)
	lastMonthDay := todayTime.AddDate(0, 0, -29)
	nowWeek := int(time.Now().Weekday())
	if nowWeek == 0 {
		nowWeek = 7
	}
	lastWeekDay := todayTime.AddDate(0, 0, -(nowWeek - 1))
	var monthRoom int64 = 0
	var weekRoom int64 = 0
	var room30Count int64 = 0
	for _, v := range authorBase.RoomList {
		if v.CreateTime >= monthTime.Unix() {
			monthRoom++
		}
		if v.CreateTime >= lastWeekDay.Unix() {
			weekRoom++
		}
		if v.CreateTime >= lastMonthDay.Unix() {
			room30Count++
		}
	}
	productCount := dy2.DyAuthorBaseProductCount{}
	startTime := time.Now().AddDate(0, 0, -31)
	yesterday := time.Now().AddDate(0, 0, -1)
	cacheKey := cache.GetCacheKey(cache.AuthorViewProductAllList, authorId, startTime.Format("20060102"), yesterday.Format("20060102"))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &productCount)
	} else {
		hbaseDataList, err := authorBusiness.GetAuthorProductHbaseList(authorId, "", startTime, yesterday)
		if err == nil {
			brandSaleMap := map[string]float64{}
			brandNumMap := map[string]int{}
			priceNumMap := map[string]int{}
			priceSaleMap := map[string]float64{}
			var totalGmv float64
			var totalSales float64
			for _, v := range hbaseDataList {
				v.Gmv = v.LivePredictGmv + v.VideoPredictGmv
				v.Sales = math.Floor(v.LivePredictSales) + math.Floor(v.VideoPredictGmv)
				totalGmv += v.Gmv
				totalSales += v.Sales
				category := v.DcmLevelFirst
				if category == "" || category == "null" {
					category = "其他"
				}
				var priceStr string
				if v.Price > 500 {
					priceStr = "500-"
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
					brandSaleMap[category] = v.Sales
				} else {
					brandSaleMap[category] += v.Sales
				}
				if _, ok := brandNumMap[category]; !ok {
					brandNumMap[category] = 1
				} else {
					brandNumMap[category] += 1
				}
				if _, ok := priceSaleMap[priceStr]; !ok {
					priceSaleMap[priceStr] = v.Sales
				} else {
					priceSaleMap[priceStr] += v.Sales
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
				topBrandSaleList = append(topBrandSaleList, dy2.NameValueInt64Chart{
					Name:  "其他",
					Value: sale,
				})
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
				topBrandNumList = append(topBrandNumList, dy2.NameValueChart{
					Name:  "其他",
					Value: num,
				})
			}
			productCount = dy2.DyAuthorBaseProductCount{
				Sales30Top3:           topSaleCates,
				ProductNum30Top3:      topNumCates,
				Sales30Top3Chart:      topBrandSaleList,
				ProductNum30Top3Chart: topBrandNumList,
				Predict30Sales:        totalSales,
				Predict30Gmv:          totalGmv,
				Sales30Chart:          []dy2.DyAuthorBaseProductPriceChart{},
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
	}
	productCount.ProductNum = authorBase.ProductCount
	data := dy2.DyAuthorBaseCount{
		LiveCount: dy2.DyAuthorBaseLiveCount{
			RoomCount:      int64(authorBase.RoomCount),
			Room30Count:    room30Count,
			Predict30Sales: math.Floor(authorBase.Predict30Sales),
			Predict30Gmv:   utils.FriendlyFloat64(authorBase.Predict30Gmv),
			AgeDuration:    authorBase.AgeLiveDuration,
			WeekRoomCount:  weekRoom,
			MonthRoomCount: monthRoom,
		},
		VideoCount: dy2.DyAuthorBaseVideoCount{
			VideoCount:       int64(authorBase.AwemeCount),
			AvgDigg:          authorBase.DiggCount,
			DiggFollowerRate: authorBase.DiggFollowerRate,
			Predict30Sales:   0,
			Predict30Gmv:     0,
			AgeDuration:      authorBase.Duration / 1000,
		},
		ProductCount: productCount,
	}
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
			data.LiveCount.AvgWeekRoomCount = utils.ToInt64(math.Ceil(float64(authorBase.RoomCount) / float64(weekNum)))
		}
		var month = utils.ToInt64(math.Ceil(float64(time.Now().Unix()-firstLiveTime.Unix()) / (30 * 86400)))
		if month > 0 {
			data.LiveCount.AvgMonthRoomCount = utils.ToInt64(math.Ceil(float64(authorBase.RoomCount) / float64(month)))
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
	receiver.SuccReturn(map[string]interface{}{
		"reputation": reputation,
	})
	return
}

//达人视频概览
func (receiver *AuthorController) AuthorAwemesByDay() {
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
	aABusiness := business.NewAuthorAwemeBusiness()
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
	//时间限制
	if t1.After(t2) || t2.After(t1.AddDate(0, 0, 90)) || t2.After(time.Now()) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	videoOverview := aABusiness.HbaseGetVideoAggRangeDate(authorId, t1, t2)
	receiver.SuccReturn(map[string]interface{}{
		"video_overview": videoOverview,
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

//粉丝分布分析
func (receiver *AuthorController) AuthorFansAnalyse() {
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	detail, comErr := hbase.GetXtAuthorDetail(authorId)
	data := map[string][]entity.XtDistributionsList{}
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
	} else {
		data["gender"] = []entity.XtDistributionsList{}
		data["city"] = []entity.XtDistributionsList{}
		data["age"] = []entity.XtDistributionsList{}
		data["province"] = []entity.XtDistributionsList{}
	}
	data["active_day"] = []entity.XtDistributionsList{}
	data["active_week"] = []entity.XtDistributionsList{}
	var countCity int64 = 0
	var countPro int64 = 0
	for _, v := range data["city"] {
		countCity += v.DistributionValue
	}
	for _, v := range data["province"] {
		countPro += v.DistributionValue
	}
	for k, v := range data["city"] {
		data["city"][k].DistributionPer = float64(v.DistributionValue) / float64(countCity)
	}
	for k, v := range data["province"] {
		data["province"][k].DistributionPer = float64(v.DistributionValue) / float64(countPro)
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
	list, total, comErr := esLiveBusiness.SearchAuthorRooms(authorId, keyword, sortStr, orderBy, page, size, t1, t2)
	for k, v := range list {
		list[k].RoomId = business.IdEncrypt(v.RoomId)
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		list[k].PredictSales = math.Floor(v.PredictSales)
		list[k].PredictGmv = math.Floor(v.PredictGmv)
		if listType == 1 {
			liveInfo, _ := hbase.GetLiveInfo(v.RoomId)
			list[k].FinishTime = liveInfo.FinishTime
		}
	}
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
	list, analysisCount, cateList, brandList, total, comErr := authorBusiness.GetAuthorProductAnalyse(authorId, keyword, firstCate, secondCate, thirdCate, brandName, sortStr, orderBy, shopType, startTime, endTime, page, pageSize)
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

//达人收录
func (receiver *AuthorController) AuthorIncome() {
	var authorId string
	var authorIncome = &dy2.DyAuthorIncome{}
	keyword := receiver.InputFormat().GetString("keyword", "")
	if keyword == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	spiderBusiness := business.NewSpiderBusiness()
	if utils.CheckType(keyword, "url") { // 抓换链接
		shortUrl, _ := business.ParseDyShortUrl(keyword)
		if shortUrl == "" {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		authorId = utils.ParseDyAuthorUrl(shortUrl) // 获取authorId
		author, err := hbase.GetAuthor(authorId)
		if err == nil {
			authorIncome = &dy2.DyAuthorIncome{
				AuthorId:     author.AuthorID,
				Avatar:       author.Data.Avatar,
				Nickname:     author.Data.Nickname,
				UniqueId:     author.Data.UniqueID,
				IsCollection: 0,
			}
		} else {
			authorIncome = spiderBusiness.GetAuthorBaseInfo(authorId)
		}
		receiver.SuccReturn(authorIncome)
		return
	} else {
		// 如果是keyword形式的，先查es，es没有数据就请求爬虫数据接口
		list, total, _ := es.NewEsAuthorBusiness().SimpleSearch(
			"", keyword, "", "",
			1, 1)
		if total == 0 {
			authorIncome := spiderBusiness.GetAuthorByKeyword(keyword)
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
				receiver.SuccReturn(authorIncome)
				return
			}
		}
	}
}
