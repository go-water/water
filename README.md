简介：go-water 是一款设计层面的 web 框架（像 gin，iris，beego，echo 一样，追求卓越）。 我们使命：更好的业务隔离，更好的系统设计，通过一系列接口、规范、约定、中间件，深度解耦业务系统。

### 星星增长趋势
[![Stargazers over time](https://starchart.cc/go-water/water.svg)](https://starchart.cc/go-water/water)

### 安装
```
go get -u github.com/go-water/water
```

### 技术概览
+ zap 日志
+ rsa 加密，openssl 生成公/私钥对
+ jwt 登陆认证
+ errors 自定义处理
+ pool 管理请求参数
+ option 配置修改
+ rate limit（限流）
+ circuit breaker（熔断）

### 样例仓库
+ [https://github.com/go-water/go-water](https://github.com/go-water/go-water)

### 官方文档
+ [https://iissy.com/go-water](https://iissy.com/go-water)

### 注意
文档暂时还未来得及更新，请以源码为准。（v1.0.0预计4月1日发版，文档，样例将同步更新）