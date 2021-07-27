package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/models/business"
	"dongchamao/services/dyimg"
	"dongchamao/structinit/repost/dy"
	"time"
)

type ProductController struct {
	controllers.ApiBaseController
}

//
func (receiver *ProductController) ProductBaseAnalysis() {
	productId := receiver.GetString(":product_id", "")
	if productId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	productBusiness := business.NewProductBusiness()
	info, _ := productBusiness.HbaseGetProductDailyRangDate(productId, startTime, endTime)
	receiver.SuccReturn(info)
	return
}

func (receiver *ProductController) ProductBase() {
	productId := receiver.GetString(":product_id", "")
	if productId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	productBusiness := business.NewProductBusiness()
	productInfo, comErr := productBusiness.HbaseGetProductInfo(productId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	yesterdayDate := time.Now().AddDate(0, 0, -1).Format("20060102")
	yesterdayTime, _ := time.ParseInLocation("20060102", yesterdayDate, time.Local)
	startTime := yesterdayTime.AddDate(0, 0, -30)
	monthData, _ := productBusiness.HbaseGetPromotionMonth(productId)
	relatedInfo, _ := productBusiness.HbaseGetProductDailyRangDate(productId, startTime, yesterdayTime)
	var roomNum int
	var awemeNum int
	authorMap := map[string]string{}
	for _, v := range relatedInfo {
		awemeNum += len(v.AwemeList)
		roomNum += len(v.LiveList)
		for _, a := range v.AwemeAuthorList {
			authorMap[a.AuthorId] = a.AuthorId
		}
		for _, a := range v.LiveAuthorList {
			authorMap[a.AuthorId] = a.AuthorId
		}
	}
	var rate30 float64 = 0
	if monthData.PvCount > 0 {
		rate30 = float64(monthData.OrderCount) / float64(monthData.PvCount)
	}
	simpleInfo := dy.SimpleDyProduct{
		ProductID:     productInfo.ProductID,
		Title:         productInfo.Title,
		MarketPrice:   productInfo.MarketPrice,
		Price:         productInfo.Price,
		URL:           productInfo.URL,
		Image:         dyimg.Product(productInfo.Image),
		Status:        productInfo.Status,
		ShopName:      productInfo.ShopName,
		Label:         productInfo.Label,
		Undercarriage: productInfo.Undercarriage,
		CrawlTime:     productInfo.CrawlTime,
		PlatformLabel: productInfo.PlatformLabel,
		MinPrice:      productInfo.MinPrice,
		CosRatio:      productInfo.CosRatio,
		CosRatioMoney: productInfo.CosRatio / 100 * productInfo.Price,
	}
	receiver.SuccReturn(map[string]interface{}{
		"pv_count_30":    monthData.PvCount,
		"order_count_30": monthData.OrderCount,
		"rate_30":        rate30,
		"aweme_num_30":   awemeNum,
		"room_num_30":    roomNum,
		"author_num_30":  len(authorMap),
		"simple_info":    simpleInfo,
	})
	return
}
