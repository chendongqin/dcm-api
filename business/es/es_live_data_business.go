package es

import (
	"dongchamao/global"
	"dongchamao/global/alias"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/models/repost/dy"
	"dongchamao/services/elasticsearch"
	"time"
)

type EsLiveDataBusiness struct {
}

func NewEsLiveDataBusiness() *EsLiveDataBusiness {
	return new(EsLiveDataBusiness)
}

//达人直播间统计
func (receiver *EsLiveDataBusiness) SumLiveData(startTime, endTime time.Time, hasProduct, living int) (total int, data es.DyLiveDataUserSumCount) {
	data = es.DyLiveDataUserSumCount{}
	esTable, connection, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	if err != nil {
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if hasProduct == 1 {
		esQuery.SetRange("num_product", map[string]interface{}{
			"gt": 0,
		})
	}
	if living == 1 {
		esQuery.SetTerm("room_status", 2)
	}
	var cacheTime time.Duration = 300
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	countResult := esMultiQuery.
		SetConnection(connection).
		SetCache(cacheTime).
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
				"total_watch_cnt": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "watch_cnt",
					},
				},
				"total_user_count": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "user_count",
					},
				},
			},
		})
	if r, ok := countResult["aggregations"]; ok {
		utils.MapToStruct(r, &data)
	}
	if h, ok := countResult["hits"]; ok {
		if t, ok2 := h.(map[string]interface{})["total"]; ok2 {
			total = utils.ToInt(t.(float64))
		}
	}
	return
}

//商品占比查询
func (receiver *EsLiveDataBusiness) LiveCompositeByCategory(startTime, endTime time.Time, rateType, living int) (total int, res []interface{}) {
	res = []interface{}{}
	esTable, connection, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	if err != nil {
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetRange("num_product", map[string]interface{}{
		"gt": 0,
	})
	if living == 1 {
		esQuery.SetTerm("room_status", 2)
	}
	var cacheTime time.Duration = 600
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	sumMap := map[string]interface{}{}
	sumTitle := ""
	if rateType == 1 {
		sumTitle = "total_watch_cnt"
		sumMap = map[string]interface{}{
			"sum": map[string]interface{}{
				"field": "watch_cnt",
			},
		}
	} else {
		sumTitle = "total_gmv"
		sumMap = map[string]interface{}{
			"sum": map[string]interface{}{
				"field": "predict_gmv",
			},
		}
	}
	countResult := esMultiQuery.
		SetConnection(connection).
		SetCache(cacheTime).
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
				"live": map[string]interface{}{
					"aggs": map[string]interface{}{
						sumTitle: sumMap,
						"r_bucket_sort": map[string]interface{}{
							"bucket_sort": map[string]interface{}{
								"sort": map[string]interface{}{
									sumTitle: map[string]interface{}{
										"order": "desc",
									},
								},
								"from": 0,
								"size": 10000,
							},
						},
					},
					"composite": map[string]interface{}{
						"size": 10000,
						"sources": map[string]interface{}{
							"dcm_level_first": map[string]interface{}{
								"terms": map[string]interface{}{
									"field": "dcm_level_first.keyword",
								},
							},
						},
					},
				},
			},
		})
	res = elasticsearch.GetBuckets(countResult, "live")
	if h, ok := countResult["hits"]; ok {
		if t, ok2 := h.(map[string]interface{})["total"]; ok2 {
			total = utils.ToInt(t.(float64))
		}
	}
	return
}

//商品占比
func (receiver *EsLiveDataBusiness) LiveCompositeByCategoryOne(startTime, endTime time.Time, rateType, living int, category string) (total int, res interface{}) {
	res = []interface{}{}
	esTable, connection, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	if err != nil {
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetRange("num_product", map[string]interface{}{
		"gt": 0,
	})
	if living == 1 {
		esQuery.SetTerm("room_status", 2)
	}
	if category != "" {
		esQuery.SetMatchPhrase("dcm_level_first", category)
	}
	var cacheTime time.Duration = 600
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	sumMap := map[string]interface{}{}
	sumTitle := ""
	if rateType == 1 {
		sumTitle = "total_watch_cnt"
		sumMap = map[string]interface{}{
			"sum": map[string]interface{}{
				"field": "watch_cnt",
			},
		}
	} else {
		sumTitle = "total_gmv"
		sumMap = map[string]interface{}{
			"sum": map[string]interface{}{
				"field": "predict_gmv",
			},
		}
	}
	countResult := esMultiQuery.
		SetConnection(connection).
		SetCache(cacheTime).
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
				sumTitle: sumMap,
			},
		})
	if v, ok := countResult["aggregations"]; ok {
		res = v
	}
	if h, ok := countResult["hits"]; ok {
		if t, ok2 := h.(map[string]interface{})["total"]; ok2 {
			total = utils.ToInt(t.(float64))
		}
	}
	return
}

