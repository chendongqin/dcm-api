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
	AvatarUri  string                   `json:"avatar_uri"`
	AvgPlay    int                      `json:"avg_play"`
	City       string                   `json:"city"`
	CoreUserId string                   `json:"core_user_id"`
	Fields     []XtHotAwemeAuthorFields `json:"fields"`
	IncRank    int                      `json:"inc_rank"`
	InitRank   int                      `json:"init_rank"`
	NewRank    int                      `json:"new_rank"`
	NickName   string                   `json:"nick_name"`
	OldRank    int                      `json:"old_rank"`
	Province   string                   `json:"province"`
}

type XtHotAwemeAuthorFields struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
