package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/models/dcm"
	"dongchamao/structinit/repost/dy"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type ProductBusiness struct {
}

func NewProductBusiness() *ProductBusiness {
	return new(ProductBusiness)
}

func (receiver *ProductBusiness) GetCacheProductCate(enableCache bool) []dy.DyCate {
	memberKey := cache.GetCacheKey(cache.ConfigKeyCache, "product_cate")
	pCate := make([]dy.DyCate, 0)
	if enableCache == true {
		jsonStr := global.Cache.Get(memberKey)
		if jsonStr != "" {
			_ = jsoniter.Unmarshal([]byte(jsonStr), &pCate)
			return pCate
		}
	}
	firstList := make([]dcm.DcProductCate, 0)
	_ = dcm.GetSlaveDbSession().Where("level=1").Desc("weight").Find(&firstList)
	for _, v := range firstList {
		secondList := make([]dcm.DcProductCate, 0)
		item := dy.DyCate{
			Name:    v.Name,
			SonCate: []dy.DyCate{},
		}
		_ = dcm.GetSlaveDbSession().Where("level=2 AND parent_id = ?", v.Id).Desc("weight").Find(&secondList)
		for _, vs := range secondList {
			item.SonCate = append(item.SonCate, dy.DyCate{
				Name:    vs.Name,
				SonCate: []dy.DyCate{},
			})
		}
		pCate = append(pCate, item)
	}
	if len(pCate) > 0 {
		userByte, _ := jsoniter.Marshal(pCate)
		_ = global.Cache.Set(memberKey, string(userByte), 86400)
	}
	return pCate
}

//获取商品url
func (receiver *ProductBusiness) GetProductUrl(platform, productId string) string {
	url := ""
	switch platform {
	case "淘宝":
		url = "https://item.taobao.com/item.htm?id=%s"
	case "京东":
		url = " https://item.m.jd.com/product/%s.html"
	case "天猫":
		url = " https://detail.tmall.com/item.htm?id=%s"
	case "苏宁":
		url = "https://m.suning.com/product/0000000000/0000000%s.html"
	case "小店":
		url = "https://haohuo.jinritemai.com/views/product/item2?id=%s"
	case "唯品会":
		url = "https://m.vip.com/public/go.html?pid=%s"
	case "考拉":
		url = "https://m-goods.kaola.com/product/%s.html"
	}
	if url != "" {
		url = fmt.Sprintf(url, productId)
	}
	return url
}
