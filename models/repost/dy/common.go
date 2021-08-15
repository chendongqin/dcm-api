package dy

type DateChart struct {
	Date       []string `json:"date"`
	CountValue []int64  `json:"count_value"`
	IncValue   []int64  `json:"inc_value"`
}

type DateCountChart struct {
	Date       []string `json:"date"`
	CountValue []int64  `json:"count_value"`
}

type TimestampCountChart struct {
	Timestamp  []int64   `json:"timestamp"`
	CountValue []float64 `json:"count_value"`
}

type DateCountFChart struct {
	Date       []string  `json:"date"`
	CountValue []float64 `json:"count_value"`
}

type NameValueChart struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type NameValueInt64Chart struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type NameValueFloat64Chart struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type DyCate struct {
	Name    string   `json:"name"`
	Num     int      `json:"num"`
	SonCate []DyCate `json:"son_cate"`
}

type DyIntCate struct {
	Name    int      `json:"name"`
	Num     int      `json:"num"`
	SonCate []DyCate `json:"son_cate"`
}
