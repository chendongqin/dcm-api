package controllers

import (
	"dongchamao/business"
	"dongchamao/business/es"
	"dongchamao/global"
	"dongchamao/services/dyimg"
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
