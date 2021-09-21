package hbase

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/entity"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"strconv"
)

//抖音视频达人热榜
func GetStartAuthorVideoRank(rankType, category string) (data []entity.XtHotAwemeAuthorData, crawlTime int64, comErr global.CommonError) {
	rowKey := utils.Md5_encode(rankType + "_" + category)
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseXtHotAwemeAuthorRank).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.XtHotAwemeAuthorMap)
	info := entity.XtHotAwemeAuthor{}
	utils.MapToStruct(detailMap, &info)
	data = info.Data
	for k, v := range data {
		if v.UniqueId == "" {
			data[k].UniqueId = v.ShortId
		}
		data[k].FieldsMap = map[string]interface{}{}
		for _, d := range v.Fields {
			data[k].FieldsMap[d.Label] = d.Value
		}
	}
	crawlTime = info.UpdateTime
	return
}

//抖音直播达人热榜
func GetStartAuthorLiveRank(rankType string) (data []entity.XtHotLiveAuthorData, crawlTime int64, comErr global.CommonError) {
	rowKey := utils.Md5_encode(rankType)
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseXtHotLiveAuthorRank).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.XtHotLiveAuthorMap)
	info := entity.XtHotLiveAuthor{}
	utils.MapToStruct(detailMap, &info)
	data = info.Data
	for k, v := range data {
		if v.UniqueId == "" {
			data[k].UniqueId = v.ShortId
		}
		data[k].FieldsMap = map[string]interface{}{}
		for _, d := range v.Fields {
			data[k].FieldsMap[d.Label] = d.Value
		}
	}
	crawlTime = info.UpdateTime
	return
}

func GetDyLiveHourRank(hour string) (data entity.DyLiveHourRanks, comErr global.CommonError) {
	rowKey := utils.Md5_encode(hour + "_小时榜")
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveHourRank).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLiveHourRankMap)
	utils.MapToStruct(detailMap, &data)
	return
}

func GetDyLiveTopRank(hour string) (data entity.DyLiveTopRanks, comErr global.CommonError) {
	rowKey := utils.Md5_encode(hour + "_实时热榜")
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveTopRank).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLiveTopMap)
	utils.MapToStruct(detailMap, &data)
	return
}

func GetDyLiveHourSellRank(hour string) (data entity.DyLiveHourSellRanks, comErr global.CommonError) {
	rowKey := utils.Md5_encode(hour + "_带货榜")
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveHourRankSell).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLiveHourSellRanksMap)
	utils.MapToStruct(detailMap, &data)
	return
}

func GetDyLiveHourPopularityRank(hour string) (data entity.DyLiveHourPopularityRanks, comErr global.CommonError) {
	rowKey := utils.Md5_encode(hour + "_人气榜")
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveHourRankPopularity).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLiveHourPopularityRanksMap)
	utils.MapToStruct(detailMap, &data)
	return
}

func GetLiveShareWeekRank(rowKey string) (data entity.DyLiveShareTops, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveShareWeekRank).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLiveShareTopMap)
	utils.MapToStruct(detailMap, &data)
	return
}

func GetAwemeShareRank(rowKey string) (data entity.DyAwemeShareTops, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAwemeShareRank).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyAwemeShareTopMap)
	utils.MapToStruct(detailMap, &data)
	return
}

//商品排行榜
func GetProductRank(day, fCate, sortStr string, hPage int) (data []entity.ShortVideoProduct, comErr global.CommonError) {
	key := day + "_" + fCate + "_" + sortStr
	rowKey := utils.Md5_encode(key)
	rowKey = rowKey + strconv.Itoa(hPage)
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseShortVideoProductRank).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity.ShortVideoCommodityTopNMap)
	hData := entity.ShortVideoCommodityTopN{}
	utils.MapToStruct(dataMap, &hData)
	data = hData.Ranks
	return
}
