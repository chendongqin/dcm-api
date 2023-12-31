package entity

type DyCommodityTopN struct {
	UpdateTime int64                   `json:"update_time"`
	Ranks      []DyProductSalesTopRank `json:"ranks"`
}

type DyProductSalesTopRank struct {
	ProductId         string  `json:"product_id"`
	DateTime          string  `json:"date_time"`
	OrderAccountCount int64   `json:"order_account_count"`
	OrderAccountPv    int64   `json:"order_account_pv"`
	OrderCount        int64   `json:"order_count"`
	OrderPv           int64   `json:"order_pv"`
	Title             string  `json:"title"`
	MarkerPrice       float64 `json:"marker_price"`
	Price             float64 `json:"price"`
	Images            string  `json:"images"`
	Gpm               float64 `json:"gpm"`
	CosRatio          float64 `json:"cos_ratio"`
	CosFee            float64 `json:"cos_fee"`
	FirstCname        string  `json:"first_cname"`
	SecondCname       string  `json:"second_cname"`
	ThirdCname        string  `json:"third_cname"`
	DcmCname          string  `json:"dcm_cname"`
	TbFirstCname      string  `json:"tb_first_cname"`
	TbSecondCname     string  `json:"tb_second_cname"`
	PlatformLabel     string  `json:"platform_label"`
	ConversionRate    float64 `json:"conversion_rate"`
}
