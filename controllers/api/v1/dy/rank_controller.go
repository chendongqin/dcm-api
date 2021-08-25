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
	es2 "dongchamao/models/es"
	"dongchamao/services/dyimg"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"math"
	"time"
)

type RankController struct {
	controllers.ApiBaseController
}

//抖音视频达人热榜
func (receiver *RankController) DyStartAuthorVideoRank() {
	rankType := receiver.GetString("rank_type", "达人指数榜")
	category := receiver.GetString("category", "全部")
	var ret map[string]interface{}
	cacheKey := cache.GetCacheKey(cache.DyStartAuthorVideoRank, rankType)
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	} else {
		data, updateTime, _ := hbase.GetStartAuthorVideoRank(rankType, category)
		for k, v := range data {
			data[k].CoreUserId = business.IdEncrypt(v.CoreUserId)
		}
		ret = map[string]interface{}{
			"list":        data,
			"update_time": updateTime,
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 86400)
	}
	receiver.SuccReturn(ret)
	return
}

//抖音直播达人热榜
func (receiver *RankController) DyStartAuthorLiveRank() {
	rankType := receiver.GetString("rank_type", "达人指数榜")
	var ret map[string]interface{}
	cacheKey := cache.GetCacheKey(cache.DyStartAuthorLiveRank, rankType)
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	} else {
		data, updateTime, _ := hbase.GetStartAuthorLiveRank(rankType)
		for k, v := range data {
			data[k].CoreUserId = business.IdEncrypt(v.CoreUserId)
		}
		ret = map[string]interface{}{
			"list":        data,
			"update_time": updateTime,
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 86400)
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
	cacheKey := cache.GetCacheKey(cache.DyLiveHourRank, dateTime.Format("2006010215"))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	} else {
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
		}
		ret = map[string]interface{}{
			"list":        data.Ranks,
			"update_time": data.CrawlTime,
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 300)
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
	var ret map[string]interface{}
	cacheKey := cache.GetCacheKey(cache.DyLiveTopRank, dateTime.Format("2006010215"))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	} else {
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
		}
		ret = map[string]interface{}{
			"list":        data.Ranks,
			"update_time": data.CrawlTime,
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 86400)
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
	var ret map[string]interface{}
	cacheKey := cache.GetCacheKey(cache.DyLiveHourSellRank, dateTime.Format("2006010215"))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	} else {
		data, _ := hbase.GetDyLiveHourSellRank(dateTime.Format("2006010215"))
		for k, v := range data.Ranks {
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
		ret = map[string]interface{}{
			"list":        data.Ranks,
			"update_time": data.CrawlTime,
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 300)
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
	var ret map[string]interface{}
	cacheKey := cache.GetCacheKey(cache.DyLiveHourPopularityRank, dateTime.Format("2006010215"))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	} else {
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
		ret = map[string]interface{}{
			"list":        data.Ranks,
			"update_time": data.CrawlTime,
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 86400)
	}
	receiver.SuccReturn(ret)
	return
}

//抖音直播达人分享周榜
func (receiver *RankController) DyLiveShareWeekRank() {
	start, end, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	if start.Weekday() != 1 || end.Weekday() != 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if end.Day()-start.Day() != 6 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var ret map[string]interface{}
	cacheKey := cache.GetCacheKey(cache.DyLiveShareWeekRank, start.Format("20060102")+"_"+end.Format("20060102"))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	} else {
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
				TotalUser:  totalUser / int64(roomNum),
				RoomNum:    roomNum,
			})
		}
		ret = map[string]interface{}{
			"list":        list,
			"update_time": data.CrawlTime,
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 86400)
	}
	receiver.SuccReturn(ret)
	return
}

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
	var ret map[string]interface{}
	cacheKey := cache.GetCacheKey(cache.DyAwemeShareRank, dateTime.Format("20060102"))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	} else {
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
		ret = map[string]interface{}{
			"list":        list,
			"update_time": data.CrawlTime,
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 86400)
	}
	receiver.SuccReturn(ret)
	return
}

type TakeGoodsRankRet struct {
	Rank             int                      `json:"rank,omitempty"`
	Nickname         string                   `json:"nickname,omitempty"`
	AuthorCover      string                   `json:"author_cover,omitempty"`
	SumGmv           float64                  `json:"sum_gmv,omitempty"`
	SumSales         float64                  `json:"sum_sales,omitempty"`
	AvgPrice         float64                  `json:"avg_price,omitempty"`
	AuthorId         string                   `json:"author_id,omitempty"`
	UniqueId         string                   `json:"unique_id,omitempty"`
	Tags             string                   `json:"tags,omitempty"`
	VerificationType int                      `json:"verification_type,omitempty"`
	VerifyName       string                   `json:"verify_name,omitempty"`
	RoomCount        int                      `json:"room_count,omitempty"`
	RoomList         []map[string]interface{} `json:"room_list"`
}

