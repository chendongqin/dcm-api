package business

import (
	"dongchamao/models/dcm"
	"time"
)

//用户等级
const (
	UserLevelDefault = 0 //普通会员
	UserLevelVip     = 1 //vip
	UserLevelSvip    = 2 //svip
	UserLevelJewel   = 3 //专业版
)

//会员平台
const (
	VipPlatformDouYin      = 1 //抖音
	VipPlatformXiaoHongShu = 2 //小红书
	VipPlatformTaoBao      = 3 //淘宝
)

type VipBusiness struct {
}

func NewVipBusiness() *VipBusiness {
	return new(VipBusiness)
}

//获取会员等级
func (this *VipBusiness) GetUserLevels() map[int]string {
	levels := map[int]string{}
	levels[UserLevelDefault] = "普通会员"
	levels[UserLevelVip] = "vip"
	levels[UserLevelSvip] = "svip"
	levels[UserLevelJewel] = "专业版"
	return levels
}

//获取等级名称
func (receiver *VipBusiness) GetUserLevel(level int) string {
	userLevels := receiver.GetUserLevels()
	if v, ok := userLevels[level]; ok {
		return v
	}
	return ""
}

//获取用户vip等级
func (receiver *VipBusiness) GetVipLevels(userId int) map[int]int {
	vipLists := make([]dcm.DcUserVip, 0)
	err := dcm.GetSlaveDbSession().Where("user_id=? ", userId).Find(&vipLists)
	vipMap := map[int]int{}
	if err == nil {
		for _, v := range vipLists {
			var level = 0
			if v.Expiration.Unix() > time.Now().Unix() {
				level = v.Level
			} else if v.OrderValidDay > 0 {
				levelTmp, res := receiver.UpdateValidDayOne(userId, v.Platform)
				if res {
					level = levelTmp
				}
			}
			vipMap[v.Platform] = level
		}
	}
	return vipMap
}

func (receiver *VipBusiness) GetVipLevel(userId, appId int) int {
	vip := &dcm.DcUserVip{}
	var level = 0
	exist, err := dcm.GetSlaveDbSession().Where("user_id=? AND platform=?", userId, appId).Get(vip)
	if err != nil {
		return 0
	}
	if !exist {
		vip.UserId = userId
		vip.Platform = appId
		vip.UpdateTime = time.Now()
		vip.Expiration = time.Now().AddDate(0, 0, -1)
		dcm.Insert(nil, vip)
	}
	if vip.Expiration.Unix() > time.Now().Unix() {
		level = vip.Level
	} else if vip.OrderValidDay > 0 {
		levelTmp, res := receiver.UpdateValidDayOne(userId, vip.Platform)
		if res {
			level = levelTmp
		}
	}
	return level
}

//更新会员等级
func (receiver *VipBusiness) UpdateValidDayOne(userId, platformId int) (int, bool) {
	vipModel := &dcm.DcUserVip{}
	dbSession := dcm.GetDbSession()
	exist, _ := dbSession.Where("user_id=? AND platform=?", userId, platformId).Get(vipModel)
	if !exist || vipModel.OrderValidDay <= 0 {
		return 0, false
	}
	whereStr := "id=? AND expiration<=?"
	updateData := map[string]interface{}{
		"order_valid_day": 0,
		"order_level":     0,
		"level":           vipModel.OrderLevel,
		"expiration":      vipModel.Expiration.AddDate(0, 0, vipModel.OrderValidDay).Format("2006-01-02 15:04:05"),
	}
	affect, err := dbSession.Table(new(dcm.DcUserVip)).Where(whereStr, vipModel.Id, time.Now().Format("2006-01-02 15:04:05")).Update(updateData)
	if affect == 0 || err != nil {
		return 0, false
	}
	return vipModel.OrderLevel, true
}
