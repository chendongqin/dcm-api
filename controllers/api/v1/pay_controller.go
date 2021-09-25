package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/models/repost"
	"dongchamao/models/repost/dy"
	"dongchamao/services/payer"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"math"
	"net/url"
	"time"
)

type PayController struct {
	controllers.ApiBaseController
}

func (receiver *PayController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
}

func (receiver *PayController) DySurplusValue() {
	var vip dcm.DcUserVip
	if _, err := dcm.GetDbSession().Where("user_id=? AND platform=1", receiver.UserId).Get(&vip); err != nil {
		return
	}
	surplusDay := vip.Expiration.Sub(time.Now()).Hours() / 24
	if surplusDay <= 0 || !receiver.HasAuth {
		receiver.SuccReturn(map[string]interface{}{
			"now_surplus_day": 0,
			"now_value":       0,
			"value":           0,
			"prime_value":     0,
			"price_config":    0,
		})
		return
	}
	//当前团队续费金额
	total := business.NewVipBusiness().GetVipLevel(receiver.UserId, 1).SubNum
	var nowValue float64
	var nowSurplusDay float64
	payBusiness := business.NewPayBusiness()
	//扩张团队单人价格
	value, primeValue := payBusiness.GetDySurplusValue(int(math.Ceil(surplusDay)))
	if total > 0 {
		//团队过期后续费重新计算时间
		startTime := vip.SubExpiration
		if startTime.Before(time.Now()) {
			nowValue = value
		} else {
			if vip.Expiration != startTime && vip.Expiration.After(time.Now()) {
				valueAdd, _ := payBusiness.GetDyAddValue(int(math.Ceil(surplusDay)))
				subTime := vip.Expiration.Sub(startTime)
				nowSurplusDay = subTime.Hours() / 24
				surplusUnit := valueAdd / math.Ceil(surplusDay)
				nowValue = utils.CeilFloat64One(surplusUnit * (math.Ceil(nowSurplusDay)))
			}
		}
	}
	//获取价格配置
	priceConfig := payBusiness.GetVipPrice()
	receiver.SuccReturn(map[string]interface{}{
		"now_surplus_day": int(math.Ceil(nowSurplusDay)),
		"now_value":       utils.CeilFloat64One(nowValue * float64(total)),
		"value":           utils.CeilFloat64One(value),
		"prime_value":     primeValue,
		"price_config":    priceConfig,
		"surplus_day":     math.Ceil(nowSurplusDay),
	})
	return
}

//抖音会员价格
func (receiver *PayController) DyPriceList() {
	payBusiness := business.NewPayBusiness()
	priceData := payBusiness.GetVipPriceConfigCheckActivity(receiver.UserId, true)
	receiver.SuccReturn(priceData)
	return
}

