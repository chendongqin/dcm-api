package hbasehelper

import (
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbase"
)

func NewQuery() *Query {
	return &Query{family: hbaseService.TableFamily, filter: hbase.NewFilters()}
}
