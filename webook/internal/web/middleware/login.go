package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @author: biao
 * @date: 2025/12/25 下午12:27
 * @description: 登录校验中间件
 */

type LoginMiddlewareBuilder struct {
	ignorePaths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (lmb *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	lmb.ignorePaths = append(lmb.ignorePaths, path)
	return lmb
}

func (lmb *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, ignorePath := range lmb.ignorePaths {
			if ignorePath == ctx.Request.URL.Path {
				return
			}
		}
		session := sessions.Default(ctx)
		if session.Get("userId") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("userId", session.Get("userId"))
	}
}