//创建抖音会员订单
func (receiver *PayController) CreateDyOrder() {
	if !business.UserActionLock("vip_order", utils.ToString(receiver.UserId), 2) {
		receiver.FailReturn(global.NewError(4211))
		return
	}
	InputData := receiver.InputFormat()
	orderType := InputData.GetInt("order_type", 0)
	referrer := InputData.GetString("referrer", "")
	groupPeople := InputData.GetInt("group_people", 0)
	buyDays := InputData.GetInt("days", 0)
	if utils.InArrayInt(orderType, []int{1, 5}) && !utils.InArrayInt(buyDays, []int{30, 180, 365}) {
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
	if orderType == 2 {
		subExpiration = time.Now()
	} else {
		if subExpiration.Before(time.Now()) {
			subExpiration = time.Now()
		}
	}
	var remark = ""
	var surplusDay int64 = 0
	var surplusUnit float64 = 0
	if userVip.Expiration.After(time.Now()) {
		surplusDay = int64(math.Ceil(userVip.Expiration.Sub(subExpiration).Hours() / 24))
		if surplusDay == 0 && orderType == 3 {
			receiver.FailReturn(global.NewMsgError("协同账号到期时间与账户会员时间一致，不需要续费～"))
			return
		}
	}
	payBusiness := business.NewPayBusiness()
	var surplusValue float64 = 0
	var trueSurplusValue float64 = 0
	vipOrderType := 1
	if userVip.Level > 0 {
		vipOrderType = 2
		trueSurplusValue, _ = payBusiness.GetDyAddValue(int(surplusDay))
		surplusValue, _ = payBusiness.GetDySurplusValue(int(surplusDay))
		if surplusDay > 0 {
			surplusUnit = trueSurplusValue / float64(surplusDay)
		}
	}
	checkActivity := false
	if orderType == 1 {
		checkActivity = true
	}
	priceData := payBusiness.GetVipPriceConfigCheckActivity(receiver.UserId, checkActivity)
	price := dy.VipPriceActive{}
	if buyDays == 30 {
		price = priceData.Month
	} else if buyDays == 180 {
		price = priceData.HalfYear
	} else if buyDays == 365 {
		price = priceData.Year
	}
	title := fmt.Sprintf("专业版%d天", buyDays)
	var amount float64 = 0
	orderInfo := repost.VipOrderInfo{
		SurplusValue: surplusValue,
	}
	//购买会员
	if orderType == 1 {
		amount = price.Price
		orderInfo.BuyDays = buyDays
		orderInfo.Amount = amount
		orderInfo.People = 1
		orderInfo.Title = "会员购买"
		remark = price.ActiveComment
	} else if orderType == 2 { //购买协同账号
		title = fmt.Sprintf("购买协同账号%d人", groupPeople)
		amount = surplusValue * float64(groupPeople)
		//先续费再购买
		if userVip.SubNum > 0 {
			if userVip.SubExpiration.Before(time.Now()) {
				amount += trueSurplusValue * float64(userVip.SubNum)
			} else {
				if userVip.Expiration.After(userVip.SubExpiration) {
					surplusSubDay := math.Ceil((userVip.Expiration.Sub(userVip.SubExpiration)).Hours() / 24)
					surplusSubsValue := surplusUnit * surplusSubDay
					tmpAmount := float64(userVip.SubNum) * utils.CeilFloat64One(surplusSubsValue)
					remark = fmt.Sprintf("已有子账号续费：%.1f元", tmpAmount)
					amount += tmpAmount
				}
			}
		}
		amount = utils.CeilFloat64One(amount)
		orderInfo.BuyDays = int(surplusDay)
		orderInfo.Amount = amount
		orderInfo.People = groupPeople
		orderInfo.Title = "协同账号购买"
		vipOrderType = 3
	} else if orderType == 3 { //协同账号续费
		totalPeople := userVip.SubNum
		title = fmt.Sprintf("协同账号续费%d人", totalPeople)
		amount = utils.CeilFloat64One(surplusValue * float64(totalPeople))
		orderInfo.BuyDays = int(surplusDay)
		orderInfo.Amount = amount
		orderInfo.People = totalPeople
		orderInfo.Title = "协同账号续费"
		vipOrderType = 4
	} else if orderType == 4 {
		title = "团队成员续费"
		totalPeople := userVip.SubNum + 1
		amount = utils.CeilFloat64One(price.Price * float64(totalPeople))
		remark = price.ActiveComment
		//先续费再购买
		if userVip.SubNum > 0 {
			if userVip.SubExpiration.Before(time.Now()) {
				amount += trueSurplusValue * float64(userVip.SubNum)
			} else {
				if userVip.Expiration.After(userVip.SubExpiration) {
					surplusSubDay := math.Ceil((userVip.Expiration.Sub(userVip.SubExpiration)).Hours() / 24)
					surplusSubsValue := surplusUnit * surplusSubDay
					tmpAmount := float64(userVip.SubNum) * utils.CeilFloat64One(surplusSubsValue)
					remark = fmt.Sprintf("已有子账号续费：%.1f元", tmpAmount)
					amount += tmpAmount
				}
			}
		}
		orderInfo.BuyDays = buyDays
		orderInfo.Amount = utils.FriendlyFloat64(amount)
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
		BuyDays:        buyDays,
		GoodsInfo:      string(orderInfoJson),
		Referrer:       referrer,
		ExpirationTime: time.Now().Add(1800 * time.Second),
		Remark:         remark,
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

//创建抖音直播监控订单
func (receiver *PayController) CreateDyMonitorOrder() {
	if !business.UserActionLock("monitor_order", utils.ToString(receiver.UserId), 2) {
		receiver.FailReturn(global.NewError(4211))
		return
	}
	InputData := receiver.InputFormat()
	number := InputData.GetInt("number", 0)
	if !utils.InArrayInt(number, []int{10, 100, 500}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	priceString := business.NewConfigBusiness().GetConfigJson("monitor_price", false)
	priceList := dy.LiveMonitorPriceList{}
	_ = jsoniter.Unmarshal([]byte(priceString), &priceList)
	var amount float64 = 0
	if utils.ToInt(priceList.MonitorPrice.Price10.Monitor) == number {
		amount = utils.ToFloat64(priceList.MonitorPrice.Price10.Price)
	} else if utils.ToInt(priceList.MonitorPrice.Price100.Monitor) == number {
		amount = utils.ToFloat64(priceList.MonitorPrice.Price100.Price)
	} else if utils.ToInt(priceList.MonitorPrice.Price500.Monitor) == number {
		amount = utils.ToFloat64(priceList.MonitorPrice.Price500.Price)
	}
	if amount == 0 {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	title := fmt.Sprintf("购买直播监控%d次", number)
	orderInfo := repost.VipOrderInfo{
		Title:      title,
		MonitorNum: number,
	}
	uniqueID, _ := utils.Snow.GetSnowflakeId()
	tradeNo := fmt.Sprintf("%s%d", time.Now().Format("060102"), uniqueID)
	orderInfoJson, _ := jsoniter.Marshal(orderInfo)
	vipOrder := dcm.DcVipOrder{
		UserId:         receiver.UserId,
		Username:       receiver.UserInfo.Username,
		TradeNo:        tradeNo,
		OrderType:      7,
		Platform:       "douyin",
		Title:          title,
		Amount:         utils.ToString(amount),
		TicketAmount:   "0",
		Level:          0,
		BuyDays:        orderInfo.BuyDays,
		GoodsInfo:      string(orderInfoJson),
		Referrer:       "",
		ExpirationTime: time.Now().Add(1800 * time.Second),
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}
	_, err := dcm.Insert(nil, &vipOrder)
	if vipOrder.Id == 0 || err != nil {
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

//订单详情
func (receiver *PayController) OrderDetail() {
	orderId := utils.ToInt(receiver.Ctx.Input.Param(":order_id"))
	if orderId == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	vipOrder := dcm.DcVipOrder{}
	exist, _ := dcm.Get(orderId, &vipOrder)
	if !exist || vipOrder.UserId != receiver.UserId {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	status := vipOrder.Status
	if vipOrder.PayStatus == 0 && vipOrder.ExpirationTime.Before(time.Now()) {
		status = 2
	}
	orderDetail := repost.VipOrderDetail{
		OrderId:      orderId,
		TradeNo:      vipOrder.TradeNo,
		OrderType:    vipOrder.OrderType,
		PayType:      vipOrder.PayType,
		Level:        vipOrder.Level,
		BuyDays:      vipOrder.BuyDays,
		Title:        vipOrder.Title,
		Amount:       vipOrder.Amount,
		TicketAmount: vipOrder.TicketAmount,
		Status:       status,
		PayStatus:    vipOrder.PayStatus,
		CreateTime:   vipOrder.CreateTime.Format("2006-01-02 15:04:05"),
		PayTime:      vipOrder.PayTime.Format("2006-01-02 15:04:05"),
		InvoiceId:    vipOrder.InvoiceId,
	}
	receiver.SuccReturn(map[string]interface{}{
		"detail": orderDetail,
	})
}

//订单删除
func (receiver *PayController) OrderDel() {
	orderId := utils.ToInt(receiver.Ctx.Input.Param(":order_id"))
	if orderId == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	vipOrder := dcm.DcVipOrder{}
	exist, _ := dcm.Get(orderId, &vipOrder)
	if !exist || vipOrder.UserId != receiver.UserId {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	affect, err := dcm.UpdateInfo(nil, orderId, map[string]interface{}{"status": -1}, new(dcm.DcVipOrder))
	if affect == 0 || err != nil {
		receiver.FailReturn(global.NewError(500))
		return
	}
	receiver.SuccReturn(nil)
}

//订单列表
func (receiver *PayController) OrderList() {
	platform := receiver.Ctx.Input.Param(":platform")
	selectStatus, _ := receiver.GetInt("select_status", 0)
	invoiceStatus, _ := receiver.GetInt("invoice_status", 0)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 30)
	vipOrderList := make([]dcm.DcVipOrder, 0)
	start := (page - 1) * pageSize
	sql := fmt.Sprintf("user_id=%d AND platform='%s'", receiver.UserId, platform)
	if selectStatus == 1 {
		sql += " AND pay_status = 1 "
	} else if selectStatus == 2 {
		sql += fmt.Sprintf(" AND pay_status = 0 AND expiration_time < '%s'", time.Now().Format("2006-01-02 15:04:05"))
	} else if selectStatus == 3 {
		sql += fmt.Sprintf(" AND pay_status = 0 AND expiration_time >= '%s'", time.Now().Format("2006-01-02 15:04:05"))
	}
	if invoiceStatus == 1 {
		sql += " AND invoice_id = 0"
	} else if invoiceStatus == 2 {
		sql += " AND invoice_id > 0"
	}
	total, _ := dcm.GetSlaveDbSession().
		Where(sql).
		Where("status!=-1").
		Limit(pageSize, start).
		Desc("create_time").
		FindAndCount(&vipOrderList)
	list := make([]repost.VipOrderDetail, 0)
	for _, v := range vipOrderList {
		status := v.Status
		if v.PayStatus == 0 && v.ExpirationTime.Before(time.Now()) {
			status = 2
		}
		list = append(list, repost.VipOrderDetail{
			OrderId:      v.Id,
			TradeNo:      v.TradeNo,
			OrderType:    v.OrderType,
			PayType:      v.PayType,
			Level:        v.Level,
			BuyDays:      v.BuyDays,
			Title:        v.Title,
			Amount:       v.Amount,
			Channel:      v.Channel,
			TicketAmount: v.TicketAmount,
			Status:       status,
			PayStatus:    v.PayStatus,
			CreateTime:   v.CreateTime.Format("2006-01-02 15:04:05"),
			PayTime:      v.PayTime.Format("2006-01-02 15:04:05"),
			InvoiceId:    v.InvoiceId,
		})
	}

	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
}
