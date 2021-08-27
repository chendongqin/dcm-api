package business

import (
	"dongchamao/business/es"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	"dongchamao/models/entity"
	"dongchamao/models/repost/dy"
	"dongchamao/services"
	"dongchamao/services/dyimg"
	jsoniter "github.com/json-iterator/go"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

const ShareUrlPrefix = "https://www.iesdouyin.com/share/user/"

type AuthorBusiness struct {
}

func NewAuthorBusiness() *AuthorBusiness {
	return new(AuthorBusiness)
}

func (a *AuthorBusiness) GetCacheAuthorLiveTags(enableCache bool) []string {
	cacheKey := cache.GetCacheKey(cache.LongTimeConfigKeyCache)
	redisService := services.NewRedisService()
	tags := make([]string, 0)
	if enableCache == true {
		jsonStr := redisService.Hget(cacheKey, "author_live_tags")
		if jsonStr != "" {
			jsonStr = utils.DeserializeData(jsonStr)
			if jsonStr != "" {
				_ = jsoniter.Unmarshal([]byte(jsonStr), &tags)
				return tags
			}
		}
	}
	list := make([]dcm.DyAuthorLiveTags, 0)
	_ = dcm.GetSlaveDbSession().Desc("weight").Find(&list)
	for _, v := range list {
		tags = append(tags, v.Name)
	}
	if len(tags) > 0 {
		_ = redisService.Hset(cacheKey, "author_live_tags", utils.SerializeData(tags))
	}
	return tags
}

//粉丝｜粉丝团趋势数据
func (a *AuthorBusiness) HbaseGetFansRangDate(authorId, startDate, endDate string) (data map[string]dy.DateChart, comErr global.CommonError) {
	data = map[string]dy.DateChart{}
	dateMap, comErr := hbase.GetFansRangDate(authorId, startDate, endDate)
	if comErr != nil {
		return
	}
	//起点补点操作
	t1, _ := time.ParseInLocation("20060102", startDate, time.Local)
	t2, _ := time.ParseInLocation("20060102", endDate, time.Local)
	beforeDate := t1.AddDate(0, 0, -1).Format("20060102")
	beforeData, _ := hbase.GetFansByDate(authorId, beforeDate)
	if _, ok := dateMap[startDate]; !ok {
		dateMap[startDate] = beforeData
	}
	//末点补点
	if endDate == time.Now().Format("20060102") {
		if _, ok := dateMap[endDate]; !ok {
			endData, _ := hbase.GetAuthor(authorId)
			dateMap[endDate] = entity.DyAuthorFans{
				FollowerCount:       endData.FollowerCount,
				TotalFansGroupCount: endData.TotalFansCount,
			}
		}
	}
	countGroupArr := make([]int64, 0)
	incGroupArr := make([]int64, 0)
	dateArr := make([]string, 0)
	countArr := make([]int64, 0)
	incArr := make([]int64, 0)
	beginDatetime := t1
	for {
		if beginDatetime.After(t2) {
			break
		}
		date := beginDatetime.Format("20060102")
		if _, ok := dateMap[date]; !ok {
			yesterday := beginDatetime.AddDate(0, 0, -1).Format("20060102")
			if _, ok := dateMap[yesterday]; ok {
				dateMap[date] = dateMap[yesterday]
			} else {
				dateMap[date] = entity.DyAuthorFans{
					FollowerCount:       0,
					TotalFansGroupCount: 0,
				}
			}
		}
		fansData := dateMap[date]
		inc1 := fansData.FollowerCount - beforeData.FollowerCount
		inc2 := fansData.TotalFansGroupCount - beforeData.TotalFansGroupCount
		dateArr = append(dateArr, beginDatetime.Format("01/02"))
		incArr = append(incArr, inc1)
		incGroupArr = append(incGroupArr, inc2)
		countArr = append(countArr, fansData.FollowerCount)
		countGroupArr = append(countGroupArr, fansData.TotalFansGroupCount)
		beforeData = fansData
		beginDatetime = beginDatetime.AddDate(0, 0, 1)
	}
	data["fans"] = dy.DateChart{
		Date:       dateArr,
		CountValue: countArr,
		IncValue:   incArr,
	}
	data["fans_group"] = dy.DateChart{
		Date:       dateArr,
		CountValue: countGroupArr,
		IncValue:   incGroupArr,
	}
	return
}

//达人数据
func (a *AuthorBusiness) HbaseGetAuthor(authorId string) (data entity.DyAuthor, comErr global.CommonError) {
	data, comErr = hbase.GetAuthor(authorId)
	if comErr != nil {
		return
	}
	data.Data.Age = GetAge(data.Data.Birthday)
	//数据做同步
	data.Data.RoomID = data.RoomId
	data.Data.Avatar = dyimg.Fix(data.Data.Avatar)
	data.Data.ShareUrl = ShareUrlPrefix + data.Data.ID
	if data.Data.UniqueID == "" {
		data.Data.UniqueID = data.Data.ShortID
	}
	return
}

//达人基础数据趋势
func (a *AuthorBusiness) HbaseGetAuthorBasicRangeDate(authorId string, startTime, endTime time.Time) (data map[string]dy.DateChart, comErr global.CommonError) {
	data = map[string]dy.DateChart{}
	dateMap, comErr := hbase.GetAuthorBasicRangeDate(authorId, startTime, endTime)
	//起点补点操作
	startDate := startTime.Format("20060102")
	endDate := endTime.Format("20060102")
	beforeDate := startTime.AddDate(0, 0, -1).Format("20060102")
	beforeBasicData, _ := hbase.GetAuthorBasic(authorId, beforeDate)
	beforeData := dy.DyAuthorBasicChart{
		FollowerCount:  beforeBasicData.FollowerCount,
		TotalFansCount: beforeBasicData.TotalFansCount,
		TotalFavorited: beforeBasicData.TotalFavorited,
		CommentCount:   beforeBasicData.CommentCount,
		ForwardCount:   beforeBasicData.ForwardCount,
	}
	if _, ok := dateMap[startDate]; !ok {
		dateMap[startDate] = beforeData
	}
	//末点补点
	if endDate == time.Now().Format("20060102") {
		if _, ok := dateMap[endDate]; !ok {
			authorData, _ := hbase.GetAuthor(authorId)
			dateMap[endDate] = dy.DyAuthorBasicChart{
				FollowerCount:  authorData.FollowerCount,
				TotalFansCount: authorData.TotalFansCount,
				TotalFavorited: authorData.TotalFavorited,
				CommentCount:   authorData.CommentCount,
				ForwardCount:   authorData.ForwardCount,
			}
		}
	}
	dateArr := make([]string, 0)
	fansCountArr := make([]int64, 0)
	fanIncArr := make([]int64, 0)
	fansGroupCountArr := make([]int64, 0)
	fansGroupIncArr := make([]int64, 0)
	diggCountArr := make([]int64, 0)
	diggIncArr := make([]int64, 0)
	commentCountArr := make([]int64, 0)
	commentIncArr := make([]int64, 0)
	forwardCountArr := make([]int64, 0)
	forwardIncArr := make([]int64, 0)
	beginDatetime := startTime
	for {
		if beginDatetime.After(endTime) {
			break
		}
		date := beginDatetime.Format("20060102")
		if _, ok := dateMap[date]; !ok {
			dateMap[date] = beforeData
		}
		currentData := dateMap[date]
		dateArr = append(dateArr, beginDatetime.Format("01/02"))
		fansCountArr = append(fansCountArr, currentData.FollowerCount)
		fanIncArr = append(fanIncArr, currentData.FollowerCount-beforeData.FollowerCount)
		fansGroupCountArr = append(fansGroupCountArr, currentData.TotalFansCount)
		fansGroupIncArr = append(fansGroupIncArr, currentData.TotalFansCount-beforeData.TotalFansCount)
		diggCountArr = append(diggCountArr, currentData.TotalFavorited)
		diggIncArr = append(diggIncArr, currentData.TotalFavorited-beforeData.TotalFavorited)
		commentCountArr = append(commentCountArr, currentData.CommentCount)
		commentIncArr = append(commentIncArr, currentData.CommentCount-beforeData.CommentCount)
		forwardCountArr = append(forwardCountArr, currentData.ForwardCount)
		forwardIncArr = append(forwardIncArr, currentData.ForwardCount-beforeData.ForwardCount)
		beforeData = currentData
		beginDatetime = beginDatetime.AddDate(0, 0, 1)
	}
	data["fans"] = dy.DateChart{
		Date:       dateArr,
		CountValue: fansCountArr,
		IncValue:   fanIncArr,
	}
	data["fans_club"] = dy.DateChart{
		Date:       dateArr,
		CountValue: fansGroupCountArr,
		IncValue:   fansGroupIncArr,
	}
	data["digg"] = dy.DateChart{
		Date:       dateArr,
		CountValue: diggCountArr,
		IncValue:   diggIncArr,
	}
	data["forward"] = dy.DateChart{
		Date:       dateArr,
		CountValue: forwardCountArr,
		IncValue:   forwardIncArr,
	}
	data["comment"] = dy.DateChart{
		Date:       dateArr,
		CountValue: commentCountArr,
		IncValue:   commentIncArr,
	}
	return
}

//达人（带货）口碑
func (a *AuthorBusiness) HbaseGetAuthorReputation(authorId string) (data entity.DyReputation, comErr global.CommonError) {
	data, comErr = hbase.GetAuthorReputation(authorId)
	if len(data.ScoreList) == 0 {
		data.ScoreList = make([]entity.DyReputationMonthScoreList, 0)
	}
	if len(data.DtScoreList) == 0 {
		data.DtScoreList = make([]entity.DyReputationDateScoreList, 0)
	} else {
		data.DtScoreList = ReputationDtScoreListOrderByTime(data.DtScoreList)
		for k, v := range data.DtScoreList {
			data.DtScoreList[k].DateStr = utils.ToString(v.Date)
		}
	}
	data.UID = authorId
	return
}

//获取达人星图评分
func (a *AuthorBusiness) GetDyAuthorScore(liveScore entity.XtAuthorLiveScore, videoScore entity.XtAuthorScore) (authorStarScores dy.DyAuthorStarScores) {
	authorStarScores.LiveScore = dy.DyAuthorStarScore{
		CooperateIndex: utils.FriendlyFloat64(float64(liveScore.CooperateIndex) / 10000),
		CpIndex:        utils.FriendlyFloat64(float64(liveScore.CpIndex) / 10000),
		GrowthIndex:    utils.FriendlyFloat64(float64(liveScore.GrowthIndex) / 10000),
		ShoppingIndex:  utils.FriendlyFloat64(float64(liveScore.ShoppingIndex) / 10000),
		SpreadIndex:    utils.FriendlyFloat64(float64(liveScore.SpreadIndex) / 10000),
		TopScore:       utils.FriendlyFloat64(float64(liveScore.TopScore) / 10000),
	}
	authorStarScores.LiveScore.AvgLevel = dy.XtAuthorScoreAvgLevel{
		CooperateIndex: utils.FriendlyFloat64(float64(liveScore.AvgLevel.CooperateIndex) / 100),
		CpIndex:        utils.FriendlyFloat64(float64(liveScore.AvgLevel.CpIndex) / 100),
		GrowthIndex:    utils.FriendlyFloat64(float64(liveScore.AvgLevel.GrowthIndex) / 100),
		ShoppingIndex:  utils.FriendlyFloat64(float64(liveScore.AvgLevel.ShoppingIndex) / 100),
		SpreadIndex:    utils.FriendlyFloat64(float64(liveScore.AvgLevel.SpreadIndex) / 100),
		TopScore:       utils.FriendlyFloat64(float64(liveScore.AvgLevel.TopScore) / 100),
	}
	authorStarScores.VideoScore = dy.DyAuthorStarScore{
		CooperateIndex: utils.FriendlyFloat64(float64(videoScore.CooperateIndex) / 10000),
		CpIndex:        utils.FriendlyFloat64(float64(videoScore.CpIndex) / 10000),
		GrowthIndex:    utils.FriendlyFloat64(float64(videoScore.GrowthIndex) / 10000),
		ShoppingIndex:  utils.FriendlyFloat64(float64(videoScore.ShoppingIndex) / 10000),
		SpreadIndex:    utils.FriendlyFloat64(float64(videoScore.SpreadIndex) / 10000),
		TopScore:       utils.FriendlyFloat64(float64(videoScore.TopScore) / 10000),
	}
	authorStarScores.VideoScore.AvgLevel = dy.XtAuthorScoreAvgLevel{
		CooperateIndex: utils.FriendlyFloat64(float64(videoScore.AvgLevel.CooperateIndex) / 100),
		CpIndex:        utils.FriendlyFloat64(float64(videoScore.AvgLevel.CpIndex) / 100),
		GrowthIndex:    utils.FriendlyFloat64(float64(videoScore.AvgLevel.GrowthIndex) / 100),
		ShoppingIndex:  utils.FriendlyFloat64(float64(videoScore.AvgLevel.ShoppingIndex) / 100),
		SpreadIndex:    utils.FriendlyFloat64(float64(videoScore.AvgLevel.SpreadIndex) / 100),
		TopScore:       utils.FriendlyFloat64(float64(videoScore.AvgLevel.TopScore) / 100),
	}
	return
}

//直播分析
func (a *AuthorBusiness) CountLiveRoomAnalyse(authorId string, startTime, endTime time.Time) (data dy.SumDyLiveRoom) {
	data = dy.SumDyLiveRoom{
		UserTotalChart: dy.DyUserTotalChart{
			Date:       []string{},
			CountValue: []int64{},
			Rooms:      [][]dy.DyLiveRoomChart{},
		},
		OnlineTimeChart: dy.DateCountFChart{
			Date:       []string{},
			CountValue: []float64{},
		},
		UvChart: dy.DateCountFChart{
			Date:       []string{},
			CountValue: []float64{},
		},
		AmountChart: dy.DateCountFChart{
			Date:       []string{},
			CountValue: []float64{},
		},
		LiveLongTimeChart:  []dy.NameValueChart{},
		LiveStartHourChart: []dy.NameValueChart{},
	}
	roomsMap, _ := hbase.GetAuthorRoomsRangDate(authorId, startTime, endTime)
	liveDataList := make([]dy.DyLiveRoomAnalyse, 0)
	roomNum := 0
	for _, rooms := range roomsMap {
		roomNum += len(rooms)
	}
	if roomNum == 0 {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(roomNum)
	hbaseDataChan := make(chan dy.DyLiveRoomAnalyse, roomNum)
	for _, rooms := range roomsMap {
		for _, room := range rooms {
			go func(roomId string) {
				defer global.RecoverPanic()
				defer wg.Done()
				liveBusiness := NewLiveBusiness()
				roomAnalyse, comErr := liveBusiness.LiveRoomAnalyse(roomId)
				if comErr == nil {
					hbaseDataChan <- roomAnalyse
				}
			}(room.RoomID)
		}
	}
	wg.Wait()
	for i := 0; i < roomNum; i++ {
		roomAnalyse, ok := <-hbaseDataChan
		if !ok {
			break
		}
		liveDataList = append(liveDataList, roomAnalyse)
	}
	sumData := map[string]dy.DyLiveRoomAnalyse{}
	sumLongTime := map[string]int{}
	sumHourTime := map[string]int{}
	dateRoomMap := map[string][]dy.DyLiveRoomChart{}
	for _, v := range liveDataList {
		date := time.Unix(v.DiscoverTime, 0).Format("01/02")
		longStr := ""
		HourStr := time.Unix(v.LiveStartTime, 0).Format("15")
		if v.LiveLongTime > 4*3600 {
			longStr = "up_h4"
		} else if v.LiveLongTime > 2*3600 {
			longStr = "h2_h4"
		} else if v.LiveLongTime > 3600 {
			longStr = "h1_h2"
		} else if v.LiveLongTime > 1800 {
			longStr = "m30_h1"
		} else {
			longStr = "down_m30"
		}
		if _, ok := sumHourTime[HourStr]; ok {
			sumHourTime[HourStr] += 1
		} else {
			sumHourTime[HourStr] = 1
		}
		if _, ok := sumLongTime[longStr]; ok {
			sumLongTime[longStr] += 1
		} else {
			sumLongTime[longStr] = 1
		}
		if _, ok := dateRoomMap[date]; !ok {
			dateRoomMap[date] = []dy.DyLiveRoomChart{}
		}
		dateRoomMap[date] = append(dateRoomMap[date], dy.DyLiveRoomChart{
			RoomId:    v.RoomId,
			Title:     v.Title,
			UserTotal: v.TotalUserCount,
		})
		if d, ex := sumData[date]; ex {
			d.TotalUserCount += v.TotalUserCount
			d.IncFans += v.IncFans
			if d.TotalUserCount > 0 {
				d.IncFansRate = float64(d.IncFans) / float64(d.TotalUserCount)
				d.InteractRate = float64(d.BarrageCount) / float64(d.TotalUserCount)
			}
			d.BarrageCount += v.BarrageCount
			avgUserCount := (d.AvgUserCount + v.AvgUserCount) / 2
			d.AvgUserCount = avgUserCount
			d.Volume += v.Volume
			d.Amount += v.Amount
			uv := (d.Uv + v.Uv) / 2
			d.Uv = uv
			saleRate := (d.SaleRate + v.SaleRate) / 2
			d.SaleRate = saleRate
			perPrice := (d.PerPrice + v.PerPrice) / 2
			d.PerPrice = perPrice
			d.LiveLongTime += v.LiveLongTime
			d.LiveStartTime = v.LiveStartTime
			avgOnlineTime := (d.AvgOnlineTime + v.AvgOnlineTime) / 2
			d.AvgOnlineTime = avgOnlineTime
			if d.PromotionNum == 0 {
				d.PromotionNum = v.PromotionNum
			}
			sumData[date] = d
		} else {
			sumData[date] = v
		}
	}
	keys := make([]string, 0)
	for k, _ := range sumData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	data.LiveStartHourChart = make([]dy.NameValueChart, 0)
	data.LiveLongTimeChart = make([]dy.NameValueChart, 0)
	for k, v := range sumLongTime {
		data.LiveLongTimeChart = append(data.LiveLongTimeChart, dy.NameValueChart{
			Name:  k,
			Value: v,
		})
	}
	for k, v := range sumHourTime {
		data.LiveStartHourChart = append(data.LiveStartHourChart, dy.NameValueChart{
			Name:  k,
			Value: v,
		})
	}
	dateChart := make([]string, 0)
	userTotalChart := make([]int64, 0)
	onlineUserChart := make([]float64, 0)
	roomChart := make([][]dy.DyLiveRoomChart, 0)
	uvChart := make([]float64, 0)
	amountChart := make([]float64, 0)
	for _, date := range keys {
		v := sumData[date]
		data.UserData.LiveNum += 1
		data.UserData.AvgUserTotal += v.TotalUserCount
		data.UserData.AvgUserCount += v.AvgUserCount
		data.UserData.AvgInteractRate += v.InteractRate
		data.UserData.IncFans += v.IncFans
		data.UserData.AvgIncFansRate += v.IncFansRate
		data.SaleData.AvgVolume += v.Volume
		data.SaleData.AvgAmount += v.Amount
		data.SaleData.AvgUv += v.Uv
		data.SaleData.AvgPerPrice += v.PerPrice
		dateChart = append(dateChart, date)
		userTotalChart = append(userTotalChart, v.TotalUserCount)
		onlineUserChart = append(onlineUserChart, v.AvgOnlineTime)
		uvChart = append(uvChart, v.Uv)
		amountChart = append(amountChart, v.Amount)
		if v.PromotionNum > 0 {
			data.UserData.PromotionLiveNum += 1
			data.SaleData.PromotionNum += v.PromotionNum
		}
		if rv, ok := dateRoomMap[date]; ok {
			roomChart = append(roomChart, rv)
		} else {
			roomChart = append(roomChart, []dy.DyLiveRoomChart{})
		}
	}
	if data.UserData.LiveNum > 0 {
		data.UserData.AvgUserTotal /= int64(data.UserData.LiveNum)
		data.UserData.AvgUserCount /= int64(data.UserData.LiveNum)
		data.UserData.AvgInteractRate /= float64(data.UserData.LiveNum)
		data.UserData.AvgIncFansRate /= float64(data.UserData.LiveNum)
	}
	if data.UserData.PromotionLiveNum > 0 {
		data.SaleData.AvgVolume /= int64(data.UserData.PromotionLiveNum)
		data.SaleData.AvgAmount /= float64(data.UserData.PromotionLiveNum)
		data.SaleData.AvgUv /= float64(data.UserData.PromotionLiveNum)
		data.SaleData.AvgPerPrice /= float64(data.UserData.PromotionLiveNum)
	}
	if data.UserData.AvgUserTotal > 0 {
		data.SaleData.SaleRate = float64(data.SaleData.AvgVolume) / float64(data.UserData.AvgUserTotal)
	}
	data.UserTotalChart = dy.DyUserTotalChart{
		Date:       dateChart,
		CountValue: userTotalChart,
		Rooms:      roomChart,
	}
	data.OnlineTimeChart = dy.DateCountFChart{
		Date:       dateChart,
		CountValue: onlineUserChart,
	}
	data.UvChart = dy.DateCountFChart{
		Date:       dateChart,
		CountValue: uvChart,
	}
	data.AmountChart = dy.DateCountFChart{
		Date:       dateChart,
		CountValue: amountChart,
	}
	return
}

//达人电商分析
func (a *AuthorBusiness) GetAuthorProductAnalyse(authorId, keyword, firstCate, secondCate, thirdCate, brandName, sortStr, orderBy string, shopType int, startTime, endTime time.Time, page, pageSize int) (list []entity.DyAuthorProductAnalysis, analysisCount dy.DyAuthorProductAnalysisCount, cateList []dy.DyCate, brandList []dy.NameValueChart, total int, comErr global.CommonError) {
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"price", "gmv", "sales", ""}) {
		comErr = global.NewError(4000)
		return
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		comErr = global.NewError(4000)
		return
	}
	list = []entity.DyAuthorProductAnalysis{}
	cateList = []dy.DyCate{}
	brandList = []dy.NameValueChart{}
	shopId := ""
	if shopType != 0 {
		authorReputation, _ := a.HbaseGetAuthorReputation(authorId)
		shopId = authorReputation.EncryptShopID
	}
	if shopType == 1 && shopId == "" {
		return
	}
	firstCateCountMap := map[string]int{}
	brandNameCountMap := map[string]int{}
	firstCateMap := map[string]map[string]bool{}
	secondCateMap := map[string]map[string]bool{}
	videoNum := 0
	liveNum := 0
	productMapList := map[string]entity.DyAuthorProductAnalysis{}
	var sumGmv float64 = 0
	var sumSale float64 = 0
	hbaseDataList, comErr := a.GetAuthorProductHbaseList(authorId, keyword, startTime, endTime)
	if comErr != nil {
		return
	}
	for _, v := range hbaseDataList {
		//数据过滤
		if keyword != "" && strings.Index(v.Title, keyword) < 0 {
			continue
		}
		if firstCate == "其他" {
			if firstCate != v.DcmLevelFirst || v.DcmLevelFirst != "" {
				continue
			}
		} else {
			if firstCate != "" && firstCate != v.DcmLevelFirst {
				continue
			}
		}
		if secondCate != "" && secondCate != v.FirstCname {
			continue
		}
		if thirdCate != "" && thirdCate != v.SecondCname {
			continue
		}
		if brandName != "" {
			if brandName == "其他" {
				if brandName != v.BrandName && v.BrandName != "" {
					continue
				}
			} else {
				if brandName != v.BrandName {
					continue
				}
			}
		}
		if shopType == 1 && v.ShopId != shopId {
			continue
		} else if shopType == 2 && v.ShopId == shopId && shopId != "" {
			continue
		}
		//数据累加
		videoNum += v.VedioCount
		liveNum += v.RoomCount
		v.Gmv = v.VideoPredictGmv + v.LivePredictGmv
		v.Sales = math.Floor(v.VedioPredictSales) + math.Floor(v.LivePredictSales)
		sumGmv += v.Gmv
		sumSale += math.Floor(v.Sales)
		if p, ok := productMapList[v.ProductId]; ok {
			p.Gmv += v.Gmv
			p.Sales += v.Sales
			p.VideoPredictGmv += v.VideoPredictGmv
			p.LivePredictGmv += v.LivePredictGmv
			p.VedioPredictSales += math.Floor(v.VedioPredictSales)
			p.LivePredictSales += math.Floor(v.LivePredictSales)
			p.RoomCount += v.RoomCount
			p.VedioCount += v.VedioCount
		} else {
			total++
			productMapList[v.ProductId] = v
			//品牌聚合
			brand := v.BrandName
			if brand == "" {
				brand = "其他"
			}
			if _, ok := brandNameCountMap[brand]; !ok {
				brandNameCountMap[brand] = 1
			} else {
				brandNameCountMap[brand] += 1
			}
			//商品分类聚合
			if v.DcmLevelFirst == "" || v.DcmLevelFirst == "null" {
				v.DcmLevelFirst = "其他"
			}
			if _, ok := firstCateMap[v.DcmLevelFirst]; !ok {
				firstCateMap[v.DcmLevelFirst] = map[string]bool{}
			}
			if _, ok := firstCateCountMap[v.DcmLevelFirst]; !ok {
				firstCateCountMap[v.DcmLevelFirst] = 1
			} else {
				firstCateCountMap[v.DcmLevelFirst] += 1
			}
			if v.FirstCname == "" || v.DcmLevelFirst == "其他" {
				continue
			}
			firstCateMap[v.DcmLevelFirst][v.FirstCname] = true
			if _, ok := secondCateMap[v.FirstCname]; !ok {
				secondCateMap[v.FirstCname] = map[string]bool{}
			}
			if v.SecondCname == "" {
				continue
			}
			secondCateMap[v.FirstCname][v.SecondCname] = true
			//简单数据处理
			productMapList[v.ProductId] = v
		}
	}
	//分类处理
	for k, v := range firstCateMap {
		secondCateList := make([]dy.DyCate, 0)
		for ck, _ := range v {
			if cv, ok := secondCateMap[ck]; ok {
				secondCateItem := dy.DyCate{
					Name: ck,
				}
				for tk, _ := range cv {
					secondCateItem.SonCate = append(secondCateItem.SonCate, dy.DyCate{
						Name:    tk,
						Num:     0,
						SonCate: nil,
					})
				}
				if len(secondCateItem.SonCate) == 0 {
					secondCateItem.SonCate = []dy.DyCate{}
				}
				secondCateList = append(secondCateList, secondCateItem)
			}
		}
		productNumber := 0
		if n, ok := firstCateCountMap[k]; ok {
			productNumber = n
		}
		item := dy.DyCate{
			Name:    k,
			Num:     productNumber,
			SonCate: []dy.DyCate{},
		}
		if len(secondCateList) > 0 {
			item.SonCate = secondCateList
		}
		cateList = append(cateList, item)
	}
	//品牌处理
	for k, v := range brandNameCountMap {
		brandList = append(brandList, dy.NameValueChart{
			Name:  k,
			Value: v,
		})
	}
	newList := make([]entity.DyAuthorProductAnalysis, 0)
	for _, v := range productMapList {

		newList = append(newList, v)
	}
	//排序
	if sortStr != "" {
		sort.Slice(newList, func(i, j int) bool {
			var left, right float64
			switch sortStr {
			case "price":
				left = newList[i].Price
				right = newList[j].Price
			case "gmv":
				left = newList[i].Gmv
				right = newList[j].Gmv
			case "sales":
				left = newList[i].Sales
				right = newList[j].Sales
			}
			if left == right {
				return newList[i].ShelfTime > newList[i].ShelfTime
			}
			if orderBy == "desc" {
				return left > right
			}
			return right > left
		})
	}
	start := (page - 1) * pageSize
	end := start + pageSize
	listLen := len(newList)
	if listLen < end {
		end = listLen
	}
	list = newList[start:end]
	for k, v := range list {
		productInfo, _ := hbase.GetProductInfo(v.ProductId)
		if v.Avatar == "" {
			v.Avatar = productInfo.Image
		}
		list[k].Avatar = dyimg.Product(v.Avatar)
		list[k].AuthorId = IdEncrypt(v.AuthorId)
		list[k].ProductId = IdEncrypt(v.ProductId)
		list[k].ProductStatus = productInfo.Status
	}
	analysisCount.ProductNum = total
	analysisCount.RoomNum = liveNum
	analysisCount.VideoNum = videoNum
	analysisCount.Gmv = sumGmv
	analysisCount.Sales = sumSale
	return
}

