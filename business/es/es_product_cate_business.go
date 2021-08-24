package es

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/services/elasticsearch"
)

type EsProductCateBusiness struct {
}

func NewEsProductCateBusiness() *EsProductCateBusiness {
	return new(EsProductCateBusiness)
}

//商品品类库查询
func (receiver *EsProductCateBusiness) SimpleSearch(
	nickname, keyword, tags, secondTags string,
	page, pageSize int) (list []es.DyAuthor, total int, comErr global.CommonError) {
	list = []es.DyAuthor{}
	sortStr := "follower_count"
	orderBy := "desc"
	if pageSize > 100 {
		comErr = global.NewError(4000)
		return
	}
	esTable := es.DyAuthorTable
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("exist", 1)
	if tags != "" {
		if tags == "其他" {
			tags = ""
		}
		esQuery.SetTerm("tags.keyword", tags)
	}
	if secondTags != "" {
		esQuery.SetTerm("tags_level_two.keyword", secondTags)
	}
	if nickname != "" {
		esQuery.SetMatchPhrase("nickname", nickname)
	}
	if keyword != "" {
		esQuery.SetMultiMatch([]string{"unique_id", "short_id", "author_id"}, keyword)
	}
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}
