package dy

import (
	business2 "dongchamao/business"
	es2 "dongchamao/business/es"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	hbase2 "dongchamao/hbase"
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

//直播详细
func (receiver *LiveController) LiveInfoData() {
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	liveBusiness := business2.NewLiveBusiness()
	liveInfo, comErr := hbase2.GetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	authorBusiness := business2.NewAuthorBusiness()
	reputation, _ := hbase2.GetLiveReputation(roomId)
	authorInfo, _ := authorBusiness.HbaseGetAuthor(liveInfo.User.ID)
	liveUser := dy2.DyLiveUserSimple{
		Avatar:          liveInfo.User.Avatar,
		FollowerCount:   authorInfo.FollowerCount,
		ID:              liveInfo.User.ID,
		Nickname:        liveInfo.User.Nickname,
		WithCommerce:    liveInfo.User.WithCommerce,
		ReputationScore: reputation.AuthorReputation.Score,
		ReputationLevel: reputation.AuthorReputation.Level,
	}
	liveSaleData, _ := hbase2.GetLiveSalesData(roomId)
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
		if liveInfo.RealGmv > 0 {
			gmv = liveInfo.RealGmv
			sales = liveInfo.RealSales
		}
	}
	if liveInfo.TotalUser > 0 {
		incFansRate = float64(liveInfo.FollowCount) / float64(liveInfo.TotalUser)
		interactRate = float64(liveInfo.BarrageCount) / float64(liveInfo.TotalUser)
		liveSale.Uv = (gmv + float64(liveSaleData.TicketCount)/10) / float64(liveInfo.TotalUser)
		liveSale.SaleRate = gmv / float64(liveInfo.TotalUser)
	}
	avgOnlineTime := liveBusiness.CountAvgOnlineTime(liveInfo.OnlineTrends, liveInfo.CreateTime, liveInfo.TotalUser)
	returnLiveInfo := dy2.DyLiveInfo{
		Cover:               liveInfo.Cover,
		CreateTime:          liveInfo.CreateTime,
		FinishTime:          liveInfo.FinishTime,
		LikeCount:           liveInfo.LikeCount,
		RoomID:              liveInfo.RoomID,
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
		ShareUrl:            business2.LiveShareUrl + liveInfo.RoomID,
	}
	liveSale.Volume = int64(math.Floor(sales))
	liveSale.Amount = gmv
	esLiveBusiness := es2.NewEsLiveBusiness()
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
		if liveInfo.RealGmv > 0 {
			gmvChart = append(gmvChart, v.RealGmv)
			salesChart = append(salesChart, math.Floor(v.RealSales))
		} else {
			gmvChart = append(gmvChart, v.PredictGmv)
			salesChart = append(salesChart, math.Floor(v.PredictSales))
		}
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
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	livePmt, _ := hbase2.GetLivePmt(roomId)
	livePromotionsMap := map[int]entity.DyLivePromotion{}
	for _, v := range livePmt.Promotions {
		livePromotionsMap[v.Index] = v
	}
	var keys []int
	for k := range livePromotionsMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	promotionsMap := map[string][]entity.DyLivePromotion{}
	for _, k := range keys {
		if v, ok := livePromotionsMap[k]; ok {
			startFormat := time.Unix(v.StartTime, 0).Format("2006-01-02 15:04:05")
			if _, ok1 := promotionsMap[startFormat]; !ok1 {
				promotionsMap[startFormat] = make([]entity.DyLivePromotion, 0)
			}
			promotionsMap[startFormat] = append(promotionsMap[startFormat], v)
		}
	}
	dates := make([]string, 0)
	dyLivePromotions := make([][]dy2.DyLivePromotion, 0)
	promotionSales := map[string]int{}
	for k, v := range promotionsMap {
		item := make([]dy2.DyLivePromotion, 0)
		for _, v1 := range v {
			saleNum := 1
			if s, ok := promotionSales[v1.ProductID]; ok {
				saleNum = s + 1
			}
			item = append(item, dy2.DyLivePromotion{
				ProductID: v1.ProductID,
				ForSale:   v1.ForSale,
				StartTime: v1.StartTime,
				StopTime:  v1.StopTime,
				Price:     v1.Price,
				Sales:     v1.Sales,
				NowSales:  0,
				GmvSales:  0,
				Title:     v1.Title,
				Cover:     dyimg.Product(v1.Cover),
				Index:     v1.Index,
				SaleNum:   saleNum,
			})
		}
		dyLivePromotions = append(dyLivePromotions, item)
		dates = append(dates, k)
	}
	promotionsList := dy2.DyLivePromotionChart{
		StartTime:     dates,
		PromotionList: dyLivePromotions,
	}
	receiver.SuccReturn(map[string]interface{}{
		"promotions_list": promotionsList,
	})
}

//直播榜单排名趋势
func (receiver *LiveController) LiveRankTrends() {
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	liveBusiness := business2.NewLiveBusiness()
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
	hourDates = business2.DealChartInt64(hourDates, 60)
	hourRanks = business2.DealChartInt(hourRanks, 60)
	saleDates = business2.DealChartInt64(saleDates, 60)
	saleRanks = business2.DealChartInt(saleRanks, 60)
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
	roomId := receiver.Ctx.Input.Param(":room_id")
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
	roomInfo, _ := hbase2.GetLiveInfo(roomId)
	esLiveBusiness := es2.NewEsLiveBusiness()
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
		liveBusiness := business2.NewLiveBusiness()
		curMap := liveBusiness.RoomCurProductByIds(roomId, productIds)
		pmtMap := liveBusiness.RoomPmtProductByIds(roomId, productIds)
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
			if s, ok := pmtMap[v.ProductID]; ok {
				for _, s1 := range s {
					item.ProductStartSale.Timestamp = append(item.ProductStartSale.Timestamp, s1.StartTime)
					item.ProductStartSale.Sales = append(item.ProductStartSale.Sales, s1.StartSales)
					if s1.StopTime > 0 {
						item.ProductEndSale.Timestamp = append(item.ProductEndSale.Timestamp, s1.StopTime)
						item.ProductEndSale.Sales = append(item.ProductEndSale.Sales, s1.FinalSales)
					}
				}
			}
			if c, ok := curMap[v.ProductID]; ok {
				c.CurList = business2.ProductCurOrderByTime(c.CurList)
				item.ProductCur = c
			} else {
				item.ProductCur = dy2.LiveCurProductCount{
					CurList: []dy2.LiveCurProduct{},
				}
			}
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
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	roomInfo, _ := hbase2.GetLiveInfo(roomId)
	esLiveBusiness := es2.NewEsLiveBusiness()
	countData := esLiveBusiness.AllRoomProductCateByRoomId(roomInfo)
	receiver.SuccReturn(map[string]interface{}{
		"count": countData,
	})
	return
}

//全网销量趋势图
func (receiver *LiveController) LiveProductSaleChart() {
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	productId := receiver.Ctx.Input.Param(":product_id")
	if productId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	info, _ := hbase2.GetRoomProductInfo(roomId, productId)
	trends := business2.RoomProductTrendOrderByTime(info.TrendData)
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
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	info, comErr := hbase2.GetLiveInfo(roomId)
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
		followerCountTrends := business2.LiveFansTrendsListOrderByTime(info.FollowerCountTrends)
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
		fansClubCountTrends := business2.LiveClubFansTrendsListOrderByTime(info.FansClubCountTrends)
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
