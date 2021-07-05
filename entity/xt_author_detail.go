package entity

var XtAuthorDetailMap = HbaseEntity{
	"author_show_items":     {Json, "author_show_items"},
	"base_info":             {Json, "base_info"},
	"crawl_time":            {Long, "crawl_time"},
	"distributions":         {AJson, "distributions"},
	"live_focus_data":       {Json, "live_focus_data"},
	"live_ltm_description":  {Json, "live_ltm_description"},
	"live_ltm_item_statics": {Json, "live_ltm_item_statics"},
	"live_mark_info":        {Json, "live_mark_info"},
	"live_score":            {Json, "live_score"},
	"ltm_description":       {Json, "ltm_description"},
	"ltm_item_statics":      {AJson, "ltm_item_statics"},
	"recommend_author":      {AJson, "recommend_author"},
	"score":                 {Json, "score"},
	"star_id":               {Long, "star_id"},
	"uid":                   {String, "uid"},
	"video_focus_data":      {Json, "video_focus_data"},
	"video_mark_info":       {Json, "video_mark_info"},
}

type XtAuthorDetail struct {
	AuthorShowItems    XtAuthorShowItems            `json:"author_show_items"`
	BaseInfo           XtAuthorBaseInfo             `json:"base_info"`
	CrawlTime          int64                        `json:"crawl_time"`
	Distributions      []XtAuthorDistributions      `json:"distributions"`
	LiveFocusData      XtAuthorLiveFocusData        `json:"live_focus_data"`
	LiveLtmDescription XtAuthorLtmDescription       `json:"live_ltm_description"`
	LiveLtmItemStatics []XtAuthorLiveLtmItemStatics `json:"live_ltm_item_statics"`
	LiveMarkInfo       XtAuthorLiveMarkInfo         `json:"live_mark_info"`
	LiveScore          XtAuthorLiveScore            `json:"live_score"`
	LtmDescription     XtAuthorLtmDescription       `json:"ltm_description"`
	LtmItemStatics     []XtAuthorLtmItemStatics     `json:"ltm_item_statics"`
	RecommendAuthor    []XtRecommendAuthor          `json:"recommend_author"`
	Score              XtAuthorScore                `json:"score"`
	StarID             int64                        `json:"star_id"`
	UID                string                       `json:"uid"`
	VideoFocusData     XtAuthorVideoFocus           `json:"video_focus_data"`
	VideoMarkInfo      XtAuthorLiveMarkInfo         `json:"video_mark_info"`
}

type XtAuthorInteraction struct {
	CompareAuthor float64 `json:"compare_author"`
	CompareAvg    float64 `json:"compare_avg"`
	Rate          float64 `json:"rate"`
}

type XtAuthorDescription struct {
	Interaction             XtAuthorInteraction `json:"interaction"` //个人作品
	InteractionEnrollment   XtAuthorInteraction `json:"interaction_enrollment"`
	PlayMedium              XtAuthorInteraction `json:"play_medium"`
	PlayMediumEnrollment    XtAuthorInteraction `json:"play_medium_enrollment"`
	VideoViewRate           XtAuthorInteraction `json:"video_view_rate"`
	VideoViewRateEnrollment XtAuthorInteraction `json:"video_view_rate_enrollment"`
}

type XtAuthorItemInfo struct {
	Comment           int64  `json:"comment"`
	CoreUserID        string `json:"core_user_id"`
	CreateTime        int64  `json:"create_time"`
	CreateTimestamp   int64  `json:"create_timestamp"`
	Duration          int    `json:"duration"`
	DurationMin       int    `json:"duration_min"`
	HeadImageURI      string `json:"head_image_uri"`
	ItemAnimatedCover string `json:"item_animated_cover"`
	ItemCover         string `json:"item_cover"`
	ItemDate          string `json:"item_date"`
	ItemID            string `json:"item_id"`
	ItemTitle         string `json:"item_title"`
	Like              int64  `json:"like"`
	OriginalStatus    int    `json:"original_status"`
	Play              int64  `json:"play"`
	Share             int64  `json:"share"`
	Status            int    `json:"status"`
	Title             string `json:"title"`
	URL               string `json:"url"`
	VideoID           string `json:"video_id"`
}

