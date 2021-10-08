package dy

import (
	"dongchamao/business"
	"dongchamao/business/es"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	dy2 "dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"sort"
	"strings"
	"time"
)

type AwemeController struct {
	controllers.ApiBaseController
}

func (receiver *AwemeController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
}

func (receiver *AwemeController) AwemeBaseData() {
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	if awemeId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	awemeBase, comErr := hbase.GetVideo(awemeId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	awemeSimple := dy2.DySimpleAweme{
		AuthorID:        awemeBase.Data.AuthorID,
		AwemeCover:      awemeBase.Data.AwemeCover,
		AwemeTitle:      awemeBase.Data.AwemeTitle,
		AwemeCreateTime: awemeBase.Data.AwemeCreateTime,
		AwemeURL:        awemeBase.Data.AwemeURL,
		CommentCount:    awemeBase.Data.CommentCount,
		DiggCount:       awemeBase.Data.DiggCount,
		DownloadCount:   awemeBase.Data.DownloadCount,
		Duration:        awemeBase.Data.Duration,
		ForwardCount:    awemeBase.Data.ForwardCount,
		ID:              awemeBase.Data.ID,
		MusicID:         awemeBase.Data.MusicID,
		CrawlTime:       awemeBase.Data.CrawlTime,
		PromotionNum:    len(awemeBase.Data.DyPromotionID),
	}

	//昨天与前天的数据差
	now := time.Now()
	yesDateTime := now.AddDate(0, 0, -1)
	beforeDateTime := now.AddDate(0, 0, -2)

	yesData, beforeYesData := entity.DyAwemeDiggCommentForwardCount{}, entity.DyAwemeDiggCommentForwardCount{}
	yesData, _ = hbase.GetVideoCountData(awemeId, yesDateTime.Format("20060102"))
	beforeYesData, _ = hbase.GetVideoCountData(awemeId, beforeDateTime.Format("20060102"))
	awemeSimple.DiggInc = yesData.DiggCount - beforeYesData.DiggCount
	awemeSimple.CommentInc = yesData.CommentCount - beforeYesData.CommentCount
	awemeSimple.ForwardInc = yesData.ForwardCount - beforeYesData.ForwardCount

	//会员信息
	author, _ := hbase.GetAuthor(awemeBase.Data.AuthorID)
	if author.Data.UniqueID == "" {
		author.Data.UniqueID = author.Data.ShortID
	}
	author.Data.Avatar = dyimg.Fix(author.Data.Avatar)
	receiver.SuccReturn(map[string]interface{}{
		"aweme_base": awemeSimple,
		"author": map[string]interface{}{
			"follower_count": author.FollowerCount,
			"nick_name":      author.Data.Nickname,
			"avatar":         author.Data.Avatar,
		},
	})
	return
}

func (receiver *AwemeController) AwemeChart() {
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	if awemeId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	t1, t2, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	awemeBusiness := business.NewAwemeBusiness()
	awemeCount, comErr := awemeBusiness.GetAwemeChart(awemeId, t1, t2, true)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	//前一天数据，做增量计算
	beginDatetime := t1
	beforeData := entity.DyAwemeDiggCommentForwardCount{}
	beforeDay := beginDatetime.AddDate(0, 0, -1).Format("20060102")
	if _, ok := awemeCount[beforeDay]; ok {
		beforeData = awemeCount[beforeDay]
	} else {
		beforeData, _ = hbase.GetVideoCountData(awemeId, beforeDay)
	}
	dateArr := make([]string, 0)
	diggCountArr := make([]int64, 0)
	commentCountArr := make([]int64, 0)
	forwardCountArr := make([]int64, 0)
	diggIncArr := make([]int64, 0)
	commentIncArr := make([]int64, 0)
	forwardIncArr := make([]int64, 0)
	for {
		if beginDatetime.After(t2) {
			break
		}
		date := beginDatetime.Format("20060102")
		if _, ok := awemeCount[date]; !ok {
			awemeCount[date] = beforeData
		}
		currentData := awemeCount[date]
		dateArr = append(dateArr, beginDatetime.Format("01/02"))
		diggCountArr = append(diggCountArr, currentData.DiggCount)
		commentCountArr = append(commentCountArr, currentData.CommentCount)
		forwardCountArr = append(forwardCountArr, currentData.ForwardCount)
		diggIncArr = append(diggIncArr, currentData.DiggCount-beforeData.DiggCount)
		commentIncArr = append(commentIncArr, currentData.CommentCount-beforeData.CommentCount)
		forwardIncArr = append(forwardIncArr, currentData.ForwardCount-beforeData.ForwardCount)
		beforeData = currentData
		beginDatetime = beginDatetime.AddDate(0, 0, 1)
	}
	returnMap := map[string]interface{}{
		"digg": dy2.DateChart{
			Date:       dateArr,
			CountValue: diggCountArr,
			IncValue:   diggIncArr,
		},
		"forward": dy2.DateChart{
			Date:       dateArr,
			CountValue: forwardCountArr,
			IncValue:   forwardIncArr,
		},
		"comment": dy2.DateChart{
			Date:       dateArr,
			CountValue: commentCountArr,
			IncValue:   commentIncArr,
		},
	}
	receiver.SuccReturn(returnMap)
	return
}

func (receiver *AwemeController) AwemeCommentHotWords() {
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	if awemeId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	awemeBase, comErr := hbase.GetVideo(awemeId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	list := make([]dy2.NameValueInt64Chart, 0)
	for k, v := range awemeBase.HotWordShow {
		list = append(list, dy2.NameValueInt64Chart{
			Name:  k,
			Value: v,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Value > list[j].Value
	})
	receiver.SuccReturn(map[string]interface{}{
		"hot_words":  list,
		"use_id_num": awemeBase.ContextNum["use_id_num"],
		"msg_id_num": awemeBase.ContextNum["msg_id_num"],
	})
	return
}

func (receiver *AwemeController) AwemeCommentTop() {
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	if awemeId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	start := (page - 1) * pageSize
	end := page * pageSize
	awemeComment, total, comErr := hbase.GetAwemeTopComment(awemeId, start, end)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	if total > 1000 {
		total = 1000
	}
	receiver.SuccReturn(map[string]interface{}{
		"total": total,
		"page":  page,
		"size":  pageSize,
		"list":  awemeComment,
	})
	return
}

//视频商品数据
func (receiver *AwemeController) AwemeProductAnalyse() {
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 5, 10)
	cacheKey := cache.GetCacheKey(cache.AwemeProductByDate, awemeId, startTime.Format("20060102"), endTime.Format("20060102"))
	cacheData := global.Cache.Get(cacheKey)
	var hbaseList = make([]entity.DyProductAwemeDailyDistribute, 0)
	if cacheData != "" {
		cacheData = utils.DeserializeData(cacheData)
		_ = jsoniter.Unmarshal([]byte(cacheData), &hbaseList)
	} else {
		hbaseList, _ = hbase.GetDyProductAwemeDailyDistributeRange(awemeId, startTime.Format("20060102"), endTime.Format("20060102"))
		_ = global.Cache.Set(cacheKey, utils.SerializeData(hbaseList), 180)
	}
	productMap := map[string]dy2.DyAwemeProductSale{}
	for _, v := range hbaseList {
		if s, ok := productMap[v.ProductId]; ok {
			s.Sales += v.Sales
			s.Gmv += v.AwemeGmv
			productMap[v.ProductId] = s
		} else {
			productInfo, _ := hbase.GetProductInfo(v.ProductId)
			productMap[v.ProductId] = dy2.DyAwemeProductSale{
				AwemeId:       v.AwemeId,
				ProductId:     v.ProductId,
				Gmv:           v.AwemeGmv,
				Sales:         v.Sales,
				Price:         v.Price,
				Title:         productInfo.Title,
				PlatformLabel: productInfo.PlatformLabel,
				ProductStatus: productInfo.Status,
				CouponInfo:    productInfo.TbCouponInfo,
				Image:         productInfo.Image,
			}
		}
	}
	list := make([]dy2.DyAwemeProductSale, 0)
	lenNum := len(productMap)
	if lenNum == 0 {
		receiver.getAwemeProducts(awemeId, page, pageSize)
		return
	}
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > lenNum {
		receiver.SuccReturn(map[string]interface{}{
			"list":  list,
			"total": lenNum,
		})
		return
	}
	for _, v := range productMap {
		v.AwemeId = business.IdEncrypt(v.AwemeId)
		v.ProductId = business.IdEncrypt(v.ProductId)
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Gmv < list[j].Gmv
	})
	if end > lenNum {
		end = lenNum
	}
	list = list[start:end]
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": lenNum,
	})
	return
}

