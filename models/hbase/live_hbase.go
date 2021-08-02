package hbase

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/hbase/entity"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
)

//获取直播间商品讲解数据
func GetLiveCurProduct(roomId string) (data entity.DyLiveCurProduct, comErr global.CommonError) {
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
	infoMap := hbaseService.HbaseFormat(result, entity.DyLiveCurProductMap)
	utils.MapToStruct(infoMap, &data)
	return
}

//直播间全网销量
func GetRoomProductInfo(roomId, productId string) (data entity.DyRoomProduct, comErr global.CommonError) {
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
	infoMap := hbaseService.HbaseFormat(result, entity.DyRoomProductMap)
	utils.MapToStruct(infoMap, &data)
	return
}

//直播间信息
func GetLiveInfo(roomId string) (data entity.DyLiveInfo, comErr global.CommonError) {
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
	liveInfoMap := hbaseService.HbaseFormat(result, entity.DyLiveInfoMap)
	utils.MapToStruct(liveInfoMap, &data)
	data.Cover = dyimg.Fix(data.Cover)
	data.User.Avatar = dyimg.Fix(data.User.Avatar)
	data.RoomID = roomId
	return
}

//直播间带货口碑数据
func GetLiveReputation(roomId string) (data entity.DyLiveReputation, comErr global.CommonError) {
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
	infoMap := hbaseService.HbaseFormat(result, entity.DyLiveReputationMap)
	utils.MapToStruct(infoMap, &data)
	data.RoomId = roomId
	return
}

//直播间带货数据
func GetLiveSalesData(roomId string) (data entity.DyAuthorLiveSalesData, comErr global.CommonError) {
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
	detailMap := hbaseService.HbaseFormat(result, entity.DyAuthorLiveSalesDataMap)
	utils.MapToStruct(detailMap, &data)
	return
}

//直播间商品pmt信息
func GetLivePmt(roomId string) (data entity.DyLivePmt, comErr global.CommonError) {
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
	detailMap := hbaseService.HbaseFormat(result, entity.DyLivePmtMap)
	utils.MapToStruct(detailMap, &data)
	return
}

//获取直播榜单数据
func GetRankTrends(roomId string) (data entity.DyLiveRankTrends, comErr global.CommonError) {
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
	detailMap := hbaseService.HbaseFormat(result, entity.DyLiveRankTrendsMap)
	utils.MapToStruct(detailMap, &data)
	return
}
