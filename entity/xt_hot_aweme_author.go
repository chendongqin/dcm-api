package entity

var XtHotAwemeAuthorMap = HbaseEntity{
	"crawl_time":  {Long, "crawl_time"},
	"update_time": {Long, "update_time"},
	"data":        {AJson, "data"},
	"tab":         {String, "tab"},
	"category":    {String, "category"},
}

type XtHotAwemeAuthor struct {
	CrawlTime  int64                  `json:"crawl_time"`
	UpdateTime int64                  `json:"update_time"`
	Data       []XtHotAwemeAuthorData `json:"data"`
	Tab        string                 `json:"tab"`
	Category   string                 `json:"category"`
}

type XtHotAwemeAuthorData struct {
	AvatarUri  string              `json:"avatar_uri"`
	City       string              `json:"city"`
	CoreUserId string              `json:"core_user_id"`
	Fields     []XtHotAuthorFields `json:"fields"`
	IncRank    int                 `json:"inc_rank"`
	InitRank   int                 `json:"init_rank"`
	NickName   string              `json:"nick_name"`
	Province   string              `json:"province"`
}

type XtHotAuthorFields struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
