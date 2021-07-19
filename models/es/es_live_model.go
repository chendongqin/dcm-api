package es

type EsAuthorLiveRoom struct {
	RoomID       string  `json:"room_id"`
	AuthorID     string  `json:"author_id"`
	Title        string  `json:"title"`
	CreateTime   string  `json:"create_time"`
	MaxUserCount int     `json:"max_user_count"`
	NumProducts  int     `json:"num_products"`
	Sales        float64 `json:"sales"`
	Gmv          float64 `json:"gmv"`
	Cover        string  `json:"cover"`
}
