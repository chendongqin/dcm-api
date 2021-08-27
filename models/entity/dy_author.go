package entity

var DyAuthorMap = HbaseEntity{
	"id":                       {String, "author_id"},
	"collection":               {Int, "collection"},
	"crawl_time":               {Long, "crawl_time"},
	"data":                     {Json, "data"},
	"tags":                     {String, "tags"},
	"tags_level_two":           {String, "tags_level_two"},
	"other_aweme_list":         {AJson, "aweme_list"},
	"other_room_list":          {AJson, "room_list"},
	"other_room_count":         {Int, "room_count"},
	"other_live_duration":      {String, "live_duration"},
	"other_avg_live_duration":  {Long, "avg_live_duration"},
	"other_med_watch_cnt":      {Long, "med_watch_cnt"},
	"other_interaction_rate":   {Double, "interaction_rate"},
	"other_predict_30_gmv":     {Double, "predict_30_gmv"},
	"other_real_30_gmv":        {Double, "real_30_gmv"},
	"other_predict_30_sales":   {Double, "predict_30_sales"},
	"other_real_30_sales":      {Double, "real_30_sales"},
	"other_aweme_count":        {Int, "aweme_count"},
	"other_digg_count":         {Long, "digg_count"},
	"other_digg_follower_rate": {Double, "digg_follower_rate"},
	"other_duration":           {Long, "duration"},
	"other_med_digg":           {Long, "med_digg"},
	"other_first_live_time":    {Long, "first_live_time"},
	"other_first_aweme_time":   {Long, "first_aweme_time"},
	"other_product_count":      {Int, "product_count"},
	"other_first_product_time": {Long, "first_product_time"},
	"other_total_fans_count":   {Long, "total_fans_count"},
	"other_follower_count":     {Long, "follower_count"},
	"other_total_favorited":    {Long, "total_favorited"},
	"other_comment_count":      {Long, "comment_count"},
	"other_forward_count":      {Long, "forward_count"},
}

type DyAuthor struct {
	AuthorID         string          `json:"author_id"`
	Collection       int             `json:"collection"`
	CrawlTime        int64           `json:"crawl_time"`
	Data             DyAuthorData    `json:"data"`
	Tags             string          `json:"tags"`
	TagsLevelTwo     string          `json:"tags_level_two"`
	AwemeList        []DyAuthorAweme `json:"aweme_list"`
	RoomList         []DyAuthorRoom  `json:"room_list"`
	RoomCount        int             `json:"room_count"`
	RoomId           string          `json:"room_id"`
	RoomStatus       int             `json:"room_status"`
	LiveDuration     string          `json:"live_duration"`
	AgeLiveDuration  int64           `json:"avg_live_duration"`
	MedWatchCnt      int64           `json:"med_watch_cnt"`
	InteractionRate  float64         `json:"interaction_rate"`
	Predict30Gmv     float64         `json:"predict_30_gmv"`
	Real30Gmv        float64         `json:"real_30_gmv"`
	Predict30Sales   float64         `json:"predict_30_sales"`
	Real30Sales      float64         `json:"real_30_sales"`
	AwemeCount       int             `json:"aweme_count"`
	DiggCount        int64           `json:"digg_count"`
	DiggFollowerRate float64         `json:"digg_follower_rate"`
	Duration         int64           `json:"duration"`
	MedDigg          int64           `json:"med_digg"`
	ProductCount     int             `json:"product_count"`
	FirstLiveTime    int64           `json:"first_live_time"`
	FirstAwemeTime   int64           `json:"first_aweme_time"`
	FirstProductTime int64           `json:"first_product_time"`
	TotalFansCount   int64           `json:"total_fans_count"`
	FollowerCount    int64           `json:"follower_count"`
	TotalFavorited   int64           `json:"total_favorited"`
	CommentCount     int64           `json:"comment_count"`
	ForwardCount     int64           `json:"forward_count"`
}

type DyAuthorRoom struct {
	RoomId     string `json:"room_id"`
	CreateTime int64  `json:"create_time"`
	FinishTime int64  `json:"finish_time"`
}

type DyAuthorAweme struct {
	AwemeId          string  `json:"aweme_id"`
	CreateTime       int64   `json:"create_time"`
	DiggCount        int64   `json:"digg_count"`
	DiggFollowerRate float64 `json:"digg_follower_rate"`
	Duration         int64   `json:"duration"`
	FollowerCount    int64   `json:"follower_count"`
}

type DyAuthorData struct {
	Avatar                string           `json:"avatar"`
	AwemeCount            int64            `json:"aweme_count"`
	Birthday              string           `json:"birthday"`
	Age                   int              `json:"age"`
	City                  string           `json:"city"`
	Commerce              int64            `json:"commerce"`
	Country               string           `json:"country"`
	DongtaiCount          int64            `json:"dongtai_count"`
	FollowerCount         int64            `json:"follower_count"`
	FollowingCount        int64            `json:"following_count"`
	Fans                  DyAuthorInfoFans `json:"fans"`
	Gender                int64            `json:"gender"`
	ID                    string           `json:"id"`
	IsStar                int64            `json:"is_star"`
	LiveCommerce          bool             `json:"live_commerce"`
	Nickname              string           `json:"nickname"`
	PcFans                int64            `json:"pc_fans"`
	Province              string           `json:"province"`
	RoomID                string           `json:"room_id"`
	ScheduledTime         string           `json:"scheduled_time"`
	SchoolName            string           `json:"school_name"`
	SecUID                string           `json:"sec_uid"`
	ShortID               string           `json:"short_id"`
	ShowFollowingFollower int64            `json:"show_following_follower"`
	Signature             string           `json:"signature"`
	TotalFavorited        int64            `json:"total_favorited"`
	UniqueID              string           `json:"unique_id"`
	VerificationType      int64            `json:"verification_type"`
	VerifyName            string           `json:"verify_name"`
	ShareUrl              string           `json:"share_url"`
	CrawlTime             int64            `json:"crawl_time"`
}

type DyAuthorInfoFans struct {
	Douyin  DyAuthorFansCount `json:"douyin"`
	Huoshan DyAuthorFansCount `json:"huoshan"`
	Toutiao DyAuthorFansCount `json:"toutiao"`
}

type DyAuthorFansCount struct {
	Count int64 `json:"count"`
}
