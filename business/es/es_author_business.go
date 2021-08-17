package es

import (
	"dongchamao/global"
	"dongchamao/global/alias"
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

func (receiver *EsAuthorBusiness) BaseSearch(
	authorId, keyword, category, secondCategory, sellTags, province, city, fanProvince, fanCity string,
	minFollower, maxFollower, minWatch, maxWatch, minDigg, maxDigg,
	minGmv, maxGmv int64, gender, minAge, maxAge, minFanAge, maxFanAge, verification, level, isDelivery, isBrand, superSeller, fanGender, page, pageSize int,
	sortStr, orderBy string) (list []es.DyAuthor, total int, comErr global.CommonError) {
	list = []es.DyAuthor{}
	if sortStr == "" {
		sortStr = "follower_incre_count"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"follower_count", "follower_incre_count", "predict_30_gmv"}) {
		comErr = global.NewError(4000)
		return
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		comErr = global.NewError(4000)
		return
	}
	if pageSize > 100 {
		comErr = global.NewError(4000)
		return
	}
	esTable := es.DyAuthorTable
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("exist", 1)
	if sortStr == "follower_count" && minFollower == 0 && maxFollower == 0 {
		minFollower = 2600000
	}
	if keyword != "" {
		if utils.HasChinese(keyword) {
			slop := 100
			length := len([]rune(keyword))
			if length <= 3 {
				slop = 2
			}
			esMultiQuery.AddMust(elasticsearch.Query().
				SetMatchPhraseWithParams("nickname", keyword, alias.M{
					"slop": slop,
				}).Condition)
		} else {
			esQuery.SetMultiMatch([]string{"unique_id", "short_id", "nickname"}, keyword)
		}
	}
	if province != "" {
		esQuery.SetTerm("province.keyword", province)
	}
	if city != "" {
		esQuery.SetTerm("city.keyword", city)
	}
	if authorId != "" {
		esQuery.SetTerm("author_id", authorId)
	}
	if category != "" {
		if category == "其他" {
			category = ""
		}
		esQuery.SetTerm("tags.keyword", category)
	}
	if sellTags != "" {
		esQuery.SetTerm("rank_sell_tags.keyword", sellTags)
	}
	if gender == 1 {
		esQuery.SetTerm("gender", 0)
	} else if gender == 2 {
		esQuery.SetTerm("gender", 1)
	}
	if level != 0 {
		esQuery.SetTerm("level", level)
	}
	if isDelivery == 1 {
		esQuery.SetTerm("is_delivery", 1)
	} else if isDelivery == 2 {
		esQuery.SetTerm("is_delivery", 0)
	}
	if isBrand == 1 {
		esQuery.SetTerm("brand", 1)
	}
	if verification != 0 {
		if verification == 1 {
			esQuery.SetTerm("verification_type", 0)
		} else if verification == 2 {
			esQuery.SetTerm("verification_type", 1)
		}
	}
	if secondCategory != "" {
		esQuery.SetTerm("tags_level_two.keyword", secondCategory)
	}
	if minGmv > 0 || maxGmv > 0 {
		rangeMap := map[string]interface{}{}
		if minGmv > 0 {
			rangeMap["gte"] = minGmv
		}
		if maxGmv > 0 {
			rangeMap["lt"] = maxGmv
		}
		esQuery.SetRange("predict_30_gmv", rangeMap)
	}
	if minAge > 0 || maxAge > 0 {
		rangeMap := map[string]interface{}{}
		if minAge > 0 {
			rangeMap["gte"] = minAge
			rangeMap["lt"] = 0
		}
		if maxAge > 0 {
			rangeMap["lt"] = maxAge
		}
		esQuery.SetRange("birthday", rangeMap)
	}
	if minDigg > 0 || maxDigg > 0 {
		rangeMap := map[string]interface{}{}
		if minDigg > 0 {
			rangeMap["gte"] = minDigg
		}
		if maxDigg > 0 {
			rangeMap["lt"] = maxDigg
		}
		esQuery.SetRange("med_digg", rangeMap)
	}
	if minFollower > 0 || maxFollower > 0 {
		rangeMap := map[string]interface{}{}
		if minFollower > 0 {
			rangeMap["gte"] = minFollower
		}
		if maxFollower > 0 {
			rangeMap["lt"] = maxFollower
		}
		esQuery.SetRange("follower_count", rangeMap)
	}
	if minWatch > 0 || maxWatch > 0 {
		rangeMap := map[string]interface{}{}
		if minWatch > 0 {
			rangeMap["gte"] = minWatch
		}
		if maxWatch > 0 {
			rangeMap["lt"] = maxWatch
		}
		esQuery.SetRange("med_watch_cnt", rangeMap)
	}
	if superSeller == 1 {
		esQuery.SetRange("follower_count", map[string]interface{}{
			"lt": 100000,
		})
		esQuery.SetRange("predict_30_gmv", map[string]interface{}{
			"gte": 100000,
		})
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
		SetOrderBy(elasticsearch.NewElasticOrder().Add("author_date_product.keyword", "asc").Order).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result, &startRow)
	_, esMultiQuery2 := elasticsearch.NewElasticQueryGroup()
	result2 := esMultiQuery2.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("author_date_product.keyword", "desc").Order).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result2, &endRow)
	return
}
