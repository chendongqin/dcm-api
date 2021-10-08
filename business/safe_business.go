package business

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	aliLog "dongchamao/services/ali_log"
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"time"
)

type SafeBusiness struct {
}

func (s *SafeBusiness) NewSafeBusiness() *SafeBusiness {
	return new(SafeBusiness)
}

func (s *SafeBusiness) AnalyseLogs() (comErr global.CommonError) {
	//t1 := time.Now().Year()  //年
	//t2 := time.Now().Month() //月
	//t3 := time.Now().Day()   //日
	//end_time := time.Date(t1, t2, t3, 0, 0, 0, 0, time.Local).Unix() - 1
	//start_time := end_time - 86400 + 1
	////分析前一天用户数据，并记入数据库
	//go logic.NewHandleLogic().AnalyseUidLogs(start_time, end_time)
	////分析前一天IP数据前50，取他们的运营商，并记入数据库
	//go logic.NewHandleLogic().AnalyseIpLogs(start_time, end_time)
	//分析前一天访问最多的播主、视频，并记入数据库
	Client := sls.CreateNormalInterface(aliLog.Endpoint, global.Cfg.String("ali_secret"), global.Cfg.String("ali_accessKey"), "")

	//白名单
	//whiteLists := []int64{1}

	//根据请求量封禁
	timeEnd := utils.Time() - 120
	var point1, point2, point3, point4 int
	hour := utils.ParseInt(time.Unix(timeEnd, 0).Format("15"), 0)
	if hour >= 1 && hour <= 6 {
		point1 = 100 * 1.5
		point2 = 300 * 1.5
		point3 = 500 * 1.5
		point4 = 250 * 1.5
	} else {
		point1 = 200 * 1.5
		point2 = 600 * 1.5
		point3 = 1000 * 1.5
		point4 = 450 * 1.5
	}
	sql1 := fmt.Sprintf(`env:"prod" and log_type:"Format" and url not "/v1/" and url not "/v1/" and url not "/v1/" and url not "/v1/" | select uid,COUNT(DISTINCT request_id) as pv group by uid HAVING pv>=%d order by pv desc`, point1)
	sql2 := fmt.Sprintf(`env:"prod" and log_type:"Format" and url not "/v1/" and url not "/v1/" and url not "/v1/" and url not "/v1/" | select uid,COUNT(DISTINCT request_id) as pv group by uid HAVING pv>=%d order by pv desc`, point2)
	sql3 := fmt.Sprintf(`env:"prod" and log_type:"Format" and url not "/v1/" and url not "/v1/" and url not "/v1/" and url not "/v1/" | select uid,COUNT(DISTINCT request_id) as pv group by uid HAVING pv>=%d order by pv desc`, point3)
	sql4 := fmt.Sprintf(`env:"prod" and log_type:"Format" and url not "/v1/" and url not "/v1/" and url not "/v1/" and url not "/v1/" | select uid,url,COUNT(DISTINCT request_id) as pv group by uid,url HAVING pv>=%d order by pv desc`, point4)

	//fmt.Println(sql1,sql2,sql3)
	res1, err := Client.GetLogs("dongchamao-api-log", "dongchamao-log-api-history", "", timeEnd-300, timeEnd, sql1, 100, 0, true)
	if err != nil {
		comErr = global.NewMsgError("获取五分钟内统计数据失败")
		return
	}
	fmt.Println(res1)

	res2, err := Client.GetLogs("dongchamao-api-log", "dongchamao-log-api-history", "", timeEnd-1800, timeEnd, sql2, 100, 0, true)
	if err != nil {
		comErr = global.NewMsgError("获取半小时内统计数据失败")
		return
	}
	//fmt.Println(res1)

	res3, err := Client.GetLogs("dongchamao-api-log", "dongchamao-log-api-history", "", timeEnd-3600, timeEnd, sql3, 100, 0, true)
	if err != nil {
		comErr = global.NewMsgError("获取一小时内统计数据失败")
		return
	}
	//fmt.Println(res2)
	res4, err := Client.GetLogs("dongchamao-api-log", "dongchamao-log-api-history", "", timeEnd-3600, timeEnd, sql4, 100, 0, true)
	if err != nil {
		comErr = global.NewMsgError("获取一小时内按URL统计数据失败")
		return
	}

	uidArr2 := make([]int64, 0)
	uidArr3 := make([]int64, 0)
	uidArr4 := make([]int64, 0)

	for _, re2 := range res2.Logs {
		uid2 := utils.ParseInt64String(re2["uid"])
		uidArr2 = append(uidArr2, uid2)
	}

	for _, re3 := range res3.Logs {
		uid3 := utils.ParseInt64String(re3["uid"])
		uidArr3 = append(uidArr3, uid3)
	}

	for _, re4 := range res4.Logs {
		uid4 := utils.ParseInt64String(re4["uid"])
		uidArr4 = append(uidArr4, uid4)
	}
	//todo 处理对应数据
	return
}
