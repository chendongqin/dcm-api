package business

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	aliLog "dongchamao/services/ali_log"
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"math"
	"strings"
	"time"
)

type SafeBusiness struct {
}

func NewSafeBusiness() *SafeBusiness {
	return new(SafeBusiness)
}

type AnaListStruct struct {
	SqlString string `json:"sql_string"` //查询语句
	Desc      string `json:"desc"`       //描述
	Point     int64  `json:"point"`      //限制数
	TimeEndS  int64  `json:"time_end_s"` //timend 需要减去的秒
}

//轮询url结构体
type ExtraUrl struct {
	Url    string `json:"url"`    //url
	Second int    `json:"second"` //频率-秒
}

//日志查询到的统计字段
type LogField struct {
	Uid int64  `json:"uid"`
	Pv  int64  `json:"pv"`
	Url string `json:"url"`
	Uri string `json:"uri"`
}

//配置
type SaleBusinessConfigStruct struct {
	RequestStages        map[string]AnaListStruct `json:"request_stages"`          //不同时间段的请求
	CommonUrlPoint       map[string]int64         `json:"common_url_point"`        //普通url白天点击量监控pv限制量
	CommonUrlNightPoints map[string]int64         `json:"common_url_night_points"` //普通url夜间1-6点点击量监控pv限制量
	ExtraUrls            []ExtraUrl               `json:"extra_url"`               //轮询url
	WhiteLists           []int64                  `json:"white_list"`              //白名单列表
}

//加速请求结构体
type SpeedUrlStruct struct {
	Uri        string `json:"uri"`
	LevelTimes int64  `json:"level_times"`
	SpeedDays  int    `json:"speed_days"`
	Desc       string `json:"desc"`
}

var SaleBusinessConfig = SaleBusinessConfigStruct{}       //通用url监控的配置
var SpeedSaleBusinessConfig = map[string]SpeedUrlStruct{} //加速的配置

// InitCommonUrlConfig 初始化通用url的配置
func (s *SafeBusiness) InitCommonUrlConfig() {
	SaleBusinessConfig = SaleBusinessConfigStruct{
		RequestStages: map[string]AnaListStruct{
			"min_5":      {Desc: "五分钟内统计数据", TimeEndS: 300},
			"min_30":     {Desc: "半小时内统计数据", TimeEndS: 1800},
			"min_60":     {Desc: "一小时内统计数据", TimeEndS: 3600},
			"min_url_60": {Desc: "一小时内按URL统计数据", TimeEndS: 3600},
		},
		CommonUrlPoint: map[string]int64{
			"min_5":      200 * 1.5,
			"min_30":     600 * 1.5,
			"min_60":     1000 * 1.5,
			"min_url_60": 450 * 1.5,
		},
		CommonUrlNightPoints: map[string]int64{
			"min_5":      100 * 1.5,
			"min_30":     300 * 1.5,
			"min_60":     500 * 1.5,
			"min_url_60": 250 * 1.5,
		},
		ExtraUrls: []ExtraUrl{
			{"/v1/dy/live/fans/data", 30},
			{"/v1/dy/living/sale", 30},
			{"/v1/dy/living/watch/chart", 30},
			{"/v1/dy/living/product/", 10},
			{"/v1/dy/living/message", 0},
			{"/v1/wechat/check", 0},
		},
		WhiteLists: []int64{1},
	}
}

