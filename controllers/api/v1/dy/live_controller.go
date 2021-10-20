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
	dy2 "dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	jsoniter "github.com/json-iterator/go"
	"math"
	"sort"
	"strings"
	"time"
)

type LiveController struct {
	controllers.ApiBaseController
}

func (receiver *LiveController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseMinShowNum, business.DyJewelBaseShowNum)
}

//直播库
func (receiver *LiveController) SearchRoom() {
	startDay := receiver.GetString("start", "")
	endDay := receiver.GetString("end", "")
	if startDay == "" {
		startDay = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	}
	if endDay == "" {
		endDay = time.Now().Format("2006-01-02")
	}
	pslTime := "2006-01-02"
	startTime, err := time.ParseInLocation(pslTime, startDay, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	endTime, err := time.ParseInLocation(pslTime, endDay, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if startTime.After(endTime) || endTime.After(startTime.AddDate(0, 0, 90)) || endTime.After(time.Now()) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	keyword := receiver.GetString("keyword", "")
	category := receiver.GetString("category", "")
	firstName := receiver.GetString("first_name", "")
	secondName := receiver.GetString("second_name", "")
	thirdName := receiver.GetString("third_name", "")
	sortStr := receiver.GetString("sort", "predict_gmv")
	orderBy := receiver.GetString("order_by", "desc")
	minAmount, _ := receiver.GetInt64("min_amount", 0)
	maxAmount, _ := receiver.GetInt64("max_amount", 0)
	minAvgUserCount, _ := receiver.GetInt64("min_avg_user_count", 0)
	maxAvgUserCount, _ := receiver.GetInt64("max_avg_user_count", 0)
	minUv, _ := receiver.GetInt("min_uv", 0)
	maxUv, _ := receiver.GetInt("max_uv", 0)
	hasProduct, _ := receiver.GetInt("has_product", 0)
	isBrand, _ := receiver.GetInt("is_brand", 0)
	keywordType, _ := receiver.GetInt("keyword_type", 0)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	receiver.KeywordBan(keyword)
	if !receiver.HasAuth {
		today := time.Now().Format("20060102")
		lastDay := time.Now().AddDate(0, 0, -6).Format("20060102")
		start := startTime.Format("20060102")
		end := endTime.Format("20060102")
		if lastDay != start || today != end || category != "" || sortStr != "predict_gmv" || orderBy != "desc" || minAmount > 0 || maxAmount > 0 || minUv > 0 || maxUv > 0 || minAvgUserCount > 0 || maxAvgUserCount > 0 || hasProduct == 1 || isBrand == 1 || page != 1 {
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
	esLiveBusiness := es.NewEsLiveBusiness()
	list, total, comErr := esLiveBusiness.SearchLiveRooms(keyword, category, firstName, secondName, thirdName, minAmount, maxAmount, minAvgUserCount, maxAvgUserCount, minUv, maxUv, hasProduct, isBrand, keywordType, sortStr, orderBy, page, pageSize, startTime, endTime)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		list[k].RoomId = business.IdEncrypt(v.RoomId)
		list[k].Cover = dyimg.Fix(v.Cover)
		list[k].Avatar = dyimg.Fix(v.Avatar)
		////todo gmv处理
		//if v.RealGmv > 0 {
		//	list[k].PredictGmv = v.RealGmv
		//}
		//if v.RealUvValue > 0 {
		//	list[k].PredictUvValue = v.RealUvValue
		//}
		list[k].AvgUserCount = math.Floor(v.AvgUserCount)
		if v.DisplayId == "" {
			list[k].DisplayId = v.ShortId
		}
		list[k].TagsArr = v.GetTagsArr()
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
	})
	return
}

//直播详细
func (receiver *LiveController) LiveInfoData() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	liveBusiness := business.NewLiveBusiness()
	liveInfo, comErr := hbase.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	reputation, _ := hbase.GetLiveReputation(roomId)
	authorInfo, _ := authorBusiness.HbaseGetAuthor(liveInfo.User.ID)
	uniqueId := authorInfo.Data.UniqueID
	if uniqueId == "" {
		uniqueId = authorInfo.Data.ShortID
	}
	liveUser := dy2.DyLiveUserSimple{
		Avatar:          liveInfo.User.Avatar,
		FollowerCount:   authorInfo.Data.FollowerCount,
		ID:              business.IdEncrypt(liveInfo.User.ID),
		UniqueId:        uniqueId,
		Nickname:        liveInfo.User.Nickname,
		WithCommerce:    liveInfo.User.WithCommerce,
		ReputationScore: reputation.AuthorReputation.Score,
		ReputationLevel: reputation.AuthorReputation.Level,
		RoomId:          authorInfo.RoomId,
	}
	//liveSaleData, _ := hbase.GetLiveSalesData(roomId)
	incOnlineTrends, maxOnlineTrends, avgUserCount := liveBusiness.DealOnlineTrends(liveInfo)
	var incFansRate, interactRate float64
	incFansRate = 0
	interactRate = 0
	liveSale := dy2.DyLiveRoomSaleData{}
	//todo gmv数据兼容
	//gmv := liveSaleData.Gmv
	//sales := liveSaleData.Sales
	//if liveSaleData.Gmv == 0 {
	gmv := liveInfo.PredictGmv
	sales := liveInfo.PredictSales
	//if liveInfo.RealGmv > 0 {
	//	gmv = liveInfo.RealGmv
	//	sales = liveInfo.RealSales
	//}
	//}
	if liveInfo.TotalUser > 0 {
		incFansRate = float64(liveInfo.FollowCount) / float64(liveInfo.TotalUser)
		interactRate = float64(liveInfo.BarrageUserCount) / float64(liveInfo.TotalUser)
		liveSale.Uv = (gmv + float64(liveInfo.RoomTicketCount)/10) / float64(liveInfo.TotalUser)
		liveSale.SaleRate = utils.RateMin(sales / float64(liveInfo.TotalUser))
	}
	avgOnlineTime := liveBusiness.CountAvgOnlineTime(liveInfo.OnlineTrends, liveInfo.CreateTime, liveInfo.TotalUser)
	returnLiveInfo := dy2.DyLiveInfo{
		Cover:               liveInfo.Cover,
		CreateTime:          liveInfo.CreateTime,
		FinishTime:          liveInfo.FinishTime,
		LikeCount:           liveInfo.LikeCount,
		RoomID:              business.IdEncrypt(liveInfo.RoomID),
		RoomStatus:          liveInfo.RoomStatus,
		Title:               liveInfo.Title,
		TotalUser:           liveInfo.TotalUser,
		User:                liveUser,
		UserCount:           liveInfo.UserCount,
		TrendsCrawlTime:     liveInfo.TrendsCrawlTime,
		IncFans:             liveInfo.FollowCount,
		IncFansRate:         utils.RateMin(incFansRate),
		InteractRate:        utils.RateMin(interactRate),
		AvgUserCount:        avgUserCount,
		MaxWatchOnlineTrend: maxOnlineTrends,
		OnlineTrends:        incOnlineTrends,
		RenewalTime:         liveInfo.CrawlTime,
		AvgOnlineTime:       avgOnlineTime,
		LiveUrl:             liveInfo.PlayURL,
		ShareUrl:            business.LiveShareUrl + liveInfo.RoomID,
	}
	liveSale.Volume = int64(math.Floor(sales))
	liveSale.Amount = gmv
	esLiveBusiness := es.NewEsLiveBusiness()
	liveSale.PromotionNum = esLiveBusiness.CountRoomProductByRoomId(liveInfo)
	if sales > 0 {
		liveSale.PerPrice = gmv / sales
	}
	//dateChart := make([]int64, 0)
	//gmvChart := make([]float64, 0)
	//salesChart := make([]float64, 0)
	//salesTrends := liveInfo.SalesTrends
	////排序
	//sort.Slice(salesTrends, func(i, j int) bool {
	//	var left, right int64
	//	left = salesTrends[i].CrawlTime
	//	right = salesTrends[j].CrawlTime
	//	return right > left
	//})
	//for _, v := range salesTrends {
	//	dateChart = append(dateChart, v.CrawlTime)
	//	//if liveInfo.RealGmv > 0 {
	//	//	gmvChart = append(gmvChart, v.RealGmv)
	//	//	salesChart = append(salesChart, math.Floor(v.RealSales))
	//	//} else {
	//	gmvChart = append(gmvChart, v.PredictGmv)
	//	salesChart = append(salesChart, math.Floor(v.PredictSales))
	//	//}
	//}
	//gmvChart = business.DealIncDirtyFloat64Chart(gmvChart)
	//salesChart = business.DealIncDirtyFloat64Chart(salesChart)

	//处理直播间大盘数据
	esInfo, _ := es.NewEsLiveBusiness().SearchRoomById(&liveInfo)
	liveLevel := map[string]interface{}{
		"date":           time.Unix(liveInfo.DiscoverTime, 0).Format("2006-01-02"),
		"flow_rates":     esInfo.FlowRates,
		"avg_stay_index": esInfo.AvgStayIndex,
		"tags":           esInfo.Tags,
		"tags_arr":       strings.Split(esInfo.Tags, "_"),
	}
	receiver.SuccReturn(map[string]interface{}{
		"live_info":              returnLiveInfo,
		"live_sale":              liveSale,
		"user_count_composition": liveInfo.UserCountComposition,
		"live_level":             liveLevel,
		//"sales_chart": map[string]interface{}{
		//	"time":  dateChart,
		//	"gmv":   gmvChart,
		//	"sales": salesChart,
		//},
	})
	return
}

//直播商品明细
func (receiver *LiveController) LivePromotions() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	liveInfo, comErr := hbase.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	promotionsList := dy2.DyLivePromotionChart{
		StartTime:     []string{},
		PromotionList: [][]dy2.DyLivePromotion{},
	}
	startRowKey, stopRowKey, err := es.NewEsLiveBusiness().ScanProductByRoomId(liveInfo)
	cacheKey := cache.GetCacheKey(cache.LivePromotionsDetailList, startRowKey, stopRowKey)
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &promotionsList)
		receiver.SuccReturn(map[string]interface{}{
			"promotions_list": promotionsList,
		})
		return
	}
	if err == nil {
		promotionsMap := map[string][]entity.DyLivePromotion{}
		promotionsSalesMap := map[string]int64{}
		stopRow, err1 := hbase.GetRoomProductInfo(stopRowKey)
		if err1 == nil {
			tmpProductId := ""
			for _, v := range stopRow.PtmPromotion {
				startFormat := time.Unix(v.StartTime, 0).Format("2006-01-02 15:04:05")
				if _, ok := promotionsMap[startFormat]; !ok {
					promotionsMap[startFormat] = []entity.DyLivePromotion{}
				}
				promotionsMap[startFormat] = append(promotionsMap[startFormat], v)
				tmpProductId = v.ProductID
			}
			promotionsSalesMap[tmpProductId] = utils.ToInt64(math.Floor(stopRow.PredictSales))
		}
		if startRowKey != stopRowKey {
			mapData, _ := hbase.GetRoomProductInfoRangDate(startRowKey, stopRowKey)
			for k, d := range mapData {
				promotionsSalesMap[k] = utils.ToInt64(math.Floor(d.PredictSales))
				for _, v := range d.PtmPromotion {
					startFormat := time.Unix(v.StartTime, 0).Format("2006-01-02 15:04:05")
					if _, ok := promotionsMap[startFormat]; !ok {
						promotionsMap[startFormat] = []entity.DyLivePromotion{}
					}
					promotionsMap[startFormat] = append(promotionsMap[startFormat], v)
				}
			}
		}
		dates := make([]string, 0)
		dyLivePromotions := make([][]dy2.DyLivePromotion, 0)
		promotionSales := map[string]int{}
		//按时间排序
		var keys []string
		for k := range promotionsMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := promotionsMap[k]
			item := make([]dy2.DyLivePromotion, 0)
			for _, v1 := range v {
				saleNum := 1
				if s, ok := promotionSales[v1.ProductID]; ok {
					saleNum = s + 1
				}
				promotionSales[v1.ProductID] = saleNum
				var sales int64 = 0
				if sa, ok := promotionsSalesMap[v1.ProductID]; ok {
					sales = sa
				}
				item = append(item, dy2.DyLivePromotion{
					ProductID: business.IdEncrypt(v1.ProductID),
					ForSale:   v1.ForSale,
					StartTime: v1.StartTime,
					StopTime:  v1.StopTime,
					Price:     v1.Price,
					Sales:     v1.Sales,
					NowSales:  sales,
					GmvSales:  sales,
					Title:     v1.Title,
					Cover:     dyimg.Product(v1.Cover),
					Index:     v1.Index,
					SaleNum:   saleNum,
				})
			}
			dyLivePromotions = append(dyLivePromotions, item)
			dates = append(dates, k)
		}
		promotionsList = dy2.DyLivePromotionChart{
			StartTime:     dates,
			PromotionList: dyLivePromotions,
		}
		var cacheTime time.Duration = 60
		if liveInfo.RoomStatus == 4 {
			cacheTime = 1800
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(promotionsList), cacheTime)
	}
	receiver.SuccReturn(map[string]interface{}{
		"promotions_list": promotionsList,
	})
	return
}

