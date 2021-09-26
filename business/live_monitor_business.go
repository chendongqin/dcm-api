package business

import (
	"dongchamao/global"
	"dongchamao/global/alias"
	"dongchamao/global/cache"
	"dongchamao/global/logger"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	"dongchamao/models/entity"
	"dongchamao/services/dyimg"
	"dongchamao/services/task"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-xorm/xorm"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"math"
	"strings"
	"time"
)

const (
	LiveMonitorStatusPending    = 0  //等待监控
	LiveMonitorStatusProcessing = 1  //监控中
	LiveMonitorStatusFinished   = 2  //监控完成
	LiveMonitorStatusCanceled   = 10 //取消监控

	LiveMonitorSourceDcm = 0 //来源洞察猫
)

type LiveMonitorBusiness struct {
}

func NewLiveMonitorBusiness() *LiveMonitorBusiness {
	return new(LiveMonitorBusiness)
}

func (receiver *LiveMonitorBusiness) ScanLiveRoom() {
	global.RecoverPanic()
	startTime := time.Now()
	totalCount := 0
	defer func() {
		logs.Info("[直播间监控] 记录数: %d, 耗时: %s", totalCount, time.Since(startTime))
	}()
	list, err := receiver.getNeedNoticeRooms()
	if err != nil {
		return
	}
	totalCount = len(list)
	if totalCount < 0 {
		return
	}
	taskPool := task.NewPool(3, 1024)
	for _, v := range list {
		monitorRoom := v
		room, err := hbase.GetLiveInfo(v.RoomId)
		if err != nil {
			continue
		}
		job := task.NewJob(func(job *task.Job) {
			receiver.checkRoom(monitorRoom, room)
		})
		taskPool.Push(job)
	}
	taskPool.PushDone()
	taskPool.Wait()
}

// 获取需要检查的直播间
func (receiver *LiveMonitorBusiness) getNeedNoticeRooms() (list []dcm.DcLiveMonitorRoom, err error) {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	before24HourTime := time.Now().Add(-86400 * time.Second)
	list = make([]dcm.DcLiveMonitorRoom, 0)
	err = dbSession.Where("create_time > ? AND status = 2", before24HourTime.Format("2006-01-02 15:04:05")).Find(&list)
	if err != nil {
		return
	}
	return
}

//检查直播间
func (receiver *LiveMonitorBusiness) checkRoom(monitorRoom dcm.DcLiveMonitorRoom, roomInfo entity.DyLiveInfo) {
	//在播
	//if roomInfo.RoomStatus == 2 {
	//	//检查商品上架
	//	productIds := strings.Split(monitorRoom.ProductId, ",")
	//	if len(productIds) < 0 {
	//		return
	//	}
	//	products, _ := roomInfo.GetProductListV2()
	//	productsMap := make(map[string]*entity.LiveRoomProductInfoV2)
	//	for _, product := range products {
	//		productsMap[product.ProductId] = product
	//	}
	//	for _, productId := range productIds {
	//		if product, exists := productsMap[productId]; exists {
	//			//半小时内检测到的商品才推送，减少重复推送验证次数
	//			if time.Now().Unix()-product.StartTime > 1800 {
	//				continue
	//			}
	//			//推送上架商品
	//			receiver.SendProductNotice(monitorRoom, product, roomInfo)
	//		}
	//	}
	//} else
	if roomInfo.RoomStatus == 4 { //下播
		_ = receiver.UpdateLiveRoomMonitor(&roomInfo)
		// 不存在微信openId则不继续推送
		if monitorRoom.OpenId == "" {
			return
		}
		//推送下播通知
		if monitorRoom.FinishNotice != 1 {
			return
		}
		receiver.SendFinishNotice(monitorRoom, roomInfo)
	} else {
		cacheKey := cache.GetCacheKey(cache.DyMonitorUpdateRoomLock, roomInfo.RoomID)
		cacheData := global.Cache.Get(cacheKey)
		if cacheData == "" {
			err := receiver.UpdateLiveRoomMonitor(&roomInfo)
			if err == nil {
				_ = global.Cache.Set(cacheKey, "1", 600)
			}
		}
	}

}