//榜单
func (receiver *EsLiveDataBusiness) LiveRankByCategory(startTime, endTime time.Time, category, sortStr string, living int) (list []es.EsDyLiveDetail, comErr global.CommonError) {
	if sortStr == "" {
		sortStr = "avg_user_count"
	}
	if !utils.InArrayString(sortStr, []string{"predict_gmv", "watch_cnt"}) {
		comErr = global.NewError(4000)
		return
	}
	esTable, connection, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if category != "" {
		esQuery.SetMatchPhrase("dcm_level_first", category)
	}
	if living == 1 {
		esQuery.SetTerm("room_status", 2)
	}
	var cacheTime time.Duration = 180
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetLimit(0, 5).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sortStr, "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	return
}

//带货行业数据分类统计
func (receiver *EsLiveDataBusiness) ProductLiveDataByCategory(startTime, endTime time.Time, category string, living int) (total int, uv, buyRate float64, data es.DyLiveDataCategorySumCount) {
	data = es.DyLiveDataCategorySumCount{}
	esTable, connection, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	if err != nil {
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetRange("num_product", map[string]interface{}{
		"gt": 0,
	})
	if category != "" {
		esQuery.SetMatchPhrase("dcm_level_first", category)
	}
	if living == 1 {
		esQuery.SetTerm("room_status", 2)
	}
	var cacheTime time.Duration = 300
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	countResult := esMultiQuery.
		SetConnection(connection).
		SetCache(cacheTime).
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
				"total_watch_cnt": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "watch_cnt",
					},
				},
				"total_user_count": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "user_count",
					},
				},
				"total_gmv": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "predict_gmv",
					},
				},
				"total_sales": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "predict_sales",
					},
				},
				"total_ticket_count": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "ticket_count",
					},
				},
			},
		})
	if r, ok := countResult["aggregations"]; ok {
		utils.MapToStruct(r, &data)
	}
	if h, ok := countResult["hits"]; ok {
		if t, ok2 := h.(map[string]interface{})["total"]; ok2 {
			total = utils.ToInt(t.(float64))
		}
	}
	uv = 0
	buyRate = 0
	if data.TotalWatchCnt.Value > 0 {
		uv = (data.TotalGmv.Value + data.TotalTicketCount.Value/10) / data.TotalWatchCnt.Value
		buyRate = data.TotalSales.Value / data.TotalWatchCnt.Value
	}
	return
}

//带货行业数据分类分级统计
func (receiver *EsLiveDataBusiness) ProductLiveDataCategoryLevel(startTime, endTime time.Time, category string, living int) (total int, data []dy.EsLiveSumDataCategoryLevel) {
	data = []dy.EsLiveSumDataCategoryLevel{}
	esTable, connection, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	if err != nil {
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetRange("num_product", map[string]interface{}{
		"gt": 0,
	})
	if category != "" {
		esQuery.SetMatchPhrase("dcm_level_first", category)
	}
	if living == 1 {
		esQuery.SetTerm("room_status", 2)
	}
	var cacheTime time.Duration = 300
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	countResult := esMultiQuery.
		SetConnection(connection).
		SetCache(cacheTime).
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
				"live": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "flow_rates_index.keyword",
						"size":  10,
					},
					"aggs": map[string]interface{}{
						"total_watch_cnt": map[string]interface{}{
							"sum": map[string]interface{}{
								"field": "watch_cnt",
							},
						},
						"total_gmv": map[string]interface{}{
							"sum": map[string]interface{}{
								"field": "predict_gmv",
							},
						},
						"stats_customer_unit_price": map[string]interface{}{
							"stats": map[string]interface{}{
								"field": "customer_unit_price",
							},
						},
						"customer_unit_price": map[string]interface{}{
							"percentiles": map[string]interface{}{
								"field":    "customer_unit_price",
								"percents": []int{50},
								"keyed":    false,
							},
						},
					},
				},
			},
		})
	res := elasticsearch.GetBuckets(countResult, "live")
	utils.MapToStruct(res, &data)
	if h, ok := countResult["hits"]; ok {
		if t, ok2 := h.(map[string]interface{})["total"]; ok2 {
			total = utils.ToInt(t.(float64))
		}
	}
	return
}

func (receiver *EsLiveDataBusiness) ProductLiveDataCategoryCustomerUnitPriceLevel(startTime, endTime time.Time, category string, living int) (data []dy.EsLiveSumDataCategoryCustomerUnitPriceLevel) {
	data = []dy.EsLiveSumDataCategoryCustomerUnitPriceLevel{}
	esTable, connection, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	if err != nil {
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetRange("num_product", map[string]interface{}{
		"gt": 0,
	})
	esQuery.SetRange("customer_unit_price", map[string]interface{}{
		"gt": 0,
	})
	if category != "" {
		esQuery.SetMatchPhrase("dcm_level_first", category)
	}
	if living == 1 {
		esQuery.SetTerm("room_status", 2)
	}
	var cacheTime time.Duration = 300
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	countResult := esMultiQuery.
		SetConnection(connection).
		SetCache(cacheTime).
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
				"live": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "flow_rates_index.keyword",
						"size":  10,
					},
					"aggs": map[string]interface{}{
						"customer_unit_price": map[string]interface{}{
							"percentiles": map[string]interface{}{
								"field":    "customer_unit_price",
								"percents": []float64{50},
								"keyed":    false,
							},
						},
					},
				},
			},
		})
	res := elasticsearch.GetBuckets(countResult, "live")
	utils.MapToStruct(res, &data)
	return
}

