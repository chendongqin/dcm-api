package es

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/models/repost/dy"
	"dongchamao/services/elasticsearch"
	"math"
	"time"
)

type EsVideoBusiness struct {
}

func NewEsVideoBusiness() *EsVideoBusiness {
	return new(EsVideoBusiness)
}

func (e *EsVideoBusiness) SearchAwemeByProduct(productId, keyword, sortStr, orderBy string,
	startTime, endTime time.Time, page, pageSize int) (list []es.DyProductVideo, total int, totalSales int64, totalGmv float64, comErr global.CommonError) {
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
	esTable, connection, err := GetESTableByMonthTime(es.DyProductVideoTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
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
	var cacheTime time.Duration = 180
	newEsMultiQuery := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition)
	result := newEsMultiQuery.
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetLimit((page-1)*pageSize, pageSize).
		SetMultiQuery().
		Query()
	utils.MapToStruct(result, &list)
	total = esMultiQuery.Count

	countResult := newEsMultiQuery.
		RawQuery(map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": esQuery.Condition,
				},
			},
			"size": 0,
			"aggs": map[string]interface{}{
				"total_gmv": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "aweme_gmv",
					},
				},
				"total_sales": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "sales",
					},
				},
			},
		})

	if h, ok := countResult["aggregations"]; ok {
		if t, ok2 := h.(map[string]interface{})["total_sales"]; ok2 {
			if t1, ok3 := t.(map[string]interface{})["value"]; ok3 {
				totalSales = utils.ToInt64(math.Floor(t1.(float64)))
			}
		}
		if t, ok2 := h.(map[string]interface{})["total_gmv"]; ok2 {
			if t1, ok3 := t.(map[string]interface{})["value"]; ok3 {
				totalGmv = t1.(float64)
			}
		}
	}
	return
}

//获取达人视频列表
func (e *EsVideoBusiness) SearchByAuthor(authorId, keyword, sortStr, orderBy string, hasProduct, page, pageSize int, startTime, endTime time.Time) (list []es.DyAweme, total int, totalSales int64, totalGmv float64, comErr global.CommonError) {
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
	esTable, connection, err := GetESTableByTime(es.DyVideoTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("author_id", authorId)
	esQuery.SetExist("field", "aweme_title")
	esQuery.SetTerm("exist", 1)
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if keyword != "" {
		esQuery.SetMatchPhrase("aweme_title", keyword)
	}
	if hasProduct == 1 {
		esQuery.SetExist("field", "product_ids")
	}
	var cacheTime time.Duration = 180
	newEsMultiQuery := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition)
	result := newEsMultiQuery.
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetLimit((page-1)*pageSize, pageSize).
		SetMultiQuery().
		Query()
	utils.MapToStruct(result, &list)
	total = esMultiQuery.Count

	//计算汇总数据
	countResult := newEsMultiQuery.
		RawQuery(map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": esQuery.Condition,
				},
			},
			"size": 0,
			"aggs": map[string]interface{}{
				"total_gmv": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "aweme_gmv",
					},
				},
				"total_sales": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "sales",
					},
				},
			},
		})
	if h, ok := countResult["aggregations"]; ok {
		if t, ok2 := h.(map[string]interface{})["total_sales"]; ok2 {
			if t1, ok3 := t.(map[string]interface{})["value"]; ok3 {
				totalSales = utils.ToInt64(math.Floor(t1.(float64)))
			}
		}
		if t, ok2 := h.(map[string]interface{})["total_gmv"]; ok2 {
			if t1, ok3 := t.(map[string]interface{})["value"]; ok3 {
				totalGmv = t1.(float64)
			}
		}
	}
	return
}

//统计销售额
func (e *EsVideoBusiness) SumDataByAuthor(authorId string, startTime, endTime time.Time) (countData dy.AuthorAwemeSum) {
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetTerm("author_id", authorId)
	esQuery.SetTerm("exist", 1)
	esTable, connection, err := GetESTableByTime(es.DyVideoTable, startTime, endTime)
	if err != nil {
		return
	}
	var cacheTime time.Duration = 300
	countResult := esMultiQuery.
		SetCache(cacheTime).
		SetConnection(connection).
		SetTable(esTable).
		SetMust(esQuery.Condition).
		RawQuery(map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": esQuery.Condition,
				},
			},
			"size": 0,
			"aggs": map[string]interface{}{
				"total_gmv": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "aweme_gmv",
					},
				},
				"total_sales": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "sales",
					},
				},
				"avg_comment": map[string]interface{}{
					"avg": map[string]interface{}{
						"field": "comment_count",
					},
				},
				"avg_digg": map[string]interface{}{
					"avg": map[string]interface{}{
						"field": "digg_count",
					},
				},
				"avg_share": map[string]interface{}{
					"avg": map[string]interface{}{
						"field": "share_count",
					},
				},
			},
		})
	if r, ok := countResult["aggregations"]; ok {
		data := es.DyAwemeSumCount{}
		utils.MapToStruct(r, &data)
		countData.Gmv = data.TotalGmv.Value
		countData.Sales = utils.ToInt64(math.Floor(data.TotalSales.Value))
		countData.AvgDigg = utils.ToInt64(math.Floor(data.AvgDigg.Value))
		countData.AvgShare = utils.ToInt64(math.Floor(data.AvgShare.Value))
		countData.AvgComment = utils.ToInt64(math.Floor(data.AvgComment.Value))
	}
	if hits, ok := countResult["hits"]; ok {
		if v, ok := hits.(map[string]interface{})["total"]; ok {
			countData.Total = int(v.(float64))
		}
	}
	return
}

