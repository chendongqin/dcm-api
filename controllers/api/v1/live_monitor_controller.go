package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/alias"
	"dongchamao/global/logger"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	"dongchamao/models/repost/dy"
	jsoniter "github.com/json-iterator/go"
	"math"
	"sort"
	"strconv"
	"time"
)

type LiveMonitorController struct {
	controllers.ApiBaseController
}

func (receiver *LiveMonitorController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
}

//直播监控价格
func (receiver *LiveMonitorController) MonitorPrice() {
	priceString := business.NewConfigBusiness().GetConfigJson("monitor_price", true)
	priceList := dy.LiveMonitorPriceList{}
	_ = jsoniter.Unmarshal([]byte(priceString), &priceList)
	receiver.SuccReturn(map[string]interface{}{
		"list": priceList.MonitorPrice,
	})
	return
}

//添加监控
func (receiver *LiveMonitorController) AddLiveMonitor() {
	inputData := receiver.InputFormat()
	authorId := business.IdDecrypt(inputData.GetString("author_id", ""))
	startTimestamp := inputData.GetInt64("start", 0)
	endTimestamp := inputData.GetInt64("end", 0)
	notice := inputData.GetBool("notice", false)
	if startTimestamp == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if endTimestamp == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	startTime := time.Unix(startTimestamp, 0)
	endTime := time.Unix(endTimestamp, 0)
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	_, comErr := hbase.GetAuthor(authorId)
	if comErr != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	now := time.Now().Local()
	if startTime.Before(now.AddDate(0, 0, -1)) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if endTime.Before(now) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if endTime.Before(startTime) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	liveMonitorBusiness := business.NewLiveMonitorBusiness()
	if liveMonitorBusiness.CheckRepeat(receiver.UserId, authorId, startTime, endTime) {
		receiver.FailReturn(global.NewMsgError("所选时段存在相同达人监控，无法重复添加。"))
		return
	}
	userInfo := dcm.DcUser{}
	exist, err := dcm.Get(receiver.UserId, &userInfo)
	if !exist || logger.CheckError(err) != nil {
		receiver.FailReturn(global.NewMsgError("获取用户信息失败"))
		return
	}
	// 当前可用次数
	remainFreeCount, remainPurchaseCount := liveMonitorBusiness.GetRemainingCount(receiver.UserId)
	currentCount := remainPurchaseCount + remainFreeCount
	// 本次需要消耗的次数
	nowCount := liveMonitorBusiness.CalcSpendCount(startTime, endTime)
	countEnough, useFreeCount, usePurchaseCount := liveMonitorBusiness.ExplainUseCount(nowCount, remainFreeCount, remainPurchaseCount)
	if !countEnough {
		remainCount := strconv.Itoa(currentCount)
		receiver.FailReturn(global.NewCodeError(4400, "监测失败，您的监控次数剩余["+remainCount+"]次，所选时间需要消耗["+strconv.Itoa(nowCount)+"]次直播监测次数。"))
		return
	}
	monitor := dcm.DcLiveMonitor{}
	monitor.OpenId = userInfo.Openid
	monitor.StartTime = startTime
	monitor.FreeCount = useFreeCount
	monitor.PurchaseCount = usePurchaseCount
	monitor.EndTime = endTime
	monitor.Notice = utils.BoolToInt(notice)
	monitor.FinishNotice = utils.BoolToInt(notice)
	monitor.ProductId = ""
	monitor.AuthorId = authorId
	monitor.UserId = receiver.UserId
	monitor.NextTime = startTime.Unix()
	monitor.CreateTime = now
	monitor.UpdateTime = now
	monitor.Source = business.LiveMonitorSourceDcm
	lastId, err := liveMonitorBusiness.AddLiveMonitor(&monitor)
	if logger.CheckError(err) != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	//	//看是否需要推送到粉丝画像加速计算队列
	//	go func() {
	//		defer global.RecoverPanic()
	//		err := logic.NewFansFeature().AutoRefresh(authorId)
	//		if err != nil {
	//			logs.Error("[粉丝画像] 直播监控推送粉丝画像失败 authorId: %s, err: %s", authorId, err)
	//		}
	//	}()
	//	c.SuccReturn(alias.M{
	//		"monitor_id": lastId,
	//	})
	receiver.SuccReturn(alias.M{
		"monitor_id": lastId,
	})
	return
}

