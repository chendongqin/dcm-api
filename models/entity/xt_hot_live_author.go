package entity

var XtHotLiveAuthorMap = HbaseEntity{
	"crawl_time":  {Long, "crawl_time"},
	"update_time": {Long, "update_time"},
	"data":        {AJson, "data"},
	"category":    {String, "category"},
}

type XtHotLiveAuthor struct {
	CrawlTime  int64                 `json:"crawl_time"`
	UpdateTime int64                 `json:"update_time"`
	Data       []XtHotLiveAuthorData `json:"data"`
	Tab        string                `json:"tab"`
	Category   string                `json:"category"`
}

type XtHotLiveAuthorData struct {
	AvatarUri  string                 `json:"avatar_uri"`
	AvgPlay    int                    `json:"avg_play"`
	City       string                 `json:"city"`
	CoreUserId string                 `json:"core_user_id"`
	Fields     []XtHotAuthorFields    `json:"fields"`
	FieldsMap  map[string]interface{} `json:"fields_map"`
	InitRank   int                    `json:"init_rank"`
	NickName   string                 `json:"nick_name"`
	Province   string                 `json:"province"`
	IncRank    int                    `json:"inc_rank"`
	ShortId    string                 `json:"short_id"`
	UniqueId   string                 `json:"unique_id"`
}
