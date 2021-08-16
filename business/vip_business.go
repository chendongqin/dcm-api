package business

import (
	"dongchamao/global/logger"
	"dongchamao/models/dcm"
	"dongchamao/models/repost/dy"
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
func (receiver *VipBusiness) GetVipLevels(userId int) []dy.AccountVipLevel {
	vipLists := make([]dcm.DcUserVip, 0)
	err := dcm.GetSlaveDbSession().Where("user_id=?", userId).Find(&vipLists)
	vipList := make([]dy.AccountVipLevel, 0)
	if err == nil {
		for _, v := range vipLists {
			var level = 0
			if v.Expiration.After(time.Now()) {
				level = v.Level
			}
			isSub := 0
			if v.ParentId > 0 && level > 0 {
				isSub = 1
			}
			vipList = append(vipList, dy.AccountVipLevel{
				PlatForm:          v.Platform,
				Level:             level,
				SubNum:            v.SubNum,
				IsSub:             isSub,
				ExpirationTime:    v.Expiration,
				SubExpirationTime: v.SubExpiration,
			})
		}
	}
	return vipList
}

func (receiver *VipBusiness) GetVipLevel(userId, appId int) dy.AccountVipLevel {
	vip := dcm.DcUserVip{}
	var level = 0
	exist, _ := dcm.GetSlaveDbSession().Where("user_id=? AND platform=?", userId, appId).Get(&vip)
	if !exist {
		vip.UserId = userId
		vip.Platform = appId
		vip.UpdateTime = time.Now()
		vip.Expiration = time.Now().AddDate(0, 0, -1)
		vip.SubExpiration = time.Now().AddDate(0, 0, -1)
		_, _ = dcm.Insert(nil, &vip)
	}
	if vip.Expiration.Unix() > time.Now().Unix() {
		level = vip.Level
	}
	isSub := 0
	if vip.ParentId > 0 && level > 0 {
		isSub = 1
	}
	info := dy.AccountVipLevel{
		PlatForm:          vip.Platform,
		Level:             level,
		SubNum:            vip.SubNum,
		IsSub:             isSub,
		ExpirationTime:    vip.Expiration,
		SubExpirationTime: vip.SubExpiration,
	}
	return info
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
	if affect == 0 || logger.CheckError(err) != nil {
		return 0, false
	}
	return vipModel.OrderLevel, true
}
