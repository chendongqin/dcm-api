package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//用户等级
const (
	UserLevelDefault = 0 //普通会员
	UserLevelVip     = 1 //vip
	UserLevelSvip    = 2 //svip
	UserLevelJewel   = 3 //专业版
)

const (
	VipPlatformDouYin      = 1 //抖音
	VipPlatformXiaoHongShu = 2 //小红书
	VipPlatformTaoBao      = 3 //淘宝
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

//获取会员等级
func (this *UserBusiness) GetUserLevels() map[int]string {
	levels := map[int]string{}
	levels[UserLevelDefault] = "普通会员"
	levels[UserLevelVip] = "vip"
	levels[UserLevelSvip] = "svip"
	levels[UserLevelJewel] = "专业版"
	return levels
}

//获取等级名称
func (receiver *UserBusiness) GetUserLevel(level int) string {
	userLevels := receiver.GetUserLevels()
	if v, ok := userLevels[level]; ok {
		return v
	}
	return ""
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
	exist, _ := dcm.GetDbSession().Where("user_id = ?", userId).
		And("app_platform = ?", platFormId).
		Get(userTokenM)

	if exist {
		_ = global.Cache.Set(cKey, userTokenM.Token, 1800)
	}

	return userTokenM.Token, exist
}

//获取用户vip等级
func (receiver *UserBusiness) GetVipLevel(userId int) map[int]int {
	vipLists := make([]dcm.DcUserVip, 0)
	err := dcm.GetDbSession().Where("user_id=? ", userId).Find(&vipLists)
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

//更新会员等级
func (receiver *UserBusiness) UpdateValidDayOne(userId, platformId int) (int, bool) {
	vipModel := &dcm.DcUserVip{}
	exist, _ := dcm.GetDbSession().Where("user_id=? AND platform=?", userId, platformId).Get(vipModel)
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
	affect, err := dcm.GetDbSession().Table(new(dcm.DcUserVip)).Where(whereStr, vipModel.Id, time.Now().Format("2006-01-02 15:04:05")).Update(updateData)
	if affect == 0 || err != nil {
		return 0, false
	}
	return vipModel.OrderLevel, true
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