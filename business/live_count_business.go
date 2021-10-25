package business

import (
	"dongchamao/business/es"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/repost/dy"
	"time"
)

type LiveCountBusiness struct {
}

func NewLiveCountBusiness() *LiveCountBusiness {
	return new(LiveCountBusiness)
}

func (receiver *LiveCountBusiness) getLastMonth(startTime, endTime time.Time) (time.Time, time.Time) {
	return startTime.AddDate(0, -1, 0), endTime.AddDate(0, -1, 0)
}

func (receiver *LiveCountBusiness) getLastDate(startTime, endTime time.Time) (time.Time, time.Time) {
	days := endTime.Day() - startTime.Day()
	return startTime.AddDate(0, 0, days-1), endTime.AddDate(0, 0, days-1)
}

//同比上月数据
func (receiver *LiveCountBusiness) CountMonthInc(startTime, endTime time.Time, category string) (returnData dy.LiveSumCountByCategoryBase, comErr global.CommonError) {
	startTime, endTime = receiver.getLastMonth(startTime, endTime)
	firstTime := GetFirstDay()
	if endTime.Before(firstTime) {
		comErr = global.NewError(4000)
		return
	}
	if startTime.Before(firstTime) {
		startTime = firstTime
	}
	total, uv, buyRate, data := es.NewEsLiveDataBusiness().ProductLiveDataByCategory(startTime, endTime, category, 0)
	returnData = dy.LiveSumCountByCategoryBase{
		RoomNum:   total,
		WatchCnt:  utils.ToInt64(data.TotalWatchCnt.Value),
		UserCount: utils.ToInt64(data.TotalUserCount.Value),
		Gmv:       data.TotalGmv.Value,
		BuyRate:   buyRate,
		Uv:        uv,
	}
	return
}

//同比上月数据
func (receiver *LiveCountBusiness) CountGmvMonthInc(startTime, endTime time.Time, category string) (value float64, comErr global.CommonError) {
	startTime, endTime = receiver.getLastMonth(startTime, endTime)
	firstTime := GetFirstDay()
	if endTime.Before(firstTime) {
		comErr = global.NewError(4000)
		return
	}
	if startTime.Before(firstTime) {
		startTime = firstTime
	}
	data := es.NewEsLiveDataBusiness().RoomProductDataByCategory(startTime, endTime, category, 0)
	value = data.TotalGmv.Value
	return
}

//环比上期数据
func (receiver *LiveCountBusiness) CountLastInc(startTime, endTime time.Time, category string) (returnData dy.LiveSumCountByCategoryBase, comErr global.CommonError) {
	startTime, endTime = receiver.getLastDate(startTime, endTime)
	firstTime := GetFirstDay()
	if endTime.Before(firstTime) {
		comErr = global.NewError(4000)
		return
	}
	if startTime.Before(firstTime) {
		startTime = firstTime
	}
	total, uv, buyRate, data := es.NewEsLiveDataBusiness().ProductLiveDataByCategory(startTime, endTime, category, 0)
	returnData = dy.LiveSumCountByCategoryBase{
		RoomNum:   total,
		WatchCnt:  utils.ToInt64(data.TotalWatchCnt.Value),
		UserCount: utils.ToInt64(data.TotalUserCount.Value),
		Gmv:       data.TotalGmv.Value,
		BuyRate:   buyRate,
		Uv:        uv,
	}
	return
}
