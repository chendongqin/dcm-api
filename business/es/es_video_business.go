package es

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/models/repost/dy"
	"dongchamao/services/elasticsearch"
	"math"
	"strings"
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
	esTable, connection, err := GetESTableByMonthTime(es.DyAuthorAwemeProductTable, startTime, endTime)
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
	tmpPage := page
	tmpPageSize := pageSize
	if keyword != "" {
		tmpPage = 1
		tmpPageSize = 10000
	}
	var cacheTime time.Duration = 180
	result := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		SetLimit((tmpPage-1)*tmpPageSize, tmpPageSize).
		SetMultiQuery().
		Query()
	utils.MapToStruct(result, &list)
	total = esMultiQuery.Count
	if keyword != "" {
		keyword = strings.ToLower(keyword)
		newList := []es.DyProductVideo{}
		for _, v := range list {
			if strings.Index(strings.ToLower(v.AwemeTitle), keyword) < 0 && strings.Index(strings.ToLower(v.Nickname), keyword) < 0 {
				continue
			}
			newList = append(newList, v)
		}
		total = len(newList)
		if total == 0 {
			list = newList
			return
		}
		start := (page - 1) * pageSize
		end := start + pageSize
		if total < end {
			end = total
		}
		list = newList[start:end]
	}

	return
}

func (e *EsVideoBusiness) SearchAwemeByProductTotal(productId, keyword string,
	startTime, endTime time.Time) (totalSales int64, totalGmv float64, comErr global.CommonError) {
	esTable, connection, err := GetESTableByMonthTime(es.DyAuthorAwemeProductTable, startTime, endTime)
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
	var cacheTime time.Duration = 180
	if keyword != "" {
		keyword = strings.ToLower(keyword)
		list := []es.DyProductVideo{}
		esQuery.SetExist("filed", "aweme_gmv")
		results := esMultiQuery.
			SetConnection(connection).
			SetTable(esTable).
			SetCache(cacheTime).
			AddMust(esQuery.Condition).
			SetOrderBy(elasticsearch.NewElasticOrder().Add("aweme_gmv", "desc").Order).
			SetLimit(0, 10000).
			SetMultiQuery().
			Query()
		utils.MapToStruct(results, &list)
		for _, v := range list {
			if strings.Index(strings.ToLower(v.AwemeTitle), keyword) < 0 && strings.Index(strings.ToLower(v.Nickname), keyword) < 0 {
				continue
			}
			totalGmv += v.AwemeGmv
			totalSales += v.Sales
		}
		return
	}
	countResult := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
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
	//if keyword != "" {
	//	esQuery.SetMatchPhrase("aweme_title", keyword)
	//}
	if hasProduct == 1 {
		esQuery.SetExist("field", "product_ids")
	}
	var cacheTime time.Duration = 180
	tmpPage := page
	tmpPageSize := pageSize
	if keyword != "" {
		tmpPage = 1
		tmpPageSize = 10000
	}
	//if keyword != "" {
	//	esQuery.SetMatchPhrase("aweme_title", keyword)
	//}
	result := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, orderBy).Order).
		//SetLimit((page-1)*pageSize, pageSize).
		SetLimit((tmpPage-1)*tmpPageSize, tmpPageSize).
		SetMultiQuery().
		Query()
	utils.MapToStruct(result, &list)
	total = esMultiQuery.Count
	if keyword != "" {
		keyword = strings.ToLower(keyword)
		newList := []es.DyAweme{}
		for _, v := range list {
			if strings.Index(strings.ToLower(v.AwemeTitle), keyword) < 0 {
				continue
			}
			newList = append(newList, v)
		}
		total = len(newList)
		if total == 0 {
			list = newList
			return
		}
		start := (page - 1) * pageSize
		end := start + pageSize
		if total < end {
			end = total
		}
		list = newList[start:end]
	}

	return
}

