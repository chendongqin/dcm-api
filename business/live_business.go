package business

import (
	"dongchamao/business/es"
	"dongchamao/global"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	"dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	"math"
	"sort"
	"time"
)

const LiveShareUrl = "https://www.iesdouyin.com/share/live/"

type LiveBusiness struct {
}

func NewLiveBusiness() *LiveBusiness {
	return new(LiveBusiness)
}

func (l *LiveBusiness) RoomCurAndPmtProductById(roomId, productId string) (curProductCount dy.LiveCurProductCount, ptmSales []dy.LiveRoomProductSaleStatus, comErr global.CommonError) {
	roomProduct, comErr := hbase.GetRoomProductInfo(roomId + "_" + productId)
	if comErr != nil {
		return
	}
	curProductCount = dy.LiveCurProductCount{
		CurList: []dy.LiveCurProduct{},
	}
	for k, v := range roomProduct.Promotion {
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
		curProductCount.CurSecond += curSecond
		if k == 0 {
			curProductCount.MaxPrice = v.PriceMax
			curProductCount.MinPrice = v.PriceMin
		} else {
			if curProductCount.MaxPrice < v.PriceMax {
				curProductCount.MaxPrice = v.PriceMax
			}
			if curProductCount.MinPrice > v.PriceMin {
				curProductCount.MinPrice = v.PriceMin
			}
		}
		curProductCount.CurNum += 1
		curProductCount.CurList = append(curProductCount.CurList, cur)
	}
	sort.Slice(curProductCount.CurList, func(i, j int) bool {
		return curProductCount.CurList[i].StartTime < curProductCount.CurList[j].StartTime
	})
	ptmSales = []dy.LiveRoomProductSaleStatus{}
	for _, v := range roomProduct.PtmPromotion {
		if v.StartTime > 0 {
			ptmSales = append(ptmSales, dy.LiveRoomProductSaleStatus{
				StartTime:  v.StartTime,
				StopTime:   v.StopTime,
				StartSales: v.InitialSales,
				FinalSales: v.FinalSales,
				Sales:      v.Sales,
			})
		}
	}
	sort.Slice(ptmSales, func(i, j int) bool {
		return ptmSales[i].StartTime < ptmSales[j].StartTime
	})
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
	liveInfo, comErr := hbase.GetLiveInfo(roomId)
	data = dy.DyLiveRoomAnalyse{
		RoomId: liveInfo.RoomID,
		Title:  liveInfo.Title,
		Cover:  dyimg.Fix(liveInfo.Cover),
	}
	if comErr != nil {
		return
	}
	data.DiscoverTime = liveInfo.DiscoverTime
	data.TotalUserCount = liveInfo.TotalUser
	data.LiveStartTime = liveInfo.CreateTime
	data.BarrageCount = liveInfo.BarrageCount
	data.AvgOnlineTime = l.CountAvgOnlineTime(liveInfo.OnlineTrends, liveInfo.CreateTime, liveInfo.TotalUser)
	liveInfo.OnlineTrends = OnlineTrendOrderByTime(liveInfo.OnlineTrends)
	lenNum := len(liveInfo.OnlineTrends)
	data.IncFans = liveInfo.FollowCount
	if liveInfo.RoomStatus == 2 {
		data.LiveLongTime = time.Now().Unix() - liveInfo.CreateTime
	} else {
		data.LiveLongTime = liveInfo.FinishTime - liveInfo.CreateTime
	}
	salesData, _ := hbase.GetLiveSalesData(roomId)
	if salesData.Gmv == 0 {
		salesData.Gmv = liveInfo.PredictGmv
		salesData.Sales = liveInfo.PredictSales
		//if liveInfo.RealGmv > 0 {
		//	gmv = liveInfo.RealGmv
		//	sales = liveInfo.RealSales
		//}
	}
	if liveInfo.TotalUser > 0 {
		data.Uv = (salesData.Gmv + float64(salesData.TicketCount)/10) / float64(liveInfo.TotalUser)
		data.SaleRate = salesData.Sales / float64(liveInfo.TotalUser)
		data.IncFansRate = float64(data.IncFans) / float64(liveInfo.TotalUser)
		data.InteractRate = float64(liveInfo.BarrageCount) / float64(liveInfo.TotalUser)
	}
	data.Volume = int64(math.Floor(salesData.Sales))
	data.Amount = salesData.Gmv
	esLiveBusiness := es.NewEsLiveBusiness()
	data.PromotionNum = esLiveBusiness.CountRoomProductByRoomId(liveInfo)
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

//获取销售额、销量
func (l *LiveBusiness) LiveSalesData(roomId string) (gmv, sales float64) {
	salesData, _ := hbase.GetLiveSalesData(roomId)
	gmv = salesData.Gmv
	sales = salesData.Sales
	if salesData.Gmv == 0 {
		liveInfo, _ := hbase.GetLiveInfo(roomId)
		gmv = liveInfo.PredictGmv
		sales = liveInfo.PredictSales
	}
	return
}

//平均停留时长计算
func (l *LiveBusiness) CountAvgOnlineTime(onlineTrends []entity.DyLiveOnlineTrends, startTime, totalUser int64) float64 {
	lenNum := len(onlineTrends)
	var avgOnlineTime float64 = 0
	if lenNum == 0 || totalUser == 0 {
		return 0
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

//获取直播榜单数据
func (l *LiveBusiness) HbaseGetRankTrends(roomId string) (data []entity.DyLiveRankTrend, comErr global.CommonError) {
	hData, comErr := hbase.GetRankTrends(roomId)
	data = RankTrendOrderByTime(hData.RankData)
	return
}
