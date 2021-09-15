package dy

import "time"

type LiveMonitor struct {
	Id            int       `json:"id" `
	UserId        int       `json:"user_id"`
	FreeCount     int       `json:"free_count"`
	PurchaseCount int       `json:"purchase_count"`
	AuthorId      string    `json:"author_id"`
	OpenId        string    `json:"-"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	NextTime      int64     `json:"next_time"`
	HasNew        int       `json:"has_new"`
	DelStatus     int       `json:"-"`
	Notice        int       `json:"notice"`
	FinishNotice  int       `json:"finish_notice"`
	ProductId     string    `json:"-"`
	Status        int       `json:"status"`
	Source        int       `json:"source"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	RoomId        string    `json:"room_id"`
	RoomCount     int       `json:"room_count"`
	Nickname      string    `json:"nickname"`
	Avatar        string    `json:"avatar"`
	UniqueID      string    `json:"unique_id"`
}

type LiveMonitorPriceList struct {
	MonitorPrice struct {
		Price10  LiveMonitorPrice `json:"price_10"`
		Price100 LiveMonitorPrice `json:"price_100"`
		Price500 LiveMonitorPrice `json:"price_500"`
	} `json:"monitor_price"`
}

type LiveMonitorPrice struct {
	OriginalPrice string `json:"original_price"`
	Price         string `json:"price"`
	OnePrice      string `json:"one_price"`
	Monitor       string `json:"monitor"`
}
