package dy

import "dongchamao/entity"

type DyLiveInfo struct {
	Cover               string                            `json:"cover"`       //封面
	CreateTime          int64                             `json:"create_time"` //开播时间
	FinishTime          int64                             `json:"finish_time"` //结束时间
	LikeCount           int64                             `json:"like_count"`  //点赞数
	RoomID              string                            `json:"room_id"`
	RoomStatus          int                               `json:"room_status"` //直播状态 2:在播 4:下播
	Title               string                            `json:"title"`
	TotalUser           int64                             `json:"total_user"` //总pv
	User                DyLiveUserSimple                  `json:"user"`
	UserCount           int64                             `json:"user_count"`        //当前在线人数
	AvgUserCount        int64                             `json:"avg_user_count"`    //平均当前在线人数
	TrendsCrawlTime     int64                             `json:"trends_crawl_time"` //更新时间
	IncFans             int64                             `json:"inc_fans"`
	IncFansRate         float64                           `json:"inc_fans_rate"`
	InteractRate        float64                           `json:"interact_rate"`
	OnlineTrends        entity.DyLiveIncOnlineTrendsChart `json:"online_trends"`
	MaxWatchOnlineTrend entity.DyLiveOnlineTrends         `json:"max_watch_online_trend"`
	RenewalTime         int64                             `json:"renewal_time"`
	AvgOnlineTime       float64                           `json:"avg_online_time"`
}

type DyLiveUserSimple struct {
	Avatar          string  `json:"avatar"`
	FollowerCount   int64   `json:"follower_count"`
	ID              string  `json:"id"`
	Nickname        string  `json:"nickname"`
	WithCommerce    bool    `json:"with_commerce"`
	ReputationScore float64 `json:"reputation_score"`
	ReputationLevel int     `json:"reputation_level"`
}

type DyLivePromotion struct {
	ProductID string  `json:"product_id"` //第三方商品id
	ForSale   int     `json:"for_sale"`   //商品状态 0:刚上架 2:在售 4:下架
	StartTime int64   `json:"start_time"` //上架时间
	StopTime  int64   `json:"stop_time"`  //下架时间
	Price     float64 `json:"price"`      //价格
	Sales     int64   `json:"sales"`      //全网销量
	NowSales  int64   `json:"now_sales"`  //本场当前实时销量
	GmvSales  int64   `json:"gmv_sales"`  //当前销量
	Title     string  `json:"title"`      //标题
	Cover     string  `json:"cover"`      //封面
	Index     int     `json:"index"`      //第几个商品
	SaleNum   int     `json:"sale_num"`   //上架次数
}

type DyLivePromotionChart struct {
	StartTime     []string            `json:"start_time"`
	PromotionList [][]DyLivePromotion `json:"promotion_list"`
}

type DyLiveRoomAnalyse struct {
	TotalUserCount int64   `json:"total_user_count"`
	IncFans        int64   `json:"inc_fans"`
	IncFansRate    float64 `json:"inc_fans_rate"`
	InteractRate   float64 `json:"interact_rate"`
	AvgUserCount   int64   `json:"avg_user_count"`
	Volume         int64   `json:"volume"`
	Amount         float64 `json:"amount"`
	Uv             float64 `json:"uv"`
	PromotionNum   int64   `json:"promotion_num"`
	SaleRate       float64 `json:"sale_rate"`
	PerPrice       float64 `json:"per_price"`
	LiveLongTime   int64   `json:"live_long_time"`
	LiveStartTime  int64   `json:"live_start_time"`
	AvgOnlineTime  float64 `json:"avg_online_time"`
}

type DyLiveRoomSaleData struct {
	Volume       int64   `json:"volume"`
	Amount       float64 `json:"amount"`
	Uv           float64 `json:"uv"`
	PromotionNum int64   `json:"promotion_num"`
	SaleRate     float64 `json:"sale_rate"`
	PerPrice     float64 `json:"per_price"`
}
