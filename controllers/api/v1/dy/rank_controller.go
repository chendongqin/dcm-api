package dy

import (
	"dongchamao/business"
	"dongchamao/business/es"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	"dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	"math"
	"strconv"
	"time"
)

type RankController struct {
	controllers.ApiBaseController
}

func (receiver *RankController) Prepare() {
	receiver.InitApiController()
	receiver.CheckDyUserGroupRight(business.DyRankUnLogin, business.DyRankLogin, business.DyJewelRankShowNum)
	//receiver.lockAction()
}

func (receiver *RankController) lockAction() {
	ip := receiver.Ctx.Input.IP()
	if !business.UserActionLock(receiver.TrueUri, ip, 1) {
		receiver.FailReturn(global.NewError(4211))
		return
	}
}

//抖音视频达人热榜
func (receiver *RankController) DyStartAuthorVideoRank() {
	rankType := receiver.GetString("rank_type", "达人指数榜")
	category := receiver.GetString("category", "全部")
	var ret map[string]interface{}
	data, updateTime, _ := hbase.GetStartAuthorVideoRank(rankType, category)
	if !receiver.HasAuth && len(data) > receiver.MaxTotal {
		data = data[0:receiver.MaxTotal]
	}
	for k, v := range data {
		data[k].CoreUserId = business.IdEncrypt(v.CoreUserId)
	}
	ret = map[string]interface{}{
		"has_login":   receiver.HasLogin,
		"has_auth":    receiver.HasAuth,
		"list":        data,
		"update_time": updateTime,
	}
	receiver.SuccReturn(ret)
	return
}

//抖音直播达人热榜
func (receiver *RankController) DyStartAuthorLiveRank() {
	rankType := receiver.GetString("rank_type", "达人指数榜")
	var ret map[string]interface{}
	data, updateTime, _ := hbase.GetStartAuthorLiveRank(rankType)
	for k, v := range data {
		data[k].CoreUserId = business.IdEncrypt(v.CoreUserId)
	}
	if !receiver.HasAuth && len(data) > receiver.MaxTotal {
		data = data[0:receiver.MaxTotal]
	}
	ret = map[string]interface{}{
		"has_login":   receiver.HasLogin,
		"has_auth":    receiver.HasAuth,
		"list":        data,
		"update_time": updateTime,
	}
	receiver.SuccReturn(ret)
	return
}

