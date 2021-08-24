package test

import (
	"dongchamao/services/umeng"
	"fmt"
	"log"
	"testing"
)

func TestIosUMengPush(t *testing.T) {
	extrasInfo := map[string]interface{}{}
	extrasInfo["sai"] = 123
	extrasInfo["lin"] = 123
	extrasInfo["mu"] = 123
	payload := &umeng.IosPayload{
		Aps: struct {
			Alert struct {
				Title    string `json:"title,omitempty"`
				Subtitle string `json:"subtitle,omitempty"`
				Body     string `json:"body,omitempty"`
			} `json:"alert,omitempty"`
			Badge            int64  `json:"badge,omitempty"`
			Sound            string `json:"sound,omitempty"`
			ContentAvailable int64  `json:"content-available,omitempty"`
			Category         string `json:"category,omitempty"`
		}{Alert: struct {
			Title    string `json:"title,omitempty"`
			Subtitle string `json:"subtitle,omitempty"`
			Body     string `json:"body,omitempty"`
		}{
			Title:    "标题123",
			Subtitle: "子标题123",
			Body:     "内容123",
		}},
		Extra: extrasInfo,
	}
	param := &umeng.SendParam{
		Types:        umeng.TypeUnicast,
		DeviceTokens: "77d5b9f9786212d475f6dcd071315a8cfc5b08d4952007d1404475fcf0da5114",
		AliasType:    "",
		Alias:        "",
		FileId:       "",
		Filter:       "",
		Payload:      payload,
		Policy: struct {
			StartTime      string `json:"start_time,omitempty"`
			ExpireTime     string `json:"expire_time,omitempty"`
			MaxSendNum     int64  `json:"max_send_num,omitempty"`
			OutBizNo       string `json:"out_biz_no,omitempty"`
			ApnsCollapseId string `json:"apns_collapse_id,omitempty"`
		}{},
		ProductionMode: "false",
		Description:    "测试推送",
	}
	appKey := "60ed3a2da6f90557b7b591a9"
	appMasterKey := "durb5igu52ziiygozrnnl9tgoyiencjp"
	umengPush := umeng.NewUmengPush(appKey, appMasterKey)
	result, err := umengPush.Send(param)
	fmt.Println(result.Ret)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(result.Data.MsgId, result.Data.ErrorMsg)
}

func TestAndroidUMengPush(t *testing.T) {
	extrasInfo := map[string]interface{}{}
	extrasInfo["sai"] = 123
	extrasInfo["lin"] = 123
	extrasInfo["mu"] = 123
	androidLoad := &umeng.AndroidPayload{
		DisplayType: "notification",
		Body: struct {
			Ticker string `json:"ticker"` // 必填,通知栏提示文字
			Title  string `json:"title"`  // 必填,通知标题
			Text   string `json:"text"`   // 必填,通知文字描述

			// 可选,状态栏图标ID,R.drawable.[smallIcon],
			// 如果没有,默认使用应用图标
			// 图片要求为24*24dp的图标,或24*24px放在drawable-mdpi下
			// 注意四周各留1个dp的空白像素
			Icon string `json:"icon,omitempty"`

			// 可选,通知栏拉开后左侧图标ID,R.drawable.[largeIcon],
			// 图片要求为64*64dp的图标,
			// 可设计一张64*64px放在drawable-mdpi下,
			// 注意图片四周留空,不至于显示太拥挤
			LargeIcon string `json:"largeIcon,omitempty"`

			// 可选,通知栏大图标的URL链接,该字段的优先级大于largeIcon
			// 该字段要求以http或者https开头
			Img string `json:"img,omitempty"`

			// 可选,通知声音,R.raw.[sound]
			// 如果该字段为空,采用SDK默认的声音,即res/raw/下的
			// umeng_push_notification_default_sound声音文件,如果
			// SDK默认声音文件不存在,则使用系统默认Notification提示音
			Sound string `json:"sound,omitempty"`

			BuilderId   string `json:"builder_id,omitempty"`   // 可选,默认为0,用于标识该通知采用的样式,使用该参数时,开发者必须在SDK里面实现自定义通知栏样式
			PlayVibrate string `json:"play_vibrate,omitempty"` // 可选,收到通知是否震动,默认为"true"
			PlayLights  string `json:"play_lights,omitempty"`  // 可选,收到通知是否闪灯,默认为"true"
			PlaySound   string `json:"play_sound,omitempty"`   // 可选,收到通知是否发出声音,默认为"true"

			// 点击"通知"的后续行为,默认为打开app
			// 可选,默认为"go_app",值可以为:
			//   "go_app": 打开应用
			//   "go_url": 跳转到URL
			//   "go_activity": 打开特定的activity
			//   "go_custom": 用户自定义内容
			AfterOpen string `json:"after_open,omitempty"`

			// 当after_open=go_url时,必填
			// 通知栏点击后跳转的URL,要求以http或者https开头
			Url string `json:"url,omitempty"`

			// 当after_open=go_activity时,必填
			// 通知栏点击后打开的Activity
			Activity string `json:"activity,omitempty"`

			// 当display_type=message时, 必填
			// 当display_type=notification且
			// after_open=go_custom时,必填
			Custom interface{} `json:"custom,omitempty"`
		}{Ticker: "测试Ticker", Text: "测试Text", Title: "测试Title"},
		Extra: extrasInfo,
	}
	param := &umeng.SendParam{
		Types:        umeng.TypeUnicast,
		DeviceTokens: "AgIRmlFt7w6t7JJXdayIMsbf4md926PCYbGoEgfBRBLK",
		AliasType:    "",
		Alias:        "",
		FileId:       "",
		Filter:       "",
		Payload:      androidLoad,
		Policy: struct {
			StartTime      string `json:"start_time,omitempty"`
			ExpireTime     string `json:"expire_time,omitempty"`
			MaxSendNum     int64  `json:"max_send_num,omitempty"`
			OutBizNo       string `json:"out_biz_no,omitempty"`
			ApnsCollapseId string `json:"apns_collapse_id,omitempty"`
		}{},
		ProductionMode: "false",
		Description:    "测试推送",
	}
	appKey := "60ed3a84a6f90557b7b591cf"
	appMasterKey := "rf7n3bkq4hz8v8u9cyuate3r9w9sa9zm"
	umengPush := umeng.NewUmengPush(appKey, appMasterKey)
	result, err := umengPush.Send(param)
	fmt.Println(result.Ret)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(result.Data.MsgId, result.Data.ErrorMsg)
}

