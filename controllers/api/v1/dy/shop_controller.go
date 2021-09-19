package dy

import (
	"dongchamao/business"
	controllers "dongchamao/controllers/api"
)

type ShopController struct {
	controllers.ApiBaseController
}

func (receiver *ShopController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseShowNum)
}

func (receiver *ShopController) ShopBase() {
	shopId := business.IdDecrypt(receiver.Ctx.Input.Param(":shop_id"))
	receiver.SuccReturn(shopId)
}
