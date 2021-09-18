package es

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/services/elasticsearch"
	"time"
)

type EsLiveDataBusiness struct {
}

func NewEsLiveDataBusiness() *EsLiveDataBusiness {
	return new(EsLiveDataBusiness)
}

//达人直播间统计
func (receiver *EsLiveDataBusiness) SumLiveData(startTime, endTime time.Time, hasProduct int) (total int, data es.DyLiveDataUserSumCount) {
	data = es.DyLiveDataUserSumCount{}
	esTable, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
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
	var cacheTime time.Duration = 300
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	countResult := esMultiQuery.
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
func (receiver *EsLiveDataBusiness) LiveCompositeByCategory(startTime, endTime time.Time, rateType int) (total int, res []interface{}) {
	res = []interface{}{}
	esTable, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
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

//榜单
func (receiver *EsLiveDataBusiness) LiveRankByCategory(startTime, endTime time.Time, category, sortStr string) (list []es.EsDyLiveRank, comErr global.CommonError) {
	if sortStr == "" {
		sortStr = "avg_user_count"
	}
	if !utils.InArrayString(sortStr, []string{"predict_gmv", "watch_cnt"}) {
		comErr = global.NewError(4000)
		return
	}
	esTable, err := GetESTableByTime(es.DyLiveInfoBaseTable, startTime, endTime)
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
	var cacheTime time.Duration = 180
	today := time.Now().Format("20060102")
	if today != endTime.Format("20060102") {
		cacheTime = 86400
	}
	results := esMultiQuery.
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
