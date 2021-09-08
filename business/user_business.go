package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	"dongchamao/models/repost"
	"dongchamao/services/dyimg"
	"dongchamao/services/mutex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"strconv"
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
func (receiver *UserBusiness) SmsLogin(mobile, code, unionid string, appId int) (user dcm.DcUser, tokenString string, expire int64, isNew int, comErr global.CommonError) {
	if mobile == "" || code == "" {
		comErr = global.NewError(4000)
		return
	}
	codeKey := cache.GetCacheKey(cache.SmsCodeVerify, "login", mobile)
	if verifyCode := global.Cache.Get(codeKey); verifyCode != code {
		comErr = global.NewError(4209)
		return
	}
	if unionid != "" {
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

//获取抖音收藏
func (receiver *UserBusiness) GetDyCollect(tagId, collectType int, keywords, label string, userId, page, pageSize int) (interface{}, int64, global.CommonError) {
	var (
		total    int64
		comErr   global.CommonError
		collects []dcm.DcUserDyCollect
	)
	dbCollect := dcm.GetDbSession()
	defer dbCollect.Close()
	var query string
	query = fmt.Sprintf("collect_type=%v", collectType)
	if tagId != 0 {
		query = fmt.Sprintf("tag_id=%v", tagId)
	}
	if keywords != "" {
		query += " AND (unique_id LIKE '%" + keywords + "%' or nickname LIKE '%" + keywords + "%')"
	}
	if label != "" {
		query += " AND FIND_IN_SET('" + label + "',label)"
	}
	query += " AND user_id=" + strconv.Itoa(userId) + " AND status=1"
	err := dbCollect.Table(dcm.DcUserDyCollect{}).Where(query).Limit(pageSize, (page-1)*pageSize).Find(&collects)
	if total, err = dbCollect.Table(dcm.DcUserDyCollect{}).Where(query).Count(); err != nil {
		comErr = global.NewError(5000)
		return nil, total, comErr
	}
	switch collectType {
	case 1:
		data := make([]repost.CollectAuthorRet, len(collects))
		for k, v := range collects {
			data[k].DcUserDyCollect = v
			data[k].DcUserDyCollect.CollectId = IdEncrypt(v.CollectId)
			dyAuthor, _ := hbase.GetAuthor(v.CollectId)
			basicData, _ := hbase.GetAuthorBasic(v.CollectId, time.Now().AddDate(0, 0, -1).Format("20060102"))
			data[k].FollowerCount = dyAuthor.Data.Fans.Douyin.Count
			data[k].FollowerIncreCount = dyAuthor.FollowerCount - basicData.FollowerCount
			data[k].Avatar = dyimg.Avatar(dyAuthor.Data.Avatar)
		}
		return data, total, nil
	case 2:
		data := make([]repost.CollectProductRet, len(collects))
		for k, v := range collects {
			data[k].DcUserDyCollect = v
			data[k].DcUserDyCollect.CollectId = IdEncrypt(v.CollectId)
		}
		return data, total, nil
	case 3:
		data := make([]repost.CollectAwemeRet, len(collects))
		for k, v := range collects {
			data[k].DcUserDyCollect = v
			data[k].DcUserDyCollect.CollectId = IdEncrypt(v.CollectId)
		}
		return data, total, nil
	}
	return nil, 0, nil
}

//获取分组收藏数量
func (receiver *UserBusiness) GetDyCollectCount(userId int) (data []repost.CollectCount, comErr global.CommonError) {
	dbCollect := dcm.GetDbSession()
	defer dbCollect.Close()
	if err := dbCollect.Table(dcm.DcUserDyCollect{}).Where("user_id=? AND status=1", userId).Select("tag_id,count(collect_id) as count").GroupBy("tag_id").Find(&data); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//获取已收藏达人标签
func (receiver *UserBusiness) GetDyCollectLabel(userId int) (data []string, comErr global.CommonError) {
	dbCollect := dcm.GetDbSession()
	defer dbCollect.Close()
	if err := dbCollect.Table(dcm.DcUserDyCollect{}).Where("user_id=? AND label<>'' AND status=1", userId).Select("label").Find(&data); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//收藏达人
func (receiver *UserBusiness) AddDyCollect(collectId string, collectType, tagId, userId int) (comErr global.CommonError) {
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
	collect.TagId = tagId
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
		if _, err := dbCollect.ID(collect.Id).Update(&collect); err != nil {
			comErr = global.NewError(5000)
			return
		}
	} else {
		collect.CreateTime = time.Now()
		collect.UserId = userId
		collect.CollectType = collectType
		if _, err := dbCollect.Insert(&collect); err != nil {
			comErr = global.NewError(5000)
			return
		}
	}
	return
}

func (receiver *UserBusiness) DyCollectExist(collectId string, collectType, userId int) (exist int) {
	collect := dcm.DcUserDyCollect{}
	dbCollect := dcm.GetDbSession().Table(collect)
	defer dbCollect.Close()
	_, _ = dbCollect.Where("user_id=? AND collect_type=? AND collect_id=?", userId, collectType, collectId).Get(&collect)
	return collect.Status
}

//取消收藏
func (receiver *UserBusiness) CancelDyCollect(id, userId int) (comErr global.CommonError) {
	dbCollect := dcm.GetDbSession().Table(dcm.DcUserDyCollect{})
	defer dbCollect.Close()
	exist, err := dbCollect.Where("id=? and status=? and user_id=?", id, 1, userId).Exist()
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

//修改收藏分组
func (receiver *UserBusiness) UpdCollectTag(id, tagId int) (comErr global.CommonError) {
	dbCollect := dcm.GetDbSession().Table(dcm.DcUserDyCollect{})
	defer dbCollect.Close()
	affect, err := dcm.UpdateInfo(dbCollect, id, map[string]interface{}{"tag_id": tagId, "update_time": time.Now()}, new(dcm.DcUserDyCollect))
	if err != nil || affect == 0 {
		comErr = global.NewError(5000)
		return
	}
	return
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

//获取分组列表
func (receiver *UserBusiness) GetDyCollectTags(userId int) (tags []dcm.DcUserDyCollectTag, comErr global.CommonError) {
	db := dcm.GetDbSession().Table(dcm.DcUserDyCollectTag{})
	if err := db.Where("user_id=? and delete_time is NULL", userId).Find(&tags); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//创建分组
func (receiver *UserBusiness) AddDyCollectTag(userId int, name string) (comErr global.CommonError) {
	db := dcm.GetDbSession().Table(dcm.DcUserDyCollectTag{})
	tag := dcm.DcUserDyCollectTag{
		Name:       name,
		UserId:     userId,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if _, err := db.Insert(tag); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//编辑分组
func (receiver *UserBusiness) UpdDyCollectTag(id int, name string) (comErr global.CommonError) {
	db := dcm.GetDbSession().Table(dcm.DcUserDyCollectTag{})
	if _, err := db.Where("id=?", id).Update(map[string]interface{}{"name": name}); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}

//删除分组
func (receiver *UserBusiness) DelDyCollectTag(id int) (comErr global.CommonError) {
	db := dcm.GetDbSession().Table(dcm.DcUserDyCollectTag{})
	if _, err := db.Where("id=?", id).Update(map[string]interface{}{"delete_time": time.Now()}); err != nil {
		comErr = global.NewError(5000)
		return
	}
	return
}
