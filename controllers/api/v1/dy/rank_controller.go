package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/hbase"
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
