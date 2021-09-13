package v1

import (
	"dongchamao/business"
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/models/dcm"
	"dongchamao/models/repost"
	"strconv"
)

type VipController struct {
	controllers.ApiBaseController
}

//获取抖音团队列表
func (receiver *VipController) GetDyTeam() {
	var userVip dcm.DcUserVip
	if _, err := dcm.GetDbSession().Where("user_id=? and platform =?", receiver.UserId, 1).Get(&userVip); err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	page := receiver.GetPage("page")
	pageSize := receiver.GetPageSize("page", 10, 100)
	list, total, err := business.NewVipBusiness().GetDyTeam(userVip.Id, page, pageSize)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	var subUserIds []string
	for _, v := range list {
		subUserIds = append(subUserIds, strconv.Itoa(v.UserId))
	}
	userInfo, comErr := business.NewUserBusiness().GetUserList(subUserIds)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	var userInfoMap = make(map[int]dcm.DcUser, len(userInfo))
	for _, v := range userInfo {
		userInfoMap[v.Id] = v
	}
	var ret = make([]repost.DyTeamSubRet, 0, len(userInfoMap))
	for _, v := range list {
		ret = append(ret, repost.DyTeamSubRet{
			UserVipId:  v.Id,
			Remark:     v.Remark,
			UpdateTime: v.UpdateTime,
			Id:         userInfoMap[v.UserId].Id,
			Username:   userInfoMap[v.UserId].Username,
			LoginTime:  userInfoMap[v.UserId].LoginTime,
		})
	}
	receiver.SuccReturn(map[string]interface{}{"list": ret, "page": page, "pageSize": pageSize, "total": total})
}

//添加抖音子账号
func (receiver *VipController) AddDyTeamSub() {
	if receiver.DyLevel == 0 {
		receiver.FailReturn(global.NewMsgError("非专业版会员无法添加"))
	}
	inputData := receiver.InputFormat()
	username := inputData.GetString("username", "")
	if username == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var subUser dcm.DcUser
	_, err := dcm.GetBy("username", username, &subUser)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	if subUser.Id == 0 {
		receiver.FailReturn(global.NewMsgError("用户不存在"))
		return
	}
	comErr := business.NewVipBusiness().AddDyTeamSub(receiver.UserId, subUser.Id)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(nil)
}

//移除抖音子账号
func (receiver *VipController) RemoveDyTeam() {
	inputData := receiver.InputFormat()
	subUserId := inputData.GetInt("user_vip_id", 0)
	if subUserId == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	err := business.NewVipBusiness().RemoveDyTeamSub(subUserId)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(nil)
}

//抖音子账号备注
func (receiver *VipController) AddDySubRemark() {
	if receiver.DyLevel == 0 {
		receiver.FailReturn(global.NewMsgError("非专业版会员无法添加"))
	}
	inputData := receiver.InputFormat()
	userVipId := inputData.GetInt("user_vip_id", 0)
	if userVipId == 0 {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	remark := inputData.GetString("remark", "")
	if remark == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	var subUser dcm.DcUserVip
	_, err := dcm.UpdateInfo(nil, userVipId, map[string]interface{}{"remark": remark}, subUser)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.SuccReturn(nil)
}
