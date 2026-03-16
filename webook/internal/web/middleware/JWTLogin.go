package middleware

import (
	iJwt "github.com/bbbbbbbbiao/WeBook/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
)

/**
 * @author: biao
 * @date: 2026/1/5 下午9:46
 * @description: 使用JWT进行鉴权中间件
 */

type JWTLoginMiddlewareBuilder struct {
	cmd         redis.Cmdable
	IgnorePaths []string
	iJwt.Handler
}

func NewJWTLoginMiddlewareBuilder(handler iJwt.Handler) *JWTLoginMiddlewareBuilder {
	return &JWTLoginMiddlewareBuilder{
		Handler: handler,
	}
}

func (jl *JWTLoginMiddlewareBuilder) IgnorePath(path string) *JWTLoginMiddlewareBuilder {
	jl.IgnorePaths = append(jl.IgnorePaths, path)
	return jl
}

func (jl *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, val := range jl.IgnorePaths {
			if ctx.Request.URL.Path == val {
				return
			}
		}
		tokenStr := jl.Handler.ExtractToken(ctx)

		var userClaims iJwt.AccessClaims
		token, err := jwt.ParseWithClaims(tokenStr, &userClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte("3E7QYaUxM5tMhDWwd5HphdYWND7WR2Vx"), nil
		})

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 增强系统登录安全（User-Agent）
		if ctx.Request.UserAgent() != userClaims.UserAgent {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 查询一下是否已经退出登录
		err = jl.CheckSession(ctx, userClaims.SsId)

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("userClaims", userClaims)
	}
}
