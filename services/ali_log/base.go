package aliLog

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"encoding/json"
	"fmt"
	_ "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
)

const (
	MQ_CMM_LOG_COMMON string = "DCM:QUENE:LOG:COMMON"
	Endpoint          string = "cn-shanghai.log.aliyuncs.com"
)

// accesskey配置
var (
	AccessKeyID     = global.Cfg.String("ali_accessKey")
	AccessKeySecret = global.Cfg.String("ali_secret")
)

//通用日志结构
type CommonAliLog struct {
	Logstore string            `json:"logstore"`
	Logtime  int64             `json:"logtime"`
	Logmsg   map[string]string `json:"logmsg"`
}

func EsLog(table string, queryString string, file string, line int, spendTime float64) {
	insertData := map[string]string{
		"index":        table,
		"query":        queryString,
		"caller_file":  file,
		"caller_line":  strconv.Itoa(line),
		"request_time": fmt.Sprintf("%.4f", spendTime),
	}
	go ReportAli(MQ_CMM_LOG_COMMON, "dongchamao-es-log", insertData)
}

//通用日志，最多支持5个额外参数,格式map[string]interface{}
func CommonLog(logType string, logAction string, logContent string, logExtraContents ...map[string]interface{}) {
	insertData := make(map[string]string)
	insertData["log_type"] = logType
	insertData["log_action"] = logAction
	insertData["log_content"] = logContent
	if len(logExtraContents) > 0 {
		logExtraContent := logExtraContents[0]
		extraData := make(map[string]interface{}, 0)
		count := 0
		for k, v := range logExtraContent {
			switch count {
			case 0:
				extraData["a_f"] = k
				extraData["a_c"] = v
			case 1:
				extraData["b_f"] = k
				extraData["b_c"] = v
			case 2:
				extraData["c_f"] = k
				extraData["c_c"] = v
			case 3:
				extraData["d_f"] = k
				extraData["d_c"] = v
			case 4:
				extraData["e_f"] = k
				extraData["e_c"] = v
			}
			count++
		}
		if len(extraData) > 0 {
			logExtra, err := json.Marshal(extraData)
			if err == nil {
				insertData["log_extra_content"] = string(logExtra)
			}
		}
	}
	go ReportAli(MQ_CMM_LOG_COMMON, "dongchamao-log-common", insertData)
	return
}

func ReportAli(queueName, logstore string, msg map[string]string) error {
	if _, ok := msg["env"]; !ok {
		msg["env"] = beego.BConfig.RunMode
	}
	//newInfo := CommonAliLog{
	//	Logstore: logstore,
	//	Logmsg:   msg,
	//	Logtime:  time.Now().Unix(),
	//}
	Log(logstore, msg)
	//reMsg, err := jsoniter.Marshal(newInfo)
	//if err != nil {
	//	return err
	//}
	//err = global.MqChannelPool.Publish("", queueName, false, false, false, amqp.Publishing{
	//	ContentType: "text/plain",
	//	Body:        reMsg,
	//})
	return nil
}

//定制log
//输入输出日志
func LogInput(requestId string, clientId string, logType string, appId int, uid int, url string, method string, ip string, useragent string, refer string, apidatas interface{}, userGroupId int, remote_addr string) {
	var realUrl string
	urldata := strings.Split(url, "?")
	if len(urldata) > 0 {
		realUrl = urldata[0]
	}
	insertData := make(map[string]string)
	insertData["request_id"] = requestId
	insertData["client_id"] = clientId
	insertData["log_type"] = logType
	insertData["appid"] = utils.ToString(appId)
	insertData["uid"] = utils.ToString(uid)
	insertData["group_id"] = utils.ToString(userGroupId)
	insertData["method"] = method
	insertData["ip"] = ip
	insertData["remote_addr"] = remote_addr
	insertData["useragent"] = useragent
	insertData["refer"] = refer
	insertData["url"] = realUrl
	insertData["args"] = ""
	insertData["timestamp"] = utils.ToString(time.Now().Unix())
	args, err := json.Marshal(apidatas)
	if err == nil {
		insertData["args"] = string(args)
	}
	go ReportAli(MQ_CMM_LOG_COMMON, "dongchamao-log-api-history", insertData)
}

//记录登录日志
func LoginLog(requestId string, clientId string, appId int, uid int64, token string, ip string, useragent string, grant_type string, action string) {
	if action == "" {
		action = "login"
	}
	insertData := make(map[string]string)
	insertData["request_id"] = requestId
	insertData["client_id"] = clientId
	insertData["appid"] = utils.ToString(appId)
	insertData["uid"] = utils.ToString(uid)
	insertData["token"] = token
	insertData["grant_type"] = grant_type
	insertData["login_ip"] = ip
	insertData["action"] = action
	insertData["useragent"] = useragent
	go ReportAli(MQ_CMM_LOG_COMMON, "dongchamao-log-user-active", insertData)
}

func SendCode(appId int, uid int64, phone string, ip string, action string) {
	insertData := make(map[string]interface{})
	insertData["appid"] = appId
	insertData["uid"] = uid
	insertData["send_ip"] = ip
	CommonLog("sendcode", action, phone, insertData)
}
