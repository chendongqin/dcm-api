package es

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/models/es"
	"dongchamao/services/elasticsearch"
	"fmt"
	"math"
	"strings"
	"time"
)

type EsLiveBusiness struct {
}

func NewEsLiveBusiness() *EsLiveBusiness {
	return new(EsLiveBusiness)
}

//达人直播间搜索
func (receiver *EsLiveBusiness) SearchAuthorRooms(authorId, keyword, sort, orderBy string, page, size int, startDate, endDate time.Time) (list []es.EsAuthorLiveRoom, total int, comErr global.CommonError) {
	if sort == "" {
		sort = "create_timestamp"
	}
	if orderBy == "" {
		orderBy = "desc"
	}
	if !utils.InArrayString(sort, []string{"create_timestamp", "gmv", "sales", "max_user_count"}) {
		comErr = global.NewError(4000)
		return
	}
	if !utils.InArrayString(orderBy, []string{"desc", "asc"}) {
		comErr = global.NewError(4000)
		return
	}
	if size > 50 {
		comErr = global.NewError(4000)
		return
	}
	//兼容数据 2021-06-29
	firstDay, _ := time.ParseInLocation("20060102", "20210701", time.Local)
	if startDate.Before(firstDay) {
		startDate = firstDay
	}
	tableArr := make([]string, 0)
	begin := startDate
	for {
		if begin.After(endDate) {
			break
		}
		tableArr = append(tableArr, fmt.Sprintf(es.DyAuthorLiveRecords, begin.Format("20060102")))
		begin = begin.AddDate(0, 0, 1)
	}
	if len(tableArr) == 0 {
		return
	}
	esTable := strings.Join(tableArr, ",")
	esQuery, esMultiQuery := elasticsearch.NewElasticQueryGroup()
	esQuery.SetTerm("author_id", authorId)
	if keyword != "" {
		esQuery.AddCondition(map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match_phrase": map[string]interface{}{
							"title": keyword,
						},
					},
					map[string]interface{}{
						"match_phrase": map[string]interface{}{
							"product_title": keyword,
						},
					},
				},
			},
		})
	}
	results := esMultiQuery.
		SetTable(esTable).
		AddMust(esQuery.Condition).
		SetLimit((page-1)*size, size).
		SetOrderBy(elasticsearch.NewElasticOrder().Add(sort, orderBy).Order).
		SetMultiQuery().
		Query()
	utils.MapToStruct(results, &list)
	for k, v := range list {
		list[k].Sales = math.Floor(v.Sales)
	}
	total = esMultiQuery.Count
	return
}
