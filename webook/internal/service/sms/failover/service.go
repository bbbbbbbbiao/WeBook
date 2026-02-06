package failover

import (
	"context"
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"sync/atomic"
)

/**
 * @author: biao
 * @date: 2026/2/3 下午10:05
 * @description: 用装饰器实现自动切换短信供应商
 */

type FailOverService struct {
	smsSvcs []sms.Service
	idx     uint64 // 为啥这里不初始化？
	// 这里默认值去0，而且项目允许时只会初始化一次，所以他会一直往上加
}

func NewFailOverService(smsSvcs []sms.Service) sms.Service {
	return &FailOverService{
		smsSvcs: smsSvcs,
	}
}

func (f *FailOverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, smsSvc := range f.smsSvcs {
		err := smsSvc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		// 这里要记录一下哪些供应商没有成功
	}

	// 一般不可能全部供应商都失败，最有可能就是我们这边网络不好
	return errors.New("发送失败，所有供应商都尝试过了")
}

// SendV1 修改有二
// 1. 修改从下标为idx（取余）开始作为供应商
// 2. 细分了错误
func (f *FailOverService) SendV1(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.smsSvcs))
	for i := idx; i < length+idx; i++ {
		err := f.smsSvcs[i%length].Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			// 超时 或者 主动取消了
			return err
		}
	}
	// 一般不可能全部供应商都失败，最有可能就是我们这边网络不好
	return errors.New("发送失败，所有供应商都尝试过了")
}
