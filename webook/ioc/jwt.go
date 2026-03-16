package ioc

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web/jwt"
	"github.com/redis/go-redis/v9"
)

/**
 * @author: biao
 * @date: 2026/3/8 下午8:54
 * @description:
 */

func InitJwtHandler(cmd redis.Cmdable) jwt.Handler {
	return jwt.NewRedisJwtHandler(cmd)
}