type XtAuthorItemStaticsAntispam struct {
	AvgComment int64 `json:"avg_comment"`
	AvgLike    int64 `json:"avg_like"`
	AvgPlay    int64 `json:"avg_play"`
	AvgShare   int64 `json:"avg_share"`
}

type XtAuthorShowItems struct {
	DataDescription           XtAuthorDescription         `json:"data_description"`
	LatestItemInfo            []XtAuthorItemInfo          `json:"latest_item_info"`
	LatestItemStaticsAntispam XtAuthorItemStaticsAntispam `json:"latest_item_statics_antispam"`
	LatestStarItemInfo        []XtAuthorItemInfo          `json:"latest_star_item_info"`
}

type XtAuthorBaseInfo struct {
	AvgPlay         int64               `json:"avg_play"`
	AwemeTags       []interface{}       `json:"aweme_tags"`
	CategoryID      string              `json:"category_id"`
	CreateTime      string              `json:"create_time"`
	ECommerceEnable bool                `json:"e_commerce_enable"`
	ID              string              `json:"id"`
	IsStar          bool                `json:"is_star"`
	LowestPrice     float64             `json:"lowest_price"`
	McnIntroduction string              `json:"mcn_introduction"`
	McnLogo         string              `json:"mcn_logo"`
	McnName         string              `json:"mcn_name"`
	ModifyTime      string              `json:"modify_time"`
	NickName        string              `json:"nick_name"`
	Platform        []int               `json:"platform"`
	PlatformChannel []int               `json:"platform_channel"`
	PlatformSource  int                 `json:"platform_source"`
	StarTags        []interface{}       `json:"star_tags"`
	Tags            []string            `json:"tags"`
	TagsAuthorStyle []string            `json:"tags_author_style"`
	TagsIds         []int               `json:"tags_ids"`
	TagsIdsLevelTwo []int               `json:"tags_ids_level_two"`
	TagsLevelTwo    []string            `json:"tags_level_two"`
	TagsRelation    map[string][]string `json:"tags_relation"`
}

type XtDistributionsList struct {
	DistributionKey   string `json:"distribution_key"`
	DistributionValue int    `json:"distribution_value"`
}

type XtAuthorDistributions struct {
	Description      string                `json:"description"`
	DistributionList []XtDistributionsList `json:"distribution_list"`
	Image            []string              `json:"image"`
	OriginType       int                   `json:"origin_type"`
	Type             int                   `json:"type"`
	TypeDisplay      string                `json:"type_display"`
}

type XtAuthorLiveFocusData struct {
	AvgDurationMin       int     `json:"avg_duration_min"`
	AvgPlay              int64   `json:"avg_play"`
	ClubFollower         int64   `json:"club_follower"`
	ExpectedCpm          float64 `json:"expected_cpm"`
	ExpectedMaxWatch     int64   `json:"expected_max_watch"`
	ExpectedPlayNum      int64   `json:"expected_play_num"`
	MiddlePlay           int64   `json:"middle_play"`
	OrderAvgTimeCost     int64   `json:"order_avg_time_cost"`
	OrderCompleteCnt     int64   `json:"order_complete_cnt"`
	OrderCompleteRate    float64 `json:"order_complete_rate"`
	PersonalInterateRate float64 `json:"personal_interate_rate"`
	TotalFavourCnt       string  `json:"total_favour_cnt"`
}

