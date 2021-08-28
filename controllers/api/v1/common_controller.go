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
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type CommonController struct {
	controllers.ApiBaseController
}

func (receiver *CommonController) Sms() {
	InputData := receiver.InputFormat()
	grantType := InputData.GetString("grant_type", "")
	mobile := InputData.GetString("mobile", receiver.UserInfo.Username)
	if !utils.VerifyMobileFormat(mobile) {
		receiver.FailReturn(global.NewError(4205))
		return
	}
	if !utils.InArrayString(grantType, []string{"login", "findpwd", "change_mobile", "bind_mobile"}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	limitIpKey := cache.GetCacheKey(cache.SmsCodeLimitBySome, "ip", receiver.Ip)
	verifyData := global.Cache.Get(limitIpKey)
	if verifyData != "" {
		receiver.FailReturn(global.NewError(4211))
		return
	}
	limitMobileKey := cache.GetCacheKey(cache.SmsCodeLimitBySome, "mobile", mobile)
	verifyData = global.Cache.Get(limitMobileKey)
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
	err := global.Cache.Set(cacheKey, code, 300)
	if logger.CheckError(err) != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	res, smsErr := aliSms.SmsCode(mobile, code)
	if !res || logger.CheckError(smsErr) != nil {
		receiver.FailReturn(global.NewError(6000))
		return
	}
	global.Cache.Set(limitIpKey, "1", 60)
	global.Cache.Set(limitMobileKey, "1", 60)
	receiver.SuccReturn(nil)
	return
}

//验证码校验
func (receiver *CommonController) CheckSmsCode() {
	mobile := receiver.GetString(":username", "")
	grantType := receiver.GetString(":grant_type", "")
	if !utils.VerifyMobileFormat(mobile) {
		receiver.FailReturn(global.NewError(4205))
		return
	}
	if !utils.InArrayString(grantType, []string{"findpwd", "change_mobile"}) {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	code := receiver.GetString(":code", "")
	codeKey := cache.GetCacheKey(cache.SmsCodeVerify, grantType, mobile)
	verifyCode := global.Cache.Get(codeKey)
	if verifyCode != code {
		receiver.FailReturn(global.NewError(4209))
		return
	}
	//换绑手机号时验证旧手机凭证
	if grantType == "change_mobile" {
		valid := utils.GetRandomString(10)
		cacheKey := cache.GetCacheKey(cache.OldMobileVerify, valid)
		err := global.Cache.Set(cacheKey, mobile, 300)
		if logger.CheckError(err) != nil {
			receiver.FailReturn(global.NewError(5000))
			return
		}
		receiver.SuccReturn(map[string]interface{}{"valid": valid})
		return
	}
	receiver.SuccReturn(nil)
	return
}

func (receiver *CommonController) IdEncryptDecrypt() {
	id := receiver.Ctx.Input.Param(":id")
	id1 := ""
	if strings.Index(id, "=") < 0 {
		id1 = business.IdEncrypt(id)
	}
	id2 := business.IdDecrypt(id)
	receiver.SuccReturn(map[string]string{
		"id":      id,
		"encrypt": id1,
		"decrypt": id2,
	})
	return
}

func (receiver *CommonController) Test() {
	return
}

func (receiver *CommonController) GetConfig() {
	var configJson dcm.DcConfigJson
	keyName := receiver.GetString(":key_name")
	exist, err := dcm.GetBy("key_name", keyName, &configJson)
	if !exist || err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	if configJson.Auth == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(configJson.Value), &data); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(data)
	return
}

func (receiver *CommonController) GetConfigList() {
	var config []dcm.DcConfigJson
	if err := dcm.GetDbSession().Table(dcm.DcConfigJson{}).Where("auth=1").Find(&config); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	var data = make([]map[string]interface{}, len(config))
	utils.MapToStruct(config, &data)
	var ret = make(map[string]map[string]interface{}, len(config))
	for _, v := range data {
		var jsonMap map[string]interface{}
		if err := json.Unmarshal([]byte(v["Value"].(string)), &jsonMap); err != nil {
			receiver.FailReturn(global.NewError(5000))
			return
		}
		ret[v["KeyName"].(string)] = jsonMap
	}
	receiver.SuccReturn(ret)
	return
}

func (receiver *CommonController) RedAuthorRoom() {
	listType := receiver.Ctx.Input.Param(":type")
	sql := "status = 1"
	if listType == "advance" {
		sql += fmt.Sprintf(" AND living_time > '%s' ", time.Now().Format("2006-01-02 15:04:05"))
	} else {
		sql += fmt.Sprintf(" AND living_time <= '%s' AND living_time >'%s' ", time.Now().Format("2006-01-02 15:04:05"), time.Now().AddDate(0, 0, -6).Format("2006-01-02 15:04:05"))
	}
	list := make([]dcm.DcAuthorRoom, 0)
	_ = dcm.GetSlaveDbSession().
		Table(new(dcm.DcAuthorRoom)).
		Where(sql).
		Desc("weight").
		Find(&list)
	if listType == "advance" {
		data := make([]dy2.RedAuthorRoom, 0)
		for _, v := range list {
			authorData, _ := hbase.GetAuthor(v.AuthorId)
			data = append(data, dy2.RedAuthorRoom{
				AuthorId:   business.IdEncrypt(v.AuthorId),
				Avatar:     dyimg.Fix(authorData.Data.Avatar),
				Sign:       authorData.Data.Signature,
				Nickname:   authorData.Data.Nickname,
				LivingTime: v.LivingTime.Unix(),
				RoomId:     business.IdEncrypt(authorData.RoomId),
				Tags:       authorData.Tags,
				RoomCount:  authorData.RoomCount,
			})
		}
		receiver.SuccReturn(map[string]interface{}{
			"list":  data,
			"total": len(list),
		})
		return
	}
	liveBusiness := business.NewLiveBusiness()
	dateMap := map[string][]dy2.RedAuthorRoom{}
	for _, v := range list {
		date := v.LivingTime.Format("2006-01-02")
		if _, ok := dateMap[date]; !ok {
			dateMap[date] = []dy2.RedAuthorRoom{}
		}
		authorData, _ := hbase.GetAuthor(v.AuthorId)
		var gmv, sales float64
		if v.RoomId == "" {
			rooms, _ := hbase.GetAuthorRoomsByDate(v.AuthorId, v.LivingTime.Format("20060102"))
			room := entity.DyAuthorLiveRoom{}
			for _, r := range rooms {
				if r.CreateTime > v.LivingTime.Unix() {
					room = r
					break
				}
				if r.CreateTime > room.CreateTime {
					room = r
				}
			}
			v.RoomId = room.RoomID
			if room.RoomID != "" {
				go func(id int, roomId string) {
					_, _ = dcm.UpdateInfo(nil, id, map[string]interface{}{"room_id": roomId}, new(dcm.DcAuthorRoom))
				}(v.Id, room.RoomID)
			}
		}
		gmv, sales = liveBusiness.LiveSalesData(v.RoomId)
		dateMap[date] = append(dateMap[date], dy2.RedAuthorRoom{
			AuthorId:           business.IdEncrypt(v.AuthorId),
			Sign:               authorData.Data.Signature,
			Avatar:             dyimg.Fix(authorData.Data.Avatar),
			Nickname:           authorData.Data.Nickname,
			AuthorLivingRoomId: business.IdEncrypt(authorData.RoomId),
			RoomId:             business.IdEncrypt(v.RoomId),
			Gmv:                gmv,
			Sales:              sales,
			Tags:               authorData.Tags,
			RoomCount:          authorData.RoomCount,
		})
	}
	data := make([]dy2.RedAuthorRoomBox, 0)
	for k, v := range dateMap {
		data = append(data, dy2.RedAuthorRoomBox{
			Date: k,
			List: v,
		})
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  data,
		"total": len(list),
	})
	return
}

//抖音首页查询
func (receiver *CommonController) DyUnionSearch() {
	if !business.UserActionLock(receiver.TrueUri, receiver.Ip, 2) {
		receiver.FailReturn(global.NewError(4211))
		return
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
	ret := map[string]interface{}{
		"author":  authorList,
		"live":    liveList,
		"product": productList,
	}
	receiver.SuccReturn(ret)
	return

}
