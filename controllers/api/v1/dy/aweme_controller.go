package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/models/business"
)

type AwemeController struct {
	controllers.ApiBaseController
}

func (receiver *AwemeController) AwemeBaseData() {
	awemeId := receiver.Ctx.Input.Param(":aweme_id")
	if awemeId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	awemeBusiness := business.NewAwemeBusiness()
	awemeBase, comErr := awemeBusiness.HbaseGetAweme(awemeId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"aweme_base": awemeBase.Data,
	})
	return
}
