package business

import (
	"dongchamao/global"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	"time"
)

type AwemeBusiness struct {
}

func NewAwemeBusiness() *AwemeBusiness {
	return new(AwemeBusiness)
}

//获取视频趋势数据
func (a *AwemeBusiness) GetAwemeChart(awemeId string, startTime, endTime time.Time, beforeGet bool) (data map[string]entity.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	data, comErr = hbase.GetVideoCountDataRangeDate(awemeId, startTime, endTime)
	start := startTime.Format("20060102")
	end := endTime.Format("20060102")
	if comErr == nil {
		yesterday := startTime.AddDate(0, 0, -1).Format("20060102")
		if beforeGet {
			beforeData, _ := hbase.GetVideoCountData(awemeId, yesterday)
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
			awemeBase, _ := hbase.GetVideo(awemeId)
			data[end] = entity.DyAwemeDiggCommentForwardCount{
				DiggCount:    awemeBase.Data.DiggCount,
				CommentCount: awemeBase.Data.CommentCount,
				ForwardCount: awemeBase.Data.ForwardCount,
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
