package business

import (
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"dongchamao/structinit/repost/dy"
	"strings"
	"time"
)

const ShareUrlPrefix = "https://www.iesdouyin.com/share/user/"

type AuthorBusiness struct {
}

func NewAuthorBusiness() *AuthorBusiness {
	return new(AuthorBusiness)
}

//粉丝｜粉丝团趋势数据
func (a *AuthorBusiness) HbaseGetFansRangDate(authorId, startDate, endDate string) (data map[string]dy.DateChart, comErr global.CommonError) {
	data = map[string]dy.DateChart{}
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
	dateMap := map[string]entity.DyAuthorFans{}
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
		dateMap[date] = hData
	}
	//起点补点操作
	t1, _ := time.ParseInLocation("20060102", startDate, time.Local)
	t2, _ := time.ParseInLocation("20060102", endDate, time.Local)
	beforeDate := t1.AddDate(0, 0, -1).Format("20060102")
	beforeData, _ := a.HbaseGetFansByDate(authorId, beforeDate)
	if _, ok := dateMap[startDate]; !ok {
		dateMap[startDate] = beforeData
	}
	//末点补点
	if endDate == time.Now().Format("20060102") {
		if _, ok := dateMap[endDate]; !ok {
			endData, _ := a.HbaseGetAuthorBasic(authorId, "")
			dateMap[endDate] = entity.DyAuthorFans{
				FollowerCount:       endData.FollowerCount,
				TotalFansGroupCount: endData.TotalFansCount,
			}
		}
	}
	countGroupArr := make([]int64, 0)
	incGroupArr := make([]int64, 0)
	dateArr := make([]string, 0)
	countArr := make([]int64, 0)
	incArr := make([]int64, 0)
	beginDatetime := t1
	for {
		if beginDatetime.After(t2) {
			break
		}
		date := beginDatetime.Format("20060102")
		if _, ok := dateMap[date]; !ok {
			yesterday := beginDatetime.AddDate(0, 0, -1).Format("20060102")
			if _, ok := dateMap[yesterday]; ok {
				dateMap[date] = dateMap[yesterday]
			} else {
				dateMap[date] = entity.DyAuthorFans{
					FollowerCount:       0,
					TotalFansGroupCount: 0,
				}
			}
		}
		fansData := dateMap[date]
		inc1 := fansData.FollowerCount - beforeData.FollowerCount
		inc2 := fansData.TotalFansGroupCount - beforeData.TotalFansGroupCount
		dateArr = append(dateArr, beginDatetime.Format("01/02"))
		incArr = append(incArr, inc1)
		incGroupArr = append(incGroupArr, inc2)
		countArr = append(countArr, fansData.FollowerCount)
		countGroupArr = append(countGroupArr, fansData.TotalFansGroupCount)
		beforeData = fansData
		beginDatetime = beginDatetime.AddDate(0, 0, 1)
	}
	data["fans"] = dy.DateChart{
		Date:       dateArr,
		CountValue: countArr,
		IncValue:   incArr,
	}
	data["fans_group"] = dy.DateChart{
		Date:       dateArr,
		CountValue: countGroupArr,
		IncValue:   incGroupArr,
	}
	return
}

//获取达人粉丝数据
func (a *AuthorBusiness) HbaseGetFansByDate(authorId, date string) (data entity.DyAuthorFans, comErr global.CommonError) {
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
	infoMap := hbaseService.HbaseFormat(result, entity.DyAuthorFansMap)
	utils.MapToStruct(infoMap, &data)
	return
}

//达人数据
func (a *AuthorBusiness) HbaseGetAuthor(authorId string) (data entity.DyAuthorData, comErr global.CommonError) {
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
	authorMap := hbaseService.HbaseFormat(result, entity.DyAuthorMap)
	author := &entity.DyAuthor{}
	utils.MapToStruct(authorMap, author)
	author.AuthorID = author.Data.ID
	author.Data.Age = GetAge(author.Data.Birthday)
	author.Data.Avatar = dyimg.Fix(author.Data.Avatar)
	author.Data.ShareUrl = ShareUrlPrefix + author.AuthorID
	if author.Data.UniqueID == "" {
		author.Data.UniqueID = author.Data.ShortID
	}
	data = author.Data
	return
}

//达人基础数据
func (a *AuthorBusiness) HbaseGetAuthorBasic(authorId, date string) (data entity.DyAuthorBasic, comErr global.CommonError) {
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
	basicMap := hbaseService.HbaseFormat(result, entity.DyAuthorBasicMap)
	utils.MapToStruct(basicMap, &data)
	return
}

