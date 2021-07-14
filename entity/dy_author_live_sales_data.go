package entity

var DyAuthorLiveSalesDataMap = HbaseEntity{
	"crawl_time":   {Long, "crawl_time"},
	"gmv":          {Double, "gmv"},
	"num_product":  {Long, "num_product"},
	"sales":        {Double, "sales"},
	"ticket_count": {Long, "ticket_count"},
}

type DyAuthorLiveSalesData struct {
	CrawlTime   int64   `json:"crawl_time"`
	Gmv         float64 `json:"gmv"`
	NumProduct  int64   `json:"num_product"`
	Sales       float64 `json:"sales"`
	TicketCount int64   `json:"ticket_count"`
}
