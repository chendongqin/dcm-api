package dy

import (
	"dongchamao/business"
	"dongchamao/business/es"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	dy2 "dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	"math"
	"sort"
	"time"
)

type LiveController struct {
	controllers.ApiBaseController
}

func (receiver *LiveController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseShowNum)
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
	if !receiver.HasAuth {
		today := time.Now().Format("20060102")
		lastDay := time.Now().AddDate(0, 0, -6).Format("20060102")
		start := startTime.Format("20060102")
		end := endTime.Format("20060102")
		if lastDay != start || today != end || keyword != "" || category != "" || sortStr != "predict_gmv" || orderBy != "desc" || minAmount > 0 || maxAmount > 0 || minUv > 0 || maxUv > 0 || minAvgUserCount > 0 || maxAvgUserCount > 0 || hasProduct == 1 || isBrand == 1 || page != 1 {
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
	liveUser := dy2.DyLiveUserSimple{
		Avatar:          liveInfo.User.Avatar,
		FollowerCount:   authorInfo.Data.FollowerCount,
		ID:              business.IdEncrypt(liveInfo.User.ID),
		Nickname:        liveInfo.User.Nickname,
		WithCommerce:    liveInfo.User.WithCommerce,
		ReputationScore: reputation.AuthorReputation.Score,
		ReputationLevel: reputation.AuthorReputation.Level,
	}
	liveSaleData, _ := hbase.GetLiveSalesData(roomId)
	incOnlineTrends, maxOnlineTrends, avgUserCount := liveBusiness.DealOnlineTrends(liveInfo)
	var incFansRate, interactRate float64
	incFansRate = 0
	interactRate = 0
	liveSale := dy2.DyLiveRoomSaleData{}
	//todo gmv数据兼容
	gmv := liveSaleData.Gmv
	sales := liveSaleData.Sales
	if liveSaleData.Gmv == 0 {
		gmv = liveInfo.PredictGmv
		sales = liveInfo.PredictSales
		//if liveInfo.RealGmv > 0 {
		//	gmv = liveInfo.RealGmv
		//	sales = liveInfo.RealSales
		//}
	}
	if liveInfo.TotalUser > 0 {
		incFansRate = float64(liveInfo.FollowCount) / float64(liveInfo.TotalUser)
		interactRate = float64(liveInfo.BarrageCount) / float64(liveInfo.TotalUser)
		liveSale.Uv = (gmv + float64(liveSaleData.TicketCount)/10) / float64(liveInfo.TotalUser)
		liveSale.SaleRate = sales / float64(liveInfo.TotalUser)
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
		IncFansRate:         incFansRate,
		InteractRate:        interactRate,
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
	dateChart := make([]int64, 0)
	gmvChart := make([]float64, 0)
	salesChart := make([]float64, 0)
	salesTrends := liveInfo.SalesTrends
	//排序
	sort.Slice(salesTrends, func(i, j int) bool {
		var left, right int64
		left = salesTrends[i].CrawlTime
		right = salesTrends[j].CrawlTime
		return right > left
	})
	for _, v := range salesTrends {
		dateChart = append(dateChart, v.CrawlTime)
		//if liveInfo.RealGmv > 0 {
		//	gmvChart = append(gmvChart, v.RealGmv)
		//	salesChart = append(salesChart, math.Floor(v.RealSales))
		//} else {
		gmvChart = append(gmvChart, v.PredictGmv)
		salesChart = append(salesChart, math.Floor(v.PredictSales))
		//}
	}
	receiver.SuccReturn(map[string]interface{}{
		"live_info": returnLiveInfo,
		"live_sale": liveSale,
		"sales_chart": map[string]interface{}{
			"time":  dateChart,
			"gmv":   gmvChart,
			"sales": salesChart,
		},
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
	}
	receiver.SuccReturn(map[string]interface{}{
		"promotions_list": promotionsList,
	})
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
	sortStr := InputData.GetString("sort", "shelf_time")
	orderBy := InputData.GetString("order_by", "desc")
	page := InputData.GetInt("page", 1)
	pageSize := InputData.GetInt("page_size", 10)
	firstLabel := InputData.GetString("first_label", "")
	secondLabel := InputData.GetString("second_label", "")
	thirdLabel := InputData.GetString("third_label", "")
	roomInfo, _ := hbase.GetLiveInfo(roomId)
	esLiveBusiness := es.NewEsLiveBusiness()
	list, productCount, total, err := esLiveBusiness.RoomProductByRoomId(roomInfo, keyword, sortStr, orderBy, firstLabel, secondLabel, thirdLabel, page, pageSize)
	if err != nil {
		receiver.FailReturn(err)
		return
	}
	countList := make([]dy2.LiveRoomProductCount, 0)
	if len(list) > 0 {
		productIds := make([]string, 0)
		for _, v := range list {
			productIds = append(productIds, v.ProductID)
		}
		liveBusiness := business.NewLiveBusiness()
		//curMap := liveBusiness.RoomCurProductByIds(roomId, productIds)
		//pmtMap := liveBusiness.RoomPmtProductByIds(roomId, productIds)
		for _, v := range list {
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
			curCount, pmtStatus, err1 := liveBusiness.RoomCurAndPmtProductById(roomId, v.ProductID)
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
	info, _ := hbase.GetRoomProductInfo(roomId + "_" + productId)
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

//数据大屏基础数据
func (receiver *LiveController) LivingBaseData() {
	roomId := business.IdDecrypt(receiver.Ctx.Input.Param(":room_id"))
	liveInfo, comErr := hbase.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	var gmv = liveInfo.PredictGmv
	liveSaleData, _ := hbase.GetLiveSalesData(roomId)
	if liveInfo.RoomStatus == 4 {
		if liveSaleData.Gmv > 0 {
			gmv = liveSaleData.Gmv
		}
	}
	authorData, _ := hbase.GetAuthor(liveInfo.User.ID)
	livingInfo := dy2.LivingInfo{
		RoomId:   business.IdDecrypt(liveInfo.RoomID),
		AuthorId: business.IdDecrypt(liveInfo.User.ID),
		Author: dy2.LivingAuthorInfo{
			Avatar:        dyimg.Fix(authorData.Data.Avatar),
			Nickname:      authorData.Data.Nickname,
			FollowerCount: authorData.Data.FollowerCount,
			RoomId:        business.IdDecrypt(authorData.RoomId),
		},
		Title:          liveInfo.Title,
		Cover:          dyimg.Fix(liveInfo.Cover),
		CreateTime:     liveInfo.CreateTime,
		Gmv:            gmv,
		UserCount:      liveInfo.UserCount,
		TotalUserCount: liveInfo.TotalUser,
		RoomStatus:     liveInfo.RoomStatus,
		FinishTime:     liveInfo.FinishTime,
		RoomShareUrl:   business.LiveShareUrl + roomId,
	}
	if liveInfo.FinishTime > 0 {
		livingInfo.LiveTime = liveInfo.FinishTime - liveInfo.CreateTime
	} else {
		livingInfo.LiveTime = time.Now().Unix() - liveInfo.CreateTime
	}
	if liveInfo.TotalUser > 0 {
		livingInfo.Uv = (gmv + float64(liveSaleData.TicketCount)/10) / float64(liveInfo.TotalUser)
		livingInfo.BarrageRate = float64(liveInfo.BarrageCount) / float64(liveInfo.TotalUser)
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
		"gmv":     liveInfo.PredictGmv,
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
		if v.Pv > 0 {
			list[k].BuyRate = v.PredictSales / float64(v.Pv)
		}
		list[k].CurList = []dy2.LiveCurProduct{}
		curCount, pmtStatus, err := liveBusiness.RoomCurAndPmtProductById(roomId, v.ProductID)
		if err == nil {
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
