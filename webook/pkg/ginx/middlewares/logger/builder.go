package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"io"
)

/**
 * @author: biao
 * @date: 2026/3/14 下午8:26
 * @description:
 */

// 实现系统出入记录日志功能
// 出口：是指调用系统调用第三方
// 入口：是指系统收到了请求，并返回响应

// 注意点
// 1. 小心日志内容过多，URL可能很长，请求体、响应体都可能很大，需要考虑是否全部输出到日志中
// 2. 考虑1的问题，以及用户可能换用不同的的日志框架，所以需要由足够的灵活性
// 3. 考虑动态开关（考虑使用Viper去监听配置文件），注意并发安全

type MiddlewareBuilder struct {
	allowReqBody  *atomic.Bool // 请求和响应是否需要打印，自定义（可能会很长，所以自定义）
	allowRespBody bool
	loggerFunc    func(ctx *gin.Context, al *AccessLog) // 使用自定义的方法输出日志，包括了选用的日志框架也是由外部调用者决定
}

func NewMiddlewareBuilder(loggerFunc func(ctx *gin.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		allowReqBody: atomic.NewBool(false),
		loggerFunc:   loggerFunc,
	}
}

func (b *MiddlewareBuilder) AllowReqBody(allowReqBody bool) *MiddlewareBuilder {
	b.allowReqBody.Store(allowReqBody)
	return b
}

func (b *MiddlewareBuilder) AllowRespBody(allowRespBody bool) *MiddlewareBuilder {
	b.allowRespBody = allowRespBody
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var al AccessLog

		al.Method = ctx.Request.Method
		// Url 也有可能很长的
		path := ctx.Request.URL.String()
		if len(path) > 1024 {
			path = path[:1024]
		}
		al.Path = path

		// 请求体
		if b.allowReqBody.Load() && ctx.Request.Body != nil {
			//reqBody, _ := io.ReadAll(ctx.Request.Body)
			reqBody, _ := ctx.GetRawData() // 这两个都是读请求体的
			// 因为这些都是流，读一次就没有了，所以需要放回去
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			if len(reqBody) > 1024 {
				reqBody = reqBody[:1024]
			}

			al.ReqBody = string(reqBody)
		}

		// 为了不会被Next给panic掉，我们使用defer
		defer func() {

			// 响应体 (ctx.Writer中没有方法提供对应的响应体，所以需要装饰器)
			// ctx.Writer 支持将数据回写，但是不支持中途获取到数据
			if b.allowRespBody {
				ctx.Writer = &responseWriter{
					ResponseWriter: ctx.Writer,
					al:             &al,
				}
			}

			// 这里进行自定义的日志输出
			b.loggerFunc(ctx, &al)
		}()

		ctx.Next()
	}
}

// 这里使用组合类型的装饰器，因为我们不需要装饰其中所有的方法
type responseWriter struct {
	gin.ResponseWriter // 将响应回写回去
	al                 *AccessLog
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.al.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseWriter) Write(data []byte) (int, error) {
	d := string(data)
	if len(d) > 1024 {
		d = d[:1024]
	}
	r.al.RespBody = d
	return r.ResponseWriter.Write(data)
}

func (r *responseWriter) WriteString(data string) (int, error) {
	if len(data) > 1024 {
		data = data[:1024]
	}
	r.al.RespBody = data
	return r.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	ReqBody    string `json:"req_body"`
	RespBody   string `json:"resp_body"`
	StatusCode int    `json:"status_code"`
}
