package controllers

import (
	"dongchamao/business"
	"dongchamao/business/es"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/models/dcm"
	"dongchamao/models/repost"
	"dongchamao/services/dyimg"
	"encoding/json"
	"strconv"
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
	dcmLevelFirst := receiver.GetString("dcm_level_first", "")
	firstCname := receiver.GetString("first_cname", "")
	secondCname := receiver.GetString("second_cname", "")
	thirdCname := receiver.GetString("third_cname", "")
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 100)
	list, total, comErr := es.NewEsProductBusiness().InternalSearch(productId, title, dcmLevelFirst, firstCname, secondCname, thirdCname, page, pageSize)
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
	firstCate := receiver.InputFormat().GetString("first_cate", "")
	secondCate := receiver.InputFormat().GetString("second_cate", "")
	thirdCate := receiver.InputFormat().GetString("third_cate", "")
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
	data, _ := business.NewWechatBusiness().GetMenus()
	//fmt.Printf("data:%+v\n",data)
	//if err != nil {
	//	receiver.FailReturn(global.NewError(5000))
	//}
	str := "{\n    \"button\": [\n        {\n            \"name\": \"抖音数据\",\n            \"sub_button\": [\n                {\n                    \"type\": \"view\",\n                    \"name\": \"搜索\",\n                    \"url\": \"https://douyin.dongchamao.com\"\n                },\n                {\n                    \"type\": \"view\",\n                    \"name\": \"行业成交大盘\",\n                    \"url\": \"https://douyin.dongchamao.com\"\n                },\n                {\n                    \"type\": \"view\",\n                    \"name\": \"直播带货榜\",\n                    \"url\": \"https://douyin.dongchamao.com/#/talent/ranking/commerce\"\n                },\n                {\n                    \"type\": \"view\",\n                    \"name\": \"红人直播\",\n                    \"url\": \"https://douyin.dongchamao.com/#/report/\"\n                },\n                {\n                    \"type\": \"view\",\n                    \"name\": \"APP下载\",\n                    \"url\": \"https://douyin.dongchamao.com\"\n                }\n           ]\n        },\n        {\n            \"name\": \"资源图书馆\",\n            \"sub_button\": [\n                {\n                    \"type\": \"click\",\n                    \"name\": \"功能上新\",\n                    \"key\": \"new_function\"\n                },\n                {\n                    \"type\": \"view_limited\",\n                    \"name\": \"运营干货\",\n                    \"media_id\": \"UQotsct61iSBt9RSogbMYyLZ60ZbdETWeWv5kzX5ML0\"\n                }\n            ]\n        },\n        {\n            \"name\": \"资源合作\",\n            \"sub_button\": [\n                {\n                    \"type\": \"click\",\n                    \"name\": \"免费试用\",\n                    \"key\": \"try_use\"\n                },\n                {\n                    \"type\": \"click\",\n                    \"name\": \"商务合作\",\n                    \"key\": \"business\"\n                }\n            ]\n        }\n    ]\n}"
	json.Unmarshal([]byte(str), &data.Menu)
	receiver.SuccReturn(data.Menu)
	return
}

func (receiver *InternalController) SetWeChatMenu() {
	input := receiver.InputFormat()
	menu := input.GetString("menu", "")
	if menu == "" {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	println(menu)
	receiver.SuccReturn(menu)
	return
}
