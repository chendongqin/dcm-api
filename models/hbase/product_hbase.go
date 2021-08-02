package hbase

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/hbase/entity"
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