//获取达人视频列表
func (e *EsVideoBusiness) SearchByAuthorTotal(authorId, keyword string, hasProduct int, startTime, endTime time.Time) (totalSales int64, totalGmv float64, comErr global.CommonError) {
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
	//if keyword != "" {
	//	esQuery.SetMatchPhrase("aweme_title", keyword)
	//}
	if hasProduct == 1 {
		esQuery.SetExist("field", "product_ids")
	}
	var cacheTime time.Duration = 180
	if keyword != "" {
		keyword = strings.ToLower(keyword)
		list := []es.DyAweme{}
		esQuery.SetExist("field", "aweme_gmv")
		result := esMultiQuery.
			SetConnection(connection).
			SetTable(esTable).
			SetCache(cacheTime).
			AddMust(esQuery.Condition).
			SetLimit(0, 10000).
			SetMultiQuery().
			Query()
		utils.MapToStruct(result, &list)
		for _, v := range list {
			if strings.Index(strings.ToLower(v.AwemeTitle), keyword) < 0 {
				continue
			}
			totalSales += v.Sales
			totalGmv += v.AwemeGmv
		}
		return
	}
	countResult := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
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
	var cacheTime time.Duration = 300
	countResult := esMultiQuery.
		SetCache(cacheTime).
		SetConnection(connection).
		SetTable(esTable).
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
						"size":  10000,
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
	esTable, connection, err := GetESTableByMonthTime(es.DyAuthorAwemeProductTable, startTime, endTime)
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
	esTable, connection, err := GetESTableByMonthTime(es.DyAuthorAwemeProductTable, startTime, endTime)
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

//获取达人带货视频聚合
func (e *EsVideoBusiness) AuthorProductAwemeSumList(authorId, productId, shopId, sortStr, orderBy string, startTime, endTime time.Time, page, pageSize int) (list []es.DyProductAwemeSum, total int, comErr global.CommonError) {
	esTable, connection, err := GetESTableByMonthTime(es.DyAuthorAwemeProductTable, startTime, endTime)
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
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
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
								"field": "aweme_gmv",
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

func (e *EsVideoBusiness) SearchAwemeAuthor(productId, shopId, tag string, minFollow, maxFollow int64,
	startTime, endTime time.Time, scoreType int) (list []es.DyProductVideo, total int, comErr global.CommonError) {
	esTable, connection, err := GetESTableByMonthTime(es.DyAuthorAwemeProductTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	if productId != "" {
		esQuery.SetTerm("product_id", productId)
	}
	if shopId != "" {
		esQuery.SetTerm("shop_id", shopId)
	}
	if scoreType > 0 {
		esQuery.SetTerm("level", scoreType)
	}
	if tag != "" {
		esQuery.SetMatchPhrase("tags", tag)
	}
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if minFollow > 0 || maxFollow > 0 {
		rangeMap := map[string]interface{}{}
		if minFollow > 0 {
			rangeMap["gte"] = minFollow
		}
		if maxFollow > 0 {
			rangeMap["lt"] = maxFollow
		}
		esQuery.SetRange("follower_count", rangeMap)
	}
	var cacheTime time.Duration = 600
	result := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("aweme_create_time", "desc").Order).
		SetLimit(0, 10000).
		SetMultiQuery().
		Query()
	utils.MapToStruct(result, &list)
	total = esMultiQuery.Count
	//if keyword != "" {
	//	keyword = strings.ToLower(keyword)
	//	newList := []es.DyProductVideo{}
	//	for _, v := range list {
	//		if strings.Index(strings.ToLower(v.AwemeTitle), keyword) < 0 && strings.Index(strings.ToLower(v.Nickname), keyword) < 0 {
	//			continue
	//		}
	//		newList = append(newList, v)
	//	}
	//	list = newList
	//}
	return
}

func (receiver *EsVideoBusiness) SumSearchAwemeAuthor(productId, shopId string, startTime, endTime time.Time) (list []es.SumProductVideoAuthor, total int, comErr global.CommonError) {
	esTable, connection, err := GetESTableByMonthTime(es.DyAuthorAwemeProductTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
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
		RawQuery(map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": esQuery.Condition,
				},
			},
			"size": 0,
			"aggs": map[string]interface{}{
				"authors": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "author_id.keyword",
						"size":  10000,
					},
					"aggs": map[string]interface{}{
						"data": map[string]interface{}{
							"top_hits": map[string]interface{}{
								"sort": []map[string]interface{}{
									{
										"aweme_create_time": map[string]interface{}{
											"order": "desc",
										},
									},
								},
								"size": 1,
							},
						},
						"sales": map[string]interface{}{
							"sum": map[string]interface{}{
								"field": "sales",
							},
						},
						"aweme_gmv": map[string]interface{}{
							"sum": map[string]interface{}{
								"field": "aweme_gmv",
							},
						},
						"aweme_create_time": map[string]interface{}{
							"max": map[string]interface{}{
								"field": "aweme_create_time",
							},
						},
						"r_sort": map[string]interface{}{
							"bucket_sort": map[string]interface{}{
								"sort": []map[string]interface{}{
									{
										"aweme_create_time": map[string]interface{}{
											"order": "desc",
										},
									},
								},
							},
						},
					},
				},
			},
		})
	buckets := elasticsearch.GetBuckets(results, "authors")
	utils.MapToStruct(buckets, &list)
	total = len(buckets)
	return
}

func (receiver *EsVideoBusiness) CountSearchAuthorAwemeProductNum(productId, shopId string,
	authorIds []string, startTime, endTime time.Time) (awemeMap map[string]int, productMap map[string]int, comErr global.CommonError) {
	if len(authorIds) == 0 {
		return
	}
	esTable, connection, err := GetESTableByMonthTime(es.DyAuthorAwemeProductTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("aweme_create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if productId != "" {
		esQuery.SetTerm("product_id", productId)
	}
	if shopId != "" {
		esQuery.SetTerm("shop_id", shopId)
	}
	esQuery.SetTerms("author_id", authorIds)
	var cacheTime time.Duration = 180
	aggsMap := map[string]interface{}{
		"awemes": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "aweme_id.keyword",
				"size":  10000,
			},
		},
	}
	if shopId != "" {
		aggsMap["products"] = map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "product_id.keyword",
				"size":  10000,
			},
		}
	}
	var outTime = 10 * time.Second
	esMultiQuery.Timeout = &outTime
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		RawQuery(map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": esQuery.Condition,
				},
			},
			"size": 0,
			"aggs": map[string]interface{}{
				"authors": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "author_id.keyword",
						"size":  100,
					},
					"aggs": aggsMap,
				},
			},
		})
	buckets := elasticsearch.GetBuckets(results, "authors")
	list := []es.CountAuthorProductAweme{}
	utils.MapToStruct(buckets, &list)
	awemeMap = map[string]int{}
	productMap = map[string]int{}
	for _, v := range list {
		awemeMap[v.Key] = len(v.Awemes.Buckets)
		productMap[v.Key] = len(v.Products.Buckets)
	}
	return
}
