package middleware

import (
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

/**
 * @author: biao
 * @date: 2026/2/1 下午3:34
 * @description: 限流中间件
 */

type Builder struct {
	limiter   ratelimit.Limiter
	genKeyFun func(ctx *gin.Context) string
	logFun    func(msg any, args ...any)
}

// NewBuilder 初始化Builder(使用默认限流方式、日志方式)，
// genKeyFun: 默认使用 IP 限流
// logFun: 默认使用 log.Println()
func NewBuilder(limiter ratelimit.Limiter) *Builder {
	return &Builder{
		limiter: limiter,
		genKeyFun: func(ctx *gin.Context) string {
			var v strings.Builder
			v.WriteString("ip-limiter")
			v.WriteString(":")
			v.WriteString(ctx.ClientIP())
			return v.String()
		},
		logFun: func(msg any, args ...any) {
			v := make([]any, 0, len(args)+1)
			v = append(v, msg)
			v = append(v, args...)
			log.Println(v...)
		},
	}
}

// 自定义限流方式
func (b *Builder) SetKeyFunc(fn func(ctx *gin.Context) string) {
	b.genKeyFun = fn
}

// 自定义日志方式
func (b *Builder) SetLogFunc(fn func(msg any, args ...any)) {
	b.logFun = fn
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			b.logFun(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// limited 为 true 表示限流，直接不处理该请求
		if limited {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Next()
	}
}

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	return b.limiter.Limit(ctx, b.genKeyFun(ctx))
}
