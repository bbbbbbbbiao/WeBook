package async

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms/ratelimit"
	"sync/atomic"
)

/**
 * @author: biao
 * @date: 2026/2/6 下午3:38
 * @description: 装饰器-将请求同步转异步
 */

var ErrLimited = ratelimit.ErrLimited

type AsyncSmsService struct {
	smsSvc sms.Service

	errRate float64 // 错误率阈值

	reqCnt uint64 // 请求数
	errCnt uint64 // 错误数
}

func NewAsyncSmsService(smsSvc sms.Service, errRate float64) sms.Service {
	return &AsyncSmsService{
		smsSvc:  smsSvc,
		errRate: errRate,
	}
}

func (a *AsyncSmsService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {

	atomic.AddUint64(&a.reqCnt, 1)
	err := a.smsSvc.Send(ctx, tpl, args, numbers...)

	switch err {
	case nil:
		return nil
	case ErrLimited:
		// 触发了限流
		return err
	default:
		atomic.AddUint64(&a.errCnt, 1)

		errRate := a.errRate
		reqCnt := atomic.LoadUint64(&a.reqCnt)
		errCnt := atomic.LoadUint64(&a.errCnt)
		actualErrRate := float64(errCnt) / float64(reqCnt)

		// 错误率大于了 阈值 -> 系统崩了
		if actualErrRate > errRate {

		}

		return err
	}
	// TODO: 同步转异步后，还需要返回err错误码吗
}
