package pools

import (
	"dongchamao/services/hbaseService/hbase"
	"errors"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/astaxie/beego/logs"
	"github.com/silenceper/pool"
	"net/http"
	"sync"
	"time"
)

const (
	HOST = "http://ld-uf6w1y03mb950v52e-proxy-hbaseue-pub.hbaseue.rds.aliyuncs.com:9190"
	// 用户名
	USER = "root"
	// 密码
	PASSWORD = "root"
)

type ThriftHbaseConfig struct {
	Host        string
	User        string
	Password    string
	MinConn     int
	MaxConn     int
	MaxIdle     int
	IdleTimeout time.Duration
}

func init() {
	hbasePoolFactory = new(HBasePoolFactory)
}

type HBasePoolFactory struct {
	pools map[string]*ThriftHbasePools
	mutex sync.Mutex
}

var hbasePoolFactory *HBasePoolFactory

func GetHBaseFactory() *HBasePoolFactory {
	return hbasePoolFactory
}

func (factory *HBasePoolFactory) Add(name string, pool *ThriftHbasePools) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()
	if factory.pools == nil {
		factory.pools = make(map[string]*ThriftHbasePools)
	}
	factory.pools[name] = pool
}

func (factory *HBasePoolFactory) Get(names ...string) *ThriftHbasePoolsClient {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}
	if name == "" {
		name = "default"
	}
	hbasePool := factory.Pool(name)
	if hbasePool == nil {
		hbasePool = factory.Pool("default")
		logs.Warn("[hbase] 未找到 [%s] 对应的 hbase 连接池，使用 [default] 继续提供服务", name)
		if hbasePool == nil {
			return nil
		}
	}
	return hbasePool.Get()
}

func (factory *HBasePoolFactory) Pool(name string) *ThriftHbasePools {
	hbasePool, ok := factory.pools[name]
	if !ok {
		return nil
	}
	return hbasePool
}

func NewThriftHbasePools(option *ThriftHbaseConfig) (*ThriftHbasePools, error) {
	if option.Host == "" {
		return nil, errors.New("need host")
	}
	//默认值
	if option.MinConn == 0 {
		option.MinConn = 5
	}
	//创建的最大连接数
	if option.MaxConn == 0 {
		option.MaxConn = option.MinConn * 4
	}
	//池子里的最大连接数
	if option.MaxIdle == 0 {
		option.MaxIdle = option.MinConn * 4
	}

	if option.IdleTimeout == 0 {
		option.IdleTimeout = 15
	}

	factory := func() (interface{}, error) {
		trans, err := thrift.NewTHttpClientWithOptions(option.Host, thrift.THttpClientOptions{Client: &http.Client{
			Timeout: time.Second * 10,
		}})
		if err != nil {
			return nil, err
		}
		// 设置用户名密码
		httClient := trans.(*thrift.THttpClient)
		if option.User != "" {
			httClient.SetHeader("ACCESSKEYID", option.User)
		}
		if option.Password != "" {
			httClient.SetHeader("ACCESSSIGNATURE", option.Password)
		}
		return httClient, nil
	}
	close := func(v interface{}) error {
		return v.(*thrift.THttpClient).Close()
	}

	//创建一个连接池： 初始化5，最大连接30
	poolConfig := &pool.Config{
		InitialCap: option.MinConn,
		MaxCap:     option.MaxConn,
		//MaxIdle: 	option.MaxIdle,
		Factory:    factory,
		Close:      close,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: (option.IdleTimeout) * time.Second,
	}
	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		return nil, err
	}
	return &ThriftHbasePools{
		tranPool: p,
	}, nil
}

type ThriftHbasePools struct {
	tranPool pool.Pool
}

type ThriftHbasePoolsClient struct {
	*hbase.THBaseServiceClient
	trans    *thrift.THttpClient
	tranPool *pool.Pool
}

func (this *ThriftHbasePools) Get() *ThriftHbasePoolsClient {
	//current := this.tranPool.Len()
	//logs.Debug("hbase pool len: ", current)
	trans, err := this.tranPool.Get()
	if err != nil {
		return nil
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := hbase.NewTHBaseServiceClientFactory(trans.(*thrift.THttpClient), protocolFactory)
	return &ThriftHbasePoolsClient{
		client,
		trans.(*thrift.THttpClient),
		&this.tranPool,
	}
}

func (c *ThriftHbasePoolsClient) Close() error {
	return (*c.tranPool).Put(c.trans)
}
