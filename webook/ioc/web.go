package ioc

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web"
	iJwt "github.com/bbbbbbbbiao/WeBook/webook/internal/web/jwt"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web/middleware"
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/ginx/middlewares/logger"
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/ginx/middlewares/ratelimit"
	logger2 "github.com/bbbbbbbbiao/WeBook/webook/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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

func InitMiddlewares(cmd redis.Cmdable,
	l logger2.LoggerV2,
	jwtHandler iJwt.Handler) []gin.HandlerFunc {

	// 不在中间件中进行监听，而是在初始化的时候进行监听
	mb := logger.NewMiddlewareBuilder(func(ctx *gin.Context, al *logger.AccessLog) {
		l.Info("系统出入口记录", logger2.String("al", al))
	}).AllowReqBody(true).AllowRespBody(true)

	viper.OnConfigChange(func(in fsnotify.Event) {
		ok := viper.GetBool("web.logReq")
		mb.AllowReqBody(ok)
	})
	return []gin.HandlerFunc{
		CorsHdl(),
		mb.Build(),
		middleware.NewJWTLoginMiddlewareBuilder(jwtHandler).
			IgnorePath("/users/login").
			IgnorePath("/users/signup").
			IgnorePath("/users/login_sms/code/send").
			IgnorePath("/users/login_sms").
			IgnorePath("/users/refresh_token").
			Build(),
		RateLimitHdl(cmd),
	}
}

func CorsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token", "Content-Length"},
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
