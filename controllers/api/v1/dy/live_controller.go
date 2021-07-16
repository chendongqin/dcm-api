package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/models/business"
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
	liveUser := dy.DyLiveUserSimple{
		Avatar:          liveInfo.User.Avatar,
		FollowerCount:   liveInfo.User.FollowerCount,
		ID:              liveInfo.User.ID,
		Nickname:        liveInfo.User.Nickname,
		WithCommerce:    liveInfo.User.WithCommerce,
		ReputationScore: reputation.Score,
		ReputationLevel: reputation.Level,
	}
	liveSaleData, _ := liveBusiness.HbaseGetLiveSalesData(roomId)
	incOnlineTrends, maxOnlineTrends, incFans, avgUserCount := liveBusiness.DealOnlineTrends(liveInfo.OnlineTrends)
	var incFansRate, interactRate float64
	incFansRate = 0
	interactRate = 0
	liveSale := dy.DyLiveRoomSaleData{}
	if liveInfo.TotalUser > 0 {
		incFansRate = float64(incFans) / float64(liveInfo.TotalUser)
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
		IncFans:             incFans,
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
	for _, v := range liveRankTrends {
		if v.Type == 8 {
			saleDates = append(saleDates, v.CrawlTime)
			saleRanks = append(saleRanks, v.Rank)
		} else if v.Type == 1 {
			hourDates = append(hourDates, v.CrawlTime)
			hourRanks = append(hourRanks, v.Rank)
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"hour_rank": map[string]interface{}{
			"time":  hourDates,
			"ranks": hourRanks,
		},
		"sale_rank": map[string]interface{}{
			"time":  saleDates,
			"ranks": saleRanks,
		},
	})
}
