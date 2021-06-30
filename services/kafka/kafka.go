package kafka

import (
	"dongchamao/global/alias"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"time"
)

var producer sarama.AsyncProducer

func Init(hosts []string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.ClientID = "dongchamao"
	config.Net.DialTimeout = 1 * time.Second
	config.Producer.Timeout = 10 * time.Second
	config.Version = sarama.V2_0_0_0
	var err error
	//sarama.Logger = logs.GetLogger()
	producer, err = sarama.NewAsyncProducer(hosts, config)
	if err != nil {
		logs.Error("[kafka] ", err)
		panic(err)
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logs.Critical(fmt.Sprintf("err: %s", err))
			}
		}()
		for {
			select {
			case success := <-producer.Successes():
				logs.Info("[kafka] producer push success, topic: %s, offset: %d,  timestamp: %s", success.Topic, success.Offset, success.Timestamp.String())
			case errs := <-producer.Errors():
				logs.Error("[kafka] producer err: %s\n", errs.Err.Error())
			}
		}
	}()
}

func NewTalentMsg(authorId string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{}
	msg.Topic = `dy-talent`
	msg.Value = sarama.StringEncoder("{\"author_id\": \"" + authorId + "\"}")
	return msg
}

func SendMessage(message *sarama.ProducerMessage) {
	defer func() {
		if err := recover(); err != nil {
			logs.Critical(fmt.Sprintf("err: %s", err))
		}
	}()
	producer.Input() <- message
}

func SendTalent(authorId string) {
	SendMessage(NewTalentMsg(authorId))
}

func NewLiveFansFeatureMsg(authorId, roomId string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic: "dy-api-message",
		Value: sarama.StringEncoder(pack(alias.M{
			"author_id": authorId,
			"room_id":   roomId,
			"type":      "liveImage",
			"time":      time.Now().Format("2006-01-02 15:04:05"),
		})),
	}
	return msg
}

func NewLiveWordCloudMsg(authorId, roomId string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic: "dy-api-message",
		Value: sarama.StringEncoder(pack(alias.M{
			"author_id": authorId,
			"room_id":   roomId,
			"type":      "liveWordCloud",
			"time":      time.Now().Format("2006-01-02 15:04:05"),
		})),
	}
	return msg
}

func NewAuthorFansFeatureMsg(authorId string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic: "dy-api-message",
		Value: sarama.StringEncoder(pack(alias.M{
			"author_id": authorId,
			"type":      "authorImage",
		})),
	}
	return msg
}

func SendLiveFansFeature(authorId string, roomId string) {
	SendMessage(NewLiveFansFeatureMsg(authorId, roomId))
}

func SendLiveWordCloud(authorId string, roomId string) {
	SendMessage(NewLiveWordCloudMsg(authorId, roomId))
}

func SendAuthorFansFeature(authorId string) {
	SendMessage(NewAuthorFansFeatureMsg(authorId))
}

func NewHighLevelIntentionImageMsg(roomId string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic: "dy-api-message-high-level",
		Value: sarama.StringEncoder(pack(alias.M{
			"room_id": roomId,
			"type":    "intentionImage",
		})),
	}
	return msg
}

func NewHighLevelFansPurchaseMsg(roomId string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic: "dy-api-message-high-level",
		Value: sarama.StringEncoder(pack(alias.M{
			"room_id": roomId,
			"type":    "fansPurchase",
		})),
	}
	return msg
}

func NewHighLevelFansFlowMsg(roomIds []string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic: "dy-api-message-high-level",
		Value: sarama.StringEncoder(pack(alias.M{
			"room_id_list": roomIds,
			"type":         "fansFlow",
		})),
	}
	return msg
}

func SendHighLevelIntentionImage(roomId string) {
	SendMessage(NewHighLevelIntentionImageMsg(roomId))
}

func SendHighLevelFansPurchaseImage(roomId string) {
	SendMessage(NewHighLevelFansPurchaseMsg(roomId))
}

func SendHighLevelFansFlowImage(roomIds []string) {
	SendMessage(NewHighLevelFansFlowMsg(roomIds))
}

func pack(m alias.M) string {
	str, _ := jsoniter.MarshalToString(m)
	return str
}

//我的抖音号数据 推送kafka
func SendAuthorMineData(dataType string, data string) {

	if dataType == "DyCreatorApiAuthorInfo" || dataType == "DyCreatorApiFansPortrait" {
		dataMap := make(map[string]interface{}, 0)
		_ = jsoniter.Unmarshal([]byte(data), &dataMap)
		msg := &sarama.ProducerMessage{
			Topic: "my-douyin-data",
			Value: sarama.StringEncoder(pack(alias.M{
				"data": dataMap,
				"type": dataType,
				"time": time.Now().Format("2006-01-02 15:04:05"),
			})),
		}
		SendMessage(msg)
	}

	if dataType == "DyCreatorApiAwemeList" || dataType == "DyCreatorApiLiveRoomList" {
		dataList := make([]interface{}, 0)
		_ = jsoniter.Unmarshal([]byte(data), &dataList)
		msg := &sarama.ProducerMessage{
			Topic: "my-douyin-data",
			Value: sarama.StringEncoder(pack(alias.M{
				"data": dataList,
				"type": dataType,
				"time": time.Now().Format("2006-01-02 15:04:05"),
			})),
		}
		SendMessage(msg)
	}

}
