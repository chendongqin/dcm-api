package entity

var DyShopMap = HbaseEntity{
	"avg_cos_ratio":   {Double, "avg_cos_ratio"},
	"brand":           {Int, "is_brand"},
	"coo_kol_num":     {Long, "coo_kol_num"},
	"expr_score":      {Double, "expr_score"},
	"logistics_score": {Json, "logistics_score"},
	"product_score":   {Json, "product_score"},
	"service_score":   {Json, "service_score"},
	"logo":            {String, "logo"},
	"name":            {String, "name"},
	"order_num":       {Long, "order_num"},
	"sales":           {Long, "sales"},
	"unique_id":       {String, "unique_id"},
	"short_id":        {String, "short_id"},
	"crawl_time":      {Long, "crawl_time"},
}

type DyShop struct {
	ShopId         string      `json:"shop_id"`
	AvgCosRatio    float64     `json:"avg_cos_ratio"`   //平均佣金
	CooKolNum      float64     `json:"coo_kol_num"`     //总合作达人数
	ExprScore      float64     `json:"expr_score"`      //体验分
	LogisticsScore DyShopScore `json:"logistics_score"` //物流体验
	ProductScore   DyShopScore `json:"product_score"`   //商品体验
	ServiceScore   DyShopScore `json:"service_score"`   //商家服务
	Logo           string      `json:"logo"`            //logo
	Name           string      `json:"name"`            //名称
	OrderNum       int64       `json:"order_num"`       //合作商品商量
	Sales          int64       `json:"sales"`           //月销量
	CrawlTime      int64       `json:"crawl_time"`      //爬虫时间
	IsBrand        int64       `json:"is_brand"`        //是否品牌
	Level          string      `json:"level"`           //店铺口碑
	UniqueId       string      `json:"unique_id"`       //抖音号
	ShortId        string      `json:"short_id"`
}

type DyShopScore struct {
	Level string `json:"level"`
	Score string `json:"score"`
}

type DyShopBaseBasic struct {
	BaseData   DyShop           ` json:"base_data"`
	DetailData DyShopBaseDetail `json:"detail_data"`
}

type DyShopBaseDetail struct {
	Sales        int64       `json:"sales"`          //日销量
	Gmv          float64     `json:"gmv"`            //每日gmv
	D30AwemeCnt  int64       `json:"30d_aweme_cnt"`  //30天视频数
	D30LiveCnt   int64       `json:"30d_live_cnt"`   //30天直播数
	D30AuthorCnt int64       `json:"30d_author_cnt"` //30天达人数
	D30Sales     int64       `json:"30d_sales"`      //30天销量
	D30Gmv       float64     `json:"30d_gmv"`        //30天gmv
	D30Pct       float64     `json:"30d_pct"`        //30天客单价
	D30Rate      float64     `json:"30d_rate"`       //30天转化率
	ProductCnt   int64       `json:"product_cnt"`    //商品数
	ShopCName    []ShopCName `json:"shop_cname"`
}

//小店商品品类销售额TOP
type GoodsCatTop struct {
	Name  string `json:"name"`  //品类名称
	Value int64  `json:"value"` //商品数量
	Gmv   int64  `json:"gmv"`   //销售额
	Sales int64  `json:"sales"` //销量

}
