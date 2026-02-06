package failover

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"sync/atomic"
)

/**
 * @author: biao
 * @date: 2026/2/4 下午4:08
 * @description: 连续超时n次自动切换供应商
 */

type TimeOutFailoverService struct {
	smsSvcs []sms.Service
	idx     uint64
	// 超时次数
	cnt uint64
	// 阈值
	threshold uint64
}

func NewTimeOutFailoverService(smsSvcs []sms.Service, threshold uint64) sms.Service {
	return &TimeOutFailoverService{
		smsSvcs:   smsSvcs,
		threshold: threshold,
	}
}

// Send 连续超时n次自动切换
func (t *TimeOutFailoverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadUint64(&t.idx)
	cnt := atomic.LoadUint64(&t.cnt)
	length := uint64(len(t.smsSvcs))

	if cnt > t.threshold {
		newIdx := (idx + 1) % length
		if atomic.CompareAndSwapUint64(&t.idx, idx, newIdx) {
			atomic.StoreUint64(&t.cnt, 0)
		}
		// 为啥还要重新取值呢
		// 因为会有并发问题，导致短时间idx加了很多，所以再取一次
		idx = atomic.LoadUint64(&t.idx)
	}

	svc := t.smsSvcs[idx%length]
	err := svc.Send(ctx, tpl, args, numbers...)

	switch err {
	case nil:
		atomic.StoreUint64(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		atomic.AddUint64(&t.cnt, 1)
		return err
	}
	return err
}
