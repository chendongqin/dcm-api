package business

import (
	"dongchamao/entity"
	"dongchamao/global"
	"dongchamao/global/utils"
	"dongchamao/services/hbaseService"
	"dongchamao/services/hbaseService/hbasehelper"
	"math"
)

type AwemeBusiness struct {
}

func NewAwemeBusiness() *AwemeBusiness {
	return new(AwemeBusiness)
}

func (a *AwemeBusiness) HbaseGetAweme(awemeId string) (data entity.DyAwemeData, comErr global.CommonError) {
	query := hbasehelper.NewQuery()
	result, err := query.SetTable(hbaseService.HbaseDyAweme).GetByRowKey([]byte(awemeId))
	if err != nil {
		comErr = global.NewMsgError(err.Error())
		return
	}
	if result.Row == nil {
		comErr = global.NewError(4040)
		return
	}
	authorMap := hbaseService.HbaseFormat(result, entity.DyAwemeMap)
	aweme := &entity.DyAweme{}
	utils.MapToStruct(authorMap, aweme)
	aweme.AwemeID = aweme.Data.ID
	duration := math.Ceil(float64(aweme.Data.Duration) / 1000)
	aweme.Data.Duration = utils.ToInt(duration)
	aweme.Data.AwemeTitle = aweme.AwemeTitle
	data = aweme.Data
	return
}