//达人带货榜
func (receiver *RankController) DyAuthorTakeGoodsRank() {
	date := receiver.GetString("date")
	startDate, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	dateType, _ := receiver.GetInt("date_type", 1)
	tags := receiver.GetString("tags")
	verified, _ := receiver.GetInt("verified")
	sortStr := receiver.GetString("sort", "sum_gmv")
	orderBy := receiver.GetString("order_by", "desc")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPage("page_size")
	var ret map[string]interface{}
	cacheKey := cache.GetCacheKey(cache.DyAuthorTakeGoodsRank, startDate, dateType, tags, sortStr, orderBy, verified, page, pageSize)
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	} else {
		list, total, updateTime, _ := es.NewEsAuthorBusiness().SaleAuthorRankCount(startDate, dateType, tags, sortStr, orderBy, verified, page, pageSize)
		var structData []es2.DyAuthorTakeGoodsCount
		utils.MapToStruct(list, &structData)
		data := make([]TakeGoodsRankRet, len(structData))
		for k, v := range structData {
			hits := v.Hit.Hits.Hits
			uniqueId := hits[0].Source.UniqueID
			if uniqueId == "" {
				uniqueId = hits[0].Source.ShortID
			}
			var roomList = make([]map[string]interface{}, 0, len(hits))
			for _, v := range hits {
				roomList = append(roomList, map[string]interface{}{
					"room_cover":     dyimg.Fix(v.Source.RoomCover),
					"room_id":        business.IdEncrypt(v.Source.RoomID),
					"room_title":     v.Source.RoomTitle,
					"date_time":      v.Source.CreateTime,
					"max_user_count": v.Source.MaxUserCount,
					"real_gmv":       v.Source.RealGmv,
					"real_sales":     v.Source.RealSales,
				})
			}
			data[k] = TakeGoodsRankRet{
				Rank:             (page-1)*pageSize + k + 1,
				Nickname:         hits[0].Source.Nickname,
				VerificationType: hits[0].Source.VerificationType,
				VerifyName:       hits[0].Source.VerifyName,
				AuthorCover:      dyimg.Avatar(hits[0].Source.AuthorCover),
				SumGmv:           v.SumGmv.Value,
				SumSales:         v.SumSales.Value,
				AvgPrice:         v.AvgPrice.Value,
				AuthorId:         business.IdEncrypt(utils.ToString(v.Key.AuthorID)),
				RoomCount:        len(hits),
				Tags:             hits[0].Source.Tags,
				UniqueId:         business.IdEncrypt(utils.ToString(uniqueId)),
				RoomList:         roomList,
			}
		}
		ret = map[string]interface{}{
			"list":        data,
			"total":       total,
			"update_time": updateTime,
		}
		_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 86400)
	}
	receiver.SuccReturn(ret)
	return
}

//达人涨粉榜
func (receiver *RankController) DyAuthorFollowerRank() {
	date := receiver.Ctx.Input.Param(":date")
	tags := receiver.GetString("tags", "")
	province := receiver.GetString("province", "")
	sortStr := receiver.GetString("sort", "")
	orderBy := receiver.GetString("order_by", "")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	dateTime, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	ret := map[string]interface{}{}
	cacheKey := cache.GetCacheKey(cache.DyRankCache, "author_follower_inc", utils.Md5_encode(fmt.Sprintf("%s%s%s%s%s%d%d", dateTime.Format("20060102"), tags, province, sortStr, orderBy, page, pageSize)))
	cacheStr := global.Cache.Get(cacheKey)
	if cacheStr != "" {
		cacheStr = utils.DeserializeData(cacheStr)
		_ = jsoniter.Unmarshal([]byte(cacheStr), &ret)
	} else {
		data, total, comErr := es.NewEsAuthorBusiness().DyAuthorFollowerIncRank(dateTime.Format("20060102"), tags, province, sortStr, orderBy, page, pageSize)
		var structData []es2.DyAuthorFollowerTop
		utils.MapToStruct(data, &structData)
		for k := range structData {
			if structData[k].UniqueID == "" {
				structData[k].UniqueID = structData[k].ShortID
			}
			structData[k].AuthorID = business.IdEncrypt(structData[k].AuthorID)
			structData[k].AuthorCover = dyimg.Fix(structData[k].AuthorCover)
		}
		if comErr != nil {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		ret = map[string]interface{}{
			"list":  data,
			"total": total,
		}
		if dateTime.Format("20060102") != time.Now().Format("20060102") {
			_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 86400)
		}
	}
	receiver.SuccReturn(ret)
	return
}
