package business

import (
	"context"
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/hbaseService/hbasehelper"
	"dongchamao/structinit/repost/dy"
)

const ShareUrlPrefix = "https://www.iesdouyin.com/share/user/"

type AuthorBusiness struct {
}

func NewAuthorBusiness() *AuthorBusiness {
	return new(AuthorBusiness)
}

func (a *AuthorBusiness) HbaseGetAuthors(rowKeys []*hbase.TGet) (data []entity.DyAuthorData) {
	client := global.HbasePools.Get("default")
	tableName := hbaseService.HbaseDyAuthor
	tableBytes := []byte(tableName)
	results, err := client.GetMultiple(context.Background(), tableBytes, rowKeys)
	if err != nil {
		return
	}
	for _, v := range results {
		authorMap := hbaseService.HbaseFormat(v, entity.DyAuthorMap)
		author := &entity.DyAuthor{}
		utils.MapToStruct(authorMap, author)
		author.AuthorID = author.Data.ID
		author.Data.Age = GetAge(author.Data.Birthday)
		author.Data.Avatar = dyimg.Fix(author.Data.Avatar)
		author.Data.ShareUrl = ShareUrlPrefix + author.AuthorID
		if author.Data.UniqueID == "" {
			author.Data.UniqueID = author.Data.ShortID
		}
		data = append(data, author.Data)
	}
	return
}

//达人基础数据
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
