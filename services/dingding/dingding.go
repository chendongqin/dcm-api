package dingding

import (
	"dongchamao/global/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/**
 * dingding send
 * demo
 * err := dingding.NewDingDing().SetTokenUrl("https://oapi.dingtalk.com/robot/send?access_token=50dd4a79339b6b184c57dd45db26208cf594c28d4a01d185f0f052cc70627b48").SendMarkDown("123","456")
 */

type DDMdtemplate struct {
	Msgtype  string                 `json:"msgtype"`
	Markdown map[string]interface{} `json:"markdown"`
	At       DDAt                   `json:"at"`
}

type DDAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DDRet struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type dingDing struct {
	tokenUrl string
	at       DDAt
}

func NewDingDing() *dingDing {
	return new(dingDing)
}

func NewWithTokenUrl(tokenUrl string) *dingDing {
	ding := NewDingDing()
	ding.SetTokenUrl(tokenUrl)
	return ding
}

func (this *dingDing) SetTokenUrl(tokenUrl string) *dingDing {
	this.tokenUrl = tokenUrl
	return this
}

func (this *dingDing) SetAtAll() *dingDing {
	this.at.IsAtAll = true
	return this
}

func (this *dingDing) SetAt(data ...string) *dingDing {
	if this.at.IsAtAll == false {
		this.at.AtMobiles = data
	}
	return this
}

func (this *dingDing) SendMarkDown(title, content string) error {
	if this.tokenUrl == "" {
		return errors.New("webhook不能为空")
	}
	contentData := make(map[string]interface{}, 0)
	contentData["title"] = title
	contentData["text"] = content
	tpl := DDMdtemplate{
		"markdown",
		contentData,
		this.at,
	}
	if jsonData, err := json.Marshal(tpl); err == nil {
		fmt.Println(string(jsonData))
		ret, err := CurlData(this.tokenUrl, "POST", string(jsonData), "application/json")
		if logger.CheckError(err) != nil {
			return err
		}
		return this._retHandle(ret)
	} else {
		return err
	}
}

func (this *dingDing) _retHandle(ret string) error {
	var retCurl DDRet
	if err := json.Unmarshal([]byte(ret), &retCurl); err == nil {
		if retCurl.Errcode == 0 {
			return nil
		} else {
			return errors.New(retCurl.Errmsg)
		}
	} else {
		return err
	}
}

func CurlData(url string, method string, postData string, contentType string) (string, error) {
	var resp *http.Response
	var err error
	if method == "POST" {
		resp, err = http.Post(url, contentType, strings.NewReader(postData))
	} else {
		resp, err = http.Get(url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if logger.CheckError(err) != nil {
		return "", err
	}
	return string(body), nil
}
