package business

import (
	"dongchamao/entity"
	"dongchamao/global/utils"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"dongchamao/structinit/repost/dy"
	"fmt"
	"time"
)

type AuthorAwemeBusiness struct {
}

func NewAuthorAwemeBusiness() *AuthorAwemeBusiness {
	return new(AuthorAwemeBusiness)
}

func (a *AuthorAwemeBusiness) HbaseGetVideoAgg(authorId, startDate, endDate string) (data dy.AuthorVideoOverview) {
	query := hbasehelper.NewQuery()
	startRow := authorId + "_" + startDate
	endRow := authorId + "_" + endDate
	results, err := query.
		SetTable(hbaseService.HbaseDyAuthorAwemeAgg).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		return
	}
	var videoNum, productVideo int64
	var diggMax, diggMin, diggCount, commentMax, commentMin, commentCount, forwardMax, forwardMin, forwardCount int64
	durationChartMap := map[string]int{
		"up_120": 0,
		"up_60":  0,
		"up_30":  0,
		"up_15":  0,
		"up_0":   0,
	}
	publishChartMap := map[string]int{}
	for _, v := range results {
		dataMap := hbaseService.HbaseFormat(v, entity.DyAuthorAwemeAggMap)
		hData := &entity.DyAuthorAwemeAggData{}
		utils.MapToStruct(dataMap, hData)
		for _, agg := range hData.Data {
			aggCreateTime := time.Unix(agg.AwemeCreateTime, 0)
			hour := aggCreateTime.Format("15")
			videoNum++
			if agg.DyPromotionId != "0" {
				productVideo++
			}
			diggCount += agg.DiggCount
			commentCount += agg.CommentCount
			forwardCount += agg.ForwardCount
			//map处理
			//时长时间
			var durationLab string
			if agg.Duration > 120000 {
				durationLab = "up_120"
			} else if agg.Duration > 60000 {
				durationLab = "up_60"
			} else if agg.Duration > 30000 {
				durationLab = "up_30"
			} else if agg.Duration > 15000 {
				durationLab = "up_15"
			} else {
				durationLab = "up_0"
			}
			durationChartMap[durationLab] += 1
			//发布时间
			if _, ok := publishChartMap[hour]; ok {
				publishChartMap[hour] += 1
			} else {
				publishChartMap[hour] = 1
			}
			//峰值
			if agg.DiggCount > diggMax {
				diggMax = agg.DiggCount
			}
			if diggMin == 0 || diggMin > agg.DiggCount {
				diggMin = agg.DiggCount
			}
			if agg.CommentCount > commentMax {
				commentMax = agg.CommentCount
			}
			if commentMin == 0 || commentMin > agg.CommentCount {
				commentMin = agg.CommentCount
			}
			if agg.ForwardCount > forwardMax {
				forwardMax = agg.ForwardCount
			}
			if forwardMin == 0 || forwardMin > agg.ForwardCount {
				forwardMin = agg.ForwardCount
			}
		}
	}
	if videoNum > 0 {
		data.AvgDiggCount = diggCount / videoNum
		data.AvgCommentCount = commentCount / videoNum
		data.AvgForwardCount = forwardCount / videoNum
	}
	data.DiggMax = diggMax
	data.DiggMin = diggMin
	data.CommentMax = commentMax
	data.CommentMin = commentMin
	data.ForwardMax = forwardMax
	data.ForwardMin = forwardMin
	data.VideoNum = videoNum
	data.ProductVideo = productVideo
	for k, v := range durationChartMap {
		data.DurationChart = append(data.DurationChart, dy.VideoChart{
			Name:  k,
			Value: v,
		})
	}
	for i := 0; i <= 23; i++ {
		hour := fmt.Sprintf("%02d", i)
		if _, ok := publishChartMap[hour]; !ok {
			publishChartMap[hour] = 0
		}
		data.PublishChart = append(data.PublishChart, dy.VideoChart{
			Name:  hour,
			Value: publishChartMap[hour],
		})
	}
	//总增量图表
	//t1, _ := time.ParseInLocation("20060102", startDate, time.Local)
	//t2, _ := time.ParseInLocation("20060102", endDate, time.Local)
	//beginDatetime := t1
	//for {
	//	if beginDatetime.After(t2) {
	//		break
	//	}
	//	beginDatetime = beginDatetime.AddDate(0,0,1)
	//}
	return
}
