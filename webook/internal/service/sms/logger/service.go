package logger

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"go.uber.org/zap"
)

/**
 * @author: biao
 * @date: 2026/3/14 下午5:32
 * @description:
 */

type LoggerService struct {
	SmsSvc sms.Service
}

func NewLoggerService(smsSvc sms.Service) sms.Service {
	return &LoggerService{
		SmsSvc: smsSvc,
	}
}

// 当调用第三方需要查看发送和返回参数时，可以使用装饰器模式
func (l *LoggerService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	zap.L().Info("发送短信",
		zap.String("tpl", tpl),
		zap.Any("args", args))

	err := l.SmsSvc.Send(ctx, tpl, args, numbers...)
	if err != nil {
		zap.L().Error("发送短信失败", zap.Error(err))
		return err
	}
	return nil
}