//统计销售额
func (e *EsVideoBusiness) SumDiggByAuthors(authorIds []string, startTime, endTime time.Time) (countData map[string]float64) {
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetTerms("author_id", authorIds)
	esQuery.SetTerm("exist", 1)
	esTable, connection, err := GetESTableByTime(es.DyVideoTable, startTime, endTime)
	if err != nil {
		return
	}
	countResult := esMultiQuery.
		SetCache(300).
		SetConnection(connection).
		SetTable(esTable).
		SetMust(esQuery.Condition).
		RawQuery(map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": esQuery.Condition,
				},
			},
			"size": 0,
			"aggs": map[string]interface{}{
				"videos": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "author_id.keyword",
						"size":  100000,
					},
					"aggs": map[string]interface{}{
						"total_digg": map[string]interface{}{
							"stats": map[string]interface{}{
								"field": "digg_count",
							},
						},
					},
				},
			},
		})
	res := elasticsearch.GetBuckets(countResult, "videos")
	var dataMap []es.DyAwemeDiggCount
	countData = make(map[string]float64)
	utils.MapToStruct(res, &dataMap)
	for _, v := range dataMap {
		countData[v.Key] = v.TotalDigg.Avg
	}
	return
}

//查询带货视频数
func (e *EsVideoBusiness) CountAwemeByAuthor(authorId string, hasProduct int, startTime, endTime time.Time) (int64, error) {
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetTerm("author_id", authorId)
	if hasProduct == 1 {
		esQuery.SetExist("field", "product_ids")
	}
	esTable, connection, err := GetESTableByTime(es.DyVideoTable, startTime, endTime)
	if err != nil {
		return 0, err
	}
	var cacheTime time.Duration = 300
	return esMultiQuery.
		SetConnection(connection).
		SetCache(cacheTime).
		SetMust(esQuery.Condition).
		SetTable(esTable).FindCount()
}

//获取视频同款视频
func (e *EsVideoBusiness) SearchByProductId(productId, awemeId, keyword, sortStr, orderBy string, page, pageSize int, startTime, endTime time.Time) (list []es.DyAweme, total int, comErr global.CommonError) {
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
	if !utils.InArrayString(sortStr, []string{"aweme_create_time", "digg_count", "comment_count", "share_count"}) {
		comErr = global.NewError(4000)
		return
	}
	esTable, connection, err := GetESTableByTime(es.DyVideoTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetMatchPhrase("product_ids", productId)
	esQuery.SetExist("field", "aweme_title")
	esQuery.SetTerm("exist", 1)
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if awemeId != "" {
		esQuery.AddCondition(map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": map[string]interface{}{
					"term": map[string]interface{}{
						"aweme_id": awemeId,
					},
				},
			},
		})
	}
	if keyword != "" {
		esQuery.SetMatchPhrase("aweme_title", keyword)
	}
	var cacheTime time.Duration = 180
	result := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetLimit((page-1)*pageSize, pageSize).
		SetMultiQuery().
		Query()
	utils.MapToStruct(result, &list)
	total = esMultiQuery.Count
	return
}

