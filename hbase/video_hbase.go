package hbase

import (
	"context"
	"dongchamao/global"
	"dongchamao/global/logger"
	"dongchamao/global/utils"
	"dongchamao/models/entity"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/hbaseService/hbasehelper"
	"encoding/json"
	"math"
	"sort"
	"strings"
	"time"
)

func GetVideoByIds(awemeIds []string) (map[string]entity.DyAweme, error) {
	rowKeys := make([]*hbase.TGet, 0)
	for _, id := range awemeIds {
		rowKeys = append(rowKeys, &hbase.TGet{Row: []byte(id)})
	}
	client := global.HbasePools.Get("default")
	defer client.Close()
	results, err := client.GetMultiple(context.Background(), []byte(hbaseService.HbaseDyAweme), rowKeys)
	if err != nil {
		return nil, err
	}
	infoMap := map[string]entity.DyAweme{}
	for _, v := range results {
		data := entity.DyAweme{}
		detailMap := hbaseService.HbaseFormat(v, entity.DyAwemeMap)
		utils.MapToStruct(detailMap, &data)
		duration := math.Ceil(float64(data.Data.Duration) / 1000)
		data.Data.Duration = utils.ToInt(duration)
		data.AwemeID = string(v.Row)
		data.AwemeTitle = data.Data.AwemeTitle
		data.Data.CrawlTime = data.CrawlTime
		data.Data.AwemeCover = dyimg.Fix(data.Data.AwemeCover)
		infoMap[data.AwemeID] = data
	}
	return infoMap, nil
}

//视频详情
func GetVideo(awemeId string) (data entity.DyAweme, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAweme).GetByRowKey([]byte(awemeId))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyAwemeMap)
	utils.MapToStruct(detailMap, &data)
	duration := math.Ceil(float64(data.Data.Duration) / 1000)
	data.Data.Duration = utils.ToInt(duration)
	data.AwemeID = string(result.Row)
	data.AwemeTitle = data.Data.AwemeTitle
	data.Data.CrawlTime = data.CrawlTime
	data.Data.AwemeCover = dyimg.Fix(data.Data.AwemeCover)
	return
}

//视频某天数据
func GetVideoCountDataRangeDate(awemeId string, startTime, endTime time.Time) (data map[string]entity.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := awemeId + "_" + startTime.Format("20060102")
	endRow := awemeId + "_" + endTime.AddDate(0, 0, 1).Format("20060102")
	results, err := query.
		SetTable(hbaseService.HbaseDyAwemeDiggCommentForwardCount).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
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
func GetVideoCountData(awemeId, date string) (data entity.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := awemeId + "_" + date
	result, err := query.
		SetTable(hbaseService.HbaseDyAwemeDiggCommentForwardCount).
		GetByRowKey([]byte(startRow))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity.DyAwemeDiggCommentForwardCountMap)
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
func GetAuthorVideoCountDataRangeDate(awemeId, startDate, endDate string) (data map[string]entity.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := awemeId + "_" + startDate
	endRow := awemeId + "_" + endDate
	results, err := query.
		SetTable(hbaseService.HbaseDyAwemeDiggCommentForwardCount).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
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
func GetAuthorVideoCountData(awemeId, date string) (data entity.DyAwemeDiggCommentForwardCount, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := awemeId + "_" + date
	result, err := query.
		SetTable(hbaseService.HbaseDyAwemeDiggCommentForwardCount).
		GetByRowKey([]byte(startRow))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity.DyAwemeDiggCommentForwardCountMap)
	utils.MapToStruct(dataMap, &data)
	return
}

//获取视频评论列表
func GetAwemeTopComment(awemeId string, start, end int) (data []entity.DyAwemeCommentTop, total int, comErr global.CommonError) {
	data = make([]entity.DyAwemeCommentTop, 0)
	query := hbasehelper.NewQuery()
	result, err := query.
		SetTable(hbaseService.HbaseDyAwemeTopComment).
		GetByRowKey([]byte(awemeId))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity.DyAwemeCommentTopMap)
	var commentStruct entity.DyAwemeCommentTopStruct
	utils.MapToStruct(dataMap, &commentStruct)
	if commentStruct.DiggInfo != "" {
		commentStruct.DiggInfo = "[" + strings.Replace(commentStruct.DiggInfo, "=----=", ",", -1) + "]"
		_ = json.Unmarshal([]byte(commentStruct.DiggInfo), &data)
		sort.Slice(data, func(i, j int) bool {
			return utils.ToInt(data[j].DiggCount) < utils.ToInt(data[i].DiggCount)
		})
		total = len(data)
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}
		data = data[start:end]
	}
	return
}
