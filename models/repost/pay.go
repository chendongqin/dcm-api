package repost

type VipOrderInfo struct {
	SurplusValue     float64 `json:"surplus_value"`
	BuyDays          int     `json:"buy_days"`
	Amount           float64 `json:"amount"`
	People           int     `json:"people"`
	IosPayProductId  string  `json:"ios_pay_product_id"`
	IosPayProductNum int     `json:"ios_pay_product_num"`
	MonitorNum       int     `json:"monitor_num"`
	Title            string  `json:"title"`
}
