package business

import (
	"dongchamao/entity"
	"dongchamao/structinit/repost/dy"
	"sort"
)

//直播流按时间倒序
type OnlineTrendSortDescList []entity.DyLiveOnlineTrends

func OnlineTrendOrderByTimeDesc(onlineTrends []entity.DyLiveOnlineTrends) []entity.DyLiveOnlineTrends {
	sort.Sort(OnlineTrendSortDescList(onlineTrends))
	return onlineTrends
}

func (I OnlineTrendSortDescList) Len() int {
	return len(I)
}

func (I OnlineTrendSortDescList) Less(i, j int) bool {
	return I[i].CrawlTime > I[j].CrawlTime
}

func (I OnlineTrendSortDescList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//直播流按时间排序
type OnlineTrendSortList []entity.DyLiveOnlineTrends

func OnlineTrendOrderByTime(onlineTrends []entity.DyLiveOnlineTrends) []entity.DyLiveOnlineTrends {
	sort.Sort(OnlineTrendSortList(onlineTrends))
	return onlineTrends
}

func (I OnlineTrendSortList) Len() int {
	return len(I)
}

func (I OnlineTrendSortList) Less(i, j int) bool {
	return I[i].CrawlTime < I[j].CrawlTime
}

func (I OnlineTrendSortList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
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

//直播讲解按时间排序
type ProductCurSortList []dy.LiveCurProduct

func ProductCurOrderByTime(curList []dy.LiveCurProduct) []dy.LiveCurProduct {
	sort.Sort(ProductCurSortList(curList))
	return curList
}

func (I ProductCurSortList) Len() int {
	return len(I)
}

func (I ProductCurSortList) Less(i, j int) bool {
	return I[i].StartTime < I[j].StartTime
}

func (I ProductCurSortList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//直播商品销量按时间排序
type RoomProductTrendSortList []entity.DyRoomProductTrend

func RoomProductTrendOrderByTime(trendList []entity.DyRoomProductTrend) []entity.DyRoomProductTrend {
	sort.Sort(RoomProductTrendSortList(trendList))
	return trendList
}

func (I RoomProductTrendSortList) Len() int {
	return len(I)
}

func (I RoomProductTrendSortList) Less(i, j int) bool {
	return I[i].CrawlTime < I[j].CrawlTime
}

func (I RoomProductTrendSortList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//口碑日榜按时间排序
type ReputationDtScoreList []entity.DyReputationDateScoreList

func ReputationDtScoreListOrderByTime(list []entity.DyReputationDateScoreList) []entity.DyReputationDateScoreList {
	sort.Sort(ReputationDtScoreList(list))
	return list
}

func (I ReputationDtScoreList) Len() int {
	return len(I)
}

func (I ReputationDtScoreList) Less(i, j int) bool {
	return I[i].Date < I[j].Date
}

func (I ReputationDtScoreList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//直播粉丝数推流按时间排序
type LiveFansTrendsList []entity.LiveFollowerCountTrends

func LiveFansTrendsListOrderByTime(list []entity.LiveFollowerCountTrends) []entity.LiveFollowerCountTrends {
	sort.Sort(LiveFansTrendsList(list))
	return list
}

func (I LiveFansTrendsList) Len() int {
	return len(I)
}

func (I LiveFansTrendsList) Less(i, j int) bool {
	return I[i].CrawlTime < I[j].CrawlTime
}

func (I LiveFansTrendsList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//直播粉丝数推流按时间排序
type LiveClubFansTrendsList []entity.LiveAnsClubCountTrends

func LiveClubFansTrendsListOrderByTime(list []entity.LiveAnsClubCountTrends) []entity.LiveAnsClubCountTrends {
	sort.Sort(LiveClubFansTrendsList(list))
	return list
}

func (I LiveClubFansTrendsList) Len() int {
	return len(I)
}

func (I LiveClubFansTrendsList) Less(i, j int) bool {
	return I[i].CrawlTime < I[j].CrawlTime
}

func (I LiveClubFansTrendsList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//直播粉丝数推流按时间排序
type ProductPriceTrendsList []entity.DyProductPriceTrend

func ProductPriceTrendsListOrderByTime(list []entity.DyProductPriceTrend) []entity.DyProductPriceTrend {
	sort.Sort(ProductPriceTrendsList(list))
	return list
}

func (I ProductPriceTrendsList) Len() int {
	return len(I)
}

func (I ProductPriceTrendsList) Less(i, j int) bool {
	return I[i].StartTime < I[j].StartTime
}

func (I ProductPriceTrendsList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
