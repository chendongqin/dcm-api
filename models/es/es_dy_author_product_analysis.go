package es

type EsDyAuthorProductAnalysis struct {
	AuthorProductDate string `json:"author_product_date"`
	AuthorId          string `json:"author_id"`
	ProductId         string `json:"product_id"`
	Title             string `json:"title"`
	CreateTime        string `json:"create_time"`
	ShelfTime         int64  `json:"shelf_time"`
}
