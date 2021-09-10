package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/logger"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/models/repost"
	"dongchamao/models/repost/dy"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"
)

type AccountController struct {
	controllers.ApiBaseController
}

func (receiver *AccountController) Prepare() {
	receiver.InitApiController()
	receiver.CheckToken()
}

//重置密码
func (receiver *AccountController) ResetPwd() {
	InputData := receiver.InputFormat()
	if receiver.UserInfo.SetPassword == 1 {
		oldPwd := InputData.GetString("old_pwd", "")
		if oldPwd == "" {
			receiver.FailReturn(global.NewError(4214))
			return
		}
		oldPwd = utils.Base64Decode(oldPwd)
		if utils.Md5_encode(oldPwd+receiver.UserInfo.Salt) != receiver.UserInfo.Password {
			receiver.FailReturn(global.NewError(4214))
			return
		}
	}
	newPwd := InputData.GetString("new_pwd", "")
	newPwd = utils.Base64Decode(newPwd)
	surePwd := InputData.GetString("sure_pwd", "")
	surePwd = utils.Base64Decode(surePwd)
	if newPwd == "" || surePwd == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	if newPwd != surePwd {
		receiver.FailReturn(global.NewError(4207))
		return
	}
	pwdLen := strings.Count(newPwd, "")
	if pwdLen > 24 || pwdLen < 6 {
		receiver.FailReturn(global.NewError(4210))
		return
	}
	pwd := utils.Md5_encode(newPwd + receiver.UserInfo.Salt)
	userBusiness := business.NewUserBusiness()
	updateData := map[string]interface{}{
		"password":     pwd,
		"set_password": 1,
		"update_time":  utils.GetNowTimeStamp(),
	}
	affect, _ := userBusiness.UpdateUserAndClearCache(nil, receiver.UserId, updateData)
	if affect == 0 {
		receiver.FailReturn(global.NewError(4213))
		return
	}
	receiver.Logout()
	receiver.SuccReturn(nil)
	return
}

//绑定、修改手机号
func (receiver *AccountController) ChangeMobile() {
	InputData := receiver.InputFormat()
	mobile := InputData.GetString("mobile", "")
	oldCode := InputData.GetString("old_code", "")
	code := InputData.GetString("code", "")
	userBusiness := business.NewUserBusiness()
	//新手机号存在校验
	exist, comErr := userBusiness.MobileExist(mobile)
	if comErr != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	if exist {
		receiver.FailReturn(global.NewMsgError("该手机号已存在"))
		return
	}
	//旧手机验证
	codeKey := cache.GetCacheKey(cache.SmsCodeVerify, "change_mobile", receiver.UserInfo.Username)
	verifyCode := global.Cache.Get(codeKey)
	if verifyCode != oldCode {
		receiver.FailReturn(global.NewError(4209))
		return
	}
	//新手机验证码校验
	codeKey1 := cache.GetCacheKey(cache.SmsCodeVerify, "bind_mobile", mobile)
	verifyCode1 := global.Cache.Get(codeKey1)
	if verifyCode1 != code {
		receiver.FailReturn(global.NewError(4209))
		return
	}
	//修改手机号
	dbSession := dcm.GetDbSession()
	_ = dbSession.Begin()
	updateMap := map[string]interface{}{
		"username": mobile,
	}
	if strings.Index(receiver.UserInfo.Nickname, "****") >= 0 {
		updateMap["nickname"] = mobile[:3] + "****" + mobile[7:]
	}
	affect, err := userBusiness.UpdateUserAndClearCache(dbSession, receiver.UserId, updateMap)
	if affect == 0 || logger.CheckError(err) != nil {
		_ = dbSession.Rollback()
		receiver.FailReturn(global.NewError(5000))
		return
	}
	_, err = dbSession.Table(new(dcm.DcVipOrder)).Where("user_id=?", receiver.UserId).Update(map[string]interface{}{"username": mobile})
	if logger.CheckError(err) != nil {
		_ = dbSession.Rollback()
		receiver.FailReturn(global.NewError(5000))
		return
	}
	_ = dbSession.Commit()
	_ = global.Cache.Delete(codeKey)
	_ = global.Cache.Delete(codeKey1)
	receiver.Logout()
	receiver.SuccReturn(nil)
	return
}

//手机号存在校验
func (receiver *AccountController) MobileExist() {
	mobile := receiver.GetString("mobile", "")
	//新手机号存在校验
	exist, comErr := business.NewUserBusiness().MobileExist(mobile)
	if comErr != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(exist)
	return
}

