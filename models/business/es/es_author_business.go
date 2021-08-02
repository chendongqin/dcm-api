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

func (receiver *EsAuthorBusiness) AuthorProductAnalysis(authorId, keyword string, startTime, endTime time.Time) (list []es.EsDyAuthorProductAnalysis, comErr global.CommonError) {
	esTable := GetESTableByTime(es.DyAuthorProductAnalysis, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	//esQuery.SetTerm("author_id", authorId)
	esQuery.SetRange("shelf_time", map[string]interface{}{
		"gte": 0,
		//"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if keyword != "" {
		esQuery.SetMatchPhrase("title", keyword)
	}
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit(0, 5000).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	return
}