//下播微信提醒
func (receiver *LiveMonitorBusiness) SendFinishNotice(monitorRoom dcm.DcLiveMonitorRoom, roomInfo entity.DyLiveInfo) {
	if monitorRoom.OpenId == "" {
		return
	}
	msgMap := map[string]*message.TemplateDataItem{
		"first": {
			Value: "下播提醒",
			Color: "red",
		},
		"keyword1": {
			Value: roomInfo.User.Nickname,
			Color: "",
		},
		"keyword2": {
			Value: roomInfo.Title,
			Color: "",
		},
		"keyword3": {
			Value: time.Unix(roomInfo.CreateTime, 0).Format("2006-01-02 15:04:05") + " 至 " + time.Unix(roomInfo.FinishTime, 0).Format("2006-01-02 15:04:05"),
			Color: "",
		},
		"remark": {
			Value: "",
			Color: "",
		},
	}
	err := NewWechatBusiness().SendMsg(monitorRoom.OpenId, WechatMsgTemplateLiveMonitorFinish, msgMap, DyDcmUrl)
	if err != nil {
		return
	}
	_, _ = dcm.GetDbSession().Table(new(dcm.DcLiveMonitorRoom)).Where("id=?", monitorRoom.Id).Update(map[string]interface{}{
		"finish_notice": -1,
	})
	return
}

//开播微信提醒
func (receiver *LiveMonitorBusiness) SendLiveMonitorMsg(openId string, room *entity.DyLiveInfo) {
	if openId == "" {
		return
	}
	sendKey := "wechat:track:live:room:send:" + openId + ":" + utils.ToString(room.RoomID)
	hasSend := global.Cache.Get(sendKey)
	if hasSend != "" {
		return
	}
	//TrackTemplateId
	timeNowStr := utils.Date("", room.CreateTime)
	hourMinute := time.Unix(room.CreateTime, 0).Format("15:04")
	subTitle := fmt.Sprintf("您监控的达人[%s]在%s开始直播了", room.User.Nickname, hourMinute)
	liveInfoUrl := DyDcmUrl + "/#/live/detail/" + room.RoomID
	msgMap := map[string]*message.TemplateDataItem{
		"first": {
			Value: subTitle,
			Color: "#ED7D31",
		},
		"keyword1": {
			Value: room.Title,
			Color: "",
		},
		"keyword2": {
			Value: timeNowStr,
			Color: "",
		},
		"remark": {
			Value: "\n" + liveInfoUrl + "\n",
			Color: "",
		},
	}
	_ = NewWechatBusiness().SendMsg(openId, WechatMsgTemplateLiveMonitorBegin, msgMap, liveInfoUrl)
	_ = global.Cache.Set(sendKey, "1", 3600*24*7)
}

func (receiver *LiveMonitorBusiness) getMaxCountFromRights(level int) int {
	// 获取最大可以使用的次数
	if level == 0 {
		return LiveMonitorMonthMinNum
	}
	return LiveMonitorMonthMaxNum
}

func (receiver *LiveMonitorBusiness) GetMaxCount(userId int) (freeCount int, purchaseCount int) {
	vipBusiness := NewVipBusiness()
	vipLevel := vipBusiness.GetVipLevel(userId, VipPlatformDouYin)
	freeCount = receiver.getMaxCountFromRights(vipLevel.Level)
	purchaseCount = vipLevel.FeeLiveMonitor
	return
}

func (receiver *LiveMonitorBusiness) GetRemainingCount(userID int) (remainFreeCount int, remainPurchaseCount int) {
	useFreeCount, _ := receiver.GetCurrentCount(userID)
	// 总共可以使用的次数
	maxFreeCount, maxPurchaseCount := receiver.GetMaxCount(userID)
	remainFreeCount = maxFreeCount - useFreeCount
	if remainFreeCount <= 0 {
		remainFreeCount = 0
	}
	remainPurchaseCount = maxPurchaseCount
	if remainPurchaseCount <= 0 {
		remainPurchaseCount = 0
	}
	return
}

