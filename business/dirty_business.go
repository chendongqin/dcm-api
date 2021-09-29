package business

import (
	"dongchamao/global"
	"dongchamao/global/logger"
	hbase2 "dongchamao/hbase"
	"dongchamao/models/entity"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/kafka"
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
	if logger.CheckError(err) != nil {
		return global.NewError(5000)
	}
	//ret, _ := NewSpiderBusiness().SpiderSpeedUp("author", authorId)
	//logs.Info("达人分类修改，爬虫推送结果：", ret)
	kafka.SendAuthorCateChange(authorId)
	return nil
}

//修改商品分类
func (receiver *DirtyBusiness) ChangeProductCate(productId, dcmLevelFirst, firstCate, secondCate, thirdCate string) global.CommonError {
	productInfo, comErr := hbase2.GetProductInfo(productId)
	if comErr != nil {
		return comErr
	}
	//测试环境不处理
	if global.IsDev() {
		return nil
	}
	hbaseBusiness := NewHbaseBusiness()
	artificialData := map[string]interface{}{
		"first_cname":  firstCate,
		"second_cname": secondCate,
		"third_cname":  thirdCate,
	}
	tableName := hbaseService.HbaseDyProduct
	if productInfo.PlatformLabel == "小店" {
		tableName = hbaseService.HbaseDyProductBrand
	}
	jsonByte, _ := jsoniter.Marshal(artificialData)
	columnL := hbaseBusiness.BuildColumnValue("other", "manmade_category", string(jsonByte), entity.String)
	columnL2 := hbaseBusiness.BuildColumnValue("info", "dcm_level_first", dcmLevelFirst, entity.String)
	err := hbaseBusiness.PutByRowKey(tableName, productId, []*hbase.TColumnValue{columnL, columnL2})
	if logger.CheckError(err) != nil {
		return global.NewError(5000)
	}
	//ret, _ := NewSpiderBusiness().SpiderSpeedUp("product", productId)
	//logs.Info("商品分类修改，爬虫推送结果：", ret)
	kafka.SendProductCateChange(productId)
	return nil
}
