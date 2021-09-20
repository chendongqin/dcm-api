package hbase

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/entity"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"strings"
	"time"
)

//小店数据
func GetShop(shopId string) (data entity.DyShop, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyShop).GetByRowKey([]byte(shopId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyShopMap)
	utils.MapToStruct(detailMap, &data)
	data.ShopId = string(result.Row)
	data.Logo = dyimg.Fix(data.Logo)
	return
}

//获取小店某天数据
func GetShopDetailByDate(shopId, date string) (data entity.DyShopDetail, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	rowKey := shopId + "_" + date
	result, err := query.SetTable(hbaseService.HbaseDyShopDetail).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyShopDetailMap)
	utils.MapToStruct(infoMap, &data)
	return
}

//获取小店区间内数据
func GetShopDetailRangDate(shopId string, startTime, endTime time.Time) (data map[string]entity.DyShopDetail, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := shopId + "_" + startTime.Format("20060102")
	endRow := shopId + "_" + endTime.Format("20060102")
	results, err := query.
		SetTable(hbaseService.HbaseDyShopDetail).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		return
	}
	data = map[string]entity.DyShopDetail{}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity.DyShopDetailMap)
		hData := entity.DyShopDetail{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData
	}
	return
}

//获取小店商品分析
func GetShopProductAnalysisByDate(shopId, productId, date string) (data entity.DyShopProductAnalysis, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	rowKey := shopId + "_" + date + "_" + productId
	result, err := query.SetTable(hbaseService.HbaseDyShopProductAnalysis).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyShopProductAnalysisMap)
	utils.MapToStruct(infoMap, &data)
	data.Date = date
	return
}

//获取小店商品分析
func GetShopProductAnalysisRangDate(shopId, starProductId, endProductId string, startTime, endTime time.Time) (data []entity.DyShopProductAnalysis, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := shopId + "_" + startTime.Format("20060102") + "_" + starProductId
	endRow := shopId + "_" + endTime.Format("20060102") + "_" + endProductId
	results, err := query.
		SetTable(hbaseService.HbaseDyShopProductAnalysis).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		return
	}
	data = make([]entity.DyShopProductAnalysis, 0)
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 3 {
			continue
		}
		dataMap := hbaseService.HbaseFormat(v, entity.DyShopProductAnalysisMap)
		hData := entity.DyShopProductAnalysis{}
		utils.MapToStruct(dataMap, &hData)
		hData.Date = utils.ToString(rowKeyArr[1])
		data = append(data, hData)
	}
	return
}
