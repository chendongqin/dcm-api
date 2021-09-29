package hbase

import (
	"context"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/entity"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/hbaseService/hbasehelper"
	"math"
	"strings"
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
func GetRoomProductInfo(rowKey string) (data entity.DyRoomProduct, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
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

//直播间全网销量
func GetRoomProductTrend(rowKey string) (data entity.DyRoomProductTrendInfo, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyRoomProductTrend).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyRoomProductTrendMap)
	utils.MapToStruct(infoMap, &data)
	return
}

func GetRoomCurProduct(rowKey string) (data entity.DyRoomCurProduct, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyRoomCurProduct).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyRoomCurProductMap)
	utils.MapToStruct(infoMap, &data)
	return
}

func GetRoomProductInfoRangDate(startRowKey, stopRowKey string) (data map[string]entity.DyRoomProduct, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	results, err := query.
		SetTable(hbaseService.HbaseDyRoomProduct).
		SetStartRow([]byte(startRowKey)).
		SetStopRow([]byte(stopRowKey)).
		Scan(10000)
	if err != nil {
		return
	}
	data = map[string]entity.DyRoomProduct{}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		productId := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity.DyRoomProductMap)
		detail := entity.DyRoomProduct{}
		utils.MapToStruct(dataMap, &detail)
		data[productId] = detail
	}
	return
}

//直播间粉丝
func GetDyLiveRoomUserInfo(roomId string) (data entity.DyLiveRoomUserInfo, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveRoomUserInfo).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyLiveRoomUserInfoMap)
	utils.MapToStruct(infoMap, &data)
	return
}

func GetLiveInfoByIds(roomIds []string) (map[string]entity.DyLiveInfo, error) {
	rowKeys := make([]*hbase.TGet, 0)
	for _, roomId := range roomIds {
		rowKeys = append(rowKeys, &hbase.TGet{Row: []byte(roomId)})
	}
	client := global.HbasePools.Get("default")
	defer client.Close()
	results, err := client.GetMultiple(context.Background(), []byte(hbaseService.HbaseDyLiveInfo), rowKeys)
	if err != nil {
		return nil, err
	}
	roomMap := map[string]entity.DyLiveInfo{}
	for _, v := range results {
		data := entity.DyLiveInfo{}
		liveInfoMap := hbaseService.HbaseFormat(v, entity.DyLiveInfoMap)
		utils.MapToStruct(liveInfoMap, &data)
		data.Cover = dyimg.Fix(data.Cover)
		data.User.Avatar = dyimg.Fix(data.User.Avatar)
		data.RealSales = math.Floor(data.RealSales)
		data.PredictSales = math.Floor(data.PredictSales)
		//todo 套头gmv
		//if data.TotalGmv > data.PredictGmv {
		data.PredictGmv = data.TotalGmv
		//}
		data.RoomID = string(v.Row)
		roomMap[data.RoomID] = data
	}
	return roomMap, nil
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
	data.RealSales = math.Floor(data.RealSales)
	data.PredictSales = math.Floor(data.PredictSales)
	//todo 套头gmv
	//if data.TotalGmv > data.PredictGmv {
	data.PredictGmv = data.TotalGmv
	//}
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
	data.Sales = math.Floor(data.Sales)
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

//获取直播间商品讲解数据
func GetLiveChatMessage(roomId string) (data entity.DyLiveChatMessage, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveChatMessage).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyLiveChatMessageMap)
	utils.MapToStruct(infoMap, &data)
	return
}
