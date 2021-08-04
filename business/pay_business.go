package business

import "math"

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
