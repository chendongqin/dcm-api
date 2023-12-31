package main

import (
	"dongchamao/business"
	"dongchamao/command"
	"dongchamao/global"
	"github.com/json-iterator/go/extra"
	"github.com/urfave/cli"
	"os"
)

func init() {
	extra.RegisterFuzzyDecoders()
}

func main() {
	cliApp := cli.NewApp()
	cliApp.Commands = getCommands()
	cliApp.Run(os.Args)
}

func getCommands() []cli.Command {
	command := cli.Command{
		Name:   "start",
		Usage:  "run command",
		Action: runCMD,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "ac",
				Value: "",
				Usage: "执行内容（user_help）",
			},
			cli.StringFlag{
				Name:  "env,e",
				Value: "prod",
				Usage: "runtime environment, dev|test|prod",
			},
			cli.IntFlag{
				Name:  "add",
				Value: 0,
				Usage: "extra add",
			},
		},
	}

	return []cli.Command{command}
}

func runCMD(ctx *cli.Context) {
	global.InitEnv()
	if !ctx.IsSet("ac") {
		panic("m, ac为必填项")
	}
	switch ctx.String("ac") {
	case "live_room_monitor":
		command.LiveRoomMonitor()
	case "live_monitor":
		command.LiveMonitor()
	case "update_live_monitor_status":
		command.UpdateLiveMonitorStatus()
	case "amount_expire_wechat_notice":
		command.AmountExpireWechatNotice()
	case "check_rank":
		//监控-每隔小时的30分监控（目前是除了商品榜和小时榜的榜单）
		command.CheckRank()
	case "check_rank_hour":
		//监控-小时榜
		command.CheckRankHour()
	case "check_rank_goods":
		//监控-每隔小时的10分监控（目前只有商品榜）
		command.CheckGoodsRank()
	case "ali_log_ana":
		//阿里云日志分析
		safe := business.NewSafeBusiness()
		safe.CommonAnalyseLogs()
	default:
		panic("undefined ac")
	}
	os.Exit(0)
}
