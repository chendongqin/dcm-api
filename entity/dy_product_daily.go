package entity

var DyProductDailyMap = HbaseEntity{
	"live_author_list":  {AJson, "live_author_list"},
	"live_list":         {AJson, "live_list"},
	"aweme_author_list": {AJson, "aweme_author_list"},
	"aweme_list":        {AJson, "aweme_list"},
}

type DyProductDaily struct {
	LiveAuthorList  []DyProductDailyLiveAuthor  `json:"live_author_list"`
	LiveList        []DyProductDailyLiveRoom    `json:"live_list"`
	AwemeAuthorList []DyProductDailyAwemeAuthor `json:"aweme_author_list"`
	AwemeList       []DyProductDailyAweme       `json:"aweme_list"`
}

type DyProductDailyLiveAuthor struct {
	StartTime int64  `json:"start_time"`
	AuthorId  string `json:"author_id"`
}

type DyProductDailyLiveRoom struct {
	RoomId    string `json:"room_id"`
	StartTime int64  `json:"start_time"`
}

type DyProductDailyAwemeAuthor struct {
	AwemeCreateTime int64  `json:"aweme_create_time"`
	AuthorId        string `json:"author_id"`
}

type DyProductDailyAweme struct {
	AwemeId         string `json:"aweme_id"`
	AwemeCreateTime int    `json:"aweme_create_time"`
}
