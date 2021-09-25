package dy

type TakeGoodsRankRet struct {
	Rank             int                      `json:"rank,omitempty"`
	Nickname         string                   `json:"nickname,omitempty"`
	AuthorCover      string                   `json:"author_cover,omitempty"`
	SumGmv           float64                  `json:"sum_gmv,omitempty"`
	SumSales         float64                  `json:"sum_sales,omitempty"`
	AvgPrice         float64                  `json:"avg_price,omitempty"`
	AuthorId         string                   `json:"author_id,omitempty"`
	UniqueId         string                   `json:"unique_id,omitempty"`
	Tags             string                   `json:"tags,omitempty"`
	VerificationType int                      `json:"verification_type,omitempty"`
	VerifyName       string                   `json:"verify_name,omitempty"`
	RoomCount        int                      `json:"room_count,omitempty"`
	RoomList         []map[string]interface{} `json:"room_list"`
}

type AuthorFansRankRet struct {
	AuthorCover           string `json:"author_cover"`
	AuthorId              string `json:"author_id"`
	AwemeIncFollowerCount int    `json:"aweme_inc_follower_count"`
	City                  string `json:"city"`
	Country               string `json:"country"`
	DateTime              string `json:"date_time"`
	FollowerCount         int    `json:"follower_count"`
	Gender                int64  `json:"gender"`
	Id                    string `json:"id"`
	IncFollowerCount      int    `json:"inc_follower_count"`
	IsDelivery            int64  `json:"is_delivery"`
	LiveIncFollowerCount  int    `json:"live_inc_follower_count"`
	Nickname              string `json:"nickname"`
	Province              string `json:"province"`
	Rank                  int    `json:"rank"`
	ShortId               string `json:"short_id"`
	Tags                  string `json:"tags"`
	TagsLevelTwo          string `json:"tags_level_two"`
	UniqueId              string `json:"unique_id"`
	VerificationType      int    `json:"verification_type"`
	VerifyName            string `json:"verify_name"`
}
