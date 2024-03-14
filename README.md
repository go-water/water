简介：go-water 是一款设计层面的 web 框架（类似 gin，iris，beego，echo），更好的业务隔离，更好的系统设计，通过一系列接口、规范、约定、中间件，深度解耦业务系统。

### 星星增长趋势
[![Stargazers over time](https://starchart.cc/go-water/water.svg)](https://starchart.cc/go-water/water)

### 安装
```
go get -u github.com/go-water/water
```

### 技术概览
+ 支持原生路由(1.22)
+ slog 日志
+ 中间件
+ 多模板支持
+ rsa 加密，openssl 生成公/私钥对
+ jwt 登陆认证
+ pool 管理请求参数
+ option 配置修改
+ rate limit（限流）
+ circuit breaker（熔断）

用例
```go
package main

import (
	"github.com/go-water/water"
)

func main() {
	r := water.New()

	r.GET("/", func(c *water.Context) {
		c.Text(200, "Hello, World!")
	})

	r.Run(":8080")
}
```
在浏览器输入
```
http://localhost:8080/
```

### 样例仓库
+ [https://github.com/go-water/go-water](https://github.com/go-water/go-water)

### 官方文档
+ [https://go-water.cn](https://go-water.cn)

### 参考仓库
+ [kit](https://github.com/go-kit/kit)
+ [gin](https://github.com/gin-gonic/gin)
+ [grape](https://github.com/hossein1376/grape)