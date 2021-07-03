package entity

var DyLiveInfoMap = HbaseEntity{
	"add_time":               {Long, "add_time"},
	"challenge":              {Json, "challenge"},
	"cover":                  {String, "cover"},
	"crawl_time":             {Long, "crawl_time"},
	"create_time":            {Long, "create_time"},
	"discover_time":          {Long, "discover_time"},
	"fans_club_count":        {Long, "fans_club_count"},
	"finish_time":            {Long, "finish_time"},
	"follow_count":           {Long, "follow_count"},
	"gift_uv_count":          {Long, "gift_uv_count"},
	"hour_rank":              {Json, "hour_rank"},
	"hour_sales_rank":        {Json, "hour_sales_rank"},
	"level":                  {Int, "level"},
	"like_count":             {Long, "like_count"},
	"play_flv_url":           {String, "play_flv_url"},
	"play_url":               {String, "play_url"},
	"pmt_cnt":                {Long, "pmt_cnt"},
	"room_id":                {String, "room_id"},
	"room_status":            {Int, "room_status"},
	"room_ticket_count":      {Long, "room_ticket_count"},
	"tag":                    {String, "tag"},
	"title":                  {String, "title"},
	"top_fans":               {AJson, "top_fans"},
	"total_user":             {Long, "total_user"},
	"user":                   {Json, "user"},
	"user_count":             {Long, "user_count"},
	"user_count_composition": {Json, "user_count_composition"},
	"watch_cnt":              {Long, "watch_cnt"},
}

type DyLiveInfo struct {
	AddTime      int64               `json:"add_time"`
	Challenge    DyLiveInfoChallenge `json:"challenge"`
	Cover        string              `json:"cover"`
	CrawlTime    int64               `json:"crawl_time"`
	CreateTime   int64               `json:"create_time"`
	DiscoverTime int64               `json:"discover_time"`
	//DiscoverTimeOriginal int64                      `json:"discover_time_original"`
	//CreateTimeFixed      int64                      `json:"create_time_fixed"`
	FansClubCount        int64                      `json:"fans_club_count"`
	FinishTime           int64                      `json:"finish_time"`
	FollowCount          int64                      `json:"follow_count"`
	GiftUvCount          int64                      `json:"gift_uv_count"`
	HourRank             DyLiveInfoRank             `json:"hour_rank"`
	HourSalesRank        DyLiveInfoRank             `json:"hour_sales_rank"`
	Level                int                        `json:"level"`
	LikeCount            int64                      `json:"like_count"`
	PlayFlvURL           string                     `json:"play_flv_url"`
	PlayURL              string                     `json:"play_url"`
	PmtCnt               int64                      `json:"pmt_cnt"`
	RoomID               string                     `json:"room_id"`
	RoomStatus           string                     `json:"room_status"`
	RoomTicketCount      int                        `json:"room_ticket_count"`
	Tag                  string                     `json:"tag"`
	Title                string                     `json:"title"`
	TopFans              []DyLiveInfoTopFan         `json:"top_fans"`
	TotalUser            int64                      `json:"total_user"`
	User                 DyLiveInfoUser             `json:"user"`
	UserCount            int64                      `json:"user_count"`
	UserCountComposition DyLiveUserCountComposition `json:"user_count_composition"`
	WatchCnt             int64                      `json:"watch_cnt"`
}

type DyLiveInfoChallenge struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserCount int64  `json:"user_count"`
	ViewCount int64  `json:"view_count"`
}

type DyLiveInfoRank struct {
	CrawlTime    int64  `json:"crawl_time"`
	DiscoverTime int64  `json:"discover_time"`
	Rank         int    `json:"rank"`
	RoomID       string `json:"room_id"`
	Type         int    `json:"type"`
}

type DyLiveInfoTopFan struct {
	Avatar    string `json:"avatar"`
	FanTicket int64  `json:"fan_ticket"`
	ID        string `json:"id"`
	Nickname  string `json:"nickname"`
}

type DyLiveInfoUser struct {
	Avatar        string `json:"avatar"`
	FollowerCount int64  `json:"follower_count"`
	ID            string `json:"id"`
	Nickname      string `json:"nickname"`
	PayScore      int64  `json:"pay_score"`
	TicketCount   int64  `json:"ticket_count"`
	WithCommerce  bool   `json:"with_commerce"`
}

type DyLiveUserCountComposition struct {
	City        float64 `json:"city"`
	MyFollow    float64 `json:"my_follow"`
	Other       float64 `json:"other"`
	VideoDetail float64 `json:"video_detail"`
}
