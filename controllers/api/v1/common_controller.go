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
	"dongchamao/services/ali_tools"
	"dongchamao/services/dyimg"
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
	sig := InputData.GetString("sig", "")
	sessionId := InputData.GetString("session_id", "")
	token := InputData.GetString("token", "")
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
	if receiver.AppId == 10000 && global.IsDev() {
		if sig == "" || sessionId == "" || token == "" {
			receiver.FailReturn(global.NewError(4000))
			return
		}
		scene := "nc_message"
		if receiver.AppId != 10000 {
			scene = "nc_message_h5"
		}
		appKey := "FFFF0N0000000000A2FA"
		err1 := ali_tools.IClientProfile(sig, sessionId, token, receiver.Ip, scene, appKey)
		if err1 != nil {
			receiver.FailReturn(global.NewError(4000))
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

func (receiver *CommonController) GetConfigList() {
	var ret = make(map[string]interface{}, 0)
	cacheKey := cache.GetCacheKey(cache.ConfigKeyCache, "all")
	cacheData := global.Cache.Get(cacheKey)
	if cacheData != "" {
		cacheData = utils.DeserializeData(cacheData)
		_ = jsoniter.Unmarshal([]byte(cacheData), &ret)
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
	receiver.SuccReturn(ret)
	return
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
		start, _ := time.ParseInLocation("20060102", time.Now().Format("20060102"), time.Local)
		for i := 0; i < 7; i++ {
			dateTime := start.AddDate(0, 0, -i)
			date := dateTime.Format("2006-01-02")
			tmpList := make([]dy2.RedAuthorRoom, 0)
			roomList := authorBusiness.RedAuthorRoomByDate(authorIds, dateTime.Format("20060102"))
			if len(roomList) == 0 {
				continue
			}
			roomIds := []string{}
			for _, v := range roomList {
				roomIds = append(roomIds, v.RoomId)
			}
			roomInfos, _ := hbase.GetLiveInfoByIds(roomIds)
			for _, v := range roomList {
				if a, ok := authorDataMap[v.AuthorId]; ok {
					v.RoomCount = a.LiveCount
					v.AuthorLivingRoomId = a.RoomId
				}
				if weight, ok := authorSortMap[v.AuthorId]; ok {
					v.Weight = weight
				}
				if r, exist := roomInfos[v.RoomId]; exist {
					v.Gmv = r.TotalGmv
				}
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
		liveList := es.NewEsLiveBusiness().GetRoomsByAuthorIds(authorIds, time.Now().Format("20060102"), 3)
		for _, v := range liveList {
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
	if keyword == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
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
