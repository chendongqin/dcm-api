package controllers

import (
	"dongchamao/business"
	"dongchamao/business/es"
	"dongchamao/global"
	"dongchamao/global/cache"
	utils2 "dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	"dongchamao/models/entity"
	"dongchamao/models/repost"
	"dongchamao/services/dyimg"
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"os"
	"strconv"
	"strings"
	"time"
)

type InternalController struct {
	ApiBaseController
}

func (receiver *InternalController) AuthorSearch() {
	nickname := receiver.GetString("nickname", "")
	keyword := receiver.GetString("keyword", "")
	tags := receiver.GetString("tags", "")
	secondTags := receiver.GetString("second_tags", "")
	minFollower, _ := receiver.GetInt64("min_follower", 0)
	maxFollower, _ := receiver.GetInt64("max_follower", 0)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	list, total, comErr := es.NewEsAuthorBusiness().SimpleSearch("", nickname, keyword, tags, secondTags, minFollower, maxFollower, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	for k, v := range list {
		list[k].Avatar = dyimg.Fix(v.Avatar)
		if v.UniqueId == "" || v.UniqueId == "0" {
			list[k].UniqueId = v.ShortId
		}
	}
	if total > 10000 {
		total = 10000
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
	return
}

func (receiver *InternalController) AuthorInfo() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	data, err := hbase.GetAuthor(authorId)
	if err != nil {
		receiver.FailReturn(err)
		return
	}
	authorData := entity.DyAuthorSimple{}
	utils2.MapToStruct(data, &authorData)
	receiver.SuccReturn(authorData)
	return
}

//修改达人分类
func (receiver *InternalController) ChangeAuthorCate() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	tags := receiver.InputFormat().GetString("tags", "")
	tagsTow := receiver.InputFormat().GetString("tags_two", "")
	if tags == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	dirtyBusiness := business.NewDirtyBusiness()
	comErr := dirtyBusiness.ChangeAuthorCate(authorId, tags, tagsTow)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(nil)
	return
}

//获取商品列表
func (receiver *InternalController) ProductSearch() {
	productId := receiver.GetString("product_id")
	title := receiver.GetString("title")
	platformLabel := receiver.GetString("platform_label")
	dcmLevelFirst := receiver.GetString("dcm_level_first", "")
	firstCname := receiver.GetString("first_cname", "")
	secondCname := receiver.GetString("second_cname", "")
	thirdCname := receiver.GetString("third_cname", "")
	minPrice, _ := receiver.GetFloat("min_price", 0)
	maxPrice, _ := receiver.GetFloat("max_price", 0)
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	list, total, comErr := es.NewEsProductBusiness().SimpleSearch(productId, title, platformLabel, dcmLevelFirst, firstCname, secondCname, thirdCname, minPrice, maxPrice, page, pageSize)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	if total > 10000 {
		total = 10000
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  list,
		"total": total,
	})
	return
}

//修改商品分类
func (receiver *InternalController) ChangeProductCate() {
	productId := receiver.Ctx.Input.Param(":product_id")
	if productId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	dcmLevelFirst := receiver.InputFormat().GetString("dcm_level_first", "")
	firstCate := receiver.InputFormat().GetString("first_cname", "")
	secondCate := receiver.InputFormat().GetString("second_cname", "")
	thirdCate := receiver.InputFormat().GetString("third_cname", "")
	if dcmLevelFirst == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	dirtyBusiness := business.NewDirtyBusiness()
	comErr := dirtyBusiness.ChangeProductCate(productId, dcmLevelFirst, firstCate, secondCate, thirdCate)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(nil)
	return
}

//删除缓存
func (receiver *InternalController) ClearCache() {
	input := receiver.InputFormat()
	cacheType := input.GetString("cacheType", "")
	if cacheType == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	val := input.GetString("val", "")
	if val == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var cacheKey string
	switch cacheType {
	case "userInfo":
		{
			userId, _ := strconv.Atoi(val)
			cacheKey = cache.GetCacheKey(cache.UserInfo, userId)
			break
		}
	case "userLevel":
		{
			println(val)
			var userLevel repost.UserLevelCache
			json.Unmarshal([]byte(val), &userLevel)
			cacheKey = cache.GetCacheKey(cache.UserLevel, userLevel.UserId, userLevel.Platform)
			break
		}
	case "configKey":
		{
			cacheKey = cache.GetCacheKey(cache.ConfigKeyCache, val)
			break
		}
	case "cate":
		{
			cacheKey = cache.GetCacheKey(cache.LongTimeConfigKeyCache)
			break
		}
	}
	global.Cache.Delete(cacheKey)
	receiver.SuccReturn(nil)
	return
}

