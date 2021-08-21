package business

import (
	"dongchamao/global"
	"dongchamao/models/entity"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbase"
	jsoniter "github.com/json-iterator/go"
)

type DirtyBusiness struct {
}

func NewDirtyBusiness() *DirtyBusiness {
	return new(DirtyBusiness)
}

//修改达人分类
func (receiver *DirtyBusiness) ChangeAuthorCate(authorId, tags, tagsTow string) global.CommonError {
	hbaseBusiness := NewHbaseBusiness()
	artificialData := map[string]interface{}{
		"tags":           tags,
		"tags_level_tow": tagsTow,
	}
	jsonByte, _ := jsoniter.Marshal(artificialData)
	columnL := hbaseBusiness.BuildColumnValue("other", "artificial_data", string(jsonByte), entity.String)
	err := hbaseBusiness.PutByRowKey(hbaseService.HbaseDyAuthor, authorId, []*hbase.TColumnValue{columnL})
	if err != nil {
		return global.NewError(5000)
	}
	return nil
}
