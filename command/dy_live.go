package command

import (
	"dongchamao/business"
	"dongchamao/models/dcm"
	"dongchamao/models/entity"
	"dongchamao/services/task"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math"
	"time"
)

type DyLive struct {
}

func UpdateLiveMonitorStatus() {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	var err error
	now := time.Now().Format("2006-01-02 15:04:05")
	// 将进行中的置为结束
	toFinished := &dcm.DcLiveMonitor{Status: business.LiveMonitorStatusFinished}
	_, err = dbSession.Where("end_time < ? AND status = ?", now, business.LiveMonitorStatusProcessing).Update(toFinished)
	// 将到点的任务置为开始
	toBegin := &dcm.DcLiveMonitor{Status: business.LiveMonitorStatusProcessing}
	_, err = dbSession.Where("start_time < ? AND status = ?", now, business.LiveMonitorStatusPending).Update(toBegin)
	if err != nil {
		logs.Error("[live monitor] 更新监控状态 err: %s", err)
	}
	return
}

// 直播间轮询
func LiveRoomMonitor() {
	business.NewLiveMonitorBusiness().ScanLiveRoom()
	return
}

func LiveMonitor() {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	list := make([]*dcm.DcLiveMonitor, 0)
	now := time.Now().Format("2006-01-02 15:04:05")
	err := dbSession.Where("start_time < ? AND end_time > ? AND status = ? AND del_status = 0", now, now, business.LiveMonitorStatusProcessing).Find(&list)
	if err != nil {
		logs.Error("[live monitor] 获取直播数失败 err: %s", err)
		return
	}
	startTime := time.Now()
	taskPool := task.NewPool(10, 1024)
	for _, v := range list {
		monitor := v
		job := task.NewJob(func(job *task.Job) {
			business.NewLiveMonitorBusiness().ScanLive(monitor)
		})
		taskPool.Push(job)
	}
	taskPool.PushDone()
	taskPool.Wait()
	spendTime := time.Since(startTime)
	recordCount := len(list)
	if spendTime.Seconds() >= 60 {
		business.NewMonitorBusiness().SendErr("直播监控超时", fmt.Sprintf("### 提醒\n\n直播监控[%d]条记录，消耗时间%s，需要尝试优化", recordCount, spendTime))
	}
	logs.Info("[live monitor] 直播监控记录 [%d] 条，耗时：%s", recordCount, spendTime)
	return
}

func DealLiveRankCharts(OriginList []entity.DyLiveRankTrend, chatNum int) (rank []int, date []int64) {
	dataLen := len(OriginList)
	rank = make([]int, 0)
	date = make([]int64, 0)
	if dataLen < 60 {
		for _, v := range OriginList {
			rank = append(rank, v.Rank)
			date = append(date, v.CrawlTime)
		}
	} else {
		stepFloat := math.Floor(float64(dataLen / chatNum))
		step := int(stepFloat)
		start := 0
		for i := 0; i < chatNum; i++ {
			start = i * step
			end := start + step
			tempMax := 1000000
			tempTime := int64(0)
			//分区取最高排名
			for _, v := range OriginList[start:end] {
				if v.Rank < tempMax {
					tempMax = v.Rank
					tempTime = v.CrawlTime
				}
			}
			rank = append(rank, tempMax)
			date = append(date, tempTime)
		}
	}
	return
}
