package failover

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"sync/atomic"
)

/**
 * @author: biao
 * @date: 2026/2/6 下午3:44
 * @description: 装饰器-超过对应错误率，自动切换
 */

type ErrorRateFailover struct {
	smsSvcs []sms.Service
	idx     uint64

	errorRate float64 // 错误率阈值

	reqCnt uint64 // 请求数
	errCnt uint64 // 错误数

}

func NewErrorRateFailover(smsSvcs []sms.Service, errorRate float64) sms.Service {
	return &ErrorRateFailover{
		smsSvcs:   smsSvcs,
		errorRate: errorRate,
	}
}

func (e *ErrorRateFailover) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadUint64(&e.idx)
	length := uint64(len(e.smsSvcs))
	errRate := e.errorRate
	reqCnt := atomic.LoadUint64(&e.reqCnt)
	errCnt := atomic.LoadUint64(&e.errCnt)

	// 计算实际错误率
	var actualErrRate float64
	if errCnt != 0 {
		actualErrRate = float64(errCnt) / float64(reqCnt)
	} else {
		actualErrRate = 0.0
	}

	// 错误率高于阈值
	if actualErrRate > errRate {
		newIdx := (idx + 1) % length
		if atomic.CompareAndSwapUint64(&e.idx, idx, newIdx) {
			atomic.StoreUint64(&e.reqCnt, 0)
			atomic.StoreUint64(&e.errCnt, 0)
		}
		idx = atomic.LoadUint64(&e.idx)
	}

	atomic.AddUint64(&e.reqCnt, 1)
	err := e.smsSvcs[idx%length].Send(ctx, tpl, args, numbers...)

	switch err {
	case nil:
		return nil
	default:
		atomic.AddUint64(&e.errCnt, 1)
		return err
	}
}
