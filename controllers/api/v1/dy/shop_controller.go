package dy

import (
	"dongchamao/business"
	"dongchamao/business/es"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/hbase"
	"dongchamao/models/entity"
	"dongchamao/models/repost/dy"
	dy2 "dongchamao/models/repost/dy"
	"dongchamao/services/dyimg"
	"math"
	"sort"
	"time"
)

type ShopController struct {
	controllers.ApiBaseController
}

func (receiver *ShopController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseLoginMinShowNum, business.DyJewelBaseShowNum)
}

//小店基本数据
func (receiver *ShopController) SearchBase() {
	keyword := receiver.GetString("keyword", "")
	category := receiver.GetString("category", "")
	secondCategory := receiver.GetString("second_category", "")
	thirdCategory := receiver.GetString("third_category", "")
	min30Sales, _ := receiver.GetInt64("min_30sales", 0)
	max30Sales, _ := receiver.GetInt64("max_30sales", 0)
	min30Gmv, _ := receiver.GetFloat("min_30gmv", 0)
	max30Gmv, _ := receiver.GetFloat("max_30gmv", 0)
	min30UnitPrice, _ := receiver.GetFloat("min_30unit_price", 0)
	max30UnitPrice, _ := receiver.GetFloat("max_30unit_price", 0)
	minScore, _ := receiver.GetFloat("min_score", 0)
	maxScore, _ := receiver.GetFloat("max_score", 0)
	isBrand, _ := receiver.GetInt("is_brand", 0)
	isLive, _ := receiver.GetInt("is_live", 0)
	isVideo, _ := receiver.GetInt("is_video", 0)
	sortStr := receiver.GetString("sort", "month_sales")
	orderBy := receiver.GetString("order_by", "desc")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	receiver.KeywordBan(keyword)
	if !receiver.HasLogin && keyword != "" {
		receiver.FailReturn(global.NewError(4001))
		return
	}
	if !receiver.HasAuth {
		if category != "" || secondCategory != "" || thirdCategory != "" || sortStr != "month_sales" || orderBy != "desc" ||
			minScore > 0 || maxScore > 0 || min30Gmv > 0 || max30Gmv > 0 || min30Sales > 0 || max30Sales > 0 || min30UnitPrice > 0 || max30UnitPrice > 0 ||
			isLive == 1 || isVideo == 1 || isBrand == 1 || page != 1 {
			if !receiver.HasLogin {
				receiver.FailReturn(global.NewError(4001))
				return
			}
			receiver.FailReturn(global.NewError(4004))
			return
		}
		if pageSize > receiver.MaxTotal {
			pageSize = receiver.MaxTotal
		}
	}
	formNum := (page - 1) * pageSize
	if formNum > receiver.MaxTotal {
		receiver.FailReturn(global.NewError(4004))
		return
	}
	list, total, comErr := es.NewEsShopBusiness().BaseSearch(keyword, category, secondCategory, thirdCategory,
		min30Sales, max30Sales, min30Gmv, max30Gmv, min30UnitPrice, max30UnitPrice, minScore, maxScore,
		isLive, isBrand, isVideo, page, pageSize,
		sortStr, orderBy)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	shopIds := make([]string, 0)
	for _, v := range list {
		shopIds = append(shopIds, v.ShopId)
	}
	if receiver.HasLogin {
		collectBusiness := business.NewCollectBusiness()
		collect, comErr := collectBusiness.DyListCollect(4, receiver.UserId, shopIds)
		if comErr != nil {
			receiver.FailReturn(comErr)
		}
		for k, v := range list {
			list[k].IsCollect = collect[v.ShopId]
		}
	}
	for k, v := range list {
		list[k].Logo = dyimg.Fix(v.Logo)
		list[k].ShopId = business.IdEncrypt(v.ShopId)
	}
	totalPage := math.Ceil(float64(total) / float64(pageSize))
	maxPage := math.Ceil(float64(receiver.MaxTotal) / float64(pageSize))
	if totalPage > maxPage {
		totalPage = maxPage
	}
	maxTotal := receiver.MaxTotal
	if maxTotal > total {
		maxTotal = total
	}
	business.NewUserBusiness().KeywordsRecord(keyword)
	receiver.SuccReturn(map[string]interface{}{
		"list":       list,
		"total":      total,
		"total_page": totalPage,
		"max_num":    maxTotal,
		"has_auth":   receiver.HasAuth,
		"has_login":  receiver.HasLogin,
	})
	return
}

