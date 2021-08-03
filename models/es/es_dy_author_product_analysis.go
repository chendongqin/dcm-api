package es

type EsDyAuthorProductAnalysis struct {
	AuthorDateProduct string `json:"author_date_product"`
	AuthorId          string `json:"author_id"`
	ProductId         string `json:"product_id"`
	Title             string `json:"title"`
	CreateTime        string `json:"create_time"`
	ShelfTime         int64  `json:"shelf_time"`
}
