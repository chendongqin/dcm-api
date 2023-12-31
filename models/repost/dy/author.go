package dy

type RepostSimpleReputation struct {
	Score         float64 `json:"score"`
	Level         int     `json:"level"`
	EncryptShopID string  `json:"encrypt_shop_id"`
	ShopName      string  `json:"shop_name"`
	ShopLogo      string  `json:"shop_logo"`
}

type DyAuthorStarScores struct {
	LiveScore  DyAuthorStarScore `json:"live_score"`
	VideoScore DyAuthorStarScore `json:"video_score"`
}

type DyAuthorStarScore struct {
	AvgLevel       XtAuthorScoreAvgLevel `json:"avg_level"`       //平均指数
	CooperateIndex float64               `json:"cooperate_index"` //合作指数
	CpIndex        float64               `json:"cp_index"`        //性价比指数
	GrowthIndex    float64               `json:"growth_index"`    //涨粉指数
	ShoppingIndex  float64               `json:"shopping_index"`  //种草指数
	SpreadIndex    float64               `json:"spread_index"`    //传播指数
	TopScore       float64               `json:"top_score"`       //总分
}

type XtAuthorScoreAvgLevel struct {
	CooperateIndex float64 `json:"cooperate_index"` //超过合作指数比例
	CpIndex        float64 `json:"cp_index"`        //超过性价比指数比例
	GrowthIndex    float64 `json:"growth_index"`    //超过涨粉指数比例
	ShoppingIndex  float64 `json:"shopping_index"`  //超过种草指数比例
	SpreadIndex    float64 `json:"spread_index"`    //超过传播指数比例
	TopScore       float64 `json:"top_score"`       //超过总分比例
}

type DyAuthorBasicChart struct {
	FollowerCount  int64 `json:"follower_count"`
	TotalFansCount int64 `json:"total_fans_count"`
	TotalFavorited int64 `json:"total_favorited"`
	CommentCount   int64 `json:"comment_count"`
	ForwardCount   int64 `json:"forward_count"`
}

type DyAuthorProductAnalysisCount struct {
	ProductNum  int     `json:"product_num"`
	Sales       float64 `json:"sales"`
	Gmv         float64 `json:"gmv"`
	RoomNum     int     `json:"room_num"`
	VideoNum    int     `json:"video_num"`
	HasShop     bool    `json:"has_shop"`
	IsRecommend bool    `json:"is_recommend"`
}

type DyAuthorProductRoom struct {
	RoomId       string  `json:"room_id"`
	Cover        string  `json:"cover"`
	CreateTime   int64   `json:"create_time"`
	Title        string  `json:"title"`
	MaxUserCount int64   `json:"max_user_count"`
	Gmv          float64 `json:"gmv"`
	Sales        float64 `json:"sales"`
}

type DyAuthorBaseCount struct {
	LiveCount    DyAuthorBaseLiveCount    `json:"live_count"`
	VideoCount   DyAuthorBaseVideoCount   `json:"video_count"`
	ProductCount DyAuthorBaseProductCount `json:"product_count"`
}

type DyAuthorBaseLiveCount struct {
	RoomCount         int64   `json:"room_count"`
	Room30Count       int64   `json:"room_30_count"`
	Predict30Sales    float64 `json:"predict_30_sales"`
	Predict30Gmv      float64 `json:"predict_30_gmv"`
	AgeDuration       int64   `json:"age_duration"`
	WeekRoomCount     int64   `json:"week_room_count"`
	AvgWeekRoomCount  int64   `json:"avg_week_room_count"`
	MonthRoomCount    int64   `json:"month_room_count"`
	AvgMonthRoomCount int64   `json:"avg_month_room_count"`
}

type DyAuthorBaseVideoCount struct {
	VideoCount       int64   `json:"video_count"`
	Video30Count     int64   `json:"video_30_count"`
	AvgDigg          int64   `json:"avg_digg"`
	Avg30Digg        int64   `json:"avg_30_digg"`
	DiggFollowerRate float64 `json:"digg_follower_rate"`
	Predict30Sales   float64 `json:"predict_30_sales"`
	Predict30Gmv     float64 `json:"predict_30_gmv"`
	AgeDuration      int64   `json:"age_duration"`
	WeekVideoCount   int64   `json:"week_video_count"`
	MonthVideoCount  int64   `json:"month_video_count"`
}

