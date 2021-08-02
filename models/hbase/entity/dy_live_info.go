package entity

var DyLiveInfoMap = HbaseEntity{
	"add_time":                      {Long, "add_time"},
	"challenge":                     {Json, "challenge"},
	"cover":                         {String, "cover"},
	"crawl_time":                    {Long, "crawl_time"},
	"create_time":                   {Long, "create_time"},
	"discover_time":                 {Long, "discover_time"},
	"fans_club_count":               {Long, "fans_club_count"},
	"finish_time":                   {Long, "finish_time"},
	"follow_count":                  {Long, "follow_count"},
	"gift_uv_count":                 {Long, "gift_uv_count"},
	"hour_rank":                     {Json, "hour_rank"},
	"hour_sales_rank":               {Json, "hour_sales_rank"},
	"level":                         {Int, "level"},
	"like_count":                    {Long, "like_count"},
	"play_flv_url":                  {String, "play_flv_url"},
	"play_url":                      {String, "play_url"},
	"pmt_cnt":                       {Long, "pmt_cnt"},
	"room_id":                       {String, "room_id"},
	"room_status":                   {Int, "room_status"},
	"room_ticket_count":             {Long, "room_ticket_count"},
	"tag":                           {String, "tag"},
	"title":                         {String, "title"},
	"top_fans":                      {AJson, "top_fans"},
	"total_user":                    {Long, "total_user"},
	"user":                          {Json, "user"},
	"user_count":                    {Long, "user_count"},
	"user_count_composition":        {Json, "user_count_composition"},
	"watch_cnt":                     {Long, "watch_cnt"},
	"trends_crawl_time":             {Long, "trends_crawl_time"},
	"trends_online_trends":          {AJson, "online_trends"},
	"other_barrage_count":           {Long, "barrage_count"},
	"trends_follower_count_trends":  {AJson, "follower_count_trends"},
	"trends_fans_club_count_trends": {AJson, "fans_club_count_trends"},
}

type DyLiveInfo struct {
	AddTime              int64                      `json:"add_time"`  //添加到直播库时间
	Challenge            DyLiveInfoChallenge        `json:"challenge"` //话题
	Cover                string                     `json:"cover"`     //封面
	CrawlTime            int64                      `json:"crawl_time"`
	CreateTime           int64                      `json:"create_time"`     //开播时间
	DiscoverTime         int64                      `json:"discover_time"`   //发现时间
	FansClubCount        int64                      `json:"fans_club_count"` //粉丝团数目
	FinishTime           int64                      `json:"finish_time"`     //结束时间
	FollowCount          int64                      `json:"follow_count"`    //粉丝数
	GiftUvCount          int64                      `json:"gift_uv_count"`   //送礼人数
	HourRank             DyLiveInfoRank             `json:"hour_rank"`       //小时榜排行
	HourSalesRank        DyLiveInfoRank             `json:"hour_sales_rank"` //带货小时榜
	Level                int                        `json:"level"`           //是否抓取弹幕 1:是 0:否
	LikeCount            int64                      `json:"like_count"`      //点赞数
	PlayFlvURL           string                     `json:"play_flv_url"`    //视频流地址
	PlayURL              string                     `json:"play_url"`        //直播地址
	PmtCnt               int64                      `json:"pmt_cnt"`
	RoomID               string                     `json:"room_id"`
	RoomStatus           int                        `json:"room_status"`       //直播状态 2:在播 4:下播
	RoomTicketCount      int                        `json:"room_ticket_count"` //该场直播音浪数
	Title                string                     `json:"title"`
	TopFans              []DyLiveInfoTopFan         `json:"top_fans"`   //送礼top3
	TotalUser            int64                      `json:"total_user"` //总pv
	User                 DyLiveInfoUser             `json:"user"`
	UserCount            int64                      `json:"user_count"`             //当前在线人数
	UserCountComposition DyLiveUserCountComposition `json:"user_count_composition"` //用户来源
	WatchCnt             int64                      `json:"watch_cnt"`              //总pv
	TrendsCrawlTime      int64                      `json:"trends_crawl_time"`      //更新时间
	BarrageCount         int64                      `json:"barrage_count"`          //弹幕人数
	OnlineTrends         []DyLiveOnlineTrends       `json:"online_trends"`
	FollowerCountTrends  []LiveFollowerCountTrends  `json:"follower_count_trends"`
	FansClubCountTrends  []LiveAnsClubCountTrends   `json:"fans_club_count_trends"`
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

type DyLiveOnlineTrends struct {
	CrawlTime     int64 `json:"crawl_time"`
	WatchCnt      int64 `json:"watch_cnt"`
	UserCount     int64 `json:"user_count"`
	FollowerCount int64 `json:"follower_count"`
}

type DyLiveIncOnlineTrends struct {
	UserCount int64 `json:"user_count"`
	WatchInc  int64 `json:"watch_inc"`
}

type DyLiveIncOnlineTrendsChart struct {
	Date            []string                `json:"date"`
	IncOnlineTrends []DyLiveIncOnlineTrends `json:"inc_online_trends"`
}

type LiveFollowerCountTrends struct {
	CrawlTime      int64 `json:"crawl_time"`
	FollowerCount  int64 `json:"follower_count"`
	NewFollowCount int64 `json:"new_follow_count"`
}

type LiveAnsClubCountTrends struct {
	FansClubCount     int64 `json:"fans_club_count"`
	TodayNewFansCount int64 `json:"today_new_fans_count"`
	CrawlTime         int64 `json:"crawl_time"`
}
