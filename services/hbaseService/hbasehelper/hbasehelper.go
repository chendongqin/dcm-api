package hbasehelper

import (
	"douyin-api/entity"
	"douyin-api/services/hbaseService/hbase"
)

func NewQuery() *Query {
	return &Query{family: entity.TableFamily, filter: hbase.NewFilters()}
}