//直播榜单排名趋势
func (receiver *LiveController) LiveRankTrends() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	liveBusiness := business.NewLiveBusiness()
	liveRankTrends, _ := liveBusiness.HbaseGetRankTrends(roomId)
	saleDates := make([]int64, 0)
	hourDates := make([]int64, 0)
	saleRanks := make([]int, 0)
	hourRanks := make([]int, 0)
	maxSaleRank := 1000000
	maxHourRank := 1000000
	for _, v := range liveRankTrends {
		if v.Type == 8 {
			saleDates = append(saleDates, v.CrawlTime)
			saleRanks = append(saleRanks, v.Rank)
			if v.Rank < maxSaleRank {
				maxSaleRank = v.Rank
			}
		} else if v.Type == 1 {
			hourDates = append(hourDates, v.CrawlTime)
			hourRanks = append(hourRanks, v.Rank)
			if v.Rank < maxHourRank {
				maxHourRank = v.Rank
			}
		}
	}
	hourDates = business.DealChartInt64(hourDates, 60)
	hourRanks = business.DealChartInt(hourRanks, 60)
	saleDates = business.DealChartInt64(saleDates, 60)
	saleRanks = business.DealChartInt(saleRanks, 60)
	receiver.SuccReturn(map[string]interface{}{
		"hour_rank": map[string]interface{}{
			"time":  hourDates,
			"ranks": hourRanks,
		},
		"sale_rank": map[string]interface{}{
			"time":  saleDates,
			"ranks": saleRanks,
		},
		"max_hour_rank": maxHourRank,
		"max_sale_rank": maxSaleRank,
	})
}

