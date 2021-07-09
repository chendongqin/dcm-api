package dcm

import (
	"dongchamao/global/mysql"
	"github.com/go-xorm/xorm"
)

func getDb() *xorm.Engine {
	return mysql.GetSession("default").Master()
}

func getSlaveDb() *xorm.Engine {
	return mysql.GetSession("default").Slave()
}

func GetDbSession() *xorm.Session {
	return getDb().NewSession()
}

func GetSlaveDbSession() *xorm.Session {
	return getSlaveDb().NewSession()
}

func Insert(session *xorm.Session, beans interface{}) (int64, error) {
	if session == nil {
		return getDb().Insert(beans)
	}
	return session.Insert(beans)
}

func Get(id int, bean interface{}) (bool, error) {
	return getSlaveDb().Id(id).Get(bean)
}

func GetBy(queryParam string, queryValue interface{}, bean interface{}) (bool, error) {
	return getSlaveDb().Where(queryParam+"=?", queryValue).Desc("id").Get(bean)
}

func UpdateInfo(session *xorm.Session, id int, params map[string]interface{}, tableNameOrBean interface{}) (int64, error) {
	if session == nil {
		return getDb().Table(tableNameOrBean).ID(id).Update(params)
	} else {
		return session.Table(tableNameOrBean).ID(id).Update(params)
	}
}

type CommonMap map[string]interface{}

func ParsePageAndPageSize(params CommonMap) (rePage, rePageSize, limitOffSet int) {
	rePage = params.GetInt("Page")
	rePageSize = params.GetInt("PageSize")
	if rePage <= 0 {
		rePage = 1
	}
	if rePageSize <= 0 {
		rePageSize = 100000
	}
	limitOffSet = (rePage - 1) * rePageSize
	return
}

func (cm *CommonMap) GetInt(key string) int {
	if v, ok := (*cm)[key]; ok {
		return v.(int)
	}
	return 0
}

func (cm *CommonMap) GetInt64(key string) int64 {
	if v, ok := (*cm)[key]; ok {
		return v.(int64)
	}
	return 0
}

func (cm *CommonMap) GetString(key string) string {
	if v, ok := (*cm)[key]; ok {
		return v.(string)
	}
	return ""
}