//小店基本数据
func (receiver *ShopController) ShopBase() {
	var returnRes entity.DyShopBaseBasic
	var comErr global.CommonError
	shopId := business.IdDecrypt(receiver.Ctx.Input.Param(":shop_id"))
	if shopId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	returnRes.BaseData, comErr = hbase.GetShop(shopId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	shopDetailData, comErr := hbase.GetShopDetailByDate(shopId, time.Now().Format("20060102"))
	if comErr != nil { //今天取不到，取昨日数据
		shopDetailData, _ = hbase.GetShopDetailByDate(shopId, time.Now().AddDate(0, 0, -1).Format("20060102"))
	}

	returnRes.DetailData.Sales = shopDetailData.Sales
	returnRes.DetailData.Gmv = shopDetailData.Gmv
	returnRes.DetailData.D30LiveCnt = shopDetailData.D30LiveCnt
	returnRes.DetailData.D30AuthorCnt = shopDetailData.D30AuthorCnt
	returnRes.DetailData.D30AwemeCnt = shopDetailData.D30AwemeCnt
	returnRes.DetailData.D30Sales = shopDetailData.D30Sales
	returnRes.DetailData.D30Gmv = shopDetailData.D30Gmv
	returnRes.DetailData.D30Pct = shopDetailData.D30Pct

	receiver.SuccReturn(returnRes)
	return
}

/**小店数据基础分析 **/
func (receiver *ShopController) ShopBaseAnalysis() {
	shopId := business.IdDecrypt(receiver.GetString(":shop_id", ""))
	if shopId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	info, _ := hbase.GetShopDetailRangDate(shopId, startTime, endTime)
	beginTime := startTime

	var date []string    //日期
	var sale []int64     //销量
	var gmv []float64    //销售额
	var awemeNum []int64 //视频数
	var liveNum []int64  //直播数

	priceSectionMap := map[string]int64{} //价格区间
	goodsCatTopMap := map[string]int64{}  //价格区间
	for {
		if beginTime.After(endTime) {
			break
		}
		dateKey := beginTime.Format("20060102")
		if v, ok := info[dateKey]; ok {
			sale = append(sale, v.Sales)
			gmv = append(gmv, v.Gmv)
			awemeNum = append(awemeNum, v.AwemeNum)
			liveNum = append(liveNum, v.LiveNum)
			for k, num := range v.PriceDist {
				if _, exist := priceSectionMap[k]; exist {
					priceSectionMap[k] += num
				} else {
					priceSectionMap[k] = num
				}
			}
			for k, num := range v.Classifications {
				if _, exist := goodsCatTopMap[k]; exist {
					goodsCatTopMap[k] += num
				} else {
					goodsCatTopMap[k] = num
				}
			}
		} else {
			sale = append(sale, 0)
			gmv = append(gmv, 0)
			awemeNum = append(awemeNum, 0)
			liveNum = append(liveNum, 0)
		}
		date = append(date, beginTime.Format("01/02"))

		beginTime = beginTime.AddDate(0, 0, 1)
	}
	priceSection := make([]dy.NameValueInt64Chart, 0)
	goodsCatTop := make([]dy.NameValueInt64Chart, 0)

	for k, v := range goodsCatTopMap {
		goodsCatTop = append(priceSection, dy.NameValueInt64Chart{
			Name:  k,
			Value: v,
		})
	}
	sort.Slice(goodsCatTop, func(i, j int) bool {
		return goodsCatTop[i].Value > goodsCatTop[j].Value
	})
	if len(goodsCatTop) > 5 {
		goodsCatTop = goodsCatTop[:5]
	}
	var priceMap = map[string]string{
		"lt50":   "0-50",
		"lt100":  "50-100",
		"lt300":  "100-300",
		"lt500":  "300-500",
		"lt1000": "500-1000",
		"gt1000": ">1000",
	}

	for k, v := range priceMap {
		priceSection = append(priceSection, dy.NameValueInt64Chart{
			Name:  priceMap[v],
			Value: priceSectionMap[k],
		})
	}
	receiver.SuccReturn(map[string]interface{}{
		"sales_chart": dy2.ShopSaleChart{
			Date:       date,
			SalesCount: sale,
			GmvCount:   gmv,
		},
		"live_aweme_chart": dy2.ShopLiveAwemeChart{
			Date:       date,
			LiveCount:  liveNum,
			AwemeCount: awemeNum,
		},
		"price_chart":   priceSection,
		"goods_cat_top": goodsCatTop,
	})
	return
}

//达人销售额占比
func (receiver *ShopController) ShopAuthorGmvRate() {
	shopId := business.IdDecrypt(receiver.GetString(":shop_id", ""))
	if shopId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	allTop5, comErr := business.NewShopBusiness().ShopAuthorView(shopId, startTime, endTime)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"top5": allTop5,
	})
	return
}

//小店商品分析
func (receiver *ShopController) ShopProductAnalysis() {
	shopId := business.IdDecrypt(receiver.GetString(":shop_id"))
	startTime, stopTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	category := receiver.GetString("category", "")
	orderBy := receiver.GetString("order_by", "")
	sortStr := receiver.GetString("sort", "")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	list, total, comErr := business.NewShopBusiness().ShopProductAnalysis(shopId, keyword, category, sortStr, orderBy, startTime, stopTime, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

//小店商品分析统计
func (receiver *ShopController) ShopProductAnalysisCount() {
	shopId := business.IdDecrypt(receiver.GetString(":shop_id"))
	startTime, stopTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	count, comErr := business.NewShopBusiness().ShopProductAnalysisCount(shopId, keyword, startTime, stopTime)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"count": count,
	})
}

