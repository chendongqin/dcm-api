package hbase

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/entity"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"strings"
	"time"
)

//商品详情
func GetProductInfo(productId string) (data entity.DyProduct, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyProduct).GetByRowKey([]byte(productId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyProductMap)
	utils.MapToStruct(detailMap, &data)
	data.ProductID = productId
	if data.TbCouponInfo == "null" {
		data.TbCouponInfo = ""
	}
	return
}

//获取商品品牌数据
func GetDyProductBrand(productId string) (data entity.DyProductBrand, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyProductBrand).GetByRowKey([]byte(productId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyProductBrandMap)
	utils.MapToStruct(detailMap, &data)
	return
}

//商品月销量数据
func GetPromotionMonth(productId string) (data entity.DyLivePromotionMonth, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLivePromotionMonth).GetByRowKey([]byte(productId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLivePromotionMonthMap)
	utils.MapToStruct(detailMap, &data)
	return
}

//获取商品日销量
func GetProductDailyRangDate(productId string, startTime, endTime time.Time) (data map[string]entity.DyProductDaily, comErr global.CommonError) {
	data = map[string]entity.DyProductDaily{}
	query := hbasehelper.NewQuery()
	startRow := productId + "_" + startTime.Format("20060102")
	endRow := productId + "_" + endTime.AddDate(0, 0, 1).Format("20060102")
	results, err := query.
		SetTable(hbaseService.HbaseDyProductDaily).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity.DyProductDailyMap)
		hData := entity.DyProductDaily{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData
	}
	return
}

//直播间销量趋势
func GetProductLiveSalesRangDate(productId string, startTime, endTime time.Time) (data map[string]entity.DyProductLiveSalesTrend, comErr global.CommonError) {
	data = map[string]entity.DyProductLiveSalesTrend{}
	query := hbasehelper.NewQuery()
	startRow := productId + "_" + startTime.Format("20060102")
	endRow := productId + "_" + endTime.AddDate(0, 0, 1).Format("20060102")
	results, err := query.
		SetTable(hbaseService.HbaseDyProductLiveSalesTrend).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity.DyProductLiveSalesTrendMap)
		hData := entity.DyProductLiveSalesTrend{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData
	}
	return
}

func GetProductAuthorAnalysis(rowKey string) (data entity.DyProductAuthorAnalysis, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyProductAuthorAnalysis).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyProductAuthorAnalysisMap)
	utils.MapToStruct(detailMap, &data)
	return
}

func GetProductAuthorAwemes(rowKey string) (data entity.DyProductAuthorAnalysis, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyProductAuthorAnalysis).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyProductAuthorAnalysisMap)
	utils.MapToStruct(detailMap, &data)
	return
}

func GetProductAuthorAnalysisRange(startRowKey, stopRowKey string) (data []entity.DyProductAuthorAnalysis, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	results, err := query.
		SetTable(hbaseService.HbaseDyProductAuthorAnalysis).
		SetStartRow([]byte(startRowKey)).
		SetStopRow([]byte(stopRowKey)).
		Scan(50000)
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	for _, v := range results {
		dataMap := hbaseService.HbaseFormat(v, entity.DyProductAuthorAnalysisMap)
		hData := entity.DyProductAuthorAnalysis{}
		utils.MapToStruct(dataMap, &hData)
		data = append(data, hData)
	}
	return
}

func GetProductAwemeAuthorAnalysis(rowKey string) (data entity.DyProductAwemeAuthorAnalysis, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyProductAwemeAuthorAnalysis).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyProductAwemeAuthorAnalysisMap)
	utils.MapToStruct(detailMap, &data)
	return
}

func GetProductAwemeAuthorAnalysisRange(startRowKey, stopRowKey string) (data []entity.DyProductAwemeAuthorAnalysis, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	results, err := query.
		SetTable(hbaseService.HbaseDyProductAwemeAuthorAnalysis).
		SetStartRow([]byte(startRowKey)).
		SetStopRow([]byte(stopRowKey)).
		Scan(50000)
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	for _, v := range results {
		dataMap := hbaseService.HbaseFormat(v, entity.DyProductAwemeAuthorAnalysisMap)
		hData := entity.DyProductAwemeAuthorAnalysis{}
		utils.MapToStruct(dataMap, &hData)
		data = append(data, hData)
	}
	return
}

//商品视频某时间段分销数据
func GetDyProductAwemeSalesTrendRangeDate(productId string, startTime, endTime time.Time) (data map[string]entity.DyProductAwemeSalesTrend, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := productId + "_" + startTime.Format("20060102")
	endRow := productId + "_" + endTime.AddDate(0, 0, 1).Format("20060102")
	results, err := query.
		SetTable(hbaseService.HbaseDyProductAwemeSalesTrend).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	data = map[string]entity.DyProductAwemeSalesTrend{}
	for _, v := range results {
		dataMap := hbaseService.HbaseFormat(v, entity.DyProductAwemeSalesTrendMap)
		hData := entity.DyProductAwemeSalesTrend{}
		utils.MapToStruct(dataMap, &hData)
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		data[date] = hData
	}
	return
}

//商品视频某天分销数据
func GetDyProductAwemeSalesTrend(product, date string) (data entity.DyProductAwemeSalesTrend, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := product + "_" + date
	result, err := query.
		SetTable(hbaseService.HbaseDyProductAwemeSalesTrend).
		GetByRowKey([]byte(startRow))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity.DyProductAwemeSalesTrendMap)
	utils.MapToStruct(dataMap, &data)
	return
}

//视频商品数据
func GetDyProductAwemeDailyDistributeRange(awemeId, beginDate, endDate string) (data []entity.DyProductAwemeDailyDistribute, comErr global.CommonError) {
	startRowKey := awemeId + "_" + beginDate + "_"
	stopRowKey := awemeId + "_" + endDate + "_9999999999999999"
	query := hbasehelper.NewQuery()
	results, err := query.
		SetTable(hbaseService.HbaseDyProductAwemeDailyDistribute).
		SetStartRow([]byte(startRowKey)).
		SetStopRow([]byte(stopRowKey)).
		Scan(1000)
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	for _, v := range results {
		dataMap := hbaseService.HbaseFormat(v, entity.DyProductAwemeDailyDistributeMap)
		hData := entity.DyProductAwemeDailyDistribute{}
		utils.MapToStruct(dataMap, &hData)
		data = append(data, hData)
	}
	return
}
