package dy

type LivingInfo struct {
	RoomId         string  `json:"room_id"`
	AuthorId       string  `json:"author_id"`
	Title          string  `json:"title"`
	Cover          string  `json:"cover"`
	CreateTime     int64   `json:"create_time"`
	Gmv            float64 `json:"gmv"`
	UserCount      int64   `json:"user_count"`
	TotalUserCount int64   `json:"total_user_count"`
	RoomStatus     int     `json:"room_status"`
	FinishTime     int64   `json:"finish_time"`
	LiveTime       int64   `json:"live_time"`
	Uv             float64 `json:"uv"`
	AvgOnlineTime  float64 `json:"avg_online_time"`
	BarrageRate    float64 `json:"barrage_rate"`
}
