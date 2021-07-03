package business

import (
	"context"
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/hbaseService/hbasehelper"
)

type AuthorBusiness struct {
}

func NewAuthorBusiness() *AuthorBusiness {
	return new(AuthorBusiness)
}

func (a *AuthorBusiness) HbaseGetAuthors(rowKeys []*hbase.TGet) (data []*entity.DyAuthor) {
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
		data = append(data, author)
	}
	return
}

//达人基础数据
func (a *AuthorBusiness) HbaseGetAuthor(authorId string) (data *entity.DyAuthor, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAuthor).GetByRowKey([]byte(authorId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	authorMap := hbaseService.HbaseFormat(result, entity.DyAuthorMap)
	author := &entity.DyAuthor{}
	utils.MapToStruct(authorMap, author)
	author.AuthorID = author.Data.ID
	data = author
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
	detailMap := hbaseService.HbaseFormat(result, entity.XtAuthorDetailMap)
	detail := &entity.XtAuthorDetail{}
	utils.MapToStruct(detailMap, detail)
	detail.UID = authorId
	data = detail
	return
}
