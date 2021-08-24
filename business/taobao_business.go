package business

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"github.com/astaxie/beego/logs"
	"github.com/bitly/go-simplejson"
	"github.com/nilorg/go-opentaobao"
	"strings"
	"time"
)

/*调用淘宝开放平台接口的逻辑封装*/

type TaoBaoBusiness struct {
	AppKey    string
	AppSecret string
}

func NewTaoBaoBusiness() *TaoBaoBusiness {
	logic := new(TaoBaoBusiness)
	logic.init()
	return logic
}

func (this *TaoBaoBusiness) init() {
	this.AppKey = global.Cfg.String("tb_app_key")
	this.AppSecret = global.Cfg.String("tb_app_secret")
	opentaobao.AppKey = global.Cfg.String("tb_app_key")
	opentaobao.AppSecret = global.Cfg.String("tb_app_secret")
	opentaobao.Router = global.Cfg.String("tb_api_url")
}

func (this *TaoBaoBusiness) ChangeAppSecret(key string, secret string) {
	this.AppKey = key
	this.AppSecret = secret
	opentaobao.AppKey = key
	opentaobao.AppSecret = secret
}

// 获取商品详情信息（简版）
func (this *TaoBaoBusiness) getItemInfo(itemId string) (res *simplejson.Json, err error) {
	res, err = opentaobao.Execute("taobao.tbk.item.info.get", opentaobao.Parameter{
		"num_iids": itemId,
	})
	return
}

func (this *TaoBaoBusiness) GetItemInfo(itemId string) (map[string]interface{}, error) {
	item, err := this.getItemInfo(itemId)
	if err != nil {
		logs.Error("商品库查询失败:", err)
		return nil, err
	}
	list, err := item.Get("tbk_item_info_get_response").Get("results").Get("n_tbk_item").Array()
	if err != nil {
		logs.Error("商品库解析失败:", err)
		return nil, err
	}
	if len(list) <= 0 {
		return nil, nil
	}
	product, _ := list[0].(map[string]interface{})

	return product, nil
}

// 是否是内容物料
func (this *TaoBaoBusiness) IsContentMaterial(itemId string) (bool, error) {
	item, err := this.getItemInfo(itemId)
	if err != nil {
		logs.Error("商品库查询失败:", err)
		return false, err
	}
	list, err := item.Get("tbk_item_info_get_response").Get("results").Get("n_tbk_item").Array()
	if err != nil {
		logs.Error("商品库解析失败:", err)
		return false, err
	}
	if len(list) <= 0 {
		return false, nil
	}
	product, ok := list[0].(map[string]interface{})
	if !ok {
		return false, nil
	}
	materialLibType := utils.ToString(product["material_lib_type"])
	types := strings.Split(materialLibType, ",")
	if utils.InArray("2", types) {
		return true, nil
	}
	return false, nil
}

// 搜索淘宝物料
func (this *TaoBaoBusiness) SearchMaterial(category string, materialId int, pageNo int, pageSize int) (res *simplejson.Json, err error) {
	adzoneId := 109456050192
	res, err = opentaobao.Execute("taobao.tbk.dg.material.optional", opentaobao.Parameter{
		"page_size":   pageSize,
		"page_no":     pageNo,
		"q":           category,
		"adzone_id":   adzoneId,
		"material_id": materialId,
	})
	return
}

func (this *TaoBaoBusiness) GetToken(code string) (*simplejson.Json, error) {
	opentaobao.Router = global.Cfg.String("tb_token_url")
	res, err := opentaobao.Execute("taobao.top.auth.token.create", opentaobao.Parameter{
		"client_id":     this.AppKey,
		"client_secret": this.AppSecret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  global.Cfg.String("tb_oauth_redirect_url") + "/v1/boss/auth/confirm",
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *TaoBaoBusiness) GetOrder(token string, startTime time.Time, endTime time.Time, positionIndex string, pageSize int) (response []interface{}, hasNext bool, returnPositionIndex string, err error) {
	startTimeStr := startTime.Format("2006-01-02 15:04:05")
	endTimeStr := endTime.Format("2006-01-02 15:04:05")

	duration, _ := time.ParseDuration("1m")
	opentaobao.Timeout = duration
	opentaobao.Session = token
	res, err := opentaobao.Execute("taobao.tbk.sc.order.details.get", opentaobao.Parameter{
		"query_type":     1,
		"start_time":     startTimeStr,
		"end_time":       endTimeStr,
		"position_index": positionIndex,
		"page_size":      pageSize,
	})
	if err != nil {
		return
	}
	if res == nil {
		return
	}

	response, _ = res.Get("tbk_sc_order_details_get_response").Get("data").Get("results").
		Get("publisher_order_dto").Array()
	hasNext, _ = res.Get("tbk_sc_order_details_get_response").Get("data").Get("has_next").Bool()
	returnPositionIndex, _ = res.Get("tbk_sc_order_details_get_response").Get("data").Get("position_index").String()

	return
}

func (this *TaoBaoBusiness) GetPrivilege(itemId string) (map[string]interface{}, error) {
	res, err := opentaobao.Execute("taobao.tbk.privilege.get", opentaobao.Parameter{
		"item_id":   itemId,
		"adzone_id": "109455450332",
		"site_id":   "856050117",
	})
	if err != nil {
		logs.Error("商品库查询失败:", err)
		return nil, err
	}
	println(res)

	return nil, nil

}

//淘宝客-服务商-淘口令解析&转链
func (this *TaoBaoBusiness) TpwdConvert(pwc string) (productId string, url string) {
	//tbOauth, _ := douyinmodelsV2.NewSvTbOauthModel().GetByTbUserId("2201435067094")
	//sessionKey := tbOauth.AccessToken
	sessionKey := ""

	opentaobao.Timeout = time.Second
	opentaobao.Session = sessionKey
	res, err := opentaobao.Execute("taobao.tbk.sc.tpwd.convert", opentaobao.Parameter{
		"password_content": pwc,
		"adzone_id":        "109753600219",
		"site_id":          "856050117",
	})
	if err != nil {
		//todo 写入抓取接口失败的日志
		//println(err.Error())
		return
	}
	if res == nil {
		return
	}

	response := res.Get("tbk_sc_tpwd_convert_response").Get("data")
	url, _ = response.Get("click_url").String()
	productId, _ = response.Get("num_iid").String()

	return
}