//达人基础数据趋势
func (a *AuthorBusiness) HbaseGetAuthorBasicRangeDate(authorId, startDate, endDate string) (data map[string]dy.DateChart, comErr global.CommonError) {
	data = map[string]dy.DateChart{}
	query := hbasehelper.NewQuery()
	startRow := authorId + "_" + startDate
	endRow := authorId + "_" + endDate
	results, err := query.
		SetTable(hbaseService.HbaseDyAuthorBasic).
		SetStartRow([]byte(startRow)).
		SetStopRow([]byte(endRow)).
		Scan(1000)
	if err != nil {
		return
	}
	dateMap := map[string]dy.DyAuthorBasicChart{}
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
		dateMap[date] = hData
	}
	//起点补点操作
	t1, _ := time.ParseInLocation("20060102", startDate, time.Local)
	t2, _ := time.ParseInLocation("20060102", endDate, time.Local)
	beforeDate := t1.AddDate(0, 0, -1).Format("20060102")
	beforeBasicData, _ := a.HbaseGetAuthorBasic(authorId, beforeDate)
	beforeData := dy.DyAuthorBasicChart{
		FollowerCount:  beforeBasicData.FollowerCount,
		TotalFansCount: beforeBasicData.TotalFansCount,
		TotalFavorited: beforeBasicData.TotalFavorited,
		CommentCount:   beforeBasicData.CommentCount,
		ForwardCount:   beforeBasicData.ForwardCount,
	}
	if _, ok := dateMap[startDate]; !ok {
		dateMap[startDate] = beforeData
	}
	//末点补点
	if endDate == time.Now().Format("20060102") {
		if _, ok := dateMap[endDate]; !ok {
			basicData, _ := a.HbaseGetAuthorBasic(authorId, "")
			dateMap[endDate] = dy.DyAuthorBasicChart{
				FollowerCount:  basicData.FollowerCount,
				TotalFansCount: basicData.TotalFansCount,
				TotalFavorited: basicData.TotalFavorited,
				CommentCount:   basicData.CommentCount,
				ForwardCount:   basicData.ForwardCount,
			}
		}
	}
	dateArr := make([]string, 0)
	fansCountArr := make([]int64, 0)
	fanIncArr := make([]int64, 0)
	fansGroupCountArr := make([]int64, 0)
	fansGroupIncArr := make([]int64, 0)
	diggCountArr := make([]int64, 0)
	diggIncArr := make([]int64, 0)
	commentCountArr := make([]int64, 0)
	commentIncArr := make([]int64, 0)
	forwardCountArr := make([]int64, 0)
	forwardIncArr := make([]int64, 0)
	beginDatetime := t1
	for {
		if beginDatetime.After(t2) {
			break
		}
		date := beginDatetime.Format("20060102")
		if _, ok := dateMap[date]; !ok {
			dateMap[date] = beforeData
		}
		currentData := dateMap[date]
		dateArr = append(dateArr, beginDatetime.Format("01/02"))
		fansCountArr = append(fansCountArr, currentData.FollowerCount)
		fanIncArr = append(fanIncArr, currentData.FollowerCount-beforeData.FollowerCount)
		fansGroupCountArr = append(fansGroupCountArr, currentData.TotalFansCount)
		fansGroupIncArr = append(fansGroupIncArr, currentData.TotalFansCount-beforeData.TotalFansCount)
		diggCountArr = append(diggCountArr, currentData.TotalFavorited)
		diggIncArr = append(diggIncArr, currentData.TotalFavorited-beforeData.TotalFavorited)
		commentCountArr = append(commentCountArr, currentData.CommentCount)
		commentIncArr = append(commentIncArr, currentData.CommentCount-beforeData.CommentCount)
		forwardCountArr = append(forwardCountArr, currentData.ForwardCount)
		forwardIncArr = append(forwardIncArr, currentData.ForwardCount-beforeData.ForwardCount)
		beforeData = currentData
		beginDatetime = beginDatetime.AddDate(0, 0, 1)
	}
	data["fans"] = dy.DateChart{
		Date:       dateArr,
		CountValue: fansCountArr,
		IncValue:   fanIncArr,
	}
	data["fans_club"] = dy.DateChart{
		Date:       dateArr,
		CountValue: fansGroupCountArr,
		IncValue:   fansGroupIncArr,
	}
	data["digg"] = dy.DateChart{
		Date:       dateArr,
		CountValue: diggCountArr,
		IncValue:   diggIncArr,
	}
	data["forward"] = dy.DateChart{
		Date:       dateArr,
		CountValue: forwardCountArr,
		IncValue:   forwardIncArr,
	}
	data["comment"] = dy.DateChart{
		Date:       dateArr,
		CountValue: commentCountArr,
		IncValue:   commentIncArr,
	}
	return
}

//获取达人粉丝团数据
func (a *AuthorBusiness) HbaseGetAuthorFansClub(authorId string) (data entity.DyLiveFansClub, comErr global.CommonError) {
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
	dataMap := hbaseService.HbaseFormat(result, entity.DyLiveFansClubMap)
	utils.MapToStruct(dataMap, &data)
	return
}

