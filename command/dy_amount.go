package command

import (
	"dongchamao/business"
	"dongchamao/models/dcm"
	"fmt"
	"time"
)

type DyAmount struct {
}

func getNoticeUser(days int) (list []dcm.UserVipJpinCombine) {
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	users := make([]dcm.UserVipJpinCombine, 0)
	beforeDay := fmt.Sprintf("%s%s", time.Now().AddDate(0, 0, days).Format("2006-01-02"), " 00:00:00")
	after30DaysTime := fmt.Sprintf("%s%s", time.Now().AddDate(0, 0, days).Format("2006-01-02"), " 23:59:59")
	whereString := "u.status = ? AND u.openid <> '' AND vip.Expiration BETWEEN ? AND ? AND vip.parent_id = ? AND vip.platform = ? "
	err := dbSession.Table("dc_user").Alias("u").
		Join("LEFT", []string{"dc_user_vip", "vip"}, "vip.user_id = u.id").
		Where(whereString, business.UserStatusNormal, beforeDay, after30DaysTime, 0, business.VipPlatformDouYin).
		Find(&users)
	if err != nil {
		//logs.Error("[notice account] 获取到期会员失败 err: %s", err)
		//fmt.Println(err)
		return
	}
	list = users
	return
}

//会员到期通知todo
func AmountExpireWechatNotice() {

	//30天内到期
	users := getNoticeUser(30)
	for _, v := range users {
		business.NewWechatBusiness().AmountExpireWechatNotice(&v)
	}
	users = getNoticeUser(7)
	for _, v := range users {
		business.NewWechatBusiness().AmountExpireWechatNotice(&v)
	}
	users = getNoticeUser(0)
	for _, v := range users {
		business.NewWechatBusiness().AmountExpireWechatNotice(&v)
	}

	return
}
