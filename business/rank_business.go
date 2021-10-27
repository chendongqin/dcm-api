package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type RankBusiness struct {
}
type PathDesc struct {
	Path string `json:"path"`
	Desc string `json:"desc"`
}
type MonitorParam struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type goodsCat struct {
	Name string `json:"name"`
}

func NewRankBusiness() *RankBusiness {
	return new(RankBusiness)
}

/**
**遍历执行监控
**hour:小时字符串
 */
func (t *RankBusiness) LoopCheck(hour string) {
	monitorBusiness := NewMonitorBusiness()
	var monitorEvents []string
	taskList := t.getRow(hour)
	length := len(taskList)
	//pathInfo := getRoute(v)
	//monitorEvents = append(monitorEvents, pathInfo.Desc)
	//checkLiveHotRank(pathInfo)
	desc := ""
	param := MonitorParam{
		Type: "monitor",
		Time: "",
	}
	if length > 0 {
		for _, v := range taskList {
			desc, _ = t.monitorKey(v, param)
			if desc != "" {
				monitorEvents = append(monitorEvents, desc)
			}
		}
		monitorBusiness.SendErr("直播监控", strings.Join(monitorEvents, ","))
	}
}

//根据key执行对应的榜单
func (t *RankBusiness) monitorKey(key string, param MonitorParam) (desc string, res bool) {
	desc = t.getRoute(key)
	//rankMap := map[string]interface{}{
	//	"live_hour":t.monitorLiveHour,
	//	"live_top":t.monitorLiveTop,
	//	"product_sale":t.monitorProductSale,
	//	"product_share":t.monitorProductShare,
	//	"product_live_sale":t.monitorProductLiveSale,
	//	"product":t.monitorProduct,
	//	"product_week":t.monitorProductWeek,
	//	"author_follower_inc":t.monitorAuthorFollowerInc,
	//	"author_goods":t.monitorAuthorGoods,
	//	"live_share":t.monitorAuthorGoods,
	//	"author_aweme_rank":t.monitorAuthorAwemeRank,
	//	"author_aweme_live":t.monitorAuthorAwemeLive,
	//}
	//row := rankMap[key]
	if desc != "" {
		switch key {
		case "live_hour":
			res = t.monitorLiveHour(desc, param)
		case "live_top":
			res = t.monitorLiveTop(desc, param)
		case "product_sale":
			res = t.monitorProductSale(desc, param)
		case "product_share":
			res = t.monitorProductShare(desc, param)
		case "product_live_sale":
			res = t.monitorProductLiveSale(desc, param)
		case "product_live_sale_week":
			res = t.monitorProductLiveSaleWeek(desc, param)
		case "product_live_sale_month":
			res = t.monitorProductLiveSaleMonth(desc, param)
		case "product":
			res = t.monitorProduct(desc, param)
		case "product_week":
			res = t.monitorProductWeek(desc, param)
		case "product_month":
			res = t.monitorProductMonth(desc, param)
		case "author_follower_inc":
			res = t.monitorAuthorFollowerInc(desc, param)
		case "author_goods":
			res = t.monitorAuthorGoods(desc, param)
		case "live_share":
			res = t.monitorLiveShare(desc, param)
		case "author_aweme_rank":
			res = t.monitorAuthorAwemeRank(desc, param)
		case "author_aweme_live":
			res = t.monitorAuthorAwemeLive(desc, param)
		}
	}

	return
}

//直播小时榜监控
func (t *RankBusiness) monitorLiveHour(desc string, param MonitorParam) (res bool) {

	now := time.Now()
	hour := now.Hour()
	toDate := now.Format("2006-01-02")
	if param.Type == "get_time" {
		hour = utils.ToInt(param.Time)
	}
	hourString := strconv.Itoa(hour)
	path := fmt.Sprintf("/v1/dy/rank/live/hour/%s/%s:00", toDate, hourString)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	return
}

//直播热榜监控
func (t *RankBusiness) monitorLiveTop(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	hour := now.Hour()
	toDate := now.Format("2006-01-02")
	if param.Type == "get_time" {
		hour = utils.ToInt(param.Time)
	}
	hourString := strconv.Itoa(hour)
	path := fmt.Sprintf("/v1/dy/rank/live/top/%s/%s:00", toDate, hourString)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	return
}

