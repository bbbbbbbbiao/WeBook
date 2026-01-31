package cache

import (
	"context"
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/cache/redismocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"testing"
)

/**
 * @author: biao】、
 * @date: 2026/1/29 下午10:21
 * @description:
 */

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable
		// 输入
		ctx   context.Context
		phone string
		biz   string
		code  string

		// 输出
		wantErr error
	}{
		{
			name: "缓存设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(
					gomock.Any(),
					luaSetCode,
					[]string{"phone_code:login:123456789"},
					[]any{"336600"},
				).Return(res)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "123456789",
			biz:     "login",
			code:    "336600",
			wantErr: nil,
		},
		{
			name: "验证码发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(
					gomock.Any(),
					luaSetCode,
					[]string{"phone_code:login:123456789"},
					[]any{"336600"},
				).Return(res)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "123456789",
			biz:     "login",
			code:    "336600",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "redis返回其他值，系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-2))
				cmd.EXPECT().Eval(
					gomock.Any(),
					luaSetCode,
					[]string{"phone_code:login:123456789"},
					[]any{"336600"},
				).Return(res)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "123456789",
			biz:     "login",
			code:    "336600",
			wantErr: errors.New("系统错误"),
		},
		{
			name: "缓存设置失败",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("缓存设置失败"))
				cmd.EXPECT().Eval(
					gomock.Any(),
					luaSetCode,
					[]string{"phone_code:login:123456789"},
					[]any{"336600"},
				).Return(res)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "123456789",
			biz:     "login",
			code:    "336600",
			wantErr: errors.New("缓存设置失败"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cache := NewRedisCodeCache(tc.mock(ctrl))

			err := cache.Set(tc.ctx, tc.phone, tc.biz, tc.code)

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
