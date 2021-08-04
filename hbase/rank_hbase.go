package hbase

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	entity2 "dongchamao/models/entity"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
)

//抖音视频达人热榜
func GetStartAuthorVideoRank(rankType, category string) (data []entity2.XtHotAwemeAuthorData, crawlTime int64, comErr global.CommonError) {
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
	detailMap := hbaseService.HbaseFormat(result, entity2.XtHotAwemeAuthorMap)
	info := entity2.XtHotAwemeAuthor{}
	utils.MapToStruct(detailMap, &info)
	data = info.Data
	crawlTime = info.UpdateTime
	return
}

//抖音直播达人热榜
func GetStartAuthorLiveRank(rankType string) (data []entity2.XtHotLiveAuthorData, crawlTime int64, comErr global.CommonError) {
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
	detailMap := hbaseService.HbaseFormat(result, entity2.XtHotLiveAuthorMap)
	info := entity2.XtHotLiveAuthor{}
	utils.MapToStruct(detailMap, &info)
	data = info.Data
	crawlTime = info.UpdateTime
	return
}
