# 限流使用指南

go-water 框架支持三种限流方式：全局限流、基于IP的限流和基于用户的限流。

## 1. 全局限流（现有功能）

对整个处理器应用统一的限流策略。

```go
package main

import (
	"context"
	"time"
	"github.com/go-water/water"
)

// 定义请求和响应类型
type TestReq struct {
	Message string `json:"message"`
}

type TestResp struct {
	Message string `json:"message"`
}

// 定义服务
type TestService struct {
	water.ServerBase
}

func (s *TestService) Handle(ctx context.Context, req *TestReq) (*TestResp, error) {
	return &TestResp{Message: "Hello: " + req.Message}, nil
}

func main() {
	app := water.NewWater()
	
	// 全局错误限流：每秒最多5个请求，突发10个
	handler := water.NewHandler(
		&TestService{},
		water.ServerErrorLimiter(time.Second, 5),
	)
	
	app.POST("/test", handler.ServerWater)
	
	app.Run(":8080")
}
```

## 2. 基于IP的限流

为不同的客户端IP分别限流，防止单个IP的恶意请求。

```go
package main

import (
	"context"
	"net/http"
	"time"
	"github.com/go-water/water"
	"github.com/go-water/water/ratelimit"
)

type TestReq struct {
	Message string `json:"message"`
}

type TestResp struct {
	Message string `json:"message"`
}

type TestService struct {
	water.ServerBase
}

func (s *TestService) Handle(ctx context.Context, req *TestReq) (*TestResp, error) {
	return &TestResp{Message: "Hello: " + req.Message}, nil
}

func main() {
	app := water.NewWater()
	
	// 创建基于IP的限流器：每秒最多10个请求，突发20个
	ipLimiter := ratelimit.NewIPBasedLimiter(time.Second, 10)
	
	// 从Context中获取客户端IP的函数
	getIP := func(ctx context.Context) string {
		// 从Context中提取water.Context
		val := ctx.Value("_go-water/context-key")
		if waterCtx, ok := val.(*water.Context); ok {
			return waterCtx.ClientIP()
		}
		return ""
	}
	
	service := &TestService{}
	handler := water.NewHandler(service)
	
	// 应用IP限流中间件到服务的端点中间件
	handler.Middlewares(
		ipLimiter.IPErrorLimiter(getIP),
	)
	
	app.POST("/test", handler.ServerWater)
	
	app.Run(":8080")
}
```

## 3. 基于用户的限流

为不同的用户分别限流，常用于API配额管理。

```go
package main

import (
	"context"
	"time"
	"github.com/go-water/water"
	"github.com/go-water/water/ratelimit"
)

type TestReq struct {
	Message string `json:"message"`
}

type TestResp struct {
	Message string `json:"message"`
}

type TestService struct {
	water.ServerBase
}

func (s *TestService) Handle(ctx context.Context, req *TestReq) (*TestResp, error) {
	return &TestResp{Message: "Hello: " + req.Message}, nil
}

func main() {
	app := water.NewWater()
	
	// 创建基于用户的限流器：每秒最多100个请求，突发200个
	userLimiter := ratelimit.NewUserBasedLimiter(time.Second, 100)
	
	// 从Context中获取用户ID的函数
	// 通常从JWT token或Session中获取
	getUserID := func(ctx context.Context) string {
		val := ctx.Value("_go-water/context-key")
		if waterCtx, ok := val.(*water.Context); ok {
			// 从context Key中获取用户ID
			userID, _ := waterCtx.Get("userID")
			if uid, ok := userID.(string); ok {
				return uid
			}
		}
		return ""
	}
	
	service := &TestService{}
	handler := water.NewHandler(service)
	
	// 应用用户限流中间件
	handler.Middlewares(
		userLimiter.UserErrorLimiter(getUserID),
	)
	
	app.POST("/test", handler.ServerWater)
	
	app.Run(":8080")
}
```

## 4. 延迟限流模式

不是拒绝超过限流的请求，而是让它们等待，直到满足限流条件。

