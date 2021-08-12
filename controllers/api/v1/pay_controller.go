package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/models/repost"
	"dongchamao/services/payer"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"time"
)

type PayController struct {
	controllers.ApiBaseController
}

func (receiver *PayController) CreateDyOrder() {
	if !business.UserActionLock("vip_order", 2) {
		receiver.FailReturn(global.NewError(4211))
		return
	}
	InputData := receiver.InputFormat()
	orderType := InputData.GetInt("order_type", 0)
	referrer := InputData.GetString("referrer", "")
	groupPeople := InputData.GetInt("group_people", 0)
	buyDays := InputData.GetInt("days", 0)
	if !utils.InArrayInt(buyDays, []int{30, 180, 365}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if orderType == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	plat := business.VipPlatformDouYin
	userVip := dcm.DcUserVip{}
	dbSession := dcm.GetDbSession()
	exist, err := dbSession.Where("user_id=? AND platform = ?", receiver.UserId, plat).Get(&userVip)
	if !exist || err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	if userVip.ParentId > 0 && orderType != 1 {
		receiver.FailReturn(global.NewMsgError("协同子账号只能购买会员业务～"))
		return
	}
	if userVip.Expiration != userVip.SubExpiration && userVip.SubExpiration.After(time.Now()) && orderType == 3 {
		receiver.FailReturn(global.NewMsgError("请先续费协同子账号才可购买协同子账号～"))
		return
	}
	if userVip.Expiration.Before(time.Now()) {
		userVip.Level = 0
	}
	if orderType == 2 && userVip.Level == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if orderType == 5 && userVip.Expiration.Unix() != userVip.SubExpiration.Unix() {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if orderType > 2 && userVip.Level == 0 {
		receiver.FailReturn(global.NewMsgError("购买协同账号请先开通会员"))
		return
	}
	subExpiration := userVip.SubExpiration
	if subExpiration.Before(time.Now()) {
		subExpiration = time.Now()
	}
	surplusDay := (userVip.Expiration.Unix() - subExpiration.Unix()) / 86400
	if surplusDay == 0 {
		receiver.FailReturn(global.NewMsgError("协同账号到期时间与账户会员时间一致，不需要续费～"))
		return
	}
	payBusiness := business.NewPayBusiness()
	var surplusValue float64 = 0
	vipOrderType := 1
	if userVip.Level > 0 {
		vipOrderType = 2
		surplusValue = payBusiness.CountDySurplusValue(int(surplusDay))
	}
	dyVipValue := business.DyVipPayMoney
	title := fmt.Sprintf("专业版%d天", buyDays)
	var amount float64 = 0
	orderInfo := repost.VipOrderInfo{
		SurplusValue: surplusValue,
	}
	//购买会员
	if orderType == 1 {
		amount = dyVipValue[buyDays]
		orderInfo.BuyDays = buyDays
		orderInfo.Amount = amount
		orderInfo.People = 1
		orderInfo.Title = "会员购买"
	} else if orderType == 2 { //购买协同账号
		title = fmt.Sprintf("购买协同账号%d人", groupPeople)
		amount = utils.FriendlyFloat64(surplusValue * float64(groupPeople))
		orderInfo.BuyDays = int(surplusDay)
		orderInfo.Amount = amount
		orderInfo.People = groupPeople
		orderInfo.Title = "协同账号购买"
		vipOrderType = 3
	} else if orderType == 3 { //协同账号续费
		totalPeople := userVip.SubNum
		title = fmt.Sprintf("协同账号续费%d人", totalPeople)
		amount = utils.FriendlyFloat64(surplusValue * float64(totalPeople))
		orderInfo.BuyDays = int(surplusDay)
		orderInfo.Amount = amount
		orderInfo.People = totalPeople
		orderInfo.Title = "协同账号续费"
		vipOrderType = 4
	} else {
		title = "团队成员续费"
		totalPeople := userVip.SubNum + 1
		amount = utils.FriendlyFloat64(dyVipValue[buyDays] * float64(totalPeople))
		orderInfo.BuyDays = buyDays
		orderInfo.Amount = amount
		orderInfo.People = totalPeople
		orderInfo.Title = "团队成员续费"
		vipOrderType = 5
	}
	uniqueID, _ := utils.Snow.GetSnowflakeId()
	tradeNo := fmt.Sprintf("%s%d", time.Now().Format("060102"), uniqueID)
	orderInfoJson, _ := jsoniter.Marshal(orderInfo)
	vipOrder := dcm.DcVipOrder{
		UserId:         receiver.UserId,
		Username:       receiver.UserInfo.Username,
		TradeNo:        tradeNo,
		OrderType:      vipOrderType,
		Platform:       "douyin",
		Title:          title,
		Amount:         utils.ToString(amount),
		TicketAmount:   "0",
		Level:          business.UserLevelJewel,
		BuyDays:        orderInfo.BuyDays,
		GoodsInfo:      string(orderInfoJson),
		Referrer:       referrer,
		ExpirationTime: time.Now().Add(1800 * time.Second),
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}
	_, err = dcm.Insert(nil, &vipOrder)
	if vipOrder.Id == 0 {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"order_id": vipOrder.Id,
	})
	return
}

//微信支付
func (receiver *PayController) WechatPay() {
	channel := receiver.Ctx.Input.Param(":channel")
	orderId := utils.ToInt(receiver.Ctx.Input.Param(":order_id"))
	if orderId == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if !utils.InArrayString(channel, []string{"app", "native"}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	vipOrder := dcm.DcVipOrder{}
	exist, _ := dcm.Get(orderId, &vipOrder)
	if !exist || vipOrder.UserId != receiver.UserId {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if vipOrder.PayType == "wechat" {
		if (channel == "native" && vipOrder.Channel == 2) || (channel == "app" && vipOrder.Channel == 1) {
			receiver.FailReturn(global.NewMsgError("订单不可支付，请刷新重试～"))
			return
		}
	}
	if vipOrder.PayStatus == 1 {
		receiver.FailReturn(global.NewMsgError("请勿重复付款～"))
		return
	}
	if vipOrder.ExpirationTime.Before(time.Now()) {
		receiver.FailReturn(global.NewMsgError("订单已失效～"))
		return
	}
	amount := utils.ToFloat64(vipOrder.Amount) * float64(100)
	amountInt := utils.ToInt64(amount)
	if global.IsDev() {
		amountInt = 1
	}
	exp := vipOrder.ExpirationTime.Unix() - time.Now().Unix()
	if channel == "native" {
		codeUrl, err := payer.NativePay(amountInt, vipOrder.TradeNo, vipOrder.Title, "/v1/pay/notify/wechat", time.Duration(exp))
		if err != nil {
			receiver.FailReturn(global.NewError(5000))
			return
		}
		_, _ = dcm.UpdateInfo(nil, orderId, map[string]interface{}{"pay_type": "wechat", "channel": 1}, new(dcm.DcVipOrder))
		receiver.SuccReturn(map[string]interface{}{
			"code_url": codeUrl,
		})
		return
	}
	prepayId, err := payer.AppPay(amountInt, vipOrder.TradeNo, vipOrder.Title, "/v1/pay/notify/wechat", time.Duration(exp))
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	prepayIdString := fmt.Sprintf("%v", *prepayId)
	timestamp := time.Now().Unix()
	nonceStr := utils.GetRandomInt(16)
	appId := global.Cfg.String("wechat_pay_app_appId")
	mchId := global.Cfg.String("wechat_pay_mchid")
	signStr := fmt.Sprintf("%s%d%s%s", appId, timestamp, nonceStr, prepayIdString)
	sign, _ := payer.Sha256WithRsa(signStr)
	_, _ = dcm.UpdateInfo(nil, orderId, map[string]interface{}{"pay_type": "wechat", "channel": 2}, new(dcm.DcVipOrder))
	receiver.SuccReturn(map[string]interface{}{
		"appid":          appId,
		"partnerid":      mchId,
		"prepayid":       prepayId,
		"wechat_package": "Sign=WXPay",
		"noncestr":       nonceStr,
		"timestamp":      timestamp,
		"sign":           sign,
	})
	return
}

func (receiver *PayController) AliPay() {
	channel := receiver.Ctx.Input.Param(":channel")
	orderId := utils.ToInt(receiver.Ctx.Input.Param(":order_id"))
	if orderId == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	returnUrl := receiver.GetString("return_url", "")
	returnUrl, _ = url.QueryUnescape(returnUrl)
	if !utils.InArrayString(channel, []string{"app", "page"}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	vipOrder := dcm.DcVipOrder{}
	exist, _ := dcm.Get(orderId, &vipOrder)
	if !exist || vipOrder.UserId != receiver.UserId {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if vipOrder.PayStatus == 1 {
		receiver.FailReturn(global.NewMsgError("请勿重复付款～"))
		return
	}
	if vipOrder.ExpirationTime.Before(time.Now()) {
		receiver.FailReturn(global.NewMsgError("订单已失效～"))
		return
	}
	amount := utils.ToFloat64(vipOrder.Amount)
	if global.IsDev() {
		amount = 0.01
	}
	timeOutExp := "30m"
	if channel == "page" {
		payUrl, err := payer.AliTradePagePay(amount, vipOrder.TradeNo, vipOrder.Title, "/v1/pay/notify/alipay", returnUrl, timeOutExp)
		if err != nil {
			receiver.FailReturn(global.NewError(5000))
			return
		}
		_, _ = dcm.UpdateInfo(nil, orderId, map[string]interface{}{"pay_type": "alipay", "channel": 1}, new(dcm.DcVipOrder))
		receiver.SuccReturn(map[string]interface{}{
			"pay_url": payUrl,
		})
		return
	}
	payParam, err := payer.AliTradeAppPay(amount, vipOrder.TradeNo, vipOrder.Title, "/v1/pay/notify/alipay", returnUrl, timeOutExp)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	_, _ = dcm.UpdateInfo(nil, orderId, map[string]interface{}{"pay_type": "alipay", "channel": 2}, new(dcm.DcVipOrder))
	receiver.SuccReturn(map[string]interface{}{
		"pay_param": payParam,
	})
	return
}
