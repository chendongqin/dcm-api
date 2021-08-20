package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	"dongchamao/services/mutex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"time"
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
		comErr = global.NewError(4212)
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
func (receiver *UserBusiness) SmsLogin(mobile, code string, appId int) (user dcm.DcUser, tokenString string, expire int64, isNew int, comErr global.CommonError) {
	codeKey := cache.GetCacheKey(cache.SmsCodeVerify, "login", mobile)
	verifyCode := global.Cache.Get(codeKey)
	if verifyCode != code {
		comErr = global.NewError(4209)
		return
	}
	exist, err := dcm.GetBy("username", mobile, &user)
	if exist && err == nil {
		if user.Status != 1 {
			comErr = global.NewError(4212)
			return
		}
	} else {
		user.Username = mobile
		user.Nickname = mobile[:3] + "****" + mobile[7:]
		user.Salt = utils.GetRandomString(4)
		user.Password = utils.Md5_encode(utils.GetRandomString(16) + user.Salt)
		user.Status = 1
		user.CreateTime = time.Now()
		user.UpdateTime = time.Now()
		//来源
		user.Entrance = AppIdMap[appId]
		affect, err := dcm.Insert(nil, &user)
		if affect == 0 || err != nil {
			comErr = global.NewError(5000)
			return
		}
		isNew = 1
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
		"prev_time": utils.GetNowTimeStamp(),
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
func (receiver *UserBusiness) UpdateUserAndClearCache(dbSession *xorm.Session, userId int, updateData map[string]interface{}) (int64, error) {
	affect, err := dcm.UpdateInfo(dbSession, userId, updateData, new(dcm.DcUser))
	if affect != 0 && err == nil {
		receiver.DeleteUserInfoCache(userId)
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
	_ = global.Cache.Set(memberKey, utils.ToString(vipLevel.Level), 1800)
	return vipLevel.Level
}

type CollectRet struct {
	dcm.DcUserDyCollect
	FollowerCount      int64
	FollowerIncreCount int64
	Predict7Gmv        float64
	Predict7Digg       float64
}

func (receiver *UserBusiness) GetDyCollect(tagId, collectType int, keywords string) (data []CollectRet, comErr global.CommonError) {
	var collects []dcm.DcUserDyCollect
	dbCollect := dcm.GetDbSession().Table(dcm.DcUserDyCollect{})
	defer dbCollect.Close()
	var query string
	query = fmt.Sprintf("tag_id=%v AND collect_type=%v", tagId, collectType)
	if keywords != "" {
		query += " AND (unique_id LIKE '%" + keywords + "%' or nickname LIKE '%" + keywords + "%')"
	}
	err := dbCollect.Where(query).Find(&collects)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	data = make([]CollectRet, len(collects))
	for k, v := range collects {
		data[k].DcUserDyCollect = v
		dyAuthor, _ := hbase.GetAuthor(v.CollectId)
		basic, _ := hbase.GetAuthorBasic(v.CollectId, "")
		data[k].FollowerCount = dyAuthor.Data.Fans.Douyin.Count
		data[k].FollowerIncreCount = basic.FollowerCount - basic.FollowerCountBefore
	}
	return
}

//收藏达人
func (receiver *UserBusiness) AddDyCollect(collectId string, collectType, userId int) (comErr global.CommonError) {
	collect := dcm.DcUserDyCollect{}
	dbCollect := dcm.GetDbSession().Table(collect)
	defer dbCollect.Close()
	exist, err := dbCollect.Where("user_id=? AND collect_type=? AND collect_id=?", userId, collectType, collectId).Get(&collect)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	if collect.Status == 1 {
		comErr = global.NewMsgError("您已收藏该达人，请刷新重试")
		return comErr
	}
	collect.Status = 1
	collect.CollectId = collectId
	collect.UpdateTime = time.Now()
	switch collectType {
	case 1:
		author, comErr := hbase.GetAuthor(collectId)
		if comErr != nil {
			return comErr
		}
		collect.Label = author.Tags
		collect.UniqueId = author.Data.UniqueID
		collect.Nickname = author.Data.Nickname
	}
	if exist {
		if _, err := dbCollect.Update(collect); err != nil {
			comErr = global.NewError(5000)
			return
		}
	} else {
		collect.CreateTime = time.Now()
		collect.UserId = userId
		collect.CollectType = collectType
		if _, err := dbCollect.Insert(collect); err != nil {
			comErr = global.NewError(5000)
			return
		}
	}
	return
}

//取消收藏
func (receiver *UserBusiness) CancelDyCollect(id int) (comErr global.CommonError) {
	dbCollect := dcm.GetDbSession().Table(dcm.DcUserDyCollect{})
	defer dbCollect.Close()
	exist, err := dbCollect.Where("id=? and status=?", id, 1).Exist()
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	if !exist {
		comErr = global.NewMsgError("您未收藏该达人，请刷新重试")
		return
	}
	affect, err := dcm.UpdateInfo(dbCollect, id, map[string]interface{}{"status": 0, "update_time": time.Now()}, new(dcm.DcUserDyCollect))
	if err != nil || affect == 0 {
		comErr = global.NewError(5000)
		return
	}
	return
}

func (receiver *UserBusiness) KeywordsRecord(keyword string) (comErr global.CommonError) {
	var record dcm.DcUserKeywordsRecord
	db := dcm.GetDbSession().Table(dcm.DcUserKeywordsRecord{})
	exist, err := db.Where("keyword=?", keyword).Get(&record)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	record.Count++
	if exist {
		if _, err := db.Update(record); err != nil {
			comErr = global.NewError(5000)
		}
	} else if _, err := db.Insert(record); err != nil {
		comErr = global.NewError(5000)
	}
	return
}