type DyAuthorBaseProductCount struct {
	ProductNum            int                             `json:"product_num"`
	Sales30Top3           []string                        `json:"sales_30_top_3"`
	ProductNum30Top3      []string                        `json:"product_num_30_top_3"`
	Sales30Top3Chart      []NameValueInt64Chart           `json:"sales_30_top_3_chart"`
	ProductNum30Top3Chart []NameValueChart                `json:"product_num_30_top_3_chart"`
	Predict30Sales        int64                           `json:"predict_30_sales"`
	Predict30Gmv          float64                         `json:"predict_30_gmv"`
	Sales30Chart          []DyAuthorBaseProductPriceChart `json:"sales_30_chart"`
}

type DyAuthorBaseProductPriceChart struct {
	Name       string `json:"name"`
	Sales      int64  `json:"sales"`
	ProductNum int    `json:"product_num"`
}

type RedAuthorRoom struct {
	AuthorId           string  `json:"author_id"`
	Avatar             string  `json:"avatar"`
	AuthorLivingRoomId string  `json:"author_living_room_id"`
	Sign               string  `json:"sign"`
	Nickname           string  `json:"nickname"`
	LivingTime         int64   `json:"living_time"`
	LiveTitle          string  `json:"live_title"`
	RoomId             string  `json:"room_id"`
	RoomStatus         int     `json:"room_status"`
	Gmv                float64 `json:"gmv"`
	Sales              float64 `json:"sales"`
	TotalUser          int64   `json:"total_user"`
	Tags               string  `json:"tags"`
	RoomCount          int     `json:"room_count"`
	Weight             int     `json:"-"`
	CreateTime         int64   `json:"create_time"`
}

type RedAuthorRoomBox struct {
	Date string          `json:"date"`
	List []RedAuthorRoom `json:"list"`
}

type DyAuthorIncome struct {
	AuthorId      string `json:"author_id"`
	Avatar        string `json:"avatar"`
	Nickname      string `json:"nickname"`
	UniqueId      string `json:"unique_id"`
	IsCollection  int    `json:"is_collection"`
	FollowerCount int64  `json:"follower_count"`
}

type DyAuthorStoreSimple struct {
	ShopId   string `json:"shop_id"`
	ShopName string `json:"shop_name"`
}

type DyAuthorRawData struct {
	Avatar           string `json:"avatar"`
	Birthday         int    `json:"birthday"`
	Commerce         int    `json:"commerce"`
	FollowerCount    int64  `json:"follower_count"`
	Gender           string `json:"gender"`
	Id               string `json:"id"`
	Nickname         string `json:"nickname"`
	SchoolName       string `json:"school_name"`
	SecUid           string `json:"sec_uid"`
	ShortId          string `json:"short_id"`
	Signature        string `json:"signature"`
	UniqueId         string `json:"unique_id"`
	VerificationType string `json:"verification_type"`
	VerifyName       string `json:"verify_name"`
}

type DyAuthorShopAnalysis struct {
	ShopId     string  `json:"shop_id"`
	Category   string  `json:"category"`
	ShopName   string  `json:"shop_name"`
	ShopIcon   string  `json:"shop_icon"`
	ProductNum int     `json:"product_num"`
	Gmv        float64 `json:"gmv"`
	Sales      int64   `json:"sales"`
	RoomNum    int     `json:"room_num"`
	IsCollect  int     `json:"is_collect"`
}

type DyAuthorBasic struct {
	FollowerCount        int64 `json:"follower_count"`
	FollowerCountBefore  int64 `json:"follower_count_before"`
	TotalFansCount       int64 `json:"total_fans_count"`
	TotalFansCountBefore int64 `json:"total_fans_count_before"`
	TotalFavorited       int64 `json:"total_favorited"`
	TotalFavoritedBefore int64 `json:"total_favorited_before"`
	CommentCount         int64 `json:"comment_count"`
	CommentCountBefore   int64 `json:"comment_count_before"`
	ForwardCount         int64 `json:"forward_count"`
	ForwardCountBefore   int64 `json:"forward_count_before"`
}
