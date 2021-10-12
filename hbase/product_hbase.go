package hbase

import (
	"context"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/entity"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/hbaseService/hbasehelper"
	"strings"
	"time"
)

func GetProductByIds(productIds []string) (map[string]entity.DyProduct, error) {
	rowKeys := make([]*hbase.TGet, 0)
	for _, id := range productIds {
		rowKeys = append(rowKeys, &hbase.TGet{Row: []byte(id)})
	}
	client := global.HbasePools.Get("default")
	defer client.Close()
	results, err := client.GetMultiple(context.Background(), []byte(hbaseService.HbaseDyProduct), rowKeys)
	if err != nil {
		return nil, err
	}
	infoMap := map[string]entity.DyProduct{}
	for _, v := range results {
		data := entity.DyProduct{}
		detailMap := hbaseService.HbaseFormat(v, entity.DyProductMap)
		utils.MapToStruct(detailMap, &data)
		data.ProductID = string(v.Row)
		if data.TbCouponInfo == "null" {
			data.TbCouponInfo = ""
		}
		infoMap[data.ProductID] = data
		if data.ManmadeCategory.FirstCname != "" {
			data.Label = data.ManmadeCategory.FirstCname
		} else if data.AiCategory.FirstCname != "" {
			data.Label = data.AiCategory.FirstCname
		}
		//佣金比例处理
		if data.CosRatio == 0 {
			data.CosRatio = data.SecCosRatio
		}
	}
	return infoMap, nil
}

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
	if data.ManmadeCategory.FirstCname != "" {
		data.Label = data.ManmadeCategory.FirstCname
	} else if data.AiCategory.FirstCname != "" {
		data.Label = data.AiCategory.FirstCname
	}
	//null数据初始化
	if len(data.ContextNum) == 0 {
		data.ContextNum = []entity.ContextNum{}
	}
	if len(data.DiggInfo) == 0 {
		data.DiggInfo = []entity.DiggInfo{}
	}
	//佣金比例处理
	if data.CosRatio == 0 {
		data.CosRatio = data.SecCosRatio
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
	if data.ManmadeCategory.FirstCname != "" {
		data.DcmLevelFirst = data.ManmadeCategory.FirstCname
	}
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
		rowKey := string(v.GetRow())
		rowArr := strings.Split(rowKey, "_")
		dataMap := hbaseService.HbaseFormat(v, entity.DyProductAuthorAnalysisMap)
		hData := entity.DyProductAuthorAnalysis{}
		utils.MapToStruct(dataMap, &hData)
		if len(rowArr) == 3 {
			hData.Date = rowArr[1]
		}
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
func GetDyProductAwemeSalesTrend(productId, date string) (data entity.DyProductAwemeSalesTrend, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := productId + "_" + date
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

//商品gpm数据
func GetDyProductGpmDate(productId, date string) (data entity.AdsDyProductGpmDi, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := productId + "_" + date
	result, err := query.
		SetTable(hbaseService.HbaseAdsDyProductGpmDi).
		GetByRowKey([]byte(startRow))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity.AdsDyProductGpmDiMap)
	utils.MapToStruct(dataMap, &data)
	return
}

//商品gpm数据
func GetDyProductGpmRangeDate(productId string, startTime, endTime time.Time) (data map[string]entity.AdsDyProductGpmDi, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := productId + "_" + startTime.Format("20060102")
	endRow := productId + "_" + endTime.AddDate(0, 0, 1).Format("20060102")
	results, err := query.
		SetTable(hbaseService.HbaseAdsDyProductGpmDi).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	data = map[string]entity.AdsDyProductGpmDi{}
	for _, v := range results {
		dataMap := hbaseService.HbaseFormat(v, entity.AdsDyProductGpmDiMap)
		hData := entity.AdsDyProductGpmDi{}
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
