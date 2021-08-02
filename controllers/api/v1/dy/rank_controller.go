package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/models/hbase"
)

type RankController struct {
	controllers.ApiBaseController
}

//抖音视频达人热榜
func (receiver *RankController) DyStartAuthorVideoRank() {
	rankType := receiver.GetString("rank_type", "达人指数榜")
	category := receiver.GetString("category", "全部")
	data, _ := hbase.GetStartAuthorVideoRank(rankType, category)
	receiver.SuccReturn(data)
	return
}

//抖音直播达人热榜
func (receiver *RankController) DyStartAuthorLiveRank() {
	rankType := receiver.GetString("rank_type", "达人指数榜")
	data, _ := hbase.GetStartAuthorLiveRank(rankType)
	receiver.SuccReturn(data)
	return
}