//直播间商品
func (receiver *LiveController) LiveProductList() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	InputData := receiver.InputFormat()
	keyword := InputData.GetString("keyword", "")
	productId := business.IdDecrypt(InputData.GetString("product_id", ""))
	sortStr := InputData.GetString("sort", "shelf_time")
	orderBy := InputData.GetString("order_by", "desc")
	page := InputData.GetInt("page", 1)
	pageSize := InputData.GetInt("page_size", 10)
	firstLabel := InputData.GetString("first_label", "")
	secondLabel := InputData.GetString("second_label", "")
	thirdLabel := InputData.GetString("third_label", "")
	roomInfo, _ := hbase.GetLiveInfo(roomId)
	esLiveBusiness := es.NewEsLiveBusiness()
	list, productCount, total, err := esLiveBusiness.RoomProductByRoomId(roomInfo, keyword, productId, sortStr, orderBy, firstLabel, secondLabel, thirdLabel, page, pageSize)
	if err != nil {
		receiver.FailReturn(err)
		return
	}
	countList := make([]dy2.LiveRoomProductCount, 0)
	if len(list) > 0 {
		liveBusiness := business.NewLiveBusiness()
		//curMap := liveBusiness.RoomCurProductByIds(roomId, productIds)
		//pmtMap := liveBusiness.RoomPmtProductByIds(roomId, productIds)
		for _, v := range list {
			curCount, pmtStatus, pv, err1 := liveBusiness.RoomCurAndPmtProductById(roomId, v.ProductID)
			v.Pv = pv
			if pv > 0 {
				v.BuyRate = utils.RateMin(v.PredictSales / float64(pv))
			}
			item := dy2.LiveRoomProductCount{
				ProductInfo: v,
				ProductStartSale: dy2.RoomProductSaleChart{
					Timestamp: []int64{},
					Sales:     []int64{},
				},
				ProductEndSale: dy2.RoomProductSaleChart{
					Timestamp: []int64{},
					Sales:     []int64{},
				},
			}
			if err1 == nil {
				for _, s1 := range pmtStatus {
					item.ProductStartSale.Timestamp = append(item.ProductStartSale.Timestamp, s1.StartTime)
					item.ProductStartSale.Sales = append(item.ProductStartSale.Sales, s1.StartSales)
					if s1.StopTime > 0 {
						item.ProductEndSale.Timestamp = append(item.ProductEndSale.Timestamp, s1.StopTime)
						item.ProductEndSale.Sales = append(item.ProductEndSale.Sales, s1.FinalSales)
					}
				}
				item.ProductCur = curCount
			} else {
				item.ProductCur = dy2.LiveCurProductCount{
					CurList: []dy2.LiveCurProduct{},
				}
			}
			item.ProductInfo.AuthorID = business.IdEncrypt(item.ProductInfo.AuthorID)
			item.ProductInfo.AuthorRoomID = business.IdEncrypt(item.ProductInfo.AuthorRoomID)
			item.ProductInfo.RoomID = business.IdEncrypt(item.ProductInfo.RoomID)
			item.ProductInfo.ProductID = business.IdEncrypt(item.ProductInfo.ProductID)
			countList = append(countList, item)
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":          countList,
		"product_count": productCount,
		"total":         total,
	})
	return
}

