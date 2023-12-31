package dy

import "time"

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

type NameValueInt64PercentChart struct {
	Name    string  `json:"name"`
	Value   int64   `json:"value"`
	Percent float64 `json:"percent"`
}

type NameValueFloat64PercentChart struct {
	Name    string  `json:"name"`
	Value   float64 `json:"value"`
	Percent float64 `json:"percent"`
}

type NameValueInt64ChartWithData struct {
	Name  string    `json:"name"`
	Value int64     `json:"value"`
	Data  []string  `json:"data"`
	Date  time.Time `json:"date"`
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

type AppVersion struct {
	Version string `json:"version"`
	Info    string `json:"info"`
	Force   int    `json:"force"`
	Url     string `json:"url"`
}
