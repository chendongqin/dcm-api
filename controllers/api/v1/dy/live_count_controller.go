package dy

import (
	"dongchamao/business"
	"dongchamao/business/es"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/utils"
	es2 "dongchamao/models/es"
	"dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	"sort"
	"strings"
)

type LiveCountController struct {
	controllers.ApiBaseController
}

func (receiver *LiveCountController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseShowNum)
}

//直播总览
func (receiver *LiveCountController) AllLiveCount() {
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	living, _ := receiver.GetInt("living", 0)
	esLiveDataBusiness := es.NewEsLiveDataBusiness()
	allTotal, allData := esLiveDataBusiness.SumLiveData(startTime, endTime, 0, living)
	productTotal, productData := esLiveDataBusiness.SumLiveData(startTime, endTime, 1, living)
	receiver.SuccReturn(map[string]interface{}{
		"all_live": map[string]interface{}{
			"total":            allTotal,
			"total_watch":      allData.TotalWatchCnt.Value,
			"total_user_count": allData.TotalUserCount.Value,
		},
		"product_live": map[string]interface{}{
			"total":            productTotal,
			"total_watch":      productData.TotalWatchCnt.Value,
			"total_user_count": productData.TotalUserCount.Value,
		},
	})
	return
}

//直播分类榜单
func (receiver *LiveCountController) LiveCategoryRank() {
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	category := receiver.GetString("category", "")
	if category == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	living, _ := receiver.GetInt("living", 0)
	sortStr := receiver.GetString("sort", "")
	esLiveDataBusiness := es.NewEsLiveDataBusiness()
	list, comErr := esLiveDataBusiness.LiveRankByCategory(startTime, endTime, category, sortStr, living)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		list[k].Avatar = dyimg.Fix(v.Avatar)
		list[k].Cover = dyimg.Fix(v.Cover)
		list[k].RoomId = business.IdEncrypt(v.RoomId)
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		if v.DisplayId == "" {
			list[k].DisplayId = v.ShortId
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"list": list,
	})
	return
}

//直播分类占比
func (receiver *LiveCountController) LiveCompositeByCategory() {
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	rateType := utils.ToInt(receiver.Ctx.Input.Param(":type"))
	if rateType == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	living, _ := receiver.GetInt("living", 0)
	esLiveDataBusiness := es.NewEsLiveDataBusiness()
	_, data := esLiveDataBusiness.LiveCompositeByCategory(startTime, endTime, rateType, living)
	rateData := make([]dy.NameValueFloat64Chart, 0)
	if rateType == 1 {
		list := make([]es2.DyLiveCategoryRateByWatchCnt, 0)
		utils.MapToStruct(data, &list)
		var totalValue int64 = 0
		watchCntMap := map[string]int64{}
		for _, v := range list {
			categories := strings.Split(v.Key.DcmLevelFirst, ",")
			for _, c := range categories {
				c = strings.Trim(c, " ")
				if c == "" || c == "其他" {
					continue
				}
				if _, exist := watchCntMap[c]; exist {
					watchCntMap[c] += v.TotalWatchCnt.Value
				} else {
					watchCntMap[c] = v.TotalWatchCnt.Value
				}
				totalValue += v.TotalWatchCnt.Value
			}
		}
		for k, v := range watchCntMap {
			var rate float64 = 0
			if totalValue > 0 {
				rate = float64(v) / float64(totalValue)
			}
			rateData = append(rateData, dy.NameValueFloat64Chart{
				Name:  k,
				Value: rate,
			})
		}
	} else {
		list := make([]es2.DyLiveCategoryRateByGmv, 0)
		utils.MapToStruct(data, &list)
		var totalValue float64 = 0
		gmvMap := map[string]float64{}
		for _, v := range list {
			categories := strings.Split(v.Key.DcmLevelFirst, ",")
			for _, c := range categories {
				c = strings.Trim(c, " ")
				if c == "" || c == "其他" {
					continue
				}
				if _, exist := gmvMap[c]; exist {
					gmvMap[c] += v.TotalGmv.Value
				} else {
					gmvMap[c] = v.TotalGmv.Value
				}
				totalValue += v.TotalGmv.Value
			}
		}
		for k, v := range gmvMap {
			var rate float64 = 0
			if totalValue > 0 {
				rate = v / totalValue
			}
			rateData = append(rateData, dy.NameValueFloat64Chart{
				Name:  k,
				Value: rate,
			})
		}
	}
	sort.Slice(rateData, func(i, j int) bool {
		return rateData[i].Value > rateData[j].Value
	})
	receiver.SuccReturn(rateData)
	return
}