func (receiver *AwemeController) getAwemeProducts(awemeId string, page, pageSize int) {
	info, err := hbase.GetVideo(awemeId)
	if err != nil {
		receiver.FailReturn(err)
		return
	}
	dateTime := time.Unix(info.Data.AwemeCreateTime, 0)
	awemeInfo, _ := es.NewEsVideoBusiness().GetByAwemeId(awemeId, dateTime.Format("20060102"))
	list := []dy2.DyAwemeProductSale{}
	if awemeInfo.ProductIds == "" {
		receiver.SuccReturn(map[string]interface{}{
			"list":  list,
			"total": 0,
		})
		return
	}
	productIds := strings.Split(awemeInfo.ProductIds, ",")
	total := len(productIds)
	if total > 0 {
		products, _ := hbase.GetProductByIds(productIds)
		awemeEncryptId := business.IdEncrypt(awemeId)
		for _, v := range products {
			list = append(list, dy2.DyAwemeProductSale{
				AwemeId:       awemeEncryptId,
				ProductId:     business.IdEncrypt(v.ProductID),
				Gmv:           0,
				Sales:         0,
				Price:         v.Price,
				Title:         v.Title,
				PlatformLabel: v.PlatformLabel,
				ProductStatus: v.Status,
				CouponInfo:    v.TbCouponInfo,
				Image:         dyimg.Fix(v.Image),
			})
		}
		start := (page - 1) * pageSize
		end := start + pageSize
		if start > total {
			receiver.SuccReturn(map[string]interface{}{
				"list":  list,
				"total": total,
			})
			return
		}
		if end > total {
			end = total
		}
		list = list[start:end]
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
	return
}

func (receiver *AwemeController) AwemeProductAnalyseChart() {
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	cacheKey := cache.GetCacheKey(cache.AwemeProductByDate, awemeId, startTime.Format("20060102"), endTime.Format("20060102"))
	cacheData := global.Cache.Get(cacheKey)
	var hbaseList = make([]entity.DyProductAwemeDailyDistribute, 0)
	if cacheData != "" {
		cacheData = utils.DeserializeData(cacheData)
		_ = jsoniter.Unmarshal([]byte(cacheData), &hbaseList)
	} else {
		hbaseList, _ = hbase.GetDyProductAwemeDailyDistributeRange(awemeId, startTime.Format("20060102"), endTime.Format("20060102"))
		_ = global.Cache.Set(cacheKey, utils.SerializeData(hbaseList), 180)
	}
	var allGmv float64 = 0
	var allSales int64 = 0
	var productMap = map[string]string{}
	dateMap := map[string]dy2.DyAwemeProductSale{}
	dateProductsMap := map[string]map[string]string{}
	for _, v := range hbaseList {
		if _, ok := productMap[v.ProductId]; !ok {
			productMap[v.ProductId] = v.ProductId
		}
		if _, ok := dateProductsMap[v.DistDate]; !ok {
			dateProductsMap[v.DistDate] = map[string]string{}
		}
		dateProductsMap[v.DistDate][v.ProductId] = v.ProductId
		allGmv += v.AwemeGmv
		allSales += v.Sales
		if s, ok := dateMap[v.DistDate]; ok {
			s.Sales += v.Sales
			s.Gmv += v.AwemeGmv
			dateMap[v.DistDate] = s
		} else {
			dateMap[v.DistDate] = dy2.DyAwemeProductSale{
				Gmv:   v.AwemeGmv,
				Sales: v.Sales,
			}
		}
	}
	infoMap := map[string]entity.DyProduct{}
	for k := range productMap {
		productInfo, _ := hbase.GetProductInfo(k)
		infoMap[k] = productInfo
	}
	list := make([]dy2.NameValueInt64ChartWithData, 0)
	dateTime := startTime
	for {
		if dateTime.After(endTime) {
			break
		}
		parseTime := dateTime.Format("20060102")
		valueTime, _ := time.ParseInLocation("20060102", parseTime, time.Local)
		var dateData = dy2.NameValueInt64ChartWithData{Name: valueTime.Format("01/02"), Value: 0, Data: []string{}, Date: valueTime}
		if dateMap[parseTime].Sales != 0 {
			data := make([]string, 0)
			if p, ok := dateProductsMap[parseTime]; ok {
				for k1 := range p {
					if _, ok1 := infoMap[k1]; ok1 {
						title := infoMap[k1].Title
						data = append(data, title)
					}
				}
			}
			dateData.Value = dateMap[parseTime].Sales
			dateData.Data = data
		}
		list = append(list, dateData)
		dateTime = dateTime.AddDate(0, 0, 1)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Date.Before(list[j].Date)
	})
	receiver.SuccReturn(map[string]interface{}{
		"count": map[string]interface{}{
			"gmv":         allGmv,
			"sales":       allSales,
			"product_num": len(productMap),
		},
		"list": list,
	})
	return
}

