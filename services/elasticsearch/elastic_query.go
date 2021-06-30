package elasticsearch

import "dongchamao/global/alias"

type Condition map[string]interface{}

type ElasticQuery struct {
	Condition []map[string]interface{}
}

func (this *ElasticQuery) SetTerm(field string, val interface{}) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["term"] = map[string]interface{}{
		field: val,
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) AddBoolean(boolean *Boolean) *ElasticQuery {
	this.Condition = append(this.Condition, boolean.Build())
	return this
}

func (this *ElasticQuery) SetTermBoost(field string, val interface{}, boost int) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["term"] = map[string]interface{}{
		field: map[string]interface{}{
			"value": val,
			"boost": boost,
		},
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) AddCondition(condition map[string]interface{}) *ElasticQuery {
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetId(id interface{}) *ElasticQuery {
	return this.SetTerm("_id", id)
}

func (this *ElasticQuery) SetExist(field string, val interface{}) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["exists"] = map[string]interface{}{
		field: val,
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetTerms(field string, val interface{}) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["terms"] = map[string]interface{}{
		field: val,
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetTermsBoost(field string, val interface{}, boost int) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["terms"] = map[string]interface{}{
		field:   val,
		"boost": boost,
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetMatch(field string, val interface{}) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["match"] = map[string]interface{}{
		field: val,
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetMatchPhrase(field string, val interface{}) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["match_phrase"] = map[string]interface{}{
		field: val,
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetMatchPhraseWithParams(field string, val interface{}, params alias.M) *ElasticQuery {
	condition := make(map[string]interface{})
	fieldQuery := alias.M{
		"query": val,
	}
	for key, val := range params {
		fieldQuery[key] = val
	}
	condition["match_phrase"] = map[string]interface{}{
		field: fieldQuery,
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetMatchPhraseBoost(field string, val interface{}, boost int) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["match_phrase"] = map[string]interface{}{
		field: map[string]interface{}{
			"query": val,
			"boost": boost,
		},
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetConstantScoreMatchPhrase(field string, val interface{}, boost int) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["constant_score"] = map[string]interface{}{
		"filter": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				field: map[string]interface{}{
					"query": val,
					"boost": boost,
				},
			},
		},
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetConstantScoreMatchPhraseWithSlop(field string, val interface{}, boost int, slop int) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["constant_score"] = map[string]interface{}{
		"filter": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				field: map[string]interface{}{
					"query": val,
					"boost": boost,
					"slop":  slop,
				},
			},
		},
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetConstantScoreMatch(field string, val interface{}, boost int, operator string) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["constant_score"] = map[string]interface{}{
		"filter": map[string]interface{}{
			"match": map[string]interface{}{
				field: map[string]interface{}{
					"query":                val,
					"boost":                boost,
					"operator":             operator,
					"analyzer":             "ik_smart",
					"minimum_should_match": "2<75% 7<100%",
				},
			},
		},
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetMultiMatch(fields []string, val interface{}, args ...string) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["multi_match"] = map[string]interface{}{
		"query":  val,
		"fields": fields,
	}
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetRawMultiMatch(val map[string]interface{}) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["multi_match"] = val
	this.Condition = append(this.Condition, condition)
	return this
}

func (this *ElasticQuery) SetRange(field string, val interface{}) *ElasticQuery {
	condition := make(map[string]interface{})
	condition["range"] = map[string]interface{}{
		field: val,
	}
	this.Condition = append(this.Condition, condition)
	return this
}
