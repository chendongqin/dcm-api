package controllers

import (
	"dongchamao/models/dcm"
	"dongchamao/services/payer"
	"github.com/astaxie/beego/logs"
	"time"
)

type CallbackController struct {
	ApiBaseController
}

func (receiver *CallbackController) WechatNotify() {
	payNotifyContent := &payer.PayNotifyContent{}
	_, payNotifyContent, err := payer.Notify(receiver.Ctx.Request)
	if err != nil {
		logs.Error("微信支付回调数据错误：", receiver.Ctx.Request.Header, receiver.Ctx.Request.Body)
	}
	logs.Error("微信支付回调数据：", payNotifyContent)
	if payNotifyContent.TradeState == "SUCCESS" {
		vipOrder := dcm.DcVipOrder{}
		exist, _ := dcm.GetBy("trade_no", payNotifyContent.OutTradeNo, &vipOrder)
		if exist {
			payTime, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", payNotifyContent.SuccessTime, time.Local)
			updateData := map[string]interface{}{
				"pay_status":     1,
				"inter_trade_no": payNotifyContent.TransactionId,
				"pay_time":       payTime,
			}
			affect, err2 := dcm.UpdateInfo(nil, vipOrder.Id, updateData, new(dcm.DcVipOrder))
			if affect == 0 || err2 != nil {
				logs.Error("微信支付更新失败：", vipOrder.Id, updateData)
			}
		}
	}
	receiver.Data["json"] = map[string]interface{}{
		"code":    "SUCCESS",
		"message": "成功",
	}
	receiver.ServeJSON()
	return

}
