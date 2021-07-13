package business

import (
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"sort"
	"time"
)

type LiveBusiness struct {
}

func NewLiveBusiness() *LiveBusiness {
	return new(LiveBusiness)
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

//OnlineTrends转化
func (l *LiveBusiness) DealOnlineTrends(onlineTrends []entity.DyLiveOnlineTrends) (entity.DyLiveIncOnlineTrendsChart, entity.DyLiveOnlineTrends, int64) {
	beforeTrend := entity.DyLiveOnlineTrends{}
	incTrends := make([]entity.DyLiveIncOnlineTrends, 0)
	dates := make([]string, 0)
	maxLiveOnlineTrends := entity.DyLiveOnlineTrends{}
	lenNum := len(onlineTrends)
	for k, v := range onlineTrends {
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
	var incFans int64 = 0
	if lenNum > 0 {
		incFans = onlineTrends[lenNum-1].FollowerCount - onlineTrends[0].FollowerCount
	}
	return incTrendsChart, maxLiveOnlineTrends, incFans
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

//直播排名按时间排序
type RankTrendSortList []entity.DyLiveRankTrend

func RankTrendOrderByTime(rankTrends []entity.DyLiveRankTrend) []entity.DyLiveRankTrend {
	sort.Sort(RankTrendSortList(rankTrends))
	return rankTrends
}

func (I RankTrendSortList) Len() int {
	return len(I)
}

func (I RankTrendSortList) Less(i, j int) bool {
	return I[i].CrawlTime < I[j].CrawlTime
}

func (I RankTrendSortList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
