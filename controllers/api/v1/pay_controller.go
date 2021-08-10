package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/services/payer"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type PayController struct {
	controllers.ApiBaseController
}

func (receiver *PayController) CreateDyOrder() {
	InputData := receiver.InputFormat()
	orderType := InputData.GetInt("order_type", 0)
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
	if userVip.Expiration.Before(time.Now()) {
		userVip.Level = 0
	}
	if orderType > 1 && userVip.Level == 0 {
		receiver.FailReturn(global.NewMsgError("购买协同账号请先开通会员"))
		return
	}
	surplusDay := (userVip.Expiration.Unix() - time.Now().Unix()) / 86400
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
	orderInfo := map[string]interface{}{
		"surplus_value": surplusValue,
	}
	//购买会员
	if orderType == 1 {
		amount = dyVipValue[buyDays]
		orderInfo["buy_days"] = buyDays
		orderInfo["amount"] = amount
		orderInfo["people"] = 1
		orderInfo["title"] = "会员购买"
	} else if orderType == 2 { //购买协同账号
		title = fmt.Sprintf("购买协同账号%d人", groupPeople)
		amount = utils.FriendlyFloat64(surplusValue * float64(groupPeople))
		orderInfo["buy_days"] = surplusDay
		orderInfo["amount"] = amount
		orderInfo["people"] = groupPeople
		orderInfo["title"] = "协同账号购买"
		vipOrderType = 3
	} else {
		totalPeople := userVip.SubNum + 1
		amount = utils.FriendlyFloat64(dyVipValue[buyDays] * float64(totalPeople))
		orderInfo["buy_days"] = buyDays
		orderInfo["amount"] = amount
		orderInfo["people"] = totalPeople
		orderInfo["title"] = "团队成员续费"
		vipOrderType = 4
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
		GoodsInfo:      string(orderInfoJson),
		ExpirationTime: time.Now().Add(1800),
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
	exp := vipOrder.ExpirationTime.Unix() - time.Now().Unix()
	if channel == "native" {
		codeUrl, err := payer.NativePay(amountInt, vipOrder.TradeNo, vipOrder.Title, time.Duration(exp))
		if err != nil {
			receiver.FailReturn(global.NewError(5000))
			return
		}
		receiver.SuccReturn(map[string]interface{}{
			"code_url": codeUrl,
		})
		return
	}
	prepayId, err := payer.AppPay(amountInt, vipOrder.TradeNo, vipOrder.Title, time.Duration(exp))
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
	receiver.SuccReturn(map[string]interface{}{
		"appid":      appId,
		"partnerid ": mchId,
		"prepayid ":  prepayId,
		"package":    "Sign=WXPay",
		"noncestr":   nonceStr,
		"timestamp":  timestamp,
		"sign":       sign,
	})
	return
}