//获取达人电商分析数据
func (a *AuthorBusiness) GetAuthorProductHbaseList(authorId, keyword string, startTime, endTime time.Time) (hbaseDataList []entity.DyAuthorProductAnalysis, comErr global.CommonError) {
	hbaseDataList = make([]entity.DyAuthorProductAnalysis, 0)
	esAuthorBusiness := es.NewEsAuthorBusiness()
	startRow, stopRow, tmpErr := esAuthorBusiness.AuthorProductAnalysis(authorId, keyword, startTime, endTime)
	if tmpErr != nil {
		comErr = tmpErr
		return
	}
	startRowKey := startRow.AuthorDateProduct
	stopRowKey := stopRow.AuthorDateProduct
	if startRowKey == "" || stopRowKey == "" {
		return
	}
	cacheKey := cache.GetCacheKey(cache.AuthorProductAllList, startRowKey, stopRowKey)
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &hbaseDataList)
	} else {
		hbaseData, _ := hbase.GetAuthorProductAnalysis(stopRowKey)
		if startRowKey != stopRowKey {
			hbaseDataList, _ = hbase.GetAuthorProductAnalysisRange(startRowKey, stopRowKey)
		}
		hbaseDataList = append(hbaseDataList, hbaseData)
		_ = global.Cache.Set(cacheKey, utils.SerializeData(hbaseDataList), 300)
	}
	return
}

