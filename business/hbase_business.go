package business

import (
	"context"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/entity"
	"dongchamao/services/hbaseService/hbase"
	"github.com/astaxie/beego/logs"
)

type HbaseBusiness struct {
}

func NewHbaseBusiness() *HbaseBusiness {
	return new(HbaseBusiness)
}

func (receiver *HbaseBusiness) BuildColumnValue(family string, qualifier string, value interface{}, valueTypes ...string) *hbase.TColumnValue {
	var finalValue []byte
	if val, ok := value.([]byte); ok {
		finalValue = val
	} else {
		valueType := entity.String
		if len(valueTypes) > 0 {
			valueType = valueTypes[0]
		}
		switch valueType {
		case entity.Long:
			finalValue = utils.Int64ToBytes(utils.ToInt64(value))
		case entity.String:
			finalValue = []byte(utils.ToString(value))
		case entity.Int:
			finalValue = utils.IntToBytes(utils.ToInt(value))
		case entity.Double:
			finalValue = utils.DoubleToBytes(utils.ToFloat64(value))
		default:
			return nil
		}
	}
	return &hbase.TColumnValue{
		Family:    []byte(family),
		Qualifier: []byte(qualifier),
		Value:     finalValue,
	}
}

func (receiver *HbaseBusiness) PutByRowKey(tableName string, rowKey string, values []*hbase.TColumnValue) (err error) {
	client := global.HbasePools.Get("default")
	defer client.Close()
	tableInBytes := []byte(tableName)
	data := &hbase.TPut{
		Row:          []byte(rowKey),
		ColumnValues: values,
	}
	err = client.Put(context.Background(), tableInBytes, data)
	if err != nil {
		logs.Error("Data error:", err)
	}
	return
}
