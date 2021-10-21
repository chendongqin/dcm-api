package command

import (
	"dongchamao/business"
	"dongchamao/global"
	"dongchamao/global/cache"
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
func getRoute(key string, hours ...int) (pathInfo PathDesc) {
	now := time.Now()
	toDate := now.Format("2006-01-02")
	yesDate := now.AddDate(0, 0, -1).Format("2006-01-02")
	//BeforeYesDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	weekDate := now.AddDate(0, 0, -7).Format("2006-01-02")
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartDate := now.AddDate(0, 0, offset).Format("2006-01-02") //本周第一天
	hour := time.Now().Hour()
	if len(hours) > 0 {
		hour = hours[0]
	}
	currentHourString := strconv.Itoa(hour)

	var routeMap = map[string]PathDesc{
		/*********直播*********/
		"live_hour": {fmt.Sprintf("/v1/dy/rank/live/hour/%s/%s:00", toDate, currentHourString), "直播小时榜"},
		"live_top":  {fmt.Sprintf("/v1/dy/rank/live/top/%s/%s:00", toDate, currentHourString), "直播热榜"},
		/*********商品*********/
		"product_sale":           {fmt.Sprintf("/v1/dy/rank/product/sale/%s?data_type=1&first_cate=&order_by=desc&sort=order_count&page=1&page_size=50", yesDate), "抖音销量榜"},
		"product_share":          {fmt.Sprintf("/v1/dy/rank/product/share/%s?first_cate=&data_type=1&order_by=desc&sort=share_count&page=1&page_size=50", yesDate), "抖音热推榜"},
		"product_live_sale":      {fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=1&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", yesDate), "直播商品榜"},
		"product_live_sale_week": {fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=2&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", weekStartDate), "直播商品榜-周榜"},
		"product":                {fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=1&page=1&page_size=50", yesDate), "视频商品榜"},
		"product_week":           {fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=2&page=1&page_size=50", weekStartDate), "视频商品榜-周榜"},
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

func getCheckRoute(key string, now time.Time) (pathInfo PathDesc) {
	//now := time.Now()
	toDate := now.Format("2006-01-02")
	yesDate := now.AddDate(0, 0, -1).Format("2006-01-02")
	//BeforeYesDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	weekDate := now.AddDate(0, 0, -7).Format("2006-01-02")
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartDate := now.AddDate(0, 0, offset).Format("2006-01-02") //本周第一天
	hour := now.Hour()
	currentHourString := fmt.Sprintf("%02d", hour)

	var routeMap = map[string]PathDesc{
		/*********直播*********/
		"live_hour": {fmt.Sprintf("/v1/dy/rank/live/hour/%s/%s:00", toDate, currentHourString), "直播小时榜"},
		"live_top":  {fmt.Sprintf("/v1/dy/rank/live/top/%s/%s:00", toDate, currentHourString), "直播热榜"},
		/*********商品*********/
		"product_sale":           {fmt.Sprintf("/v1/dy/rank/product/sale/%s?data_type=1&first_cate=&order_by=desc&sort=order_count&page=1&page_size=50", yesDate), "抖音销量榜"},
		"product_share":          {fmt.Sprintf("/v1/dy/rank/product/share/%s?first_cate=&data_type=1&order_by=desc&sort=share_count&page=1&page_size=50", yesDate), "抖音热推榜"},
		"product_live_sale":      {fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=1&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", yesDate), "直播商品榜"},
		"product_live_sale_week": {fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=2&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", weekStartDate), "直播商品榜-周榜"},
		"product":                {fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=1&page=1&page_size=50", yesDate), "视频商品榜"},
		"product_week":           {fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=2&page=1&page_size=50", weekStartDate), "视频商品榜-周榜"},
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

//监控报警对应的榜单
func checkLiveHotRank(pathInfo PathDesc) {
	checkRes := requestRank(pathInfo)
	if !checkRes {
		business.NewMonitorBusiness().SendTemplateMessage("S", pathInfo.Desc, fmt.Sprintf("%s挂了，请求地址：%s", pathInfo.Desc, pathInfo.Path))
	}
	return
}

//请求对应的榜单
func requestRank(pathInfo PathDesc) (checkRes bool) {
	interBusiness := business.NewInterBusiness()
	testApi := interBusiness.BuildURL(pathInfo.Path)
	res, comErr := interBusiness.HttpGet(testApi)
	if comErr != nil {
		return
	}
	checkRes = false
	resMap := map[string]interface{}{}
	utils.MapToStruct(res, &resMap)
	if v, exist := resMap["list"]; exist {
		list := make([]interface{}, 0)
		utils.MapToStruct(v, &list)
		if len(list) > 0 {
			checkRes = true
		}
	}
	return
}

// SwitchTopDateTime 根据key返回对应榜单需要展示的日期时间
func SwitchTopDateTime(key string) (main map[string][]string, hourList map[string][]string, weekList []map[string]string, monthList []string, comErr global.CommonError) {
	if key == "author_aweme_rank" || key == "author_aweme_live" {
		comErr = global.NewMsgError("传入参数错误，不存在的key")
		return
	}
	desc := getRoute(key).Desc
	if desc == "" {
		comErr = global.NewMsgError("传入参数错误，不存在的key")
		return
	}
	hourList = map[string][]string{}
	weekList = []map[string]string{}
	monthList = []string{}
	main = make(map[string][]string)
	switch key {
	case "live_hour":
		main, hourList = dateTimeLiveHour(key)
	case "live_top":
		main, hourList = dateTimeLiveHour(key)
	case "product_sale":
		main = getCheckDateList(key)
	case "product_share":
		main = getCheckDateList(key)
	case "product_live_sale":
		main = getCheckDateList(key)
		weekList = getWeekList(key)
		monthList = getMonthList()
	case "product":
		main = getCheckDateList(key)
		weekList = getWeekList(key)
		monthList = getMonthList()
	case "author_follower_inc":
		main = getCheckDateList(key)
	case "author_goods":
		main = getCheckDateList(key)
	case "live_share":
		weekList = getWeekList(key)
	}
	main["desc"] = []string{fmt.Sprintf("%s的日期时间", desc)}
	return
}

//小时榜
func dateTimeLiveHour(key string) (res map[string][]string, dateHourList map[string][]string) {
	res = map[string][]string{"date": {}, "hour": {}}
	now := time.Now()
	dateList := getDateList(7, now)
	var currentHourList, commonHourList []string
	getHourList := func(start int) (hourList []string) {
		hourList = []string{}
		for i := start; i >= 0; i-- {
			hourString := fmt.Sprintf("%02d:00", i)
			hourList = append(hourList, hourString)
		}
		return
	}
	res["date"] = dateList
	startCurrentHour := now.Hour()
	for i := startCurrentHour; i >= 0; i-- {
		if checkIsExistHour(key, i) {
			startCurrentHour = i
			break
		} else if i == 0 {
			startCurrentHour = -1
			break
		}
	}
	currentHourList = getHourList(startCurrentHour)
	commonHourList = getHourList(23)
	dateHourList = map[string][]string{}
	for k, v := range dateList {
		if k == 0 {
			if len(currentHourList) > 0 {
				dateHourList[v] = currentHourList
			}
		} else {
			dateHourList[v] = commonHourList
		}
	}
	return
}

//根据是否有数据判断日期列表
func getCheckDateList(key string) (res map[string][]string) {
	res = map[string][]string{"date": {}, "hour": {}}
	now := time.Now()
	isExist := checkIsExistDate(key)
	beforeInt := -2
	if isExist {
		beforeInt = -1
	}
	startSate := now.AddDate(0, 0, beforeInt)
	res["date"] = getDateList(30, startSate)
	return
}

//周榜日期列表获取
func getWeekList(key string) (res []map[string]string) {
	//这里仿照前段，只给三个切片
	now := time.Now()
	num := 3
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	startDateTime := time.Now().AddDate(0, 0, (offset - 1))
	isExist := checkIsExistWeek(key)
	if !isExist {
		startDateTime = startDateTime.AddDate(0, 0, -7)
	}
	dateSelectList := []map[string]string{}
	for i := 0; i < num; i++ {
		rightDate := startDateTime.AddDate(0, 0, -i*6)
		leftDate := startDateTime.AddDate(0, 0, -(i+1)*6)
		dateString := fmt.Sprintf("%s-%s", leftDate.Format("01/02"), rightDate.Format("01/02"))
		dateSelectList = append(dateSelectList, map[string]string{"week_string": dateString, "request_date": rightDate.AddDate(0, 0, 1).Format("2006-01-02")})
		startDateTime = startDateTime.AddDate(0, 0, -1)
	}
	res = dateSelectList
	return
}

//电商直播达人分享榜
func getWeekListLiveShare() (res []map[string]string) {
	//这里仿照前段，只给三个切片
	now := time.Now()
	num := 3
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	startDateTime := time.Now().AddDate(0, 0, (offset - 1))
	dateSelectList := []map[string]string{}
	for i := 0; i < num; i++ {
		rightDate := startDateTime.AddDate(0, 0, -i*6)
		leftDate := startDateTime.AddDate(0, 0, -(i+1)*6)
		dateString := fmt.Sprintf("%s-%s", leftDate.Format("01/02"), rightDate.Format("01/02"))
		requestDateString := fmt.Sprintf("%s/%s", leftDate.Format("2006-01-02"), rightDate.Format("2006-01-02"))
		dateSelectList = append(dateSelectList, map[string]string{"week_string": dateString, "request_date": requestDateString})
		startDateTime = startDateTime.AddDate(0, 0, -1)
	}
	res = dateSelectList
	return
}

//月榜列表获取
func getMonthList() (res []string) {
	//这里仿照前段，只给三个切片
	//num := 3
	//startDateTime := time.Now().AddDate(0, -1, 0)
	dateSelectList := []string{}
	//for i := 0; i < num; i++ {
	//	monthDate := startDateTime.AddDate(0, -i, 0)
	//	stopDate, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-09-01 00:00:00", time.Local)
	//	if stopDate.Before(monthDate) {
	//		dateString := monthDate.Format("2006-01")
	//		dateSelectList = append(dateSelectList, dateString)
	//	}
	//}
	res = dateSelectList
	return
}

//判断缓存是否存在
func checkcachKey(cachKey string) (isExist bool) {
	result := global.Cache.Get(cachKey)
	if result != "" {
		if result == "1" {
			isExist = true
		} else {
			isExist = false
		}
	} else {
		isExist = false
	}
	return
}

//检测该日榜周榜榜单是否已经存在了数据
func checkIsExistDate(key string) (isExist bool) {
	cachKey := cache.GetCacheKey(cache.DyRankCache, time.Now().Format("20060102"), key)
	isExist = checkcachKey(cachKey)
	if isExist == false {
		pathInfo := getRoute(key)
		isExist = requestRank(pathInfo)
		if isExist {
			//有数据情况，缓存设置到今天结束
			now := time.Now()
			dateString := fmt.Sprintf("%s 00:00:00", now.AddDate(0, 0, 1).Format("2006-01-02"))
			stopTime, _ := time.ParseInLocation("2006-01-02 15:04:05", dateString, time.Local)
			seconds := stopTime.Unix() - now.Unix()
			secondsDuration := time.Duration(seconds)
			global.Cache.Set(cachKey, "1", secondsDuration)
		}
	}
	return
}

//检测周榜榜单是否已经存在了数据
func checkIsExistWeek(key string) (isExist bool) {
	cachKey := cache.GetCacheKey(cache.DyRankCache, "week", key)
	isExist = checkcachKey(cachKey)
	if isExist == false {
		if key != "live_share" {
			key = fmt.Sprintf("%s_week", key)
		}
		pathInfo := getRoute(key)
		isExist = requestRank(pathInfo)
		if isExist {
			//有数据情况，缓存到本周结束
			now := time.Now()
			dateString := fmt.Sprintf("%s 00:00:00", now.AddDate(0, 0, (8-int(now.Weekday()))).Format("2006-01-02"))
			stopTime, _ := time.ParseInLocation("2006-01-02 15:04:05", dateString, time.Local)
			seconds := stopTime.Unix() - now.Unix()
			secondsDuration := time.Duration(seconds)
			global.Cache.Set(cachKey, "1", secondsDuration)
		}
	}
	return
}

//检测该小时榜榜单是否已经存在了数据
func checkIsExistHour(key string, currentHour int) (isExist bool) {
	cachKey := cache.GetCacheKey(cache.DyRankCache, utils.ToString(currentHour), key)
	isExist = checkcachKey(cachKey)
	if isExist == false {
		checkNowTime, _ := time.ParseInLocation("20060102150405", fmt.Sprintf("%s%02d0000", time.Now().Format("20060102"), currentHour), time.Local)
		pathInfo := getCheckRoute(key, checkNowTime)
		isExist = requestRank(pathInfo)
		if isExist {
			//有数据情况，缓存设置到今天结束
			dateString := fmt.Sprintf("%s %s:59:59", checkNowTime.Format("2006-01-02"), strconv.Itoa(checkNowTime.Hour()))
			stopTime, _ := time.ParseInLocation("2006-01-02 15:04:05", dateString, time.Local)
			seconds := stopTime.Unix() - checkNowTime.Unix()
			secondsDuration := time.Duration(seconds)
			global.Cache.Set(cachKey, "1", secondsDuration)
		}
	}
	return
}

//获取日期列表
func getDateList(daysCount int, startTime time.Time) (list []string) {
	list = []string{}
	for i := 0; i < daysCount; i++ {
		date := startTime.AddDate(0, 0, -i).Format("2006-01-02")
		list = append(list, date)
	}
	return
}
