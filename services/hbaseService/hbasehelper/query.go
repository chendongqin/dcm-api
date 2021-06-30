package hbasehelper

import (
	"context"
	"douyin-api/entity"
	"douyin-api/global"
	"douyin-api/services/hbaseService/hbase"
	"douyin-api/services/pools"
	"github.com/astaxie/beego/logs"
	"time"
)

type Query struct {
	columns     []*hbase.TColumn
	tableName   string
	family      string
	startRow    []byte
	stopRow     []byte
	batchSize   *int32
	cacheBlocks *bool
	caching     *int32
	reversed    *bool
	filter      *hbase.Filter
}

func (q *Query) SetCaching(caching int32) *Query {
	if q.caching == nil {
		q.caching = new(int32)
	}
	*q.caching = caching
	return q
}

func (q *Query) SetBatchSize(batchSize int32) *Query {
	if q.batchSize == nil {
		q.batchSize = new(int32)
	}
	*q.batchSize = batchSize
	return q
}

func (q *Query) SetCacheBlocks(cacheBlocks bool) *Query {
	if q.cacheBlocks == nil {
		q.cacheBlocks = new(bool)
	}
	*q.cacheBlocks = cacheBlocks
	return q
}

func (q *Query) SetReversed(reversed bool) *Query {
	if q.reversed == nil {
		q.reversed = new(bool)
	}
	*q.reversed = reversed
	return q
}

func (q *Query) SetTable(tableName string) *Query {
	q.tableName = tableName
	return q
}

func (q *Query) SetFamily(family string) *Query {
	q.family = family
	return q
}

func (q *Query) Filter() *hbase.Filter {
	return q.filter
}

func (q *Query) AddFilter(filterInterface hbase.FilterInterface) *Query {
	q.filter.Add(filterInterface)
	return q
}

func (q *Query) Select(columns ...string) *Query {
	if q.columns == nil {
		q.columns = make([]*hbase.TColumn, len(columns))
	}
	for key, column := range columns {
		tColumn := &hbase.TColumn{
			Family:    []byte(q.family),
			Qualifier: []byte(column),
		}
		q.columns[key] = tColumn
	}
	return q
}

func (q *Query) getClient() *pools.ThriftHbasePoolsClient {
	client := global.HbasePools.Get(entity.Mapping(q.tableName))
	return client
}

func (q *Query) Get(tGet *hbase.TGet) (result *hbase.TResult_, err error) {
	client := q.getClient()
	defer client.Close()
	result, err = client.Get(context.Background(), []byte(q.tableName), tGet)
	return
}

func (q *Query) GetByRowKey(rowKey []byte) (result *hbase.TResult_, err error) {
	tGet := &hbase.TGet{
		Row: rowKey,
	}
	if q.columns != nil {
		tGet.Columns = q.columns
	}
	return q.Get(tGet)
}

func (q *Query) GetByRowKeys(rowkeys [][]byte) (results []*hbase.TResult_, err error) {
	tGets := make([]*hbase.TGet, len(rowkeys))
	for k, v := range rowkeys {
		tGet := &hbase.TGet{
			Row: v,
		}
		if q.columns != nil {
			tGet.Columns = q.columns
		}
		tGets[k] = tGet
	}
	return q.GetMultiple(tGets)
}

func (q *Query) GetMultiple(tGets []*hbase.TGet) (results []*hbase.TResult_, err error) {
	client := q.getClient()
	defer client.Close()
	results, err = client.GetMultiple(context.Background(), []byte(q.tableName), tGets)
	return
}

func (q *Query) SetStartRow(startRow []byte) *Query {
	q.startRow = startRow
	return q
}

func (q *Query) SetStopRow(stopRow []byte) *Query {
	q.stopRow = stopRow
	return q
}

func (q *Query) Scan(num int32) (results []*hbase.TResult_, err error) {
	startTime := time.Now()
	defer func() {
		logs.Debug("HBase Scan spend time: %s [%d] [table: %s]", time.Since(startTime), num, q.tableName)
	}()
	tScan := &hbase.TScan{}
	isSettingStartRow := len(q.startRow) > 0
	if isSettingStartRow {
		tScan.StartRow = q.startRow
	}
	if len(q.stopRow) > 0 {
		tScan.StopRow = q.stopRow
	} else if isSettingStartRow {
		q.filter.Add(&hbase.PrefixFilter{Prefix: string(q.startRow)})
	}

	if q.caching != nil {
		tScan.Caching = q.caching
	}

	if q.batchSize != nil {
		tScan.BatchSize = q.batchSize
	}

	if q.reversed != nil {
		tScan.Reversed = q.reversed
	}

	if !q.filter.IsEmpty() {
		tScan.FilterString = q.filter.ToByteArray()
	}

	if q.cacheBlocks != nil {
		tScan.CacheBlocks = q.cacheBlocks
	}

	if q.columns != nil {
		tScan.Columns = q.columns
	}

	client := q.getClient()
	defer client.Close()

	results, err = client.GetScannerResults(context.Background(), []byte(q.tableName), tScan, num)
	return
}
