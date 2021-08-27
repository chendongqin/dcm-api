package global

import (
	"crypto/tls"
	"dongchamao/global/cache"
	_ "dongchamao/global/cache/redis"
	"dongchamao/global/mysql"
	aliLog "dongchamao/services/ali_log"
	"dongchamao/services/elastichelper"
	"dongchamao/services/kafka"
	"dongchamao/services/pools"
	"fmt"
	"github.com/astaxie/beego"
	beegoCache "github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/silenceper/wechat/v2"
	wxCache "github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"gopkg.in/mgo.v2"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var Cache cache.CacheInterface
var MemoryCache beegoCache.Cache

//var FileCache cache.CacheInterface
var Cfg = beego.AppConfig

//var CurrentDir string
var MongoSession *mgo.Session

//const (
//	QueueLog                        = "CMM:QUENE:LOG:COMMON"
//	QueueUpdateStarAuthor           = "CMM:QUENE:UPDATE:STAR:AUTHOR"
//	QueueDouplusDeliveryOrders      = "CMM:QUEUE:DOUPLUS:DELIVERY:ORDERS"
//	QueueDouplusDeliveryDelayOrders = "CMM:QUEUE:DOUPLUS:DELIVERY:DELAY:ORDERS"
//	QueueDouplusOrdersExchange      = "delay:douplus:orders:exchange"
//)

func InitEnv() {
	_initBConfig()
	_initLogs()
	_initCache()
	_initWxOfficialAccount()
	_initEs()
	_initHbaseThriftPool()
	_initDataBase()
	_initSlsConfig()
	//_initMongodb() // deprecated
	//_initValidate()
	//_initRabbitMqPool()
	//_initKafkaProducer()
	//初始化全局httpclient超时时间
	http.DefaultClient.Timeout = 30 * time.Second
}

func _initKafkaProducer() {
	kafkaHostsConf := Cfg.String("kafka_hosts")
	if kafkaHostsConf == "" {
		logs.Error("kafka init fail :( kafka_hosts is empty")
		os.Exit(1)
	}
	kafkaHosts := strings.Split(kafkaHostsConf, ",")
	for k, v := range kafkaHosts {
		v = strings.Trim(v, " ")
		kafkaHosts[k] = v
	}
	kafka.Init(kafkaHosts)
}

func _initBConfig() {
	//调整panic recover方法
	beego.BConfig.RecoverFunc = RequestRecoverPanic
}

func _initLogs() {
	err0 := logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"%s","level":6,"maxlines":2000000,"maxsize":0,"daily":true,"maxdays":7,"color":true}`, Cfg.String("logfile")))
	fmt.Println(err0)
	logs.Async()
}

func _initCache() {
	var err error
	default_redis_host := Cfg.String("default_redis_host")
	default_redis_passwd := Cfg.String("default_redis_passwd")
	default_redis_db := Cfg.String("default_redis_db")
	default_redis_maxidle := Cfg.String("default_redis_maxidle")

	logs.Info("cache init start")
	Cache, err = cache.CacheFactory("redis", `{"conn":"`+default_redis_host+`","dbNum":"`+default_redis_db+`","password":"`+default_redis_passwd+`","maxIdle":"`+default_redis_maxidle+`"}`)
	if err != nil {
		logs.Error("cache init fail :(", err)
		os.Exit(1)
	}

	//内存缓存，一部分频繁使用的数据(如IP白名单)使用一层本地内存缓存，避免REDIS的IO太过频繁
	MemoryCache, err = beegoCache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		logs.Error("memory cache init fail :(", err)
		os.Exit(1)
	}
}

func _initDataBase() {
	dbLink := Cfg.String("master_db")
	maxIdleConn, _ := Cfg.Int("master_db_max_idle_conn")
	maxOpenConn, _ := Cfg.Int("master_db_max_open_conn")
	slaveDbLink := Cfg.String("slave_db")
	slaveMaxidleconn, _ := Cfg.Int("slave_db_max_idle_conn")
	slaveMaxopenconn, _ := Cfg.Int("slave_db_max_open_conn")
	master := mysql.Options{
		Dns:         dbLink,
		MaxIdleConn: maxIdleConn,
		MaxOpenConn: maxOpenConn,
	}
	slaves := make([]mysql.Options, 0)
	if slaveDbLink != "" {
		slaves = append(slaves, mysql.Options{
			Dns:         slaveDbLink,
			MaxIdleConn: slaveMaxidleconn,
			MaxOpenConn: slaveMaxopenconn,
		})
	}
	err := mysql.InitMysql("default", master, slaves, IsDev())
	if err != nil {
		fmt.Println("db1 init fail :(")
		os.Exit(1)
	}

}

// deprecated
func _initMongodb() {
	var err error
	mongodb_addr := Cfg.String("default_mongodb_addr")
	mongodb_database := Cfg.String("default_mongodb_database")
	mongodb_user := Cfg.String("default_mongodb_user")
	mongodb_passwd := Cfg.String("default_mongodb_passwd")
	mongodb_maxidle, _ := Cfg.Int("default_mongodb_maxidle")

	//init mongodb
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:     []string{mongodb_addr},
		Timeout:   60 * time.Second,
		Database:  mongodb_database,
		Username:  mongodb_user,
		Password:  mongodb_passwd,
		Mechanism: "SCRAM-SHA-1",
	}

	MongoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		fmt.Println(err)
		fmt.Println("mongodb init fail :(")
		os.Exit(1)
	}
	MongoSession.SetPoolLimit(mongodb_maxidle)
	MongoSession.SetMode(mgo.Monotonic, true)
}

var esClients map[string]*elasticsearch.Client

func GetEsClient(clientTypes ...string) *elasticsearch.Client {
	name := "default"
	if len(clientTypes) > 0 {
		name = clientTypes[0]
		if name == "" {
			name = "default"
		}
	}
	return esClients[name]
}

func _initEs() {
	esClients = make(map[string]*elasticsearch.Client, 2)
	_initEsWithName()
	//_initEsWithName("fast")
}

func _initEsWithName(clientTypes ...string) {
	clientType := ""
	name := "default"
	if len(clientTypes) > 0 {
		clientType = "_" + clientTypes[0]
		name = clientTypes[0]
	}
	es_host := Cfg.String("es_host" + clientType)
	es_user := Cfg.String("es_user" + clientType)
	es_passwd := Cfg.String("es_passwd" + clientType)
	cfg := elasticsearch.Config{
		Addresses: []string{
			es_host,
		},
		Username: es_user,
		Password: es_passwd,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			MaxConnsPerHost:       10,
			ResponseHeaderTimeout: time.Second * 10,
			DialContext:           (&net.Dialer{Timeout: time.Second * 5}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logs.Error("[init] [elasticsearch] err: %s", err)
		return
	}
	esClients[name] = client
	logs.Info("register elasticsearch: [%s] %s", name, es_host)
}

//使用olivere/elastic
func _initOlivereElastic() {
	name := "default"
	host := Cfg.String("es_host")
	user := Cfg.String("es_user")
	password := Cfg.String("es_passwd")

	elastichelper.InitElasticSearchClient(name, host, user, password)
}

var HbasePools *pools.HBasePoolFactory

func _initHbaseThriftPool() {
	HbasePools = pools.GetHBaseFactory()
	user := Cfg.String("hbase_user")
	password := Cfg.String("hbase_passwd")
	host := Cfg.String("hbase_host")
	option := pools.ThriftHbaseConfig{
		Host:        host,
		User:        user,
		Password:    password,
		MinConn:     10,
		MaxConn:     40,
		MaxIdle:     40,
		IdleTimeout: 30,
	}
	pool, err := pools.NewThriftHbasePools(&option)
	if err != nil {
		fmt.Println(err)
		panic("hbase init fail :(" + err.Error())
	}
	HbasePools.Add("default", pool)
}

func _initSlsConfig() {
	cf := &aliLog.SlsConfig{
		Endpoint:        "cn-hangzhou.log.aliyuncs.com",
		AccessKeyID:     Cfg.String("ali_accessKey"),
		AccessKeySecret: Cfg.String("ali_secret"),
		DefaultProject:  "dongchamao-api-log",
	}
	if beego.BConfig.RunMode == "dev" {
		cf.Endpoint = "cn-hangzhou.log.aliyuncs.com"
	}
	aliLog.InitAliyunSls(cf)
}

func _initValidate() {
	// 设置表单验证messages
	var MessageTmpls = map[string]string{
		"Required":     "不能为空",
		"Min":          "最小值 为 %d",
		"Max":          "最大值 为 %d",
		"Range":        "范围 为 %d 到 %d",
		"MinSize":      "最短长度 为 %d",
		"MaxSize":      "最大长度 为 %d",
		"Length":       "长度必须 为 %d",
		"Alpha":        "必须是有效的字母",
		"Numeric":      "必须是有效的数字",
		"AlphaNumeric": "必须是有效的字母或数字",
		"Match":        "必须匹配 %s",
		"NoMatch":      "必须不匹配 %s",
		"AlphaDash":    "必须是有效的字母、数字或连接符号(-_)",
		"Email":        "必须是有效的电子邮件地址",
		"IP":           "必须是有效的IP地址",
		"Base64":       "必须是有效的base64字符",
		"Mobile":       "必须是有效的手机号码",
		"Tel":          "必须是有效的电话号码",
		"Phone":        "必须是有效的电话或移动电话号码",
		"ZipCode":      "必须是有效的邮政编码",
	}

	validation.SetDefaultMessage(MessageTmpls)
}

//初始化微信公众号

var WxOfficial *officialaccount.OfficialAccount

func _initWxOfficialAccount() {
	wc := wechat.NewWechat()
	wechatAppId := Cfg.String("wx_office_app_id")
	wechatAppSecret := Cfg.String("wx_office_app_secret")
	wechatToken := Cfg.String("wx_office_app_token")
	wechatEncodedAESKey := Cfg.String("wx_office_encoded_aes_key")

	db, _ := strconv.Atoi(Cfg.String("default_redis_db"))
	wxRedis := &wxCache.RedisOpts{
		Host:     Cfg.String("default_redis_host"),
		Password: Cfg.String("default_redis_passwd"),
		Database: db,
	}
	cfg := &config.Config{
		AppID:          wechatAppId,
		AppSecret:      wechatAppSecret,
		Token:          wechatToken,
		EncodingAESKey: wechatEncodedAESKey,
		Cache:          wxCache.NewRedis(wxRedis),
	}
	WxOfficial = wc.GetOfficialAccount(cfg)
}

func IsDev() bool {
	if beego.BConfig.RunMode != "product" {
		return true
	}
	return false
}