// InitSpeedConfig 初始化直播加速的配置
func (s *SafeBusiness) InitSpeedConfig(days int) {
	SpeedSaleBusinessConfig = map[string]SpeedUrlStruct{
		"speed_author":  {Uri: "/v1/dy/author/info/", Desc: "达人"},
		"speed_live":    {Uri: "/v1/dy/live/info/", Desc: "直播"},
		"speed_product": {Uri: "/v1/dy/product/base/", Desc: "商品"},
	}
	var row SpeedUrlStruct
	if days == 0 {
		row = SpeedSaleBusinessConfig["speed_author"]
		row.SpeedDays = 1
		row.LevelTimes = 2
		SpeedSaleBusinessConfig["speed_author"] = row

		row = SpeedSaleBusinessConfig["speed_live"]
		row.SpeedDays = 1
		row.LevelTimes = 2
		SpeedSaleBusinessConfig["speed_live"] = row
		row = SpeedSaleBusinessConfig["speed_product"]

		row.SpeedDays = 1
		row.LevelTimes = 2
		SpeedSaleBusinessConfig["speed_product"] = row
	}
	if days == 7 {
		row = SpeedSaleBusinessConfig["speed_author"]
		row.SpeedDays = 30
		row.LevelTimes = 3
		SpeedSaleBusinessConfig["speed_author"] = row

		row = SpeedSaleBusinessConfig["speed_live"]
		row.SpeedDays = 30
		row.LevelTimes = 3
		SpeedSaleBusinessConfig["speed_live"] = row

		row = SpeedSaleBusinessConfig["speed_product"]
		row.SpeedDays = 30
		row.LevelTimes = 3
		SpeedSaleBusinessConfig["speed_product"] = row
	}
	return
}

// CommonAnalyseLogs 通用url日志分析
func (s *SafeBusiness) CommonAnalyseLogs() (logList map[string][]LogField) {
	s.InitCommonUrlConfig()
	timeEnd := utils.Time() - 120
	commonUrlMapList := s.commonMap(timeEnd)
	logList = make(map[string][]LogField)
	for k, v := range commonUrlMapList {
		logListRes, _ := s.reqestAliLog(timeEnd, v)
		//logList = append(logList, logListRes)
		logList[k] = logListRes
	}
	return
}

func (s *SafeBusiness) SpeedLogs() (logList [][]LogField) {

	return
}

//筛选需要加速的数据
func (s *SafeBusiness) SpeedFilterLog(key string, days int, inputEndTime int64, page int) (result interface{}) {
	s.InitSpeedConfig(days)
	mapList := s.combineSpeedMap(days)
	timeEnd := inputEndTime
	todayStartStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 00:00:00", time.Now().Format("2006-01-02")), time.Local)
	if inputEndTime == 0 || inputEndTime < todayStartStamp.Unix() {
		timeEnd = utils.Time() - 120
	}
	mapRow := mapList[key]
	uri := SpeedSaleBusinessConfig[key].Uri
	logListRes, _ := s.reqestAliLog(timeEnd, mapRow)
	list := []string{}
	for _, v := range logListRes {
		url := v.Url
		if url != "null" {
			parami := strings.Replace(url, uri, "", 1)
			parami = IdDecrypt(parami)
			list = append(list, parami)
		}
	}
	total := len(list)
	totalPage, currentPage := 1, 1
	pageSize := len(list)
	if key == "speed_live" {
		currentPage, pageSize = page, 30
		totalPage = int(math.Ceil(float64(total) / float64(pageSize)))
		sliceSplitLeft := (page - 1) * pageSize
		sliceSplitRight := page * pageSize
		if page >= totalPage && totalPage > 0 {
			currentPage = totalPage
			sliceSplitLeft = (totalPage - 1) * pageSize
			sliceSplitRight = total
		}
		llen := len(list)
		if llen > 0 {
			roomInfos, _ := hbase.GetLiveInfoByIds(list[sliceSplitLeft:sliceSplitRight])
			newList := []string{}
			for _, v := range roomInfos {
				if !utils.InArrayString(v.User.ID, newList) {
					newList = append(newList, v.User.ID)
				}
			}
			list = newList
		}

	}
	result = map[string]interface{}{
		"list":         list,
		"total":        total,
		"current_page": currentPage,
		"total_page":   totalPage,
		"page_size":    pageSize,
	}
	return
}

