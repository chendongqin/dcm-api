package elasticsearch

import (
	"context"
	"crypto/md5"
	"dongchamao/global"
	"dongchamao/global/utils"
	//"dongchamao/services/cmmlog"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/bitly/go-simplejson"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	jsoniter "github.com/json-iterator/go"
	//"runtime"
	"strings"
	"time"
)

var SlowTime = time.Millisecond * 500 //慢日志记录时间

var DebugShow = true //是否显示调试输出

var SlowLogShow = true //是否显示慢日志

type Meta []map[string]interface{}

type SearchIds []string

type SuggestItem struct {
	Input  string `json:"input"`
	Weight int    `json:"weight"`
}

type ElasticMultiQuery struct {
	Connection         string //设置使用哪个客户端
	Must               []map[string]interface{}
	Filter             []map[string]interface{}
	FilterBool         *Boolean
	MustNot            []map[string]interface{}
	MinimumShouldMatch int
	Collapse           string
	FunctionScore      *FunctionScore
	Should             []map[string]interface{}
	OrderBy            []map[string]interface{}
	Source             map[string]interface{}
	MinScore           float64
	Limit              map[string]int
	SearchAfter        []interface{}
	QueryString        string
	cacheEnabled       bool
	cacheTime          time.Duration
	Timeout            *time.Duration
	tableName          string
	Count              int
	Preference         string
	FilterPath         []string
	IsDeepPage         bool
	IsTimedOut         bool
}

type Boolean struct {
	Must               []map[string]interface{}
	MustNot            []map[string]interface{}
	Should             []map[string]interface{}
	MinimumShouldMatch int
}

func NewBoolean() *Boolean {
	return &Boolean{
		Must:    make([]map[string]interface{}, 0),
		MustNot: make([]map[string]interface{}, 0),
		Should:  make([]map[string]interface{}, 0),
	}
}

func (b *Boolean) AddMust(condition []map[string]interface{}) *Boolean {
	if len(b.Must) <= 0 {
		b.Must = make([]map[string]interface{}, 0)
	}
	b.Must = append(b.Must, condition...)
	return b
}

func (b *Boolean) AddMustNot(condition []map[string]interface{}) *Boolean {
	if len(b.MustNot) <= 0 {
		b.MustNot = make([]map[string]interface{}, 0)
	}
	b.MustNot = append(b.MustNot, condition...)
	return b
}

func (b *Boolean) AddShould(condition []map[string]interface{}) *Boolean {
	if len(b.Should) <= 0 {
		b.Should = make([]map[string]interface{}, 0)
	}
	b.Should = append(b.Should, condition...)
	return b
}

func (b *Boolean) SetMinimumShouldMatch(val int) *Boolean {
	b.MinimumShouldMatch = val
	return b
}

func (b *Boolean) Build() map[string]interface{} {
	condition := make(map[string]interface{}, 0)

	if len(b.Must) > 0 {
		condition["must"] = b.Must
	}
	if len(b.MustNot) > 0 {
		condition["must_not"] = b.MustNot
	}
	if len(b.Should) > 0 {
		condition["should"] = b.Should
	}
	if b.MinimumShouldMatch > 0 {
		condition["minimum_should_match"] = b.MinimumShouldMatch
	}
	query := map[string]interface{}{
		"bool": condition,
	}
	return query
}

func (this *ElasticMultiQuery) Copy() *ElasticMultiQuery {
	esMultiQuery := *this
	return &esMultiQuery
}

func (this *ElasticMultiQuery) SetConnection(connection string) *ElasticMultiQuery {
	this.Connection = connection
	return this
}

func (this *ElasticMultiQuery) SetPreference(preference string) *ElasticMultiQuery {
	this.Preference = preference
	return this
}

func (this *ElasticMultiQuery) SetTimeout(timeout time.Duration) *ElasticMultiQuery {
	this.Timeout = &timeout
	return this
}

func (this *ElasticMultiQuery) SetFilterPath(v ...string) *ElasticMultiQuery {
	this.FilterPath = v
	return this
}

func (this *ElasticMultiQuery) SetTable(tbname string) *ElasticMultiQuery {
	this.tableName = tbname
	return this
}

func (this *ElasticMultiQuery) SetMust(condition []map[string]interface{}) *ElasticMultiQuery {
	this.Must = condition
	return this
}

func (this *ElasticMultiQuery) SetFilter(condition []map[string]interface{}) *ElasticMultiQuery {
	this.Filter = condition
	return this
}

func (this *ElasticMultiQuery) SetFilterBool(boolean *Boolean) *ElasticMultiQuery {
	this.FilterBool = boolean
	return this
}

