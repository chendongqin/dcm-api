package hbaseService

import (
	"dongchamao/global/utils"
	entity2 "dongchamao/models/hbase/entity"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/msgpack"
	jsoniter "github.com/json-iterator/go"
	"math"
)

func HbaseFormat(result *hbase.TResult_, fieldMap entity2.HbaseEntity) map[string]interface{} {
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
			if fieldType == entity2.MDouble {
				fv := v.Value
				retMap[fieldName], _ = msgpack.UnpackFloat64(fv)
			} else if fieldType == entity2.Json {
				fv := v.Value
				tmpMap := map[string]interface{}{}
				jsoniter.Unmarshal(fv, &tmpMap)
				retMap[fieldName] = tmpMap
			} else if fieldType == entity2.AJson {
				fv := v.Value
				tmpMap := make([]map[string]interface{}, 0)
				jsoniter.Unmarshal(fv, &tmpMap)
				retMap[fieldName] = tmpMap
			} else if fieldType == entity2.MString {
				fv := v.Value
				retMap[fieldName], _ = msgpack.UnpackString(fv)
			} else if fieldType == entity2.MInt {
				fv := v.Value
				val, _ := msgpack.UnpackInt32(fv)
				retMap[fieldName] = int(val)
			} else if fieldType == entity2.MLong {
				fv := v.Value
				retMap[fieldName], _ = msgpack.UnpackInt64(fv)
			} else if fieldType == entity2.Long {
				fv := utils.ParseByteInt64(v.Value)
				retMap[fieldName] = fv
			} else if fieldType == entity2.Int {
				fv := utils.ParseByteInt32(v.Value)
				retMap[fieldName] = int(fv)
			} else if fieldType == entity2.String {
				fv := string(v.Value)
				retMap[fieldName] = fv
			} else if fieldType == entity2.Float {
				fv := utils.ParseByteFloat32(v.Value)
				retMap[fieldName] = fv
			} else if fieldType == entity2.Double {
				fv := utils.ParseByteFloat64(v.Value)
				if math.IsNaN(fv) {
					fv = 0
				}
				retMap[fieldName] = fv
			} else if fieldType == entity2.Byte {
				if len(v.Value) == 1 {
					fv := v.Value[0]
					retMap[fieldName] = fv
				} else {
					retMap[fieldName] = v.Value
				}
			} else if fieldType == entity2.Bool {
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
