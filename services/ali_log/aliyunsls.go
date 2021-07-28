package aliLog

import (
	"encoding/json"
	"errors"
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/gogo/protobuf/proto"
	"time"
)

var ProducerInstance *producer.Producer
var DefaultProject string

type SlsConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	DefaultProject  string
}

func InitAliyunSls(sc *SlsConfig) {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = sc.Endpoint
	producerConfig.AccessKeyID = sc.AccessKeyID
	producerConfig.AccessKeySecret = sc.AccessKeySecret
	ProducerInstance = producer.InitProducer(producerConfig)
	ProducerInstance.Start() // 启动producer实例
	DefaultProject = sc.DefaultProject
}

//通用日志结构
type CommonAliyunSlsLog struct {
	Project  string            `json:"project"`
	Logstore string            `json:"logstore"`
	Logtime  int64             `json:"logtime"`
	Logmsg   map[string]string `json:"logmsg"`
}

func PushLog(casl CommonAliyunSlsLog, levels ...string) error {
	if ProducerInstance == nil {
		return errors.New("producer not exists")
	}
	if casl.Logstore == "" {
		return errors.New("logstore not exists")
	}
	var project = DefaultProject
	if casl.Project != "" {
		project = casl.Project
	}
	if project == "" {
		return errors.New("project not exists")
	}
	content := []*sls.LogContent{}
	for k, v := range casl.Logmsg {
		content = append(content, &sls.LogContent{
			Key:   proto.String(k),
			Value: proto.String(v),
		})
	}
	if len(levels) > 0 {
		logLevel := levels[0]
		content = append(content, &sls.LogContent{
			Key:   proto.String("log_level"),
			Value: proto.String(logLevel),
		})
	}
	pushlogtime := uint32(casl.Logtime)
	log := &sls.Log{
		Time:     proto.Uint32(pushlogtime),
		Contents: content,
	}
	err := ProducerInstance.SendLog(project, casl.Logstore, "topic", "127.0.0.1", log)
	//cback := new(aliyunSlsCallback)
	//err := ProducerInstance.SendLogWithCallBack(project, cl.Logstore, "topic", "127.0.0.1", log, cback)
	return err

}

type aliyunSlsCallback struct {
}

func (c *aliyunSlsCallback) Success(r *producer.Result) {
	fmt.Println(r.GetErrorCode(), r.GetErrorMessage())
}

func (c *aliyunSlsCallback) Fail(r *producer.Result) {
	fmt.Println(r.GetErrorCode(), r.GetErrorMessage())
}

func Info(msg string) {
	var cl CommonAliyunSlsLog
	if err := json.Unmarshal([]byte(msg), &cl); err == nil {
		_ = PushLog(cl, "INFO")
	}
}

func Warn(msg string) {
	var cl CommonAliyunSlsLog
	if err := json.Unmarshal([]byte(msg), &cl); err == nil {
		_ = PushLog(cl, "WARN")
	}
}

func Error(msg string) {
	var cl CommonAliyunSlsLog
	if err := json.Unmarshal([]byte(msg), &cl); err == nil {
		_ = PushLog(cl, "ERROR")
	}
}

func Log(logstore string, msg map[string]string) {
	newInfo := CommonAliyunSlsLog{
		Logstore: logstore,
		Logmsg:   msg,
		Logtime:  time.Now().Unix(),
	}
	_ = PushLog(newInfo)
}
