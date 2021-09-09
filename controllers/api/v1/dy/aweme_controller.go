package dy

import (
	"dongchamao/business"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	dy2 "dongchamao/models/repost/dy"
	jsoniter "github.com/json-iterator/go"
	"sort"
	"time"
)

type AwemeController struct {
	controllers.ApiBaseController
}

func (receiver *AwemeController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseShowNum)
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
		ShareCount:      awemeBase.Data.ShareCount,
		PromotionNum:    len(awemeBase.Data.DyPromotionID),
	}
	receiver.SuccReturn(map[string]interface{}{
		"aweme_base": awemeSimple,
	})
	return
}

func (receiver *AwemeController) AwemeChart() {
	awemeId := business.IdEncrypt(receiver.Ctx.Input.Param(":aweme_id"))
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
		"hot_words": list,
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
	for k, v := range dateMap {
		data := make([]string, 0)
		if p, ok := dateProductsMap[k]; ok {
			for k1 := range p {
				if _, ok1 := infoMap[k1]; ok1 {
					title := infoMap[k1].Title
					data = append(data, title)
				}
			}
		}
		valueTime, _ := time.ParseInLocation("20060102", k, time.Local)
		list = append(list, dy2.NameValueInt64ChartWithData{
			Name:  valueTime.Format("01/02"),
			Value: v.Sales,
			Data:  data,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[i].Name
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
