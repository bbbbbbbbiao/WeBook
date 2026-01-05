package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
		session.Options(sessions.Options{
			MaxAge: 60,
		})
		updateTime := session.Get("updateTime")
		if updateTime == nil {
			session.Set("updateTime", time.Now().UnixMilli())
			session.Save()
		}

		if time.Now().UnixMilli()-updateTime.(int64) > 10*1000 {
			session.Set("updateTime", time.Now().UnixMilli())
			session.Save()
		}
		ctx.Set("userId", session.Get("userId"))
	}
}