//达人（带货）口碑
func (a *AuthorBusiness) HbaseGetAuthorReputation(authorId string) (data *entity.DyReputation, comErr global.CommonError) {
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
	reputationMap := hbaseService.HbaseFormat(result, entity.DyReputationMap)
	reputation := &entity.DyReputation{}
	utils.MapToStruct(reputationMap, reputation)
	if reputation.ScoreList == nil {
		reputation.ScoreList = make([]entity.DyReputationDateScoreList, 0)
	}
	//reputation.ShopLogo = dyimg.Fix(reputation.ShopLogo)
	data = reputation
	return
}

//星图达人
func (a *AuthorBusiness) HbaseGetXtAuthorDetail(authorId string) (data *entity.XtAuthorDetail, comErr global.CommonError) {
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
	detailMap := hbaseService.HbaseFormat(result, entity.XtAuthorDetailMap)
	detail := &entity.XtAuthorDetail{}
	utils.MapToStruct(detailMap, detail)
	detail.UID = authorId
	data = detail
	return
}

//获取达人星图评分
func (a *AuthorBusiness) GetDyAuthorScore(liveScore entity.XtAuthorLiveScore, videoScore entity.XtAuthorScore) (authorStarScores dy.DyAuthorStarScores) {
	authorStarScores.LiveScore = dy.DyAuthorStarScore{
		CooperateIndex: utils.FriendlyFloat64(float64(liveScore.CooperateIndex) / 10000),
		CpIndex:        utils.FriendlyFloat64(float64(liveScore.CpIndex) / 10000),
		GrowthIndex:    utils.FriendlyFloat64(float64(liveScore.GrowthIndex) / 10000),
		ShoppingIndex:  utils.FriendlyFloat64(float64(liveScore.ShoppingIndex) / 10000),
		SpreadIndex:    utils.FriendlyFloat64(float64(liveScore.SpreadIndex) / 10000),
		TopScore:       utils.FriendlyFloat64(float64(liveScore.TopScore) / 10000),
	}
	authorStarScores.LiveScore.AvgLevel = dy.XtAuthorScoreAvgLevel{
		CooperateIndex: utils.FriendlyFloat64(float64(liveScore.AvgLevel.CooperateIndex) / 100),
		CpIndex:        utils.FriendlyFloat64(float64(liveScore.AvgLevel.CpIndex) / 100),
		GrowthIndex:    utils.FriendlyFloat64(float64(liveScore.AvgLevel.GrowthIndex) / 100),
		ShoppingIndex:  utils.FriendlyFloat64(float64(liveScore.AvgLevel.ShoppingIndex) / 100),
		SpreadIndex:    utils.FriendlyFloat64(float64(liveScore.AvgLevel.SpreadIndex) / 100),
		TopScore:       utils.FriendlyFloat64(float64(liveScore.AvgLevel.TopScore) / 100),
	}
	authorStarScores.VideoScore = dy.DyAuthorStarScore{
		CooperateIndex: utils.FriendlyFloat64(float64(videoScore.CooperateIndex) / 10000),
		CpIndex:        utils.FriendlyFloat64(float64(videoScore.CpIndex) / 10000),
		GrowthIndex:    utils.FriendlyFloat64(float64(videoScore.GrowthIndex) / 10000),
		ShoppingIndex:  utils.FriendlyFloat64(float64(videoScore.ShoppingIndex) / 10000),
		SpreadIndex:    utils.FriendlyFloat64(float64(videoScore.SpreadIndex) / 10000),
		TopScore:       utils.FriendlyFloat64(float64(videoScore.TopScore) / 10000),
	}
	authorStarScores.VideoScore.AvgLevel = dy.XtAuthorScoreAvgLevel{
		CooperateIndex: utils.FriendlyFloat64(float64(videoScore.AvgLevel.CooperateIndex) / 100),
		CpIndex:        utils.FriendlyFloat64(float64(videoScore.AvgLevel.CpIndex) / 100),
		GrowthIndex:    utils.FriendlyFloat64(float64(videoScore.AvgLevel.GrowthIndex) / 100),
		ShoppingIndex:  utils.FriendlyFloat64(float64(videoScore.AvgLevel.ShoppingIndex) / 100),
		SpreadIndex:    utils.FriendlyFloat64(float64(videoScore.AvgLevel.SpreadIndex) / 100),
		TopScore:       utils.FriendlyFloat64(float64(videoScore.AvgLevel.TopScore) / 100),
	}
	return
}

//获取达人直播间
func (a *AuthorBusiness) HbaseGetAuthorRoomsRangDate(authorId, startDate, endDate string) (data map[string][]entity.DyAuthorLiveRoom, comErr global.CommonError) {
	data = map[string][]entity.DyAuthorLiveRoom{}
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
func (a *AuthorBusiness) HbaseGetAuthorRoomsByDate(authorId, date string) (data []entity.DyAuthorLiveRoom, comErr global.CommonError) {
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
	infoMap := hbaseService.HbaseFormat(result, entity.DyAuthorFansMap)
	hData := &entity.DyAuthorRoomMapping{}
	utils.MapToStruct(infoMap, hData)
	data = hData.Data
	return
}
