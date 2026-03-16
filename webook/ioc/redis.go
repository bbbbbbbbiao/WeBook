package ioc

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

/**
 * @author: biao
 * @date: 2026/1/23 下午9:38
 * @description:
 */

func InitRedis() redis.Cmdable {

	type Config struct {
		Addr string `yaml:"addr"`
	}
	var c Config
	if err := viper.UnmarshalKey("redis", &c); err != nil {
		panic(fmt.Errorf("初始化配置失败，%v", err))
	}
	return redis.NewClient(&redis.Options{
		Addr: c.Addr,
	})
}

func InitTimeDuration() time.Duration {
	return time.Minute * 30
}