//达人电商分析直播列表
func (a *AuthorBusiness) GetAuthorProductRooms(authorId, productId string, startTime, stopTime time.Time, page, pageSize int) (list []dy.DyAuthorProductRoom, total int, comErr global.CommonError) {
	esLiveBusiness := es.NewEsLiveBusiness()
	roomIds, total, comErr := esLiveBusiness.GetAuthorProductSearchRoomIds(authorId, productId, startTime, stopTime, page, pageSize)
	if len(roomIds) == 0 || comErr != nil {
		return
	}
	list = []dy.DyAuthorProductRoom{}
	for _, roomId := range roomIds {
		liveInfo, _ := hbase.GetLiveInfo(roomId)
		liveSaleData, _ := hbase.GetLiveSalesData(roomId)
		//todo gmv数据兼容
		gmv := liveSaleData.Gmv
		sales := liveSaleData.Sales
		if liveSaleData.Gmv == 0 {
			gmv = liveInfo.PredictGmv
			sales = liveInfo.PredictSales
			//if liveInfo.RealGmv > 0 {
			//	gmv = liveInfo.RealGmv
			//	sales = liveInfo.RealSales
			//}
		}
		list = append(list, dy.DyAuthorProductRoom{
			RoomId:       IdEncrypt(roomId),
			Cover:        dyimg.Fix(liveInfo.Cover),
			CreateTime:   liveInfo.CreateTime,
			Title:        liveInfo.Title,
			MaxUserCount: liveInfo.MaxUserCount,
			Gmv:          gmv,
			Sales:        sales,
		})
	}
	return
}

//channel控制go协程获取达人信息
func (a *AuthorBusiness) GetAuthorByIdsLimitGo(authorIds []string, maxNum int) map[string]entity.DyAuthorData {
	var wg sync.WaitGroup
	authorLen := len(authorIds)
	if authorLen < maxNum {
		maxNum = authorLen
	}
	authorChan := make(chan string, maxNum)
	authors := make([]entity.DyAuthorData, 0)
	authorMap := map[string]entity.DyAuthorData{}
	for _, v := range authorIds {
		authorChan <- v
		wg.Add(1)
		go func(aCh chan string, wg *sync.WaitGroup) {
			defer wg.Done()
			authorId, ok := <-aCh
			if ok {
				authorData, comErr := hbase.GetAuthor(authorId)
				if comErr == nil {
					authors = append(authors, authorData.Data)
				}
			}
		}(authorChan, &wg)
	}
	wg.Wait()
	for _, v := range authors {
		authorMap[v.ID] = v
	}
	return authorMap
}