//直播间商品分类
func (receiver *LiveController) LiveProductCateList() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	roomInfo, _ := hbase.GetLiveInfo(roomId)
	esLiveBusiness := es.NewEsLiveBusiness()
	countData := esLiveBusiness.AllRoomProductCateByRoomId(roomInfo)
	receiver.SuccReturn(map[string]interface{}{
		"count": countData,
	})
	return
}

//全网销量趋势图
func (receiver *LiveController) LiveProductSaleChart() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	productId := business.IdDecrypt(receiver.Ctx.Input.Param(":product_id"))
	if productId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	info, _ := hbase.GetRoomProductTrend(roomId + "_" + productId)
	trends := business.RoomProductTrendOrderByTime(info.TrendData)
	timestamps := make([]int64, 0)
	sales := make([]float64, 0)
	for _, v := range trends {
		timestamps = append(timestamps, v.CrawlTime)
		sales = append(sales, math.Floor(v.Sales))
	}
	receiver.SuccReturn(dy2.TimestampCountChart{
		Timestamp:  timestamps,
		CountValue: sales,
	})
	return
}

//直播间粉丝趋势
func (receiver *LiveController) LiveFansTrends() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	info, comErr := hbase.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	fansDate := make([]int64, 0)
	clubDate := make([]int64, 0)
	fansTrends := make([]int64, 0)
	fansIncTrends := make([]int64, 0)
	clubTrends := make([]int64, 0)
	clubIncTrends := make([]int64, 0)
	if len(info.FollowerCountTrends) > 0 {
		followerCountTrends := business.LiveFansTrendsListOrderByTime(info.FollowerCountTrends)
		fansDate = append(fansDate, info.CreateTime)
		lenNum := len(followerCountTrends)
		//beforeFansTrend := entity.LiveFollowerCountTrends{
		//	CrawlTime:     info.CreateTime,
		//	FollowerCount: followerCountTrends[lenNum-1].FollowerCount - info.FollowCount,
		//	NewFollowerCount: 0,
		//}
		fansTrends = append(fansTrends, followerCountTrends[lenNum-1].FollowerCount-info.FollowCount)
		fansIncTrends = append(fansIncTrends, 0)
		for _, v := range followerCountTrends {
			fansDate = append(fansDate, v.CrawlTime)
			fansTrends = append(fansTrends, v.FollowerCount)
			//inc := v.FollowerCount - beforeFansTrend.FollowerCount
			fansIncTrends = append(fansIncTrends, v.NewFollowCount)
			//beforeFansTrend = v
		}
	}
	var clubInc int64 = 0
	if len(info.FansClubCountTrends) > 0 {
		fansClubCountTrends := business.LiveClubFansTrendsListOrderByTime(info.FansClubCountTrends)
		beforeClubTrend := entity.LiveAnsClubCountTrends{
			FansClubCount:     fansClubCountTrends[0].FansClubCount - fansClubCountTrends[0].TodayNewFansCount,
			TodayNewFansCount: fansClubCountTrends[0].TodayNewFansCount,
			CrawlTime:         info.CreateTime,
		}
		clubDate = append(clubDate, beforeClubTrend.CrawlTime)
		clubTrends = append(clubTrends, beforeClubTrend.FansClubCount)
		clubIncTrends = append(clubIncTrends, beforeClubTrend.TodayNewFansCount)
		for _, v := range fansClubCountTrends {
			clubDate = append(clubDate, v.CrawlTime)
			clubTrends = append(clubTrends, v.FansClubCount)
			inc := v.FansClubCount - beforeClubTrend.FansClubCount
			clubIncTrends = append(clubIncTrends, inc)
			beforeClubTrend = v
		}
		lenNum := len(clubTrends)
		clubInc = clubTrends[lenNum-1] - clubTrends[0]
	}
	var incFansRate float64 = 0
	if info.TotalUser > 0 {
		incFansRate = float64(info.FollowCount) / float64(info.TotalUser)
	}
	fansIncTrends = business.DealIncDirtyInt64Chart(fansIncTrends)
	receiver.SuccReturn(map[string]interface{}{
		"fans_chart": map[string]interface{}{
			"date":  fansDate,
			"count": fansTrends,
			"inc":   fansIncTrends,
		},
		"club_chart": map[string]interface{}{
			"date":  clubDate,
			"count": clubTrends,
			"inc":   clubIncTrends,
		},
		"inc_fans":      info.FollowCount,
		"inc_club":      clubInc,
		"inc_fans_rate": incFansRate,
	})
	return
}

