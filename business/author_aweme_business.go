package business

import (
	"dongchamao/business/es"
	"dongchamao/global"
	"dongchamao/models/entity"
	"dongchamao/models/repost/dy"
	"fmt"
	"sync"
	"time"
)

type AuthorAwemeBusiness struct {
}

func NewAuthorAwemeBusiness() *AuthorAwemeBusiness {
	return new(AuthorAwemeBusiness)
}

//获取视频概览数据
func (a *AuthorAwemeBusiness) HbaseGetVideoAggRangeDate(authorId string, startTime, endTime time.Time) (data dy.AuthorVideoOverview) {
	//results, err := hbase.GetAuthorVideoAggRangeDate(authorId, startTime, endTime)
	//if err != nil {
	//	return
	//}
	var videoNum, productVideo int64
	esVideoBusiness := es.NewEsVideoBusiness()
	results, total, _ := esVideoBusiness.SearchByAuthor(authorId, "", "", "", 0, 1, 10000, startTime, endTime)
	videoNum = int64(total)
	var diggMax, diggMin, diggCount, commentMax, commentMin, commentCount, forwardMax, forwardMin, forwardCount int64
	durationChartMap := map[string]int{
		"up_120": 0,
		"up_60":  0,
		"up_30":  0,
		"up_15":  0,
		"up_0":   0,
	}
	publishChartMap := map[string]int{}
	allAwemeChan := make(chan map[string]map[string]entity.DyAwemeDiggCommentForwardCount, len(results))
	allAwemeData := map[string]map[string]entity.DyAwemeDiggCommentForwardCount{}
	awemeIds := make([]string, 0)
	var wg sync.WaitGroup
	for _, v := range results {
		//dataMap := hbaseService.HbaseFormat(v, entity.DyAuthorAwemeAggMap)
		//hData := &entity.DyAuthorAwemeAggData{}
		//utils.MapToStruct(dataMap, hData)
		//for _, agg := range hData.Data {
		awemeIds = append(awemeIds, v.AwemeId)
		aggCreateTime := time.Unix(v.AwemeCreateTime, 0)
		hour := aggCreateTime.Format("15")
		if v.ProductIds != "" {
			productVideo++
		}
		diggCount += v.DiggCount
		commentCount += v.CommentCount
		forwardCount += v.ShareCount
		//map处理
		//时长时间
		var durationLab string
		if v.Duration > 120000 {
			durationLab = "up_120"
		} else if v.Duration > 60000 {
			durationLab = "up_60"
		} else if v.Duration > 30000 {
			durationLab = "up_30"
		} else if v.Duration > 15000 {
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
		if v.DiggCount > diggMax {
			diggMax = v.DiggCount
		}
		if diggMin == 0 || diggMin > v.DiggCount {
			diggMin = v.DiggCount
		}
		if v.CommentCount > commentMax {
			commentMax = v.CommentCount
		}
		if commentMin == 0 || commentMin > v.CommentCount {
			commentMin = v.CommentCount
		}
		if v.ShareCount > forwardMax {
			forwardMax = v.ShareCount
		}
		if forwardMin == 0 || forwardMin > v.ShareCount {
			forwardMin = v.ShareCount
		}
		//视频趋势数据处理
		createTime := time.Unix(v.AwemeCreateTime, 0)
		wg.Add(1)
		go func(awemeId string, startT, endT time.Time) {
			defer global.RecoverPanic()
			wg.Done()
			awemeBusiness := NewAwemeBusiness()
			awemeData, comErr := awemeBusiness.GetAwemeChart(awemeId, startT, endT, false)
			if comErr == nil {
				tmp := map[string]map[string]entity.DyAwemeDiggCommentForwardCount{
					awemeId: awemeData,
				}
				allAwemeChan <- tmp
			} else {
				allAwemeChan <- nil
			}
		}(v.AwemeId, createTime, endTime)
		//}
	}
	wg.Wait()
	//videoSumData := esVideoBusiness.SumDataByAuthor(authorId, startTime, endTime)
	if videoNum > 0 {
		data.AvgDiggCount = diggCount / videoNum
		data.AvgCommentCount = commentCount / videoNum
		data.AvgForwardCount = forwardCount / videoNum
	}
	data.VideoNum = videoNum
	data.ProductVideo = productVideo
	//data.ProductVideo, _ = esVideoBusiness.CountAwemeByAuthor(authorId, 1, startTime, endTime)
	//data.AvgDiggCount = videoSumData.AvgDigg
	//data.AvgCommentCount = videoSumData.AvgComment
	//data.AvgForwardCount = videoSumData.AvgShare
	//data.VideoNum = int64(videoSumData.Total)
	data.DiggMax = diggMax
	data.DiggMin = diggMin
	data.CommentMax = commentMax
	data.CommentMin = commentMin
	data.ForwardMax = forwardMax
	data.ForwardMin = forwardMin
	for k, v := range durationChartMap {
		data.DurationChart = append(data.DurationChart, dy.VideoChart{
			Name:  k,
			Value: v,
		})
	}
	//小时数据
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
	for i := 0; i < int(videoNum); i++ {
		tmp, ok := <-allAwemeChan
		if !ok {
			break
		}
		if tmp == nil {
			continue
		}
		for k, v := range tmp {
			allAwemeData[k] = v
		}
	}
	//前一天数据，做增量计算
	beginDatetime := startTime
	beforeSumData := entity.DyAwemeDiggCommentForwardCount{}
	beforeDay := beginDatetime.AddDate(0, 0, -1).Format("20060102")
	for _, awemeId := range awemeIds {
		if v, ok := allAwemeData[awemeId][beforeDay]; ok {
			beforeSumData.DiggCount += v.DiggCount
			beforeSumData.CommentCount += v.CommentCount
			beforeSumData.ForwardCount += v.ForwardCount
		}
	}
	dateArr := make([]string, 0)
	diggCountArr := make([]int64, 0)
	commentCountArr := make([]int64, 0)
	forwardCountArr := make([]int64, 0)
	diggIncArr := make([]int64, 0)
	commentIncArr := make([]int64, 0)
	forwardIncArr := make([]int64, 0)
	for {
		if beginDatetime.After(endTime) {
			break
		}
		sumData := entity.DyAwemeDiggCommentForwardCount{}
		date := beginDatetime.Format("20060102")
		for _, awemeId := range awemeIds {
			if v, ok := allAwemeData[awemeId][date]; ok {
				sumData.DiggCount += v.DiggCount
				sumData.CommentCount += v.CommentCount
				sumData.ForwardCount += v.ForwardCount
			}
		}
		dateArr = append(dateArr, beginDatetime.Format("01/02"))
		diggCountArr = append(diggCountArr, sumData.DiggCount)
		commentCountArr = append(commentCountArr, sumData.CommentCount)
		forwardCountArr = append(forwardCountArr, sumData.ForwardCount)
		diggIncArr = append(diggIncArr, sumData.DiggCount-beforeSumData.DiggCount)
		commentIncArr = append(commentIncArr, sumData.CommentCount-beforeSumData.CommentCount)
		forwardIncArr = append(forwardIncArr, sumData.ForwardCount-beforeSumData.ForwardCount)
		beforeSumData = sumData
		beginDatetime = beginDatetime.AddDate(0, 0, 1)
	}
	data.DiggChart = dy.DateChart{
		Date:       dateArr,
		CountValue: diggCountArr,
		IncValue:   diggIncArr,
	}
	data.ForwardChart = dy.DateChart{
		Date:       dateArr,
		CountValue: forwardCountArr,
		IncValue:   forwardIncArr,
	}
	data.CommentChart = dy.DateChart{
		Date:       dateArr,
		CountValue: commentCountArr,
		IncValue:   commentIncArr,
	}
	return
}
