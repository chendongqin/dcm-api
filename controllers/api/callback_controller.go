package controllers

import (
	"dongchamao/business"
	"dongchamao/global"
	"dongchamao/models/dcm"
	"dongchamao/services/payer"
	"github.com/astaxie/beego/logs"
	"github.com/iGoogle-ink/gopay/alipay"
	"net/http"
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
	logs.Info("微信支付回调数据：", payNotifyContent)
	if payNotifyContent.TradeState == "SUCCESS" {
		vipOrder := dcm.DcVipOrder{}
		exist, _ := dcm.GetBy("trade_no", payNotifyContent.OutTradeNo, &vipOrder)
		if exist {
			if vipOrder.PayStatus == 1 {
				logs.Error("微信支付回调：", receiver.Ctx.Request.Header, receiver.Ctx.Request.Body)
			}
			payTime, _ := time.Parse("2006-01-02T15:04:05+08:00", payNotifyContent.SuccessTime)
			updateData := map[string]interface{}{
				"pay_status":     1,
				"status":         1,
				"pay_type":       "wechat",
				"inter_trade_no": payNotifyContent.TransactionId,
				"pay_time":       payTime.Format("2006-01-02 15:04:05"),
			}
			affect, err2 := dcm.UpdateInfo(nil, vipOrder.Id, updateData, new(dcm.DcVipOrder))
			if affect == 0 || err2 != nil {
				logs.Error("微信支付更新失败：", vipOrder.Id, updateData, err2)
			}
			payBusiness := business.NewPayBusiness()
			doRes := payBusiness.DoPayDyCallback(vipOrder)
			if !doRes {
				logs.Error("会员数据更新失败：", vipOrder.Id)
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

func (receiver *CallbackController) ApipayNotify() {
	logStr := "====AliPayNotify====" + string(receiver.Ctx.Input.RequestBody)
	logs.Info(logStr)

	publicKey := global.Cfg.String("ali_pay_public_key")

	notifyReq, _ := alipay.ParseNotifyResult(receiver.Ctx.Request)
	checkSign, _ := alipay.VerifySign(publicKey, notifyReq)
	w := receiver.Ctx.ResponseWriter
	w.WriteHeader(http.StatusOK)
	if !checkSign {
		logStr = "=======签名错误=====" + string(receiver.Ctx.Input.RequestBody)
		logs.Error(logStr)
		_, _ = w.Write([]byte("fail"))
		return
	}
	vipOrder := dcm.DcVipOrder{}
	exist, _ := dcm.GetBy("trade_no", notifyReq.OutTradeNo, &vipOrder)
	if exist {
		if vipOrder.PayStatus == 1 {
			logs.Error("微信支付回调：", receiver.Ctx.Request.Header, receiver.Ctx.Request.Body)
		}
		payTime := notifyReq.NotifyTime
		updateData := map[string]interface{}{
			"pay_status":     1,
			"status":         1,
			"pay_type":       "alipay",
			"inter_trade_no": notifyReq.TradeNo,
			"pay_time":       payTime,
		}
		affect, err2 := dcm.UpdateInfo(nil, vipOrder.Id, updateData, new(dcm.DcVipOrder))
		if affect == 0 || err2 != nil {
			logs.Error("支付宝支付更新失败：", vipOrder.Id, updateData, err2)
		}
		payBusiness := business.NewPayBusiness()
		doRes := payBusiness.DoPayDyCallback(vipOrder)
		if !doRes {
			logs.Error("会员数据更新失败：", vipOrder.Id)
		}
	}
	_, _ = w.Write([]byte("success"))
}
