package business

import (
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"math"
	"time"
)

type AwemeBusiness struct {
}

func NewAwemeBusiness() *AwemeBusiness {
	return new(AwemeBusiness)
}

//视频详情
func (a *AwemeBusiness) HbaseGetAweme(awemeId string) (data entity.DyAwemeData, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAweme).GetByRowKey([]byte(awemeId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	authorMap := hbaseService.HbaseFormat(result, entity.DyAwemeMap)
	aweme := &entity.DyAweme{}
	utils.MapToStruct(authorMap, aweme)
	aweme.AwemeID = aweme.Data.ID
	duration := math.Ceil(float64(aweme.Data.Duration) / 1000)
	aweme.Data.Duration = utils.ToInt(duration)
	aweme.Data.AwemeTitle = aweme.AwemeTitle
	data = aweme.Data
	return
}

//视频某天数据
func (a *AwemeBusiness) HbaseGetAwemeCountDataRangeDate(awemeId string, startTime, endTime time.Time) (data map[string]entity.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := awemeId + "_" + startTime.Format("20060102")
	endRow := awemeId + "_" + endTime.AddDate(0, 0, 1).Format("20060102")
	results, err := query.
		SetTable(hbaseService.HbaseDyAwemeDiggCommentForwardCount).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	data = map[string]entity.DyAwemeDiggCommentForwardCount{}
	for _, v := range results {
		dataMap := hbaseService.HbaseFormat(v, entity.DyAwemeDiggCommentForwardCountMap)
		hData := entity.DyAwemeDiggCommentForwardCount{}
		utils.MapToStruct(dataMap, &hData)
		t := time.Unix(hData.CrawlTime, 0)
		date := t.Format("20060102")
		data[date] = hData
	}
	return
}

//获取视频每天详情数据
func (a *AwemeBusiness) HbaseGetAwemeCountData(awemeId, date string) (data entity.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := awemeId + "_" + date
	result, err := query.
		SetTable(hbaseService.HbaseDyAwemeDiggCommentForwardCount).
		GetByRowKey([]byte(startRow))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity.DyAwemeDiggCommentForwardCountMap)
	utils.MapToStruct(dataMap, &data)
	return
}

//获取视频趋势数据
func (a *AwemeBusiness) GetAwemeChart(awemeId string, startTime, endTime time.Time, beforeGet bool) (data map[string]entity.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	data, comErr = a.HbaseGetAwemeCountDataRangeDate(awemeId, startTime, endTime)
	start := startTime.Format("20060102")
	end := endTime.Format("20060102")
	if comErr == nil {
		yesterday := startTime.AddDate(0, 0, -1).Format("20060102")
		if beforeGet {
			beforeData, _ := a.HbaseGetAwemeCountData(awemeId, yesterday)
			data[yesterday] = beforeData
		}
		//首发补点
		if _, ok := data[start]; !ok {
			if beforeData, ok := data[yesterday]; !ok {
				data[start] = beforeData
			} else {
				data[start] = entity.DyAwemeDiggCommentForwardCount{}
			}
		}
		//末尾补点
		if _, ok := data[end]; !ok {
			awemeBusiness := NewAwemeBusiness()
			awemeBase, _ := awemeBusiness.HbaseGetAweme(awemeId)
			data[end] = entity.DyAwemeDiggCommentForwardCount{
				DiggCount:    awemeBase.DiggCount,
				CommentCount: awemeBase.CommentCount,
				ForwardCount: awemeBase.ForwardCount,
			}
		}
		beginDatetime := startTime
		beforeDay := startTime.Format("20060102")
		//空数据补点
		for {
			if beginDatetime.After(endTime) {
				break
			}
			today := beginDatetime.Format("20060102")
			//数据不存在向前补点
			if _, ok := data[today]; !ok {
				data[today] = data[beforeDay]
			}
			beforeDay = today
			beginDatetime = beginDatetime.AddDate(0, 0, 1)
		}
	}
	return
}
