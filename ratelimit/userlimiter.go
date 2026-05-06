package ratelimit

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-water/water/endpoint"
	"golang.org/x/time/rate"
)

// UserBasedLimiter 基于用户的限流器
type UserBasedLimiter struct {
	limiters      sync.Map // map[string]*rate.Limiter (UserID -> Limiter)
	interval      time.Duration
	burst         int
	mu            sync.RWMutex
	lastClean     time.Time
	cleanInterval time.Duration
}

// NewUserBasedLimiter 创建基于用户的限流器
func NewUserBasedLimiter(interval time.Duration, burst int) *UserBasedLimiter {
	return &UserBasedLimiter{
		interval:      interval,
		burst:         burst,
		cleanInterval: 10 * time.Minute,
	}
}

// getLimiter 获取或创建指定用户的限流器
func (ubl *UserBasedLimiter) getLimiter(userID string) *rate.Limiter {
	limiterAny, _ := ubl.limiters.LoadOrStore(userID, rate.NewLimiter(rate.Every(ubl.interval), ubl.burst))
	return limiterAny.(*rate.Limiter)
}

// UserErrorLimiter 返回基于用户的错误限流中间件
func (ubl *UserBasedLimiter) UserErrorLimiter(getUserID func(ctx context.Context) string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (any, error) {
			userID := getUserID(ctx)
			if userID == "" {
				// 如果无法获取用户ID，使用匿名用户处理或直接拒绝
				return nil, errors.New("user not authenticated")
			}

			limiter := ubl.getLimiter(userID)
			if !limiter.Allow() {
				return nil, errors.New("rate limit exceeded for user: " + userID)
			}

			return next(ctx, request)
		}
	}
}

// UserDelayingLimiter 返回基于用户的延迟限流中间件
func (ubl *UserBasedLimiter) UserDelayingLimiter(getUserID func(ctx context.Context) string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (any, error) {
			userID := getUserID(ctx)
			if userID == "" {
				return nil, errors.New("user not authenticated")
			}

			limiter := ubl.getLimiter(userID)
			if err := limiter.Wait(ctx); err != nil {
				return nil, err
			}

			return next(ctx, request)
		}
	}
}