func (this *ElasticMultiQuery) AddFilter(condition []map[string]interface{}) *ElasticMultiQuery {
	if len(this.Filter) <= 0 {
		this.Filter = make([]map[string]interface{}, 0)
	}
	this.Filter = append(this.Filter, condition...)
	return this
}

func (this *ElasticMultiQuery) SetCollapse(field string) *ElasticMultiQuery {
	this.Collapse = field
	return this
}

func (this *ElasticMultiQuery) SetMustNot(condition []map[string]interface{}) *ElasticMultiQuery {
	this.MustNot = condition
	return this
}

func (this *ElasticMultiQuery) AddMustNot(condition []map[string]interface{}) *ElasticMultiQuery {
	if len(this.MustNot) <= 0 {
		this.MustNot = make([]map[string]interface{}, 0)
	}
	this.MustNot = append(this.MustNot, condition...)
	return this
}

func (this *ElasticMultiQuery) AddMust(condition []map[string]interface{}) *ElasticMultiQuery {
	if len(this.Must) <= 0 {
		this.Must = make([]map[string]interface{}, 0)
	}
	this.Must = append(this.Must, condition...)
	return this
}

func (this *ElasticMultiQuery) SetShould(condition []map[string]interface{}) *ElasticMultiQuery {
	this.Should = condition
	return this
}

func (this *ElasticMultiQuery) AddShould(condition []map[string]interface{}) *ElasticMultiQuery {
	if len(this.Should) <= 0 {
		this.Should = make([]map[string]interface{}, 0)
	}
	this.Should = append(this.Should, condition...)
	return this
}

func (this *ElasticMultiQuery) SetOrderBy(order []map[string]interface{}) *ElasticMultiQuery {
	this.OrderBy = order
	return this
}

func (this *ElasticMultiQuery) SetSource(source map[string]interface{}) *ElasticMultiQuery {
	this.Source = source
	return this
}

func (this *ElasticMultiQuery) SetFields(fields ...string) *ElasticMultiQuery {
	this.Source = map[string]interface{}{
		"includes": fields,
	}
	return this
}

func (this *ElasticMultiQuery) SetLimit(start, pagesize int) *ElasticMultiQuery {
	this.Limit = map[string]int{
		"from":     start,
		"pagesize": pagesize,
	}
	return this
}

func (this *ElasticMultiQuery) RawQueryWithMap(jsonMap map[string]interface{}) map[string]interface{} {
	query, _ := json.Marshal(jsonMap)
	queryString := string(query)
	return this.RawQuery(queryString)
}

func (this *ElasticMultiQuery) Json(result map[string]interface{}) (jsonObject *simplejson.Json, err error) {
	jsonStr, err := jsoniter.Marshal(result)
	if err != nil {
		return
	}
	jsonObject, err = simplejson.NewJson(jsonStr)
	return
}

func (this *ElasticMultiQuery) RawQuery(queryString interface{}) map[string]interface{} {
	if val, ok := queryString.(string); ok {
		this.QueryString = val
	} else {
		jsonByte, err := json.Marshal(queryString)
		if err != nil {
			return nil
		}
		this.QueryString = string(jsonByte)
	}

	retMap := make(map[string]interface{}, 0)
	if this.tableName == "" {
		logs.Warn("Not specify an index name.")
		return retMap
	}

	resp := make(map[string]interface{}, 0)
	cacheKey := this._cacheKey(this.tableName, this.QueryString)
	if this.cacheEnabled == true {
		cacheData := global.Cache.Get(cacheKey)
		if cacheData != "" {
			if err := json.Unmarshal([]byte(cacheData), &resp); err == nil {
				return resp
			} else {
				logs.Error("err: %s", err)
			}
		}
	}

	res, err := this.elasticQuery()
	if err != nil {
		return retMap
	}
	defer res.Body.Close()

	respStr := res.String()
	respIsError := res.IsError()
	if respIsError == true {
		return resp
	}
	respStr = strings.Trim(respStr, "[200 OK]")

	if err := json.Unmarshal([]byte(respStr), &resp); err == nil {
		if this.cacheEnabled == true {
			_ = global.Cache.SetMap(cacheKey, resp, this.cacheTime)
		}
	}
	return resp
}

func (this *ElasticMultiQuery) SetFunctionScore(score *FunctionScore) *ElasticMultiQuery {
	this.FunctionScore = score
	return this
}

func (this *ElasticMultiQuery) SetMinimumShouldMatch(minimumShouldMatch int) *ElasticMultiQuery {
	this.MinimumShouldMatch = minimumShouldMatch
	return this
}

