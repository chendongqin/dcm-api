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

//商品分析
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
	dateChart := make([]string, 0)
	hotAuthorChart := make([]int, 0)
	liveAuthorChart := make([]int, 0)
	awemeAuthorChart := make([]int, 0)
	awemeChart := make([]int, 0)
	roomChart := make([]int, 0)
	beginTime := startTime
	for {
		if beginTime.After(endTime) {
			break
		}
		dateChart = append(dateChart, beginTime.Format("01/02"))
		dateKey := beginTime.Format("20060102")
		if v, ok := info[dateKey]; ok {
			authors := map[string]string{}
			awemeAuthorNum := 0
			liveAuthorNum := 0
			for _, a := range v.AwemeAuthorList {
				awemeAuthorNum++
				authors[a.AuthorId] = a.AuthorId
			}
			for _, a := range v.LiveAuthorList {
				liveAuthorNum++
				authors[a.AuthorId] = a.AuthorId
			}
			hotAuthorChart = append(hotAuthorChart, len(authors))
			liveAuthorChart = append(liveAuthorChart, liveAuthorNum)
			awemeAuthorChart = append(awemeAuthorChart, awemeAuthorNum)
			awemeChart = append(awemeChart, len(v.AwemeList))
			roomChart = append(roomChart, len(v.LiveList))
		} else {
			hotAuthorChart = append(hotAuthorChart, 0)
			liveAuthorChart = append(liveAuthorChart, 0)
			awemeAuthorChart = append(awemeAuthorChart, 0)
			awemeChart = append(awemeChart, 0)
			roomChart = append(roomChart, 0)
		}
		beginTime = beginTime.AddDate(0, 0, 1)
	}
	receiver.SuccReturn(map[string]interface{}{
		"date":         dateChart,
		"hot_author":   hotAuthorChart,
		"live_author":  liveAuthorChart,
		"aweme_author": awemeAuthorChart,
		"aweme":        awemeChart,
		"room":         roomChart,
	})
	return
}

//商品基础数据
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
	brandInfo, _ := productBusiness.HbaseGetDyProductBrand(productId)
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
	if productInfo.MinPrice == 0 {
		productInfo.MinPrice = productInfo.Price
	}
	shopName := brandInfo.ShopName
	if shopName == "" {
		shopName = productInfo.TbNick
	}
	label := brandInfo.DcmLevelFirst
	if label == "" {
		label = "其他"
	}
	simpleInfo := dy.SimpleDyProduct{
		ProductID:     productInfo.ProductID,
		Title:         productInfo.Title,
		MarketPrice:   productInfo.MarketPrice,
		Price:         productInfo.Price,
		URL:           productBusiness.GetProductUrl(productInfo.PlatformLabel, productInfo.ProductID),
		Image:         dyimg.Product(productInfo.Image),
		Status:        productInfo.Status,
		ShopId:        productInfo.ShopID,
		ShopName:      shopName,
		Label:         label,
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
