package ioc

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms/memory"
	r "github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms/ratelimit"
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

/**
 * @author: biao
 * @date: 2026/1/23 下午9:50
 * @description:
 */

func InitSMSService(cmd redis.Cmdable) sms.Service {
	service := memory.NewService()
	limiter := ratelimit.NewRedisSlidingWindowLimiter(cmd, time.Second, 1000)
	limitService := r.NewRateLimitService(service, limiter)
	return limitService
}
