package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/models/business"
	"dongchamao/models/hbase"
	entity2 "dongchamao/models/hbase/entity"
	"dongchamao/structinit/repost/dy"
)

type AwemeController struct {
	controllers.ApiBaseController
}

func (receiver *AwemeController) AwemeBaseData() {
	awemeId := receiver.Ctx.Input.Param(":aweme_id")
	if awemeId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	awemeBase, comErr := hbase.GetVideo(awemeId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	awemeSimple := dy.DySimpleAweme{
		AuthorID:        awemeBase.AuthorID,
		AwemeCover:      awemeBase.AwemeCover,
		AwemeTitle:      awemeBase.AwemeTitle,
		AwemeCreateTime: awemeBase.AwemeCreateTime,
		AwemeURL:        awemeBase.AwemeURL,
		CommentCount:    awemeBase.CommentCount,
		DiggCount:       awemeBase.DiggCount,
		DownloadCount:   awemeBase.DownloadCount,
		Duration:        awemeBase.Duration,
		ForwardCount:    awemeBase.ForwardCount,
		ID:              awemeBase.ID,
		MusicID:         awemeBase.MusicID,
		ShareCount:      awemeBase.ShareCount,
		PromotionNum:    len(awemeBase.DyPromotionID),
	}
	receiver.SuccReturn(map[string]interface{}{
		"aweme_base": awemeSimple,
	})
	return
}

func (receiver *AwemeController) AwemeChart() {
	awemeId := receiver.Ctx.Input.Param(":aweme_id")
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
	beforeData := entity2.DyAwemeDiggCommentForwardCount{}
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
		"digg": dy.DateChart{
			Date:       dateArr,
			CountValue: diggCountArr,
			IncValue:   diggIncArr,
		},
		"forward": dy.DateChart{
			Date:       dateArr,
			CountValue: forwardCountArr,
			IncValue:   forwardIncArr,
		},
		"comment": dy.DateChart{
			Date:       dateArr,
			CountValue: commentCountArr,
			IncValue:   commentIncArr,
		},
	}
	receiver.SuccReturn(returnMap)
	return
}
