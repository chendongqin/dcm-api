package payer

//
//import (
//	"context"
//	"crypto/rsa"
//	"crypto/x509"
//	"dongchamao/global"
//	"fmt"
//	"github.com/wechatpay-apiv3/wechatpay-go/core"
//	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
//	"github.com/wechatpay-apiv3/wechatpay-go/utils"
//	"log"
//	"net/http"
//)
//
//func Builder()  {
//	//var (
//	//	mchID                      string              // 商户号
//	//	mchCertificateSerialNumber string              // 商户证书序列号
//	//	mchPrivateKey              *rsa.PrivateKey     // 商户私钥
//	//	wechatPayCertList          []*x509.Certificate // 平台证书列表
//	//	customHTTPClient           *http.Client        // 可选，自定义客户端实例
//	//)
//	mchPrivateKeyPath := global.Cfg.String("wechat_pay_cert_key")
//	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(mchPrivateKeyPath)
//	if err != nil {
//		return nil, fmt.Errorf("商户私钥有误：%v", err)
//	}
//
//	ctx := context.Background()
//	opts := []core.ClientOption{
//		option.WithWechatPayAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, wechatPayCertList),
//		option.WithHTTPClient(customHTTPClient),
//	}
//	client, err := core.NewClient(ctx, opts...)
//	if err != nil {
//		log.Printf("new wechat pay client err:%s", err)
//		return
//	}
//	_ = client
//}
