package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
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

//func Consumer() error {
//	kafkaHotsConf := global.Cfg.String("kafka_hosts")
//	if kafkaHotsConf == "" {
//		logs.Error("kafka fail :( kafka_hosts is empty")
//		return errors.New("kafka_hosts is empty")
//	}
//	c, err := kafka.NewConsumer(&kafka.ConfigMap{
//		"bootstrap.servers":               kafkaHotsConf,
//		"group.id":                        "",
//		"socket.timeout.ms":               10000,
//		"session.timeout.ms":              10000,
//		"broker.address.family":           "v4",
//		"go.events.channel.enable":        true,
//		"go.application.rebalance.enable": true,
//	})
//	if err  != nil {
//		return err
//	}
//	defer c.Close()
//	err = c.SubscribeTopics([]string{"dy-live-chat-message"}, nil)
//}
