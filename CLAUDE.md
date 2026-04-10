# CLAUDE.md

本文档为 Claude Code (claude.ai/code) 在本代码库中工作提供指导。

## 概述

`go-water` 是一个专注于业务隔离和系统设计的 Go Web 框架。它提供了一种结构化的方法来构建 Web 应用程序，清晰地分离了控制层（处理器）、服务层（业务逻辑）和中间件。该框架内置支持路由、上下文管理、请求绑定、验证、JWT 认证、熔断、限流、结构化日志记录和模板渲染。

## 架构

### 核心组件

- **Water**: 主框架实例 (`water.Water`)，管理路由器、中间件和上下文池。
- **Router**: 处理 HTTP 路由注册，支持分组和静态文件。路由存储在按方法和路径键控的映射中。
- **Context** (`water.Context`): 包装 `http.Request` 和 `http.ResponseWriter`，提供绑定、查询参数、表单数据、JSON 响应等辅助方法。使用上下文池以提高性能。
- **Handler**: 两种类型的处理器：
  - `HandlerFunc`: 标准函数签名 `func(*water.Context)`，用于简单的 HTTP 处理器。
  - `water.Handler` 接口：用于基于服务的端点，内置中间件（限流、熔断、日志记录）。
- **Service**: 业务逻辑层。服务实现 `water.Service` 接口（`Name`, `SetLogger`）并嵌入 `water.ServerBase`。服务必须有一个 `Handle` 方法，签名为 `func(context.Context, *RequestType) (*ResponseType, error)`，其中 `RequestType` 和 `ResponseType` 是用户定义的结构体。该方法通过反射调用。
- **Endpoint**: 一个 `func(ctx context.Context, req any) (any, error)` 抽象，包装服务调用。中间件 (`endpoint.Middleware`) 可以应用于端点。
- **Binding**: 请求数据绑定，支持 JSON、表单、查询参数、URI、请求头和自定义 set 绑定。使用 `github.com/go-playground/validator/v10` 进行验证。
- **Middleware**:
  - HTTP 中间件：`func(http.Handler) http.Handler`，全局应用或按路由器应用。
  - 端点中间件：`endpoint.Middleware`，应用于服务端点（限流、熔断）。
- **Server Options**: 处理器配置 (`water.ServerOption`)，包括过滤器、终结器、错误限流器、延迟限流器和熔断器。
- **Filter & Finalizer**: 服务端点的前置和后置处理钩子 (`water.Filter`, `water.FinalizerFunc`)。
- **Circuit Breaking**: 通过 `circuitbreaker.GoBreaker` 集成 `github.com/sony/gobreaker`。
- **Rate Limiting**: 两种类型：错误限流器（超过限制时拒绝请求）和延迟限流器（等待）。使用 `golang.org/x/time/rate`。
- **Logging**: 通过 `log/slog` 进行结构化 JSON 日志记录。可通过 `logger` 包配置日志记录器。
- **Authentication**: 基于 JWT 的认证，使用 RSA 密钥 (`water.SetAuthToken`, `water.ParseFromRequest`)。
- **Rendering**: 通过 `render` 包支持 JSON、HTML、文本和重定向响应。
- **Templates**: 通过 `multitemplate` 包支持多模板。

### 请求流程

1. HTTP 请求进入 `Water.Run()` → `RouterHandler.ServeHTTP()`。
2. 从池中分配上下文，分配请求/响应。
3. 调用路由匹配的 `HandlerFunc`（简单处理器）或 `water.Handler.ServerWater`（服务端点）。
4. 对于服务端点：
   - 执行过滤器和终结器。
   - 应用端点中间件（限流、熔断）。
   - 通过反射调用服务的 `Handle` 方法。
   - 返回响应，记录错误。
5. 重置上下文并返回到池中。

### 关键设计模式

- **关注点分离**: 处理器定义 HTTP 接口，服务包含业务逻辑。
- **中间件链**: HTTP 中间件和端点中间件都允许横切关注点。
- **上下文池化**: 减少分配开销。
- **基于接口的扩展**: 绑定、验证、日志记录和渲染都基于接口，支持自定义实现。

## 开发命令

这是一个 Go 模块。使用标准 Go 工具：

```bash
# 安装依赖
go mod tidy

# 运行测试（当前没有测试）
go test ./...

# 格式化代码
go fmt ./...

# 代码检查
go vet ./...

# 构建包
go build ./...

# 生成文档
go doc ./...

# 运行示例（参见 examples/README.md 获取外部示例）
cd examples/some_example && go run main.go
```

## 项目结构

```
.
├── water.go                 # Water 主结构和路由器
├── context.go               # Context 实现
├── handler.go               # 处理器和服务端点逻辑
├── service.go               # 服务接口和 ServerBase
├── auth.go                  # JWT 认证工具
├── option.go                # 处理器的服务器选项
├── middleware.go            # HTTP 中间件类型
├── path.go                  # 路径工具
├── fs.go                    # 文件系统工具
├── errors.go                # 错误类型
├── logger.go                # 日志记录器设置
├── version.go               # 版本常量
├── go.mod                   # 模块定义
├── README.md                # 项目概述和示例
├── binding/                 # 请求绑定实现
├── circuitbreaker/          # 熔断器中间件
├── endpoint/                # 端点类型和中间件
├── logger/                  # 结构化日志记录配置
├── multitemplate/           # 多模板渲染
├── ratelimit/               # 限流中间件
├── render/                  # 响应渲染接口
└── examples/                # 外部示例参考
```

## 注意事项

- 框架使用 Go 1.22+ 原生路由 (`net/http` `PathValue`)。
- 通过 `github.com/go-playground/validator/v10` 集成验证。验证错误以 `validator.ValidationErrors` 形式返回。
- JWT 认证使用 RSA 密钥（私钥/公钥文件）。
- `Set` 绑定允许绑定通过 `c.Set()` 设置的上下文值。
- `ServerBase` 结构体为服务提供默认的 `Name` 和日志记录器方法。
- 中间件顺序：全局 HTTP 中间件按注册顺序应用；路由器中间件按相反顺序应用（最后注册的先运行）。