//直播粉丝分析
func (receiver *LiveController) LiveFanAnalyse() {
	roomType := receiver.Ctx.Input.Param(":type")
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	liveInfo, comErr := hbase.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	info, _ := hbase.GetDyLiveRoomUserInfo(roomId)
	var roomUserTotal int64 = 0
	var roomAgePeopleTotal int64 = 0
	var roomGenderTotal int64 = 0
	var roomAgeTotal int64 = 0
	var roomCityTotal int64 = 0
	var roomProvinceTotal int64 = 0
	genderChart := make([]entity.XtDistributionsList, 0)
	ageChart := make([]entity.XtDistributionsList, 0)
	cityChart := make([]entity.XtDistributionsList, 0)
	provinceChart := make([]entity.XtDistributionsList, 0)
	wordChart := make([]entity.XtDistributionsList, 0)
	for k, v := range info.Gender {
		roomUserTotal += v
		name := ""
		if k == "男" {
			name = "male"
		} else if k == "女" {
			name = "female"
		} else {
			continue
		}
		roomGenderTotal += v
		genderChart = append(genderChart, entity.XtDistributionsList{
			DistributionKey:   name,
			DistributionValue: v,
		})
	}
	ageWeightMap := map[string]int{"-18": 1, "18-23": 2, "24-30": 3, "31-40": 4, "41-50": 5, "50-": 6}
	ageMap := map[string]int64{}
	for k, v := range info.AgeDistrinbution {
		roomAgePeopleTotal += v
		if k == "" {
			continue
		}
		roomAgeTotal += v
		name := receiver.liveFansAgeMap(k)
		if _, exist := ageMap[name]; !exist {
			ageMap[name] = v
		} else {
			ageMap[name] += v
		}
	}
	for k, v := range ageMap {
		weight := 0
		if w, exist := ageWeightMap[k]; exist {
			weight = w
		}
		ageChart = append(ageChart, entity.XtDistributionsList{
			DistributionKey:   k,
			DistributionValue: v,
			Weight:            weight,
		})
	}
	for k, v := range info.City {
		if k == "" {
			continue
		}
		roomCityTotal += v
		cityChart = append(cityChart, entity.XtDistributionsList{
			DistributionKey:   k,
			DistributionValue: v,
		})
	}
	for k, v := range info.Province {
		if k == "" {
			continue
		}
		roomProvinceTotal += v
		provinceChart = append(provinceChart, entity.XtDistributionsList{
			DistributionKey:   k,
			DistributionValue: v,
		})
	}
	if roomType != "ing" {
		if len(info.Word) > 0 {
			for k, v := range info.Word {
				if k == "" {
					continue
				}
				wordChart = append(wordChart, entity.XtDistributionsList{
					DistributionKey:   k,
					DistributionValue: v,
				})
			}
			sort.Slice(wordChart, func(i, j int) bool {
				return wordChart[i].DistributionValue > wordChart[j].DistributionValue
			})
			if len(wordChart) > 100 {
				wordChart = wordChart[:100]
			}
		}
	}
	sort.Slice(cityChart, func(i, j int) bool {
		return cityChart[i].DistributionValue > cityChart[j].DistributionValue
	})
	sort.Slice(provinceChart, func(i, j int) bool {
		return provinceChart[i].DistributionValue > provinceChart[j].DistributionValue
	})
	if len(cityChart) > 10 {
		cityChart = cityChart[:10]
	}
	if len(provinceChart) > 10 {
		provinceChart = provinceChart[:10]
	}
	if roomGenderTotal > 0 {
		for k, v := range genderChart {
			genderChart[k].DistributionPer = float64(v.DistributionValue) / float64(roomGenderTotal)
		}
	}
	if roomAgeTotal > 0 {
		for k, v := range ageChart {
			ageChart[k].DistributionPer = float64(v.DistributionValue) / float64(roomAgeTotal)
		}
	}
	if roomCityTotal > 0 {
		for k, v := range cityChart {
			cityChart[k].DistributionPer = float64(v.DistributionValue) / float64(roomCityTotal)
		}
	}
	if roomProvinceTotal > 0 {
		for k, v := range provinceChart {
			provinceChart[k].DistributionPer = float64(v.DistributionValue) / float64(roomProvinceTotal)
		}
	}
	var barrageRate float64 = 0
	if liveInfo.TotalUser > 0 {
		barrageRate = float64(liveInfo.BarrageUserCount) / float64(liveInfo.TotalUser)
	}
	sort.Slice(ageChart, func(i, j int) bool {
		return ageChart[i].Weight < ageChart[j].Weight
	})
	receiver.SuccReturn(map[string]interface{}{
		"total_people":       roomUserTotal,
		"age_people":         roomAgePeopleTotal,
		"total_user":         liveInfo.TotalUser,
		"inc_fans":           liveInfo.FollowCount,
		"barrage_count":      liveInfo.BarrageCount,
		"barrage_user_count": liveInfo.BarrageUserCount,
		"barrage_rate":       barrageRate,
		"word_chart":         wordChart,
		"gender_chart":       genderChart,
		"age_chart":          ageChart,
		"city_chart":         cityChart,
		"province_chart":     provinceChart,
	})
	return
}

