package es

import (
	"dongchamao/global"
	"dongchamao/global/alias"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/services/elasticsearch"
	"fmt"
	"time"
)

type EsAuthorBusiness struct {
}

func NewEsAuthorBusiness() *EsAuthorBusiness {
	return new(EsAuthorBusiness)
}

//达人库查询
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
	if sortStr == "follower_count" && minFollower == 0 && maxFollower == 0 && keyword == "" {
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
			esQuery.AddCondition(map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{
							"match_phrase": map[string]interface{}{
								"tags": map[string]interface{}{
									"query": "",
								},
							},
						},
						{
							"match_phrase": map[string]interface{}{
								"tags": map[string]interface{}{
									"query": "0",
								},
							},
						},
						{
							"term": map[string]interface{}{
								"tags": category,
							},
						},
					},
				},
			})
		} else {
			esQuery.SetTerm("tags.keyword", category)
		}
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
		SetCache(180).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

//达人库查询
func (receiver *EsAuthorBusiness) SimpleSearch(
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
			esQuery.AddCondition(map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{
							"match_phrase": map[string]interface{}{
								"tags": map[string]interface{}{
									"query": "",
								},
							},
						},
						{
							"match_phrase": map[string]interface{}{
								"tags": map[string]interface{}{
									"query": "0",
								},
							},
						},
						{
							"term": map[string]interface{}{
								"tags": tags,
							},
						},
					},
				},
			})
		} else {
			esQuery.SetTerm("tags.keyword", tags)
		}
	}
	if secondTags != "" {
		esQuery.SetTerm("tags_level_two.keyword", secondTags)
	}
	if nickname != "" {
		esQuery.SetMatchPhrase("nickname", nickname)
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
			esQuery.SetMultiMatch([]string{"unique_id", "short_id", "nickname", "author_id"}, keyword)
		}
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

func (receiver *EsAuthorBusiness) KeywordSearch(keyword string) (list []es.DyAuthor) {
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esTable := es.DyAuthorTable
	esQuery.SetTerm("exist", 1)
	if utils.HasChinese(keyword) {
		slop := 100
		length := len([]rune(keyword))
		if length <= 3 {
			slop = 2
		}
		esQuery.SetMatchPhraseWithParams("nickname", keyword, alias.M{
			"slop": slop,
		})
	} else {
		esQuery.
			SetMultiMatch([]string{"author_id", "nickname", "unique_id", "short_id"}, keyword)
	}
	results := esMultiQuery.
		SetTable(esTable).
		SetCache(60).
		AddMust(esQuery.Condition).
		SetLimit(0, 4).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("follower_count", "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	return
}

//商品达人分析
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

//带货达人榜聚合统计
func (receiver *EsAuthorBusiness) SaleAuthorRankCount(startTime time.Time, dateType int, tags, sortStr, orderBy string, verified, page, pageSize int) ([]interface{}, int, global.CommonError) {
	if sortStr == "" {
		sortStr = "sum_gmv"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"sum_gmv", "sum_sales", "avg_price"}) {
		return nil, 0, global.NewError(4004)
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		return nil, 0, global.NewError(4004)
	}
	esQuery, _ := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("predict_gmv", map[string]interface{}{
		"gte": 100000,
	})
	if tags != "" {
		esQuery.SetTerm("tags.keyword", tags)
	}
	if verified == 1 {
		esQuery.SetTerm("verification_type", 2)
	}
	var esTable string
	cacheTime := 600 * time.Second
	today := time.Now().Format("20060102")
	switch dateType {
	case 1:
		date := startTime.Format("20060102")
		if date != today {
			cacheTime = 86400
		}
		esTable = fmt.Sprintf(es.DyAuthorTakeGoodsTopTable, date)
	case 2:
		endDate := startTime.AddDate(0, 0, 6)
		if endDate.Format("20060102") != today {
			cacheTime = 86400
		}
		esTable = GetESTableByDayTime(es.DyAuthorTakeGoodsTopTable, startTime, endDate)
	case 3:
		esTable = fmt.Sprintf(es.DyAuthorTakeGoodsTopTable+"*", startTime.Format("200601"))
	}
	countResult := elasticsearch.NewElasticMultiQuery().
		SetCache(cacheTime).
		SetTable(esTable).RawQuery(map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": esQuery.Condition,
			},
		},
		"size": 0,
		"aggs": map[string]interface{}{
			"authors": map[string]interface{}{
				"aggs": map[string]interface{}{
					"sum_gmv": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "predict_gmv",
						},
					},
					"sum_sales": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "predict_sales",
						},
					},
					"hit": map[string]interface{}{
						"top_hits": map[string]interface{}{
							"size": 100,
						},
					},
					"avg_price": map[string]interface{}{
						"bucket_script": map[string]interface{}{
							"buckets_path": map[string]interface{}{
								"all_gmv":   "sum_gmv",
								"all_sales": "sum_sales",
							},
							"script": map[string]interface{}{
								"source": "params.all_gmv / params.all_sales",
								"lang":   "painless",
							},
						},
					},
					"r_bucket_sort": map[string]interface{}{
						"bucket_sort": map[string]interface{}{
							"sort": map[string]interface{}{
								sortStr: map[string]interface{}{
									"order": orderBy,
								},
							},
							"from": (page - 1) * pageSize,
							"size": pageSize,
						},
					},
				},
				"composite": map[string]interface{}{
					"size": 10000,
					"sources": map[string]interface{}{
						"author_id": map[string]interface{}{
							"terms": map[string]interface{}{
								"field": "author_id.keyword",
							},
						},
					},
				},
			},
		},
	})
	res := elasticsearch.GetBuckets(countResult, "authors")
	var total int
	if countResult["hits"] != nil && countResult["hits"].(map[string]interface{})["total"] != nil {
		total = int(countResult["hits"].(map[string]interface{})["total"].(float64))
	}
	return res, total, nil
}

//达人涨粉榜
func (receiver *EsAuthorBusiness) DyAuthorFollowerIncRank(date, tags, province, city, sortStr, orderBy string, isDelivery, page, pageSize int) (list []es.DyAuthorFollowerTop, total int, comErr global.CommonError) {
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	if sortStr == "" {
		sortStr = "inc_follower_count"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sortStr, []string{"live_inc_follower_count", "inc_follower_count", "aweme_inc_follower_count", "follower_count"}) {
		comErr = global.NewError(4000)
		return
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		comErr = global.NewError(4000)
		return
	}
	if tags != "" {
		esQuery.SetMatchPhrase("tags", tags)
	}
	if isDelivery != 0 {
		esQuery.SetTerm("is_delivery", isDelivery)
	}
	if province != "" {
		esQuery.SetTerm("province.keyword", province)
	}
	if city != "" {
		esQuery.SetTerm("city.keyword", city)
	}
	esTable := fmt.Sprintf(es.DyAuthorFollowerTable, date)
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
