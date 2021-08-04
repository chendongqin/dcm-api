package hbaseService

import (
	"dongchamao/global/utils"
	"dongchamao/models/entity"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/msgpack"
	jsoniter "github.com/json-iterator/go"
	"math"
)

func HbaseFormat(result *hbase.TResult_, fieldMap entity.HbaseEntity) map[string]interface{} {
	retMap := make(map[string]interface{})
	for _, v := range result.ColumnValues {
		family := string(v.Family)
		//if family == "other" {
		//	continue
		//}
		fn := string(v.Qualifier)
		if family != "info" {
			fn = family + "_" + fn
		}
		if ai, ok := fieldMap[fn]; ok == true {
			fieldType := ai.FieldType
			fieldName := ai.FieldName
			if fieldType == entity.MDouble {
				fv := v.Value
				retMap[fieldName], _ = msgpack.UnpackFloat64(fv)
			} else if fieldType == entity.Json {
				fv := v.Value
				tmpMap := map[string]interface{}{}
				jsoniter.Unmarshal(fv, &tmpMap)
				retMap[fieldName] = tmpMap
			} else if fieldType == entity.AJson {
				fv := v.Value
				tmpMap := make([]map[string]interface{}, 0)
				jsoniter.Unmarshal(fv, &tmpMap)
				retMap[fieldName] = tmpMap
			} else if fieldType == entity.MString {
				fv := v.Value
				retMap[fieldName], _ = msgpack.UnpackString(fv)
			} else if fieldType == entity.MInt {
				fv := v.Value
				val, _ := msgpack.UnpackInt32(fv)
				retMap[fieldName] = int(val)
			} else if fieldType == entity.MLong {
				fv := v.Value
				retMap[fieldName], _ = msgpack.UnpackInt64(fv)
			} else if fieldType == entity.Long {
				fv := utils.ParseByteInt64(v.Value)
				retMap[fieldName] = fv
			} else if fieldType == entity.Int {
				fv := utils.ParseByteInt32(v.Value)
				retMap[fieldName] = int(fv)
			} else if fieldType == entity.String {
				fv := string(v.Value)
				retMap[fieldName] = fv
			} else if fieldType == entity.Float {
				fv := utils.ParseByteFloat32(v.Value)
				retMap[fieldName] = fv
			} else if fieldType == entity.Double {
				fv := utils.ParseByteFloat64(v.Value)
				if math.IsNaN(fv) {
					fv = 0
				}
				retMap[fieldName] = fv
			} else if fieldType == entity.Byte {
				if len(v.Value) == 1 {
					fv := v.Value[0]
					retMap[fieldName] = fv
				} else {
					retMap[fieldName] = v.Value
				}
			} else if fieldType == entity.Bool {
				if v.Value[0] == uint8(255) {
					retMap[fieldName] = true
				} else {
					retMap[fieldName] = false
				}
			} else {
				retMap[fieldName] = v.Value
			}
		}

	}
	return retMap
}
