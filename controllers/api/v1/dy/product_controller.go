package dy

import (
	"dongchamao/business"
	"dongchamao/business/es"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	dy2 "dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	"math"
	"sort"
	"time"
)

type ProductController struct {
	controllers.ApiBaseController
}

func (receiver *ProductController) GetCacheProductCate() {
	productBusiness := business.NewProductBusiness()
	cateList := productBusiness.GetCacheProductCate(true)
	receiver.SuccReturn(cateList)
	return
}

func (receiver *ProductController) Search() {
	hasAuth := false
	hasLogin := false
	if receiver.DyLevel == 3 {
		hasAuth = true
	}
	if receiver.UserId > 0 {
		hasLogin = true
	}
	keyword := receiver.GetString("keyword", "")
	category := receiver.GetString("category", "")
	secondCategory := receiver.GetString("second_category", "")
	thirdCategory := receiver.GetString("third_category", "")
	platform := receiver.GetString("platform", "")
	sortStr := receiver.GetString("sort", "")
	orderBy := receiver.GetString("order_by", "")
	minCommissionRate, _ := receiver.GetFloat("min_commission_rate", 0)
	minPrice, _ := receiver.GetFloat("min_price", 0)
	maxPrice, _ := receiver.GetFloat("max_price", 0)
	commerceType, _ := receiver.GetInt("commerce_type", 0)
	isCoupon, _ := receiver.GetInt("is_coupon", 0)
	isStar, _ := receiver.GetInt("is_star", 0)
	notStar, _ := receiver.GetInt("not_star", 0)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	if !hasAuth {
		if category != "" || secondCategory != "" || thirdCategory != "" || platform != "" || minCommissionRate > 0 || minPrice > 0 || maxPrice > 0 || commerceType > 0 ||
			isCoupon > 0 || isStar > 0 || notStar > 0 || page != 1 {
			if !hasLogin {
				receiver.FailReturn(global.NewError(4001))
				return
			}
			receiver.FailReturn(global.NewError(4004))
			return
		}
		if pageSize > 10 {
			pageSize = 10
		}
	}
	formNum := (page - 1) * pageSize
	if formNum > business.DyJewelBaseShowNum {
		receiver.FailReturn(global.NewError(4004))
		return
	}
	productId := ""
	productBusiness := business.NewProductBusiness()
	if keyword != "" {
		itemId := productBusiness.UrlExplain(keyword)
		if itemId != "" {
			if itemId != "" {
				productId = itemId
				keyword = ""
			}
		} else {
			tbShortUrl := utils.ParseTaobaoShare(keyword)
			if tbShortUrl != "" {
				url := productBusiness.ExplainTaobaoShortUrl(tbShortUrl)
				id := productBusiness.UrlExplain(url)
				if id != "" {
					productId = id
					keyword = ""
				} else {
					page = 0
					pageSize = 0
				}
			}
		}
	}
	esProductBusiness := es.NewEsProductBusiness()
	list, total, comErr := esProductBusiness.BaseSearch(productId, keyword, category, secondCategory, thirdCategory, platform,
		minCommissionRate, minPrice, maxPrice, commerceType, isCoupon, isStar, notStar, page, pageSize, sortStr, orderBy)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		list[k].Image = dyimg.Fix(v.Image)
		list[k].ProductId = business.IdEncrypt(v.ProductId)
	}
	totalPage := math.Ceil(float64(total) / float64(pageSize))
	maxPage := math.Ceil(float64(business.DyJewelBaseShowNum) / float64(pageSize))
	if totalPage > maxPage {
		totalPage = maxPage
	}
	maxTotal := business.DyJewelBaseShowNum
	if maxTotal > total {
		maxTotal = total
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":       list,
		"total":      total,
		"total_page": totalPage,
		"max_num":    maxTotal,
		"has_auth":   hasAuth,
		"has_login":  hasLogin,
	})
	return
}