//获取达人带货商品数据
func (e *EsVideoBusiness) ScanAwemeProductByAuthor(authorId, keyword, category, secondCategory, thirdCategory, brandName, shopId string, shopType int, startTime, endTime time.Time, page, pageSize int) (list []es.EsDyAuthorAwemeProduct, total int, comErr global.CommonError) {
	if shopType == 1 && shopId == "" {
		return
	}
	esTable, connection, err := GetESTableByTime(es.DyAuthorAwemeProductTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetTerm("author_id", authorId)
	if keyword != "" {
		esQuery.SetMatchPhrase("title", keyword)
	}
	if category != "" {
		esQuery.SetMatchPhrase("dcm_level_first.keyword", category)
	}
	if secondCategory != "" {
		esQuery.SetMatchPhrase("first_cname.keyword", secondCategory)
	}
	if thirdCategory != "" {
		esQuery.SetMatchPhrase("second_cname.keyword", thirdCategory)
	}
	if brandName != "" {
		esQuery.SetTerm("brand_name.keyword", brandName)
	}
	if shopType == 1 {
		esQuery.SetTerm("shop_id", shopId)
	} else if shopType == 2 {
		if shopId != "" {
			esQuery.AddCondition(map[string]interface{}{
				"bool": map[string]interface{}{
					"must_not": map[string]interface{}{
						"term": map[string]interface{}{
							"shop_id": shopId,
						},
					},
				},
			})
		}
	}
	var cacheTime time.Duration = 300
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("aweme_create_time", "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

//获取达人小店商品数据
func (e *EsVideoBusiness) ScanAwemeShopByAuthor(authorId, keyword string, startTime, endTime time.Time, page, pageSize int) (list []es.EsDyAuthorAwemeProduct, total int, comErr global.CommonError) {
	esTable, connection, err := GetESTableByTime(es.DyAuthorAwemeProductTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetTerm("author_id", authorId)
	if keyword != "" {
		esQuery.SetMatchPhrase("title", keyword)
	}
	esQuery.AddCondition(map[string]interface{}{
		"bool": map[string]interface{}{
			"must_not": map[string]interface{}{
				"term": map[string]interface{}{
					"shop_id": "",
				},
			},
		},
	})
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(300).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("aweme_create_time", "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

//获取达人带货视频聚合
func (e *EsVideoBusiness) AuthorProductAwemeSumList(authorId, productId, shopId, sortStr, orderBy string, startTime, endTime time.Time, page, pageSize int) (list []es.DyProductAwemeSum, total int, comErr global.CommonError) {
	esTable, connection, err := GetESTableByTime(es.DyAuthorAwemeProductTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if authorId != "" {
		esQuery.SetTerm("author_id", authorId)
	}
	if productId != "" {
		esQuery.SetTerm("product_id", productId)
	}
	if shopId != "" {
		esQuery.SetTerm("shop_id", shopId)

	}
	var cacheTime time.Duration = 600
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		RawQuery(map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": esQuery.Condition,
				},
			},
			"size": 0,
			"aggs": map[string]interface{}{
				"awemes": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "aweme_id.keyword",
						"size":  10000,
					},
					"aggs": map[string]interface{}{
						"total_sales": map[string]interface{}{
							"sum": map[string]interface{}{
								"field": "sales",
							},
						},
						"total_gmv": map[string]interface{}{
							"sum": map[string]interface{}{
								"field": "gmv",
							},
						},
						"r_bucket_sort": map[string]interface{}{
							"bucket_sort": map[string]interface{}{
								"sort": map[string]interface{}{
									"total_" + sortStr: map[string]interface{}{
										"order": orderBy,
									},
								},
								"from": (page - 1) * pageSize,
								"size": pageSize,
							},
						},
						"count": map[string]interface{}{
							"cardinality": map[string]interface{}{
								"field": "aweme_id.keyword",
							},
						},
					},
				},
				"count": map[string]interface{}{
					"sum_bucket": map[string]interface{}{
						"buckets_path": "awemes>count.value",
					},
				},
			},
		})
	res := elasticsearch.GetBuckets(results, "awemes")
	utils.MapToStruct(res, &list)
	total = elasticsearch.GetBucketsCount(results, "count")
	return
}

func (e *EsVideoBusiness) NewAuthorProductAwemeSumList(authorId, sortStr, orderBy string, startTime, endTime time.Time, page, pageSize int) (list []es.DyAweme, total int, comErr global.CommonError) {
	esTable, connection, err := GetESTableByTime(es.DyVideoTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if sortStr == "" {
		sortStr = "aweme_create_time"
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if authorId != "" {
		esQuery.SetTerm("author_id", authorId)
	}
	//if productId != "" {
	//	esQuery.SetMatchPhrase("product_ids", productId)
	//}
	var cacheTime time.Duration = 600
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetLimit((page-1)*pageSize, pageSize).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

//获取视频同款视频
func (e *EsVideoBusiness) GetByAwemeId(awemeId, date string) (info es.DyAweme, comErr global.CommonError) {
	esTable, connection := GetESTableByDate(es.DyVideoTable, date)
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetMatchPhrase("aweme_id", awemeId)
	esQuery.SetTerm("exist", 1)
	var cacheTime time.Duration = 180
	result := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetMultiQuery().
		QueryOne()
	utils.MapToStruct(result, &info)
	return
}
