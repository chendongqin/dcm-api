package business

import (
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
)

type LiveBusiness struct {
}

func NewLiveBusiness() *LiveBusiness {
	return new(LiveBusiness)
}

//直播间信息
func (l *LiveBusiness) HbaseGetLiveInfo(roomId string) (data *entity.DyLiveInfo, comErr global.CommonError) {
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
	liveInfo := &entity.DyLiveInfo{}
	liveInfoMap := hbaseService.HbaseFormat(result, entity.DyLiveInfoMap)
	utils.MapToStruct(liveInfoMap, liveInfo)
	liveInfo.Cover = dyimg.Fix(liveInfo.Cover)
	liveInfo.User.Avatar = dyimg.Fix(liveInfo.User.Avatar)
	liveInfo.RoomID = roomId
	data = liveInfo
	return
}

//直播间信息
func (l *LiveBusiness) HbaseGetLivePmt(roomId string) (data *entity.DyLivePmt, comErr global.CommonError) {
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
	detail := &entity.DyLivePmt{}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLivePmtMap)
	utils.MapToStruct(detailMap, detail)
	data = detail
	return
}