//抖音销量榜监控
func (t *RankBusiness) monitorProductSale(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	reqDate := now.AddDate(0, 0, -1).Format("2006-01-02")
	if param.Type == "get_time" {
		reqDate = param.Time
	}
	path := fmt.Sprintf("/v1/dy/rank/product/sale/%s?data_type=1&first_cate=&order_by=desc&sort=order_count&page=1&page_size=50", reqDate)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	return
}

//抖音热推榜监控
func (t *RankBusiness) monitorProductShare(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	reqDate := now.AddDate(0, 0, -1).Format("2006-01-02")
	if param.Type == "get_time" {
		reqDate = param.Time
	}
	path := fmt.Sprintf("/v1/dy/rank/product/share/%s?first_cate=&data_type=1&order_by=desc&sort=share_count&page=1&page_size=50", reqDate)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	return
}

//直播商品榜监控
func (t *RankBusiness) monitorProductLiveSale(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	reqDate := now.AddDate(0, 0, -1).Format("2006-01-02")
	if param.Type == "get_time" {
		reqDate = param.Time
	}
	path := fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=1&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", reqDate)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	if res {
		pathString := "/v1/dy/rank/product/live/sale/%s?data_type=1&first_cate=%s&order_by=desc&sort=sales&page=1&page_size=50"
		t.monitorCategory(desc, pathString, reqDate, param.Type)
	}
	return
}

//直播商品榜周榜监控
func (t *RankBusiness) monitorProductLiveSaleWeek(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartTime := now.AddDate(0, 0, offset)
	weekStartDate := weekStartTime.AddDate(0, 0, -7).Format("2006-01-02")
	path := fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=2&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", weekStartDate)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	if param.Type == "get_time" {
		return
	}
	if res {
		pathString := "/v1/dy/rank/product/live/sale/%s?data_type=2&first_cate=%s&order_by=desc&sort=sales&page=1&page_size=50"
		t.monitorCategory(desc, pathString, weekStartDate, param.Type)
	}
	return
}

//直播商品榜周榜监控
func (t *RankBusiness) monitorProductLiveSaleMonth(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	startDate := fmt.Sprintf("%s-01", now.AddDate(0, -1, 0).Format("2006-01"))
	path := fmt.Sprintf("/v1/dy/rank/product/live/sale/%s?data_type=3&first_cate=&order_by=desc&sort=sales&page=1&page_size=50", startDate)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	if param.Type == "get_time" {
		return
	}
	if res {
		pathString := "/v1/dy/rank/product/live/sale/%s?data_type=3&first_cate=%s&order_by=desc&sort=sales&page=1&page_size=50"
		t.monitorCategory(desc, pathString, startDate, param.Type)
	}
	return
}

//视频商品榜监控
func (t *RankBusiness) monitorProduct(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	reqDate := now.AddDate(0, 0, -1).Format("2006-01-02")
	path := fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=1&page=1&page_size=50", reqDate)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	if param.Type == "get_time" {
		return
	}
	if res {
		pathString := "/v1/dy/rank/product/%s?first_cate=%s&order_by=desc&sort=sales&data_type=1&page=1&page_size=50"
		t.monitorCategory(desc, pathString, reqDate, param.Type)
	}
	return
}

//视频商品榜周榜监控
func (t *RankBusiness) monitorProductWeek(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartTime := now.AddDate(0, 0, offset)
	weekStartDate := weekStartTime.AddDate(0, 0, -7).Format("2006-01-02")
	path := fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=2&page=1&page_size=50", weekStartDate)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	if param.Type == "get_time" {
		return
	}
	if res {
		pathString := "/v1/dy/rank/product/%s?first_cate=%s&order_by=desc&sort=sales&data_type=2&page=1&page_size=50"
		t.monitorCategory(desc, pathString, weekStartDate, param.Type)
	}
	return
}

