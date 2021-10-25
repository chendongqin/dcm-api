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
	//_, data := esLiveDataBusiness.LiveCompositeByCategory(startTime, endTime, rateType, living)
	rateData := make([]dy.NameValueFloat64Chart, 0)
	cateList := business.NewProductBusiness().GetCacheProductCate(true)
	total := 0
	var totalSum float64 = 0
	for _, v := range cateList {
		if v.Name == "其他" {
			continue
		}
		temTotal, temRes := esLiveDataBusiness.LiveCompositeByCategoryOne(startTime, endTime, rateType, living, v.Name)
		total += temTotal
		var tempSum float64 = 0
		if rateType == 1 {
			tempData := es2.LiveCategoryWatchCnt{}
			utils.MapToStruct(temRes, &tempData)
			tempSum = float64(tempData.TotalWatchCnt.Value)
		} else {
			tempData := es2.LiveCategoryGmv{}
			utils.MapToStruct(temRes, &tempData)
			tempSum = float64(tempData.TotalGmv.Value)
		}
		totalSum += tempSum
		rateData = append(rateData, dy.NameValueFloat64Chart{
			Name:  v.Name,
			Value: tempSum,
		})
	}
	newRateData := make([]dy.NameValueFloat64Chart, 0)
	for _, v := range rateData {
		if v.Value == 0 {
			continue
		}
		var value float64 = 0
		if totalSum > 0 {
			value = utils.RateMin(v.Value / totalSum)
		}
		newRateData = append(newRateData, dy.NameValueFloat64Chart{
			Name:  v.Name,
			Value: value,
		})
	}
	//if rateType == 1 {
	//	list := make([]es2.DyLiveCategoryRateByWatchCnt, 0)
	//	utils.MapToStruct(data, &list)
	//	var totalValue int64 = 0
	//	watchCntMap := map[string]int64{}
	//	for _, v := range list {
	//		categories := strings.Split(v.Key.DcmLevelFirst, ",")
	//		for _, c := range categories {
	//			c = strings.Trim(c, " ")
	//			if c == "" || c == "其他" {
	//				continue
	//			}
	//			if _, exist := watchCntMap[c]; exist {
	//				watchCntMap[c] += v.TotalWatchCnt.Value
	//			} else {
	//				watchCntMap[c] = v.TotalWatchCnt.Value
	//			}
	//			totalValue += v.TotalWatchCnt.Value
	//		}
	//	}
	//	for k, v := range watchCntMap {
	//		var rate float64 = 0
	//		if totalValue > 0 {
	//			rate = float64(v) / float64(totalValue)
	//		}
	//		rateData = append(rateData, dy.NameValueFloat64Chart{
	//			Name:  k,
	//			Value: rate,
	//		})
	//	}
	//} else {
	//	list := make([]es2.DyLiveCategoryRateByGmv, 0)
	//	utils.MapToStruct(data, &list)
	//	var totalValue float64 = 0
	//	gmvMap := map[string]float64{}
	//	for _, v := range list {
	//		categories := strings.Split(v.Key.DcmLevelFirst, ",")
	//		for _, c := range categories {
	//			c = strings.Trim(c, " ")
	//			if c == "" || c == "其他" {
	//				continue
	//			}
	//			if _, exist := gmvMap[c]; exist {
	//				gmvMap[c] += v.TotalGmv.Value
	//			} else {
	//				gmvMap[c] = v.TotalGmv.Value
	//			}
	//			totalValue += v.TotalGmv.Value
	//		}
	//	}
	//	for k, v := range gmvMap {
	//		var rate float64 = 0
	//		if totalValue > 0 {
	//			rate = v / totalValue
	//		}
	//		rateData = append(rateData, dy.NameValueFloat64Chart{
	//			Name:  k,
	//			Value: rate,
	//		})
	//	}
	//}
	sort.Slice(newRateData, func(i, j int) bool {
		return newRateData[i].Value > newRateData[j].Value
	})
	receiver.SuccReturn(newRateData)
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
	gmvData := esLiveDataBusiness.RoomProductDataByCategory(startTime, endTime, category, living)
	returnData := dy.LiveSumCountByCategory{
		RoomNum:   total,
		WatchCnt:  utils.ToInt64(data.TotalWatchCnt.Value),
		UserCount: utils.ToInt64(data.TotalUserCount.Value),
		Gmv:       gmvData.TotalGmv.Value,
		BuyRate:   utils.RateMin(buyRate),
		Uv:        uv,
	}
	if living == 0 {
		liveCountBusiness := business.NewLiveCountBusiness()
		monthData, comErr := liveCountBusiness.CountMonthInc(startTime, startTime, category)
		gmvMonthData, comErr := liveCountBusiness.CountGmvMonthInc(startTime, startTime, category)
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
			if gmvMonthData > 0 {
				returnData.GmvMonthInc = (returnData.Gmv - gmvMonthData) / gmvMonthData
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

//带货直播分类分级统计
func (receiver *LiveCountController) LiveSumByCategoryLevel() {
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
	_, dataList := esLiveDataBusiness.ProductLiveDataCategoryLevel(startTime, endTime, category, living)
	CustomerUnitPriceList := esLiveDataBusiness.ProductLiveDataCategoryCustomerUnitPriceLevel(startTime, endTime, category, living)
	levelsArr := []string{"E", "D", "C", "B", "A", "S"}
	levelMap := map[string]dy.LiveSumDataCategoryLevel{}
	var allGmv float64 = 0
	var allWatch int64 = 0
	for _, v := range dataList {
		var avgWatch int64 = 0
		var avgGmv float64 = 0
		if v.TotalGmv.Value > 0 {
			avgGmv = v.TotalGmv.Value / float64(v.DocCount)
		}
		if v.TotalWatchCnt.Value > 0 {
			avgWatch = v.TotalWatchCnt.Value / int64(v.DocCount)
		}
		allGmv += v.TotalGmv.Value
		allWatch += v.TotalWatchCnt.Value
		item := dy.LiveSumDataCategoryLevel{
			Level:      v.Key,
			RoomCount:  v.DocCount,
			TotalWatch: utils.ToInt64(v.TotalWatchCnt.Value),
			AvgWatch:   avgWatch,
			TotalGmv:   v.TotalGmv.Value,
			AvgGmv:     avgGmv,
		}
		item.CustomerUnitPrice.Min = v.StatsCustomerUnitPrice.Min
		item.CustomerUnitPrice.Max = v.StatsCustomerUnitPrice.Max
		levelMap[v.Key] = item
	}
	for _, v := range CustomerUnitPriceList {
		if item, exist := levelMap[v.Key]; exist {
			for _, c := range v.CustomerUnitPrice.Values {
				if c.Key == 50 {
					item.CustomerUnitPrice.Median = c.Value
					levelMap[v.Key] = item
				}
			}
		}
	}
	list := make([]dy.LiveSumDataCategoryLevel, 0)
	for _, key := range levelsArr {
		item := dy.LiveSumDataCategoryLevel{}
		if v, exist := levelMap[key]; exist {
			item = v
		}
		if allGmv > 0 {
			item.GmvPer = item.TotalGmv / allGmv
		}
		if allWatch > 0 {
			item.WatchPer = float64(item.TotalWatch) / float64(allWatch)
		}
		list = append(list, item)
	}
	receiver.SuccReturn(list)
	return
}

//带货直播分类分级统计
func (receiver *LiveCountController) LiveSumByCategoryLevelTwo() {
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
	keyword := receiver.GetString("keyword", "")
	esLiveDataBusiness := es.NewEsLiveDataBusiness()
	_, dataList := esLiveDataBusiness.ProductLiveDataCategoryLevelTwoShow(startTime, endTime, category, living, keyword)
	sort.Slice(dataList, func(i, j int) bool {
		if strings.Index(dataList[j].Key, "S") == 0 && strings.Index(dataList[i].Key, "S") < 0 {
			return true
		}
		if strings.Index(dataList[i].Key, "S") == 0 && strings.Index(dataList[j].Key, "S") < 0 {
			return false
		}
		return dataList[i].Key > dataList[j].Key
	})
	levelMap := map[string]map[int]map[int]dy.LiveSumDataCategoryLevelTwo{}
	for _, v := range dataList {
		keyNameArr := strings.Split(v.Key, "")
		levelKey := keyNameArr[0]
		intKey := utils.ToInt(strings.Replace(v.Key, levelKey, "", 1))
		if _, ok := levelMap[levelKey]; !ok {
			levelMap[levelKey] = map[int]map[int]dy.LiveSumDataCategoryLevelTwo{}
			if _, ok1 := levelMap[levelKey][intKey]; !ok1 {
				levelMap[levelKey][intKey] = map[int]dy.LiveSumDataCategoryLevelTwo{}
			}
		}
		item := map[int]dy.LiveSumDataCategoryLevelTwo{}
		buckets := v.LiveTwo.Buckets
		sort.Slice(buckets, func(i, j int) bool {
			return buckets[i].Key < buckets[j].Key
		})
		for _, v1 := range buckets {
			item[v1.Key] = dy.LiveSumDataCategoryLevelTwo{
				FlowLevel:     v.Key,
				StayLevel:     v1.Key,
				RoomCount:     v1.DocCount,
				TotalWatchCnt: v1.TotalWatchCnt.Value,
				TotalGmv:      v1.TotalGmv.Value,
			}
		}
		levelMap[levelKey][intKey] = item
	}
	levelsArr := []string{"E", "D", "C", "B", "A", "S"}
	list := make([][][]dy.LiveSumDataCategoryLevelTwo, 0)
	for _, level := range levelsArr {
		item := [][]dy.LiveSumDataCategoryLevelTwo{}
		if v, ok := levelMap[level]; ok {
			kBegin := 1
			if level == "E" {
				kBegin = 0
			}
			for i := kBegin; i <= 9; i++ {
				itemMap := map[int]dy.LiveSumDataCategoryLevelTwo{}
				item1 := []dy.LiveSumDataCategoryLevelTwo{}
				if d, exist := v[i]; exist {
					itemMap = d
				}
				for j := 1; j <= 10; j++ {
					tmp := dy.LiveSumDataCategoryLevelTwo{
						FlowLevel: level + utils.ToString(i),
						StayLevel: j,
					}
					if v2, exist1 := itemMap[j]; exist1 {
						tmp = v2
					}
					item1 = append(item1, tmp)
				}
				item = append(item, item1)
			}
		}
		list = append(list, item)
	}
	receiver.SuccReturn(list)
	return
}

//明细列表
func (receiver *LiveCountController) LiveSumByCategoryLevelList() {
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	category := receiver.GetString("category", "")
	keyword := receiver.GetString("keyword", "")
	if category == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	living, _ := receiver.GetInt("living", 0)
	stayLevel := utils.ToInt(receiver.Ctx.Input.Param(":stay_level"))
	level := receiver.Ctx.Input.Param(":level")
	if level == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 5, 30)
	esLiveDataBusiness := es.NewEsLiveDataBusiness()
	total, list, comErr := esLiveDataBusiness.ProductLiveDataCategoryLevelList(startTime, endTime, keyword, category, level, stayLevel, living, page, pageSize)
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
		"list":  list,
		"total": total,
	})
	return
}

