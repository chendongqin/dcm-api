package dy

type RepostSimpleReputation struct {
	Score         float64 `json:"score"`
	Level         int     `json:"level"`
	EncryptShopID string  `json:"encrypt_shop_id"`
	ShopName      string  `json:"shop_name"`
	ShopLogo      string  `json:"shop_logo"`
}