//视频商品榜月榜监控
func (t *RankBusiness) monitorProductMonth(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	NowYearMonth := now.Format("2006-01")
	startDate := fmt.Sprintf("%s-01", now.AddDate(0, -1, 0).Format("2006-01"))
	if NowYearMonth == "2021-10" {
		startDate = "2021-10-01"
	}
	path := fmt.Sprintf("/v1/dy/rank/product/%s?first_cate=&order_by=desc&sort=sales&data_type=3&page=1&page_size=50", startDate)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	if param.Type == "get_time" {
		return
	}
	if res {
		pathString := "/v1/dy/rank/product/%s?first_cate=%s&order_by=desc&sort=sales&data_type=2&page=1&page_size=50"
		t.monitorCategory(desc, pathString, startDate, param.Type)
	}
	return
}

//达人涨粉榜
func (t *RankBusiness) monitorAuthorFollowerInc(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	reqDate := now.AddDate(0, 0, -1).Format("2006-01-02")
	if param.Type == "get_time" {
		reqDate = param.Time
	}
	path := fmt.Sprintf("/v1/dy/rank/author/follower/inc/%s?tags=&province=&page=1&is_delivery=0&page_size=50&order_by=desc&sort=inc_follower_count", reqDate)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	return
}

//达人带货榜
func (t *RankBusiness) monitorAuthorGoods(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	reqDate := now.AddDate(0, 0, -1).Format("2006-01-02")
	if param.Type == "get_time" {
		reqDate = param.Time
	}
	path := fmt.Sprintf("/v1/dy/rank/author/goods/%s?date_type=1&tags=&verified=0&page=1&page_size=50&sort=sum_gmv&order_by=desc", reqDate)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	return
}

//电商达人分享榜
func (t *RankBusiness) monitorLiveShare(desc string, param MonitorParam) (res bool) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartTime := now.AddDate(0, 0, offset)
	weekStartDate := weekStartTime.AddDate(0, 0, -7).Format("2006-01-02")
	weekEndDate := weekStartTime.AddDate(0, 0, -1).Format("2006-01-02")
	reqWeekString := fmt.Sprintf("%s/%s", weekStartDate, weekEndDate)
	if param.Type == "get_time" {
		reqWeekString = param.Time
	}
	path := fmt.Sprintf("/v1/dy/rank/live/share/%s", reqWeekString)
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	return
}

//抖音短视频达人热榜
func (t *RankBusiness) monitorAuthorAwemeRank(desc string, param MonitorParam) (res bool) {
	desc = "抖音短视频达人热榜"
	path := "/v1/dy/rank/author/aweme?rank_type=达人指数榜&category=全部"
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	return
}

//抖音直播主播热榜
func (t *RankBusiness) monitorAuthorAwemeLive(desc string, param MonitorParam) (res bool) {
	path := "/v1/dy/rank/author/live?rank_type=达人指数榜"
	row := PathDesc{
		Path: path,
		Desc: desc,
	}
	res = t.checkLiveHotRank(row, param.Type)
	return
}

//分类榜的监控
func (t *RankBusiness) monitorCategory(desc, pathString, reqDate, mointType string) {
	catList, err := t.requestProductCat()
	if err != nil {
		return
	}
	for _, v := range catList {
		path := fmt.Sprintf(pathString, reqDate, v.Name)
		row := PathDesc{
			Path: path,
			Desc: desc,
		}
		t.checkLiveHotRank(row, mointType)
	}
	return
}

/**
**name:榜单名称
**hour：小时
 */
