package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/models/repost/dy"
	"dongchamao/services/mutex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
)

//用户等级
const (
	//用户状态
	UserStatusNormal  = 1  //正常状态
	UserStatusDisable = 0  //禁用状态
	UserStatusCancel  = -1 //注销状态
)

type UserBusiness struct {
}

func NewUserBusiness() *UserBusiness {
	return new(UserBusiness)
}

type TokenData struct {
	AppId      int   `json:"appId"`
	Id         int   `json:"id"`
	ExpireTime int64 `json:"expire_time"`
}

//APPID:10000:pc端,10001:h5,10002:微信小程序,10003、10004：app,10005:Wap
//用户来源0:PC,1:小程序,2:APP,3:wap
var AppIdMap = map[int]int{10000: 0, 10001: 0, 10002: 1, 10003: 2, 10004: 2, 10005: 3}

func (receiver *UserBusiness) AddOrUpdateUniqueToken(userId int, appId int, token string) error {
	platFormId := GetAppPlatFormIdWithAppId(appId)
	t := utils.GetNowTimeStamp()

	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	sql := `
INSERT INTO dc_user_token (user_id,app_platform,token,create_time,update_time) VALUES (?,?,?,?,?) 
ON DUPLICATE KEY UPDATE  update_time = values(update_time), token = VALUES(token);
`
	_, err := dbSession.Exec(sql, userId, platFormId, token, t, t)
	//更新缓存里的token
	receiver.GetUniqueToken(userId, appId, false)
	return err
}

//获取唯一token
func (receiver *UserBusiness) GetUniqueToken(userId int, appId int, enableCache bool) (string, bool) {
	platFormId := GetAppPlatFormIdWithAppId(appId)
	cKey := cache.GetCacheKey(cache.UserPlatformUniqueToken, userId, platFormId)

	if enableCache == true {
		cToken := global.Cache.Get(cKey)
		if cToken != "" {
			return cToken, true
		}
	}
	userTokenM := &dcm.DcUserToken{}
	exist, _ := dcm.GetSlaveDbSession().Where("user_id = ?", userId).
		And("app_platform = ?", platFormId).
		Get(userTokenM)

	if exist {
		_ = global.Cache.Set(cKey, userTokenM.Token, 1800)
	}

	return userTokenM.Token, exist
}

