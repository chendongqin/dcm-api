package controllers

import (
	"dongchamao/business"
	"dongchamao/global"
	"github.com/astaxie/beego"
)

type BaseController struct {
	JsonEncrypt bool
	beego.Controller
}

func (this *BaseController) SuccReturn(retData interface{}) {
	retJson := map[string]interface{}{
		"code":   0,
		"status": true,
		"msg":    "ok",
		"data":   retData,
	}
	if retData == nil {
		retJson["data"] = []int{}
	}
	this.setData(retJson)
	this.ServeJSON()
	//this.Abort("200")
}

func (this *BaseController) FailReturn(err global.CommonError) {
	retJson := make(map[string]interface{})
	if err == nil {
		retJson["code"] = 5000
		retJson["msg"] = ""
	} else {
		retJson["code"], retJson["msg"] = err.Error()
	}
	this.Data["status"] = false
	this.setData(retJson)
	this.ServeJSON()
}
func (this *BaseController) FailReturnWithData(ErrorCode int, retData interface{}) {
	retJson := map[string]interface{}{
		"code":   ErrorCode,
		"msg":    "ok",
		"status": false,
		"data":   retData,
	}
	this.setData(retJson)
	this.ServeJSON()
}

func (this *BaseController) SuccReturnWithData(code int, retData interface{}) {
	retJson := map[string]interface{}{
		"code":   code,
		"status": true,
		"msg":    "ok",
		"data":   retData,
	}

	if retData == nil {
		retJson["data"] = []int{}
	}
	this.setData(retJson)
	this.ServeJSON()
}

//是否加密解密
func (this *BaseController) setData(retJson map[string]interface{}) {
	if this.JsonEncrypt {
		if data, exist := retJson["data"]; exist {
			retJson["data"] = business.JsonEncrypt(data)
		}
	}
	this.Data["json"] = retJson
}
