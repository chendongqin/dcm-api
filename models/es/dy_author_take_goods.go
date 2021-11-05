package es

type DyAuthorTakeGoods struct {
	AuthorId             string  `json:"author_id"`
	ShortId              string  `json:"short_id"`
	Nickname             string  `json:"nickname"`
	Avatar               string  `json:"avatar"`
	VerificationType     int     `json:"verification_type"`
	VerifyName           string  `json:"verify_name"`
	UniqueId             string  `json:"unique_id"`
	PredictSalesSum      float64 `json:"predict_sales_sum"`
	PredictGmvSum        float64 `json:"predict_gmv_sum"`
	PerPrice             float64 `json:"per_price"`
	RoomIdCountWithSales string  `json:"room_id_count_with_sales"`
	Tags                 string  `json:"tags"`
	RoomIdCount          string  `json:"room_id_count"`
	Dt                   string  `json:"dt"`
}

type DyAuthorTakeGoodsCount struct {
	AvgPrice struct {
		Value float64 `json:"value"`
	} `json:"avg_price"`
	DocCount int `json:"doc_count"`
	Hit      struct {
		Hits struct {
			Hits []struct {
				ID     string  `json:"_id"`
				Index  string  `json:"_index"`
				Score  float64 `json:"_score"`
				Source struct {
					AuthorCover      string  `json:"author_cover"`
					AuthorID         string  `json:"author_id"`
					AuthorRoomID     string  `json:"author_room_id"`
					CreateTime       int     `json:"create_time"`
					DateTime         string  `json:"date_time"`
					DiscoverTime     int     `json:"discover_time"`
					MaxUserCount     int     `json:"max_user_count"`
					Nickname         string  `json:"nickname"`
					PredictGmv       float64 `json:"predict_gmv"`
					PredictSales     float64 `json:"predict_sales"`
					RealGmv          float64 `json:"real_gmv"`
					RealSales        float64 `json:"real_sales"`
					RoomCover        string  `json:"room_cover"`
					RoomID           string  `json:"room_id"`
					RoomTitle        string  `json:"room_title"`
					ShortID          string  `json:"short_id"`
					Tags             string  `json:"tags"`
					UniqueID         string  `json:"unique_id"`
					VerificationType int     `json:"verification_type"`
					VerifyName       string  `json:"verify_name"`
				} `json:"_source"`
				Type string `json:"_type"`
			} `json:"hits"`
			MaxScore float64 `json:"max_score"`
			Total    int     `json:"total"`
		} `json:"hits"`
	} `json:"hit"`
	Key struct {
		AuthorID string `json:"author_id"`
	} `json:"key"`
	SumGmv struct {
		Value float64 `json:"value"`
	} `json:"sum_gmv"`
	SumSales struct {
		Value float64 `json:"value"`
	} `json:"sum_sales"`
}