//info
func (receiver *AccountController) Info() {
	username := receiver.UserInfo.Username
	isWechat := 0
	if receiver.UserInfo.Unionid != "" {

	}
	usernameEncrypt := ""
	if len(username) >= 11 {
		usernameEncrypt = username[:3] + "****" + username[7:]
	}
	account := dy.RepostAccountData{
		UserId:      receiver.UserInfo.Id,
		Username:    usernameEncrypt,
		Nickname:    receiver.UserInfo.Nickname,
		Avatar:      receiver.UserInfo.Avatar,
		PasswordSet: receiver.UserInfo.SetPassword,
		Wechat:      isWechat,
	}
	vipBusiness := business.NewVipBusiness()
	vipLevels := vipBusiness.GetVipLevels(receiver.UserInfo.Id)
	for _, v := range vipLevels {
		expiration := "-"
		subExpiration := "-"
		if v.ExpirationTime.After(time.Now()) {
			expiration = v.ExpirationTime.Format("2006-01-02")
		}
		if v.SubExpirationTime.After(time.Now()) {
			subExpiration = v.SubExpirationTime.Format("2006-01-02 15:04:05")
		}
		vipLevel := dy.RepostAccountVipLevel{
			Level:             v.Level,
			LevelName:         vipBusiness.GetUserLevel(v.Level),
			ExpirationTime:    expiration,
			SubNum:            v.SubNum,
			IsSub:             v.IsSub,
			SubExpirationTime: subExpiration,
		}
		if v.PlatForm == business.VipPlatformDouYin {
			account.DyLevel = vipLevel
		} else if v.PlatForm == business.VipPlatformXiaoHongShu {
			account.XhsLevel = vipLevel
		} else if v.PlatForm == business.VipPlatformTaoBao {
			account.TbLevel = vipLevel
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"has_auth":  receiver.HasAuth,
		"has_login": receiver.HasLogin,
		"info":      account,
	})
	return
}

//登出
func (receiver *AccountController) Logout() {
	cacheKey := cache.GetCacheKey(cache.UserPlatformUniqueToken, receiver.AppId, receiver.UserId)
	_ = global.Cache.Delete(cacheKey)
	//执行登出事件
	receiver.RegisterLogout()
	//uniquetoken更新置为空  旧的token不可用
	userBusiness := business.NewUserBusiness()
	_ = userBusiness.AddOrUpdateUniqueToken(receiver.UserId, receiver.AppId, "")
	userBusiness.DeleteUserInfoCache(receiver.UserInfo.Id)
	receiver.SuccReturn("success")
	return

}

