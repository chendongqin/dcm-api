package elasticsearch

func SearchKeywordsOptimization(esMultiQuery *ElasticMultiQuery, keyword string, fields ...string) {
	esMultiQuery.AddMust(GetKeywordsOptimizationCondition(keyword, fields...))
}

func GetKeywordsOptimizationCondition(keyword string, fields ...string) []map[string]interface{} {
	esQuery := NewElasticQuery()
	esQuery.SetConstantScoreMatch(fields[0], keyword, 10, "or")

	condition := map[string]interface{}{
		"bool": map[string]interface{}{
			"should": []map[string]interface{}{
				esQuery.Condition[0],
				{
					"constant_score": map[string]interface{}{
						"filter": map[string]interface{}{
							"match_phrase": map[string]interface{}{
								fields[0]: map[string]interface{}{
									"query": keyword,
									//"analyzer": "ik_smart",
									"slop":  10,
									"boost": 100,
								},
							},
						},
					},
				},
				{
					"constant_score": map[string]interface{}{
						"filter": map[string]interface{}{
							"term": map[string]interface{}{
								"shop_name": map[string]interface{}{
									"value": keyword,
								},
							},
						},
					},
				},
			},
		},
	}

	return []map[string]interface{}{condition}
}

func GetProductKeywordsOptimizationCondition(keyword string, fields ...string) []map[string]interface{} {
	esQuery := NewElasticQuery()
	esQuery.SetConstantScoreMatch(fields[0], keyword, 10, "or")

	condition := map[string]interface{}{
		"bool": map[string]interface{}{
			"should": []map[string]interface{}{
				esQuery.Condition[0],
				{
					"constant_score": map[string]interface{}{
						"filter": map[string]interface{}{
							"match_phrase": map[string]interface{}{
								fields[0]: map[string]interface{}{
									"query": keyword,
									//"analyzer": "ik_smart",
									"slop":  10,
									"boost": 100,
								},
							},
						},
					},
				},
				{
					"constant_score": map[string]interface{}{
						"filter": map[string]interface{}{
							"term": map[string]interface{}{
								"shop_name": map[string]interface{}{
									"value": keyword,
								},
							},
						},
					},
				},
			},
		},
	}

	return []map[string]interface{}{condition}
}