func (receiver *InternalController) GetConfig() {
	var configJson dcm.DcConfigJson
	keyName := receiver.GetString(":key_name")
	exist, err := dcm.GetBy("key_name", keyName, &configJson)
	if !exist || err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	if configJson.Auth == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(configJson.Value), &data); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(data)
	return
}

func (receiver *InternalController) GetWeChatMenu() {
	data, err := business.NewWechatBusiness().GetMenus()
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
	}
	receiver.SuccReturn(data.Menu)
	return
}

func (receiver *InternalController) SetWeChatMenu() {
	input := receiver.InputFormat()
	menu := input.GetString("menu", "")
	business.NewMonitorBusiness().SendErr("内部接口更新菜单", "SetWeChatMenu")
	if menu == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var menuMap map[string]interface{}
	if err := json.Unmarshal([]byte(menu), &menuMap); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	if err := business.NewWechatBusiness().UpdateMenus(menuMap); err != nil {
		receiver.FailReturn(global.NewCommonError(err))
		return
	}
	receiver.SuccReturn(nil)
	return
}

func (receiver *InternalController) UploadWeChatMedia() {
	file, fileHeader, err := receiver.GetFile("file")
	var buff = make([]byte, fileHeader.Size)
	if _, err = file.Read(buff); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	s := strings.Split(fileHeader.Filename, ".")
	tag := s[len(s)-1]
	fileName := time.Now().Format("20060102101020") + utils.RandStringBytes(6) + "." + tag
	localDir := "./temp"
	_, err = os.Stat(localDir)
	if os.IsNotExist(err) {
		utils2.MakeDir(localDir)
	}
	path := localDir + "/" + fileName
	osFile, err := os.Create(path)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	//移除临时文件
	defer os.Remove(path)
	if _, err := osFile.Write(buff); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	_, _, err = business.NewWechatBusiness().AddMedia(material.MediaTypeImage, path)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(nil)
}

func (receiver *InternalController) GetWeChatMediaList() {
	mediaType := receiver.GetString("media_type")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	from := int64((page - 1) * pageSize)
	to := int64(page*pageSize - 1)
	list := business.NewWechatBusiness().GetMediaList(material.PermanentMaterialType(mediaType), from, to)
	receiver.SuccReturn(map[string]interface{}{"list": list.Item, "page": page, "pageSize": pageSize, "total": list.TotalCount})
	return
}

func (receiver *InternalController) DelWeChatMedia() {
	mediaId := receiver.GetString("media_id")
	receiver.SuccReturn(business.NewWechatBusiness().DelMedia(mediaId))
	return
}

//id加解密
func (receiver *InternalController) IdEncryptDecrypt() {
	id := receiver.Ctx.Input.Param(":id")
	id1 := ""
	if strings.Index(id, "=") < 0 {
		id1 = business.IdEncrypt(id)
	}
	id2 := business.IdDecrypt(id)
	receiver.SuccReturn(map[string]string{
		"id":      id,
		"encrypt": id1,
		"decrypt": id2,
	})
	return
}

//json解密
func (receiver *InternalController) JsonDecrypt() {
	str := receiver.InputFormat().GetString("str", "")
	decryptStr := business.JsonDecrypt(str)
	receiver.SuccReturn(map[string]interface{}{
		"decrypt_str": decryptStr,
	})
	return
}

//红人看板加速
func (receiver *InternalController) SpiderLiveSpeedUp() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	author, err := hbase.GetAuthor(authorId)
	if err != nil {
		receiver.FailReturn(err)
		return
	}
	if global.IsDev() {
		receiver.SuccReturn(nil)
		return
	}
	business.NewSpiderBusiness().AddLive(authorId, author.FollowerCount, business.AddLiveTopStar)
	_, _ = dcm.Insert(nil, &dcm.DcLiveSpiderLogs{
		AuthorId:   authorId,
		Top:        business.AddLiveTopStar,
		AddLog:     "red_author",
		CreateTime: time.Now(),
	})
	receiver.SuccReturn(nil)
	return
}

//通用的url日志处理
func (receiver *InternalController) CommonUrlLog() {
	safe := business.NewSafeBusiness()
	res := safe.CommonAnalyseLogs()
	receiver.SuccReturn(map[string]interface{}{
		"list": res,
	})
	return
}

/**根据日志筛选每小时需要加速的达人，直播，商品*/
func (receiver *InternalController) SpeedUp() {
	days := receiver.Ctx.Input.Param(":days")
	daysInt := utils2.ToInt(days)
	if !utils2.InArrayInt(daysInt, []int{0, 7}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	safe := business.NewSafeBusiness()
	var data map[string][]string
	data = make(map[string][]string)
	keys := []string{"speed_author", "speed_live", "speed_product"}
	for _, v := range keys {
		data[v] = safe.SpeedFilterLog(v, daysInt)
	}
	receiver.SuccReturn(data)
	//safe.SpeedFilterLog("speed_live",daysInt)
}