func (receiver *AccountController) DyUserSearchSave() {
	searchType := receiver.GetString(":type")
	data := receiver.ApiDatas
	dataMap, _ := utils.ToMapStringInterface(data)
	searchData := map[string]interface{}{}
	note := ""
	for k, v := range dataMap {
		if k == "note" {
			note = utils.ToString(v)
			continue
		}
		searchData[k] = v
	}
	if note == "" {
		receiver.FailReturn(global.NewMsgError("请输入筛选器昵称"))
		return
	}
	total, _ := dcm.GetSlaveDbSession().
		Table(new(dcm.DcUserSearch)).
		Where("user_id =? AND search_type = ?", receiver.UserId, searchType).
		Count()
	if total >= 10 {
		receiver.FailReturn(global.NewMsgError("最多保存10条常用筛选器"))
		return
	}
	contentByte, _ := jsoniter.Marshal(searchData)
	searchM := dcm.DcUserSearch{
		UserId:     receiver.UserId,
		SearchType: searchType,
		Note:       note,
		Content:    string(contentByte),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	affect, err := dcm.Insert(nil, &searchM)
	if affect == 0 || err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(nil)
	return
}

func (receiver *AccountController) DyUserSearchDel() {
	id := receiver.GetString(":id")
	dbSession := dcm.GetDbSession()
	searchM := dcm.DcUserSearch{}
	affect, err := dbSession.Where("id = ? AND user_id = ?", id, receiver.UserId).Delete(&searchM)
	if affect == 0 || err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(nil)
	return
}

func (receiver *AccountController) DyUserSearchList() {
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	searchType := receiver.GetString(":type")
	list := make([]dcm.DcUserSearch, 0)
	dbSession := dcm.GetSlaveDbSession()
	start := (page - 1) * pageSize
	total, err := dbSession.
		Where("user_id =? AND search_type = ?", receiver.UserId, searchType).
		Limit(pageSize, start).
		FindAndCount(&list)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	repostList := make([]repost.SearchData, 0)
	for _, v := range list {
		content := map[string]interface{}{}
		_ = jsoniter.Unmarshal([]byte(v.Content), &content)
		repostList = append(repostList, repost.SearchData{
			Id:         v.Id,
			SearchType: v.SearchType,
			Note:       v.Note,
			Content:    content,
		})
	}
	receiver.SuccReturn(map[string]interface{}{
		"list":  repostList,
		"total": total,
	})
	return
}

//判断是否收藏
func (receiver *AccountController) IsCollect() {
	platform, err := receiver.GetInt("platform", 1)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	collectType, err := receiver.GetInt("collect_type", 1)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	collectId := receiver.GetString("collect_id")
	if collectId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var id int
	switch platform {
	case 1:
		id = business.NewUserBusiness().DyCollectExist(collectType, receiver.UserId, collectId)
	}
	receiver.SuccReturn(map[string]interface{}{"is_collect": id})
	return
}

//添加收藏
func (receiver *AccountController) AddCollect() {
	//platform：1抖音
	platform, err := receiver.GetInt("platform", 1)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	tagId, err := receiver.GetInt("tag_id", 0)
	if err != nil {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	collectId := receiver.GetString("collect_id")
	if collectId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	collectType, err := receiver.GetInt("collect_type", 1)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	var comErr global.CommonError
	switch platform {
	case 1:
		comErr = business.NewUserBusiness().AddDyCollect(collectId, collectType, tagId, receiver.UserInfo.Id)
	}
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn("success")
	return
}

//取消收藏
func (receiver *AccountController) DelCollect() {
	id, err := receiver.GetInt(":id")
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	comErr := business.NewUserBusiness().CancelDyCollect(id, receiver.UserId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn("success")
	return
}

//获取收藏列表
func (receiver *AccountController) GetCollect() {
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page_size", 10, 50)
	platform, err := receiver.GetInt("platform", 1)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	collectType, err := receiver.GetInt("collect_type", 1)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	tagId, err := receiver.GetInt("tag_id", 0)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	label := receiver.GetString("label", "")
	keywords := receiver.GetString("keywords")
	switch platform {
	case 1:
		data, total, comErr := business.NewUserBusiness().GetDyCollect(tagId, collectType, keywords, label, receiver.UserId, page, pageSize)
		if comErr != nil {
			receiver.FailReturn(comErr)
			return
		}
		receiver.SuccReturn(map[string]interface{}{
			"list":     data,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		})
	}
	return
}

//更换收藏分组
func (receiver *AccountController) UpdCollectTag() {
	id, err := receiver.GetInt(":id")
	tagId, err := receiver.GetInt(":tag_id")
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	comErr := business.NewUserBusiness().UpdCollectTag(id, tagId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn("success")
	return
}

//抖音达人收藏分组
func (receiver *AccountController) GetDyCollectTags() {
	userBusiness := business.NewUserBusiness()
	data, comErr := userBusiness.GetDyCollectTags(receiver.UserId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	collectCount, comErr := userBusiness.GetDyCollectCount(receiver.UserId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	var total int64
	var countMap = make(map[int]int64, len(collectCount))
	for _, v := range collectCount {
		total += v.Count
		countMap[v.TagId] = v.Count
	}
	var countRet = make([]repost.CollectTagRet, len(data))
	for k, v := range data {
		countRet[k].DcUserDyCollectTag = v
		countRet[k].Count = countMap[v.Id]
	}
	receiver.SuccReturn(map[string]interface{}{"total": total, "list": countRet})
	return
}

func (receiver *AccountController) UpdDyCollectTag() {
	id, _ := receiver.GetInt(":id")
	InputData := receiver.InputFormat()
	name := InputData.GetString("name", "")
	comErr := business.NewUserBusiness().UpdDyCollectTag(id, name)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn("success")
	return
}

func (receiver *AccountController) AddDyCollectTag() {
	InputData := receiver.InputFormat()
	name := InputData.GetString("name", "")
	comErr := business.NewUserBusiness().AddDyCollectTag(receiver.UserId, name)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn("success")
	return
}

func (receiver *AccountController) DelDyCollectTag() {
	id, err := receiver.GetInt(":id")
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	comErr := business.NewUserBusiness().DelDyCollectTag(id)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn("success")
	return
}

//标签列表
func (receiver *AccountController) DyCollectLabel() {
	userBusiness := business.NewUserBusiness()
	collectType, err := receiver.GetInt("collect_type", 1)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	collectLabel, comErr := userBusiness.GetDyCollectLabel(receiver.UserId, collectType)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	var ret []string
	for _, v := range collectLabel {
		ret = append(ret, strings.Split(v, ",")...)
	}
	receiver.SuccReturn(utils.UniqueStringSlice(ret))
	return
}

//收藏达人备注
func (receiver *AccountController) DyCollectRemark() {
	InputData := receiver.InputFormat()
	id := InputData.GetInt("id", 0)
	var collect dcm.DcUserDyCollect
	_, err := dcm.Get(id, &collect)
	if err != nil || collect.Id == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	remark := InputData.GetString("remark", "")
	if _, err := dcm.UpdateInfo(dcm.GetDbSession(), id, map[string]interface{}{"remark": remark}, new(dcm.DcUserDyCollect)); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(nil)
	return
}