//商品分析
func (receiver *ProductController) ProductBaseAnalysis() {
	productId := business.IdDecrypt(receiver.GetString(":product_id", ""))
	if productId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	info, _ := hbase.GetProductDailyRangDate(productId, startTime, endTime)
	monthData, _ := hbase.GetPromotionMonth(productId)
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
	orderList := make([]dy2.ProductOrderDaily, 0)
	countData := dy2.ProductOrderDaily{}
	beginTime := startTime
	authorMap := map[string]string{}
	roomMap := map[string]string{}
	videoMap := map[string]string{}
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
				authorMap[a.AuthorId] = a.AuthorId
			}
			for _, a := range v.LiveAuthorList {
				liveAuthorNum++
				authors[a.AuthorId] = a.AuthorId
				authorMap[a.AuthorId] = a.AuthorId
			}
			for _, aw := range v.AwemeList {
				videoMap[aw.AwemeId] = aw.AwemeId
			}
			for _, r := range v.LiveList {
				roomMap[r.RoomId] = r.RoomId
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
		orderList = append(orderList, dy2.ProductOrderDaily{
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
	countData.AwemeNum = len(videoMap)
	countData.RoomNum = len(roomMap)
	countData.AuthorNum = len(authorMap)
	if countData.PvCount > 0 {
		countData.Rate = float64(countData.OrderCount) / float64(countData.PvCount)
	}
	sort.Slice(orderList, func(i, j int) bool {
		return orderList[i].Date > orderList[j].Date
	})
	receiver.SuccReturn(map[string]interface{}{
		"author_chart": dy2.ProductAuthorChart{
			Date:             dateChart,
			AuthorCount:      hotAuthorChart,
			AwemeAuthorCount: awemeAuthorChart,
			LiveAuthorCount:  liveAuthorChart,
		},
		"count_chart": dy2.ProductLiveAwemeChart{
			Date:       dateChart,
			LiveCount:  roomChart,
			AwemeCount: awemeChart,
		},
		"order_chart": dy2.ProductOrderChart{
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
	productId := business.IdDecrypt(receiver.GetString(":product_id", ""))
	if productId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	productBusiness := business.NewProductBusiness()
	productInfo, comErr := hbase.GetProductInfo(productId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	brandInfo, _ := hbase.GetDyProductBrand(productId)
	yesterdayDate := time.Now().AddDate(0, 0, -1).Format("20060102")
	yesterdayTime, _ := time.ParseInLocation("20060102", yesterdayDate, time.Local)
	startTime := yesterdayTime.AddDate(0, 0, -30)
	monthData, _ := hbase.GetPromotionMonth(productId)
	relatedInfo, _ := hbase.GetProductDailyRangDate(productId, startTime, yesterdayTime)
	var roomNum int
	var awemeNum int
	roomMap := map[string]string{}
	awemeMap := map[string]string{}
	authorMap := map[string]string{}
	for _, v := range relatedInfo {
		for _, aw := range v.AwemeList {
			awemeMap[aw.AwemeId] = aw.AwemeId
		}
		for _, r := range v.LiveList {
			roomMap[r.RoomId] = r.RoomId
		}
		for _, a := range v.AwemeAuthorList {
			authorMap[a.AuthorId] = a.AuthorId
		}
		for _, a := range v.LiveAuthorList {
			authorMap[a.AuthorId] = a.AuthorId
		}
	}
	roomNum = len(roomMap)
	awemeNum = len(awemeMap)
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
	simpleInfo := dy2.SimpleDyProduct{
		ProductID:     business.IdEncrypt(productInfo.ProductID),
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
	today := utils.ToInt64(time.Now().Format("20060102"))
	last30Day := utils.ToInt64(time.Now().AddDate(0, 0, -29).Format("20060102"))
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
	begin, _ := time.ParseInLocation("20060102", time.Now().AddDate(0, 0, -29).Format("20060102"), time.Local)
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
		if beforeData.StartTime >= last7Day && today != beforeData.StartTime {
			dateChart7 = append(dateChart7, begin.Unix())
			priceChart7 = append(priceChart7, beforeData.Price)
			cosPriceChart7 = append(cosPriceChart7, cosPrice)
		}
		if beforeData.StartTime >= last15Day && today != beforeData.StartTime {
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
	productId := business.IdDecrypt(receiver.GetString(":product_id", ""))
	if productId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	infoMap, _ := hbase.GetProductLiveSalesRangDate(productId, startTime, endTime)
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
	productId := business.IdDecrypt(receiver.Ctx.Input.Param(":product_id"))
	InputData := receiver.InputFormat()
	keyword := InputData.GetString("keyword", "")
	sortStr := InputData.GetString("sort", "shelf_time")
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
	countList := make([]dy2.LiveRoomProductCount, 0)
	if len(list) > 0 {
		liveBusiness := business.NewLiveBusiness()
		authorBusiness := business.NewAuthorBusiness()
		curMap := map[string]dy2.LiveCurProductCount{}
		pmtMap := map[string][]dy2.LiveRoomProductSaleStatus{}
		curChan := make(chan map[string]dy2.LiveCurProductCount, 0)
		pmtChan := make(chan map[string][]dy2.LiveRoomProductSaleStatus, 0)
		authorChan := make(chan map[string]entity.DyAuthorData, 0)
		authorMap := map[string]entity.DyAuthorData{}
		for _, v := range list {
			go func(curCh chan map[string]dy2.LiveCurProductCount, pmtCh chan map[string][]dy2.LiveRoomProductSaleStatus, authorCh chan map[string]entity.DyAuthorData, roomId, productId, authorId string) {
				curMapTmp := liveBusiness.RoomCurProductByIds(roomId, []string{productId})
				pmtMapTmp := liveBusiness.RoomPmtProductByIds(roomId, []string{productId})
				authorData, _ := authorBusiness.HbaseGetAuthor(authorId)
				authorMapTmp := map[string]entity.DyAuthorData{}
				roomProductCurMap := map[string]dy2.LiveCurProductCount{}
				roomProductPmtMap := map[string][]dy2.LiveRoomProductSaleStatus{}
				if c, ok := curMapTmp[productId]; ok {
					roomProductCurMap[roomId] = c
				}
				if p, ok := pmtMapTmp[productId]; ok {
					roomProductPmtMap[roomId] = p
				}
				authorMapTmp[authorId] = authorData.Data
				curCh <- roomProductCurMap
				pmtCh <- roomProductPmtMap
				authorChan <- authorMapTmp
			}(curChan, pmtChan, authorChan, v.RoomID, v.ProductID, v.AuthorID)
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
			if v.RoomCover == "" {
				liveInfo, _ := hbase.GetLiveInfo(v.RoomID)
				v.RoomCover = dyimg.Fix(liveInfo.Cover)
			}
			item := dy2.LiveRoomProductCount{
				ProductInfo: v,
				ProductStartSale: dy2.RoomProductSaleChart{
					Timestamp: []int64{},
					Sales:     []int64{},
				},
				ProductEndSale: dy2.RoomProductSaleChart{
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
				item.ProductCur = dy2.LiveCurProductCount{
					CurList: []dy2.LiveCurProduct{},
				}
			}
			item.ProductInfo.AuthorID = business.IdEncrypt(item.ProductInfo.AuthorID)
			item.ProductInfo.ProductID = business.IdEncrypt(item.ProductInfo.ProductID)
			item.ProductInfo.AuthorRoomID = business.IdEncrypt(item.ProductInfo.AuthorRoomID)
			item.ProductInfo.RoomID = business.IdEncrypt(item.ProductInfo.RoomID)
			countList = append(countList, item)
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  countList,
		"total": total,
	})
	return
}

func (receiver *ProductController) ProductLiveAuthorAnalysis() {
	productId := business.IdDecrypt(receiver.Ctx.Input.Param(":product_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	tag := receiver.GetString("tag", "")
	minFollow, _ := receiver.GetInt64("min_follow", 0)
	maxFollow, _ := receiver.GetInt64("max_follow", 0)
	scoreType, _ := receiver.GetInt("score_type", 5)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	productBusiness := business.NewProductBusiness()
	list, total, comErr := productBusiness.ProductAuthorAnalysis(productId, keyword, tag, startTime, endTime, minFollow, maxFollow, scoreType, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		authorInfo, _ := hbase.GetAuthor(v.AuthorId)
		list[k].Avatar = dyimg.Fix(authorInfo.Data.Avatar)
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		list[k].ProductId = business.IdEncrypt(v.ProductId)
		list[k].Nickname = authorInfo.Data.Nickname
		list[k].RoomNum = len(v.RelatedRooms)
		list[k].RelatedRooms = []entity.DyProductAuthorRelatedRoom{}
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
	return
}

func (receiver *ProductController) ProductLiveAuthorAnalysisCount() {
	productId := business.IdDecrypt(receiver.Ctx.Input.Param(":product_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	productBusiness := business.NewProductBusiness()
	countList, comErr := productBusiness.ProductAuthorAnalysisCount(productId, keyword, startTime, endTime)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"list": countList,
	})
	return
}

func (receiver *ProductController) ProductAuthorLiveRooms() {
	productId := business.IdDecrypt(receiver.Ctx.Input.Param(":product_id"))
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 5, 10)
	sortStr := receiver.GetString("sort", "start_ts")
	orderBy := receiver.GetString("order_by", "desc")
	list, total := business.NewProductBusiness().ProductAuthorLiveRooms(productId, authorId, startTime, endTime, sortStr, orderBy, page, pageSize)
	for k, v := range list {
		list[k].Cover = dyimg.Fix(v.Cover)
		list[k].RoomId = business.IdEncrypt(v.RoomId)
		endLiveTime := v.EndTs
		if endLiveTime == 0 {
			endLiveTime = time.Now().Unix()
		}
		list[k].LiveSecond = endLiveTime - v.StartTs
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
}
