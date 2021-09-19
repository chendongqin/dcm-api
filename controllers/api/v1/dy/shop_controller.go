package dy

import (
	"dongchamao/business"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	"time"
)

type ShopController struct {
	controllers.ApiBaseController
}

func (receiver *ShopController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseShowNum)
}

//小店基本数据
func (receiver *ShopController) ShopBase() {
	var returnRes entity.DyShopBaseBasic
	var comErr global.CommonError
	shopId := business.IdDecrypt(receiver.Ctx.Input.Param(":shop_id"))
	if shopId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	returnRes.BaseData, comErr = hbase.GetShop(shopId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	shopDetailData, comErr := hbase.GetShopDetailByDate(shopId, time.Now().AddDate(0, 0, 0).Format("20060102"))
	if comErr != nil { //今天取不到，取昨日数据
		shopDetailData, _ = hbase.GetShopDetailByDate(shopId, time.Now().AddDate(0, 0, -1).Format("20060102"))
	}

	returnRes.DetailData.Sales = shopDetailData.Sales
	returnRes.DetailData.Gmv = shopDetailData.Gmv
	returnRes.DetailData.D30LiveCnt = shopDetailData.D30LiveCnt
	returnRes.DetailData.D30AuthorCnt = shopDetailData.D30AuthorCnt
	returnRes.DetailData.D30AwemeCnt = shopDetailData.D30AwemeCnt
	returnRes.DetailData.D30Sales = shopDetailData.D30Sales
	returnRes.DetailData.D30Gmv = shopDetailData.D30Gmv
	returnRes.DetailData.D30Pct = shopDetailData.D30Pct

	receiver.SuccReturn(returnRes)
	return
}