//小店直播达人分析
func (receiver *ShopController) ShopLiveAuthorAnalysis() {
	shopId := business.IdDecrypt(receiver.Ctx.Input.Param(":shop_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	tag := receiver.GetString("tag", "")
	minFollow, _ := receiver.GetInt64("min_follow", 0)
	maxFollow, _ := receiver.GetInt64("max_follow", 0)
	scoreType, _ := receiver.GetInt("score_type", 5)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	list, total, comErr := business.NewShopBusiness().ShopLiveAuthorAnalysis(shopId, keyword, tag, startTime, endTime, minFollow, maxFollow, scoreType, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		authorInfo, _ := hbase.GetAuthor(v.AuthorId)
		list[k].Avatar = dyimg.Fix(authorInfo.Data.Avatar)
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		list[k].ProductId = business.IdEncrypt(v.ProductId)
		list[k].Nickname = authorInfo.Data.Nickname
		list[k].RoomNum = len(v.RelatedRooms)
		list[k].RelatedRooms = []entity.DyProductAuthorRelatedRoom{}
	}
	maxTotal := total
	if total > business.EsMaxShowNum {
		maxTotal = business.EsMaxShowNum
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":           list,
		"total":          total,
		"max_show_total": maxTotal,
	})
	return
}

//小店直播达人分析统计
func (receiver *ShopController) ShopLiveAuthorAnalysisCount() {
	shopId := business.IdDecrypt(receiver.Ctx.Input.Param(":shop_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	countList, comErr := business.NewShopBusiness().ShopLiveAuthorAnalysisCount(shopId, keyword, startTime, endTime)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"list": countList,
	})
	return
}

//小店达人直播间列表
func (receiver *ShopController) ShopLiveAuthorRooms() {
	shopId := business.IdDecrypt(receiver.Ctx.Input.Param(":shop_id"))
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 5, 10)
	sortStr := receiver.GetString("sort", "start_ts")
	orderBy := receiver.GetString("order_by", "desc")
	list, total := business.NewProductBusiness().ProductAuthorLiveRooms("", shopId, authorId, startTime, endTime, sortStr, orderBy, page, pageSize)
	for k, v := range list {
		list[k].Cover = dyimg.Fix(v.Cover)
		list[k].RoomId = business.IdEncrypt(v.RoomId)
		endLiveTime := v.EndTs
		if endLiveTime == 0 {
			endLiveTime = time.Now().Unix()
		}
		list[k].LiveSecond = endLiveTime - v.StartTs
	}
	maxTotal := total
	if total > business.EsMaxShowNum {
		maxTotal = business.EsMaxShowNum
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":           list,
		"total":          total,
		"max_show_total": maxTotal,
	})
}

//小店视频达人分析
func (receiver *ShopController) ShopAwemeAuthorAnalysis() {
	shopId := business.IdDecrypt(receiver.Ctx.Input.Param(":shop_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	tag := receiver.GetString("tag", "")
	minFollow, _ := receiver.GetInt64("min_follow", 0)
	maxFollow, _ := receiver.GetInt64("max_follow", 0)
	scoreType, _ := receiver.GetInt("score_type", 5)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	list, total, comErr := business.NewShopBusiness().ShopAwemeAuthorAnalysis(shopId, keyword, tag, startTime, endTime, minFollow, maxFollow, scoreType, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		list[k].AuthorId = business.IdEncrypt(v.AuthorId)
		list[k].ProductId = business.IdEncrypt(v.ProductId)
		list[k].Avatar = dyimg.Fix(v.Avatar)
		if v.DisplayId == "" || v.DisplayId == "0" {
			list[k].DisplayId = v.ShortId
		}
	}
	maxTotal := total
	if total > business.EsMaxShowNum {
		maxTotal = business.EsMaxShowNum
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":           list,
		"total":          total,
		"max_show_total": maxTotal,
	})
	return
}

//小店视频达人分析统计
func (receiver *ShopController) ShopAwemeAuthorAnalysisCount() {
	shopId := business.IdDecrypt(receiver.Ctx.Input.Param(":shop_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	productBusiness := business.NewProductBusiness()
	countList, comErr := productBusiness.ProductAwemeAuthorAnalysisCount(shopId, keyword, startTime, endTime)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"list": countList,
	})
	return
}

//小店达人视频列表
func (receiver *ShopController) ShopAuthorAwemes() {
	shopId := business.IdDecrypt(receiver.Ctx.Input.Param(":shop_id"))
	authorId := business.IdDecrypt(receiver.Ctx.Input.Param(":author_id"))
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 5, 10)
	sortStr := receiver.GetString("sort", "aweme_gmv")
	orderBy := receiver.GetString("order_by", "desc")
	list, total := business.NewProductBusiness().ProductAuthorAwemes("", shopId, authorId, startTime, endTime, sortStr, orderBy, page, pageSize)
	for k, v := range list {
		list[k].AwemeCover = dyimg.Fix(v.AwemeCover)
		list[k].AwemeId = business.IdEncrypt(v.AwemeId)
		list[k].AwemeUrl = business.AwemeUrl + v.AwemeId
	}
	maxTotal := total
	if total > business.EsMaxShowNum {
		maxTotal = business.EsMaxShowNum
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":           list,
		"total":          total,
		"max_show_total": maxTotal,
	})
}
