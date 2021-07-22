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