func (receiver *LiveController) liveFansAgeMap(key string) string {
	var ageMap = map[string]string{
		"小于18":  "-18",
		"18~23": "18-23",
		"24~30": "24-30",
		"31~40": "31-40",
		"41~50": "41-50",
		"50+":   "50-",
	}
	if s, ok := ageMap[key]; ok {
		return s
	}
	return key
}

//直播商品分类意向
func (receiver *LiveController) LiveProductPvAnalyse() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	liveInfo, comErr := hbase.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	cateChart := make([]dy2.ProductPvChart, 0)
	cacheKey := cache.GetCacheKey(cache.LiveRoomProductList, roomId)
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &cateChart)
		receiver.SuccReturn(map[string]interface{}{
			"list": cateChart,
		})
		return
	}
	list, total, comErr := es.NewEsLiveBusiness().GetProductByRoomId(liveInfo)
	cateMap := map[string]dy2.ProductPvChartMap{}
	var totalPv int64 = 0
	if total > 0 {
		stopRowKey := list[0].RoomID + "_" + list[0].ProductID
		startRowKey := list[total-1].RoomID + "_" + list[total-1].ProductID
		mapData := map[string]entity.DyRoomProduct{}
		if startRowKey != stopRowKey {
			mapData, _ = hbase.GetRoomProductInfoRangDate(startRowKey, stopRowKey)
		}
		stopData, err := hbase.GetRoomProductInfo(stopRowKey)
		if err == nil {
			mapData[list[total-1].ProductID] = stopData
		}
		for _, v := range list {
			var ptmPv int64 = 0
			if productData, ok := mapData[v.ProductID]; ok {
				for _, p := range productData.PtmPromotion {
					if p.FinalPv > 0 {
						ptmPv += p.FinalPv - p.InitialPv
					} else {
						ptmPv += p.Pv - p.InitialPv
					}
				}
				totalPv += ptmPv
			}
			if v.DcmLevelFirst == "" {
				v.DcmLevelFirst = "其他"
			}
			if c, ok := cateMap[v.DcmLevelFirst]; !ok {
				cateMap[v.DcmLevelFirst] = dy2.ProductPvChartMap{
					Pv:  ptmPv,
					Son: map[string]dy2.ProductPvChartMap{},
				}
			} else {
				c.Pv += ptmPv
				cateMap[v.DcmLevelFirst] = c
			}
			if v.DcmLevelFirst == "其他" {
				continue
			}
			if c, ok := cateMap[v.DcmLevelFirst].Son[v.FirstCname]; !ok {
				cateMap[v.DcmLevelFirst].Son[v.FirstCname] = dy2.ProductPvChartMap{
					Pv:  ptmPv,
					Son: map[string]dy2.ProductPvChartMap{},
				}
			} else {
				c.Pv += ptmPv
				cateMap[v.DcmLevelFirst].Son[v.FirstCname] = c
			}
			if c, ok := cateMap[v.DcmLevelFirst].Son[v.FirstCname].Son[v.SecondCname]; !ok {
				cateMap[v.DcmLevelFirst].Son[v.FirstCname].Son[v.SecondCname] = dy2.ProductPvChartMap{
					Pv:  ptmPv,
					Son: map[string]dy2.ProductPvChartMap{},
				}
			} else {
				c.Pv += ptmPv
				cateMap[v.DcmLevelFirst].Son[v.FirstCname].Son[v.SecondCname] = c
			}
		}
	}
	cateChart = receiver.getSon(cateMap, totalPv)
	cateChart = receiver.sortSon(cateChart)
	var cacheTime time.Duration = 60
	if liveInfo.RoomStatus == 4 {
		cacheTime = 1800
	}
	_ = global.Cache.Set(cacheKey, utils.SerializeData(cateChart), cacheTime)
	receiver.SuccReturn(map[string]interface{}{
		"list": cateChart,
	})
	return
}

