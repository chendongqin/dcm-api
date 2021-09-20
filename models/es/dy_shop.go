package es

type DyShop struct {
	ShopId           string  `json:"shop_id"`
	Logo             string  `json:"logo"`
	ShopName         string  `json:"shop_name"`
	Score            float64 `json:"score"`
	Level            string  `json:"level"`
	MonthSales       int64   `json:"month_sales"`
	MonthGmv         float64 `json:"month_gmv"`
	MonthSinglePrice float64 `json:"month_single_price"`
	ProductNum       int64   `json:"product_num"`
	RelateAweme      int64   `json:"relate_aweme"`
	RelateRoom       int64   `json:"relate_room"`
	RelateAuthor     int64   `json:"relate_author"`
	IsBrand          int     `json:"is_brand"`
	IsCollect        int     `json:"is_collect"`
}
