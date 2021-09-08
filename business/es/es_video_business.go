package es

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/services/elasticsearch"
	"time"
)

type EsVideoBusiness struct {
}

func NewEsVideoBusiness() *EsVideoBusiness {
	return new(EsVideoBusiness)
}

func (e *EsVideoBusiness) SearchAwemeByProduct(productId, keyword, sortStr, orderBy string,
	startTime, endTime time.Time, page, pageSize int) (list []es.DyProductVideo, total int, comErr global.CommonError) {
	if orderBy == "" {
		orderBy = "desc"
	}
	if sortStr == "" {
		sortStr = "aweme_create_time"
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		comErr = global.NewError(4000)
		return
	}
	if !utils.InArrayString(sortStr, []string{"aweme_create_time", "aweme_gmv", "sales", "comment_count", "digg_count", "forward_count"}) {
		comErr = global.NewError(4000)
		return
	}
	esTable := GetESTableByMonthTime(es.DyProductVideoTable, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("product_id", productId)
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if keyword != "" {
		esQuery.AddCondition(map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match_phrase": map[string]interface{}{
							"aweme_title": keyword,
						},
					},
					{
						"match_phrase": map[string]interface{}{
							"nickname": keyword,
						},
					},
				},
			},
		})
	}
	result := esMultiQuery.
		SetTable(esTable).
		SetCache(180).
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetLimit((page-1)*pageSize, pageSize).
		SetMultiQuery().
		Query()
	utils.MapToStruct(result, &list)
	total = esMultiQuery.Count
	return
}

//获取达人列表
func (e *EsVideoBusiness) SearchByAuthor(authorId, keyword, sortStr, orderBy string, hasProduct, page, pageSize int, startTime, endTime time.Time) (list []es.DyAweme, total int, comErr global.CommonError) {
	if orderBy == "" {
		orderBy = "desc"
	}
	if sortStr == "" {
		sortStr = "aweme_create_time"
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		comErr = global.NewError(4000)
		return
	}
	if !utils.InArrayString(sortStr, []string{"aweme_create_time", "digg_count", "comment_count", "share_count", "aweme_gmv", "sales"}) {
		comErr = global.NewError(4000)
		return
	}
	esTable := GetESTableByMonthTime(es.DyVideoTable, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("author_id", authorId)
	esQuery.SetExist("field", "aweme_title")
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if keyword != "" {
		esQuery.SetMatchPhrase("aweme_title.keyword", keyword)
	}
	if hasProduct == 1 {
		esQuery.SetExist("field", "product_ids")
	}
	result := esMultiQuery.
		SetTable(esTable).
		SetCache(180).
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetLimit((page-1)*pageSize, pageSize).
		SetMultiQuery().
		Query()
	utils.MapToStruct(result, &list)
	total = esMultiQuery.Count
	return
}
