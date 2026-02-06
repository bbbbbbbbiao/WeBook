package ratelimit

import "context"

/**
 * @author: biao
 * @date: 2026/2/1 下午2:50
 * @description:
 */

type Limiter interface {
	// Limit 有没有触发限流。
	// key 限流对象
	// bool 代表是否限流，true 就是要限流
	// err 限流器本身是否有错误
	Limit(ctx context.Context, key string) (bool, error)
}
