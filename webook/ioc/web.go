package ioc

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web/middleware"
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

/**
 * @author: biao
 * @date: 2026/1/23 下午10:06
 * @description:
 */

// 你的gin框架呢？你的路由注册呢？你的中间件呢？
// 进行路由注册
// 进行中间件的插入
func InitWebServe(middlewares []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(cmd redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		CorsHdl(),
		middleware.NewLoginMiddlewareBuilder().
			IgnorePath("/users/login").
			IgnorePath("/users/signup").
			IgnorePath("/users/login_sms/code/send").
			IgnorePath("/users/login_sms").
			Build(),
		RateLimitHdl(cmd),
	}
}

func CorsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token", "Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.HasPrefix(origin, "youCompany.com")
		},
	})
}

func RateLimitHdl(cmd redis.Cmdable) gin.HandlerFunc {
	// 为啥不单独Init一个限流器，因为其他地方也用到了-不同的-限流器，冲突了
	limiter := ratelimit.NewRedisSlidingWindowLimiter(cmd, time.Second, 100)
	return middleware.NewBuilder(limiter).Build()
}
