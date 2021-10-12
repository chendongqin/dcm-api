package command

import (
	"dongchamao/business"
	"dongchamao/global"
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

//监听除了小时榜和商品榜的（非整点）榜单
func CheckRank() {
	currentHour := time.Now().Hour()
	currentHourString := strconv.Itoa(currentHour)
	currentHourString = fmt.Sprintf("%s:30", currentHourString)
	loopCheck(currentHourString)
	return
}

//监控商品（整点）榜单
func CheckGoodsRank() {
	currentHour := time.Now().Hour()
	currentHourString := strconv.Itoa(currentHour)
	loopCheck(currentHourString)
	return
}

//监听小时榜
func CheckRankHour() {
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
			pathInfo := getRoute(v)
			monitorEvents = append(monitorEvents, pathInfo.Desc)
			checkLiveHotRank(pathInfo)
		}
		monitorBusiness.SendErr("直播监控", strings.Join(monitorEvents, ","))
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
func getRoute(key string) (pathInfo PathDesc) {
	toDate := time.Now().Format("2006-01-02")
	yesDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	//BeforeYesDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	weekDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	hour := time.Now().Hour()
	currentHourString := strconv.Itoa(hour)

	var routeMap = map[string]PathDesc{
		/*********直播*********/
		"live_hour": {fmt.Sprintf("/v1/dy/rank/live/hour/%s/%s:00", toDate, currentHourString), "直播小时榜"},
		"live_top":  {fmt.Sprintf("/v1/dy/rank/live/top/%s/%s:00", toDate, currentHourString), "直播热榜"},
		/*********商品*********/
		"product_sale":           {fmt.Sprintf("/v1/dy/rank/product/sale/%s?data_type=1&first_cate=&order_by=desc&sort=order_count&page=1&page_size=50", yesDate), "抖音销量榜"},
		"product_share":          {fmt.Sprintf("/v1/dy/rank/product/share/%s?first_cate=&data_type=1&order_by=desc&sort=share_count&page=1&page_size=50", yesDate), "抖音热推榜"},
		"product_live_sale":      {fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=1&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", yesDate), "直播商品榜"},
		"product_live_sale_week": {fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=2&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", toDate), "直播商品榜-周榜"},
		"product":                {fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=1&page=1&page_size=50", yesDate), "视频商品榜"},
		"product_week":           {fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=2&page=1&page_size=50", toDate), "视频商品榜-周榜"},
		/*********达人*********/
		"author_follower_inc": {fmt.Sprintf("/v1/dy/rank/author/follower/inc/%s?tags=&province=&page=1&is_delivery=0&page_size=50&order_by=desc&sort=inc_follower_count", yesDate), "达人涨粉榜"},
		"author_goods":        {fmt.Sprintf("/v1/dy/rank/author/goods/%s?date_type=1&tags=&verified=0&page=1&page_size=50&sort=sum_gmv&order_by=desc", yesDate), "达人带货榜"},
		//"video_share":         {fmt.Sprintf("/v1/dy/rank/video/share/%s", BeforeYesDate), "电商视频达人分享榜"},
		"live_share":        {fmt.Sprintf("/v1/dy/rank/live/share/%s/%s", weekDate, yesDate), "电商直播达人分享榜"},
		"author_aweme_rank": {"/v1/dy/rank/author/aweme?rank_type=达人指数榜&category=全部", "抖音短视频达人热榜"},
		"author_aweme_live": {"/v1/dy/rank/author/live?rank_type=达人指数榜", "抖音直播主播热榜"},
	}
	pathInfo = routeMap[key]
	return
}

func getRow(hour string) (taskList []string) {
	hourGroup := getHourGroup()
	taskList = hourGroup[hour]
	return
}
func getHourGroup() (hourGroup map[string][]string) {
	hourGroup = map[string][]string{
		"every": {"live_hour", "live_top"},
		"10":    {"product_live_sale", "product_live_sale_week", "product_sale", "product_share"},
		"10:30": {"author_follower_inc", "author_goods"},
		"12:30": {"live_share"},
		"15":    {"product_week", "product"},
		"15:30": {"author_aweme_live"},
		//"16:30": {"video_share", "author_aweme_rank"},
		"16:30": {"author_aweme_rank"},
	}
	return
}

//demo
func checkLiveHotRank(pathInfo PathDesc) {
	fmt.Println(pathInfo)
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
		business.NewMonitorBusiness().SendTemplateMessage("S", pathInfo.Desc, fmt.Sprintf("%s挂了，请求地址：%s", pathInfo.Desc, pathInfo.Path))
	}
	return
}

// SwitchTopDateTime 根据key返回对应榜单需要展示的日期时间
func SwitchTopDateTime(key string) (res map[string][]string, comErr global.CommonError) {
	if key == "author_aweme_rank" || key == "author_aweme_live" {
		comErr = global.NewMsgError("传入参数错误，不存在的key")
		return
	}
	desc := getRoute(key).Desc
	if desc == "" {
		comErr = global.NewMsgError("传入参数错误，不存在的key")
		return
	}
	switch key {
	case "live_hour":
		res = dateTimeLiveHour()
	case "live_top":
		res = dateTimeLiveHour()
	case "product_sale":
		res = getCheckDateList(key)
	case "product_share":
		res = getCheckDateList(key)
	case "product_live_sale":
		res = getCheckDateList(key)
	case "product_live_sale_week":
		res = weekList()
	case "product":
		res = getCheckDateList(key)
	case "product_week":
		res = weekList()
	case "author_follower_inc":
		res = getCheckDateList(key)
	case "author_goods":
		res = getCheckDateList(key)
	case "live_share":
		res = getCheckDateList(key)
	}
	res["desc"] = []string{fmt.Sprintf("%s的日期时间", desc)}
	return
}

//直播小时榜的日期和时间/直播全天热榜，目前这两个榜单日期时间一样
func dateTimeLiveHour() (res map[string][]string) {
	res = map[string][]string{"date": {}, "hour": {}}
	now := time.Now()
	dateList := getDateList(7, now)
	var hourList []string
	currentHour := now.Hour()
	for i := 0; i <= currentHour; i++ {
		hourString := fmt.Sprintf("%02d:00", currentHour-i)
		hourList = append(hourList, hourString)
	}
	res["date"] = dateList
	res["hour"] = hourList
	return
}

//根据监控时间判断获取日期列表获取
func getCheckDateList(key string) (res map[string][]string) {
	res = map[string][]string{"date": {}, "hour": {}}
	now := time.Now()
	isBefore := checkIsBefore(key)
	beforeInt := -2
	if isBefore {
		beforeInt = -1
	}
	startSate := now.AddDate(0, 0, beforeInt)
	res["date"] = getDateList(30, startSate)
	return
}

//周榜日期列表获取
func weekList() (res map[string][]string) {
	//这里仿照前段，只给三个切片
	res = map[string][]string{"date": {}, "hour": {}}
	num := 3
	startDateTime := time.Now().AddDate(0, 0, -2)
	var dateSelectList []string
	for i := 0; i < num; i++ {
		rightDate := startDateTime.AddDate(0, 0, -i*6)
		leftDate := startDateTime.AddDate(0, 0, -(i+1)*6)
		dateString := fmt.Sprintf("%s-%s", leftDate.Format("01/02"), rightDate.Format("01/02"))
		dateSelectList = append(dateSelectList, dateString)
		startDateTime = startDateTime.AddDate(0, 0, -1)
	}
	res["date"] = dateSelectList
	return
}

//检测当前是否已经超过了监控时间
func checkIsBefore(key string) (isBefore bool) {
	hourGroup := getHourGroup()
	groupKey := ""
	for k, v := range hourGroup {
		if utils.InArrayString(key, v) {
			groupKey = k
		}
	}
	isInclude := strings.Contains(groupKey, ":30")
	if !isInclude {
		groupKey = fmt.Sprintf("%s:00:00", groupKey)
	} else {
		groupKey = fmt.Sprintf("%s:00", groupKey)
	}
	dateTime := fmt.Sprintf("%s %s", time.Now().Format("2006-01-02"), groupKey)
	dateTimeStampObj, _ := time.ParseInLocation("2006-01-02 15:04:05", dateTime, time.Local)
	isBefore = dateTimeStampObj.Before(time.Now())
	return

}

//获取日期列表
func getDateList(daysCount int, startTime time.Time) (list []string) {
	for i := 0; i < daysCount; i++ {
		date := startTime.AddDate(0, 0, -i).Format("2006-01-02")
		list = append(list, date)
	}
	return
}
