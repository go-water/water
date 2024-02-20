简介：go-water 是一款设计层面的 web 框架（像 gin，iris，beego，echo 一样，追求卓越）。 我们使命：更好的业务隔离，更好的系统设计，通过一系列接口、规范、约定、中间件，深度解耦业务系统。

### 星星增长趋势
[![Stargazers over time](https://starchart.cc/go-water/water.svg)](https://starchart.cc/go-water/water)

### 安装
```
go get -u github.com/go-water/water
```

### 技术概览
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
```
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-water/water"
	"github.com/go-water/water/multitemplate"
)

func main() {
	router := water.New()
	router.HTMLRender = createMyRender()

	router.Use(Logger)
	router.GET("/", Index)
	v2 := router.Group("/v2")
	{
		v2.GET("/hello", GetHello)
	}

	router.Serve(":80")
}

func Index(ctx *water.Context) {
	ctx.HTML(http.StatusOK, "index", water.H{"title": "我是标题", "body": "你好，朋友。"})
}

func GetHello(ctx *water.Context) {
	ctx.JSON(http.StatusOK, water.H{"msg": "Hello World!"})
}

func Logger(handlerFunc water.HandlerFunc) water.HandlerFunc {
	return func(ctx *water.Context) {
		start := time.Now()
		defer func() {
			msg := fmt.Sprintf("[WATER] %v | %15s | %13v | %-7s %s",
				time.Now().Format("2006/01/02 - 15:04:05"),
				ctx.ClientIP(),
				time.Since(start),
				ctx.Request.Method,
				ctx.Request.URL.Path,
			)

			fmt.Println(msg)
		}()

		handlerFunc(ctx)
	}
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "views/layout.html", "views/index.html", "views/_header.html", "views/_footer.html")
	return r
}
```
views/layout.html
```
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>{{.title}}</title>
</head>
<body>
<div>
    <div>
        {{template "_header"}}
    </div>
    <div>
        {{template "content" .}}
    </div>
    <div>
        {{template "_footer"}}
    </div>
</div>
</body>
</html>
```
views/index.html
```
{{define "content"}}
我是内容：{{.body}}
{{end}}
```
views/_header.html
```
{{define "_header"}}
我是 Header。
{{end}}
```
views/_footer.html
```
{{define "_footer"}}
我是 Footer。
{{end}}
```

### 样例仓库
+ [https://github.com/go-water/go-water](https://github.com/go-water/go-water)

### 官方文档
+ [https://go-water.cn](https://go-water.cn)

### 参考仓库
+ [kit](https://github.com/go-kit/kit)
+ [gin](https://github.com/gin-gonic/gin)
+ [grape](https://github.com/hossein1376/grape)