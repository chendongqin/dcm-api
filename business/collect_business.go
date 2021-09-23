package business

import (
	"dongchamao/business/es"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	es2 "dongchamao/models/es"
	"dongchamao/models/repost"
	"dongchamao/services/dyimg"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type CollectBusiness struct {
}

func NewCollectBusiness() *CollectBusiness {
	return new(CollectBusiness)
}

//获取抖音收藏
func (receiver *CollectBusiness) GetDyCollect(tagId, collectType int, keywords, label string, userId, page, pageSize int) (interface{}, int64, global.CommonError) {
	var (
		total    int64
		comErr   global.CommonError
		collects []dcm.DcUserDyCollect
	)
	dbCollect := dcm.GetDbSession()
	defer dbCollect.Close()
	var query string
	query = fmt.Sprintf("collect_type=%v", collectType)
	if tagId != 0 {
		query = fmt.Sprintf("tag_id=%v", tagId)
	}
	if keywords != "" {
		query += " AND (unique_id LIKE '%" + keywords + "%' or nickname LIKE '%" + keywords + "%')"
	}
	if label != "" {
		query += " AND FIND_IN_SET('" + label + "',label)"
	}
	query += " AND user_id=" + strconv.Itoa(userId) + " AND status=1"
	err := dbCollect.Table(dcm.DcUserDyCollect{}).Where(query).Limit(pageSize, (page-1)*pageSize).OrderBy("update_time desc").Find(&collects)
	if total, err = dbCollect.Table(dcm.DcUserDyCollect{}).Where(query).Count(); err != nil {
		comErr = global.NewError(5000)
		return nil, total, comErr
	}
	switch collectType {
	case 1:
		data := make([]repost.CollectAuthorRet, len(collects))
		startTime := time.Now().AddDate(0, 0, -31)
		yesterday := time.Now().AddDate(0, 0, -1)
		var authorIds []string
		for k, v := range collects {
			data[k].DcUserDyCollect = v
			data[k].DcUserDyCollect.CollectId = IdEncrypt(v.CollectId)
			dyAuthor, _ := hbase.GetAuthor(v.CollectId)
			basicData, _ := hbase.GetAuthorBasic(v.CollectId, time.Now().AddDate(0, 0, -1).Format("20060102"))
			YesBasicData, _ := hbase.GetAuthorBasic(v.CollectId, time.Now().AddDate(0, 0, -2).Format("20060102"))
			data[k].FollowerCount = basicData.FollowerCount
			data[k].FollowerIncreCount = basicData.FollowerCount - YesBasicData.FollowerCount
			data[k].Avatar = dyimg.Avatar(dyAuthor.Data.Avatar)
			authorIds = append(authorIds, v.CollectId)
		}
		liveSumData := es.NewEsLiveBusiness().SumDataByAuthors(authorIds, startTime, yesterday)
		for k, v := range collects {
			authorBase, comErr := NewAuthorBusiness().HbaseGetAuthor(v.CollectId)
			if comErr != nil {
				return nil, 0, comErr
			}
			//近30日场均销售额
			data[k].Predict7Gmv = utils.FriendlyFloat64(liveSumData[v.CollectId].TotalGmv.Avg)
			//近30日视频平均点赞
			if authorBase.AwemeCount != 0 {
				data[k].Predict7Digg = float64(authorBase.DiggCount) / float64(authorBase.AwemeCount)
			}
		}
		return data, total, nil
	case 2:
		data := make([]repost.CollectProductRet, len(collects))
		var productIds []string
		for k, v := range collects {
			v.CollectId = IdEncrypt(v.CollectId)
			data[k].DcUserDyCollect = v
			productIds = append(productIds, v.CollectId)
		}
		products, _, commonError := es.NewEsProductBusiness().SearchProducts(productIds)
		var productMap = make(map[string]es2.DyProduct, len(products))
		for _, v := range products {
			v.ProductId = IdEncrypt(v.ProductId)
			v.Image = dyimg.Fix(v.Image)
			productMap[v.ProductId] = v
		}
		var labelUpdate = make(map[string]string)
		for k, v := range data {
			productInfo := productMap[v.CollectId]
			data[k].ProductId = productInfo.ProductId
			data[k].Image = productInfo.Image
			data[k].Nickname = productInfo.Title
			data[k].Price = productInfo.Price
			data[k].CouponPrice = productInfo.CouponPrice
			data[k].Pv = productInfo.Pv
			data[k].OrderAccount = productInfo.OrderAccount
			data[k].WeekOrderAccount = productInfo.MonthOrderAccount
			data[k].PlatformLabel = productInfo.PlatformLabel
			data[k].Undercarriage = productInfo.Undercarriage
			data[k].IsCoupon = productInfo.IsCoupon
			data[k].WeekRelateAuthor = productInfo.RelateAuthor
			//获取分类更新商品
			if data[k].Label != productInfo.DcmLevelFirst {
				labelUpdate[productInfo.ProductId] = productInfo.DcmLevelFirst
			}
		}
		go func() {
			//更新商品分类
			for k, v := range labelUpdate {
				if _, err := dcm.GetDbSession().Table(dcm.DcUserDyCollect{}).Where("collect_id=? and collect_type=?", k, collectType).Update(map[string]string{"label": v}); err != nil {
					log.Println("collect_product_label:", err.Error())
					return
				}
			}
		}()
		return data, total, commonError
	case 3:
		data := make([]repost.CollectAwemeRet, len(collects))
		for k, v := range collects {
			awemeBase, comErr := hbase.GetVideo(v.CollectId)
			if comErr != nil {
				return nil, 0, comErr
			}
			awemeAuthor, comErr := hbase.GetAuthor(awemeBase.Data.AuthorID)
			if comErr != nil {
				return nil, 0, comErr
			}
			v.CollectId = IdEncrypt(v.CollectId)
			yesData, _ := hbase.GetVideoCountData(v.CollectId, time.Now().AddDate(0, 0, -1).Format("20060102"))
			beforeYesData, _ := hbase.GetVideoCountData(v.CollectId, time.Now().AddDate(0, 0, -2).Format("20060102"))
			data[k].DcUserDyCollect = v
			data[k].AwemeAuthorID = IdEncrypt(awemeBase.Data.AuthorID)
			data[k].AwemeCover = awemeBase.Data.AwemeCover
			data[k].AwemeTitle = awemeBase.AwemeTitle
			data[k].AwemeCreateTime = awemeBase.Data.AwemeCreateTime
			data[k].AwemeURL = awemeBase.Data.AwemeURL
			data[k].DiggCount = beforeYesData.DiggCount
			data[k].DiggCountIncr = yesData.DiggCount - beforeYesData.DiggCount
			data[k].AuthorAvatar = dyimg.Fix(awemeAuthor.Data.Avatar)
			data[k].AuthorNickname = awemeAuthor.Data.Nickname
		}
		return data, total, nil
	}
	return nil, 0, nil
}

//获取分组收藏数量
func (receiver *CollectBusiness) GetDyCollectCount(userId, collectType int) (data []repost.CollectCount, comErr global.CommonError) {
	dbCollect := dcm.GetDbSession()
	defer dbCollect.Close()
	if err := dbCollect.Table(dcm.DcUserDyCollect{}).Where("user_id=? AND status=1 AND collect_type=?", userId, collectType).Select("tag_id,count(collect_id) as count").GroupBy("tag_id").Find(&data); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//获取已收藏达人标签
func (receiver *CollectBusiness) GetDyCollectLabel(userId, collectType int) (data []string, comErr global.CommonError) {
	dbCollect := dcm.GetDbSession()
	defer dbCollect.Close()
	if err := dbCollect.Table(dcm.DcUserDyCollect{}).Where("user_id=? AND label<>'' AND status=1 AND collect_type=?", userId, collectType).Select("label").Find(&data); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//收藏达人
func (receiver *CollectBusiness) AddDyCollect(collectId string, collectType, tagId, userId int) (comErr global.CommonError) {
	collect := dcm.DcUserDyCollect{}
	dbCollect := dcm.GetDbSession().Table(collect)
	defer dbCollect.Close()
	exist, err := dbCollect.Where("user_id=? AND collect_type=? AND collect_id=?", userId, collectType, collectId).Get(&collect)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	if collect.Status == 1 {
		comErr = global.NewMsgError("您已收藏该达人，请刷新重试")
		return comErr
	}
	collect.Status = 1
	collect.TagId = tagId
	collect.CollectId = collectId
	collect.UpdateTime = time.Now()
	switch collectType {
	case 1:
		//达人
		author, comErr := hbase.GetAuthor(collectId)
		if comErr != nil {
			return comErr
		}
		collect.Label = strings.Replace(author.Tags, "|", ",", -1)
		collect.UniqueId = author.Data.UniqueID
		collect.Nickname = author.Data.Nickname
		break
	case 2:
		//商品
		info, comErr := hbase.GetProductInfo(collectId)
		if comErr != nil {
			return comErr
		}
		collect.Nickname = info.Title
		collect.Label = info.DcmLevelFirst
	case 3:
		//视频
	}
	if exist {
		if _, err := dbCollect.ID(collect.Id).Update(&collect); err != nil {
			comErr = global.NewError(5000)
			return
		}
	} else {
		collect.CreateTime = time.Now()
		collect.UserId = userId
		collect.CollectType = collectType
		if _, err := dbCollect.Insert(&collect); err != nil {
			comErr = global.NewError(5000)
			return
		}
	}
	return
}

func (receiver *CollectBusiness) DyListCollect(collectType int, userId int, ids []string) (collectMap map[string]int, comErr global.CommonError) {
	var collectList []dcm.DcUserDyCollect
	if err := dcm.GetDbSession().Where("collect_type=? and user_id=? and status=1", collectType, userId).In("collect_id", ids).Find(&collectList); err != nil {
		return nil, global.NewCommonError(err)
	}
	collectMap = make(map[string]int, len(collectList))
	for _, v := range collectList {
		collectMap[v.CollectId] = v.Id
	}
	return
}

//取消收藏
func (receiver *CollectBusiness) CancelDyCollect(id, userId int) (comErr global.CommonError) {
	dbCollect := dcm.GetDbSession().Table(dcm.DcUserDyCollect{})
	defer dbCollect.Close()
	exist, err := dbCollect.Where("id=? and status=? and user_id=?", id, 1, userId).Exist()
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	if !exist {
		comErr = global.NewMsgError("您未收藏该达人，请刷新重试")
		return
	}
	affect, err := dcm.UpdateInfo(dbCollect, id, map[string]interface{}{"status": 0, "update_time": time.Now()}, new(dcm.DcUserDyCollect))
	if err != nil || affect == 0 {
		comErr = global.NewError(5000)
		return
	}
	return
}

//修改收藏分组
func (receiver *CollectBusiness) UpdCollectTag(id, tagId int) (comErr global.CommonError) {
	dbCollect := dcm.GetDbSession().Table(dcm.DcUserDyCollect{})
	defer dbCollect.Close()
	affect, err := dcm.UpdateInfo(dbCollect, id, map[string]interface{}{"tag_id": tagId, "update_time": time.Now()}, new(dcm.DcUserDyCollect))
	if err != nil || affect == 0 {
		comErr = global.NewError(5000)
		return
	}
	return
}

//获取分组列表
func (receiver *CollectBusiness) GetDyCollectTags(userId int) (tags []dcm.DcUserDyCollectTag, comErr global.CommonError) {
	db := dcm.GetDbSession().Table(dcm.DcUserDyCollectTag{})
	if err := db.Where("user_id=? and delete_time is NULL", userId).Find(&tags); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//创建分组
func (receiver *CollectBusiness) AddDyCollectTag(userId int, name string) (comErr global.CommonError) {
	db := dcm.GetDbSession().Table(dcm.DcUserDyCollectTag{})
	by, err := db.Where("name=? and delete_time is null", name).Exist()
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	if by {
		comErr = global.NewMsgError("分组已存在")
		return
	}
	tag := dcm.DcUserDyCollectTag{
		Name:       name,
		UserId:     userId,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if _, err := db.Insert(tag); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//编辑分组
func (receiver *CollectBusiness) UpdDyCollectTag(id int, name string) (comErr global.CommonError) {
	db := dcm.GetDbSession().Table(dcm.DcUserDyCollectTag{})
	if _, err := db.Where("id=?", id).Update(map[string]interface{}{"name": name}); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//删除分组
func (receiver *CollectBusiness) DelDyCollectTag(id int) (comErr global.CommonError) {
	db := dcm.GetDbSession().Table(dcm.DcUserDyCollectTag{})
	if _, err := db.Where("id=?", id).Update(map[string]interface{}{"delete_time": time.Now()}); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}
