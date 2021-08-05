package dy

import (
	"dongchamao/business"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/hbase"
	"dongchamao/services/dyimg"
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
	receiver.SuccReturn(map[string]interface{}{
		"list":        data,
		"update_time": updateTime,
	})
	return
}

//抖音直播达人热榜
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