type XtAuthorLtmDescription struct {
	CommentRange  int64 `json:"comment_range"`
	LikeRange     int64 `json:"like_range"`
	MaxComment    int64 `json:"max_comment"`
	MaxLike       int64 `json:"max_like"`
	MaxPlay       int64 `json:"max_play"`
	MaxShare      int64 `json:"max_share"`
	MediumComment int64 `json:"medium_comment"`
	MediumLike    int64 `json:"medium_like"`
	MediumPlay    int64 `json:"medium_play"`
	MediumShare   int64 `json:"medium_share"`
	MinComment    int64 `json:"min_comment"`
	MinLike       int64 `json:"min_like"`
	MinPlay       int64 `json:"min_play"`
	MinShare      int64 `json:"min_share"`
	ShareRange    int64 `json:"share_range"`
	VvRange       int64 `json:"vv_range"`
}

type XtAuthorLiveLtmItemStatics struct {
	AvgWatch      int64   `json:"avg_watch"`
	Barrage       int64   `json:"barrage"`
	Comment       int64   `json:"comment"`
	CreateTime    int64   `json:"create_time"`
	Duration      int     `json:"duration"`
	DurationMin   int     `json:"duration_min"`
	FansClub      int     `json:"fans_club"`
	FinishTime    int64   `json:"finish_time"`
	HeadImageURI  string  `json:"head_image_uri"`
	InteractRate  float64 `json:"interact_rate"`
	ItemID        string  `json:"item_id"`
	Like          int64   `json:"like"`
	MaxWatch      int64   `json:"max_watch"`
	Play          int64   `json:"play"`
	PlayOver5Min  int64   `json:"play_over_5min"`
	SendGift      int64   `json:"send_gift"`
	SendProp      int64   `json:"send_prop"`
	SendRedPacket int64   `json:"send_red_packet"`
	Share         int64   `json:"share"`
	Status        int     `json:"status"`
	Title         string  `json:"title"`
	URL           string  `json:"url"`
	VideoID       string  `json:"video_id"`
}

type XtAuthorHotListRanks struct {
	HotListID   string `json:"hot_list_id"`
	HotListName string `json:"hot_list_name"`
	Rank        int    `json:"rank"`
	Tag         string `json:"tag"`
	Type        int    `json:"type"`
}

type XtAuthorPriceInfo struct {
	ActivityInfo       []interface{} `json:"activity_info"`
	AuthorID           string        `json:"author_id"`
	Desc               string        `json:"desc"`
	Enable             bool          `json:"enable"`
	Field              string        `json:"field"`
	HasDiscount        bool          `json:"has_discount"`
	IsOpen             bool          `json:"is_open"`
	NeedPrice          bool          `json:"need_price"`
	OnlineStatus       int           `json:"online_status"`
	OriginPrice        int           `json:"origin_price"`
	PlatformActivityID string        `json:"platform_activity_id"`
	PlatformSource     int           `json:"platform_source"`
	Price              int           `json:"price"`
	SettlementDesc     string        `json:"settlement_desc"`
	SettlementType     int           `json:"settlement_type"`
	TaskCategory       int           `json:"task_category"`
	TaskCategoryStatus int           `json:"task_category_status"`
	VideoType          int           `json:"video_type"`
}

type XtAuthorLiveMarkInfo struct {
	ActivityInfo []interface{}          `json:"activity_info"`
	HotListRanks []XtAuthorHotListRanks `json:"hot_list_ranks"`
	PriceInfo    []XtAuthorPriceInfo    `json:"price_info"`
}

type XtAuthorAvgLevel struct {
	CooperateIndex int `json:"cooperate_index"` //超过合作指数比例
	CpIndex        int `json:"cp_index"`        //超过性价比指数比例
	GrowthIndex    int `json:"growth_index"`    //超过涨粉指数比例
	ShoppingIndex  int `json:"shopping_index"`  //超过种草指数比例
	SpreadIndex    int `json:"spread_index"`    //超过传播指数比例
	TopScore       int `json:"top_score"`       //超过总分比例
}

