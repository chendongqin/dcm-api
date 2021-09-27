package command

import (
	"dongchamao/business"
	"dongchamao/global/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PathDesc struct {
	Path string `json:"path"`
	Desc string `json:"desc"`
}

//监听除了小时榜的榜单
func CheckRank() {
	//监控固定时段榜单
	currentHour := time.Now().Hour()
	currentHourString := strconv.Itoa(currentHour)
	loopCheck(currentHourString)
	return
}

//监听小时榜
func CheckRankHour() {
	//监控小时榜
	currentHourString := "every"
	loopCheck(currentHourString)
	return
}

//遍历执行监控
func loopCheck(hour string) {
	monitorBusiness := business.NewMonitorBusiness()
	var monitorEvents []string

	taskList := getRow(hour)
	length := len(taskList)
	if length > 0 {
		for _, v := range taskList {
			pathInfo := getRoute(v, hour)

			checkLiveHotRank(pathInfo)
			monitorBusiness.SendErr(pathInfo.Desc, strings.Join(monitorEvents, ","))
			fmt.Println(pathInfo)
			fmt.Println(pathInfo.Desc)
		}
	}
}

func checkTime(currentTime time.Time, hour, minute int) bool {
	if currentTime.Hour() == hour && currentTime.Minute() == minute {
		return true
	}
	return false
}

/**
**name:榜单名称
**hour：小时
 */
func getRoute(key, hour string) (pathInfo PathDesc) {
	toDate := time.Now().Format("20060102")
	yesDate := time.Now().AddDate(0, 0, -1).Format("20060102")
	BeforeYesDate := time.Now().AddDate(0, 0, -1).Format("20060102")
	weekDate := time.Now().AddDate(0, 0, -7).Format("20060102")

	var routeMap = map[string]PathDesc{
		/*********直播*********/
		"live_hour": {fmt.Sprintf("/v1/dy/rank/live/hour/%s/%s", toDate, hour), "直播小时榜"},
		"live_top":  {fmt.Sprintf("/v1/dy/rank/live/hour/%s/%s", toDate, hour), "直播热榜"},
		/*********商品*********/
		"product_sale":           {fmt.Sprintf("/v1/dy/rank/product/sale/%s?data_type=1&first_cate=&order_by=desc&sort=order_count&page=1&page_size=50", toDate), "抖音销量榜"},
		"product_share":          {fmt.Sprintf("/v1/dy/rank/product/share/%s?first_cate=&data_type=1&order_by=desc&sort=share_count&page=1&page_size=50", toDate), "抖音热推榜"},
		"product_live_sale":      {fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=1&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", toDate), "直播商品榜"},
		"product_live_sale_week": {fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=2&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", toDate), "直播商品榜-周榜"},
		"product":                {fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=1&page=1&page_size=50", yesDate), "视频商品榜"},
		"product_week":           {fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=2&page=1&page_size=50", toDate), "视频商品榜-周榜"},
		/*********达人*********/
		"author_follower_inc": {fmt.Sprintf("/v1/dy/rank/author/follower/inc/%s?tags=&province=&page=1&is_delivery=0&page_size=50&order_by=desc&sort=inc_follower_count", yesDate), "达人涨粉榜"},
		"author_goods":        {fmt.Sprintf("/v1/dy/rank/author/goods/%s?date_type=1&tags=&verified=0&page=1&page_size=50&sort=sum_gmv&order_by=desc", yesDate), "达人带货榜"},
		"video_share":         {fmt.Sprintf("/v1/dy/rank/video/share/%s", BeforeYesDate), "电商视频达人分享榜"},
		"live_share":          {fmt.Sprintf("/v1/dy/rank/live/share/%s/%s", weekDate, yesDate), "电商直播达人分享榜"},
		"author_aweme_rank":   {"/v1/dy/rank/author/aweme?rank_type=达人指数榜&category=全部", "抖音短视频达人热榜"},
		"author_aweme_live":   {"/v1/dy/rank/author/live?rank_type=达人指数榜", "抖音直播主播热榜"},
	}
	pathInfo = routeMap[key]
	return
}

func getRow(hour string) (taskList []string) {
	hourGroup := map[string][]string{
		"every": {"live_hour", "live_top"},
		"10":    {"product_sale", "live_top", "product_live_sale", "product_live_sale_week", "author_follower_inc", "author_goods"},
		"13":    {"live_share"},
		"15":    {"product", "product_week", "author_aweme_live"},
		"16":    {"video_share", "author_aweme_rank"},
	}
	taskList = hourGroup[hour]
	return
}

//demo
func checkLiveHotRank(pathInfo PathDesc) {
	interBusiness := business.NewInterBusiness()
	testApi := interBusiness.BuildURL(pathInfo.Path)
	res, comErr := interBusiness.HttpGet(testApi)
	if comErr != nil {
		return
	}
	checkRes := false
	resMap := map[string]interface{}{}
	utils.MapToStruct(res, &resMap)
	if v, exist := resMap["list"]; exist {
		if len(v.([]interface{})) > 0 {
			checkRes = true
		}
	}
	if !checkRes {
		business.NewMonitorBusiness().SendTemplateMessage("S", pathInfo.Desc, fmt.Sprintf("%s挂了", pathInfo.Desc))
	}
	return
}
