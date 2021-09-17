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
	Year     VipPriceActive `json:"year"`
	HalfYear VipPriceActive `json:"half_year"`
	Month    VipPriceActive `json:"month"`
}

type VipPriceActive struct {
	Price         float64
	ActivePrice   float64
	ActiveComment string
}

func (v *VipPriceActive) GetPrice() float64 {
	if v.ActivePrice != 0 {
		return v.ActivePrice
	}
	return v.Price
}
