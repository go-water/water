package ratelimit

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-water/water/endpoint"
	"golang.org/x/time/rate"
)

// IPBasedLimiter 基于客户端IP的限流器
type IPBasedLimiter struct {
	limiters      sync.Map // map[string]*rate.Limiter (IP -> Limiter)
	interval      time.Duration
	burst         int
	mu            sync.RWMutex
	lastClean     time.Time
	cleanInterval time.Duration
}

// NewIPBasedLimiter 创建基于IP的限流器
// interval 是限流的间隔时间，burst 是突发大小
func NewIPBasedLimiter(interval time.Duration, burst int) *IPBasedLimiter {
	return &IPBasedLimiter{
		interval:      interval,
		burst:         burst,
		cleanInterval: 10 * time.Minute, // 每10分钟清理一次
	}
}

// getLimiter 获取或创建指定IP的限流器
func (ibl *IPBasedLimiter) getLimiter(ip string) *rate.Limiter {
	limiterAny, _ := ibl.limiters.LoadOrStore(ip, rate.NewLimiter(rate.Every(ibl.interval), ibl.burst))
	return limiterAny.(*rate.Limiter)
}

// IPErrorLimiter 返回基于IP的错误限流中间件
func (ibl *IPBasedLimiter) IPErrorLimiter(getIP func(ctx context.Context) string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (any, error) {
			ip := getIP(ctx)
			if ip == "" {
				// 如果无法获取IP，拒绝请求
				return nil, errors.New("unable to get client IP")
			}

			limiter := ibl.getLimiter(ip)
			if !limiter.Allow() {
				return nil, errors.New("rate limit exceeded for IP: " + ip)
			}

			return next(ctx, request)
		}
	}
}

// IPDelayingLimiter 返回基于IP的延迟限流中间件
func (ibl *IPBasedLimiter) IPDelayingLimiter(getIP func(ctx context.Context) string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (any, error) {
			ip := getIP(ctx)
			if ip == "" {
				return nil, errors.New("unable to get client IP")
			}

			limiter := ibl.getLimiter(ip)
			if err := limiter.Wait(ctx); err != nil {
				return nil, err
			}

			return next(ctx, request)
		}
	}
}