//组装加速的map
func (s *SafeBusiness) combineSpeedMap(days int) (mapList map[string]AnaListStruct) {
	now := time.Now()
	date := now.AddDate(0, 0, -days).Format("2006-01-02")
	date = fmt.Sprintf("%s 00:00:00", date)
	formatTime, _ := time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)
	beforeStamp := formatTime.Unix()
	nowStamp := time.Now().Unix()
	TimeEndS := nowStamp - beforeStamp
	mapList = make(map[string]AnaListStruct)
	for k, v := range SpeedSaleBusinessConfig {
		row := AnaListStruct{
			//SqlString string `json:"sql_string"`      //查询语句
			Desc:     v.Desc,
			Point:    v.LevelTimes,
			TimeEndS: TimeEndS,
		}
		row.SqlString = s.combineSpeedSqlString(v.Uri, row)
		mapList[k] = row
	}
	return
}

//组装加速的sql字符串
func (s *SafeBusiness) combineSpeedSqlString(uri string, row AnaListStruct) (sqlString string) {
	runmode := "prod"
	sqlString = fmt.Sprintf("env:%s and log_type:\"Format\" and uri:\"%s\"  | select url,COUNT(*) as pv where uid <> 0  group by url HAVING pv>=%d order by pv desc limit 1000000", runmode, uri, row.Point)
	return
}

//通用url的map
func (s *SafeBusiness) commonMap(timeEnd int64) (mapList map[string]AnaListStruct) {
	mapList = SaleBusinessConfig.RequestStages
	hour := utils.ParseInt(time.Unix(timeEnd, 0).Format("15"), 0)
	pointMap := SaleBusinessConfig.CommonUrlNightPoints
	if hour >= 1 && hour <= 6 {
		pointMap = SaleBusinessConfig.CommonUrlPoint
	}
	for k, v := range mapList {
		row := v
		row.Point = pointMap[k]
		row.SqlString = s.commonUrlCondition(k, row)
		mapList[k] = row
	}
	return
}

//排除url条件拼接字符串语句
func (s *SafeBusiness) commonUrlCondition(k string, row AnaListStruct) (sqlString string) {
	whereNotString := ""
	urlSlices := SaleBusinessConfig.ExtraUrls
	for k, v := range urlSlices {
		isplit := ""
		if k > 0 {
			isplit = ","
		}
		whereNotString = fmt.Sprintf("%s%s'%s'", whereNotString, isplit, v.Url)
	}
	if whereNotString != "" {
		whereNotString = fmt.Sprintf("where uri not in (%s)", whereNotString)
	}
	//runmode := global.Cfg.String("runmode")
	runmode := "prod"
	andString := fmt.Sprintf("env:%s and log_type:\"Format\" ", runmode)
	if k != "min_url_60" {
		sqlString = fmt.Sprintf("%s  | select uid,COUNT(*) as pv %s group by uid HAVING pv>=%d order by pv desc", andString, whereNotString, row.Point)
	} else {
		sqlString = fmt.Sprintf("%s  | select uri,COUNT(*) as pv %s group by uri HAVING pv>=%d order by pv desc", andString, whereNotString, row.Point)
	}
	return
}

/**请求阿里云**/
func (s *SafeBusiness) reqestAliLog(timeEnd int64, row AnaListStruct) (logList []LogField, comErr global.CommonError) {
	Client := sls.CreateNormalInterface(aliLog.Endpoint, global.Cfg.String("ali_accessKey"), global.Cfg.String("ali_secret"), "")
	res, err := Client.GetLogs("dongchamao-api-log", "dongchamao-log-api-history", "", timeEnd-row.TimeEndS, timeEnd, row.SqlString, 100, 0, true)
	if err != nil {
		comErr = global.NewMsgError(fmt.Sprintf("获取%s失败", row.Desc))
		return
	}
	for _, v := range res.Logs {
		uid := int64(0)
		if u, ok := v["uid"]; ok {
			uid = utils.ParseInt64String(u)
			if utils.InArrayInt64(uid, SaleBusinessConfig.WhiteLists) {
				continue
			}
		}
		pv := utils.ParseInt64String(v["pv"])
		uri := ""
		if a, ok := v["uri"]; ok {
			uri = a
		}
		url := ""
		if a, ok := v["url"]; ok {
			url = a
		}
		logList = append(logList, LogField{uid, pv, url, uri})
	}
	return
}
