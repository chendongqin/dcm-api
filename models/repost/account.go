package repost

type SearchData struct {
	Id         int         `json:"id"`
	SearchType string      `json:"search_type"`
	Note       string      `json:"note"`
	Content    interface{} `json:"content"`
}

type VipOrderDetail struct {
	OrderId      int    `json:"order_id"`
	TradeNo      string `json:"trade_no"`
	OrderType    int    `json:"order_type"`
	PayType      string `json:"pay_type"`
	Level        int    `json:"level"`
	BuyDays      int    `json:"buy_days"`
	Title        string `json:"title"`
	Amount       string `json:"amount"`
	Channel      int    `json:"channel"`
	TicketAmount string `json:"ticket_amount"`
	Status       int    `json:"status"`
	PayStatus    int    `json:"pay_status"`
	CreateTime   string `json:"create_time"`
	PayTime      string `json:"pay_time"`
	InvoiceId    int    `json:"invoice_id"`
}
