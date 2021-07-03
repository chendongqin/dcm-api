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
	liveInfo := &entity.DyLiveInfo{}
	liveInfoMap := hbaseService.HbaseFormat(result, entity.DyLiveInfoMap)
	utils.MapToStruct(liveInfoMap, liveInfo)
	liveInfo.Cover = dyimg.Fix(liveInfo.Cover)
	liveInfo.User.Avatar = dyimg.Fix(liveInfo.User.Avatar)
	liveInfo.RoomID = roomId
	data = liveInfo
	return
}
