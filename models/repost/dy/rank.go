package dy

type TakeGoodsRankRet struct {
	Rank             int                      `json:"rank,omitempty"`
	Nickname         string                   `json:"nickname,omitempty"`
	AuthorCover      string                   `json:"author_cover,omitempty"`
	SumGmv           float64                  `json:"sum_gmv,omitempty"`
	SumSales         float64                  `json:"sum_sales,omitempty"`
	AvgPrice         float64                  `json:"avg_price,omitempty"`
	AuthorId         string                   `json:"author_id,omitempty"`
	UniqueId         string                   `json:"unique_id,omitempty"`
	Tags             string                   `json:"tags,omitempty"`
	VerificationType int                      `json:"verification_type,omitempty"`
	VerifyName       string                   `json:"verify_name,omitempty"`
	RoomCount        int                      `json:"room_count,omitempty"`
	RoomList         []map[string]interface{} `json:"room_list"`
}
