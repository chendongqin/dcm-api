package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/models/business"
	"dongchamao/models/business/es"
	"dongchamao/services/dyimg"
	"dongchamao/structinit/repost/dy"
	"math"
	"sort"
	"time"
)

type LiveController struct {
	controllers.ApiBaseController
}

func (receiver *LiveController) LiveInfoData() {
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	liveBusiness := business.NewLiveBusiness()
	liveInfo, comErr := liveBusiness.HbaseGetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	reputation, _ := authorBusiness.HbaseGetAuthorReputation(liveInfo.User.ID)
	authorInfo, _ := authorBusiness.HbaseGetAuthor(liveInfo.User.ID)
	liveUser := dy.DyLiveUserSimple{
		Avatar:          liveInfo.User.Avatar,
		FollowerCount:   authorInfo.FollowerCount,
		ID:              liveInfo.User.ID,
		Nickname:        liveInfo.User.Nickname,
		WithCommerce:    liveInfo.User.WithCommerce,
		ReputationScore: reputation.Score,
		ReputationLevel: reputation.Level,
	}
	liveSaleData, _ := liveBusiness.HbaseGetLiveSalesData(roomId)
	incOnlineTrends, maxOnlineTrends, avgUserCount := liveBusiness.DealOnlineTrends(liveInfo)
	var incFansRate, interactRate float64
	incFansRate = 0
	interactRate = 0
	liveSale := dy.DyLiveRoomSaleData{}
	if liveInfo.TotalUser > 0 {
		incFansRate = float64(liveInfo.FollowCount) / float64(liveInfo.TotalUser)
		interactRate = float64(liveInfo.BarrageCount) / float64(liveInfo.TotalUser)
		liveSale.Uv = (liveSaleData.Gmv + float64(liveSaleData.TicketCount)/10) / float64(liveInfo.TotalUser)
		liveSale.SaleRate = liveSaleData.Gmv / float64(liveInfo.TotalUser)
	}
	avgOnlineTime := liveBusiness.CountAvgOnlineTime(liveInfo.OnlineTrends, liveInfo.CreateTime, liveInfo.TotalUser)
	returnLiveInfo := dy.DyLiveInfo{
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
		ShareUrl:            business.LiveShareUrl + liveInfo.RoomID,
	}
	liveSale.Volume = int64(math.Floor(liveSaleData.Sales))
	liveSale.Amount = liveSaleData.Gmv
	liveSale.PromotionNum = liveSaleData.NumProducts
	if liveSaleData.Sales > 0 {
		liveSale.PerPrice = liveSaleData.Gmv / liveSaleData.Sales
	}
	receiver.SuccReturn(map[string]interface{}{
		"live_info": returnLiveInfo,
		"live_sale": liveSale,
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
	liveBusiness := business.NewLiveBusiness()
	livePmt, _ := liveBusiness.HbaseGetLivePmt(roomId)
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
	dyLivePromotions := make([][]dy.DyLivePromotion, 0)
	promotionSales := map[string]int{}
	for k, v := range promotionsMap {
		item := make([]dy.DyLivePromotion, 0)
		for _, v1 := range v {
			saleNum := 1
			if s, ok := promotionSales[v1.ProductID]; ok {
				saleNum = s + 1
			}
			item = append(item, dy.DyLivePromotion{
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
	promotionsList := dy.DyLivePromotionChart{
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
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	InputData := receiver.InputFormat()
	keyword := InputData.GetString("keyword", "")
	sortStr := InputData.GetString("sort", "start_time")
	orderBy := InputData.GetString("order_by", "desc")
	page := InputData.GetInt("page", 1)
	pageSize := InputData.GetInt("page_size", 10)
	firstLabel := InputData.GetString("first_label", "")
	secondLabel := InputData.GetString("second_label", "")
	thirdLabel := InputData.GetString("third_label", "")
	esLiveBusiness := es.NewEsLiveBusiness()
	list, productCount, total, err := esLiveBusiness.RoomProductByRoomId(roomId, keyword, sortStr, orderBy, firstLabel, secondLabel, thirdLabel, page, pageSize)
	if err != nil {
		receiver.FailReturn(err)
		return
	}
	countList := make([]dy.LiveRoomProductCount, 0)
	if len(list) > 0 {
		productIds := make([]string, 0)
		for _, v := range list {
			productIds = append(productIds, v.ProductID)
		}
		liveBusiness := business.NewLiveBusiness()
		curMap := liveBusiness.RoomCurProductByIds(roomId, productIds)
		pmtMap := liveBusiness.RoomPmtProductByIds(roomId, productIds)
		for _, v := range list {
			item := dy.LiveRoomProductCount{
				ProductInfo: v,
				ProductStartSale: dy.RoomProductSaleChart{
					Timestamp: []int64{},
					Sales:     []int64{},
				},
				ProductEndSale: dy.RoomProductSaleChart{
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
				c.CurList = business.ProductCurOrderByTime(c.CurList)
				item.ProductCur = c
			} else {
				item.ProductCur = dy.LiveCurProductCount{
					CurList: []dy.LiveCurProduct{},
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

//直播间商品
func (receiver *LiveController) LiveProductCateList() {
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	esLiveBusiness := es.NewEsLiveBusiness()
	countData := esLiveBusiness.AllRoomProductCateByRoomId(roomId)
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
	liveBusiness := business.NewLiveBusiness()
	info, _ := liveBusiness.RoomCurProductSaleTrend(roomId, productId)
	trends := business.RoomProductTrendOrderByTime(info.TrendData)
	timestamps := make([]int64, 0)
	sales := make([]float64, 0)
	for _, v := range trends {
		timestamps = append(timestamps, v.CrawlTime)
		sales = append(sales, math.Floor(v.Sales))
	}
	receiver.SuccReturn(dy.TimestampCountChart{
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
	liveBusiness := business.NewLiveBusiness()
	info, comErr := liveBusiness.HbaseGetLiveInfo(roomId)
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
		//lenNum := len(followerCountTrends)
		//beforeFansTrend := entity.LiveFollowerCountTrends{
		//	CrawlTime:     info.CreateTime,
		//	FollowerCount: followerCountTrends[lenNum-1].FollowerCount - info.FollowCount,
		//	NewFollowerCount: 0,
		//}
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
			FansClubCount: fansClubCountTrends[0].FansClubCount,
			CrawlTime:     info.CreateTime,
		}
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
