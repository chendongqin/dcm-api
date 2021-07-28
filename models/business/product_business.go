package business

import (
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"fmt"
	"strings"
	"time"
)

type ProductBusiness struct {
}

func NewProductBusiness() *ProductBusiness {
	return new(ProductBusiness)
}

func (receiver ProductBusiness) HbaseGetProductInfo(productId string) (data entity.DyProduct, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyProduct).GetByRowKey([]byte(productId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyProductMap)
	utils.MapToStruct(detailMap, &data)
	data.ProductID = productId
	return
}

//获取商品品牌数据
func (receiver *ProductBusiness) HbaseGetDyProductBrand(productId string) (data entity.DyProductBrand, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyProductBrand).GetByRowKey([]byte(productId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyProductBrandMap)
	utils.MapToStruct(detailMap, &data)
	return
}

func (receiver *ProductBusiness) HbaseGetPromotionMonth(productId string) (data entity.DyLivePromotionMonth, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLivePromotionMonth).GetByRowKey([]byte(productId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLivePromotionMonthMap)
	utils.MapToStruct(detailMap, &data)
	return
}

func (receiver *ProductBusiness) HbaseGetProductDailyRangDate(productId string, startTime, endTime time.Time) (data map[string]entity.DyProductDaily, comErr global.CommonError) {
	data = map[string]entity.DyProductDaily{}
	query := hbasehelper.NewQuery()
	startRow := productId + "_" + startTime.Format("20060102")
	endRow := productId + "_" + endTime.AddDate(0, 0, 1).Format("20060102")
	results, err := query.
		SetTable(hbaseService.HbaseDyProductDaily).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity.DyProductDailyMap)
		hData := entity.DyProductDaily{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData
	}
	return
}

//直播间销量趋势
func (receiver *ProductBusiness) HbaseGetProductLiveSalesRangDate(productId string, startTime, endTime time.Time) (data map[string]entity.DyProductLiveSalesTrend, comErr global.CommonError) {
	data = map[string]entity.DyProductLiveSalesTrend{}
	query := hbasehelper.NewQuery()
	startRow := productId + "_" + startTime.Format("20060102")
	endRow := productId + "_" + endTime.AddDate(0, 0, 1).Format("20060102")
	results, err := query.
		SetTable(hbaseService.HbaseDyProductLiveSalesTrend).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity.DyProductLiveSalesTrendMap)
		hData := entity.DyProductLiveSalesTrend{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData
	}
	return
}

//获取商品url
func (receiver *ProductBusiness) GetProductUrl(platform, productId string) string {
	url := ""
	switch platform {
	case "淘宝":
		url = "https://item.taobao.com/item.htm?id=%s"
	case "京东":
		url = " https://item.m.jd.com/product/%s.html"
	case "天猫":
		url = " https://detail.tmall.com/item.htm?id=%s"
	case "苏宁":
		url = "https://m.suning.com/product/0000000000/0000000%s.html"
	case "小店":
		url = "https://haohuo.jinritemai.com/views/product/item2?id=%s"
	case "唯品会":
		url = "https://m.vip.com/public/go.html?pid=%s"
	case "考拉":
		url = "https://m-goods.kaola.com/product/%s.html"
	}
	if url != "" {
		url = fmt.Sprintf(url, productId)
	}
	return url
}
