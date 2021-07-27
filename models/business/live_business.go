package business

import (
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"dongchamao/structinit/repost/dy"
	"math"
	"time"
)

const LiveShareUrl = "https://www.iesdouyin.com/share/live/"

type LiveBusiness struct {
}

func NewLiveBusiness() *LiveBusiness {
	return new(LiveBusiness)
}

//获取直播间商品讲解数据
func (l *LiveBusiness) HbaseGetLiveCurProduct(roomId string) (data entity.DyLiveCurProduct, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveCurProduct).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyLiveCurProductMap)
	utils.MapToStruct(infoMap, &data)
	return
}

//直播间全网销量
func (l *LiveBusiness) GetHbaseRoomProductInfo(roomId, productId string) (data entity.DyRoomProduct, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	rowKey := roomId + "_" + productId
	result, err := query.SetTable(hbaseService.HbaseDyRoomProduct).GetByRowKey([]byte(rowKey))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	infoMap := hbaseService.HbaseFormat(result, entity.DyRoomProductMap)
	utils.MapToStruct(infoMap, &data)
	return
}

//获取讲解商品数据
func (l *LiveBusiness) RoomCurProductByIds(roomId string, productIds []string) map[string]dy.LiveCurProductCount {
	curInfo, _ := l.HbaseGetLiveCurProduct(roomId)
	productMap := map[string]dy.LiveCurProductCount{}
	for _, v := range curInfo.Promotion {
		if !utils.InArrayString(v.ProductID, productIds) {
			continue
		}
		if v.EndTime == 0 {
			continue
		}
		curSecond := v.EndTime - v.StartTime
		incSales := v.EndSales - v.StartSales
		endSales := v.EndSales
		var avgUserCount int64 = 0
		if v.TotalCrawlTimes > 0 {
			avgUserCount = v.TotalUserCount / v.TotalCrawlTimes
		}
		cur := dy.LiveCurProduct{
			StartTime:    v.StartTime,
			EndTime:      v.EndTime,
			AvgUserCount: avgUserCount,
			IncSales:     incSales,
			StartSales:   v.StartSales,
			EndSales:     endSales,
		}
		if c, ok := productMap[v.ProductID]; ok {
			c.CurSecond += curSecond
			if c.MaxPrice < v.PriceMax {
				c.MaxPrice = v.PriceMax
			}
			if c.MinPrice > v.PriceMin {
				c.MinPrice = v.PriceMin
			}
			c.CurNum += 1
			c.CurList = append(c.CurList, cur)
			productMap[v.ProductID] = c
		} else {
			productMap[v.ProductID] = dy.LiveCurProductCount{
				CurSecond: curSecond,
				MaxPrice:  v.PriceMax,
				MinPrice:  v.PriceMin,
				CurNum:    1,
				ShopId:    v.ShopId,
				ShopName:  v.ShopName,
				ShopIcon:  v.ShopIcon,
				CurList:   []dy.LiveCurProduct{cur},
			}
		}
	}
	return productMap
}

//获取商品数据
func (l *LiveBusiness) RoomPmtProductByIds(roomId string, productIds []string) map[string][]dy.LiveRoomProductSaleStatus {
	pmtInfo, _ := l.HbaseGetLivePmt(roomId)
	productMap := map[string][]dy.LiveRoomProductSaleStatus{}
	for _, v := range pmtInfo.Promotions {
		if !utils.InArrayString(v.ProductID, productIds) {
			continue
		}
		if v.StartTime > 0 {
			if _, ok := productMap[v.ProductID]; !ok {
				productMap[v.ProductID] = make([]dy.LiveRoomProductSaleStatus, 0)
			}
			productMap[v.ProductID] = append(productMap[v.ProductID], dy.LiveRoomProductSaleStatus{
				StartTime:  v.StartTime,
				StopTime:   v.StopTime,
				StartSales: v.InitialSales,
				FinalSales: v.FinalSales,
			})
		}
	}
	return productMap
}

//直播间信息
func (l *LiveBusiness) HbaseGetLiveInfo(roomId string) (data entity.DyLiveInfo, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveInfo).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	liveInfoMap := hbaseService.HbaseFormat(result, entity.DyLiveInfoMap)
	utils.MapToStruct(liveInfoMap, &data)
	data.Cover = dyimg.Fix(data.Cover)
	data.User.Avatar = dyimg.Fix(data.User.Avatar)
	data.RoomID = roomId
	return
}

//直播间商品趋势
//func (l *LiveBusiness) CountRoomProductSaleChart(roomId string) {
//	pmtInfo,comErr := l.HbaseGetLivePmt(roomId)
//	if comErr != nil {
//		return
//	}
//	productIdMap := map[string]string{}
//	for _,v := range pmtInfo.Promotions {
//		productIdMap[v.ProductID] = v.ProductID
//	}
//	productRoomChan := make(chan map[string]entity.DyRoomProduct,0)
//
//}

//直播间带货数据
func (l *LiveBusiness) HbaseGetLiveSalesData(roomId string) (data entity.DyAuthorLiveSalesData, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAuthorLiveSalesData).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyAuthorLiveSalesDataMap)
	utils.MapToStruct(detailMap, &data)
	return
}

