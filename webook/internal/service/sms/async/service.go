package async

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms/ratelimit"
	"sync/atomic"
	"time"
)

/**
 * @author: biao
 * @date: 2026/2/6 下午3:38
 * @description: 装饰器-将请求同步转异步
 */

var ErrLimited = ratelimit.ErrLimited

type AsyncSmsService struct {
	reqRepo repository.ReqRepository
	smsSvc  sms.Service

	errRate float64 // 错误率阈值

	reqCnt uint64 // 请求数
	errCnt uint64 // 错误数

	maxRetry int // 重试次数
}

func NewAsyncSmsService(reqRepo repository.ReqRepository, smsSvc sms.Service, errRate float64, maxRetry int) sms.Service {
	return &AsyncSmsService{
		reqRepo:  reqRepo,
		smsSvc:   smsSvc,
		errRate:  errRate,
		maxRetry: maxRetry,
	}
}

func (a *AsyncSmsService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {

	atomic.AddUint64(&a.reqCnt, 1)
	err := a.smsSvc.Send(ctx, tpl, args, numbers...)

	switch err {
	case nil:
		return nil
	case ErrLimited:
		// 触发了限流（将请求存入数据库，用于异步发送）
		//TODO: 这里返回时，是否需要区分是数据库错误还是限流错误

		now := time.Now().UnixMilli()
		err = a.reqRepo.Create(ctx, domain.Request{
			Tpl:     tpl,
			Args:    args,
			Numbers: numbers,
			Ctime:   now,
			Utime:   now,
		})
		return ErrLimited
	default:
		atomic.AddUint64(&a.errCnt, 1)

		errRate := a.errRate
		reqCnt := atomic.LoadUint64(&a.reqCnt)
		errCnt := atomic.LoadUint64(&a.errCnt)
		actualErrRate := float64(errCnt) / float64(reqCnt)

		// 错误率 大于了 阈值 -> 系统崩了 （将请求存入数据库，用于异步发送）
		if actualErrRate > errRate {
			now := time.Now().UnixMilli()
			// 触发了限流（将请求存入数据库）
			err = a.reqRepo.Create(ctx, domain.Request{
				Tpl:     tpl,
				Args:    args,
				Numbers: numbers,
				Ctime:   now,
				Utime:   now,
			})
		}

		return err
	}
	// TODO: 同步转异步后，还需要返回err错误码吗
}

func (a *AsyncSmsService) AsyncSend(ctx context.Context) {
	go func() {
		reqs, err := a.reqRepo.Find(ctx)
		if err != nil {
			return
		}

		for _, req := range reqs {
			for i := 0; i < a.maxRetry; i++ {
				err := a.smsSvc.Send(ctx, req.Tpl, req.Args, req.Numbers...)
				if err == nil {
					break
				}
			}
		}
	}()
}
