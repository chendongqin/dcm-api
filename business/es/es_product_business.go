package es

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/services/elasticsearch"
	"time"
)

type EsProductBusiness struct {
}

func NewEsProductBusiness() *EsProductBusiness {
	return new(EsProductBusiness)
}

func (i *EsProductBusiness) SearchRangeDateList(productId, keyword string, startTime, endTime time.Time, page, pageSize int) (list []es.DyProductAuthorAnalysis, total int, comErr global.CommonError) {
	esTable := GetESTableByTime(es.DyProductAuthorAnalysisTable, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("productId", productId)
	esQuery.SetRange("createSdf.keyword", map[string]interface{}{
		"gte": startTime.Format("20060102"),
		"lte": endTime.Format("20060102"),
	})
	if keyword != "" {
		esQuery.SetMultiMatch([]string{"nickname", "displayId", "shortId"}, keyword)
	}
	results := esMultiQuery.
		SetTable(esTable).
		SetFields("productId", "authorId", "createSdf").
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("_id", "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

func (i *EsProductBusiness) SearchRangeDateRowKey(productId, keyword string, startTime, endTime time.Time) (startRow es.DyProductAuthorAnalysis, stopRow es.DyProductAuthorAnalysis, total int, comErr global.CommonError) {
	esTable := GetESTableByTime(es.DyProductAuthorAnalysisTable, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("product_id", productId)
	esQuery.SetRange("createSdf.keyword", map[string]interface{}{
		"gte": startTime.Format("20060102"),
		"lte": endTime.Format("20060102"),
	})
	if keyword != "" {
		esQuery.SetMultiMatch([]string{"nickname", "displayId", "shortId"}, keyword)
	}
	result := esMultiQuery.
		SetTable(esTable).
		SetFields("productId", "authorId", "createSdf").
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("_id", "desc").Order).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result, &startRow)
	result2 := esMultiQuery.
		SetTable(esTable).
		SetFields("productId", "authorId", "createSdf").
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("_id", "asc").Order).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result2, &stopRow)
	total = esMultiQuery.Count
	return
}
