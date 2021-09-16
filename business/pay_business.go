package business

import (
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/models/repost"
	"dongchamao/models/repost/dy"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"math"
	"time"
)

var monthDay = 30
var halfYearDay = 180
var yearDay = 365

var DyVipPayMoney = map[int]float64{
	monthDay:    259,
	halfYearDay: 799,
	yearDay:     1199,
}

type PayBusiness struct {
}

func NewPayBusiness() *PayBusiness {
	return new(PayBusiness)
}

func (receiver *PayBusiness) DoPayDyCallback(vipOrder dcm.DcVipOrder) bool {
	dbSession := dcm.GetDbSession()
	_ = dbSession.Begin()
	userLevel := dcm.DcUserVip{}
	exist, _ := dbSession.Where("user_id=? AND platform=?", vipOrder.UserId, 1).Get(&userLevel)
	if !exist {
		userLevel.UserId = vipOrder.UserId
		userLevel.Platform = 1
		userLevel.UpdateTime = time.Now()
		userLevel.Expiration = time.Now().AddDate(0, 0, -1)
		userLevel.SubExpiration = time.Now().AddDate(0, 0, -1)
		affect, err := dcm.Insert(dbSession, &userLevel)
		if affect == 0 || err != nil {
			_ = dbSession.Rollback()
			return false
		}
	}
	if userLevel.Expiration.Before(time.Now()) {
		userLevel.Level = 0
	}
	orderInfo := repost.VipOrderInfo{}
	_ = jsoniter.Unmarshal([]byte(vipOrder.GoodsInfo), &orderInfo)
	updateMap := map[string]interface{}{}
	if userLevel.ParentId > 0 {
		userLevel.Level = 0
		updateMap["parent_id"] = 0
	}
	updateMap["level"] = vipOrder.Level
	updateMap["value_type"] = 2
	updateMap["update_time"] = time.Now().Format("2006-01-02 15:04:05")
	nowTime := time.Now()
	switch vipOrder.OrderType {
	case 1, 2:
		if userLevel.Level == 0 || userLevel.Level < vipOrder.Level {
			updateMap["expiration"] = nowTime.AddDate(0, 0, orderInfo.BuyDays).Format("2006-01-02 15:04:05")
		} else {
			updateMap["expiration"] = userLevel.Expiration.AddDate(0, 0, orderInfo.BuyDays).Format("2006-01-02 15:04:05")
		}
	case 3:
		updateMap["sub_expiration"] = userLevel.Expiration
		updateMap["sub_num"] = userLevel.SubNum + orderInfo.People
	case 4:
		updateMap["sub_expiration"] = userLevel.Expiration
	case 5:
		expiration := userLevel.Expiration.AddDate(0, 0, orderInfo.BuyDays).Format("2006-01-02 15:04:05")
		updateMap["expiration"] = expiration
		updateMap["sub_expiration"] = expiration
	}
	affect, err := dcm.UpdateInfo(dbSession, userLevel.Id, updateMap, new(dcm.DcUserVip))
	if affect == 0 || err != nil {
		_ = dbSession.Rollback()
		return false
	}
	_ = dbSession.Commit()
	NewUserBusiness().DeleteUserLevelCache(vipOrder.UserId, VipPlatformDouYin)
	return true
}

//剩余价值计算
func (receiver *PayBusiness) CountDySurplusValue(surplusDay int) float64 {
	yearPayMoney := DyVipPayMoney[yearDay]
	halfYearPayMoney := DyVipPayMoney[halfYearDay]
	monthPayMoney := DyVipPayMoney[monthDay]
	if surplusDay >= yearDay {
		value := float64(surplusDay) * yearPayMoney / float64(yearDay)
		return math.Ceil(value)
	}
	//半年剩余价值
	halfYear := surplusDay / halfYearDay
	halfYearValue := halfYearPayMoney * float64(halfYear)
	surplusDay -= halfYearDay * halfYear
	//剩余价值计算
	var dayValue float64 = 0
	if surplusDay > monthDay {
		dayValue = float64(surplusDay) * halfYearPayMoney / float64(halfYearDay)
	} else {
		dayValue = float64(surplusDay) * monthPayMoney / float64(monthDay)
	}
	value := halfYearValue + dayValue
	if value < 100 {
		value = 100
	}
	return math.Ceil(value)
}

func (receiver *PayBusiness) GetVipPriceConfig() (price dy.VipPriceConfig, initPrice dy.VipPriceConfig) {
	var configJson dcm.DcConfigJson
	_, _ = dcm.GetBy("key_name", "vip_price", &configJson)
	var config dy.VipPrice
	_ = json.Unmarshal([]byte(configJson.Value), &config)
	priceData := config.VipPrice
	var priceMap = make(map[int]float64, len(config.VipPrice))
	var initPriceMap = make(map[int]float64, len(config.VipPrice))
	for _, v := range priceData {
		priceMap[utils.ToInt(v.Days)] = utils.ToFloat64(v.Price)
		initPriceMap[utils.ToInt(v.Days)] = utils.ToFloat64(v.InitPrice)
	}
	//价格
	price = dy.VipPriceConfig{
		Year: priceMap[yearDay], HalfYear: priceMap[halfYearDay], Month: priceMap[monthDay],
	}
	//原价
	initPrice = dy.VipPriceConfig{
		Year: initPriceMap[yearDay], HalfYear: initPriceMap[halfYearDay], Month: initPriceMap[monthDay],
	}
	return
}

//扩充团队价格与原价
func (receiver *PayBusiness) GetDySurplusValue(surplusDay int) (value float64, initValue float64) {
	price, initPrice := receiver.GetVipPriceConfig()
	if surplusDay >= yearDay {
		value = float64(surplusDay) * price.Year / float64(yearDay)
		initValue = float64(surplusDay) * initPrice.Year / float64(yearDay)
		return math.Ceil(value), math.Ceil(initValue)
	}
	//半年剩余价值
	halfYear := surplusDay / halfYearDay
	halfYearValue := price.HalfYear * float64(halfYear)
	initHalfYearValue := initPrice.HalfYear * float64(halfYear)
	surplusDay -= halfYearDay * halfYear
	//剩余价值计算
	var dayValue float64 = 0
	var initDayValue float64 = 0
	if surplusDay > monthDay {
		dayValue = float64(surplusDay) * price.HalfYear / float64(halfYearDay)
		initDayValue = float64(surplusDay) * initPrice.Month / float64(halfYearDay)
	} else {
		dayValue = float64(surplusDay) * price.Month / float64(monthDay)
		initDayValue = float64(surplusDay) * initPrice.Month / float64(monthDay)
	}
	value = halfYearValue + dayValue
	initValue = initHalfYearValue + initDayValue
	if value < 100 {
		value = 100
	}
	if initValue < 100 {
		initValue = 100
	}
	return math.Ceil(value), math.Ceil(initValue)
}
