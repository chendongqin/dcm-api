package command

type DyAmount struct {
}

//会员到期通知todo
func AmountExpireWechatNotice() {
	//cacheKey30  := cache.GetCacheKey(cache.AmountExpireWechatNotice, 30)
	//cacheKey7   := cache.GetCacheKey(cache.AmountExpireWechatNotice, 7)
	//cacheKey1   := cache.GetCacheKey(cache.AmountExpireWechatNotice, 1)
	//cacheData30 := global.Cache.Get(cacheKey30)
	//cacheData7  := global.Cache.Get(cacheKey7)
	//cacheData1  := global.Cache.Get(cacheKey1)
	//
	//utils.DeserializeData(cacheData30)
	//
	//
	////_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 300)
	//
	//dbSession := dcm.GetDbSession()
	//defer dbSession.Close()
	//
	//type DcUser struct {
	//	Id int64
	//	Expiration     time.Time `xorm:"comment('过期时间') TIMESTAMP"`
	//}
	//type DcUserVip struct {
	//	Openid string
	//	Platform       int       `xorm:"not null default 1 comment('1抖音2小红书3淘宝') unique(USER_LEVEL) TINYINT(1)"`
	//	ParentId       int       `xorm:"not null default 0 comment('主账户id') INT(11)"`
	//}
	//
	//type UserVipType struct {
	//	DcUser `xorm:"extends"`
	//	DcUserVip `xorm:"extends"`
	//
	//}
	//users := make([]UserVipType, 0)
	//today := fmt.Sprintf("%s%s", time.Now().Format("2006-01-02"), " 00:00:00")
	//after30DaysTime := fmt.Sprintf("%s%s", time.Now().AddDate(0,0,30).Format("2006-01-02"), " 00:00:00")
	//whereString := "u.status = ? AND vip.Expiration BETWEEN ? AND ? AND vip.parent_id = ? AND vip.platform = ? "
	//err := dbSession.Table("dc_user").Alias("u").
	//	Join("LEFT", []string{"dc_user_vip", "vip"}, "vip.user_id = u.id").
	//	Where(whereString,business.UserStatusNormal,today,after30DaysTime,0,business.VipPlatformDouYin).Find(&users)
	//if err != nil {
	//	//logs.Error("[notice account] 获取到期会员失败 err: %s", err)
	//	return
	//}
	//
	//for k,v := range users{
	//	fmt.Println(k)
	//	fmt.Println(v)
	//}
	////fmt.Printf("%+v \n",users)
	////fmt.Println(users)
	////fmt.Println(users[0].Openid)
	////fmt.Printf("%+v \n",users[0].DcUser.Id)
	//returnRes
}
