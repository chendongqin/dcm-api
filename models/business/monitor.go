package business

import (
	"dongchamao/global"
	"dongchamao/services/dingding"
	"fmt"
	"strings"
	"time"
)

type MonitorBusiness struct {
}

func NewMonitorBusiness() *MonitorBusiness {
	return new(MonitorBusiness)
}

const MonitorPrefix = "洞察猫"

//艾特电话名单
const AtMobilesGroupGeneral = "18605971553,13665927819,18659328891,13959201421,13386936061"

func (receiver *MonitorBusiness) SendTemplateMessage(level string, event string, content ...string) {
	color := "#FF0000"
	finalContent := strings.Trim(strings.Join(content, ""), " ")
	if finalContent == "" {
		return
	}
	str := "<font color=\"" + color + "\">**监控等级：**" + level + "</font>\n\n"
	str += "**时间：**" + time.Now().Local().Format("2006-01-02 15:04:05") + "\n\n"
	str += "**类型：**" + event + "\n\n"
	str += "**说明：**" + strings.Join(content, "")
	fmt.Println(str)
	receiver.SendMarkDown("监控数据提醒", str)
}

func (receiver *MonitorBusiness) SendMarkDown(title string, content string, atss ...string) {
	ding := dingding.NewWithTokenUrl(global.Cfg.String("ding_ding_monitor"))
	ats := AtMobilesGroupGeneral
	if len(atss) > 0 {
		ats = atss[0]
	}
	if ats != "" {
		atMobiles := strings.Split(ats, ",")
		tmp := ""
		for _, atMobile := range atMobiles {
			tmp += " @" + atMobile + " "
		}
		ding.SetAt(atMobiles...)
		content = tmp + "\n\n" + content
	}
	_ = ding.SendMarkDown(MonitorPrefix+title, content)
}
