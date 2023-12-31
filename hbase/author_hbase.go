package hbase

import (
	"context"
	"dongchamao/global"
	"dongchamao/global/logger"
	"dongchamao/global/utils"
	"dongchamao/models/entity"
	"dongchamao/models/repost/dy"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/hbaseService/hbasehelper"
	"strings"
	"time"
)

func GetAuthorByIds(authorIds []string) (map[string]entity.DyAuthor, error) {
	rowKeys := make([]*hbase.TGet, 0)
	for _, authorId := range authorIds {
		rowKeys = append(rowKeys, &hbase.TGet{Row: []byte(authorId)})
	}
	client := global.HbasePools.Get("default")
	defer client.Close()
	results, err := client.GetMultiple(context.Background(), []byte(hbaseService.HbaseDyAuthor), rowKeys)
	if err != nil {
		return nil, err
	}
	authorMap := map[string]entity.DyAuthor{}
	for _, v := range results {
		data := entity.DyAuthor{}
		detailMap := hbaseService.HbaseFormat(v, entity.DyAuthorMap)
		utils.MapToStruct(detailMap, &data)
		data.AuthorID = data.Data.ID
		if data.RoomStatus != 2 {
			data.RoomId = ""
		}
		data.Data.RoomID = data.RoomId
		if data.Tags == "0" {
			data.Tags = ""
		}
		if data.TagsLevelTwo == "0" {
			data.TagsLevelTwo = ""
		}
		if data.Data.UniqueID == "0" || data.Data.UniqueID == "" {
			data.Data.UniqueID = data.Data.ShortID
		}
		data.Data.CrawlTime = data.CrawlTime
		if data.ArtificialData.Tags != "" {
			data.Tags = data.ArtificialData.Tags
			data.TagsLevelTwo = data.ArtificialData.TagsLevelTow
		}
		authorMap[data.AuthorID] = data
	}
	return authorMap, nil
}

//达人数据
func GetAuthor(authorId string) (data entity.DyAuthor, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAuthor).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	authorMap := hbaseService.HbaseFormat(result, entity.DyAuthorMap)
	utils.MapToStruct(authorMap, &data)
	data.AuthorID = data.Data.ID
	if data.RoomStatus != 2 {
		data.RoomId = ""
	}
	data.Data.RoomID = data.RoomId
	if data.Tags == "0" {
		data.Tags = ""
	}
	if data.TagsLevelTwo == "0" {
		data.TagsLevelTwo = ""
	}
	if data.Data.UniqueID == "0" || data.Data.UniqueID == "" {
		data.Data.UniqueID = data.Data.ShortID
	}
	data.Data.CrawlTime = data.CrawlTime
	if data.ArtificialData.Tags != "" {
		data.Tags = data.ArtificialData.Tags
		data.TagsLevelTwo = data.ArtificialData.TagsLevelTow
	}
	return
}

//达人基础数据
func GetAuthorBasic(authorId, date string) (data entity.DyAuthorBasic, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	rowKey := authorId
	if date != "" {
		rowKey += "_" + date
	}
	result, err := query.SetTable(hbaseService.HbaseDyAuthorBasic).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	basicMap := hbaseService.HbaseFormat(result, entity.DyAuthorBasicMap)
	utils.MapToStruct(basicMap, &data)
	return
}

//获取达人粉丝数据
func GetFansByDate(authorId, date string) (data entity.DyAuthorFans, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	rowKey := authorId + "_" + date
	result, err := query.SetTable(hbaseService.HbaseDyAuthorFans).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyAuthorFansMap)
	utils.MapToStruct(infoMap, &data)
	return
}

func GetFansRangDate(authorId, startDate, endDate string) (data map[string]entity.DyAuthorFans, comErr global.CommonError) {
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
	data = map[string]entity.DyAuthorFans{}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity.DyAuthorFansMap)
		hData := entity.DyAuthorFans{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData
	}
	return
}

//达人基础数据趋势
func GetAuthorBasicRangeDate(authorId string, startTime, endTime time.Time) (data map[string]dy.DyAuthorBasicChart, comErr global.CommonError) {
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
	data = map[string]dy.DyAuthorBasicChart{}
	for _, v := range results {
		rowKey := string(v.GetRow())
		rowKeyArr := strings.Split(rowKey, "_")
		if len(rowKeyArr) < 2 {
			comErr = global.NewError(5000)
			return
		}
		date := rowKeyArr[1]
		dataMap := hbaseService.HbaseFormat(v, entity.DyAuthorBasicMap)
		hData := dy.DyAuthorBasicChart{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData
	}
	return
}

//获取达人粉丝团数据
func GetAuthorFansClub(authorId string) (data entity.DyLiveFansClub, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveFansClub).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	dataMap := hbaseService.HbaseFormat(result, entity.DyLiveFansClubMap)
	utils.MapToStruct(dataMap, &data)
	return
}

