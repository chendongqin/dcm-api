package business

import (
	"crypto/md5"
	"dongchamao/global"
	"dongchamao/global/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	INTERAPPID    = "20000"
	API_HOST_TEST = "https://api.dongchamao.cn"
	API_HOST      = "https://api.dongchamao.com"
)

type InterBusiness struct {
}

func NewInterBusiness() *InterBusiness {
	return new(InterBusiness)
}

type JsonRet struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

func (i *InterBusiness) BuildSignHeader(req *http.Request) global.CommonError {
	timestamp := time.Now().Unix()
	random := utils.GetRandomString(6)
	secret, _ := NewAccountAuthBusiness().GetAppSecret(INTERAPPID, true)
	if secret == "" {
		return global.NewError(4000)
	}
	tmpStr := fmt.Sprintf("%d%s%s", timestamp, random, secret)
	h := md5.New()
	h.Write([]byte(tmpStr))
	sign := hex.EncodeToString(h.Sum(nil))
	req.Header.Add("APPID", INTERAPPID)
	req.Header.Add("TIMESTAMP", utils.ToString(timestamp))
	req.Header.Add("RANDOM", random)
	req.Header.Add("SIGN", sign)
	return nil
}

func (i *InterBusiness) GetHost() string {
	if global.IsDev() {
		return API_HOST_TEST
	} else {
		return API_HOST
	}
}

func (i *InterBusiness) GetSimple(url string) (interface{}, global.CommonError) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil) //建立一个请求
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	comErr := i.BuildSignHeader(req)
	if comErr != nil {
		return nil, comErr
	}
	resp, err := client.Do(req) //提交
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	retStruct := new(JsonRet)
	err = jsoniter.Unmarshal(body, retStruct)
	fmt.Println(retStruct)
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	if retStruct.Status {
		return retStruct.Data, nil
	} else {
		return nil, global.NewError(retStruct.Code)
	}
}

func (i *InterBusiness) HttpGet(url string) (interface{}, global.CommonError) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil) //建立一个请求
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	comErr := i.BuildSignHeader(req)
	if comErr != nil {
		return nil, comErr
	}
	resp, err := client.Do(req) //提交
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}

	retStruct := new(JsonRet)
	err = json.Unmarshal(body, retStruct)
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	if retStruct.Status {
		return retStruct.Data, nil
	} else {
		return nil, global.NewError(retStruct.Code)
	}
}

func (i *InterBusiness) HttpPost(url string, postJson string) (interface{}, global.CommonError) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(postJson)) //建立一个请求
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	comErr := i.BuildSignHeader(req)
	if comErr != nil {
		return nil, comErr
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req) //提交
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	retStruct := new(JsonRet)
	err = json.Unmarshal(body, retStruct)
	if err != nil {
		return nil, global.NewMsgError(err.Error())
	}
	if retStruct.Status {
		return retStruct.Data, nil
	} else {
		return nil, global.NewError(retStruct.Code)
	}
}

func (i *InterBusiness) BuildURL(path string) string {
	return i.GetHost() + "/" + strings.TrimLeft(path, "/")
}

func (i *InterBusiness) BuildGetURL(path string, params map[string]interface{}) string {
	baseURL := i.BuildURL(path)
	keys := make([]string, len(params))
	index := 0
	for key, value := range params {
		keys[index] = key + "=" + utils.ToString(value)
		index++
	}
	return baseURL + "?" + strings.Join(keys, "&")
}
