package payer

import (
	"dongchamao/global"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
)

//alipay web
func AliTradePagePay(amount float64, OrderSn, subject, notifyUrl, returnUrl, timeoutExpress string) (string, error) {
	c, err := NewAliPayClient(notifyUrl, returnUrl)
	if err != nil {
		return "", err
	}
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", subject)
	body.Set("out_trade_no", OrderSn)
	body.Set("total_amount", amount)
	body.Set("product_code", "FAST_INSTANT_TRADE_PAY")
	if timeoutExpress != "" {
		body.Set("timeout_express", timeoutExpress)
	}
	//电脑网站支付请求
	payUrl, err := c.TradePagePay(body)
	return payUrl, err
}

//alipay APP
func AliTradeAppPay(amount float64, OrderSn, subject, notifyUrl, returnUrl, timeoutExpress string) (string, error) {
	c, err := NewAliPayClient(notifyUrl, returnUrl)
	if err != nil {
		return "", err
	}
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", subject)
	body.Set("out_trade_no", OrderSn)
	body.Set("total_amount", amount)
	if timeoutExpress != "" {
		body.Set("timeout_express", timeoutExpress)
	}
	//电脑网站支付请求
	payParam, err := c.TradeAppPay(body)
	return payParam, err
}

//alipay wap
func AliTradeWapPay(amount float64, OrderSn, subject, notifyUrl, returnUrl, timeoutExpress string) (string, error) {
	c, err := NewAliPayClient(notifyUrl, returnUrl)
	if err != nil {
		return "", err
	}
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", subject)
	body.Set("out_trade_no", OrderSn)
	body.Set("total_amount", amount)
	quitUrl := global.Cfg.String("ali_pay_return_url")
	body.Set("quit_url", quitUrl)
	body.Set("product_code", "QUICK_WAP_WAY")
	if timeoutExpress != "" {
		body.Set("timeout_express", timeoutExpress)
	}
	//电脑网站支付请求
	payUrl, err := c.TradeWapPay(body)
	return payUrl, err
}

//退款操作
func AliPayRefund(refundAmount float64, orderSn, refundSn, notifyUrl, returnUrl string) (*alipay.TradeRefundResponse, error) {
	c, err := NewAliPayClient(notifyUrl, returnUrl)
	if err != nil {
		return nil, err
	}
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", orderSn)
	body.Set("out_request_no", refundSn)
	body.Set("refund_amount", refundAmount)
	return c.TradeRefund(body)
}

//支付宝生成订单的client
func NewAliPayClient(notifyUrl, returnUrl string) (*alipay.Client, error) {
	if returnUrl == "" {
		returnUrl = global.Cfg.String("ali_pay_return_url")
	}
	appId := global.Cfg.String("ali_pay_appid")
	privateKey := global.Cfg.String("ali_pay_csr_app_private_key")
	certPathRoot := global.Cfg.String("ali_pay_cert_path_root")
	certPathApp := global.Cfg.String("ali_pay_cert_path_app")
	certPathAli := global.Cfg.String("ali_pay_cert_path_ali")
	notifyUrl = global.Cfg.String("pay_notify_url") + notifyUrl
	client := alipay.NewClient(appId, privateKey, true)
	//配置公共参数
	err := client.
		SetPrivateKeyType(alipay.PKCS8).
		SetNotifyUrl(notifyUrl).
		SetReturnUrl(returnUrl).
		SetCertSnByPath(certPathApp, certPathRoot, certPathAli)

	return client, err
}

//支付宝  查询订单详情
func AliTradeQuery(orderSn, notifyUrl, returnUrl string) (*alipay.TradeQueryResponse, error) {
	c, err := NewAliPayClient(notifyUrl, returnUrl)
	if err != nil {
		return nil, err
	}
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", orderSn)

	//电脑网站支付请求
	resp, err := c.TradeQuery(body)
	return resp, err
}
