package entity

var DyAuthorMap = HbaseEntity{
	"id":         {String, "author_id"},
	"collection": {Int, "collection"},
	"crawl_time": {Long, "crawl_time"},
	"data":       {Json, "data"},
}

type DyAuthor struct {
	AuthorID   string       `json:"author_id"`
	Collection int          `json:"collection"`
	CrawlTime  int64        `json:"crawl_time"`
	Data       DyAuthorData `json:"data"`
}

type DyAuthorData struct {
	Avatar                string `json:"avatar"`
	AwemeCount            int64  `json:"aweme_count"`
	Birthday              string `json:"birthday"`
	Age                   int    `json:"age"`
	City                  string `json:"city"`
	Commerce              int64  `json:"commerce"`
	Country               string `json:"country"`
	DongtaiCount          int64  `json:"dongtai_count"`
	FollowerCount         int64  `json:"follower_count"`
	FollowingCount        int64  `json:"following_count"`
	Gender                int64  `json:"gender"`
	ID                    string `json:"id"`
	IsStar                int64  `json:"is_star"`
	LiveCommerce          bool   `json:"live_commerce"`
	Nickname              string `json:"nickname"`
	PcFans                int64  `json:"pc_fans"`
	Province              string `json:"province"`
	RoomID                string `json:"room_id"`
	ScheduledTime         string `json:"scheduled_time"`
	SchoolName            string `json:"school_name"`
	SecUID                string `json:"sec_uid"`
	ShortID               string `json:"short_id"`
	ShowFollowingFollower int64  `json:"show_following_follower"`
	Signature             string `json:"signature"`
	TotalFavorited        int64  `json:"total_favorited"`
	UniqueID              string `json:"unique_id"`
	VerificationType      int64  `json:"verification_type"`
	VerifyName            string `json:"verify_name"`
	ShareUrl              string `json:"share_url"`
}

//type DyAuthorFans struct {
//Douyin  DyAuthorFansCount `json:"douyin"`
//Huoshan DyAuthorFansCount `json:"huoshan"`
//Toutiao DyAuthorFansCount `json:"toutiao"`
//}

type DyAuthorFansCount struct {
	Count int64 `json:"count"`
}
