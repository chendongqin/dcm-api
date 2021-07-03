package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/models/business"
)

type LiveController struct {
	controllers.ApiBaseController
}

func (receiver *LiveController) LiveInfoData() {
	roomId := receiver.Ctx.Input.Param(":room_id")
	if roomId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	liveBusiness := business.NewLiveBusiness()
	liveInfo, comErr := liveBusiness.HbaseGetLiveInfo(roomId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"live_info": liveInfo,
	})
	return
}
