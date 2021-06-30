package elastichelper

import (
	"crypto/tls"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	storage = make(map[string]*elastic.Client)
)

//InitElasticSearchClient 根据配置初始化esclient
func InitElasticSearchClient(name, host, user, password string) {
	//name := "default"
	//host := global.Cfg.String("es_host")
	//user := global.Cfg.String("es_user")
	//password := global.Cfg.String("es_passwd")
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 10,
			DialContext:           (&net.Dialer{Timeout: time.Second * 5}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
		Timeout: 5 * time.Second,
	}

	options := []elastic.ClientOptionFunc{
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetBasicAuth(user, password),
		elastic.SetHttpClient(httpClient),
	}
	//如果是本地环境，打印elasticsearch日志
	if beego.BConfig.RunMode == "dev" {
		options = append(
			options,
			elastic.SetErrorLog(log.New(os.Stderr, "Err:ELASTIC ", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
			elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)),
		)
	}

	requestClient, err := elastic.NewClient(options...)
	if err != nil {
		// Handle error
		logs.Critical("elastic node [" + name + "] new client failed: " + err.Error())
		//panic("elastic node [" + node + "] new client failed: " + err.Error())
	}
	storage[name] = requestClient
}

//GetClient 获取esclient
func GetClient(args ...string) *elastic.Client {
	name := "default"
	if len(args) > 0 {
		name = args[0]
	}
	client, exists := storage[name]
	if exists {
		return client
	}
	return nil
}