type XtAuthorLiveScore struct {
	AvgLevel       XtAuthorAvgLevel `json:"avg_level"`
	BusinessIndex  int              `json:"business_index"`
	CooperateIndex int              `json:"cooperate_index"` //合作指数
	CpIndex        int              `json:"cp_index"`        //性价比指数
	GrowthIndex    int              `json:"growth_index"`    //涨粉指数
	ImpactIndex    int              `json:"impact_index"`
	Median         XtAuthorAvgLevel `json:"median"`         //行业中位数
	ShoppingIndex  int              `json:"shopping_index"` //种草指数
	SpreadIndex    int              `json:"spread_index"`   //传播指数
	TopScore       int              `json:"top_score"`      //总分
}

type XtAuthorLtmItemStatics struct {
	Comment        int64  `json:"comment"`
	CoreUserID     string `json:"core_user_id"`
	CreateTime     int64  `json:"create_time"`
	Duration       int    `json:"duration"`
	HeadImageURI   string `json:"head_image_uri"`
	ItemID         string `json:"item_id"`
	Like           int64  `json:"like"`
	OriginalStatus int    `json:"original_status"`
	Play           int64  `json:"play"`
	Share          int64  `json:"share"`
	Status         int    `json:"status"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	VideoID        string `json:"video_id"`
}

type XtRecommendAuthor struct {
	AuthorID             string   `json:"author_id"`
	AuthorScore          int64    `json:"author_score"`
	AvatarURI            string   `json:"avatar_uri"`
	AvgPlay              int64    `json:"avg_play"`
	CoreUserID           string   `json:"core_user_id"`
	ECommerceEnable      bool     `json:"e_commerce_enable"`
	ExpectedCpm          float64  `json:"expected_cpm"`
	ExpectedPlayNum      int64    `json:"expected_play_num"`
	Follower             int64    `json:"follower"`
	Grade                int      `json:"grade"`
	IsCollected          bool     `json:"is_collected"`
	IsOnline             bool     `json:"is_online"`
	IsPlanAuthor         bool     `json:"is_plan_author"`
	MinPrice             int      `json:"min_price"`
	NickName             string   `json:"nick_name"`
	OngoingOrderCnt      int      `json:"ongoing_order_cnt"`
	Overload             bool     `json:"overload"`
	PersonalInterateRate float64  `json:"personal_interate_rate"`
	RecommendTypes       []string `json:"recommend_types"`
	TagsAuthorStyle      string   `json:"tags_author_style"`
	TagsContent          []string `json:"tags_content"`
}

type XtAuthorHotList struct {
	HotListID   string `json:"hot_list_id"`
	HotListName string `json:"hot_list_name"`
	Rank        int    `json:"rank"`
	Tag         string `json:"tag"`
}

type XtAuthorScore struct {
	AvgLevel       XtAuthorAvgLevel `json:"avg_level"`
	BusinessIndex  int              `json:"business_index"`
	CooperateIndex int              `json:"cooperate_index"`
	CpIndex        int              `json:"cp_index"`
	GrowthHotlist  XtAuthorHotList  `json:"growth_hotlist"`
	GrowthIndex    int              `json:"growth_index"`
	ImpactIndex    int              `json:"impact_index"`
	Median         XtAuthorAvgLevel `json:"median"`
	ShoppingIndex  int              `json:"shopping_index"`
	SpreadHotlist  XtAuthorHotList  `json:"spread_hotlist"`
	SpreadIndex    int              `json:"spread_index"`
	TopHotlist     XtAuthorHotList  `json:"top_hotlist"`
	TopScore       int              `json:"top_score"`
}

type XtAuthorVideoFocus struct {
	AvgPlay              int64   `json:"avg_play"`
	ExpectedCpm          float64 `json:"expected_cpm"`
	ExpectedPlayNum      int64   `json:"expected_play_num"`
	MiddlePlay           int64   `json:"middle_play"`
	OrderAvgTimeCost     int64   `json:"order_avg_time_cost"`
	OrderCompleteCnt     int64   `json:"order_complete_cnt"`
	OrderCompleteRate    float64 `json:"order_complete_rate"`
	PersonalInterateRate float64 `json:"personal_interate_rate"`
	TotalFavourCnt       string  `json:"total_favour_cnt"`
}
