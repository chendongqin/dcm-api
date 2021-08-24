package controllers

import (
	"dongchamao/business"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/services/payer"
	"github.com/astaxie/beego/logs"
	"github.com/go-pay/gopay/alipay"
	"net/http"
	"time"
)

type CallbackController struct {
	ApiBaseController
}

//微信支付抖音账号回调
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
				logs.Error("微信支付回调重复")
				receiver.Data["json"] = map[string]interface{}{
					"code":    "SUCCESS",
					"message": "成功",
				}
				receiver.ServeJSON()
				return
			}
			if !global.IsDev() {
				amount := float64(payNotifyContent.Amount.Total) / float64(100)
				if utils.ToFloat64(vipOrder.Amount) != amount {
					logs.Error("支付金额与订单金额不匹配")
					receiver.Data["json"] = map[string]interface{}{
						"code":    "SUCCESS",
						"message": "成功",
					}
					receiver.ServeJSON()
					return
				}
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
			if vipOrder.Platform == "douyin" {
				doRes := payBusiness.DoPayDyCallback(vipOrder)
				if !doRes {
					logs.Error("会员数据更新失败：", vipOrder.Id)
				}
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

//支付宝抖音账号回调
func (receiver *CallbackController) AlipayNotify() {
	logStr := "====AliPayNotify====" + string(receiver.Ctx.Input.RequestBody)
	logs.Info(logStr)

	publicKey := global.Cfg.String("ali_pay_cert_path_ali")

	notifyReq, _ := alipay.ParseNotifyResult(receiver.Ctx.Request)
	checkSign, _ := alipay.VerifySignWithCert(publicKey, notifyReq)
	w := receiver.Ctx.ResponseWriter
	w.WriteHeader(http.StatusOK)
	if !checkSign {
		logStr = "=======签名错误=====" + string(receiver.Ctx.Input.RequestBody)
		logs.Error(logStr)
		_, _ = w.Write([]byte("fail"))
		return
	}
	if notifyReq.TradeStatus == "TRADE_SUCCESS" {
		vipOrder := dcm.DcVipOrder{}
		exist, _ := dcm.GetBy("trade_no", notifyReq.OutTradeNo, &vipOrder)
		if exist {
			if vipOrder.PayStatus == 1 {
				logs.Error("微信支付回调重复")
				_, _ = w.Write([]byte("success"))
				return
			}
			if !global.IsDev() {
				if vipOrder.Amount != notifyReq.BuyerPayAmount {
					logs.Error("支付金额与订单金额不匹配")
					_, _ = w.Write([]byte("success"))
					return
				}
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
			if vipOrder.Platform == "douyin" {
				doRes := payBusiness.DoPayDyCallback(vipOrder)
				if !doRes {
					logs.Error("会员数据更新失败：", vipOrder.Id)
				}
			}
		}
	}
	_, _ = w.Write([]byte("success"))
	return
}
