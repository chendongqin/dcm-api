package es

type DyAuthorTakeGoods struct {
	AuthorRoomId     string  `json:"author_room_id"`
	RoomId           string  `json:"room_id"`
	AuthorId         string  `json:"author_id"`
	RoomTitle        string  `json:"room_title"`
	RoomCover        string  `json:"room_cover"`
	CreateTime       int64   `json:"create_time"`
	DiscoverTime     int64   `json:"discover_time"`
	PredictSales     float64 `json:"predict_sales"`
	PredictGmv       float64 `json:"predict_gmv"`
	RealSales        float64 `json:"real_sales"`
	RealGmv          float64 `json:"real_gmv"`
	MaxUserCount     int     `json:"max_user_count"`
	Nickname         string  `json:"nickname"`
	ShortId          string  `json:"short_id"`
	UniqueId         string  `json:"unique_id"`
	AuthorCover      string  `json:"author_cover"`
	VerificationType int     `json:"verification_type"`
	VerifyName       string  `json:"verify_name"`
	Tags             string  `json:"tags"`
	DateTime         string  `json:"date_time"`
}
