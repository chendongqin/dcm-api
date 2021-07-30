package business

import (
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
)

type RankBusiness struct {
}

func NewRankBusiness() *RankBusiness {
	return new(RankBusiness)
}

//抖音视频达人热榜
func (receiver *RankBusiness) HbaseStartAuthorVideoRank(rankType, category string) (data []entity.XtHotAwemeAuthorData, comErr global.CommonError) {
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
	return
}

//抖音直播达人热榜
func (receiver *RankBusiness) HbaseStartAuthorLiveRank(rankType string) (data []entity.XtHotLiveAuthorData, comErr global.CommonError) {
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
	return
}
