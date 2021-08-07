package es

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/services/elasticsearch"
	"time"
)

type EsAuthorBusiness struct {
}

func NewEsAuthorBusiness() *EsAuthorBusiness {
	return new(EsAuthorBusiness)
}

func (receiver *EsAuthorBusiness) AuthorProductAnalysis(authorId, keyword string, startTime, endTime time.Time) (startRow es.EsDyAuthorProductAnalysis, endRow es.EsDyAuthorProductAnalysis, comErr global.CommonError) {
	esTable := GetESTableByTime(es.DyAuthorProductAnalysisTable, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("author_id", authorId)
	esQuery.SetRange("shelf_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if keyword != "" {
		esQuery.SetMatchPhrase("title", keyword)
	}
	result := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("author_date_product.keyword", "desc").Order).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result, &startRow)
	result2 := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("author_date_product.keyword", "asc").Order).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result2, &endRow)
	return
}
