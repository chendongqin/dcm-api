package controllers

import (
	"dongchamao/models/dcm"
	"encoding/json"
	"sort"
)

type ScriptController struct {
	ApiBaseController
}

func (receiver *ScriptController) AuthorTag() {
	var config dcm.DcConfigJson
	db := dcm.GetDbSession().Table(dcm.DcConfigJson{})
	if _, err := db.Where("key_name='author_cate'").Get(&config); err != nil {
		panic(err)
		return
	}
	var tag AuthorTag
	println(config.Value)
	if err := json.Unmarshal([]byte(config.Value), &tag); err != nil {
		panic(err)
	}
	var key []string
	var dataMap = make(map[string]string)
	for _, v := range tag.Tag {
		for kk, vv := range v.First {
			key = append(key, kk)
			dataMap[kk] = vv
		}
		for _, vv := range v.Second {
			for kkk, vvv := range vv {
				key = append(key, kkk)
				dataMap[kkk] = vvv
			}
		}
	}
	sort.Strings(key)
	for _, v := range key {
		println(dataMap[v])
	}
	receiver.SuccReturn(dataMap)
	return
}

type AuthorTag struct {
	Tag []struct {
		First  map[string]string   `json:"first"`
		Second []map[string]string `json:"second"`
	} `json:"tag"`
}