//本月剩余价值
func (receiver *LiveMonitorController) LiveMonitorNum() {
	liveMonitorBusiness := business.NewLiveMonitorBusiness()
	// 当前可用次数
	remainFreeCount, remainPurchaseCount := liveMonitorBusiness.GetRemainingCount(receiver.UserId)
	currentCount := remainPurchaseCount + remainFreeCount
	receiver.SuccReturn(alias.M{
		"current_count":         currentCount,
		"remain_free_count":     remainFreeCount,
		"remain_purchase_count": remainPurchaseCount,
	})
	return
}

// 计算需要消耗的次数
func (receiver *LiveMonitorController) LiveMonitorCalcCount() {
	startTimestamp := utils.ToInt64(receiver.Ctx.Input.Param(":start"))
	endTimestamp := utils.ToInt64(receiver.Ctx.Input.Param(":end"))
	if startTimestamp == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if endTimestamp == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	startTime := time.Unix(startTimestamp, 0)
	endTime := time.Unix(endTimestamp, 0)
	if startTime.After(endTime) || startTime.Before(time.Now().AddDate(0, 0, -1)) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	count := business.NewLiveMonitorBusiness().CalcSpendCount(startTime, endTime)
	receiver.SuccReturn(alias.M{
		"count": count,
	})
	return
}

// 获取监控列表
func (receiver *LiveMonitorController) LiveMonitorList() {
	keyword := receiver.GetString("keyword", "")
	start := receiver.GetString("start", "")
	end := receiver.GetString("end", "")
	page := receiver.GetPage("page")
	size := receiver.GetPageSize("page_size", 10, 100)
	status, _ := receiver.GetInt("status", -1)
	list, totalCount := business.NewLiveMonitorBusiness().LiveMonitorRoomList(receiver.UserId, status, keyword, page, size, start, end)
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": totalCount,
	})
	return
}

// 取消监控
func (receiver *LiveMonitorController) CancelLiveMonitor() {
	monitorId := utils.ToInt(receiver.Ctx.Input.Param(":monitor_id"))
	if business.NewLiveMonitorBusiness().CancelLiveMonitor(receiver.UserId, monitorId) {
		receiver.SuccReturn(nil)
		return
	}
	receiver.FailReturn(global.NewError(5000))
	return
}

//用户已读操作
func (receiver *LiveMonitorController) ReadLiveMonitor() {
	monitorId := receiver.Ctx.Input.Param(":monitor_id")
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	_, _ = dbSession.ID(monitorId).Where("user_id = ?", receiver.UserId).Cols("has_new").Update(&dcm.DcLiveMonitor{HasNew: 0})
	receiver.SuccReturn(nil)
	return
}

//直播间列表
func (receiver *LiveMonitorController) LiveMonitorRooms() {
	monitorId := receiver.Ctx.Input.Param(":monitor_id")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	dbSession := dcm.GetDbSession()
	start := (page - 1) * pageSize
	defer dbSession.Close()
	list := make([]dcm.DcLiveMonitorRoom, 0)
	totalCount, _ := dbSession.
		Where("user_id = ? AND monitor_id = ?", receiver.UserId, monitorId).
		OrderBy("id desc").
		Limit(pageSize, start).
		FindAndCount(&list)
	var roomIds []string
	for _, v := range list {
		roomIds = append(roomIds, v.RoomId)
	}
	//获取直播间列表
	rooms, _ := hbase.GetLiveInfoByIds(roomIds)
	roomList := make([]dy.DyLiveRoomSimple, 0)
	for _, v := range rooms {
		var liveTime int64 = 0
		if v.FinishTime > 0 {
			liveTime = v.FinishTime - v.CreateTime
		} else {
			liveTime = time.Now().Unix() - v.CreateTime
		}
		roomList = append(roomList, dy.DyLiveRoomSimple{
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
			FinishTime: v.FinishTime,
			LiveTime:   liveTime,
			LikeCount:  v.LikeCount,
			RoomID:     business.IdEncrypt(v.RoomID),
			RoomStatus: v.RoomStatus,
			Title:      v.Title,
			TotalUser:  v.TotalUser,
			Gmv:        v.PredictGmv,
			Sales:      utils.ToInt64(math.Floor(v.PredictSales)),
		})
	}
	sort.Slice(roomList, func(i, j int) bool {
		return roomList[i].CreateTime > roomList[j].CreateTime
	})
	receiver.SuccReturn(alias.M{
		"list":  roomList,
		"total": totalCount,
	})
	return
}

//删除监控
func (receiver *LiveMonitorController) DeleteLiveMonitor() {
	monitorId := utils.ToInt(receiver.Ctx.Input.Param(":monitor_id"))
	if business.NewLiveMonitorBusiness().DeleteLiveMonitor(receiver.UserId, monitorId) {
		receiver.SuccReturn(nil)
		return
	}
	receiver.FailReturn(global.NewError(5000))
}
