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

### 这个样例，复制代码，就可以直接跑
```
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-water/water"
	"github.com/sony/gobreaker"
)

func main() {
	router := water.New()
	router.GET("/", Index)
	_ = router.Run(":80")
}

// 控制层，这里定义了一个 Handlers 来管理所有业务接口
var (
	options = []water.ServerOption{
		// 一分钟内，连续10次后，将限流
		water.ServerErrorLimiter(time.Minute, 10),
		// 熔断定义，服务层异常将触发熔断
		water.ServerBreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{})),
	}

	index = water.NewHandler(&IndexService{ServerBase: &water.ServerBase{}}, options...)
)

func Index(ctx *water.Context) {
	request, err := water.BindJSON[IndexRequest](ctx)
	if err != nil {
		_ = ctx.JSON(http.StatusBadRequest, water.H{"err": err.Error()})
		return
	}

	request.Name = "Jimmy"
	resp, err := index.ServerWater(ctx, request)
	if err != nil {
		_ = ctx.JSON(http.StatusBadRequest, water.H{"err": err.Error()})
		return
	}

	_ = ctx.JSON(http.StatusOK, resp)
}

// 业务接口服务定义，结构体包含一个water.ServerBase，同时必须实现 Handle 方法
type IndexService struct {
	*water.ServerBase
}

type IndexRequest struct {
	Name string
}

type IndexResponse struct {
	Message string
}

func (s *IndexService) Handle(ctx context.Context, req *IndexRequest) (*IndexResponse, error) {
	resp := new(IndexResponse)
	resp.Message = fmt.Sprintf("Hello, %s!", req.Name)
	return resp, nil
	// 如果要测试服务熔断，可以打开下面代码，让代码返回异常，测试连续6次错误，第7次将不再进入这个方法
	// return nil, errors.New("service failure")
}
```
在浏览器输入
```
http://localhost/
```
运行结果
```
{
    "Message": "Hello, Jimmy!"
}
一分钟内连续10次后，限流
{
  "err": "rate limit exceeded"
}
熔断前，连续错误6次
{
  "err": "service failure"
}
熔断后
{
  "err": "circuit breaker is open"
}
```

### 样例仓库
+ [https://github.com/go-water/go-water](https://github.com/go-water/go-water)

### 使用的网站列表
+ [https://jitask.com](https://jitask.com)

### 参考仓库
+ [kit](https://github.com/go-kit/kit)
+ [gin](https://github.com/gin-gonic/gin)
+ [grape](https://github.com/hossein1376/grape)