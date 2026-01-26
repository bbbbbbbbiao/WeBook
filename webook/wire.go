//go:build wireinject

package main

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/cache"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web"
	"github.com/bbbbbbbbiao/WeBook/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

/**
 * @author: biao
 * @date: 2026/1/23 上午11:06
 * @description:
 */

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitTimeDuration,

		dao.NewUserDao,

		cache.NewUserCache,
		cache.NewRedisCodeCache,

		repository.NewUserRepository,
		repository.NewCacheCodeRepository,

		service.NewUserService,
		service.NewCodeService,
		ioc.InitSMSService,

		web.NewUserHandler,

		// 中间件？ 注册路由呢？
		ioc.InitWebServe,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
