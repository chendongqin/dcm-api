package hbase

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	entity2 "dongchamao/models/entity"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"math"
)

//获取直播间商品讲解数据
func GetLiveCurProduct(roomId string) (data entity2.DyLiveCurProduct, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveCurProduct).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity2.DyLiveCurProductMap)
	utils.MapToStruct(infoMap, &data)
	return
}

//直播间全网销量
func GetRoomProductInfo(roomId, productId string) (data entity2.DyRoomProduct, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	rowKey := roomId + "_" + productId
	result, err := query.SetTable(hbaseService.HbaseDyRoomProduct).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity2.DyRoomProductMap)
	utils.MapToStruct(infoMap, &data)
	return
}

//直播间信息
func GetLiveInfo(roomId string) (data entity2.DyLiveInfo, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveInfo).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	liveInfoMap := hbaseService.HbaseFormat(result, entity2.DyLiveInfoMap)
	utils.MapToStruct(liveInfoMap, &data)
	data.Cover = dyimg.Fix(data.Cover)
	data.User.Avatar = dyimg.Fix(data.User.Avatar)
	data.RealSales = math.Floor(data.RealSales)
	data.PredictSales = math.Floor(data.PredictSales)
	data.RoomID = roomId
	return
}

//直播间带货口碑数据
func GetLiveReputation(roomId string) (data entity2.DyLiveReputation, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveReputation).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity2.DyLiveReputationMap)
	utils.MapToStruct(infoMap, &data)
	data.RoomId = roomId
	return
}

//直播间带货数据
func GetLiveSalesData(roomId string) (data entity2.DyAuthorLiveSalesData, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAuthorLiveSalesData).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity2.DyAuthorLiveSalesDataMap)
	utils.MapToStruct(detailMap, &data)
	data.Sales = math.Floor(data.Sales)
	return
}

//直播间商品pmt信息
func GetLivePmt(roomId string) (data entity2.DyLivePmt, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLivePmt).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity2.DyLivePmtMap)
	utils.MapToStruct(detailMap, &data)
	return
}

//获取直播榜单数据
func GetRankTrends(roomId string) (data entity2.DyLiveRankTrends, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveRankTrend).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity2.DyLiveRankTrendsMap)
	utils.MapToStruct(detailMap, &data)
	return
}
