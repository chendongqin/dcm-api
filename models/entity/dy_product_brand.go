package entity

var DyProductBrandMap = HbaseEntity{
	"dy_promotion_id":        {String, "dy_promotion_id"},
	"crawl_time":             {Long, "crawl_time"},
	"sec_shop_id":            {String, "sec_shop_id"},
	"shop_id":                {String, "shop_id"},
	"shop_name":              {String, "shop_name"},
	"shop_tel":               {String, "shop_tel"},
	"shop_icon":              {String, "shop_icon"},
	"shop_url":               {String, "shop_url"},
	"brand_name":             {String, "brand_name"},
	"product_desc":           {Double, "product_desc"},
	"product_desc_level":     {String, "product_desc_level"},
	"product_serv":           {Double, "product_serv"},
	"product_serv_level":     {String, "product_serv_level"},
	"product_post":           {Double, "product_post"},
	"product_post_level":     {String, "product_post_level"},
	"platform":               {String, "platform"},
	"subtitle":               {String, "subtitle"},
	"product_format":         {String, "product_format"},
	"company_name":           {String, "company_name"},
	"category":               {Json, "category"},
	"dcm_level_first":        {String, "dcm_level_first"},
	"sales":                  {Long, "sales"},
	"price":                  {Double, "price"},
	"other_manmade_category": {Json, "manmade_category"},
}

type DyProductBrand struct {
	DyPromotionId    string                   `json:"dy_promotion_id"`
	CrawlTime        int64                    `json:"crawl_time"`
	SecShopId        string                   `json:"sec_shop_id"`
	ShopId           string                   `json:"shop_id"`
	ShopName         string                   `json:"shop_name"`
	ShopTel          string                   `json:"shop_tel"`
	ShopIcon         string                   `json:"shop_icon"`
	ShopUrl          string                   `json:"shop_url"`
	BrandName        string                   `json:"brand_name"`
	ProductDesc      float64                  `json:"product_desc"`
	ProductDescLevel string                   `json:"product_desc_level"`
	ProductServ      float64                  `json:"product_serv"`
	ProductServLevel string                   `json:"product_serv_level"`
	ProductPost      float64                  `json:"product_post"`
	ProductPostLevel string                   `json:"product_post_level"`
	Platform         string                   `json:"platform"`
	Subtitle         string                   `json:"subtitle"`
	ProductFormat    string                   `json:"product_format"`
	CompanyName      string                   `json:"company_name"`
	Category         DyProductAiCategory      `json:"category"`
	DcmLevelFirst    string                   `json:"dcm_level_first"`
	Sales            int64                    `json:"sales"`
	Price            float64                  `json:"price"`
	ManmadeCategory  DyProductManmadeCategory `json:"manmade_category"`
}
