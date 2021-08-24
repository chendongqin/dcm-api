package ali_tools

import (
	"dongchamao/global"
	"errors"
	"fmt"
	afs "github.com/alibabacloud-go/afs-20180112/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
)

func IClientProfileClient() (*afs.Client, error) {
	config := new(rpc.Config)
	aliAccessKey := global.Cfg.String("ali_accessKey")
	aliAccessSecret := global.Cfg.String("ali_secret")
	if aliAccessKey == "" || aliAccessSecret == "" {
		return nil, errors.New("配置加载失败")
	}
	config.SetAccessKeyId(aliAccessKey).
		SetAccessKeySecret(aliAccessSecret).
		SetRegionId("cn-hangzhou").
		SetEndpoint("afs.aliyuncs.com")
	client, _ := afs.NewClient(config)
	return client, nil
}

func IClientProfile(sig, sessionId, token, ip, scene, appKey string) error {
	client, err := IClientProfileClient()
	if err != nil {
		return err
	}
	request := new(afs.AuthenticateSigRequest)
	request.SetSig(sig)
	request.SetSessionId(sessionId)
	request.SetToken(token)
	request.SetRemoteIp(ip)
	request.SetScene(scene)
	request.SetAppKey(appKey)
	response, err := client.AuthenticateSig(request)
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}
