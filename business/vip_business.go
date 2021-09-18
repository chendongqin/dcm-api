package business

import (
	"dongchamao/global"
	"dongchamao/global/logger"
	"dongchamao/models/dcm"
	"dongchamao/models/repost/dy"
	"log"
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
	var parentVip dcm.DcUserVip
	err := dcm.GetSlaveDbSession().Where("user_id=?", userId).Find(&vipLists)
	vipList := make([]dy.AccountVipLevel, 0)
	if err == nil {
		for _, v := range vipLists {
			if v.ParentId != 0 {
				_, _ = dcm.GetSlaveDbSession().Where("id=?", v.ParentId).Get(&parentVip)
				parentVip.Expiration = parentVip.SubExpiration
				v = parentVip
			}
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
				ParentId:          v.ParentId,
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
	parentId := vip.ParentId
	expiration := vip.Expiration
	if parentId != 0 {
		parentVip := dcm.DcUserVip{}
		exist, _ = dcm.Get(vip.ParentId, &parentVip)
		if parentVip.Expiration.After(time.Now()) {
			expiration = parentVip.Expiration
			vip = parentVip
		} else {
			//父账号团队过期
			go func() {
				if _, err := dcm.UpdateInfo(nil, vip.ParentId, map[string]interface{}{"sub_num": 0}, new(dcm.DcUserVip)); err != nil {
					log.Println("parent_vip_expired:", err.Error())
					return
				}
				if _, err := dcm.UpdateInfo(nil, vip.Id, map[string]interface{}{"parent_id": 0, "remark": ""}, new(dcm.DcUserVip)); err != nil {
					log.Println("parent_vip_expired_sub:", err.Error())
					return
				}
			}()
		}
	}
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
		Id:                vip.Id,
		PlatForm:          vip.Platform,
		ParentId:          parentId,
		Level:             level,
		SubNum:            vip.SubNum,
		IsSub:             isSub,
		FeeLiveMonitor:    vip.LiveMonitorNum,
		ExpirationTime:    expiration,
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

//添加抖音团队成员
func (this *VipBusiness) AddDyTeamSub(userId, subUserId int) global.CommonError {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	var subUserVip dcm.DcUserVip
	var now = time.Now()
	if _, err := dbSession.Where("user_id=? and platform=1", subUserId).Get(&subUserVip); err != nil {
		return global.NewError(5000)
	}
	if subUserVip.Level != 0 && subUserVip.Expiration.After(now) {
		return global.NewMsgError("专业版账号无法添加")
	}
	if subUserVip.ParentId != 0 {
		return global.NewMsgError("已在团队中")
	}
	var userVip dcm.DcUserVip
	if _, err := dbSession.Where("user_id=? and platform=1", userId).Get(&userVip); err != nil {
		return global.NewError(5000)
	}
	subCount, err := dbSession.Table("dc_user_vip").Where("parent_id=?", userId).Count()
	if err != nil {
		return global.NewError(5000)
	}
	if int(subCount) >= userVip.SubNum {
		return global.NewMsgError("人数已满")
	}
	if _, err := dcm.UpdateInfo(dbSession, subUserVip.Id, map[string]interface{}{"parent_id": userVip.Id, "update_time": now.Format("2006-01-02 15:04:05")}, new(dcm.DcUserVip)); err != nil {
		return global.NewError(5000)
	}
	return nil
}

//添加抖音团队成员
func (this *VipBusiness) RemoveDyTeamSub(subUserId int) global.CommonError {
	if _, err := dcm.UpdateInfo(dcm.GetDbSession(), subUserId, map[string]interface{}{"parent_id": 0, "remark": ""}, new(dcm.DcUserVip)); err != nil {
		return global.NewError(5000)
	}
	return nil
}

func (this *VipBusiness) GetDyTeam(parentId, page, pageSize int) (list []dcm.DcUserVip, total int64, comErr global.CommonError) {
	err := dcm.GetDbSession().Table(dcm.DcUserVip{}).Where("parent_id=?", parentId).Limit(pageSize, (page-1)*pageSize).Find(&list)
	if err != nil {
		return nil, 0, global.NewError(5000)
	}
	total, err = dcm.GetDbSession().Table(dcm.DcUserVip{}).Where("parent_id=?", parentId).Count()
	if err != nil {
		return nil, 0, global.NewError(5000)
	}
	return
}