//明细统计
func (receiver *LiveCountController) LiveSumByCategoryLevelCount() {
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	category := receiver.GetString("category", "")
	keyword := receiver.GetString("keyword", "")
	if category == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	living, _ := receiver.GetInt("living", 0)
	stayLevel := utils.ToInt(receiver.Ctx.Input.Param(":stay_level"))
	level := receiver.Ctx.Input.Param(":level")
	if level == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	esLiveDataBusiness := es.NewEsLiveDataBusiness()
	total, data, comErr := esLiveDataBusiness.ProductLiveDataCategoryLevelCount(startTime, endTime, keyword, category, level, stayLevel, living)
	var avgWatch int64 = 0
	var avgGmv float64 = 0
	if total > 0 {
		avgWatch = data.TotalWatchCnt.Value / int64(total)
		avgGmv = data.TotalGmv.Value / float64(total)
	}
	receiver.SuccReturn(map[string]interface{}{
		"room_count":          total,
		"total_watch":         data.TotalWatchCnt.Value,
		"avg_watch":           avgWatch,
		"total_gmv":           data.TotalGmv.Value,
		"avg_gmv":             avgGmv,
		"min_user_uint_price": data.StatsCustomerUnitPrice.Min,
		"max_user_uint_price": data.StatsCustomerUnitPrice.Max,
	})
	return
}