//循环遍历
func (receiver *LiveController) getSon(son map[string]dy2.ProductPvChartMap, pv int64) []dy2.ProductPvChart {
	data := make([]dy2.ProductPvChart, 0)
	for k, v := range son {
		item := dy2.ProductPvChart{
			LabelName: k,
			Pv:        v.Pv,
			LabelSon:  []dy2.ProductPvChart{},
		}
		if pv > 0 {
			item.Percent = float64(item.Pv) / float64(pv)
		}
		if len(v.Son) > 0 {
			item.LabelSon = receiver.getSon(v.Son, v.Pv)
		} else {
			item.LabelSon = []dy2.ProductPvChart{}
		}
		data = append(data, item)
	}
	return data
}

func (receiver *LiveController) sortSon(data []dy2.ProductPvChart) []dy2.ProductPvChart {
	sort.Slice(data, func(i, j int) bool {
		if data[i].LabelName == "其他" {
			return false
		}
		if data[j].LabelName == "其他" {
			return true
		}
		return data[i].Pv > data[j].Pv
	})
	if len(data) > 5 {
		data = data[:5]
	}
	for k, v := range data {
		if len(v.LabelSon) > 0 {
			data[k].LabelSon = receiver.sortSon(v.LabelSon)
		}
	}
	return data
}

//数据大屏基础数据
func (receiver *LiveController) LivingBaseData() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	liveInfo, comErr := hbase.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	authorData, _ := hbase.GetAuthor(liveInfo.User.ID)
	livingInfo := dy2.LivingInfo{
		RoomId:   business.IdEncrypt(liveInfo.RoomID),
		AuthorId: business.IdEncrypt(liveInfo.User.ID),
		Author: dy2.LivingAuthorInfo{
			Avatar:        dyimg.Fix(authorData.Data.Avatar),
			Nickname:      authorData.Data.Nickname,
			FollowerCount: authorData.Data.FollowerCount,
			RoomId:        business.IdEncrypt(authorData.RoomId),
		},
		Title:        liveInfo.Title,
		Cover:        dyimg.Fix(liveInfo.Cover),
		CreateTime:   liveInfo.CreateTime,
		RoomShareUrl: business.LiveShareUrl + roomId,
	}
	receiver.SuccReturn(livingInfo)
	return
}

//数据大屏销售数据
func (receiver *LiveController) LivingSaleData() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	liveInfo, comErr := hbase.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	var gmv = liveInfo.PredictGmv
	//liveSaleData, _ := hbase.GetLiveSalesData(roomId)
	//if liveInfo.RoomStatus == 4 {
	//	if liveSaleData.Gmv > 0 {
	//		gmv = liveSaleData.Gmv
	//	}
	//}
	livingInfo := dy2.LivingSale{
		RoomId:         business.IdDecrypt(liveInfo.RoomID),
		CreateTime:     liveInfo.CreateTime,
		Gmv:            gmv,
		UserCount:      liveInfo.UserCount,
		TotalUserCount: liveInfo.TotalUser,
		RoomStatus:     liveInfo.RoomStatus,
		FinishTime:     liveInfo.FinishTime,
	}
	if liveInfo.FinishTime > 0 {
		livingInfo.LiveTime = liveInfo.FinishTime - liveInfo.CreateTime
	} else {
		livingInfo.LiveTime = time.Now().Unix() - liveInfo.CreateTime
	}
	if liveInfo.TotalUser > 0 {
		livingInfo.Uv = utils.RateMin((gmv + float64(liveInfo.RoomTicketCount)/10) / float64(liveInfo.TotalUser))
		livingInfo.BarrageRate = utils.RateMin(float64(liveInfo.BarrageUserCount) / float64(liveInfo.TotalUser))
	}
	livingInfo.AvgOnlineTime = business.NewLiveBusiness().CountAvgOnlineTime(liveInfo.OnlineTrends, liveInfo.CreateTime, liveInfo.TotalUser)
	receiver.SuccReturn(livingInfo)
	return
}

//数据大屏观看趋势数据
func (receiver *LiveController) LivingWatchChart() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	liveInfo, comErr := hbase.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	incOnlineTrends, _, _ := business.NewLiveBusiness().DealOnlineTrends(liveInfo)
	receiver.SuccReturn(map[string]interface{}{
		"room_id": business.IdEncrypt(roomId),
		"trends":  incOnlineTrends,
	})
	return
}

