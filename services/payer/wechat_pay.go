package payer

import (
	"context"
	"crypto/x509"
	"dongchamao/global"
	"errors"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"net/http"
	"time"
)

func Builder() (*core.Client, error) {
	mchPrivateKeyPath := global.Cfg.String("wechat_pay_cert_key")
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(mchPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("商户配置有误：%v", err)
	}
	mchID := global.Cfg.String("wechat_pay_mchid")
	mchCertificateSerialNumber := global.Cfg.String("wechat_pay_cert_sn")
	wechatPayCertPath := global.Cfg.String("wechat_pay_cert")
	wechatPayCert, err := utils.LoadCertificateWithPath(wechatPayCertPath)
	if err != nil {
		return nil, fmt.Errorf("商户配置有误：%v", err)
	}
	wechatPayCertList := []*x509.Certificate{wechatPayCert}
	customHTTPClient := &http.Client{}
	ctx := context.Background()
	opts := []core.ClientOption{
		option.WithWechatPayAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, wechatPayCertList),
		option.WithHTTPClient(customHTTPClient),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func AppPay(amount int64, tradeNo, description string, expire time.Duration) (*string, error) {
	client, err := Builder()
	svc := app.AppApiService{Client: client}
	ctx := context.Background()
	appID := global.Cfg.String("wechat_pay_app_appId")
	mchId := global.Cfg.String("wechat_pay_mchid")
	notifyUrl := global.Cfg.String("wechat_pay_notify_url")
	resp, result, err := svc.Prepay(ctx,
		app.PrepayRequest{
			Appid:         core.String(appID),
			Mchid:         core.String(mchId),
			Description:   core.String(description),
			OutTradeNo:    core.String(tradeNo),
			TimeExpire:    core.Time(time.Now().Add(expire)),
			NotifyUrl:     core.String(notifyUrl),
			SupportFapiao: core.Bool(false),
			Amount: &app.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(amount),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != 200 {
		return nil, errors.New(result.Response.Status)
	}
	return resp.PrepayId, nil
}

func NativePay(amount int64, tradeNo, description string, expire time.Duration) (*string, error) {
	client, err := Builder()
	svc := native.NativeApiService{Client: client}
	ctx := context.Background()
	appID := global.Cfg.String("wechat_pay_appId")
	mchId := global.Cfg.String("wechat_pay_mchid")
	notifyUrl := global.Cfg.String("wechat_pay_notify_url")
	resp, result, err := svc.Prepay(ctx,
		native.PrepayRequest{
			Appid:         core.String(appID),
			Mchid:         core.String(mchId),
			Description:   core.String(description),
			OutTradeNo:    core.String(tradeNo),
			TimeExpire:    core.Time(time.Now().Add(expire)),
			NotifyUrl:     core.String(notifyUrl),
			SupportFapiao: core.Bool(false),
			Amount: &native.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(amount),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != 200 {
		return nil, errors.New(result.Response.Status)
	}
	return resp.CodeUrl, nil
}

func H5Pay(amount int64, tradeNo, description string, expire time.Duration) (*string, error) {
	client, err := Builder()
	svc := h5.H5ApiService{Client: client}
	ctx := context.Background()
	appID := global.Cfg.String("wechat_pay_appId")
	mchId := global.Cfg.String("wechat_pay_mchid")
	notifyUrl := global.Cfg.String("wechat_pay_notify_url")
	resp, result, err := svc.Prepay(ctx,
		h5.PrepayRequest{
			Appid:       core.String(appID),
			Mchid:       core.String(mchId),
			Description: core.String(description),
			OutTradeNo:  core.String(tradeNo),
			TimeExpire:  core.Time(time.Now().Add(expire)),
			NotifyUrl:   core.String(notifyUrl),
			Amount: &h5.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(100),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != 200 {
		return nil, errors.New(result.Response.Status)
	}
	return resp.H5Url, nil
}

func Sha256WithRsa(rsaStr string) (string, error) {
	mchPrivateKeyPath := global.Cfg.String("wechat_pay_cert_key")
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(mchPrivateKeyPath)
	if err != nil {
		return "", err
	}
	return utils.SignSHA256WithRSA(rsaStr, mchPrivateKey)
}