//带货直播分类统计
func (receiver *LiveCountController) LiveSumByCategory() {
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	category := receiver.GetString("category", "")
	if category == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	living, _ := receiver.GetInt("living", 0)
	esLiveDataBusiness := es.NewEsLiveDataBusiness()
	total, uv, buyRate, data := esLiveDataBusiness.ProductLiveDataByCategory(startTime, endTime, category, living)
	returnData := dy.LiveSumCountByCategory{
		RoomNum:   total,
		WatchCnt:  utils.ToInt64(data.TotalWatchCnt.Value),
		UserCount: utils.ToInt64(data.TotalUserCount.Value),
		Gmv:       data.TotalGmv.Value,
		BuyRate:   buyRate,
		Uv:        uv,
	}
	if living == 0 {
		liveCountBusiness := business.NewLiveCountBusiness()
		monthData, comErr := liveCountBusiness.CountMonthInc(startTime, startTime, category)
		if comErr == nil {
			if monthData.RoomNum > 0 {
				returnData.RoomNumMonthInc = float64(returnData.RoomNum-monthData.RoomNum) / float64(monthData.RoomNum)
			} else {
				returnData.RoomNumMonthInc = 1
			}
			if monthData.WatchCnt > 0 {
				returnData.WatchCntMonthInc = float64(returnData.WatchCnt-monthData.WatchCnt) / float64(monthData.WatchCnt)
			} else {
				returnData.WatchCntMonthInc = 1
			}
			if monthData.UserCount > 0 {
				returnData.UserCountMonthInc = float64(returnData.UserCount-monthData.UserCount) / float64(monthData.UserCount)
			} else {
				returnData.UserCountMonthInc = 1
			}
			if monthData.Gmv > 0 {
				returnData.GmvMonthInc = (returnData.Gmv - monthData.Gmv) / monthData.Gmv
			} else {
				returnData.GmvMonthInc = 1
			}
			if monthData.Uv > 0 {
				returnData.UvMonthInc = (returnData.Uv - monthData.Uv) / monthData.Uv
			} else {
				returnData.UvMonthInc = 1
			}
			if monthData.BuyRate > 0 {
				returnData.BuyRateMonthInc = (returnData.BuyRate - monthData.BuyRate) / monthData.BuyRate
			} else {
				returnData.BuyRateMonthInc = 1
			}
		}
		lastData, comErr := liveCountBusiness.CountLastInc(startTime, startTime, category)
		if comErr == nil {
			if lastData.RoomNum > 0 {
				returnData.RoomNumLastInc = float64(returnData.RoomNum-lastData.RoomNum) / float64(lastData.RoomNum)
			} else {
				returnData.RoomNumLastInc = 1
			}
			if lastData.WatchCnt > 0 {
				returnData.WatchCntLastInc = float64(returnData.WatchCnt-lastData.WatchCnt) / float64(lastData.WatchCnt)
			} else {
				returnData.WatchCntLastInc = 1
			}
			if lastData.UserCount > 0 {
				returnData.UserCountLastInc = float64(returnData.UserCount-lastData.UserCount) / float64(lastData.UserCount)
			} else {
				returnData.UserCountLastInc = 1
			}
			if lastData.Gmv > 0 {
				returnData.GmvLastInc = (returnData.Gmv - lastData.Gmv) / lastData.Gmv
			} else {
				returnData.GmvLastInc = 1
			}
			if lastData.Uv > 0 {
				returnData.UvLastInc = (returnData.Uv - lastData.Uv) / lastData.Uv
			} else {
				returnData.UvLastInc = 1
			}
			if lastData.BuyRate > 0 {
				returnData.BuyRateLastInc = (returnData.BuyRate - lastData.BuyRate) / lastData.BuyRate
			} else {
				returnData.BuyRateLastInc = 1
			}
		}
	}
	receiver.SuccReturn(returnData)
	return
}