//OnlineTrends转化
func (l *LiveBusiness) DealOnlineTrends(liveInfo entity.DyLiveInfo) (entity.DyLiveIncOnlineTrendsChart, entity.DyLiveOnlineTrends, int64) {
	onlineTrends := OnlineTrendOrderByTime(liveInfo.OnlineTrends)
	beforeTrend := entity.DyLiveOnlineTrends{}
	incTrends := make([]entity.DyLiveIncOnlineTrends, 0)
	dates := make([]string, 0)
	dates = append(dates, time.Unix(liveInfo.CreateTime, 0).Format("2006-01-02 15:04:05"))
	incTrends = append(incTrends, entity.DyLiveIncOnlineTrends{
		UserCount: 0,
		WatchInc:  0,
	})
	maxLiveOnlineTrends := entity.DyLiveOnlineTrends{}
	lenNum := len(onlineTrends)
	//平均在线人数
	var sumUserCount int64 = 0
	for k, v := range onlineTrends {
		sumUserCount += v.UserCount
		var inc int64 = 0
		if k != 0 {
			inc = v.WatchCnt - beforeTrend.WatchCnt
		} else {
			maxLiveOnlineTrends = v
		}
		if v.UserCount > maxLiveOnlineTrends.UserCount {
			maxLiveOnlineTrends = v
		}
		startFormat := time.Unix(v.CrawlTime, 0).Format("2006-01-02 15:04:05")
		dates = append(dates, startFormat)
		incTrends = append(incTrends, entity.DyLiveIncOnlineTrends{
			UserCount: v.UserCount,
			WatchInc:  inc,
		})
		beforeTrend = v
	}
	incTrendsChart := entity.DyLiveIncOnlineTrendsChart{
		Date:            dates,
		IncOnlineTrends: incTrends,
	}
	var avgUserCount int64 = 0
	if lenNum > 0 {
		avgUserCount = sumUserCount / int64(lenNum)
	}
	return incTrendsChart, maxLiveOnlineTrends, avgUserCount
}

//直播间分析
func (l *LiveBusiness) LiveRoomAnalyse(roomId string) (data dy.DyLiveRoomAnalyse, comErr global.CommonError) {
	data = dy.DyLiveRoomAnalyse{}
	liveInfo, comErr := l.HbaseGetLiveInfo(roomId)
	if comErr != nil {
		return
	}
	data.TotalUserCount = liveInfo.TotalUser
	data.LiveStartTime = liveInfo.CreateTime
	data.BarrageCount = liveInfo.BarrageCount
	data.AvgOnlineTime = l.CountAvgOnlineTime(liveInfo.OnlineTrends, liveInfo.CreateTime, liveInfo.TotalUser)
	liveInfo.OnlineTrends = OnlineTrendOrderByTime(liveInfo.OnlineTrends)
	lenNum := len(liveInfo.OnlineTrends)
	if lenNum > 1 {
		data.IncFans = liveInfo.OnlineTrends[lenNum-1].FollowerCount - liveInfo.OnlineTrends[0].FollowerCount
	}
	if liveInfo.RoomStatus == 2 {
		data.LiveLongTime = time.Now().Unix() - liveInfo.CreateTime
	} else {
		data.LiveLongTime = liveInfo.FinishTime - liveInfo.CreateTime
	}
	salesData, _ := l.HbaseGetLiveSalesData(roomId)
	if liveInfo.TotalUser > 0 {
		data.Uv = (salesData.Gmv + float64(salesData.TicketCount)/10) / float64(liveInfo.TotalUser)
		data.SaleRate = salesData.Sales / float64(liveInfo.TotalUser)
		data.IncFansRate = float64(data.IncFans) / float64(liveInfo.TotalUser)
		data.InteractRate = float64(liveInfo.BarrageCount) / float64(liveInfo.TotalUser)
	}
	data.Volume = int64(math.Floor(salesData.Sales))
	data.Amount = salesData.Gmv
	data.PromotionNum = salesData.NumProducts
	if salesData.Sales > 0 {
		data.PerPrice = salesData.Gmv / salesData.Sales
	}
	var sumUserCount int64 = 0
	for _, v := range liveInfo.OnlineTrends {
		sumUserCount += v.UserCount
	}
	if lenNum > 0 {
		data.AvgUserCount = sumUserCount / int64(lenNum)
	}
	return
}

//平均停留时长计算
func (l *LiveBusiness) CountAvgOnlineTime(onlineTrends []entity.DyLiveOnlineTrends, startTime, totalUser int64) float64 {
	lenNum := len(onlineTrends)
	var avgOnlineTime float64 = 0
	if lenNum == 0 {
		return avgOnlineTime
	} else if lenNum == 1 {
		onlineTrend := onlineTrends[0]
		avgOnlineTime = float64(onlineTrend.CrawlTime-startTime) * (float64(onlineTrend.UserCount) / 2) / float64(totalUser)
		return avgOnlineTime
	}
	onlineTrends = OnlineTrendOrderByTimeDesc(onlineTrends)
	for k, v := range onlineTrends {
		if k == lenNum-1 {
			avgOnlineTimeTmp := float64(v.CrawlTime-startTime) * (float64(v.UserCount) / 2) / float64(totalUser)
			avgOnlineTime += avgOnlineTimeTmp
			break
		}
		next := onlineTrends[k+1]
		avgOnlineTimeTmp := float64(v.CrawlTime-next.CrawlTime) * (float64(v.UserCount+next.UserCount) / 2) / float64(totalUser)
		avgOnlineTime += avgOnlineTimeTmp
	}
	return avgOnlineTime
}

//直播间信息
func (l *LiveBusiness) HbaseGetLivePmt(roomId string) (data entity.DyLivePmt, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLivePmt).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLivePmtMap)
	utils.MapToStruct(detailMap, &data)
	return
}

//获取直播榜单数据
func (l *LiveBusiness) HbaseGetRankTrends(roomId string) (data []entity.DyLiveRankTrend, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyLiveRankTrend).GetByRowKey([]byte(roomId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	detailMap := hbaseService.HbaseFormat(result, entity.DyLiveRankTrendsMap)
	hData := &entity.DyLiveRankTrends{}
	utils.MapToStruct(detailMap, hData)
	data = RankTrendOrderByTime(hData.RankData)
	return
}