//达人小店|品牌
func GetAuthorStore(authorId string) (data entity.DyAuthorStore, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAuthorStore).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	reputationMap := hbaseService.HbaseFormat(result, entity.DyAuthorStoreMap)
	utils.MapToStruct(reputationMap, &data)
	return
}

//达人（带货）口碑
func GetAuthorReputation(authorId string) (data entity.DyReputation, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyReputation).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	reputationMap := hbaseService.HbaseFormat(result, entity.DyReputationMap)
	utils.MapToStruct(reputationMap, &data)
	return
}

//星图达人
func GetXtAuthorDetail(authorId string) (data entity.XtAuthorDetail, comErr global.CommonError) {
	data = entity.XtAuthorDetail{}
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseXtAuthorDetail).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.XtAuthorDetailMap)
	utils.MapToStruct(detailMap, &data)
	data.UID = authorId
	return
}

//获取达人直播间
func GetAuthorRoomsRangDate(authorId string, startTime, endTime time.Time) (data map[string][]entity.DyAuthorLiveRoom, comErr global.CommonError) {
	data = map[string][]entity.DyAuthorLiveRoom{}
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
		dataMap := hbaseService.HbaseFormat(v, entity.DyAuthorRoomMappingMap)
		hData := entity.DyAuthorRoomMapping{}
		utils.MapToStruct(dataMap, &hData)
		data[date] = hData.Data
	}
	return
}

//获取达人当日直播间
func GetAuthorRoomsByDate(authorId, date string) (data []entity.DyAuthorLiveRoom, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	rowKey := authorId + "_" + date
	result, err := query.SetTable(hbaseService.HbaseDyAuthorRoomMapping).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyAuthorRoomMappingMap)
	hData := &entity.DyAuthorRoomMapping{}
	utils.MapToStruct(infoMap, hData)
	data = hData.Data
	return
}

func GetAuthorProductAnalysis(rowKey string) (data entity.DyAuthorProductAnalysis, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAuthorProductAnalysis).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyAuthorProductAnalysisMap)
	utils.MapToStruct(infoMap, &data)
	return
}

func GetAuthorProductAnalysisRange(startRowKey, stopRowKey string) (data []entity.DyAuthorProductAnalysis, comErr global.CommonError) {
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
		infoMap := hbaseService.HbaseFormat(v, entity.DyAuthorProductAnalysisMap)
		detail := entity.DyAuthorProductAnalysis{}
		utils.MapToStruct(infoMap, &detail)
		data = append(data, detail)
	}
	return
}

//达人带货行业
func GetAuthorLiveTags() (data []entity.DyAuthorLiveTags, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	results, err := query.
		SetTable(hbaseService.HbaseDyAuthorLiveTags).
		Scan(10000)
	if err != nil {
		return
	}
	for _, v := range results {
		infoMap := hbaseService.HbaseFormat(v, entity.DyAuthorLiveTagsMap)
		detail := entity.DyAuthorLiveTags{}
		utils.MapToStruct(infoMap, &detail)
		if detail.Tags == "" || detail.Tags == "null" {
			continue
		}
		data = append(data, detail)
	}
	return
}

//达人带货商品
func GetAuthorDateProduct(rowKey string) (data entity.DyAuthorDateProductData, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAuthorProductDateMapping).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyAuthorProductDateMappingMap)
	utils.MapToStruct(infoMap, &data)
	return
}

//达人带货商品
func GetAuthorProductRangeDate(authorId string, startTime, stopTime time.Time) (data []entity.DyAuthorDateProductData, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	startKey := authorId + "_" + startTime.Format("20060102")
	stopKey := authorId + "_" + stopTime.Format("20060102")
	results, err := query.
		SetStartRow([]byte(startKey)).
		SetStopRow([]byte(stopKey)).
		SetTable(hbaseService.HbaseDyAuthorProductDateMapping).
		Scan(1000)
	if err != nil {
		return
	}
	for _, v := range results {
		infoMap := hbaseService.HbaseFormat(v, entity.DyAuthorProductDateMappingMap)
		detail := entity.DyAuthorDateProductData{}
		utils.MapToStruct(infoMap, &detail)
		data = append(data, detail)
	}
	return
}

func GetAuthorLiveFansClubUser(authorId string) (data entity.DyLiveFansClubUser, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAuthorLiveFansClubUser).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewError(5000)
		logger.Error(err)
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLiveFansClubUserMap)
	utils.MapToStruct(detailMap, &data)
	return
}

// GetAuthorRank 达人排行数据
func GetAuthorRank(authorId string) (data entity.DyAuthorRank, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.DyAuthorPersonalTag).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	authorMap := hbaseService.HbaseFormat(result, entity.DyAuthorRankMap)
	utils.MapToStruct(authorMap, &data)

	return
}
