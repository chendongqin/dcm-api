package v1dy

import (
	controllers "dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/models/business"
)

type AuthorController struct {
	controllers.ApiBaseController
}

func (receiver *AuthorController) AuthorBaseData() {
	authorId := receiver.Ctx.Input.Param(":author_id")
	if authorId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	authorBusiness := business.NewAuthorBusiness()
	authorBase, comErr := authorBusiness.HbaseGetAuthor(authorId)
	if comErr != nil {
		receiver.FailReturn(comErr)
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"author_base": authorBase.Data,
	})
	return
}
