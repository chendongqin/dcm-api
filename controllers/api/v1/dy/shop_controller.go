package dy

import (
	"dongchamao/business"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/hbase"
	"dongchamao/models/repost/dy"
	dy2 "dongchamao/models/repost/dy"
	"sort"
	"dongchamao/models/entity"
	"time"
)

type ShopController struct {
	controllers.ApiBaseController
}

func (receiver *ShopController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
	receiver.CheckDyUserGroupRight(business.DyJewelBaseMinShowNum, business.DyJewelBaseShowNum)
}

//小店基本数据
func (receiver *ShopController) ShopBase() {
	shopId := business.IdDecrypt(receiver.Ctx.Input.Param(":shop_id"))
	receiver.SuccReturn(shopId)
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

	receiver.SuccReturn(returnRes)
	return
}
