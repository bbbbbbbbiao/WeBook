package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

/**
 * @author: biao
 * @date: 2026/2/1 下午2:57
 * @description: 基于滑动窗口的限流
 */

//go:embed slide_window.lua
var luaSlideWindow string

type RedisSlidingWindowLimiter struct {
	Cmd redis.Cmdable

	// 窗口大小（时间跨度大小）
	Interval time.Duration
	// 阈值
	Rate int

	// Interval 时间内 最多允许 Rate 个请求
	// 1s 时间内 最多允许 3000 个请求
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &RedisSlidingWindowLimiter{
		Cmd:      cmd,
		Interval: interval,
		Rate:     rate,
	}
}

func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return false, fmt.Errorf("generate uuid falied %w", err)
	}

	return r.Cmd.Eval(ctx, luaSlideWindow, []string{key}, r.Interval.Milliseconds(), r.Rate, time.Now().UnixMilli(), uid.String()).Bool()
}
