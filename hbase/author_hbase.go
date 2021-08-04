package hbase

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	entity2 "dongchamao/models/entity"
	dy2 "dongchamao/models/repost/dy"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"strings"
	"time"
)

//达人数据
func GetAuthor(authorId string) (data entity2.DyAuthorData, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAuthor).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	authorMap := hbaseService.HbaseFormat(result, entity2.DyAuthorMap)
	author := &entity2.DyAuthor{}
	utils.MapToStruct(authorMap, author)
	data = author.Data
	data.CrawlTime = author.CrawlTime
	return
}

//达人基础数据
func GetAuthorBasic(authorId, date string) (data entity2.DyAuthorBasic, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	rowKey := authorId
	if date != "" {
		rowKey += "_" + date
	}
	result, err := query.SetTable(hbaseService.HbaseDyAuthorBasic).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	basicMap := hbaseService.HbaseFormat(result, entity2.DyAuthorBasicMap)
	utils.MapToStruct(basicMap, &data)
	return
}

//获取达人粉丝数据
func GetFansByDate(authorId, date string) (data entity2.DyAuthorFans, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	rowKey := authorId + "_" + date
	result, err := query.SetTable(hbaseService.HbaseDyAuthorFans).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity2.DyAuthorFansMap)
	utils.MapToStruct(infoMap, &data)
	return
}

func GetFansRangDate(authorId, startDate, endDate string) (data map[string]entity2.DyAuthorFans, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := authorId + "_" + startDate
	endRow := authorId + "_" + endDate
	results, err := query.
		SetTable(hbaseService.HbaseDyAuthorFans).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		return
	}
	data = map[string]entity2.DyAuthorFans{}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity2.DyAuthorFansMap)
		hData := entity2.DyAuthorFans{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData
	}
	return
}

//达人基础数据趋势
func GetAuthorBasicRangeDate(authorId string, startTime, endTime time.Time) (data map[string]dy2.DyAuthorBasicChart, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startRow := authorId + "_" + startTime.Format("20060102")
	endRow := authorId + "_" + endTime.AddDate(0, 0, 1).Format("20060102")
	results, err := query.
		SetTable(hbaseService.HbaseDyAuthorBasic).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		return
	}
	data = map[string]dy2.DyAuthorBasicChart{}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity2.DyAuthorBasicMap)
		hData := dy2.DyAuthorBasicChart{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData
	}
	return
}

//获取达人粉丝团数据
func GetAuthorFansClub(authorId string) (data entity2.DyLiveFansClub, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveFansClub).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity2.DyLiveFansClubMap)
	utils.MapToStruct(dataMap, &data)
	return
}

//达人（带货）口碑
func GetAuthorReputation(authorId string) (data *entity2.DyReputation, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyReputation).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	reputationMap := hbaseService.HbaseFormat(result, entity2.DyReputationMap)
	utils.MapToStruct(reputationMap, &data)
	return
}

//星图达人
func GetXtAuthorDetail(authorId string) (data *entity2.XtAuthorDetail, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseXtAuthorDetail).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity2.XtAuthorDetailMap)
	detail := &entity2.XtAuthorDetail{}
	utils.MapToStruct(detailMap, detail)
	detail.UID = authorId
	data = detail
	return
}

//获取达人直播间
func GetAuthorRoomsRangDate(authorId string, startTime, endTime time.Time) (data map[string][]entity2.DyAuthorLiveRoom, comErr global.CommonError) {
	data = map[string][]entity2.DyAuthorLiveRoom{}
	startDate := startTime.Format("20060102")
	endDate := endTime.AddDate(0, 0, 1).Format("20060102")
	query := hbasehelper.NewQuery()
	startRow := authorId + "_" + startDate
	endRow := authorId + "_" + endDate
	results, err := query.
		SetTable(hbaseService.HbaseDyAuthorRoomMapping).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		return
	}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity2.DyAuthorRoomMappingMap)
		hData := entity2.DyAuthorRoomMapping{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData.Data
	}
	return
}

//获取达人当日直播间
func GetAuthorRoomsByDate(authorId, date string) (data []entity2.DyAuthorLiveRoom, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	rowKey := authorId + "_" + date
	result, err := query.SetTable(hbaseService.HbaseDyAuthorRoomMapping).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity2.DyAuthorFansMap)
	hData := &entity2.DyAuthorRoomMapping{}
	utils.MapToStruct(infoMap, hData)
	data = hData.Data
	return
}

func GetAuthorProductAnalysis(rowKey string) (data entity2.DyAuthorProductAnalysis, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAuthorProductAnalysis).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity2.DyAuthorProductAnalysisMap)
	utils.MapToStruct(infoMap, &data)
	return
}

func GetAuthorProductAnalysisRange(startRowKey, stopRowKey string) (data []entity2.DyAuthorProductAnalysis, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	results, err := query.
		SetTable(hbaseService.HbaseDyAuthorProductAnalysis).
		SetStartRow([]byte(startRowKey)).
		SetStopRow([]byte(stopRowKey)).
		Scan(10000)
	if err != nil {
		return
	}
	for _, v := range results {
		infoMap := hbaseService.HbaseFormat(v, entity2.DyAuthorProductAnalysisMap)
		detail := entity2.DyAuthorProductAnalysis{}
		utils.MapToStruct(infoMap, &detail)
		data = append(data, detail)
	}
	return
}