// ExplainUseCount 计算需要的次数
// 传入需要消耗的次数，剩余的免费次数和付费次数
// 得到次数是否足够使用，需要消耗多少的免费和付费次数
// can = false 表示剩余次数不足够
func (receiver *LiveMonitorBusiness) ExplainUseCount(needCount, remainFreeCount, remainPurchaseCount int) (can bool, useFreeCount, usePurchaseCount int) {
	can = true
	if remainFreeCount > needCount {
		useFreeCount = needCount
		return
	}
	useFreeCount = remainFreeCount
	usePurchaseCount = needCount - remainFreeCount
	if usePurchaseCount > remainPurchaseCount {
		can = false
	}
	return
}

// 获取本月已使用的次数
func (receiver *LiveMonitorBusiness) GetCurrentCount(userId int) (freeCount int, purchaseCount int) {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	beginDate, endDate := utils.GetFirstDateAndLastDate(time.Now())
	endDate = endDate.AddDate(0, 0, 1)
	sumFreeCount, err := dbSession.Where("user_id = ? AND create_time >= ? AND create_time < ? AND status <?", userId, beginDate, endDate, LiveMonitorStatusCanceled).Sum(new(dcm.DcLiveMonitor), "free_count")
	if logger.CheckError(err) != nil {
		return 0, 0
	}
	sumPurchaseCount, err := dbSession.Where("user_id = ? AND create_time >= ? AND create_time < ? AND status <?", userId, beginDate, endDate, LiveMonitorStatusCanceled).Sum(new(dcm.DcLiveMonitor), "purchase_count")
	if logger.CheckError(err) != nil {
		return 10000, 10000
	}
	freeCount = utils.ToInt(sumFreeCount)
	purchaseCount = utils.ToInt(sumPurchaseCount)
	return freeCount, purchaseCount
}

// 计算两个时间点需要消耗多少次监控次数
func (receiver *LiveMonitorBusiness) CalcSpendCount(startTime, endTime time.Time) int {
	spendTime := endTime.Unix() - startTime.Unix()
	// 不足6小时当6小时算，所以向上进位
	count := int(math.Ceil(float64(spendTime) / float64(6*3600)))
	return count
}

func (receiver *LiveMonitorBusiness) GetLast24HourOnlineRooms() []entity.DyLiveInfo {
	list, err := receiver.Get24HourRooms()
	if err != nil {
		return nil
	}
	var roomIds []string
	for _, v := range list {
		roomIds = append(roomIds, v.RoomId)
	}
	roomIds = utils.UniqueStringSlice(roomIds)
	if len(roomIds) <= 0 {
		return nil
	}
	finalList := make([]entity.DyLiveInfo, 0)
	for _, v := range roomIds {
		liveInfo, _ := hbase.GetLiveInfo(v)
		if liveInfo.RoomStatus == 2 {
			finalList = append(finalList, liveInfo)
		}
	}
	return finalList
}