//商品带货同款视频
func (receiver *AwemeController) AwemeProductSameAweme() {
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	productId := business.IdDecrypt(receiver.Ctx.Input.Param(":product_id"))
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 5, 30)
	keyword := receiver.GetString("keyword", "")
	orderBy := receiver.GetString("order_by", "")
	sortStr := receiver.GetString("sort", "")
	list, total, comErr := es.NewEsVideoBusiness().SearchByProductId(productId, awemeId, keyword, sortStr, orderBy, page, pageSize, startTime, endTime)
	for k, v := range list {
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		list[k].AwemeId = business.IdEncrypt(v.AwemeId)
		list[k].AwemeCover = dyimg.Fix(v.AwemeCover)
		list[k].Avatar = dyimg.Fix(v.Avatar)
		if v.UniqueId == "" {
			list[k].UniqueId = v.ShortId
		}
		list[k].AwemeUrl = business.AwemeUrl + v.AwemeId
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
	return
}

//直播粉丝分析
func (receiver *AwemeController) AwemeFanAnalyse() {
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	info, comErr := hbase.GetVideo(awemeId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	var genderTotal int64 = 0
	var ageTotal int64 = 0
	var cityTotal int64 = 0
	var provinceTotal int64 = 0
	genderChart := make([]entity.XtDistributionsList, 0)
	genderMap := make(map[string]entity.XtDistributionsList, 0)
	ageChart := make([]entity.XtDistributionsList, 0)
	ageMap := make(map[string]entity.XtDistributionsList, 0)
	cityChart := make([]entity.XtDistributionsList, 0)
	cityMap := make(map[string]entity.XtDistributionsList, 0)
	provinceChart := make([]entity.XtDistributionsList, 0)
	provinceMap := make(map[string]entity.XtDistributionsList, 0)
	var genderFormat = map[string]string{"男": "male", "女": "female"}
	for _, v := range info.Gender {
		gender := genderFormat[v.Gender]
		if gender == "" {
			continue
		}
		genderNum := utils.ToInt64(v.GenderNum)
		genderData := genderMap[gender]
		genderTotal += genderNum
		genderData.DistributionKey = gender
		genderData.DistributionValue = genderNum
		genderMap[gender] = genderData
	}
	for _, v := range genderMap {
		genderChart = append(genderChart, v)
	}
	for _, v := range info.AgeDistrinbution {
		if v.AgeDistrinbution == "" {
			continue
		}
		ageDistributionNum := utils.ToInt64(v.AgeDistrinbutionNum)
		ageTotal += ageDistributionNum
		age := ageMap[v.AgeDistrinbution]
		age.DistributionKey = v.AgeDistrinbution
		age.DistributionValue += ageDistributionNum
		ageMap[v.AgeDistrinbution] = age
	}
	for _, v := range info.City {
		if v.City == "" {
			continue
		}
		v.City = strings.Replace(v.City, "市", "", -1)
		cityNum := utils.ToInt64(v.CityNum)
		cityTotal += cityNum
		city := cityMap[v.City]
		city.DistributionKey = v.City
		city.DistributionValue += cityNum
		cityMap[v.City] = city
	}
	for _, v := range ageMap {
		ageChart = append(ageChart, v)
	}
	for _, v := range cityMap {
		cityChart = append(cityChart, v)
	}
	for _, v := range info.Province {
		if v.Province == "" {
			continue
		}
		v.Province = strings.Replace(v.Province, "省", "", -1)
		distributionValue := utils.ToInt64(v.ProvinceNum)
		provinceTotal += distributionValue
		province := provinceMap[v.Province]
		province.DistributionKey = v.Province
		province.DistributionValue += distributionValue
		provinceMap[v.Province] = province
	}
	for _, v := range provinceMap {
		provinceChart = append(provinceChart, v)
	}
	sort.Slice(cityChart, func(i, j int) bool {
		return cityChart[i].DistributionValue > cityChart[j].DistributionValue
	})
	sort.Slice(provinceChart, func(i, j int) bool {
		return provinceChart[i].DistributionValue > provinceChart[j].DistributionValue
	})
	if genderTotal > 0 {
		for k, v := range genderChart {
			genderChart[k].DistributionPer = float64(v.DistributionValue) / float64(genderTotal)
		}
	}
	if ageTotal > 0 {
		for k, v := range ageChart {
			ageChart[k].DistributionPer = float64(v.DistributionValue) / float64(ageTotal)
		}
	}
	if cityTotal > 0 {
		for k, v := range cityChart {
			cityChart[k].DistributionPer = float64(v.DistributionValue) / float64(cityTotal)
		}
	}
	if provinceTotal > 0 {
		for k, v := range provinceChart {
			provinceChart[k].DistributionPer = float64(v.DistributionValue) / float64(provinceTotal)
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"age_people":      ageTotal,
		"age_chart":       ageChart,
		"gender_chart":    genderChart,
		"gender_total":    genderTotal,
		"city_chart":      cityChart,
		"city_total":      cityTotal,
		"province_chart":  provinceChart,
		"province_people": provinceTotal,
	})
	return
}

//视频加速
func (receiver *AwemeController) AwemeSpeed() {

	if !business.UserActionLock(receiver.TrueUri, utils.ToString(receiver.UserId), 5) {
		receiver.FailReturn(global.NewError(6000))
		return
	}
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	if awemeId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	_, comErr := hbase.GetVideo(awemeId)
	if comErr != nil {
		receiver.FailReturn(global.NewError(4000))
	}

	spriderName := "aweme"
	cacheKey := cache.GetCacheKey(cache.SpiderSpeedUpLimit, spriderName, awemeId)
	cacheData := global.Cache.Get(cacheKey)
	if cacheData != "" {
		//缓存存在
		receiver.FailReturn(global.NewError(6000))
		return
	}
	//加速
	ret, _ := business.NewSpiderBusiness().SpiderSpeedUp(spriderName, awemeId)
	global.Cache.Set(cacheKey, "1", 300)

	logs.Info("视频加速，爬虫推送结果：", ret)
	receiver.SuccReturn([]string{})
	return
}
