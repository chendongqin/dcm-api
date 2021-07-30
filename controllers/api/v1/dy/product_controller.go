package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/business"
	"dongchamao/models/business/es"
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
	monthData, _ := productBusiness.HbaseGetPromotionMonth(productId)
	dailyMapData := map[string]entity.DyLiveProductDaily{}
	for _, v := range monthData.DailyList {
		t, _ := time.ParseInLocation("2006/01/02", v.StatisticsTime, time.Local)
		if t.Before(startTime) || t.After(endTime) {
			continue
		}
		dailyMapData[t.Format("20060102")] = v
	}
	dateChart := make([]string, 0)
	hotAuthorChart := make([]int, 0)
	liveAuthorChart := make([]int, 0)
	awemeAuthorChart := make([]int, 0)
	awemeChart := make([]int, 0)
	roomChart := make([]int, 0)
	orderChart := make([]int64, 0)
	pvChart := make([]int64, 0)
	rateChart := make([]float64, 0)
	orderList := make([]dy.ProductOrderDaily, 0)
	countData := dy.ProductOrderDaily{}
	beginTime := startTime
	for {
		if beginTime.After(endTime) {
			break
		}
		dateStr := beginTime.Format("01/02")
		dateChart = append(dateChart, dateStr)
		dateKey := beginTime.Format("20060102")
		awemeAuthorNum := 0
		liveAuthorNum := 0
		authorNum := 0
		awemeNum := 0
		roomNum := 0
		var order int64 = 0
		var pv int64 = 0
		var rate float64 = 0
		if v, ok := info[dateKey]; ok {
			authors := map[string]string{}
			for _, a := range v.AwemeAuthorList {
				awemeAuthorNum++
				authors[a.AuthorId] = a.AuthorId
			}
			for _, a := range v.LiveAuthorList {
				liveAuthorNum++
				authors[a.AuthorId] = a.AuthorId
			}
			authorNum = len(authors)
			awemeNum = len(v.AwemeList)
			roomNum = len(v.LiveList)
		}
		if d, ok := dailyMapData[dateKey]; ok {
			order = d.ProductOrderAccount
			pv = d.Pv
			if d.Pv > 0 {
				rate = float64(d.ProductOrderAccount) / float64(d.Pv)
			}
		}
		hotAuthorChart = append(hotAuthorChart, authorNum)
		liveAuthorChart = append(liveAuthorChart, liveAuthorNum)
		awemeAuthorChart = append(awemeAuthorChart, awemeAuthorNum)
		awemeChart = append(awemeChart, awemeNum)
		roomChart = append(roomChart, roomNum)
		orderChart = append(orderChart, order)
		pvChart = append(pvChart, pv)
		rateChart = append(rateChart, rate)
		countData.OrderCount += order
		countData.PvCount += pv
		countData.AwemeNum += awemeNum
		countData.RoomNum += roomNum
		countData.AuthorNum += authorNum
		orderList = append(orderList, dy.ProductOrderDaily{
			Date:       dateStr,
			OrderCount: order,
			PvCount:    pv,
			Rate:       rate,
			AwemeNum:   awemeNum,
			RoomNum:    roomNum,
			AuthorNum:  authorNum,
		})
		beginTime = beginTime.AddDate(0, 0, 1)
	}
	if countData.PvCount > 0 {
		countData.Rate = float64(countData.OrderCount) / float64(countData.PvCount)
	}
	receiver.SuccReturn(map[string]interface{}{
		"author_chart": dy.ProductAuthorChart{
			Date:             dateChart,
			AuthorCount:      hotAuthorChart,
			AwemeAuthorCount: awemeAuthorChart,
			LiveAuthorCount:  liveAuthorChart,
		},
		"count_chart": dy.ProductLiveAwemeChart{
			Date:       dateChart,
			LiveCount:  roomChart,
			AwemeCount: awemeChart,
		},
		"order_chart": dy.ProductOrderChart{
			Date:       dateChart,
			OrderCount: orderChart,
			PvCount:    pvChart,
			Rate:       rateChart,
		},
		"daily_list":  orderList,
		"order_count": countData,
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
	dateChart7 := make([]int64, 0)
	priceChart7 := make([]float64, 0)
	cosPriceChart7 := make([]float64, 0)
	dateChart15 := make([]int64, 0)
	priceChart15 := make([]float64, 0)
	cosPriceChart15 := make([]float64, 0)
	dateChart30 := make([]int64, 0)
	priceChart30 := make([]float64, 0)
	cosPriceChart30 := make([]float64, 0)
	last30Day := utils.ToInt64(time.Now().AddDate(0, 0, -30).Format("20060102"))
	last15Day := utils.ToInt64(time.Now().AddDate(0, 0, -15).Format("20060102"))
	last7Day := utils.ToInt64(time.Now().AddDate(0, 0, -7).Format("20060102"))
	priceTrends := business.ProductPriceTrendsListOrderByTime(productInfo.PriceTrends)
	priceMap := map[int64]entity.DyProductPriceTrend{}
	for _, v := range priceTrends {
		if last30Day > v.StartTime {
			continue
		}
		priceMap[v.StartTime] = v
	}
	begin, _ := time.ParseInLocation("20060102", time.Now().AddDate(0, 0, -30).Format("20060102"), time.Local)
	beforeData := entity.DyProductPriceTrend{}
	for {
		if begin.After(time.Now()) {
			break
		}
		nowDate := utils.ToInt64(begin.Format("20060102"))
		if v, ok := priceMap[nowDate]; ok {
			beforeData = v
		} else {
			beforeData.StartTime = nowDate
		}
		cosPrice := beforeData.Price * productInfo.CosRatio / 100
		if beforeData.StartTime > last7Day {
			dateChart7 = append(dateChart7, begin.Unix())
			priceChart7 = append(priceChart7, beforeData.Price)
			cosPriceChart7 = append(cosPriceChart7, cosPrice)
		}
		if beforeData.StartTime > last15Day {
			dateChart15 = append(dateChart15, begin.Unix())
			priceChart15 = append(priceChart15, beforeData.Price)
			cosPriceChart15 = append(cosPriceChart15, cosPrice)
		}
		dateChart30 = append(dateChart30, begin.Unix())
		priceChart30 = append(priceChart30, beforeData.Price)
		cosPriceChart30 = append(cosPriceChart30, cosPrice)
		begin = begin.AddDate(0, 0, 1)
	}
	receiver.SuccReturn(map[string]interface{}{
		"pv_count_30":    monthData.PvCount,
		"order_count_30": monthData.OrderCount,
		"rate_30":        rate30,
		"aweme_num_30":   awemeNum,
		"room_num_30":    roomNum,
		"author_num_30":  len(authorMap),
		"simple_info":    simpleInfo,
		"chart_7": map[string]interface{}{
			"date":      dateChart7,
			"price":     priceChart7,
			"cos_price": cosPriceChart7,
		},
		"chart_15": map[string]interface{}{
			"date":      dateChart15,
			"price":     priceChart15,
			"cos_price": cosPriceChart15,
		},
		"chart_30": map[string]interface{}{
			"date":      dateChart30,
			"price":     priceChart30,
			"cos_price": cosPriceChart30,
		},
	})
	return
}

//商品销量趋势
func (receiver *ProductController) ProductLiveChart() {
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
	infoMap, _ := productBusiness.HbaseGetProductLiveSalesRangDate(productId, startTime, endTime)
	dateChart := make([]string, 0)
	saleChart := make([]int64, 0)
	roomNumChart := make([]int, 0)
	priceChart := make([]float64, 0)
	beginTime := startTime
	for {
		if beginTime.After(endTime) {
			break
		}
		var sale int64 = 0
		roomNum := 0
		var price float64 = 0
		if v, ok := infoMap[beginTime.Format("20060102")]; ok {
			sale = v.Sales
			roomNum = v.RoomNum
			price = v.Price
		}
		dateChart = append(dateChart, beginTime.Format("01/02"))
		saleChart = append(saleChart, sale)
		roomNumChart = append(roomNumChart, roomNum)
		priceChart = append(priceChart, price)
		beginTime = beginTime.AddDate(0, 0, 1)
	}
	receiver.SuccReturn(map[string]interface{}{
		"date":     dateChart,
		"price":    priceChart,
		"sale":     saleChart,
		"room_num": roomNumChart,
	})
	return
}

//商品直播间列表
func (receiver *ProductController) ProductLiveRoomList() {
	productId := receiver.Ctx.Input.Param(":product_id")
	InputData := receiver.InputFormat()
	keyword := InputData.GetString("keyword", "")
	sortStr := InputData.GetString("sort", "predict_gmv")
	orderBy := InputData.GetString("order_by", "desc")
	page := InputData.GetInt("page", 1)
	size := InputData.GetInt("page_size", 10)
	if productId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	t1, t2, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	esLiveBusiness := es.NewEsLiveBusiness()
	list, total, comErr := esLiveBusiness.SearchProductRooms(productId, keyword, sortStr, orderBy, page, size, t1, t2)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	countList := make([]dy.LiveRoomProductCount, 0)
	if len(list) > 0 {
		liveBusiness := business.NewLiveBusiness()
		authorBusiness := business.NewAuthorBusiness()
		curMap := map[string]dy.LiveCurProductCount{}
		pmtMap := map[string][]dy.LiveRoomProductSaleStatus{}
		curChan := make(chan map[string]dy.LiveCurProductCount, 0)
		pmtChan := make(chan map[string][]dy.LiveRoomProductSaleStatus, 0)
		authorChan := make(chan map[string]entity.DyAuthorData, 0)
		authorMap := map[string]entity.DyAuthorData{}
		for _, v := range list {
			go func(curChan chan map[string]dy.LiveCurProductCount, pmtChan chan map[string][]dy.LiveRoomProductSaleStatus, roomId, productId, authorId string) {
				curMapTmp := liveBusiness.RoomCurProductByIds(v.RoomID, []string{productId})
				pmtMapTmp := liveBusiness.RoomPmtProductByIds(v.RoomID, []string{productId})
				authorData, _ := authorBusiness.HbaseGetAuthor(authorId)
				authorMapTmp := map[string]entity.DyAuthorData{}
				roomProductCurMap := map[string]dy.LiveCurProductCount{}
				roomProductPmtMap := map[string][]dy.LiveRoomProductSaleStatus{}
				if v, ok := curMapTmp[productId]; ok {
					roomProductCurMap[roomId] = v
				}
				if v, ok := pmtMapTmp[productId]; ok {
					roomProductPmtMap[roomId] = v
				}
				authorMapTmp[authorId] = authorData
				curChan <- roomProductCurMap
				pmtChan <- roomProductPmtMap
				authorChan <- authorMapTmp
			}(curChan, pmtChan, v.RoomID, v.ProductID, v.AuthorID)
		}
		for i := 0; i < len(list); i++ {
			cur, ok := <-curChan
			if !ok {
				break
			}
			for k, v := range cur {
				curMap[k] = v
			}
		}
		for i := 0; i < len(list); i++ {
			pmt, ok := <-pmtChan
			if !ok {
				break
			}
			for k, v := range pmt {
				pmtMap[k] = v
			}
		}
		for i := 0; i < len(list); i++ {
			aMap, ok := <-authorChan
			if !ok {
				break
			}
			for k, v := range aMap {
				authorMap[k] = v
			}
		}
		for _, v := range list {
			if author, ok := authorMap[v.AuthorID]; ok {
				v.AuthorRoomID = author.RoomID
				v.Avatar = dyimg.Fix(author.Avatar)
			}
			item := dy.LiveRoomProductCount{
				ProductInfo: v,
				ProductStartSale: dy.RoomProductSaleChart{
					Timestamp: []int64{},
					Sales:     []int64{},
				},
				ProductEndSale: dy.RoomProductSaleChart{
					Timestamp: []int64{},
					Sales:     []int64{},
				},
			}
			if s, ok := pmtMap[v.RoomID]; ok {
				for _, s1 := range s {
					item.ProductStartSale.Timestamp = append(item.ProductStartSale.Timestamp, s1.StartTime)
					item.ProductStartSale.Sales = append(item.ProductStartSale.Sales, s1.StartSales)
					if s1.StopTime > 0 {
						item.ProductEndSale.Timestamp = append(item.ProductEndSale.Timestamp, s1.StopTime)
						item.ProductEndSale.Sales = append(item.ProductEndSale.Sales, s1.FinalSales)
					}
				}
			}
			if c, ok := curMap[v.RoomID]; ok {
				c.CurList = business.ProductCurOrderByTime(c.CurList)
				item.ProductCur = c
			} else {
				item.ProductCur = dy.LiveCurProductCount{
					CurList: []dy.LiveCurProduct{},
				}
			}
			countList = append(countList, item)
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  countList,
		"total": total,
	})
	return
}
