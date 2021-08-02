package dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/models/business"
	"dongchamao/models/business/es"
	"dongchamao/models/hbase"
	entity2 "dongchamao/models/hbase/entity"
	"dongchamao/structinit/repost/dy"
	"time"
)

type AuthorController struct {
	controllers.ApiBaseController
}

//达人分类
func (receiver *AuthorController) AuthorCate() {
	configBusiness := business.NewConfigBusiness()
	cateJson := configBusiness.GetConfigJson("author_cate", true)
	cate := business.DealAuthorCateJson(cateJson)
	receiver.SuccReturn(cate)
	return
}

//达人数据
func (receiver *AuthorController) AuthorBaseData() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	authorBase, comErr := authorBusiness.HbaseGetAuthor(authorId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	reputation, comErr := authorBusiness.HbaseGetAuthorReputation(authorId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	fansClub, _ := hbase.GetAuthorFansClub(authorId)
	basic, _ := hbase.GetAuthorBasic(authorId, "")
	returnMap := map[string]interface{}{
		"author_base": authorBase,
		"reputation": dy.RepostSimpleReputation{
			Score:         reputation.Score,
			Level:         reputation.Level,
			EncryptShopID: reputation.EncryptShopID,
			ShopName:      reputation.ShopName,
			ShopLogo:      reputation.ShopLogo,
		},
		"fans_club": fansClub.TotalFansCount,
		"rank":      nil,
		"basic":     basic,
	}
	receiver.SuccReturn(returnMap)
	return
}

//星图指数数据
func (receiver *AuthorController) AuthorStarSimpleData() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	returnMap := map[string]interface{}{
		"has_star_detail": false,
		"star_detail":     nil,
	}
	authorBusiness := business.NewAuthorBusiness()
	xtDetail, comErr := hbase.GetXtAuthorDetail(authorId)
	if comErr == nil {
		returnMap["has_star_detail"] = true
		returnMap["star_detail"] = authorBusiness.GetDyAuthorScore(xtDetail.LiveScore, xtDetail.Score)
	}
	receiver.SuccReturn(returnMap)
	return
}

//达人口碑
func (receiver *AuthorController) Reputation() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	reputation, comErr := authorBusiness.HbaseGetAuthorReputation(authorId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"reputation": reputation,
	})
	return
}

//达人视频概览
func (receiver *AuthorController) AuthorAwemesByDay() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	startDay := receiver.Ctx.Input.Param(":start")
	endDay := receiver.Ctx.Input.Param(":end")
	if authorId == "" || startDay == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if endDay == "" {
		endDay = time.Now().Format("2006-01-02")
	}
	aABusiness := business.NewAuthorAwemeBusiness()
	pslTime := "2006-01-02"
	t1, err := time.ParseInLocation(pslTime, startDay, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	t2, err := time.ParseInLocation(pslTime, endDay, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	//时间限制
	if t1.After(t2) || t2.After(t1.AddDate(0, 0, 90)) || t2.After(time.Now()) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	videoOverview := aABusiness.HbaseGetVideoAggRangeDate(authorId, t1, t2)
	receiver.SuccReturn(map[string]interface{}{
		"video_overview": videoOverview,
	})
	return
}

//基础数据趋势图
func (receiver *AuthorController) AuthorBasicChart() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	startDay := receiver.Ctx.Input.Param(":start")
	endDay := receiver.Ctx.Input.Param(":end")
	if authorId == "" || startDay == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if endDay == "" {
		endDay = time.Now().Format("2006-01-02")
	}
	pslTime := "2006-01-02"
	t1, err := time.ParseInLocation(pslTime, startDay, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	t2, err := time.ParseInLocation(pslTime, endDay, time.Local)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if t1.After(t2) || t2.After(t1.AddDate(0, 0, 90)) || t2.After(time.Now()) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	data, comErr := authorBusiness.HbaseGetAuthorBasicRangeDate(authorId, t1, t2)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(data)
	return
}

//粉丝分布分析
func (receiver *AuthorController) AuthorFansAnalyse() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	detail, comErr := hbase.GetXtAuthorDetail(authorId)
	data := map[string][]entity2.XtDistributionsList{}
	if comErr == nil {
		for _, v := range detail.Distributions {
			name := ""
			switch v.Type {
			case entity2.XtGenderDistribution:
				name = "gender"
			case entity2.XtCityDistribution:
				name = "city"
			case entity2.XtAgeDistribution:
				name = "age"
			case entity2.XtProvinceDistribution:
				name = "province"
			default:
				continue
			}
			data[name] = v.DistributionList
		}
	} else {
		data["gender"] = []entity2.XtDistributionsList{}
		data["city"] = []entity2.XtDistributionsList{}
		data["age"] = []entity2.XtDistributionsList{}
		data["province"] = []entity2.XtDistributionsList{}
	}
	data["active_day"] = []entity2.XtDistributionsList{}
	data["active_week"] = []entity2.XtDistributionsList{}
	var countCity int64 = 0
	var countPro int64 = 0
	for _, v := range data["city"] {
		countCity += v.DistributionValue
	}
	for _, v := range data["province"] {
		countPro += v.DistributionValue
	}
	for k, v := range data["city"] {
		data["city"][k].DistributionPer = float64(v.DistributionValue) / float64(countCity)
	}
	for k, v := range data["province"] {
		data["province"][k].DistributionPer = float64(v.DistributionValue) / float64(countPro)
	}
	receiver.SuccReturn(data)
	return
}

//达人直播分析
func (receiver *AuthorController) CountLiveRoomAnalyse() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	t1, t2, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	data := authorBusiness.CountLiveRoomAnalyse(authorId, t1, t2)
	receiver.SuccReturn(data)
	return
}

//达人直播间列表
func (receiver *AuthorController) AuthorLiveRooms() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	InputData := receiver.InputFormat()
	keyword := InputData.GetString("keyword", "")
	sortStr := InputData.GetString("sort", "create_timestamp")
	orderBy := InputData.GetString("order_by", "desc")
	page := InputData.GetInt("page", 1)
	size := InputData.GetInt("page_size", 10)
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	t1, t2, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	esLiveBusiness := es.NewEsLiveBusiness()
	list, total, comErr := esLiveBusiness.SearchAuthorRooms(authorId, keyword, sortStr, orderBy, page, size, t1, t2)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
	return
}

//达人电商分析
func (receiver *AuthorController) AuthorProductAnalyse() {
	authorId := receiver.GetString(":author_id")
	startTime, endTime, comErr := receiver.GetRangeDate()
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	keyword := receiver.GetString("keyword", "")
	firstCate := receiver.GetString("first_cate", "")
	secondCate := receiver.GetString("second_cate", "")
	thirdCate := receiver.GetString("third_cate", "")
	brandName := receiver.GetString("brand_name", "")
	sortStr := receiver.GetString("sort", "")
	orderBy := receiver.GetString("order_by", "")
	shopType, _ := receiver.GetInt("shop_type", 0)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	if brandName == "不限" {
		brandName = ""
	}
	if firstCate == "不限" {
		firstCate = ""
	}
	authorBusiness := business.NewAuthorBusiness()
	list, analysisCount, cateList, brandList, total, comErr := authorBusiness.GetAuthorProductAnalyse(authorId, keyword, firstCate, secondCate, thirdCate, brandName, sortStr, orderBy, shopType, startTime, endTime, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":           list,
		"cate_list":      cateList,
		"brand_list":     brandList,
		"analysis_count": analysisCount,
		"total":          total,
	})
	return
}
