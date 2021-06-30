钉钉易玩

#### git commit 提交规范

> + type: commit 的类型
> + feat: 新特性
> + fix: 修改Bug问题
> + perf: 优化相关问题
> + refactor: 代码重构
> + docs: 文档修改
> + test: 测试用例修改
> + modify: 调整代码
> + footer: 一些备注, 通常是 BREAKING CHANGE 或修复的 bug 的链接.

#### 代码注释规范
```
//常量一样用小驼峰标准，首字母大写
const appEnv = "env"
```
函数方法 注释
```
//@Title 用户退出登录
//@Success 1 {map[string]interface{}}
//@Failure 0|5000 {map[string]interface{}}
//@router /v1/login/logout [get]
```

参数规范
GET POST 参数统一下划线风格
```
nick_name  昵称
phone_number 手机号
```
Redis Key 规范 
```
dp 统一前缀  
dp:code:login:%s  code 业务类型  login 业务细分  %s 具体参数
dp:limit:action:task_push
```

