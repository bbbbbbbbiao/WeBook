package retryable

import (
	"context"
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
)

/**
 * @author: biao
 * @date: 2026/2/5 下午4:12
 * @description: 装饰器：重试策略
 */

type RetryService struct {
	smsSvc sms.Service

	// 重试次数
	retryMax int
}

func NewRetryService(smsSvc sms.Service, retryMax int) sms.Service {
	return &RetryService{
		smsSvc:   smsSvc,
		retryMax: retryMax,
	}
}

func (r *RetryService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for i := 0; i < r.retryMax; i++ {
		err := r.smsSvc.Send(ctx, tpl, args, numbers...)

		if err == nil {
			return nil
		}
	}

	return errors.New("重试也失败了")
}
