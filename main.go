package main

import (
	"dongchamao/global"
	_ "dongchamao/routers"
	"github.com/astaxie/beego"
)

func main() {
	global.InitEnv()
	beego.Run()
}
