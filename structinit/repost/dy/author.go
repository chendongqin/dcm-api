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
