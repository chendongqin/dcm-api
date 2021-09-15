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
