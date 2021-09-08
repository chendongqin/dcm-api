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
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	list, total, comErr := es.NewEsAuthorBusiness().SimpleSearch(nickname, keyword, tags, secondTags, page, pageSize)
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
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	list, total, comErr := es.NewEsProductBusiness().InternalSearch(productId, title, platformLabel, dcmLevelFirst, firstCname, secondCname, thirdCname, page, pageSize)
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
		receiver.FailReturn(global.NewError(5000))
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
	mediaID, url, err := business.NewWechatBusiness().AddMedia(material.MediaTypeImage, path)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	if comErr := business.NewFileBusiness().InsertFile(fileName, url, tag, mediaID); comErr != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(nil)
}
