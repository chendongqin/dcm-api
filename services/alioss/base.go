package alioss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego"
	"hash"
	"io"
	"time"
)

// 请填写您的AccessKeyId。
var accessKeyId string = "LTAI4G13UrqKAYH7ngd17X9C"

// 请填写您的AccessKeySecret。
var accessKeySecret string = "Dtovz4xb9je7gVHGQVQ1Qmi3UimqKZ"

// host的格式为 bucketname.endpoint ，请替换为您的真实信息。

var Host = beego.AppConfig.String("oss_url")

// 用户上传文件时指定的前缀。
// var uploadDir string = "user-dir-prefix/"
var expireTime int64 = 30

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func getGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}


func GetUserAvatarPolicyToken(uploadDir, tmpKey, upType string, uid int64) *PolicyToken {
	callBackBody := fmt.Sprintf("filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}&tmp_key=%s&up_type=%s&uid=%d", tmpKey, upType, uid)
	avatarCallBackUrl := beego.AppConfig.String("oss_callback_url")
	// fmt.Println(avatarCallBackUrl)
	return commUploadPolicyToken(uploadDir, callBackBody, avatarCallBackUrl)
}

func commUploadPolicyToken(uploadDir, callBackBody, callBackUrlstring string) *PolicyToken {
	now := time.Now().Unix()
	expireEnd := now + expireTime
	var tokenExpire = getGmtIso8601(expireEnd)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, uploadDir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, _ := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = callBackUrlstring
	callbackParam.CallbackBody = callBackBody
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callbackStr, _ := json.Marshal(callbackParam)

	callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)

	policyToken := &PolicyToken{
		AccessKeyId: accessKeyId,
		Host:        Host,
		Expire:      expireEnd,
		Signature:   signedStr,
		Directory:   uploadDir,
		Policy:      debyte,
		Callback:    callbackBase64,
	}

	return policyToken
}

func UploadLocalFile(localFile, uploadUrl string) error {
	client, err := oss.New("https://oss-cn-hangzhou.aliyuncs.com", accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	bucketName := beego.AppConfig.String("bucketName")
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}
	err = bucket.PutObjectFromFile(uploadUrl, localFile)
	if err != nil {
		return err
	}
	return nil
}
