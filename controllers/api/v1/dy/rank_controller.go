package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/models/business"
)

type RankController struct {
	controllers.ApiBaseController
}

//抖音视频达人热榜
func (receiver *RankController) DyStartAuthorVideoRank() {
	rankType := receiver.GetString("rank_type", "达人指数榜")
	category := receiver.GetString("category", "全部")
	rankBusiness := business.NewRankBusiness()
	data, _ := rankBusiness.HbaseStartAuthorVideoRank(rankType, category)
	receiver.SuccReturn(data)
	return
}

//抖音直播达人热榜
func (receiver *RankController) DyStartAuthorLiveRank() {
	rankType := receiver.GetString("rank_type", "达人指数榜")
	rankBusiness := business.NewRankBusiness()
	data, _ := rankBusiness.HbaseStartAuthorLiveRank(rankType)
	receiver.SuccReturn(data)
	return
}
