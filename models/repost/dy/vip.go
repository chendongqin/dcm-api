package dy

type VipPrice struct {
	VipPrice []struct {
		Unit      string `json:"unit"`
		Month     string `json:"month"`
		Days      string `json:"days"`
		Price     string `json:"price"`
		InitPrice string `json:"init_price"`
	} `json:"vip_price"`
}

type VipPriceConfig struct {
	Year     float64 `json:"year"`
	HalfYear float64 `json:"half_year"`
	Month    float64 `json:"month"`
}
