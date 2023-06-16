package mysql_dcm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	"sync"
	"time"
	"xorm.io/core"
)

type Options struct {
	Dns         string
	MaxIdleConn int
	MaxOpenConn int
}

var sessionList = sync.Map{}

func InitMysql(dbName string, master Options, Slave []Options, isDev bool) error {
	fmt.Printf("Load db %s config...", dbName)
	masterEngine, err := initDbHandle(master, isDev)
	if err != nil {
		return err
	}
	if len(Slave) == 0 {
		eg, err := xorm.NewEngineGroup(masterEngine, []*xorm.Engine{masterEngine})
		if err == nil {
			sessionList.Store(dbName, eg)
		} else {
			return err
		}
	} else {
		fmt.Printf("Load db %s slave config...", dbName)
		engines := make([]*xorm.Engine, 0)
		for _, v := range Slave {
			engine, err := initDbHandle(v, isDev)
			if err == nil {
				engines = append(engines, engine)
			}
		}
		eg, err := xorm.NewEngineGroup(masterEngine, engines)
		if err == nil {
			sessionList.Store(dbName, eg)
		} else {
			return err
		}
	}
	return nil
}

func initDbHandle(options Options, isDev bool) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", options.Dns)
	if err != nil {
		return nil, err
	}
	engine.SetMaxIdleConns(options.MaxIdleConn)
	engine.SetMaxOpenConns(options.MaxOpenConn)
	engine.SetConnMaxLifetime(8 * time.Second)

	// 本地环境开启日志
	if isDev {
		engine.ShowExecTime(true)
		engine.ShowSQL(true)
		engine.Logger().SetLevel(core.LOG_DEBUG)
	}
	err = engine.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("MysqlClient ping [%s], %s\n", options.Dns, err.Error()))
	}
	return engine, nil
}

func GetSession(name string) *xorm.EngineGroup {
	if v, ok := sessionList.Load(name); ok {
		return v.(*xorm.EngineGroup)
	}
	return nil
}