//token创建
func (receiver *UserBusiness) CreateToken(appId int, userId int, expire int64) (string, int64, error) {
	if expire == 0 {
		expire = 3600
	}
	//如果是app登录，token有效期是30天
	if utils.InArrayInt(appId, []int{10003, 10004}) == true {
		expire += 23 * 86400
	}

	expire_time := utils.Time() + expire
	auth_code := []byte(global.Cfg.String("auth_code"))
	claims := jwt.MapClaims{
		"appId":       appId,
		"id":          userId,
		"expire_time": expire_time,
		"iat":         time.Now().Unix(), //签发时间,确保每次登录得到的token不同
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	authToken, err := token.SignedString(auth_code)

	//写入用户唯一token表
	//_ = receiver.AddOrUpdateUniqueToken(userId, appId, authToken)

	return authToken, expire_time, err
}

//密码登陆
func (receiver *UserBusiness) LoginByPwd(username, pwd string, appId int) (user dcm.DcUser, tokenString string, expire int64, comErr global.CommonError) {
	if pwd == "" {
		comErr = global.NewError(4208)
		return
	}
	exist, _ := dcm.GetBy("username", username, &user)
	if !exist {
		comErr = global.NewError(4208)
		return
	}
	tmpPwd := utils.Md5_encode(pwd + user.Salt)
	if tmpPwd != user.Password {
		comErr = global.NewError(4208)
		return
	}
	if user.Status != 1 {
		if user.Status == 0 {
			comErr = global.NewError(4212)
			return
		}
		comErr = global.NewError(4217)
		return
	}
	expire = 604800
	tokenString, expire, err := receiver.CreateToken(appId, user.Id, expire)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	err = receiver.AddOrUpdateUniqueToken(user.Id, appId, tokenString)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//密码登陆
func (receiver *UserBusiness) AppleLogin(appleId string, appId int) (user dcm.DcUser, tokenString string, expire int64, isNew int, comErr global.CommonError) {
	if appleId == "" {
		comErr = global.NewError(4208)
		return
	}
	exist, _ := dcm.GetBy("apple_id", appleId, &user)
	if !exist {
		user.Status = 1
		user.CreateTime = time.Now()
		user.UpdateTime = time.Now()
		//来源
		user.Entrance = AppIdMap[appId]
		user.AppleId = appleId
		affect, err := dcm.Insert(nil, &user)
		if affect == 0 || err != nil {
			comErr = global.NewError(5001)
			return
		}
		user.Username = fmt.Sprintf("8888%07d", user.Id)
		user.Nickname = "appleUser" + strconv.Itoa(user.Id)
		_, err = dcm.GetDbSession().Where("id=?", user.Id).Update(&user)
		if err != nil {
			comErr = global.NewError(5001)
			return
		}
		isNew = 1
	}
	if user.Status != 1 {
		if user.Status == 0 {
			comErr = global.NewError(4212)
			return
		}
		comErr = global.NewError(4217)
		return
	}
	expire = 604800
	tokenString, expire, err := receiver.CreateToken(appId, user.Id, expire)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	err = receiver.AddOrUpdateUniqueToken(user.Id, appId, tokenString)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//验证码登陆
func (receiver *UserBusiness) SmsLogin(mobile, code, password, unionid, appleId string, appId int) (user dcm.DcUser, tokenString string, expire int64, isNew int, comErr global.CommonError) {
	if mobile == "" || code == "" {
		comErr = global.NewError(4000)
		return
	}
	codeKey := cache.GetCacheKey(cache.SmsCodeVerify, "login", mobile)
	if verifyCode := global.Cache.Get(codeKey); verifyCode != code {
		comErr = global.NewError(4209)
		return
	}
	user = dcm.DcUser{}
	exist, err := dcm.GetBy("username", mobile, &user)
	if exist && err == nil {
		if user.Status != 1 {
			comErr = global.NewError(4212)
			return
		}
	} else {
		user.Username = mobile
		user.Nickname = mobile[:3] + "****" + mobile[7:]
		if password != "" {
			user.SetPassword = 1
		}
		user.Salt = utils.GetRandomString(4)
		user.Password = utils.Md5_encode(password + user.Salt)
		user.Status = 1
		user.CreateTime = time.Now()
		user.UpdateTime = time.Now()
		//来源
		user.Entrance = AppIdMap[appId]
		affect, err := dcm.Insert(nil, &user)
		if affect == 0 || err != nil {
			comErr = global.NewError(5001)
			return
		}
		isNew = 1
	}
	if user.AppleId != "" && appleId != "" {
		comErr = global.NewError(4406)
		return
	}
	user.BindPhone = 1
	user.AppleId = appleId
	if unionid != "" {
		if user.Unionid != "" {
			comErr = global.NewError(4305)
			return
		}
		wechatModel := dcm.DcWechat{} //如果有微信信息 头像/昵称 默认用微信
		if exist, _ := dcm.GetSlaveDbSession().Where("unionid = ?", unionid).Get(&wechatModel); !exist {
			comErr = global.NewError(4304)
			return
		}
		user.Nickname = wechatModel.NickName
		user.Avatar = wechatModel.Avatar
		user.Openid = wechatModel.Openid
		user.OpenidApp = wechatModel.OpenidApp
		user.Unionid = wechatModel.Unionid
	} else {
		user.Nickname = mobile[:3] + "****" + mobile[7:]
	}
	_, err = dcm.GetDbSession().Where("id=?", user.Id).Update(&user)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	tokenString, expire, err = receiver.CreateToken(appId, user.Id, 604800)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	err = receiver.AddOrUpdateUniqueToken(user.Id, appId, tokenString)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	_ = global.Cache.Delete(codeKey)
	return
}

//微信登录 1.客户端 2.扫码两种方式
func (receiver *UserBusiness) WechatLogin(unionid, source string, appId int) (user dcm.DcUser, tokenString string, expire int64, comErr global.CommonError) {
	if unionid == "" {
		comErr = global.NewError(4301)
		return
	}
	//判断 dc_wechat 是否有openid信息
	wechatModel := dcm.DcWechat{}
	if exist, _ := dcm.GetSlaveDbSession().Where("unionid = ?", unionid).Get(&wechatModel); !exist {
		comErr = global.NewError(4300)
		return
	}
	if source == "wechat" {
		if wechatModel.Subscribe != 1 {
			comErr = global.NewError(4300)
			return
		}
	}
	if source == "wechat_app" {
		if wechatModel.OpenidApp == "" {
			comErr = global.NewError(4302)
			return
		}
	}
	var err error
	//查询是否绑定用户
	userModel := dcm.DcUser{}
	if exist, _ := dcm.GetSlaveDbSession().Where("unionid = ?", unionid).Get(&userModel); exist {
		if userModel.Status != 1 {
			comErr = global.NewError(4212)
			return
		}
		tokenString, expire, err = receiver.CreateToken(appId, userModel.Id, 604800)
		if err != nil {
			comErr = global.NewError(5000)
			return
		}
		err = receiver.AddOrUpdateUniqueToken(userModel.Id, appId, tokenString)
		if err != nil {
			comErr = global.NewError(5000)
			return
		}
	}
	return userModel, tokenString, expire, nil
}

func (receiver *UserBusiness) DeleteUserInfoCache(userid int) bool {
	memberKey := cache.GetCacheKey(cache.UserInfo, userid)
	err := global.Cache.Delete(memberKey)
	if err == nil {
		return true
	} else {
		return false
	}
}

//清除等级缓存
func (receiver *UserBusiness) DeleteUserLevelCache(userid, levelType int) bool {
	memberKey := cache.GetCacheKey(cache.UserLevel, userid, levelType)
	err := global.Cache.Delete(memberKey)
	if err == nil {
		return true
	} else {
		return false
	}
}

//更新活跃时间等
func (receiver *UserBusiness) UpdateVisitedTimes(userAccount dcm.DcUser) bool {
	if userAccount.Id == 0 {
		return false
	}
	today := time.Now().Format("20060102")
	yesterdayTime := time.Now().AddDate(0, 0, -1)
	prevTimeDate := userAccount.PrevTime.Format("20060102")
	//如果上次登录时间是今天，则不处理
	if today == prevTimeDate {
		return false
	}
	//加锁防止多个请求进来都进行处理
	lockKey := cache.GetCacheKey(cache.UserPrevTimeLock, userAccount.Id)
	lockSecret := utils.GetRandomInt(4)
	lock, ok, err := mutex.TryLockWithTimeout(global.Cache.GetInstance().(redis.Conn), lockKey, lockSecret, 30)
	if err != nil {
		return false
	}
	if !ok {
		return false
	}
	defer lock.Unlock()
	dbSession := dcm.GetDbSession()
	defer dbSession.Close()
	updateData := map[string]interface{}{
		"prev_time":   utils.GetNowTimeStamp(),
		"successions": 1,
	}
	if yesterdayTime.Format("20060102") == prevTimeDate || userAccount.Successions == 0 {
		successions := userAccount.Successions + 1
		totalSuccessions := userAccount.TotalSuccessions + 1
		updateData["successions"] = successions
		updateData["total_successions"] = totalSuccessions
		if successions > userAccount.MaxSuccessions {
			updateData["max_successions"] = successions
		}
	}
	affect, err := dcm.UpdateInfo(nil, userAccount.Id, updateData, new(dcm.DcUser))
	if affect == 0 || err != nil {
		return false
	}
	return true
}

//更新用户数据
func (receiver *UserBusiness) MobileExist(mobile string) (bool, global.CommonError) {
	//新手机重复校验
	var exist dcm.DcUser
	dbSession := dcm.GetDbSession()
	if _, err := dbSession.Where("username=?", mobile).Get(&exist); err != nil {
		return false, global.NewError(4000)
	}
	return exist.Id != 0, nil
}

//更新用户数据
func (receiver *UserBusiness) UpdateUserAndClearCache(dbSession *xorm.Session, userId int, updateData map[string]interface{}) (int64, error) {
	affect, err := dcm.UpdateInfo(dbSession, userId, updateData, new(dcm.DcUser))
	if affect != 0 && err == nil {
		receiver.DeleteUserInfoCache(userId)
		receiver.DeleteUserLevelCache(userId, 1)
	}
	return affect, err
}

func (receiver *UserBusiness) GetCacheUser(userId int, enableCache bool) (dcm.DcUser, bool) {
	memberKey := cache.GetCacheKey(cache.UserInfo, userId)
	user := dcm.DcUser{}
	if enableCache == true {
		userJson := global.Cache.Get(memberKey)
		if userJson != "" {
			_ = jsoniter.Unmarshal([]byte(userJson), &user)
			return user, true
		}
	}
	exist, _ := dcm.GetSlaveDbSession().Where("id = ?", userId).
		Get(&user)
	if exist {
		userByte, _ := jsoniter.Marshal(user)
		_ = global.Cache.Set(memberKey, string(userByte), 1800)
	}
	return user, exist
}

func (receiver *UserBusiness) GetCacheUserLevel(userId, levelType int, enableCache bool) int {
	memberKey := cache.GetCacheKey(cache.UserLevel, userId, levelType)
	if enableCache == true {
		level := global.Cache.Get(memberKey)
		if level != "" {
			return utils.ToInt(level)
		}
	}
	vipBusiness := NewVipBusiness()

	vipLevel := vipBusiness.GetVipLevel(userId, levelType)
	expireTime := 1800 * time.Second
	if time.Now().Unix()-vipLevel.ExpirationTime.Unix() < 1800 {
		expireTime = time.Now().Sub(vipLevel.ExpirationTime)
	}
	_ = global.Cache.Set(memberKey, utils.ToString(vipLevel.Level), expireTime)
	return vipLevel.Level
}

//关键词统计
func (receiver *UserBusiness) KeywordsRecord(keyword string) (comErr global.CommonError) {
	if keyword == "" {
		return
	}
	var record dcm.DcUserKeywordsRecord
	db := dcm.GetDbSession()
	defer db.Close()
	exist, err := db.Where("keyword=?", keyword).Get(&record)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	if exist {
		if _, err := db.Table(new(dcm.DcUserKeywordsRecord)).Where("id=?", record.Id).Incr("count", 1).Update(map[string]interface{}{
			"update_time": time.Now().Format("2006-01-02 15:04:05"),
		}); err != nil {
			comErr = global.NewError(5000)
		}
		return
	}
	record.Keyword = keyword
	record.CreateTime = time.Now()
	record.UpdateTime = time.Now()
	record.Count = 1
	_, err = dcm.Insert(db, &record)
	if err != nil {
		comErr = global.NewError(5000)
	}
	return
}

//关键词统计
func (receiver *UserBusiness) GetUserList(userIds []string) (userList []dcm.DcUser, comErr global.CommonError) {
	if err := dcm.GetDbSession().In("id", userIds).Find(&userList); err != nil {
		return nil, global.NewCommonError(err)
	}
	return
}

//新用户注册赠送vip
func (receiver *UserBusiness) SendUserVip(user *dcm.DcUser, buyDays int) {
	uniqueID, _ := utils.Snow.GetSnowflakeId()
	now := time.Now()
	title := fmt.Sprintf("赠送%d天专业版", buyDays)
	remark := "新用户注册"
	ExpirationDate := now.AddDate(0, 0, buyDays)
	var VipOrder = dcm.DcVipOrder{
		UserId:       user.Id,
		Username:     user.Username,
		TradeNo:      fmt.Sprintf("%s%d", time.Now().Format("060102"), uniqueID),
		Channel:      0,
		InterTradeNo: "",
		OrderType:    6,
		Platform:     "douyin",
		Level:        UserLevelJewel,
		BuyDays:      buyDays,
		Title:        title,
		Status:       1,
		PayStatus:    1,
		GoodsInfo:    "",
		Remark:       remark,
		CreateTime:   now,
		UpdateTime:   now,
	}

	var UserVip = dcm.DcUserVip{
		UserId:     user.Id,
		Platform:   VipPlatformDouYin,
		Level:      UserLevelJewel,
		Expiration: ExpirationDate,
		Remark:     "新用户注册赠送vip",
		UpdateTime: now,
	}
	dbSession := dcm.GetDbSession()
	affectOrder, errOrder := dcm.Insert(dbSession, &VipOrder)
	if affectOrder == 0 || errOrder != nil {
		//logger.CheckError(errOrder)
		NewMonitorBusiness().SendErr("注册赠送vip", fmt.Sprintf("vip_order数据表插入报错:%s", errOrder))
		_ = dbSession.Rollback()
		return
	} else {
		affectUser, errUser := dcm.Insert(dbSession, &UserVip)
		if affectUser == 0 || errUser != nil {
			NewMonitorBusiness().SendErr("注册赠送vip", fmt.Sprintf("user_vip数据表插入报错:%s", errOrder))
			_ = dbSession.Rollback()
			return
		} else {
			_ = dbSession.Commit()
		}
	}
}

func (receiver *UserBusiness) CollectSum(userId int) (sum dy.CollectSum) {
	var collectList []dcm.DcUserDyCollect
	if err := dcm.GetDbSession().Table("dc_user_dy_collect").Where("user_id=? and status=1", userId).Find(&collectList); err != nil {
		return
	}
	for _, v := range collectList {
		//1达人2商品3视频
		switch v.CollectType {
		case 1:
			sum.Author++
		case 2:
			sum.Product++
		case 3:
			sum.Aweme++
		}
	}
	return
}
