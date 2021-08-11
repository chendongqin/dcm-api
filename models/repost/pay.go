package repost

type VipOrderInfo struct {
	SurplusValue float64 `json:"surplus_value"`
	BuyDays      int     `json:"buy_days"`
	Amount       float64 `json:"amount"`
	People       int     `json:"people"`
	Title        string  `json:"title"`
}
