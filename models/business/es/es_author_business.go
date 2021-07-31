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

func (receiver *EsAuthorBusiness) AuthorProductAnalysis(authorId, keyword, firstCate, secondCate, thirdCate, brandName, shopId string, shopType int, startTime, endTime time.Time) (list []es.EsDyAuthorProductAnalysis, comErr global.CommonError) {
	esTable := GetESTableByTime(es.DyAuthorProductAnalysis, startTime, endTime)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("author_id", authorId)
	esQuery.SetRange("shelf_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if shopType == 1 && shopId == "" {
		return
	} else if shopType == 2 {
		esQuery.SetTerm("shop_id.keyword", shopId)
	}
	if keyword != "" {
		esQuery.SetMatchPhrase("title", keyword)
	}
	if firstCate != "" {
		if firstCate == "其他" {
			esQuery.AddCondition(map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{
							"trems": map[string]interface{}{"dcm_level_first.keyword": []string{firstCate, ""}},
						},
						{
							"bool": map[string]interface{}{
								"must_not": map[string]interface{}{
									"exists": map[string]interface{}{
										"field": "dcm_level_first",
									},
								},
							},
						},
					},
				},
			})
			secondCate = ""
			thirdCate = ""
		} else {
			esQuery.SetTerm("dcm_level_first.keyword", firstCate)
		}
	}
	if secondCate != "" {
		esQuery.SetTerm("first_cname.keyword", secondCate)
	}
	if thirdCate != "" {
		esQuery.SetTerm("second_cname.keyword", thirdCate)
	}
	if brandName != "" {
		if brandName == "其他" {
			esQuery.SetTerms("brand_name.keyword", []string{brandName, ""})
		}
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
