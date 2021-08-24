package controllers

import (
	"dongchamao/models/dcm"
	"encoding/json"
	"sort"
	"strconv"
)

type ScriptController struct {
	ApiBaseController
}

func (receiver *ScriptController) AuthorTag() {
	var config dcm.DcConfigJson
	var cate dcm.DcAuthorCate
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
	var keySlice []int
	var dataMap = make(map[int]dcm.DcAuthorCate)
	for _, v := range tag.Tag {
		var parentId int
		for kk, vv := range v.First {
			firstKey, _ := strconv.Atoi(kk)
			keySlice = append(keySlice, firstKey)
			parentId, _ = strconv.Atoi(kk)
			cate.Id = parentId
			cate.Name = vv
			cate.Level = 1
			cate.ParentId = 0
			dataMap[firstKey] = cate
		}
		for _, vv := range v.Second {
			for kkk, vvv := range vv {
				secondKey, _ := strconv.Atoi(kkk)
				keySlice = append(keySlice, secondKey)
				cate.Id = secondKey
				cate.Name = vvv
				cate.Level = 2
				cate.ParentId = parentId
				dataMap[secondKey] = cate
			}
		}
	}
	sort.Ints(keySlice)
	for _, v := range keySlice {
		if _, err := db.Insert(dataMap[v]); err != nil {
			return
		}
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