// 获取近24小时监控到的直播间
func (receiver *LiveMonitorBusiness) Get24HourRooms() (list []dcm.DcLiveMonitorRoom, err error) {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	list = make([]dcm.DcLiveMonitorRoom, 0)
	err = dbSession.Table(new(dcm.DcLiveMonitorRoom)).Where("create_time >= ?", time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")).Find(&list)
	return
}

func (receiver *LiveMonitorBusiness) ScanLive(monitor *dcm.DcLiveMonitor) {
	startTime := time.Unix(monitor.NextTime, 0)
	endTime := startTime.AddDate(0, 0, 30)
	results, _ := hbase.GetAuthorRoomsRangDate(monitor.AuthorId, startTime, endTime)
	if len(results) <= 0 {
		//如果没找到数据，并且是第一次扫描
		if monitor.StartTime.Unix() == monitor.NextTime {
			receiver.firstScan(monitor)
		}
		return
	}
	item := entity.DyAuthorLiveRoom{}
	for _, v := range results {
		for _, r := range v {
			if r.CreateTime > item.CreateTime {
				item = r
			}
		}
	}
	receiver.foundNewLive(monitor, item.RoomID)
}

// 首次扫描
// 首次扫描无数据时触发从达人表中查询正在直播的直播间
func (receiver *LiveMonitorBusiness) firstScan(monitor *dcm.DcLiveMonitor) {
	author, _ := hbase.GetAuthor(monitor.AuthorId)
	liveRoomId := utils.ToString(author.RoomId)
	liveRoomStatus := utils.ToInt(author.RoomStatus)
	if liveRoomStatus != 2 {
		return
	}
	receiver.foundNewLive(monitor, liveRoomId)
}

func (receiver *LiveMonitorBusiness) foundNewLive(monitor *dcm.DcLiveMonitor, roomId string) {
	roomInfo, err := hbase.GetLiveInfo(roomId)
	if err != nil {
		logs.Error("[live monitor] 直播监控获取直播间数据失败，err: %s", err)
		return
	}
	if roomInfo.RoomID == "" {
		return
	}
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	_ = dbSession.Begin()
	// 新增直播记录
	if !receiver.AddByMonitor(dbSession, monitor, &roomInfo) {
		_ = dbSession.Rollback()
		return
	}

	// 更新下一次时间
	monitor.NextTime = roomInfo.CreateTime + 1
	if !receiver.UpdateNextTime(dbSession, monitor) {
		_ = dbSession.Rollback()
		return
	}
	_ = dbSession.Commit()

	// 发送提醒 & 推送
	receiver.liveNotice(monitor, &roomInfo)
}

func (receiver *LiveMonitorBusiness) liveNotice(monitor *dcm.DcLiveMonitor, roomInfo *entity.DyLiveInfo) {
	// 需要微信提醒
	if monitor.Notice == 1 && monitor.OpenId != "" {
		receiver.SendLiveMonitorMsg(monitor.OpenId, roomInfo)
	}
	//var wg sync.WaitGroup
	//wg.Add(2)
	//go func() {
	//	defer wg.Done()
	//	defer global.RecoverPanic()
	//	// 短信推送
	//	if monitor.Notice == 1 {
	//		user := dcm.DcUser{}
	//		_,err := dcm.Get(monitor.UserId,&user)
	//		mobile := user.Username
	//		if err == nil && mobile != "" {
	//			//todo 短信发送
	//		}
	//	}
	//}()
	//go func() {
	//	defer wg.Done()
	//	defer global.RecoverPanic()
	//	// APP推送
	//	userDeviceModel := apiv1models.NewVoUserDevices()
	//	if exists, _ := userDeviceModel.LoadByUser(monitor.UserId); exists {
	//		subTitle := "您监控的播主【" + roomInfo.Nickname + "】正在直播"
	//		if monitor.Source == douyinmodelsV2.LiveMonitorSourcePartner {
	//			subTitle = "您投放的达人【" + roomInfo.Nickname + "】开始直播带货开始了"
	//		}
	//		body := roomInfo.RoomTitle
	//		title := "直播监控"
	//		userIdStr := strconv.Itoa(monitor.UserId)
	//		roomId := monitor.RoomId
	//		lockResource := fmt.Sprintf("notice:push:favorite:author:%s:%s:%s", monitor.AuthorId, userDeviceModel.DeviceToken, roomInfo.RoomId)
	//		_, ok, err := mutex.TryLockWithTimeout(global.Cache.GetInstance().(redis.Conn), lockResource, "monitor"+utils.ToString(time.Now().Unix()), 43200)
	//		if err != nil {
	//			logs.Error("推送失败，获取锁失败[%s] [%s] [%s] [%s]", monitor.AuthorId, userIdStr, roomId, err)
	//			return
	//		}
	//		if !ok {
	//			logs.Debug("已经发送 [%s] [%s] [%s]", monitor.AuthorId, userIdStr, roomId)
	//			return
	//		}
	//		result := upush.SendByDevice(userDeviceModel, title, subTitle, body, "author", monitor.AuthorId)
	//		if !result {
	//			logs.Error("推送发送失败：user: [%d] token: [%s]", monitor.UserId, userDeviceModel.DeviceToken)
	//		}
	//	}
	//}()
	//wg.Wait()
}

// 新增直播记录
func (receiver *LiveMonitorBusiness) AddByMonitor(dbSession *xorm.Session, monitor *dcm.DcLiveMonitor, roomInfo *entity.DyLiveInfo) bool {
	if dbSession == nil {
		dbSession = dcm.GetDbSession()
		defer dbSession.Close()
	}
	now := time.Now()
	room := &dcm.DcLiveMonitorRoom{
		MonitorId:    monitor.Id,
		UserId:       monitor.UserId,
		AuthorId:     monitor.AuthorId,
		RoomId:       roomInfo.RoomID,
		Status:       roomInfo.RoomStatus,
		OpenId:       monitor.OpenId,
		FinishNotice: monitor.FinishNotice,
		ProductId:    monitor.ProductId,
		CreateTime:   now,
		UpdateTime:   now,
	}
	//if roomInfo.RoomStatus == 4 {
	room.Gmv = utils.ToString(roomInfo.PredictGmv)
	room.Sales = utils.ToInt(roomInfo.PredictSales)
	room.UserTotal = utils.ToInt(roomInfo.TotalUser)
	//}
	_, err := dbSession.InsertOne(room)
	if err != nil {
		return false
	}
	return true
}

func (receiver *LiveMonitorBusiness) CheckRepeat(userId int, authorId string, startTime time.Time, endTime time.Time) bool {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	rows, err := dbSession.
		Where("user_id = ? AND end_time > ? AND start_time < ? AND author_id = ? AND del_status = 0 AND status <?",
			userId,
			startTime.Format("2006-01-02 15:04:05"),
			endTime.Format("2006-01-02 15:04:05"),
			authorId,
			LiveMonitorStatusCanceled,
		).
		Count(&dcm.DcLiveMonitor{})
	if err != nil {
		logs.Error("[live monitor] 排重错误, err: %s", err)
	}
	return rows >= 1
}

// 新增一条监控记录
func (receiver *LiveMonitorBusiness) AddLiveMonitor(liveMonitor *dcm.DcLiveMonitor) (lastId int64, err error) {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	now := time.Now()
	if liveMonitor.StartTime.Before(now) {
		liveMonitor.Status = 1
	}
	_ = dbSession.Begin()
	_, err = dbSession.InsertOne(liveMonitor)
	if err != nil {
		_ = dbSession.Rollback()
		NewMonitorBusiness().SendErr("直播监控", fmt.Sprintf("[live monitor] add live monitor failed. err: %s", err))
		return
	}
	if liveMonitor.PurchaseCount > 0 {
		affect, err1 := dbSession.Table(new(dcm.DcUserVip)).
			Where("user_id=? AND platform=? AND live_monitor_num >= ?", liveMonitor.UserId, VipPlatformDouYin, liveMonitor.PurchaseCount).
			Cols("live_monitor_num").
			Decr("live_monitor_num", liveMonitor.PurchaseCount).
			Update(new(dcm.DcUserVip))
		if affect == 0 || err1 != nil {
			NewMonitorBusiness().SendErr("直播监控", fmt.Sprintf("[live monitor] add live monitor failed. err: %s", err))
			_ = dbSession.Rollback()
			err = err1
			return
		}
	}
	lastId = int64(liveMonitor.Id)
	author, _ := hbase.GetAuthor(liveMonitor.AuthorId)
	if liveMonitor.StartTime.Before(now) && liveMonitor.EndTime.After(now) && author.RoomStatus == 2 {
		if room, err1 := hbase.GetLiveInfo(author.RoomId); err1 == nil {
			if liveMonitor.NextTime <= room.CreateTime {
				receiver.AddByMonitor(dbSession, liveMonitor, &room)
				affect, err2 := dbSession.
					Table(new(dcm.DcLiveMonitor)).
					Where("id=?", liveMonitor.Id).
					Update(map[string]interface{}{"next_time": room.CreateTime + 1})
				if affect == 0 || err2 != nil {
					NewMonitorBusiness().SendErr("直播监控", fmt.Sprintf("[live monitor] add live monitor failed. err: %s", err1))
					_ = dbSession.Rollback()
					err = err2
					return
				}
			}
		}
	}
	_ = dbSession.Commit()
	go NewSpiderBusiness().AddLive(liveMonitor.AuthorId, author.FollowerCount, AddLiveTopMonitored, liveMonitor.EndTime.Unix())
	return
}

// 更新直播间记录
func (receiver *LiveMonitorBusiness) UpdateLiveRoomMonitor(roomInfo *entity.DyLiveInfo) (err error) {
	updateMap := map[string]interface{}{
		"gmv":         roomInfo.PredictGmv,
		"sales":       roomInfo.PredictSales,
		"user_total":  roomInfo.TotalUser,
		"update_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	if roomInfo.RoomStatus == 4 {
		updateMap["finish_time"] = roomInfo.FinishTime
	}
	_, err = dcm.GetDbSession().
		Table(new(dcm.DcLiveMonitorRoom)).
		Cols("gmv", "sales", "user_total", "update_time").
		Where("room_id=?", roomInfo.RoomID).
		Update(updateMap)
	if err != nil {
		NewMonitorBusiness().SendErr("更新直播间记录", err.Error())
	}
	return
}

// 取消监控
func (receiver *LiveMonitorBusiness) CancelLiveMonitor(userId int, monitorId int) bool {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	item := &dcm.DcLiveMonitor{
		Status:    LiveMonitorStatusCanceled,
		DelStatus: 1,
	}
	_ = dbSession.Begin()
	rows, err := dbSession.Where("user_id = ? AND id = ? AND status = ?", userId, monitorId, LiveMonitorStatusPending).Cols("status", "del_status").Update(item)
	if err != nil {
		_ = dbSession.Rollback()
		return false
	}
	_ = dbSession.Commit()
	return rows >= 1
}

// 删除监控
func (receiver *LiveMonitorBusiness) DeleteLiveMonitor(userId int, monitorId int) bool {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	item := &dcm.DcLiveMonitor{
		DelStatus: 1,
	}
	rows, err := dbSession.Where("user_id = ? AND id = ?", userId, monitorId).In("status", LiveMonitorStatusProcessing, LiveMonitorStatusFinished, LiveMonitorStatusCanceled).Cols("del_status").Update(item)
	if err != nil {
		return false
	}
	return rows >= 1
}

func (receiver *LiveMonitorBusiness) GetLiveMonitorAuthors(dbSession *xorm.Session, userId int, keyword string) []string {
	var authorIds []string
	tempList := make([]dcm.DcLiveMonitor, 0)
	_ = dbSession.Select("author_id").Where("user_id = ?", userId).GroupBy("author_id").Find(&tempList)
	for _, v := range tempList {
		authorIds = append(authorIds, v.AuthorId)
	}
	list, _ := hbase.GetAuthorByIds(authorIds)
	var finalAuthorsId []string
	for _, v := range list {
		shortId := v.Data.ShortID
		uniqueId := v.Data.UniqueID
		nickname := v.Data.Nickname
		authorId := v.AuthorID
		if strings.Contains(nickname, keyword) || strings.Contains(shortId, keyword) || strings.Contains(uniqueId, keyword) {
			finalAuthorsId = append(finalAuthorsId, authorId)
		}
	}
	return finalAuthorsId
}

//
func (receiver *LiveMonitorBusiness) LiveMonitorRoomList(userId int, status int, keyword string, page int, size int, start, end string) (list []dcm.DcLiveMonitor, totalCount int64) {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	if keyword != "" {
		authorIds := receiver.GetLiveMonitorAuthors(dbSession, userId, keyword)
		dbSession.In("author_id", authorIds)
	}
	pageStart := (page - 1) * size
	dbSession.Where("user_id = ? AND del_status = 0 AND source = ?", userId, LiveMonitorSourceDcm).
		Limit(size, pageStart).
		OrderBy("id desc")
	if status >= 0 {
		dbSession.And("status = ?", status)
	}
	if start != "" {
		startTime, err := time.ParseInLocation("2006-01-02", start, time.Local)
		if err != nil {
			return
		}
		dbSession.And("create_time > ?", startTime.Format("2006-01-02 15:04:05"))
	}
	if end != "" {
		endTime, err := time.ParseInLocation("2006-01-02", end, time.Local)
		if err != nil {
			return
		}
		dbSession.And("create_time < ?", endTime.AddDate(0, 0, 1).Format("2006-01-02 15:04:05"))
	}
	list = make([]dcm.DcLiveMonitor, 0)
	totalCount, _ = dbSession.FindAndCount(&list)

	var authorIds []string
	var monitorIds []int
	existsAuthorIds := make(map[string]bool)
	for _, v := range list {
		monitorIds = append(monitorIds, v.Id)
		if _, exists := existsAuthorIds[v.AuthorId]; !exists {
			authorIds = append(authorIds, v.AuthorId)
			existsAuthorIds[v.AuthorId] = true
		}
	}
	rooms := make([]map[string]interface{}, 0)
	_ = dbSession.Table(&dcm.DcLiveMonitorRoom{}).
		In("monitor_id", monitorIds).
		Select("monitor_id, count(1) as num, max(room_id) as room_id,sum(gmv) as gmv,sum(sales) as sales,sum(user_total) as total_user").
		GroupBy("monitor_id").
		Find(&rooms)
	roomsGroup := make(map[int]alias.M)
	for _, v := range rooms {
		id := utils.ToInt(v["monitor_id"])
		roomsGroup[id] = v
	}
	authorsMap, _ := hbase.GetAuthorByIds(authorIds)
	for k, v := range list {
		detail, ok := authorsMap[v.AuthorId]
		if !ok {
			detail, _ = hbase.GetAuthor(v.AuthorId)
			authorsMap[v.AuthorId] = detail
		}

		shortId := detail.Data.ShortID
		uniqueId := detail.Data.UniqueID
		avatar := detail.Data.Avatar
		nickname := detail.Data.Nickname
		// 修正抖音号
		finalUniqueId := uniqueId
		if finalUniqueId == "" {
			if shortId == "" {
				finalUniqueId = v.AuthorId
			} else {
				finalUniqueId = shortId
			}
		}
		// 计算房间数
		if roomsInfo, exists := roomsGroup[v.Id]; exists {
			list[k].RoomId = IdEncrypt(utils.ToString(roomsInfo["room_id"]))
			list[k].RoomCount = utils.ToInt(roomsInfo["num"])
			list[k].TotalUser = utils.ToInt64(utils.ToFloat64(roomsInfo["total_user"]))
			list[k].Sales = utils.ToInt64(utils.ToFloat64(roomsInfo["sales"]))
			list[k].Gmv = utils.ToFloat64(roomsInfo["gmv"])
			if list[k].TotalUser > 0 {
				list[k].Uv = list[k].Gmv / float64(list[k].TotalUser)
			}
		}
		// 填充达人信息
		list[k].CreateTimeString = v.CreateTime.Format("2006-01-02 15:04:05")
		list[k].UpdateTimeString = v.UpdateTime.Format("2006-01-02 15:04:05")
		list[k].StartTimeString = v.StartTime.Format("2006-01-02 15:04:05")
		list[k].EndTimeString = v.EndTime.Format("2006-01-02 15:04:05")
		list[k].Nickname = nickname
		list[k].UniqueID = finalUniqueId
		list[k].Avatar = dyimg.Avatar(avatar, dyimg.AvatarLittle)
		list[k].AuthorId = IdEncrypt(v.AuthorId)
	}
	return
}

//更新直播监控
func (receiver *LiveMonitorBusiness) UpdateNextTime(dbSession *xorm.Session, monitor *dcm.DcLiveMonitor) bool {
	if dbSession == nil {
		dbSession = dcm.GetDbSession()
		defer dbSession.Close()
	}

	monitor.UpdateTime = time.Now()
	monitor.HasNew = 1

	rows, err := dbSession.ID(monitor.Id).Cols("has_new, updated_at, next_time").Update(monitor)
	if err != nil {
		logs.Error("[live monitor] update live monitor next time failed. err: %s", err)
		return false
	}
	return rows >= 1
}
