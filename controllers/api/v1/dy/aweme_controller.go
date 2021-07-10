package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/models/business"
	"dongchamao/structinit/repost/dy"
	"time"
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
	awemeBusiness := business.NewAwemeBusiness()
	awemeBase, comErr := awemeBusiness.HbaseGetAweme(awemeId)
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
	startDay := receiver.Ctx.Input.Param(":start")
	endDay := receiver.Ctx.Input.Param(":end")
	if awemeId == "" || startDay == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if endDay == "" {
		endDay = time.Now().Format("2006-01-02")
	}
	pslTime := "2006-01-02"
	t1, err := time.ParseInLocation(pslTime, startDay, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	t2, err := time.ParseInLocation(pslTime, endDay, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if t1.After(t2) || t2.After(t1.AddDate(0, 0, 90)) || t2.After(time.Now()) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	awemeBusiness := business.NewAwemeBusiness()
	awemeCount, comErr := awemeBusiness.GetAwemeChart(awemeId, t1.Format("20060102"), t2.Format("20060102"), true)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	//前一天数据，做增量计算
	beginDatetime := t1
	beforeData := entity.DyAwemeDiggCommentForwardCount{}
	beforeDay := beginDatetime.AddDate(0, 0, -1).Format("20060102")
	if _, ok := awemeCount[beforeDay]; ok {
		beforeData = awemeCount[beforeDay]
	} else {
		beforeData, _ = awemeBusiness.HbaseGetAwemeCountData(awemeId, beforeDay)
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
