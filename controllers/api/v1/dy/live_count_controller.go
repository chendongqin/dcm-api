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
	esLiveDataBusiness := es.NewEsLiveDataBusiness()
	allTotal, allData := esLiveDataBusiness.SumLiveData(startTime, endTime, 0)
	productTotal, productData := esLiveDataBusiness.SumLiveData(startTime, endTime, 1)
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

//直播总览
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
	sortStr := receiver.GetString("sort", "")
	esLiveDataBusiness := es.NewEsLiveDataBusiness()
	list, comErr := esLiveDataBusiness.LiveRankByCategory(startTime, endTime, category, sortStr)
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

//直播总览
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
	esLiveDataBusiness := es.NewEsLiveDataBusiness()
	_, data := esLiveDataBusiness.LiveCompositeByCategory(startTime, endTime, rateType)
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
