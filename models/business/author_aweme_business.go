package business

import (
	"dongchamao/entity"
	"dongchamao/global"
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

//获取视频概览数据
func (a *AuthorAwemeBusiness) HbaseGetVideoAggRangeDate(authorId, startDate, endDate string) (data dy.AuthorVideoOverview) {
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
	allAwemeChan := make(chan map[string]map[string]entity.DyAwemeDiggCommentForwardCount, 0)
	allAwemeData := map[string]map[string]entity.DyAwemeDiggCommentForwardCount{}
	awemeIds := make([]string, 0)
	for _, v := range results {
		dataMap := hbaseService.HbaseFormat(v, entity.DyAuthorAwemeAggMap)
		hData := &entity.DyAuthorAwemeAggData{}
		utils.MapToStruct(dataMap, hData)
		for _, agg := range hData.Data {
			awemeIds = append(awemeIds, agg.AwemeID)
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
			//视频趋势数据处理
			createTime := time.Unix(agg.AwemeCreateTime, 0)
			go func(ch chan map[string]map[string]entity.DyAwemeDiggCommentForwardCount, awemeId, start, end string) {
				awemeBusiness := NewAwemeBusiness()
				awemeData, comErr := awemeBusiness.GetAwemeChart(awemeId, start, end, false)
				if comErr == nil {
					allAwemeDataMap := map[string]map[string]entity.DyAwemeDiggCommentForwardCount{}
					allAwemeDataMap[awemeId] = awemeData
					ch <- allAwemeDataMap
				}
			}(allAwemeChan, agg.AwemeID, createTime.Format("20060102"), endDate)
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
		for k, v := range tmp {
			allAwemeData[k] = v
		}
	}
	t1, _ := time.ParseInLocation("20060102", startDate, time.Local)
	t2, _ := time.ParseInLocation("20060102", endDate, time.Local)
	//前一天数据，做增量计算
	beginDatetime := t1
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
		if beginDatetime.After(t2) {
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

//粉丝某天数据
func (a *AuthorAwemeBusiness) HbaseGetAwemeCountDataRangeDate(awemeId, startDate, endDate string) (data map[string]entity.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := awemeId + "_" + startDate
	endRow := awemeId + "_" + endDate
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
func (a *AuthorAwemeBusiness) HbaseGetAwemeCountData(awemeId, date string) (data entity.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
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
