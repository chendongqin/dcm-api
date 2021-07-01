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
