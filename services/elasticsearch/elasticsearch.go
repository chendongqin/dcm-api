package elasticsearch

import (
	"crypto/tls"
	"dongchamao/global"
	"dongchamao/global/utils"
	"github.com/elastic/go-elasticsearch/v8"
	"net"
	"net/http"
	"time"
)

type ElasticsearchService struct {
	queryString string
}

type ElasticResp struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     ElasticHits
}

type ElasticCountResp struct {
	Count int64 `json:"count"`
}

type ElasticHits struct {
	Total int             `json:"total"`
	Hits  []ElasticSource `json:"hits"`
}

type ElasticSource struct {
	Index  interface{}            `json:"_index"`
	Type   interface{}            `json:"_type"`
	Id     interface{}            `json:"_id"`
	Source map[string]interface{} `json:"_source"`
}

type MapObject map[string]interface{}

func NewElasticMultiQuery() *ElasticMultiQuery {
	return new(ElasticMultiQuery)
}

func NewElasticQueryGroup() (*ElasticQuery, *ElasticMultiQuery) {
	return NewElasticQuery(), NewElasticMultiQuery()
}

func NewElasticQuery() *ElasticQuery {
	return new(ElasticQuery)
}

func Query() *ElasticQuery {
	return NewElasticQuery()
}

func GetBuckets(resp map[string]interface{}, name string) []interface{} {
	if aggregations, ok := resp["aggregations"].(map[string]interface{}); ok {
		data := aggregations[name].(map[string]interface{})
		buckets := utils.ToInterfaceSlice(data["buckets"])
		return buckets
	} else if data, ok := resp[name].(map[string]interface{}); ok {
		return data["buckets"].([]interface{})
	}
	return nil
}

func NewWritableClient() (client *elasticsearch.Client, err error) {
	host := global.Cfg.String("es_host")
	user := global.Cfg.String("es_write_user")
	passwd := global.Cfg.String("es_write_passwd")
	cfg := elasticsearch.Config{
		Addresses: []string{
			host,
		},
		Username: user,
		Password: passwd,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 10,
			DialContext:           (&net.Dialer{Timeout: time.Second * 5}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
	}
	client, err = elasticsearch.NewClient(cfg)
	return
}