//抖音直播小时热榜
func (receiver *RankController) DyLiveHourRank() {
	date := receiver.GetString(":date", "")
	hour := receiver.GetString(":hour", "")
	dateTime, err := time.ParseInLocation("2006-01-02 15:04", date+" "+hour, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var ret map[string]interface{}
	data, _ := hbase.GetDyLiveHourRank(dateTime.Format("2006010215"))
	for k, v := range data.Ranks {
		data.Ranks[k].LiveInfo.User.Id = business.IdEncrypt(v.LiveInfo.User.Id)
		data.Ranks[k].RoomId = business.IdEncrypt(v.RoomId)
		data.Ranks[k].LiveInfo.Cover = dyimg.Fix(v.LiveInfo.Cover)
		data.Ranks[k].LiveInfo.User.Avatar = dyimg.Fix(v.LiveInfo.User.Avatar)
		if v.LiveInfo.User.DisplayId == "" {
			data.Ranks[k].LiveInfo.User.DisplayId = v.LiveInfo.User.ShortId
		}
		data.Ranks[k].ShareUrl = business.LiveShareUrl + v.RoomId
		if v.Category == "0" {
			data.Ranks[k].Category = ""
		}
	}
	ret = map[string]interface{}{
		"update_time": data.CrawlTime,
		"has_login":   receiver.HasLogin,
		"has_auth":    receiver.HasAuth,
		"list":        data.Ranks,
	}
	if !receiver.HasAuth && len(data.Ranks) > receiver.MaxTotal {
		ret["list"] = data.Ranks[0:receiver.MaxTotal]
	}
	receiver.SuccReturn(ret)
	return
}

//抖音直播实时榜
func (receiver *RankController) DyLiveTopRank() {
	date := receiver.GetString(":date", "")
	hour := receiver.GetString(":hour", "")
	dateTime, err := time.ParseInLocation("2006-01-02 15:04", date+" "+hour, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	data, _ := hbase.GetDyLiveTopRank(dateTime.Format("2006010215"))
	for k, v := range data.Ranks {
		data.Ranks[k].LiveInfo.User.Id = business.IdEncrypt(v.LiveInfo.User.Id)
		data.Ranks[k].RoomId = business.IdEncrypt(v.RoomId)
		data.Ranks[k].LiveInfo.Cover = dyimg.Fix(v.LiveInfo.Cover)
		data.Ranks[k].LiveInfo.User.Avatar = dyimg.Fix(v.LiveInfo.User.Avatar)
		if v.LiveInfo.User.DisplayId == "" {
			data.Ranks[k].LiveInfo.User.DisplayId = v.LiveInfo.User.ShortId
		}
		data.Ranks[k].ShareUrl = business.LiveShareUrl + v.RoomId
		if v.Category == "0" {
			data.Ranks[k].Category = ""
		}
	}
	if !receiver.HasAuth {
		if len(data.Ranks) > receiver.MaxTotal {
			data.Ranks = data.Ranks[0:receiver.MaxTotal]
		}
	}
	ret := map[string]interface{}{
		"list":        data.Ranks,
		"update_time": data.CrawlTime,
		"has_login":   receiver.HasLogin,
		"has_auth":    receiver.HasAuth,
	}
	receiver.SuccReturn(ret)
	return
}

//抖音直播小时带货榜
func (receiver *RankController) DyLiveHourSellRank() {
	date := receiver.GetString(":date", "")
	hour := receiver.GetString(":hour", "")
	dateTime, err := time.ParseInLocation("2006-01-02 15:04", date+" "+hour, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	data, _ := hbase.GetDyLiveHourSellRank(dateTime.Format("2006010215"))
	for k, v := range data.Ranks {
		data.Ranks[k].Rank = k + 1
		data.Ranks[k].LiveInfo.User.Id = business.IdEncrypt(v.LiveInfo.User.Id)
		data.Ranks[k].RoomId = business.IdEncrypt(v.RoomId)
		data.Ranks[k].LiveInfo.Cover = dyimg.Fix(v.LiveInfo.Cover)
		data.Ranks[k].LiveInfo.User.Avatar = dyimg.Fix(v.LiveInfo.User.Avatar)
		if v.LiveInfo.User.DisplayId == "" {
			data.Ranks[k].LiveInfo.User.DisplayId = v.LiveInfo.User.ShortId
		}
		shopTags := make([]string, 0)
		for _, s := range v.ShopTags {
			if s == "" {
				continue
			}
			shopTags = append(shopTags, s)
		}
		data.Ranks[k].ShopTags = shopTags
		data.Ranks[k].ShareUrl = business.LiveShareUrl + v.RoomId
		//if v.RealGmv > 0 {
		//	data.Ranks[k].PredictGmv = v.RealGmv
		//	data.Ranks[k].PredictSales = v.RealSales
		//}
	}
	var ret map[string]interface{}
	if !receiver.HasAuth && len(data.Ranks) > receiver.MaxTotal {
		data.Ranks = data.Ranks[0:receiver.MaxTotal]
	}
	ret = map[string]interface{}{
		"list":        data.Ranks,
		"update_time": data.CrawlTime,
		"has_login":   receiver.HasLogin,
		"has_auth":    receiver.HasAuth,
	}
	receiver.SuccReturn(ret)
	return
}

//抖音直播实时榜
func (receiver *RankController) DyLiveHourPopularityRank() {
	date := receiver.GetString(":date", "")
	hour := receiver.GetString(":hour", "")
	dateTime, err := time.ParseInLocation("2006-01-02 15:04", date+" "+hour, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}

	data, _ := hbase.GetDyLiveHourPopularityRank(dateTime.Format("2006010215"))
	for k, v := range data.Ranks {
		data.Ranks[k].LiveInfo.User.Id = business.IdEncrypt(v.LiveInfo.User.Id)
		data.Ranks[k].RoomId = business.IdEncrypt(v.RoomId)
		data.Ranks[k].LiveInfo.Cover = dyimg.Fix(v.LiveInfo.Cover)
		data.Ranks[k].LiveInfo.User.Avatar = dyimg.Fix(v.LiveInfo.User.Avatar)
		if v.LiveInfo.User.DisplayId == "" {
			data.Ranks[k].LiveInfo.User.DisplayId = v.LiveInfo.User.ShortId
		}
		data.Ranks[k].ShareUrl = business.LiveShareUrl + v.RoomId
	}
	if !receiver.HasAuth {
		if len(data.Ranks) > receiver.MaxTotal {
			data.Ranks = data.Ranks[0:receiver.MaxTotal]
		}
	}
	var ret = map[string]interface{}{
		"list":        data.Ranks,
		"update_time": data.CrawlTime,
		"has_login":   receiver.HasLogin,
		"has_auth":    receiver.HasAuth,
	}
	receiver.SuccReturn(ret)
	return
}

//抖音直播达人分享周榜
func (receiver *RankController) DyLiveShareWeekRank() {
	start, end, comErr := receiver.GetRealRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	if start.Weekday() != 1 || end.Weekday() != 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if start.AddDate(0, 0, 6) != end {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	data, _ := hbase.GetLiveShareWeekRank(start.Format("20060102") + "_" + end.Format("20060102"))
	list := make([]entity.DyLiveShareWeekData, 0)
	for _, v := range data.Data {
		var gmv float64 = 0
		var sales float64 = 0
		var totalUser int64 = 0
		for _, r := range v.Rooms {
			//if r.RealSales > 0 {
			//	gmv += r.RealGmv
			//	sales += math.Floor(r.RealSales)
			//} else {
			gmv += r.PredictGmv
			sales += math.Floor(r.PredictSales)
			//}
			totalUser += r.TotalUser
		}
		uniqueId := v.UniqueId
		if uniqueId == "" || uniqueId == "0" {
			uniqueId = v.ShortId
		}
		roomNum := len(v.Rooms)
		var TotalUser int64
		if roomNum != 0 {
			TotalUser = totalUser / int64(roomNum)
		} else {
			TotalUser = int64(0)
		}
		list = append(list, entity.DyLiveShareWeekData{
			AuthorId:   business.IdEncrypt(utils.ToString(v.AuthorId)),
			Avatar:     dyimg.Avatar(v.Avatar),
			Category:   v.Category,
			InitRank:   v.InitRank,
			Name:       v.Name,
			RankChange: v.RankChange,
			Score:      v.Score,
			UniqueId:   uniqueId,
			Gmv:        gmv,
			Sales:      sales,
			TotalUser:  TotalUser,
			RoomNum:    roomNum,
		})
	}
	if !receiver.HasAuth {
		if len(list) > receiver.MaxTotal {
			list = list[0:receiver.MaxTotal]
		}
	}
	var ret = map[string]interface{}{
		"list":        list,
		"update_time": data.CrawlTime,
		"has_login":   receiver.HasLogin,
		"has_auth":    receiver.HasAuth,
	}
	receiver.SuccReturn(ret)
	return
}

//抖音视频达人分享日榜
func (receiver *RankController) DyAwemeShareRank() {
	date := receiver.Ctx.Input.Param(":date")
	if date == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	pslTime := "2006-01-02"
	dateTime, err := time.ParseInLocation(pslTime, date, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	data, _ := hbase.GetAwemeShareRank(dateTime.Format("20060102"))
	list := make([]entity.DyAwemeShareTopCopy, 0)
	for _, v := range data.Data {
		uniqueId := v.UniqueId
		if uniqueId == "" || uniqueId == "0" {
			uniqueId = v.ShortId
		}
		list = append(list, entity.DyAwemeShareTopCopy{
			AuthorId:      business.IdEncrypt(utils.ToString(v.AuthorId)),
			Category:      v.Category,
			Avatar:        dyimg.Avatar(v.Avatar),
			InitRank:      v.InitRank,
			Name:          v.Name,
			RankChange:    v.RankChange,
			Score:         v.Score,
			UniqueId:      uniqueId,
			FollowerCount: v.FollowerCount,
			IncDiggCount:  v.IncDiggCount,
		})
	}
	if !receiver.HasAuth && len(list) > receiver.MaxTotal {
		list = list[0:receiver.MaxTotal]
	}
	var ret = map[string]interface{}{
		"list":        list,
		"update_time": data.CrawlTime,
		"has_login":   receiver.HasLogin,
		"has_auth":    receiver.HasAuth,
	}
	receiver.SuccReturn(ret)
	return
}

//抖音销量日榜
func (receiver *RankController) ProductSalesTopDayRank() {
	date := receiver.Ctx.Input.Param(":date")
	fCate := receiver.GetString("first_cate", "")
	sCate := receiver.GetString("second_cate", "")
	tCate := receiver.GetString("third_cate", "")
	sortStr := receiver.GetString("sort", "")
	orderBy := receiver.GetString("order", "")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	if date == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	pslTime := "2006-01-02"
	dateTime, err := time.ParseInLocation(pslTime, date, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if !receiver.HasAuth {
		if page != 1 {
			receiver.FailReturn(global.NewError(4004))
			return
		}
		pageSize = receiver.MaxTotal
	}
	list, total, comErr := es.NewEsProductBusiness().ProductSalesTopDayRank(dateTime.Format("20060102"), fCate, sCate, tCate, sortStr, orderBy, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		list[k].ProductId = business.IdEncrypt(v.ProductId)
		list[k].Images = dyimg.Fix(v.Images)
	}
	if total > receiver.MaxTotal {
		total = receiver.MaxTotal
	}
	receiver.SuccReturn(map[string]interface{}{
		"has_login": receiver.HasLogin,
		"has_auth":  receiver.HasAuth,
		"list":      list,
		"total":     total,
	})
	return
}

//抖音热推日榜
func (receiver *RankController) ProductShareTopDayRank() {
	date := receiver.Ctx.Input.Param(":date")
	fCate := receiver.GetString("first_cate", "")
	sCate := receiver.GetString("second_cate", "")
	tCate := receiver.GetString("third_cate", "")
	sortStr := receiver.GetString("sort", "")
	orderBy := receiver.GetString("order", "")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	if date == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	pslTime := "2006-01-02"
	dateTime, err := time.ParseInLocation(pslTime, date, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if !receiver.HasAuth {
		if page != 1 {
			receiver.FailReturn(global.NewError(4004))
			return
		}
		pageSize = receiver.MaxTotal
	}
	list, total, comErr := es.NewEsProductBusiness().ProductShareTopDayRank(dateTime.Format("20060102"), fCate, sCate, tCate, sortStr, orderBy, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		list[k].ProductId = business.IdEncrypt(v.ProductId)
		list[k].Images = dyimg.Fix(v.Images)
	}
	if total > receiver.MaxTotal {
		total = receiver.MaxTotal
	}
	receiver.SuccReturn(map[string]interface{}{
		"has_login": receiver.HasLogin,
		"has_auth":  receiver.HasAuth,
		"list":      list,
		"total":     total,
	})
	return
}

//达人带货榜
func (receiver *RankController) DyAuthorTakeGoodsRank() {
	date := receiver.Ctx.Input.Param(":date")
	startDate, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	//dateType, _ := receiver.GetInt("date_type", 1)
	tags := receiver.GetString("tags", "all")
	//verified, _ := receiver.GetInt("verified")
	sortStr := receiver.GetString("sort", "sum_sales")
	orderBy := receiver.GetString("order_by", "desc")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	switch sortStr {
	case "sum_sales":
		sortStr = "predict_sales_sum"
		break
	case "sum_gmv":
		sortStr = "predict_gmv_sum"
		break
	case "avg_price":
		sortStr = "per_price"
		break
	}
	if tags == "其他" {
		tags = "null"
	}
	if !receiver.HasAuth {
		//dateType = 1
		page = 1
		if pageSize > receiver.MaxTotal {
			pageSize = receiver.MaxTotal
		}
	}
	start := (page - 1) * pageSize
	end := page * pageSize
	if end > receiver.MaxTotal {
		end = receiver.MaxTotal
	}
	var ret map[string]interface{}
	//cacheKey := cache.GetCacheKey(cache.DyRankCache, "author_take_goods", utils.Md5_encode(fmt.Sprintf("%s%d%s%s%s%d%d%d", startDate, dateType, tags, sortStr, orderBy, verified, page, pageSize)))
	//cacheStr := global.Cache.Get(cacheKey)
	//if cacheStr != "" {
	//	cacheStr = utils.DeserializeData(cacheStr)
	//	_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	//} else {
	var originList []entity.DyAuthorDaySalesRank
	key := sortStr + "_" + startDate.Format("20060102") + "_" + tags
	rowKeys := make([][]byte, 0)
	//0925无数据特殊处理
	//if firstRow.AuthorId == "" && startDate.Format("20060102") == "20210925" {
	//	key = sortStr + "_" + startDate.AddDate(0, 0, -1).Format("20060102") + "_" + tags
	//	rowKey = utils.Md5_encode(key) + "_" + strconv.Itoa(1)
	//}
	if orderBy == "desc" {
		for i := start + 1; i <= end; i++ {
			rowKeys = append(rowKeys, []byte(utils.Md5_encode(key)+"_"+strconv.Itoa(i)))
		}
		originList, _ = hbase.GetSaleAuthorRank(rowKeys)
	} else {
		firstRow, _ := hbase.GetSaleAuthorRow(utils.Md5_encode(key) + "_" + strconv.Itoa(1))
		maxRow, _ := strconv.Atoi(firstRow.RnMax)
		if maxRow > 0 {
			for i := maxRow - start; i >= maxRow-end+1; i-- {
				rowKeys = append(rowKeys, []byte(utils.Md5_encode(key)+"_"+strconv.Itoa(i)))
			}
			originList, _ = hbase.GetSaleAuthorRank(rowKeys)
		}
	}
	data := make([]dy.TakeGoodsRankRet, 0)
	total := 0
	for k, v := range originList {
		total, _ = strconv.Atoi(v.RnMax)
		tempData := dy.TakeGoodsRankRet{}
		tempData.Rank = (page-1)*pageSize + k + 1
		tempData.AuthorId = business.IdEncrypt(v.AuthorId)
		tempData.UniqueId = v.ShortId
		tempData.Nickname = v.Nickname
		tempData.AuthorCover = dyimg.Fix(v.Avatar)
		tempData.VerificationType, _ = strconv.Atoi(v.VerificationType)
		tempData.VerifyName = v.VerifyName
		tempData.Tags = v.Tags
		tempData.SumSales, _ = strconv.ParseFloat(v.PredictSalesSum, 64)
		tempData.SumGmv, _ = strconv.ParseFloat(v.PredictGmvSum, 64)
		tempData.AvgPrice, _ = strconv.ParseFloat(v.PerPrice, 64)
		tempData.RoomCount, _ = strconv.Atoi(v.RoomIdCount)
		var roomList = []map[string]interface{}{}
		tempData.RoomList = roomList
		data = append(data, tempData)
	}

	//list, total, _ := es.NewEsAuthorBusiness().SaleAuthorRankCount(startDate, dateType, tags, sortStr, orderBy, verified, page, pageSize)
	//var structData []es2.DyAuthorTakeGoodsCount
	//utils.MapToStruct(list, &structData)
	//data := make([]dy.TakeGoodsRankRet, len(structData))
	//for k, v := range structData {
	//	hits := v.Hit.Hits.Hits
	//	uniqueId := hits[0].Source.UniqueID
	//	if uniqueId == "" {
	//		uniqueId = hits[0].Source.ShortID
	//	}
	//	var roomList = make([]map[string]interface{}, 0, len(hits))
	//	for _, v := range hits {
	//		roomList = append(roomList, map[string]interface{}{
	//			"room_cover":     dyimg.Fix(v.Source.RoomCover),
	//			"room_id":        business.IdEncrypt(v.Source.RoomID),
	//			"room_title":     v.Source.RoomTitle,
	//			"date_time":      v.Source.CreateTime,
	//			"max_user_count": v.Source.MaxUserCount,
	//			"gmv":            v.Source.PredictGmv,
	//			"sales":          v.Source.PredictSales,
	//		})
	//	}
	//	data[k] = dy.TakeGoodsRankRet{
	//		Rank:             (page-1)*pageSize + k + 1,
	//		Nickname:         hits[0].Source.Nickname,
	//		VerificationType: hits[0].Source.VerificationType,
	//		VerifyName:       hits[0].Source.VerifyName,
	//		AuthorCover:      dyimg.Avatar(hits[0].Source.AuthorCover),
	//		SumGmv:           v.SumGmv.Value,
	//		SumSales:         v.SumSales.Value,
	//		AvgPrice:         v.AvgPrice.Value,
	//		AuthorId:         business.IdEncrypt(utils.ToString(v.Key.AuthorID)),
	//		RoomCount:        len(hits),
	//		Tags:             hits[0].Source.Tags,
	//		UniqueId:         business.IdEncrypt(utils.ToString(uniqueId)),
	//		RoomList:         roomList,
	//	}
	//}
	if total > receiver.MaxTotal {
		total = receiver.MaxTotal
	}

	list := make([]dy.TakeGoodsRankRet, 0)
	if total > 0 {
		if start > 0 {
			lens := end - start
			list = data[0:lens]
		} else {
			list = data[start:end]
		}
	}

	ret = map[string]interface{}{
		"list":  list,
		"total": total,
	}
	//	if startDate.Format("20060102") != time.Now().Format("20060102") {
	//		_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 86400)
	//	}
	//}
	ret["has_login"] = receiver.HasLogin
	ret["has_auth"] = receiver.HasAuth
	if !receiver.HasAuth && utils.ToInt(ret["total"]) > receiver.MaxTotal {
		ret["total"] = receiver.MaxTotal
	}
	receiver.SuccReturn(ret)
	return
}

//达人涨粉榜
func (receiver *RankController) DyAuthorFollowerRank() {
	date := receiver.Ctx.Input.Param(":date")
	startDate, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	tags := receiver.GetString("tags", "all")
	sortStr := receiver.GetString("sort", "sum_sales")
	orderBy := receiver.GetString("order_by", "desc")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	var sortMap = map[string]string{"live_inc_follower_count": "live_fans_inc", "inc_follower_count": "fans_inc", "aweme_inc_follower_count": "aweme_fans_inc"}
	if sortMap[sortStr] != "" {
		sortStr = sortMap[sortStr]
	}
	if tags == "其他" {
		tags = "null"
	}
	if !receiver.HasAuth {
		page = 1
		if pageSize > receiver.MaxTotal {
			pageSize = receiver.MaxTotal
		}
	}
	start := (page - 1) * pageSize
	end := page * pageSize
	if end > receiver.MaxTotal {
		end = receiver.MaxTotal
	}
	var ret map[string]interface{}
	var originList []entity.DyAuthorDayFansIncrease
	rowKey := sortStr + "_" + startDate.Format("20060102") + "_" + tags
	var rowKeys [][]byte
	if orderBy == "desc" {
		startTemp := start
		for {
			rowKeys = append(rowKeys, []byte(utils.Md5_encode(rowKey)+"_"+strconv.Itoa(startTemp+1)))
			startTemp++
			if startTemp >= end {
				break
			}
		}
		originList, _ = hbase.GetFansAuthorRank(rowKeys)
	} else {
		firstRow, _ := hbase.GetFansAuthorRow(utils.Md5_encode(rowKey) + "_" + strconv.Itoa(1))
		maxRow, _ := strconv.Atoi(firstRow.RnMax)
		if maxRow > 0 {
			startTemp := maxRow - end
			endTemp := maxRow - start
			for {
				rowKeys = append(rowKeys, []byte(utils.Md5_encode(rowKey)+"_"+strconv.Itoa(startTemp+1)))
				startTemp++
				if startTemp >= endTemp {
					break
				}
			}
			originList, _ = hbase.GetFansAuthorRank(rowKeys)
		}
	}
	data := make([]dy.AuthorFansRankRet, 0)
	total := 0
	//0未认证；1蓝v；2黄v
	var VerificationTypeMap = map[string]int{"没有认证": 0, "蓝v": 1, "黄v": 2}
	for k, v := range originList {
		total, _ = strconv.Atoi(v.RnMax)
		tempData := dy.AuthorFansRankRet{}
		tempData.Rank = (page-1)*pageSize + k + 1
		tempData.AuthorId = business.IdEncrypt(v.AuthorId)
		if v.UniqueId != "" {
			tempData.UniqueId = v.UniqueId
		} else {
			tempData.UniqueId = v.ShortId
		}
		tempData.Nickname = v.Nickname
		tempData.AuthorCover = dyimg.Fix(v.Avatar)
		tempData.VerificationType = VerificationTypeMap[v.VerificationType]
		tempData.VerifyName = v.VerifyName
		tempData.FollowerCount, _ = strconv.Atoi(v.FollowerCount)
		tempData.IncFollowerCount, _ = strconv.Atoi(v.FansInc)
		tempData.LiveIncFollowerCount, _ = strconv.Atoi(v.LiveFansInc)
		tempData.AwemeIncFollowerCount, _ = strconv.Atoi(v.AwemeFansInc)
		tempData.Tags = v.Tags
		data = append(data, tempData)
	}
	if total > receiver.MaxTotal {
		total = receiver.MaxTotal
	}
	list := make([]dy.AuthorFansRankRet, 0)
	if total > 0 {
		if start > 0 {
			lens := end - start
			list = data[0:lens]
		} else {
			list = data[start:end]
		}
	}
	ret = map[string]interface{}{
		"list":  list,
		"total": total,
	}
	ret["has_login"] = receiver.HasLogin
	ret["has_auth"] = receiver.HasAuth
	if !receiver.HasAuth && utils.ToInt(ret["total"]) > receiver.MaxTotal {
		ret["total"] = receiver.MaxTotal
	}
	receiver.SuccReturn(ret)
	return
}

////达人涨粉榜
//func (receiver *RankController) DyAuthorFollowerRank() {
//	date := receiver.Ctx.Input.Param(":date")
//	tags := receiver.GetString("tags", "")
//	province := receiver.GetString("province", "")
//	city := receiver.GetString("city", "")
//	sortStr := receiver.GetString("sort", "")
//	orderBy := receiver.GetString("order_by", "")
//	isDelivery, _ := receiver.GetInt("is_delivery", 0)
//	page := receiver.GetPage("page")
//	pageSize := receiver.GetPageSize("page_size", 10, 100)
//	if !receiver.HasAuth {
//		page = 1
//		if pageSize > receiver.MaxTotal {
//			pageSize = receiver.MaxTotal
//		}
//	}
//	dateTime, err := time.ParseInLocation("2006-01-02", date, time.Local)
//	if err != nil {
//		receiver.FailReturn(global.NewError(4000))
//		return
//	}
//	data, total, comErr := es.NewEsAuthorBusiness().DyAuthorFollowerIncRank(dateTime.Format("20060102"), tags, province, city, sortStr, orderBy, isDelivery, page, pageSize)
//	form := (page - 1) * pageSize
//	for k, v := range data {
//
//		data[k].Rank = k + form + 1
//		if data[k].UniqueID == "" {
//			data[k].UniqueID = v.ShortID
//		}
//		data[k].AuthorID = business.IdEncrypt(v.AuthorID)
//		data[k].AuthorCover = dyimg.Fix(v.AuthorCover)
//	}
//	if comErr != nil {
//		receiver.FailReturn(global.NewError(4000))
//		return
//	}
//	if total > receiver.MaxTotal {
//		total = receiver.MaxTotal
//	}
//	ret := map[string]interface{}{
//		"list": data,
//	}
//	ret["has_login"] = receiver.HasLogin
//	ret["has_auth"] = receiver.HasAuth
//	if !receiver.HasAuth && total > receiver.MaxTotal {
//		ret["total"] = receiver.MaxTotal
//	}
//	receiver.SuccReturn(ret)
//	return
//}

//短视频商品排行榜
func (receiver *RankController) VideoProductRank() {
	date := receiver.Ctx.Input.Param(":date")
	dataType, _ := receiver.GetInt("data_type", 1)
	fCate := receiver.GetString("first_cate", "all")
	sortStr := receiver.GetString("sort", "sales")
	orderBy := receiver.GetString("order_by", "desc")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	if date == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	pslTime := "2006-01-02"
	dateTime, err := time.ParseInLocation(pslTime, date, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var rowKey = ""
	switch dataType {
	case 1: //日榜
		day := dateTime.Format("20060102")
		key := day + "_" + fCate + "_" + sortStr
		rowKey = utils.Md5_encode(key)
		break
	case 2: //周榜
		startTime := dateTime
		lastWeekStartTime := dateTime.AddDate(0, 0, -7)
		firstDay := lastWeekStartTime.Format("02")
		lastWeekEndTime := dateTime.AddDate(0, 0, -1)
		endDay := lastWeekEndTime.Format("02")
		if startTime.Weekday() != 1 {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		weekRange := startTime.Format("20060102") + firstDay + endDay
		key := weekRange + "_" + fCate + "_" + sortStr
		rowKey = utils.Md5_encode(key)
		break
	case 3: //月榜
		month := dateTime.Format("200601")
		key := month + "_" + fCate + "_" + sortStr
		rowKey = utils.Md5_encode(key)
		break
	}
	if !receiver.HasAuth {
		if page != 1 || dataType > 1 {
			receiver.FailReturn(global.NewError(4004))
			return
		}
		pageSize = receiver.MaxTotal
	}
	start := (page - 1) * pageSize
	end := page * pageSize
	total := 0
	finished := false
	list := make([]entity.ShortVideoProduct, 0)
	if orderBy == "asc" {
		for i := 0; i < 5; i++ {
			tempData, _ := hbase.GetVideoProductRank(rowKey, i)
			lenNum := len(tempData)
			tmpTotal := total
			total += lenNum
			if finished {
				continue
			}
			if total > start {
				if end <= total {
					list = append(list, tempData[start-tmpTotal:end-tmpTotal]...)
					finished = true
				} else {
					list = append(list, tempData[start-tmpTotal:]...)
					start = total
				}
			}
		}
	} else {
		for i := 4; i >= 0; i-- {
			tempData, _ := hbase.GetVideoProductRank(rowKey, i)
			lenNum := len(tempData)
			for j := 0; j < lenNum/2; j++ { //倒序
				temp := tempData[lenNum-1-j]
				tempData[lenNum-1-j] = tempData[j]
				tempData[j] = temp
			}
			tmpTotal := total
			total += lenNum
			if finished {
				continue
			}
			if total > start {
				if end <= total {
					list = append(list, tempData[start-tmpTotal:end-tmpTotal]...)
					finished = true
				} else {
					list = append(list, tempData[start-tmpTotal:]...)
					start = total
				}
			}
		}
	}
	if !receiver.HasAuth && total > receiver.MaxTotal {
		total = receiver.MaxTotal
	}
	for k, v := range list {
		list[k].ProductId = business.IdEncrypt(v.ProductId)
	}
	ret := map[string]interface{}{
		"list":      list,
		"has_login": receiver.HasLogin,
		"has_auth":  receiver.HasAuth,
		"total":     total,
	}
	receiver.SuccReturn(ret)
	return
}

//直播商品榜
func (receiver *RankController) LiveProductRank() {
	date := receiver.Ctx.Input.Param(":date")
	dataType, _ := receiver.GetInt("data_type", 1)
	fCate := receiver.GetString("first_cate", "all")
	sortStr := receiver.GetString("sort", "sales")
	orderBy := receiver.GetString("order_by", "desc")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	switch sortStr {
	case "live_count":
		sortStr = "awemenum"
		break
	case "gmv":
		sortStr = "saleroom"
		break
	case "cos_fee":
		sortStr = "fee"
		break
	}
	if date == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	pslTime := "2006-01-02"
	dateTime, err := time.ParseInLocation(pslTime, date, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var rowKey = ""
	switch dataType {
	case 1: //日榜
		day := dateTime.Format("20060102")
		key := day + "_" + fCate + "_" + sortStr
		rowKey = utils.Md5_encode(key)
		break
	case 2: //周榜
		startTime := dateTime
		lastWeekStartTime := dateTime.AddDate(0, 0, -7)
		firstDay := lastWeekStartTime.Format("02")
		lastWeekEndTime := dateTime.AddDate(0, 0, -1)
		endDay := lastWeekEndTime.Format("02")
		if startTime.Weekday() != 1 {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		weekRange := startTime.Format("20060102") + firstDay + endDay
		key := weekRange + "_" + fCate + "_" + sortStr
		rowKey = utils.Md5_encode(key)
		break
	case 3: //月榜
		month := dateTime.Format("200601")
		key := month + "_" + fCate + "_" + sortStr
		rowKey = utils.Md5_encode(key)
		break
	}
	if !receiver.HasAuth {
		if page != 1 || dataType > 1 {
			receiver.FailReturn(global.NewError(4004))
			return
		}
		pageSize = receiver.MaxTotal
	}
	start := (page - 1) * pageSize
	end := page * pageSize
	total := 0
	finished := false
	orginList := make([]entity.LiveProduct, 0)
	if orderBy == "asc" {
		for i := 0; i < 5; i++ {
			tempData, _ := hbase.GetLiveProductRank(rowKey, i)
			lenNum := len(tempData)
			tmpTotal := total
			total += lenNum
			if finished {
				continue
			}
			if total > start {
				if end <= total {
					orginList = append(orginList, tempData[start-tmpTotal:end-tmpTotal]...)
					finished = true
				} else {
					orginList = append(orginList, tempData[start-tmpTotal:]...)
					start = total
				}
			}
		}
	} else {
		for i := 4; i >= 0; i-- {
			tempData, _ := hbase.GetLiveProductRank(rowKey, i)
			lenNum := len(tempData)
			for j := 0; j < lenNum/2; j++ { //倒序
				temp := tempData[lenNum-1-j]
				tempData[lenNum-1-j] = tempData[j]
				tempData[j] = temp
			}
			tmpTotal := total
			total += lenNum
			if finished {
				continue
			}
			if total > start {
				if end <= total {
					orginList = append(orginList, tempData[start-tmpTotal:end-tmpTotal]...)
					finished = true
				} else {
					orginList = append(orginList, tempData[start-tmpTotal:]...)
					start = total
				}
			}
		}
	}
	if !receiver.HasAuth && total > receiver.MaxTotal {
		total = receiver.MaxTotal
	}
	list := make([]entity.DyLiveProductSaleTopRank, 0)
	for _, v := range orginList {
		tempData := entity.DyLiveProductSaleTopRank{}
		tempData.ProductId = business.IdEncrypt(v.ProductId)
		tempData.Images = v.Image
		tempData.Title = v.Title
		tempData.LiveCount = v.RoomNum
		tempData.PlatformLabel = v.PlatformLabel
		tempData.Price = v.Price
		tempData.CosFee = v.CosFee
		tempData.CosRatio = v.CosRatio
		tempData.Gmv = v.Saleroom
		tempData.Sales = float64(v.Sales)
		list = append(list, tempData)
	}

	ret := map[string]interface{}{
		"list":      list,
		"has_login": receiver.HasLogin,
		"has_auth":  receiver.HasAuth,
		"total":     total,
	}
	receiver.SuccReturn(ret)
	return
}