```go
package main

import (
	"context"
	"time"
	"github.com/go-water/water"
	"github.com/go-water/water/ratelimit"
)

type TestReq struct {
	Message string `json:"message"`
}

type TestResp struct {
	Message string `json:"message"`
}

type TestService struct {
	water.ServerBase
}

func (s *TestService) Handle(ctx context.Context, req *TestReq) (*TestResp, error) {
	return &TestResp{Message: "Hello: " + req.Message}, nil
}

func main() {
	app := water.NewWater()
	
	// 创建基于IP的延迟限流器
	ipLimiter := ratelimit.NewIPBasedLimiter(time.Second, 5)
	
	getIP := func(ctx context.Context) string {
		val := ctx.Value("_go-water/context-key")
		if waterCtx, ok := val.(*water.Context); ok {
			return waterCtx.ClientIP()
		}
		return ""
	}
	
	service := &TestService{}
	handler := water.NewHandler(service)
	
	// 使用延迟限流：请求会等待而不是被拒绝
	handler.Middlewares(
		ipLimiter.IPDelayingLimiter(getIP),
	)
	
	app.POST("/test", handler.ServerWater)
	
	app.Run(":8080")
}
```

## 5. 组合多种限流策略

可以组合全局限流、IP限流和用户限流。

```go
package main

import (
	"context"
	"time"
	"github.com/go-water/water"
	"github.com/go-water/water/ratelimit"
)

type TestReq struct {
	Message string `json:"message"`
}

type TestResp struct {
	Message string `json:"message"`
}

type TestService struct {
	water.ServerBase
}

func (s *TestService) Handle(ctx context.Context, req *TestReq) (*TestResp, error) {
	return &TestResp{Message: "Hello: " + req.Message}, nil
}

func main() {
	app := water.NewWater()
	
	ipLimiter := ratelimit.NewIPBasedLimiter(time.Second, 50)   // 每IP每秒50个
	userLimiter := ratelimit.NewUserBasedLimiter(time.Second, 100) // 每用户每秒100个
	
	getIP := func(ctx context.Context) string {
		val := ctx.Value("_go-water/context-key")
		if waterCtx, ok := val.(*water.Context); ok {
			return waterCtx.ClientIP()
		}
		return ""
	}
	
	getUserID := func(ctx context.Context) string {
		val := ctx.Value("_go-water/context-key")
		if waterCtx, ok := val.(*water.Context); ok {
			userID, _ := waterCtx.Get("userID")
			if uid, ok := userID.(string); ok {
				return uid
			}
		}
		return ""
	}
	
	service := &TestService{}
	handler := water.NewHandler(
		service,
		// 全局限流：整个服务最多200个并发请求
		water.ServerErrorLimiter(time.Second, 200),
	)
	
	// 组合IP限流和用户限流
	handler.Middlewares(
		ipLimiter.IPErrorLimiter(getIP),
		userLimiter.UserErrorLimiter(getUserID),
	)
	
	app.POST("/test", handler.ServerWater)
	
	app.Run(":8080")
}
```

## 工作原理

### IP 限流
- 为每个客户端IP维护独立的限流器
- 通过 `ClientIP()` 方法获取客户端IP（支持代理转发）
- 使用 `sync.Map` 实现高效的并发访问

### 用户限流
- 为每个已认证用户维护独立的限流器
- 通过自定义 `getUserID` 函数从Context中提取用户信息
- 适用于API Key或会话认证的场景

## 最佳实践

1. **选择合适的限流策略**
   - 全局限流：保护服务资源
   - IP限流：防止单个IP的恶意请求
   - 用户限流：实施API配额策略

2. **设置合理的限流参数**
   - `interval`: 限流周期（通常为秒级）
   - `burst`: 允许的突发请求数

3. **结合使用多种策略**
   - 不同限流方式可以同时应用
   - 限流中间件按顺序执行

4. **处理限流错误**
   - 返回 429 Too Many Requests HTTP状态码
   - 提供友好的错误提示

5. **监控限流指标**
   - 记录被限流的请求
   - 分析使用模式，调整限流参数
