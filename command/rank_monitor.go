package command

import (
	"dongchamao/business"
	"dongchamao/global/utils"
	"strings"
	"time"
)

func CheckRank() {
	monitorBusiness := business.NewMonitorBusiness()
	currentTime := time.Now().Local()

	var monitorEvents []string

	//每天上午5点执行的监控
	if checkTime(currentTime, 5, 0) {
		//go checkLiveHotRank()
		//monitorEvents = append(monitorEvents, []string{"播主日榜监控", "商品日榜监控", "品牌日榜监控", "电商视频监控", "淘客推广排行榜监控"}...)
	}
	go checkLiveHotRank()
	monitorBusiness.SendErr("监控榜单执行", strings.Join(monitorEvents, ","))
	return
}

func checkTime(currentTime time.Time, hour, minute int) bool {
	if currentTime.Hour() == hour && currentTime.Minute() == minute {
		return true
	}
	return false
}

//demo
func checkLiveHotRank() {
	interBusiness := business.NewInterBusiness()
	testApi := interBusiness.BuildURL("/v1/dy/rank/live/top/2021-09-27/10:00")
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
		business.NewMonitorBusiness().SendTemplateMessage("S", "直播热榜", "榜单挂了")
	}
	return
}
