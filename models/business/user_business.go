package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"github.com/dgrijalva/jwt-go"
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
func (receiver *UserBusiness) SmsLogin(mobile, code string, appId int) (user dcm.DcUser, tokenString string, expire int64, comErr global.CommonError) {
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
		affect, err := dcm.Insert(nil, &user)
		if affect == 0 || err != nil {
			comErr = global.NewError(5000)
			return
		}
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
	return
}

//收藏达人
func (receiver *UserBusiness) CollectAuthor(userId, authorId string, platform int) (comErr global.CommonError) {
	dbSession := dcm.GetDbSession()
	collect := dcm.DcUserCollect{}
	exist, err := dbSession.Where("user_id=? AND platform=? AND author_id=?", userId, platform, authorId).Get(&collect)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	if exist || collect.Status == 1 {
		comErr = global.NewMsgError("您已收藏该达人，请刷新重试")
		return
	}
	//达人查询写入

	return
}

//取消收藏
func (receiver *UserBusiness) CancelCollectAuthor(userId, authorId string, platform int) (comErr global.CommonError) {
	dbSession := dcm.GetDbSession()
	collect := dcm.DcUserCollect{}
	exist, err := dbSession.Where("user_id=? AND platform=? AND author_id=?", userId, platform, authorId).Get(&collect)
	if err != nil {
		comErr = global.NewError(5000)
		return
	}
	if !exist || collect.Status == 0 {
		comErr = global.NewMsgError("您未收藏了该达人，请刷新重试")
		return
	}
	affect, err := dcm.UpdateInfo(dbSession, collect.Id, map[string]interface{}{"status": 0}, new(dcm.DcUserCollect))
	if err != nil || affect == 0 {
		comErr = global.NewError(5000)
		return
	}
	return
}
