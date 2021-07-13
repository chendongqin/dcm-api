package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/models/business"
	"dongchamao/services/dyimg"
	"dongchamao/structinit/repost/dy"
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
	}
	incOnlineTrends, maxOnlineTrends, incFans := liveBusiness.DealOnlineTrends(liveInfo.OnlineTrends)
	incFansRate := float64(incFans) / float64(liveInfo.TotalUser)
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
		InteractRate:        0,
		MaxWatchOnlineTrend: maxOnlineTrends,
		OnlineTrends:        incOnlineTrends,
	}
	receiver.SuccReturn(map[string]interface{}{
		"live_info": returnLiveInfo,
	})
	return
}

//
func (receiver *LiveController) LivePmt() {
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	liveBusiness := business.NewLiveBusiness()
	livePmt, comErr := liveBusiness.HbaseGetLivePmt(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
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
	for k, v := range promotionsMap {
		item := make([]dy.DyLivePromotion, 0)
		for _, v1 := range v {
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