//数据大屏商品数据
func (receiver *LiveController) LivingProduct() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	sortStr := receiver.GetString("sort", "")
	orderBy := receiver.GetString("order_by", "")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	liveInfo, comErr := hbase.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	esLiveBusiness := es.NewEsLiveBusiness()
	sales, total := esLiveBusiness.SumRoomProductByRoomId(liveInfo)
	originalList, total, comErr := esLiveBusiness.LivingProductList(liveInfo, sortStr, orderBy, page, pageSize)
	list := make([]dy2.LivingProducts, 0)
	utils.MapToStruct(originalList, &list)
	liveBusiness := business.NewLiveBusiness()
	for k, v := range list {
		list[k].RoomID = business.IdEncrypt(v.RoomID)
		list[k].ProductID = business.IdEncrypt(v.ProductID)
		list[k].AuthorID = business.IdEncrypt(v.AuthorID)
		if v.IsReturn == 1 && v.StartTime == v.ShelfTime {
			list[k].IsReturn = 0
		}
		list[k].Cover = dyimg.Product(v.Cover)
		list[k].PredictSales = math.Floor(v.PredictSales)
		list[k].CurList = []dy2.LiveCurProduct{}
		curCount, pmtStatus, pv, err := liveBusiness.RoomCurAndPmtProductById(roomId, v.ProductID)
		if err == nil {
			v.Pv = pv
			if v.Pv > 0 {
				list[k].BuyRate = utils.RateMin(v.PredictSales / float64(v.Pv))
			}
			list[k].CurSecond = curCount.CurSecond
			pmtStatusLen := len(pmtStatus)
			if pmtStatusLen > 0 {
				startPmt := pmtStatus[0]
				stopPmt := pmtStatus[pmtStatusLen-1]
				list[k].StartPmtSales = startPmt.StartSales
				list[k].EndPmtSales = stopPmt.FinalSales
			}
			if len(curCount.CurList) > 0 {
				list[k].CurList = curCount.CurList
				cur := curCount.CurList[len(curCount.CurList)-1]
				list[k].StartCurTime = cur.StartTime
				list[k].EndCurTime = cur.EndTime
			}
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"sales": sales,
		"total": total,
	})
	return
}

//直播间弹幕
func (receiver LiveController) LivingMessage() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	beginNum, _ := receiver.GetInt64("begin", 0)
	visitNum, _ := receiver.GetInt64("visit_begin", 0)
	pageSize := receiver.GetPageSize("page_size", 30, 200)
	data, _ := hbase.GetLiveChatMessage(roomId)
	list := make([]entity.LivingChatMessage, 0)
	visitList := make([]entity.LivingChatVisit, 0)
	lenNum := len(data.Latest500Msg)
	visitLenNum := len(data.Visits)
	var endNum int64 = 0
	var endVisitNum int64 = 0
	if data.EndNum > beginNum && lenNum > 0 {
		firstId := data.Latest500Msg[0].RankId
		if beginNum <= firstId {
			beginNum = 0
		}
		start := 0
		end := 0
		if beginNum > 0 {
			start = int(beginNum-firstId) + 1
		} else {
			start = lenNum - pageSize
			if start < 0 {
				start = 0
			}
		}
		end = start + pageSize
		if lenNum < end {
			end = lenNum
		}
		if start > lenNum {
			start = lenNum
		}
		if lenNum > 0 {
			list = data.Latest500Msg[start:end]
		}
		lastKey := -1
		for k, v := range list {
			list[k].Avatar = dyimg.Fix(v.Avatar)
			lastKey = k
		}
		if lastKey < 0 {
			endNum = beginNum
		} else {
			endNum = list[lastKey].RankId
		}
	}
	if data.VisitNum > visitNum && visitLenNum > 0 {
		firstId := data.Visits[0].RankId
		if visitNum <= firstId {
			visitNum = 0
		}
		start := 0
		end := 0
		if visitNum > 0 {
			start = int(visitNum-firstId) + 1
		} else {
			start = visitLenNum - pageSize
			if start < 0 {
				start = 0
			}
		}
		end = start + pageSize
		if visitLenNum < end {
			end = visitLenNum
		}
		if start > visitLenNum {
			start = visitLenNum
		}
		if visitLenNum > 0 {
			visitList = data.Visits[start:end]
		}
		lastKey := -1
		for k := range visitList {
			lastKey = k
		}
		if lastKey < 0 {
			endVisitNum = visitNum
		} else {
			endVisitNum = visitList[lastKey].RankId
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":          list,
		"end_num":       endNum,
		"visit_list":    visitList,
		"visit_end_num": endVisitNum,
	})
	return
}

//直播加速
func (receiver *LiveController) LiveSpeed() {
	isScreen, _ := receiver.GetInt("is_screen", 0)
	if !business.UserActionLock(receiver.TrueUri, utils.ToString(receiver.UserId), 5) {
		receiver.FailReturn(global.NewError(6000))
		return
	}

	AuthorId := business.IdDecrypt(receiver.GetString(":author_id", ""))
	if AuthorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	spriderName := "live"
	cacheKey := cache.GetCacheKey(cache.SpiderSpeedUpLimit, spriderName, AuthorId)
	cacheData := global.Cache.Get(cacheKey)
	if cacheData != "" {
		//缓存存在
		receiver.FailReturn(global.NewError(6000))
		return
	}
	//加速
	top := business.AddLiveTopConcerned
	expireTime := time.Now().AddDate(0, 0, 7).Unix()
	if isScreen == 1 {
		top = business.AddLiveTopHighLevelStar
		expireTime = time.Now().AddDate(0, 0, 1).Unix()
	}
	author, _ := hbase.GetAuthor(AuthorId)
	go business.NewSpiderBusiness().AddLive(AuthorId, author.FollowerCount, top, expireTime)
	global.Cache.Set(cacheKey, "1", 300)
	receiver.SuccReturn([]string{})
	return
}
