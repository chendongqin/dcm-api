package dy

import (
	"dongchamao/business"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	dy2 "dongchamao/models/repost/dy"
	"sort"
)

type AwemeController struct {
	controllers.ApiBaseController
}

func (receiver *AwemeController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseShowNum)
}

func (receiver *AwemeController) AwemeBaseData() {
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	if awemeId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	awemeBase, comErr := hbase.GetVideo(awemeId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	awemeSimple := dy2.DySimpleAweme{
		AuthorID:        awemeBase.Data.AuthorID,
		AwemeCover:      awemeBase.Data.AwemeCover,
		AwemeTitle:      awemeBase.Data.AwemeTitle,
		AwemeCreateTime: awemeBase.Data.AwemeCreateTime,
		AwemeURL:        awemeBase.Data.AwemeURL,
		CommentCount:    awemeBase.Data.CommentCount,
		DiggCount:       awemeBase.Data.DiggCount,
		DownloadCount:   awemeBase.Data.DownloadCount,
		Duration:        awemeBase.Data.Duration,
		ForwardCount:    awemeBase.Data.ForwardCount,
		ID:              awemeBase.Data.ID,
		MusicID:         awemeBase.Data.MusicID,
		ShareCount:      awemeBase.Data.ShareCount,
		PromotionNum:    len(awemeBase.Data.DyPromotionID),
	}
	receiver.SuccReturn(map[string]interface{}{
		"aweme_base": awemeSimple,
	})
	return
}

func (receiver *AwemeController) AwemeChart() {
	awemeId := business.IdEncrypt(receiver.Ctx.Input.Param(":aweme_id"))
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
	beforeData := entity.DyAwemeDiggCommentForwardCount{}
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
		"digg": dy2.DateChart{
			Date:       dateArr,
			CountValue: diggCountArr,
			IncValue:   diggIncArr,
		},
		"forward": dy2.DateChart{
			Date:       dateArr,
			CountValue: forwardCountArr,
			IncValue:   forwardIncArr,
		},
		"comment": dy2.DateChart{
			Date:       dateArr,
			CountValue: commentCountArr,
			IncValue:   commentIncArr,
		},
	}
	receiver.SuccReturn(returnMap)
	return
}

func (receiver *AwemeController) AwemeCommentHotWords() {
	awemeId := business.IdDecrypt(receiver.Ctx.Input.Param(":aweme_id"))
	if awemeId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	awemeBase, comErr := hbase.GetVideo(awemeId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	list := make([]dy2.NameValueInt64Chart, 0)
	for k, v := range awemeBase.HotWordShow {
		list = append(list, dy2.NameValueInt64Chart{
			Name:  k,
			Value: v,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Value > list[j].Value
	})
	receiver.SuccReturn(map[string]interface{}{
		"hot_words": list,
	})
	return
}
