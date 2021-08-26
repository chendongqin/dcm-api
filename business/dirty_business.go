package business

import (
	"dongchamao/global"
	hbase2 "dongchamao/hbase"
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
	_, comErr := hbase2.GetAuthor(authorId)
	if comErr != nil {
		return comErr
	}
	//测试环境不处理
	if global.IsDev() {
		return nil
	}
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

//修改商品分类
func (receiver *DirtyBusiness) ChangeProductCate(productId, dcmLevelFirst, firstCate, secondCate, thirdCate string) global.CommonError {
	_, comErr := hbase2.GetProductInfo(productId)
	if comErr != nil {
		return comErr
	}
	//测试环境不处理
	if global.IsDev() {
		return nil
	}
	hbaseBusiness := NewHbaseBusiness()
	artificialData := map[string]interface{}{
		"dcm_level_first": dcmLevelFirst,
		"first_cname":     firstCate,
		"second_cname":    secondCate,
		"third_cname":     thirdCate,
	}
	jsonByte, _ := jsoniter.Marshal(artificialData)
	columnL := hbaseBusiness.BuildColumnValue("other", "artificial_data", string(jsonByte), entity.String)
	err := hbaseBusiness.PutByRowKey(hbaseService.HbaseDyProduct, productId, []*hbase.TColumnValue{columnL})
	if err != nil {
		return global.NewError(5000)
	}
	return nil
}