func (this *ElasticMultiQuery) SetMultiQuery() *ElasticMultiQuery {
	query := make(map[string]interface{})
	if len(this.Filter) <= 0 && len(this.Must) <= 0 && len(this.MustNot) <= 0 && len(this.Should) <= 0 && this.FilterBool == nil {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	} else {
		condition := make(map[string]interface{})
		if len(this.Must) > 0 {
			condition["must"] = this.Must
		}
		if len(this.Filter) > 0 {
			condition["filter"] = this.Filter
		}
		if this.FilterBool != nil {
			condition["filter"] = this.FilterBool.Build()
		}
		if len(this.MustNot) > 0 {
			condition["must_not"] = this.MustNot
		}
		if len(this.Should) > 0 {
			condition["should"] = this.Should
		}
		if this.MinimumShouldMatch > 0 {
			condition["minimum_should_match"] = this.MinimumShouldMatch
		}
		if this.FunctionScore == nil {
			query = map[string]interface{}{
				"query": map[string]interface{}{
					"bool": condition,
				},
			}
		} else {
			data := make(map[string]interface{})
			data["query"] = map[string]interface{}{
				"bool": condition,
			}
			functionScore := this.FunctionScore.Build()
			for key, val := range functionScore {
				data[key] = val
			}
			query = map[string]interface{}{
				"query": map[string]interface{}{
					"function_score": data,
				},
			}
		}
	}
	if len(this.OrderBy) > 0 {
		query["sort"] = this.OrderBy
	}
	if len(this.Source) > 0 {
		query["_source"] = this.Source
	}

	if len(this.Limit) > 0 {
		this.IsDeepPage = false
		if this.Limit["from"] >= 0 {
			query["from"] = this.Limit["from"]
		}
		query["size"] = this.Limit["pagesize"]
		from := utils.ToInt64(query["from"])
		size := utils.ToInt64(query["size"])
		if from > 10000 {
			from = 0
			size = 0
			this.IsDeepPage = true
		} else if from+size > 10000 {
			size = 10000 - from
		}
		query["from"] = from
		query["size"] = size
	}

	if this.MinScore > 0 {
		query["min_score"] = this.MinScore
	}

	if len(this.SearchAfter) > 0 {
		query["search_after"] = this.SearchAfter
	}

	if this.Collapse != "" {
		query["collapse"] = map[string]interface{}{
			"field": this.Collapse,
		}
	}

	querys, err := json.Marshal(query)
	if err != nil {
		this.QueryString = ""
	} else {
		this.QueryString = string(querys)
	}
	return this
}

func (this *ElasticMultiQuery) SetMinScore(score float64) *ElasticMultiQuery {
	this.MinScore = score
	return this
}

func (this *ElasticMultiQuery) SetCache(cacheTime time.Duration) *ElasticMultiQuery {
	this.cacheEnabled = true
	this.cacheTime = cacheTime
	return this
}

func (this *ElasticMultiQuery) SetSearchAfter(searchAfter []interface{}) *ElasticMultiQuery {
	this.SearchAfter = searchAfter
	return this
}

func (this *ElasticMultiQuery) FindCount() (int64, error) {
	startTime := time.Now()
	defer func() {
		if DebugShow {
			logs.Info("elasticsearch find count spend time:%s", time.Since(startTime))
		}
	}()
	this.Limit = make(map[string]int)
	this.OrderBy = make([]map[string]interface{}, 0)
	this.SetMultiQuery()
	if DebugShow {
		logs.Info("elastic search find count query string:%s", this.QueryString)
	}
	resp := &ElasticCountResp{}
	cacheKey := this._cacheKey(this.tableName, "count"+this.QueryString)
	if this.cacheEnabled == true {
		cacheData := global.Cache.Get(cacheKey)
		if cacheData != "" {
			if err := json.Unmarshal([]byte(cacheData), resp); err == nil {
				return resp.Count, nil
			}
		}
	}

	client := global.GetEsClient(this.Connection)

	res, err := client.Count(
		client.Count.WithContext(context.Background()),
		client.Count.WithIndex(this.tableName),
		client.Count.WithBody(strings.NewReader(this.QueryString)),
		client.Count.WithPretty(),
	)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(resp); err == nil {
		if this.cacheEnabled == true {
			_ = global.Cache.SetMap(cacheKey, resp, this.cacheTime)
		}
		return resp.Count, nil
	} else {
		return 0, err
	}
}

