package es

type DyAuthor struct {
	AuthorId           string  `json:"author_id"`
	Avatar             string  `json:"avatar"`
	FollowerCount      int64   `json:"follower_count"`
	FollowerIncreCount int64   `json:"follower_incre_count"`
	Gender             string  `json:"gender"`
	Birthday           int     `json:"birthday"`
	Province           string  `json:"province"`
	City               string  `json:"city"`
	VerificationType   string  `json:"verification_type"` //0未认证；1蓝v；2黄v
	VerifyName         string  `json:"verify_name"`
	ShortId            string  `json:"short_id"`
	CreateTime         string  `json:"create_time"`
	RowTime            string  `json:"row_time"`
	RankSellTags       string  `json:"rank_sell_tags"`
	Exist              int     `json:"exist"`
	Nickname           string  `json:"nickname"`
	Tags               string  `json:"tags"`
	TagsLevelTwo       string  `json:"tags_level_two"`
	DiggFollowerRate   float64 `json:"digg_follower_rate"`
	MedDigg            int64   `json:"med_digg"`
	MedWatchCnt        int64   `json:"med_watch_cnt"`
	InteractionRate    float64 `json:"interaction_rate"`
	Predict30Gmv       float64 `json:"predict_30_gmv"`
	Real30Gmv          float64 `json:"real_30_gmv"`
	Score              float64 `json:"score"`
	Level              int     `json:"level"`
	Brand              int     `json:"brand"`
	IsDelivery         int     `json:"is_delivery"`
	UniqueId           string  `json:"unique_id"`
	RoomId             string  `json:"room_id"`
}
