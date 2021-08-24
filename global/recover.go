package global

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"runtime"
	"strings"
)

func RecoverPanic(fixer ...func(err interface{})) {
	if err := recover(); err != nil {
		stacks := GetStacks()
		logs.Critical(fmt.Sprintf("err: %s, stacks: %s", err, stacks))
		if len(fixer) > 0 {
			fixer[0](err)
		}
	}
}

func RequestRecoverPanic(ctx *context.Context) {
	if err := recover(); err != nil {
		if err == beego.ErrAbort {
			return
		}
		if !beego.BConfig.RecoverPanic {
			panic(err)
		}
		stacks := GetStacks()
		logs.Critical(fmt.Sprintf("request url: %s, err: %s,  stack: %s", ctx.Input.URL(), err, stacks))
		if ctx.Output.Status != 0 {
			ctx.ResponseWriter.WriteHeader(ctx.Output.Status)
		} else {
			ctx.ResponseWriter.WriteHeader(500)
		}
		response := map[string]interface{}{
			"errCode": 5000,
			"errMsg":  "系统错误",
		}
		jb, _ := json.Marshal(response)
		ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = ctx.ResponseWriter.Write(jb)
	}
}

func GetStacks() string {
	var stack []string
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stack = append(stack, fmt.Sprintf("%s:%d", file, line))
	}
	joinStr := ", "
	if beego.BConfig.RunMode == beego.DEV { // 测试环境换行显示
		joinStr = "\r\n"
	}
	return strings.Join(stack, joinStr)
}
