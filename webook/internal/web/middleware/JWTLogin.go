package middleware

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

/**
 * @author: biao
 * @date: 2026/1/5 下午9:46
 * @description: 使用JWT进行鉴权中间件
 */

type JWTLoginMiddlewareBuilder struct {
	IgnorePaths []string
}

func NewJWTLoginMiddlewareBuilder() *JWTLoginMiddlewareBuilder {
	return &JWTLoginMiddlewareBuilder{}
}

func (jl *JWTLoginMiddlewareBuilder) IgnorePath(path ...string) *JWTLoginMiddlewareBuilder {
	jl.IgnorePaths = append(jl.IgnorePaths, path...)
	return jl
}

func (jl *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, val := range jl.IgnorePaths {
			if ctx.Request.URL.Path == val {
				return
			}
		}
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(header, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		userClaims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, userClaims, func(token *jwt.Token) (interface{}, error) {
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
		ctx.Set("userClaims", userClaims)
	}
}
