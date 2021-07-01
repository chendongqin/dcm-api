package hbaseService

import (
	"dongchamao/global/utils"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/msgpack"
	jsoniter "github.com/json-iterator/go"
	"math"
)

type HbaseField struct {
	FieldType string
	FieldName string
}

type HbaseEntity map[string]HbaseField

func HbaseFormat(result *hbase.TResult_, fieldMap HbaseEntity) map[string]interface{} {
	retMap := make(map[string]interface{})
	for _, v := range result.ColumnValues {
		fn := string(v.Qualifier)
		if ai, ok := fieldMap[fn]; ok == true {
			fieldType := ai.FieldType
			fieldName := ai.FieldName
			if fieldType == "m_double" {
				fv := v.Value
				retMap[fieldName], _ = msgpack.UnpackFloat64(fv)
			} else if fieldType == "m_json" {
				fv := v.Value
				tmpMap := map[string]interface{}{}
				jsoniter.Unmarshal(fv, &tmpMap)
				retMap[fieldName] = tmpMap
			} else if fieldType == "m_string" {
				fv := v.Value
				retMap[fieldName], _ = msgpack.UnpackString(fv)
			} else if fieldType == "m_int" {
				fv := v.Value
				val, _ := msgpack.UnpackInt32(fv)
				retMap[fieldName] = int(val)
			} else if fieldType == "m_long" {
				fv := v.Value
				retMap[fieldName], _ = msgpack.UnpackInt64(fv)
			} else if fieldType == "long" {
				fv := utils.ParseByteInt64(v.Value)
				retMap[fieldName] = fv
			} else if fieldType == "int" {
				fv := utils.ParseByteInt32(v.Value)
				retMap[fieldName] = int(fv)
			} else if fieldType == "string" {
				fv := string(v.Value)
				retMap[fieldName] = fv
			} else if fieldType == "float" {
				fv := utils.ParseByteFloat32(v.Value)
				retMap[fieldName] = fv
			} else if fieldType == "double" {
				fv := utils.ParseByteFloat64(v.Value)
				if math.IsNaN(fv) {
					fv = 0
				}
				retMap[fieldName] = fv
			} else if fieldType == "byte" {
				if len(v.Value) == 1 {
					fv := v.Value[0]
					retMap[fieldName] = fv
				} else {
					retMap[fieldName] = v.Value
				}
			} else if fieldType == "bool" {
				if v.Value[0] == uint8(255) {
					retMap[fieldName] = true
				} else {
					retMap[fieldName] = false
				}
			} else {
				retMap[fieldName] = v.Value
			}
		} else {
			retMap[fn] = v.Value
		}
	}
	return retMap
}
