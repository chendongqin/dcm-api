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

func (receiver *PayBusiness) GetVipPriceConfig() (price dy.VipPriceConfig, primePrice dy.VipPriceConfig) {
	priceMap, primePriceMap := receiver.GetVipPriceConfigMap()
	price = dy.VipPriceConfig{
		Year: priceMap[yearDay], HalfYear: priceMap[halfYearDay], Month: priceMap[monthDay],
	}
	primePrice = dy.VipPriceConfig{
		Year: primePriceMap[yearDay], HalfYear: primePriceMap[halfYearDay], Month: primePriceMap[monthDay],
	}
	return
}

//获取抖音会员价格Map数据
func (receiver *PayBusiness) GetVipPriceConfigMap() (priceMap map[int]float64, primePriceMap map[int]float64) {
	var configJson dcm.DcConfigJson
	_, _ = dcm.GetBy("key_name", "vip_price", &configJson)
	var config dy.VipPrice
	_ = json.Unmarshal([]byte(configJson.Value), &config)
	priceData := config.VipPrice
	priceMap = make(map[int]float64, len(config.VipPrice))
	primePriceMap = make(map[int]float64, len(config.VipPrice))
	for _, v := range priceData {
		priceMap[utils.ToInt(v.Days)] = utils.ToFloat64(v.Price)
		primePriceMap[utils.ToInt(v.Days)] = utils.ToFloat64(v.InitPrice)
	}
	return
}

//扩充团队价格与原价
func (receiver *PayBusiness) GetDySurplusValue(surplusDay int) (value float64, primeValue float64) {
	price, primePrice := receiver.GetVipPriceConfig()
	if surplusDay >= yearDay {
		value = float64(surplusDay) * price.Year / float64(yearDay)
		primeValue = float64(surplusDay) * primePrice.Year / float64(yearDay)
		return math.Ceil(value), math.Ceil(primeValue)
	}
	//半年剩余价值
	halfYear := surplusDay / halfYearDay
	halfYearValue := price.HalfYear * float64(halfYear)
	primeHalfYearValue := primePrice.HalfYear * float64(halfYear)
	surplusDay -= halfYearDay * halfYear
	//剩余价值计算
	var dayValue float64 = 0
	var primeDayValue float64 = 0
	if surplusDay > monthDay {
		dayValue = float64(surplusDay) * price.HalfYear / float64(halfYearDay)
		primeDayValue = float64(surplusDay) * primePrice.Month / float64(halfYearDay)
	} else {
		dayValue = float64(surplusDay) * price.Month / float64(monthDay)
		primeDayValue = float64(surplusDay) * primePrice.Month / float64(monthDay)
	}
	value = halfYearValue + dayValue
	primeValue = primeHalfYearValue + primeDayValue
	if value < 100 {
		value = 100
	}
	if primeValue < 100 {
		primeValue = 100
	}
	return math.Ceil(value), math.Ceil(primeValue)
}

//首月首次月销量处理
func (receiver *PayBusiness) BirthdayPriceActivity(userId int, dyVipValue map[int]float64) map[int]float64 {
	if time.Now().Unix() >= 1635696000 {
		return dyVipValue
	}
	if userId > 0 {
		exist, _ := dcm.GetSlaveDbSession().
			Where("user_id=? AND status>=0 AND expiration_time > ?", userId, time.Now().Format("2006-01-02 15:04:05")).
			Get(new(dcm.DcVipOrder))
		if exist {
			return dyVipValue
		}
	}
	dyVipValue[monthDay] = 99
	return dyVipValue
}
