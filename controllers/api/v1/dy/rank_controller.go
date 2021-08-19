package dy

import (
	"dongchamao/business"
	"dongchamao/business/es"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	"dongchamao/services/dyimg"
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
	data, updateTime, _ := hbase.GetStartAuthorVideoRank(rankType, category)
	for k, v := range data {
		data[k].CoreUserId = business.IdEncrypt(v.CoreUserId)
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":        data,
		"update_time": updateTime,
	})
	return
}

//抖音直播达人热榜
func (receiver *RankController) DyStartAuthorLiveRank() {
	rankType := receiver.GetString("rank_type", "达人指数榜")
	data, updateTime, _ := hbase.GetStartAuthorLiveRank(rankType)
	for k, v := range data {
		data[k].CoreUserId = business.IdEncrypt(v.CoreUserId)
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":        data,
		"update_time": updateTime,
	})
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
	receiver.SuccReturn(map[string]interface{}{
		"list":        data.Ranks,
		"update_time": data.CrawlTime,
	})
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
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":        data.Ranks,
		"update_time": data.CrawlTime,
	})
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
		data.Ranks[k].LiveInfo.User.Id = business.IdEncrypt(v.LiveInfo.User.Id)
		data.Ranks[k].RoomId = business.IdEncrypt(v.RoomId)
		data.Ranks[k].LiveInfo.Cover = dyimg.Fix(v.LiveInfo.Cover)
		data.Ranks[k].LiveInfo.User.Avatar = dyimg.Fix(v.LiveInfo.User.Avatar)
		if v.LiveInfo.User.DisplayId == "" {
			data.Ranks[k].LiveInfo.User.DisplayId = v.LiveInfo.User.ShortId
		}
		data.Ranks[k].ShareUrl = business.LiveShareUrl + v.RoomId
		if v.RealGmv > 0 {
			data.Ranks[k].PredictGmv = v.RealGmv
			data.Ranks[k].PredictSales = v.RealSales
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":        data.Ranks,
		"update_time": data.CrawlTime,
	})
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
	receiver.SuccReturn(map[string]interface{}{
		"list":        data.Ranks,
		"update_time": data.CrawlTime,
	})
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
	data, _ := hbase.GetLiveShareWeekRank(start.Format("20060102") + "_" + end.Format("20060102"))
	list := make([]entity.DyLiveShareWeekData, 0)
	for _, v := range data.Data {
		var gmv float64 = 0
		var sales float64 = 0
		var totalUser int64 = 0
		for _, r := range v.Rooms {
			if r.RealSales > 0 {
				gmv += r.RealGmv
				sales += math.Floor(r.RealSales)
			} else {
				gmv += r.PredictGmv
				sales += math.Floor(r.PredictSales)
			}
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
	receiver.SuccReturn(map[string]interface{}{
		"list":        list,
		"update_time": data.CrawlTime,
	})
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
	receiver.SuccReturn(map[string]interface{}{
		"list":        list,
		"update_time": data.CrawlTime,
	})
	receiver.SuccReturn(data)
	return
}

//达人带货榜
func (receiver *RankController) DyAuthorTakeGoodsRank() {
	date := receiver.GetString("date", "*")
	tags := receiver.GetString("tags", "综合")
	list, _ := es.NewEsAuthorBusiness().AuthorTakeGoodsRank(date, tags)
	for k, v := range list {
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		if v.UniqueId == "" || v.UniqueId == "0" {
			list[k].UniqueId = v.ShortId
		}
		list[k].AuthorCover = dyimg.Avatar(v.AuthorCover)
		list[k].RoomCover = dyimg.Avatar(v.RoomCover)
	}
	receiver.SuccReturn(map[string]interface{}{
		"list": list,
	})
	return
}