func (t *RankBusiness) getRoute(key string) (desc string) {

	var routeMap = map[string]string{
		"live_hour":               "直播小时榜",
		"live_top":                "直播热榜",
		"product_sale":            "抖音销量榜",
		"product_share":           "抖音热推榜",
		"product_live_sale":       "直播商品榜",
		"product_live_sale_week":  "直播商品榜-周榜",
		"product_live_sale_month": "直播商品榜-月榜",
		"product":                 "视频商品榜",
		"product_week":            "视频商品榜-周榜",
		"product_month":           "视频商品榜-月榜",
		"author_follower_inc":     "达人涨粉榜",
		"author_goods":            "达人带货榜",
		"live_share":              "电商直播达人分享榜",
		"author_aweme_rank":       "抖音短视频达人热榜",
		"author_aweme_live":       "抖音直播主播热榜",
	}
	desc = ""
	if _, ok := routeMap[key]; !ok {
		return
	}
	desc = routeMap[key]
	return
}
func (t *RankBusiness) getRow(hour string) (taskList []string) {
	hourGroup := t.getHourGroup()
	taskList = hourGroup[hour]
	return
}
func (t *RankBusiness) getHourGroup() (hourGroup map[string][]string) {
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
func (t *RankBusiness) checkLiveHotRank(pathInfo PathDesc, monitorType string) (checkRes bool) {
	checkRes = t.requestRank(pathInfo)
	if !checkRes && monitorType == "monitor" {
		NewMonitorBusiness().SendTemplateMessage("S", pathInfo.Desc, fmt.Sprintf("%s挂了，请求地址：%s", pathInfo.Desc, pathInfo.Path))
	}
	return
}

//请求对应的榜单
func (t *RankBusiness) requestRank(pathInfo PathDesc) (checkRes bool) {
	interBusiness := NewInterBusiness()
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

//获取商品分类
func (t *RankBusiness) requestProductCat() (catList []goodsCat, comErr global.CommonError) {
	interBusiness := NewInterBusiness()
	api := interBusiness.BuildURL("/v1/dy/product/cate")
	var res interface{}
	res, comErr = interBusiness.HttpGet(api)
	if comErr != nil {
		return
	}
	utils.MapToStruct(res, &catList)
	return
}

//根据key返回对应榜单需要展示的日期时间
func (t *RankBusiness) SwitchTopDateTime(key string) (main map[string][]string, hourList map[string][]string, weekList []map[string]string, monthList []string, comErr global.CommonError) {
	if key == "author_aweme_rank" || key == "author_aweme_live" {
		comErr = global.NewMsgError("传入参数错误，不存在的key")
		return
	}
	desc := t.getRoute(key)
	if desc == "" {
		comErr = global.NewMsgError("传入参数错误，不存在的key")
		return
	}
	hourList = map[string][]string{}
	weekList = []map[string]string{}
	monthList = []string{}
	main = make(map[string][]string)
	switch key {
	case "live_hour", "live_top":
		main, hourList = t.getHourDateList(key)
	case "product_sale", "product_share", "author_follower_inc", "author_goods":
		main = t.getCheckDateList(key)
	case "product_live_sale", "product":
		main = t.getCheckDateList(key)
		weekList = t.getWeekList(key)
		monthList = t.getMonthList(key)
	case "live_share":
		main = map[string][]string{"date": {}, "hour": {}}
		weekList = t.getWeekListLiveShare(key)
	}
	main["desc"] = []string{fmt.Sprintf("%s的日期时间", desc)}
	return
}

//小时榜
func (t *RankBusiness) getHourDateList(key string) (res map[string][]string, dateHourList map[string][]string) {
	res = map[string][]string{"date": {}, "hour": {}}
	now := time.Now()
	dateList := t.getDateList(7, now)
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
		if t.checkIsExistHour(key, i) {
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

//获取榜单日期列表
func (t *RankBusiness) getCheckDateList(key string) (res map[string][]string) {
	res = map[string][]string{"date": {}, "hour": {}}
	now := time.Now()
	isExist := t.checkIsExistDate(key)
	beforeInt := -2
	if isExist {
		beforeInt = -1
	}
	startSate := now.AddDate(0, 0, beforeInt)
	res["date"] = t.getDateList(30, startSate)
	return
}

//周榜日期列表获取
func (t *RankBusiness) getWeekList(key string) (res []map[string]string) {
	//这里仿照前段，只给三个切片
	now := time.Now()
	num := 3
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	startDateTime := time.Now().AddDate(0, 0, (offset - 1))
	isExist := t.checkIsExistWeek(key)
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
func (t *RankBusiness) getWeekListLiveShare(key string) (res []map[string]string) {
	//这里仿照前段，只给三个切片
	now := time.Now()
	num := 3
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	startDateTime := time.Now().AddDate(0, 0, (offset - 1))
	isExist := t.checkIsExistWeek(key)
	if !isExist {
		startDateTime = startDateTime.AddDate(0, 0, -7)
	}
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
func (t *RankBusiness) getMonthList(key string) (res []string) {
	num := 3
	now := time.Now()
	beforeInt := -1
	if key == "product" {
		if now.Format("2006-01") == "2021-10" {
			beforeInt = 0
		}
	}
	if !t.checkIsExistMonth(key) {
		beforeInt = beforeInt - 1
	}
	startDateTime := now.AddDate(0, beforeInt, 0)
	dateSelectList := []string{}
	for i := 0; i < num; i++ {
		monthDate := startDateTime.AddDate(0, -i, 0)
		stopDate, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-09-01 00:00:00", time.Local)
		if key == "product" {
			stopDate, _ = time.ParseInLocation("2006-01-02 15:04:05", "2021-10-01 00:00:00", time.Local)
		}
		if stopDate.Before(monthDate) {
			dateString := monthDate.Format("2006-01")
			dateSelectList = append(dateSelectList, dateString)
		}
	}
	res = dateSelectList
	return
}

//获取日期列表
func (t *RankBusiness) getDateList(daysCount int, startTime time.Time) (list []string) {
	list = []string{}
	for i := 0; i < daysCount; i++ {
		date := startTime.AddDate(0, 0, -i).Format("2006-01-02")
		list = append(list, date)
	}
	return
}

//检测该日榜周榜榜单是否已经存在了数据
func (t *RankBusiness) checkIsExistDate(key string) (isExist bool) {
	cachKey := cache.GetCacheKey(cache.DyRankCache, time.Now().Format("20060102"), key)
	isExist = t.checkcachKey(cachKey)
	if isExist == false {
		param := MonitorParam{
			Type: "get_time",
			Time: "",
		}
		_, isExist = t.monitorKey(key, param)
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
func (t *RankBusiness) checkIsExistWeek(key string) (isExist bool) {
	cachKey := cache.GetCacheKey(cache.DyRankCache, "week", key)
	isExist = t.checkcachKey(cachKey)
	if isExist == false {
		if key != "live_share" {
			key = fmt.Sprintf("%s_week", key)
		}
		param := MonitorParam{
			Type: "get_time",
			Time: "",
		}
		_, isExist = t.monitorKey(key, param)
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

//检测月榜榜单是否已经存在了数据
func (t *RankBusiness) checkIsExistMonth(key string) (isExist bool) {
	now := time.Now()
	Month := now.Format("2006-01")
	cachKey := cache.GetCacheKey(cache.DyRankCache, Month, key)
	isExist = t.checkcachKey(cachKey)
	if isExist == false {
		key = fmt.Sprintf("%s_month", key)

		param := MonitorParam{
			Type: "get_time",
			Time: "",
		}

		_, isExist = t.monitorKey(key, param)
		if isExist {
			//有数据情况，缓存到本月结束
			dateString := fmt.Sprintf("%s-01 00:00:00", now.AddDate(0, 1, 0).Format("2006-01"))
			stopTime, _ := time.ParseInLocation("2006-01-02 15:04:05", dateString, time.Local)
			seconds := stopTime.Unix() - now.Unix()
			secondsDuration := time.Duration(seconds)
			global.Cache.Set(cachKey, "1", secondsDuration)
		}
	}
	return
}

//检测该小时榜榜单是否已经存在了数据
func (t *RankBusiness) checkIsExistHour(key string, currentHour int) (isExist bool) {
	cachKey := cache.GetCacheKey(cache.DyRankCache, utils.ToString(currentHour), key)
	isExist = t.checkcachKey(cachKey)
	if isExist == false {
		checkNowTime, _ := time.ParseInLocation("20060102150405", fmt.Sprintf("%s%02d0000", time.Now().Format("20060102"), currentHour), time.Local)
		param := MonitorParam{
			Type: "get_time",
			Time: utils.ToString(checkNowTime),
		}
		_, isExist = t.monitorKey(key, param)
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

//判断缓存是否存在
func (t *RankBusiness) checkcachKey(cachKey string) (isExist bool) {
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
