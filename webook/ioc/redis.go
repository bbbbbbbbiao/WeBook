package ioc

import (
	"github.com/bbbbbbbbiao/WeBook/webook/config"
	"github.com/redis/go-redis/v9"
	"time"
)

/**
 * @author: biao
 * @date: 2026/1/23 下午9:38
 * @description:
 */

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
}

func InitTimeDuration() time.Duration {
	return time.Minute * 30
}
