package test

import (
	"dongchamao/business"
	"fmt"
	"regexp"
	"testing"
)

func TestSpiderApi(t *testing.T) {
	spiderBusiness := business.NewSpiderBusiness()
	res := spiderBusiness.SpiderSpeedUp("author", "193965930847527")
	fmt.Println(res)
}

func TestSpiderAddLive(t *testing.T) {
	spiderBusiness := business.NewSpiderBusiness()
	spiderBusiness.AddLive("4195355415549012", 1962290, 4, 1633009421)
}

func TestGetAuthorByKeyword(t *testing.T) {
	spiderBusiness := business.NewSpiderBusiness()
	data := spiderBusiness.GetAuthorByKeyword("luoyonghao")
	fmt.Println(data)
}

func TestGetRoomPmt(t *testing.T) {
	spiderBusiness := business.NewSpiderBusiness()
	body := spiderBusiness.GetRoomPmt("73589350397")
	fmt.Println(body)
}

func TestParseUrl(t *testing.T) {
	pattern := `\/user\/(\d+)`
	reg := regexp.MustCompile(pattern)
	da := reg.FindAllStringSubmatch("/share/user/3202198174179895", -1)
	if len(da) > 0 {
		if len(da[0]) > 1 {
			ret := da[0][1]
			fmt.Println(ret)
		}
		return
	}
	return
}