func (this *ElasticMultiQuery) Query() []map[string]interface{} {
	retMap := make([]map[string]interface{}, 0)
	if this.tableName == "" {
		return retMap
	}
	if this.QueryString == "" {
		this.QueryString = `{"query":{"match_all" : {}}}`
	}

	resp := &ElasticResp{}
	cacheKey := this._cacheKey(this.tableName, this.QueryString)
	if this.cacheEnabled == true {
		cacheData := global.Cache.Get(cacheKey)
		if cacheData != "" {
			if err := json.Unmarshal([]byte(cacheData), resp); err == nil {
				this.Count = resp.Hits.Total
				this.IsTimedOut = resp.TimedOut
				if this.Count > 0 {
					for _, v := range resp.Hits.Hits {
						source := v.Source
						source["id"] = v.Id
						retMap = append(retMap, source)
					}
					return retMap
				}
			}
		}
	}

	res, err := this.elasticQuery()
	if err != nil {
		return retMap
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(resp); err == nil {
		if this.cacheEnabled == true && !resp.TimedOut {
			_ = global.Cache.SetMap(cacheKey, resp, this.cacheTime)
		}
		this.Count = resp.Hits.Total
		this.IsTimedOut = resp.TimedOut
		for _, v := range resp.Hits.Hits {
			source := v.Source
			source["id"] = v.Id
			retMap = append(retMap, source)
		}
		return retMap
	}
	return retMap
}

//获取es返回的所有字段
func (this *ElasticMultiQuery) QueryRawResp() (jsonObj *simplejson.Json) {
	retMap := make([]map[string]interface{}, 0)
	if this.tableName == "" {
		return
	}
	if this.QueryString == "" {
		this.QueryString = `{"query":{"match_all" : {}}}`
	}

	resp := &ElasticResp{}
	cacheKey := this._cacheKey(this.tableName, this.QueryString)
	if this.cacheEnabled == true {
		cacheData := global.Cache.Get(cacheKey)
		if cacheData != "" {
			if err := json.Unmarshal([]byte(cacheData), resp); err == nil {
				this.Count = resp.Hits.Total
				for _, v := range resp.Hits.Hits {
					retMap = append(retMap, v.Source)
				}
				return
			}
		}
	}

	res, err := this.elasticQuery()
	if err != nil {
		return
	}
	defer res.Body.Close()

	jsonObj, _ = simplejson.NewFromReader(res.Body)

	return jsonObj
}

func (this *ElasticMultiQuery) elasticQuery() (*esapi.Response, error) {

	//return nil, errors.New("服务故障")

	if this.IsDeepPage {
		return nil, errors.New("超过翻页最大限制")
	}
	defaultTimeout := 5 * time.Second
	if this.Timeout != nil {
		defaultTimeout = *this.Timeout
	}
	startTime := time.Now()
	//慢记录
	defer func() {
		spendTime := time.Since(startTime)
		if spendTime >= SlowTime && SlowLogShow {
			logs.Warn("[elasticsearch] [slow query] [table:%s] [querystring:%s] [took: %s]", this.tableName, this.QueryString, spendTime)
		}
	}()
	var res *esapi.Response
	if this.Preference == "hash" {
		this.Preference = utils.Md5_encode(this.QueryString)
	}
	client := global.GetEsClient(this.Connection)
	//err := retry.Do(func() error {
	var err error
	searchRequestFunc := []func(request *esapi.SearchRequest){client.Search.WithContext(context.Background()),
		client.Search.WithIndex(this.tableName),
		client.Search.WithTimeout(defaultTimeout),
		client.Search.WithBody(strings.NewReader(this.QueryString)),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
		client.Search.WithFilterPath(this.FilterPath...),
	}
	if this.Preference != "" {
		logs.Info("[elasticsearch] enable preference: [%s]", this.Preference)
		searchRequestFunc = append(searchRequestFunc, client.Search.WithPreference(this.Preference))
	}
	res, err = client.Search(searchRequestFunc...)
	if err != nil && DebugShow {
		logs.Error("[elasticsearch] [err:%s] [querystring:%s]", err, this.QueryString)
	}
	//	return err
	//}, 3, time.Millisecond*100)
	utils.DevRun(func() {
		if DebugShow {
			logs.Debug("[elasticsearch] [index:%s] [querystring:%s] [took:%s]", this.tableName, this.QueryString, time.Since(startTime))
		}
	})
	if err != nil {
		logs.Error("[elasticsearch] [err:%s], [querystring:%s]", err, this.QueryString)
	}
	//spendTime := time.Since(startTime)
	//_, file, line, ok := runtime.Caller(3)
	//if !ok {
	//	file = "unknown"
	//	line = 0
	//}
	//cmmlog.EsLog(this.tableName, this.QueryString, file, line, spendTime.Seconds())
	return res, err
}

func (this *ElasticMultiQuery) QueryOne() map[string]interface{} {
	this.SetLimit(0, 1)
	retMap := this.Query()
	if len(retMap) >= 1 {
		return retMap[0]
	}
	return nil
}

func (this *ElasticMultiQuery) Get() map[string]interface{} {
	return this.QueryOne()
}

func (this *ElasticMultiQuery) _cacheKey(tableName, str string) string {
	h := md5.New()
	h.Write([]byte(tableName + str))
	return "es_cache:" + hex.EncodeToString(h.Sum(nil))
}
