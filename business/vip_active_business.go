package business

import (
	"dongchamao/models/dcm"
	"dongchamao/models/repost/dy"
	"time"
)

type VipActiveBusiness struct {
}

func NewVipActiveBusiness() *VipActiveBusiness {
	return new(VipActiveBusiness)
}

func (receiver *VipActiveBusiness) CheckDyVipActive(userId int, price dy.VipPriceConfig) dy.VipPriceConfig {
	receiver.BirthdayPriceActivity(userId, price)
	return price
}

//首月首次月销量处理
func (receiver *VipActiveBusiness) BirthdayPriceActivity(userId int, dyVipValue dy.VipPriceConfig) dy.VipPriceConfig {
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
	dyVipValue.Month.ActivePrice = 99
	dyVipValue.Month.ActiveComment += ""
	return dyVipValue
}
