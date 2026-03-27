//go:build wireinject

package main

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/article"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/cache"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao"
	article2 "github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao/article"
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
		article2.NewArticleDaoImpl,

		cache.NewUserCache,
		cache.NewRedisCodeCache,

		repository.NewUserRepository,
		repository.NewCacheCodeRepository,
		article.NewArticleRepositoryImpl,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleServiceImpl,
		ioc.InitSMSService,
		ioc.InitJwtHandler,

		web.NewUserHandler,
		web.NewArticleHandler,

		// 中间件？ 注册路由呢？
		ioc.InitWebServe,
		ioc.InitMiddlewares,
		//ioc.InitRedisSlidingWindowLimiter,

		ioc.InitLogger,
	)
	return new(gin.Engine)
}
