package ratelimit

import (
	"context"
	"fmt"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/ginx/middlewares/ratelimit"
)

/**
 * @author: biao
 * @date: 2026/2/2 下午9:04
 * @description: 装饰器，给短信服务加上限流
 */

// 如果要查询是具体哪个错误时，则暴露出去
var ErrLimited = fmt.Errorf("触发了限流")

type RateLimitService struct {
	SmsSvc sms.Service
	limit  ratelimit.Limiter
}

func NewRateLimitService(smsSvc sms.Service, limit ratelimit.Limiter) sms.Service {
	return &RateLimitService{
		SmsSvc: smsSvc,
		limit:  limit,
	}
}

func (r *RateLimitService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limit, err := r.limit.Limit(ctx, "sms:tencent")
	if err != nil {
		// 系统错误
		// 可以限流：保守策略，第三方很坑的时候
		// 可以不限：尽量容错策略，第三方很强，这样可以保证我们业务可用性很高
		return fmt.Errorf("limit.Limit: %w", err)
	}

	if limit {
		return ErrLimited
	}
	err = r.SmsSvc.Send(ctx, tpl, args, numbers...)
	return err
}
