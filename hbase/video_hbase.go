package hbase

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	entity2 "dongchamao/models/entity"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/hbaseService/hbasehelper"
	"math"
	"time"
)

//视频详情
func GetVideo(awemeId string) (data entity2.DyAwemeData, comErr global.CommonError) {
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
	authorMap := hbaseService.HbaseFormat(result, entity2.DyAwemeMap)
	aweme := &entity2.DyAweme{}
	utils.MapToStruct(authorMap, aweme)
	duration := math.Ceil(float64(aweme.Data.Duration) / 1000)
	data = aweme.Data
	data.Duration = utils.ToInt(duration)
	data.AwemeTitle = aweme.AwemeTitle
	data.AwemeCover = dyimg.Fix(data.AwemeCover)
	return
}

//视频某天数据
func GetVideoCountDataRangeDate(awemeId string, startTime, endTime time.Time) (data map[string]entity2.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
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
	data = map[string]entity2.DyAwemeDiggCommentForwardCount{}
	for _, v := range results {
		dataMap := hbaseService.HbaseFormat(v, entity2.DyAwemeDiggCommentForwardCountMap)
		hData := entity2.DyAwemeDiggCommentForwardCount{}
		utils.MapToStruct(dataMap, &hData)
		t := time.Unix(hData.CrawlTime, 0)
		date := t.Format("20060102")
		data[date] = hData
	}
	return
}

//获取视频每天详情数据
func GetVideoCountData(awemeId, date string) (data entity2.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := awemeId + "_" + date
	result, err := query.
		SetTable(hbaseService.HbaseDyAwemeDiggCommentForwardCount).
		GetByRowKey([]byte(startRow))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity2.DyAwemeDiggCommentForwardCountMap)
	utils.MapToStruct(dataMap, &data)
	return
}

//获取视频概览数据
func GetAuthorVideoAggRangeDate(authorId string, startTime, endTime time.Time) (results []*hbase.TResult_, err error) {
	query := hbasehelper.NewQuery()
	startRow := authorId + "_" + startTime.Format("20060102")
	endRow := authorId + "_" + endTime.AddDate(0, 0, 1).Format("20060102")
	results, err = query.
		SetTable(hbaseService.HbaseDyAuthorAwemeAgg).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		return
	}
	return
}

//视频某天数据
func GetAuthorVideoCountDataRangeDate(awemeId, startDate, endDate string) (data map[string]entity2.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
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
	data = map[string]entity2.DyAwemeDiggCommentForwardCount{}
	for _, v := range results {
		dataMap := hbaseService.HbaseFormat(v, entity2.DyAwemeDiggCommentForwardCountMap)
		hData := entity2.DyAwemeDiggCommentForwardCount{}
		utils.MapToStruct(dataMap, &hData)
		t := time.Unix(hData.CrawlTime, 0)
		date := t.Format("20060102")
		data[date] = hData
	}
	return
}

//获取视频每天详情数据
func GetAuthorVideoCountData(awemeId, date string) (data entity2.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := awemeId + "_" + date
	result, err := query.
		SetTable(hbaseService.HbaseDyAwemeDiggCommentForwardCount).
		GetByRowKey([]byte(startRow))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity2.DyAwemeDiggCommentForwardCountMap)
	utils.MapToStruct(dataMap, &data)
	return
}
