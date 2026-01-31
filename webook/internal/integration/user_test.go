package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web"
	"github.com/bbbbbbbbiao/WeBook/webook/ioc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/**
* @author: biao
* @date: 2026/1/30 下午4:14
* @description: 集成测试
 */

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string

		// 要考虑准备数据、验证数据和清理数据
		before func(t *testing.T) // 准备数据
		after  func(t *testing.T) // 验证、清理数据

		reqBody string

		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {

			},
			// 清理和验证数据
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 保证可以多次运行 -> 清理数据
				result, err := rdb.GetDel(ctx, "phone_code:login:123456789").Result()
				cancel()
				assert.NoError(t, err)
				// 验证数据
				assert.True(t, len(result) == 6)
			},

			reqBody: `{"phone":"123456789"}`,

			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				rdb.Set(ctx, "phone_code:login:123456789", "111111", time.Minute*9+time.Second*30)
				cancel()
			},
			// 清理和验证数据
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 保证可以多次运行 -> 清理数据
				result, err := rdb.Get(ctx, "phone_code:login:123456789").Result()
				cancel()
				assert.NoError(t, err)
				// 验证数据
				assert.True(t, result == "111111")
			},

			reqBody: `{"phone":"123456789"}`,

			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "发送太频繁，请稍后再试",
			},
		},
		{
			name: "没有设置过期时间",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				rdb.Set(ctx, "phone_code:login:123456789", "111111", 0)
				cancel()
			},
			// 清理和验证数据
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 保证可以多次运行 -> 清理数据
				result, err := rdb.Get(ctx, "phone_code:login:123456789").Result()
				cancel()
				assert.NoError(t, err)
				// 验证数据
				assert.True(t, result == "111111")
			},

			reqBody: `{"phone":"123456789"}`,

			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			req, err := http.NewRequest(
				http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			var respRes web.Result
			_ = json.NewDecoder(resp.Body).Decode(&respRes)
			assert.Equal(t, tc.wantBody, respRes)

			tc.after(t)
		})
	}
}
