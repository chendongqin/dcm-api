package dy

type VipPrice struct {
	VipPrice []struct {
		Title           string `json:"title"`
		Desc            string `json:"desc"`
		Tag             string `json:"tag"`
		PrimePriceValue string `json:"prime_price_value"`
		Days            string `json:"days"`
		Price           string `json:"price"`
		PrimePrice      string `json:"prime_price"`
	} `json:"vip_price"`
}

type VipPriceConfig struct {
	Year     float64 `json:"year"`
	HalfYear float64 `json:"half_year"`
	Month    float64 `json:"month"`
}
