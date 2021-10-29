package v1

import (
	"dongchamao/business"
	"dongchamao/business/es"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/logger"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	"dongchamao/models/entity"
	dy2 "dongchamao/models/repost/dy"
	"dongchamao/services/ali_sms"
	"dongchamao/services/dyimg"
	"dongchamao/services/tencent"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"math"
	"sort"
	"strings"
	"time"
)

type CommonController struct {
	controllers.ApiBaseController
}

func (receiver *CommonController) Sms() {
	InputData := receiver.InputFormat()
	grantType := InputData.GetString("grant_type", "")
	mobile := InputData.GetString("mobile", "")
	if !utils.InArrayString(grantType, []string{"login", "findpwd", "change_mobile", "bind_mobile"}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if utils.InArrayString(grantType, []string{"change_mobile"}) {
		mobile = receiver.UserInfo.Username
	}
	if !utils.VerifyMobileFormat(mobile) {
		receiver.FailReturn(global.NewError(4205))
		return
	}
	if business.NewAccountAuthBusiness().CheckSmsSend(receiver.Clientos, receiver.AppVersion) {
		ticket := receiver.InputFormat().GetString("ticket", "")
		randStr := receiver.InputFormat().GetString("randstr", "")
		if !tencent.TencentCaptcha(ticket, randStr, receiver.Ip) {
			receiver.FailReturn(global.NewError(8001))
			return
		}
	}
	//limitIpKey := cache.GetCacheKey(cache.SmsCodeLimitBySome, grantType, "ip", receiver.Ip)
	//verifyData := global.Cache.Get(limitIpKey)
	//if verifyData != "" {
	//	receiver.FailReturn(global.NewError(4211))
	//	return
	//}
	if grantType == "bind_mobile" {
		var user dcm.DcUser
		exist, _ := dcm.GetBy("username", mobile, &user)
		if exist && user.Openid != "" {
			receiver.FailReturn(global.NewError(4215))
			return
		}
	}
	limitMobileKey := cache.GetCacheKey(cache.SmsCodeLimitBySome, grantType, "mobile", mobile)
	verifyData := global.Cache.Get(limitMobileKey)
	if verifyData != "" {
		receiver.FailReturn(global.NewError(4211))
		return
	}
	if utils.InArrayString(grantType, []string{"findpwd"}) {
		user := dcm.DcUser{}
		exist, _ := dcm.GetBy("username", mobile, &user)
		if !exist {
			receiver.FailReturn(global.NewError(4204))
			return
		}
	}
	cacheKey := cache.GetCacheKey(cache.SmsCodeVerify, grantType, mobile)
	code := utils.GetRandomInt(6)
	res, smsErr := aliSms.SmsCode(mobile, code)
	if !res || logger.CheckError(smsErr) != nil {
		business.NewMonitorBusiness().SendErr("短信验证码错误", fmt.Sprintf("env:%s,error:%s,mobile:%s,code:%s", global.IsDev(), smsErr.Error(), mobile, code))
		receiver.FailReturn(global.NewError(6000))
		return
	}
	err := global.Cache.Set(cacheKey, code, 300)
	if logger.CheckError(err) != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	//_ = global.Cache.Set(limitIpKey, "1", 60)
	_ = global.Cache.Set(limitMobileKey, "1", 60)
	receiver.SuccReturn(nil)
	return
}

//验证码校验
func (receiver *CommonController) CheckSmsCode() {
	mobile := receiver.GetString(":username", "")
	grantType := receiver.GetString(":grant_type", "")
	if !utils.InArrayString(grantType, []string{"findpwd", "change_mobile", "bind_mobile"}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if utils.InArrayString(grantType, []string{"change_mobile"}) {
		mobile = receiver.UserInfo.Username
	}
	if !utils.VerifyMobileFormat(mobile) {
		receiver.FailReturn(global.NewError(4205))
		return
	}
	code := receiver.GetString(":code", "")
	codeKey := cache.GetCacheKey(cache.SmsCodeVerify, grantType, mobile)
	verifyCode := global.Cache.Get(codeKey)
	if verifyCode != code {
		receiver.FailReturn(global.NewError(4209))
		return
	}
	receiver.SuccReturn(nil)
	return
}

func (receiver *CommonController) Test() {
	receiver.SuccReturn(nil)
	return
}

func (receiver *CommonController) GetConfig() {
	keyName := receiver.GetString(":key_name")
	cacheKey := cache.GetCacheKey(cache.ConfigKeyCache, keyName)
	cacheData := global.Cache.Get(cacheKey)
	var isString = true
	if cacheData == "" {
		var configJson dcm.DcConfigJson
		exist, err := dcm.GetBy("key_name", keyName, &configJson)
		if !exist || err != nil {
			receiver.FailReturn(global.NewError(5000))
			return
		}
		if configJson.Auth == 0 {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		if configJson.ContentType != 2 {
			isString = false
		}
		cacheData = configJson.Value
	}
	if isString {
		receiver.SuccReturn(cacheData)
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(cacheData), &data); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(data)
	return
}

func (receiver *CommonController) InvitePhone() {
	var configJson dcm.DcConfigJson
	if !receiver.HasLogin {
		receiver.FailReturn(global.NewError(4001))
		return
	}
	input := receiver.InputFormat()
	keyName := input.GetString("key_name", "")
	exist, err := dcm.GetDbSession().Where("key_name=? and auth=0", keyName).Get(&configJson)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	userPhone := input.GetString("user_phone", "")
	invitePhone := input.GetString("invite_phone", "")
	platform := input.GetString("platform", "")
	if userPhone == "" || invitePhone == "" || platform == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var value = make(map[string]map[string]string)
	json.Unmarshal([]byte(configJson.Value), &value)
	if value[userPhone] == nil {
		value[userPhone] = make(map[string]string)
	}
	value[userPhone][invitePhone] = platform
	marshal, _ := json.Marshal(value)
	configJson.KeyName = keyName
	configJson.Value = string(marshal)
	if !exist {
		if _, err = dcm.Insert(dcm.GetDbSession(), &configJson); err != nil {
			receiver.FailReturn(global.NewError(5000))
			return
		}
	} else {
		_, err = dcm.GetDbSession().Where("key_name=?", keyName).Update(&configJson)
		if err != nil {
			receiver.FailReturn(global.NewError(5000))
			return
		}
	}
	receiver.SuccReturn(nil)
	return
}

func (receiver *CommonController) GetConfigList() {
	var ret = make(map[string]interface{}, 0)
	cacheKey := cache.GetCacheKey(cache.ConfigKeyCache, "all")
	cacheData := global.Cache.Get(cacheKey)
	if cacheData != "" {
		cacheData = utils.DeserializeData(cacheData)
		_ = jsoniter.Unmarshal([]byte(cacheData), &ret)
		//校验是否下发支付
		iosPayOpen := 1
		if receiver.checkIosPay() {
			iosPayOpen = 0
		}
		ret["ios_pay"].(map[string]interface{})["ios_pay"] = iosPayOpen
		ret["ios_pay"].(map[string]interface{})["open"] = 1
		receiver.SuccReturn(ret)
		return
	}
	var configList []dcm.DcConfigJson
	if err := dcm.GetDbSession().Table(dcm.DcConfigJson{}).Where("auth=1").Find(&configList); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	for _, v := range configList {
		if v.ContentType != 1 {
			var jsonMap map[string]interface{}
			if err := json.Unmarshal([]byte(v.Value), &jsonMap); err != nil {
				receiver.FailReturn(global.NewError(5000))
				return
			}
			ret[v.KeyName] = jsonMap
		} else {
			ret[v.KeyName] = v.Value
		}
	}
	_ = global.Cache.Set(cacheKey, utils.SerializeData(ret), 300)
	//校验是否下发支付
	iosPayOpen := 1
	if receiver.checkIosPay() {
		iosPayOpen = 0
	}
	ret["ios_pay"] = map[string]interface{}{
		"ios_pay": iosPayOpen,
		"open":    1,
	}
	receiver.SuccReturn(ret)
	return
}

func (receiver *CommonController) checkIosPay() bool {
	checkIosVersion := business.GetConfig("check_ios_pay")
	iosVersion := receiver.Ctx.Input.Header("APPVERSION")
	if iosVersion == checkIosVersion {
		return true
	}
	return false
}

func (receiver *CommonController) RedAuthorRoom() {
	listType := receiver.Ctx.Input.Param(":type")
	sql := "status = 1"
	if listType == "advance" {
		sql += fmt.Sprintf(" AND living_time > '%s' ", time.Now().Format("2006-01-02 15:04:05"))
	}
	list := make([]dcm.DcAuthorRoom, 0)
	_ = dcm.GetSlaveDbSession().
		Table(new(dcm.DcAuthorRoom)).
		Where(sql).
		Desc("weight").
		Find(&list)
	authorIds := make([]string, 0)
	authorSortMap := map[string]int{}
	for _, v := range list {
		authorSortMap[v.AuthorId] = v.Weight
		authorIds = append(authorIds, v.AuthorId)
	}
	authorBusiness := business.NewAuthorBusiness()
	authorCacheKey := cache.GetCacheKey(cache.RedAuthorMapCache, utils.Md5_encode(strings.Join(authorIds, "")))
	authorCacheData := global.Cache.Get(authorCacheKey)
	var authorDataMap = map[string]entity.DyAuthorSimple{}
	if authorCacheData != "" {
		authorCacheData = utils.DeserializeData(authorCacheData)
		_ = jsoniter.Unmarshal([]byte(authorCacheData), &authorDataMap)
	} else {
		authorMap, _ := hbase.GetAuthorByIds(authorIds)
		utils.MapToStruct(authorMap, &authorDataMap)
		_ = global.Cache.Set(authorCacheKey, utils.SerializeData(authorDataMap), 600)
	}
	if listType == "advance" {
		data := make([]dy2.RedAuthorRoom, 0)
		today := time.Now().Format("20060102")
		for _, v := range list {
			authorData := entity.DyAuthorSimple{}
			if a, ok := authorDataMap[v.AuthorId]; ok {
				authorData = a
			}
			if authorData.RoomId != "" && v.LivingTime.Format("2006012") == today {
				continue
			}
			data = append(data, dy2.RedAuthorRoom{
				AuthorId:           business.IdEncrypt(v.AuthorId),
				Avatar:             dyimg.Fix(authorData.Data.Avatar),
				Sign:               authorData.Data.Signature,
				Nickname:           authorData.Data.Nickname,
				LivingTime:         v.LivingTime.Unix(),
				AuthorLivingRoomId: business.IdEncrypt(authorData.RoomId),
				RoomId:             business.IdEncrypt(authorData.RoomId),
				Tags:               authorData.Tags,
				RoomCount:          authorData.LiveCount,
			})
		}
		receiver.SuccReturn(map[string]interface{}{
			"list":  data,
			"total": len(list),
		})
		return
	}
	data := make([]dy2.RedAuthorRoomBox, 0)
	total := 0
	if len(list) > 0 {
		today := time.Now().Format("2006-01-02")
		start, _ := time.ParseInLocation("20060102", time.Now().Format("20060102"), time.Local)
		for i := 0; i < 7; i++ {
			dateTime := start.AddDate(0, 0, -i)
			date := dateTime.Format("2006-01-02")
			tmpList := make([]dy2.RedAuthorRoom, 0)
			roomList := authorBusiness.RedAuthorRoomByDate(authorIds, dateTime.Format("20060102"))
			if date == today && len(roomList) == 0 {
				start = start.AddDate(0, 0, -1)
				date = start.Format("2006-01-02")
				roomList = authorBusiness.RedAuthorRoomByDate(authorIds, start.Format("20060102"))
			}
			if len(roomList) == 0 {
				continue
			}
			roomInfos := map[string]entity.DyLiveInfo{}
			hasLiving := false
			roomIds := []string{}
			for _, v := range roomList {
				roomIds = append(roomIds, v.RoomId)
				if v.RoomStatus == 2 {
					hasLiving = true
				}
			}
			if hasLiving {
				roomInfos, _ = hbase.GetLiveInfoByIds(roomIds)
			}
			for _, v := range roomList {
				if a, ok := authorDataMap[v.AuthorId]; ok {
					v.RoomCount = a.LiveCount
					v.AuthorLivingRoomId = a.RoomId
				}
				if weight, ok := authorSortMap[v.AuthorId]; ok {
					v.Weight = weight
				}
				if r, exist := roomInfos[v.RoomId]; exist {
					v.Gmv = r.PredictGmv
					v.TotalUser = r.WatchCnt
					v.RoomStatus = r.RoomStatus
				}
				v.AuthorId = business.IdEncrypt(v.AuthorId)
				v.RoomId = business.IdEncrypt(v.RoomId)
				tmpList = append(tmpList, v)
				total++
			}
			sort.Slice(tmpList, func(i, j int) bool {
				if tmpList[i].RoomStatus == 2 && tmpList[j].RoomStatus == 4 {
					return true
				}
				if tmpList[i].RoomStatus == 4 && tmpList[j].RoomStatus == 2 {
					return false
				}
				return tmpList[i].Weight > tmpList[j].Weight
			})
			data = append(data, dy2.RedAuthorRoomBox{
				Date: date,
				List: tmpList,
			})
		}
		sort.Slice(data, func(i, j int) bool {
			return data[i].Date > data[j].Date
		})
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  data,
		"total": total,
	})
	return
}

//红人看榜正在直播top3
func (receiver *CommonController) RedAuthorLivingRoom() {
	cacheKey := cache.GetCacheKey(cache.RedAuthorLivingRooms)
	cacheData := global.Cache.Get(cacheKey)
	list := make([]dcm.DcAuthorRoom, 0)
	if cacheData != "" {
		cacheData = utils.DeserializeData(cacheData)
		_ = jsoniter.Unmarshal([]byte(cacheData), &list)
		receiver.SuccReturn(map[string]interface{}{
			"list": list,
		})
		return
	}
	_ = dcm.GetSlaveDbSession().
		Table(new(dcm.DcAuthorRoom)).
		Where("status=?", 1).
		Desc("weight").
		Find(&list)
	data := make([]dy2.RedAuthorRoom, 0)
	if len(list) > 0 {
		authorIds := make([]string, 0)
		for _, v := range list {
			authorIds = append(authorIds, v.AuthorId)
		}
		dateStr := time.Now().Format("20060102")
		liveList := es.NewEsLiveBusiness().GetRoomsByAuthorIds(authorIds, dateStr, 3)
		if len(liveList) == 0 {
			dateStr = time.Now().AddDate(0, 0, -1).Format("20060102")
			liveList = es.NewEsLiveBusiness().GetRoomsByAuthorIds(authorIds, dateStr, 3)
		}
		roomInfos := map[string]entity.DyLiveInfo{}
		hasLiving := false
		roomIds := []string{}
		for _, v := range liveList {
			roomIds = append(roomIds, v.RoomId)
			if v.RoomStatus == 2 {
				hasLiving = true
			}
		}
		if hasLiving {
			roomInfos, _ = hbase.GetLiveInfoByIds(roomIds)
		}
		for _, v := range liveList {
			if r, exist := roomInfos[v.RoomId]; exist {
				v.PredictGmv = r.PredictGmv
				v.PredictSales = r.PredictSales
				v.WatchCnt = r.WatchCnt
			}
			data = append(data, dy2.RedAuthorRoom{
				AuthorId:   business.IdEncrypt(v.AuthorId),
				Avatar:     dyimg.Fix(v.Avatar),
				Nickname:   v.Nickname,
				LiveTitle:  v.Title,
				RoomId:     business.IdEncrypt(v.RoomId),
				RoomStatus: v.RoomStatus,
				Gmv:        v.PredictGmv,
				Sales:      math.Floor(v.PredictSales),
				TotalUser:  v.WatchCnt,
				Tags:       v.Tags,
				CreateTime: v.CreateTime,
			})
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"list": data,
	})
	return
}

//抖音首页查询
func (receiver *CommonController) DyUnionSearch() {
	if !receiver.HasLogin {
		if !business.UserActionLock(receiver.TrueUri, receiver.Ip, 2) {
			receiver.FailReturn(global.NewError(4211))
			return
		}
	}
	keyword := receiver.GetString("keyword", "")
	//if keyword == "" {
	//	receiver.FailReturn(global.NewError(4000))
	//	return
	//}
	receiver.KeywordBan(keyword)
	authorList := es.NewEsAuthorBusiness().KeywordSearch(keyword)
	for k, v := range authorList {
		authorData, _ := hbase.GetAuthor(v.AuthorId)
		authorList[k].AuthorId = business.IdEncrypt(v.AuthorId)
		authorList[k].Avatar = dyimg.Fix(v.Avatar)
		authorList[k].RoomId = business.IdEncrypt(authorData.RoomId)
		if v.UniqueId == "" || v.UniqueId == "0" {
			authorList[k].UniqueId = v.ShortId
		}
	}
	liveList := es.NewEsLiveBusiness().KeywordSearch(keyword)
	for k, v := range liveList {
		liveList[k].AuthorId = business.IdEncrypt(v.AuthorId)
		liveList[k].RoomId = business.IdEncrypt(v.RoomId)
		if v.DisplayId == "" || v.DisplayId == "0" {
			liveList[k].DisplayId = v.ShortId
		}
	}
	productList := es.NewEsProductBusiness().KeywordSearch(keyword)
	for k, v := range productList {
		productList[k].ProductId = business.IdEncrypt(v.ProductId)
	}
	shopList, _, _ := es.NewEsShopBusiness().SimpleSearch(keyword, "", "", "", 1, 2, "", "")
	for k, v := range shopList {
		shopList[k].ShopId = business.IdEncrypt(v.ShopId)
	}
	ret := map[string]interface{}{
		"author":  authorList,
		"live":    liveList,
		"product": productList,
		"shop":    shopList,
	}
	receiver.SuccReturn(ret)
	return

}

//获取当前版本接口
func (receiver *CommonController) CheckAppVersion() {
	appType := receiver.Ctx.Input.Param(":type")
	platform := 0
	if strings.ToLower(appType) == "ios" {
		platform = 1
	} else if strings.ToLower(appType) == "android" {
		platform = 2
	} else {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	row := dcm.DcAppVersion{}
	_, _ = dcm.GetSlaveDbSession().
		Where("platform=?", platform).
		Desc("id").
		Get(&row)
	receiver.SuccReturn(dy2.AppVersion{
		Version: row.Version,
		Info:    row.Info,
		Force:   row.Force,
		Url:     row.Url,
	})
}

//获取时间戳
func (receiver *CommonController) CheckTime() {
	receiver.SuccReturn(map[string]interface{}{
		"time": time.Now().Unix(),
	})
}

//获取当前版本接口
func (receiver *CommonController) CountChannelClick() {
	if receiver.Channel != "" {
		clickLog := dcm.DcUserChannelLogs{
			UserId:      receiver.UserId,
			Channel:     receiver.Channel,
			ChannelWord: receiver.ChannelWords,
			AppId:       receiver.AppId,
			Ip:          receiver.Ip,
			CreateTime:  time.Now(),
		}
		_, _ = dcm.Insert(nil, &clickLog)
	}
	receiver.SuccReturn(nil)
}

//滑块处理
func (receiver *CommonController) ClearAcfVerify() {
	ticket := receiver.InputFormat().GetString("ticket", "")
	randStr := receiver.InputFormat().GetString("randstr", "")
	if !tencent.TencentCaptcha(ticket, randStr, receiver.Ip) {
		receiver.FailReturn(global.NewError(8001))
		return
	}
	if receiver.UserId > 0 {
		_ = global.Cache.Delete(cache.GetCacheKey(cache.SecurityVerifyCodeUid, receiver.UserId))
	}
	_ = global.Cache.Delete(cache.GetCacheKey(cache.SecurityVerifyCodeIp, receiver.Ip))
	receiver.SuccReturn(nil)
	return
}
