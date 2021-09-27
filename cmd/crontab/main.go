package main

import (
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
		command.CheckRank()
	default:
		panic("undefined ac")
	}
	os.Exit(0)
}
