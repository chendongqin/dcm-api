package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/models/business"
)

type RankController struct {
	controllers.ApiBaseController
}

func (receiver *RankController) DyStartAuthorVideoRank() {
	rankType := receiver.GetString("rank_type", "达人指数榜")
	category := receiver.GetString("category", "全部")
	rankBusiness := business.NewRankBusiness()
	data, _ := rankBusiness.HbaseStartAuthorVideoRank(rankType, category)
	receiver.SuccReturn(data)
	return
}