//带货行业数据分类分级分布数据
func (receiver *EsLiveDataBusiness) ProductLiveDataCategoryLevelTwoShow(startTime, endTime time.Time, category string, living int, keyword string) (total int, data []dy.EsLiveSumDataCategoryLevelTwo) {
	data = []dy.EsLiveSumDataCategoryLevelTwo{}
	esTable, connection, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	if err != nil {
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	esQuery.SetRange("num_product", map[string]interface{}{
		"gt": 0,
	})
	esQuery.SetRange("avg_stay_index", map[string]interface{}{
		"gt": 0,
	})
	if category != "" {
		esQuery.SetMatchPhrase("dcm_level_first", category)
	}
	if living == 1 {
		esQuery.SetTerm("room_status", 2)
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
			esQuery.SetMultiMatch([]string{"display_id", "short_id", "nickname"}, keyword)
		}
	}
	var cacheTime time.Duration = 300
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	countResult := esMultiQuery.
		SetConnection(connection).
		SetCache(cacheTime).
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
				"live": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "flow_rates.keyword",
						"size":  10000,
					},
					"aggs": map[string]interface{}{
						"live_tow": map[string]interface{}{
							"terms": map[string]interface{}{
								"field": "avg_stay_index",
							},
							"aggs": map[string]interface{}{
								"total_watch_cnt": map[string]interface{}{
									"sum": map[string]interface{}{
										"field": "watch_cnt",
									},
								},
								"total_gmv": map[string]interface{}{
									"sum": map[string]interface{}{
										"field": "predict_gmv",
									},
								},
							},
						},
					},
				},
			},
		})
	res := elasticsearch.GetBuckets(countResult, "live")
	utils.MapToStruct(res, &data)
	if h, ok := countResult["hits"]; ok {
		if t, ok2 := h.(map[string]interface{})["total"]; ok2 {
			total = utils.ToInt(t.(float64))
		}
	}
	return
}

//等级分布明细列表
func (receiver *EsLiveDataBusiness) ProductLiveDataCategoryLevelList(startTime, endTime time.Time, keyword, category, level string, stayLevel, living, page, pageSize int) (total int, list []es.EsDyLiveDetail, comErr global.CommonError) {
	esTable, connection, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	if pageSize > 30 {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if category != "" {
		esQuery.SetMatchPhrase("dcm_level_first", category)
	}
	if living == 1 {
		esQuery.SetTerm("room_status", 2)
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
			esQuery.SetMultiMatch([]string{"display_id", "short_id", "nickname"}, keyword)
		}
	}
	esQuery.SetTerm("flow_rates.keyword", level)
	esQuery.SetTerm("avg_stay_index", stayLevel)
	var cacheTime time.Duration = 180
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	results := esMultiQuery.
		SetConnection(connection).
		SetTable(esTable).
		SetCache(cacheTime).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*pageSize, pageSize).
		SetOrderBy(elasticsearch.NewElasticOrder().Add("predict_gmv", "desc").Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	total = esMultiQuery.Count
	return
}

//等级分布明细统计
func (receiver *EsLiveDataBusiness) ProductLiveDataCategoryLevelCount(startTime, endTime time.Time, keyword, category, level string, stayLevel, living int) (total int, data dy.EsLiveSumDataCategoryLevel, comErr global.CommonError) {
	data = dy.EsLiveSumDataCategoryLevel{}
	esTable, connection, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
	if err != nil {
		comErr = global.NewError(4000)
		return
	}
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetRange("create_time", map[string]interface{}{
		"gte": startTime.Unix(),
		"lt":  endTime.AddDate(0, 0, 1).Unix(),
	})
	if category != "" {
		esQuery.SetMatchPhrase("dcm_level_first", category)
	}
	if living == 1 {
		esQuery.SetTerm("room_status", 2)
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
			esQuery.SetMultiMatch([]string{"display_id", "short_id", "nickname"}, keyword)
		}
	}
	esQuery.SetTerm("flow_rates.keyword", level)
	esQuery.SetTerm("avg_stay_index", stayLevel)
	var cacheTime time.Duration = 180
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	countResult := esMultiQuery.
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
				"total_watch_cnt": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "watch_cnt",
					},
				},
				"total_gmv": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "predict_gmv",
					},
				},
				"stats_customer_unit_price": map[string]interface{}{
					"stats": map[string]interface{}{
						"field": "customer_unit_price",
					},
				},
			},
		})
	if aggregations, ok := countResult["aggregations"].(map[string]interface{}); ok {
		utils.MapToStruct(aggregations, &data)
	}
	if h, ok := countResult["hits"]; ok {
		if t, ok2 := h.(map[string]interface{})["total"]; ok2 {
			total = utils.ToInt(t.(float64))
		}
	}
	return
}
