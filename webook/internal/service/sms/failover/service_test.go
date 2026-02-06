package failover

import (
	"context"
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	smsmocks "github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"testing"
)

/**
 * @author: biao
 * @date: 2026/2/4 下午2:26
 * @description:
 */

func TestFailOverService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		ctx     context.Context
		tpl     string
		args    []string
		numbers []string

		wantErr error
	}{
		{
			name: "一次成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svcs := make([]sms.Service, 0, 2)
				svcs = append(svcs, svc1, svc2)

				svc1.EXPECT().Send(gomock.Any(), "xxxx", []string{"1111"}, []string{"123456"}).
					Return(nil)

				return svcs
			},

			ctx:     context.Background(),
			tpl:     "xxxx",
			args:    []string{"1111"},
			numbers: []string{"123456"},

			wantErr: nil,
		},
		{
			name: "重试成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svcs := make([]sms.Service, 0, 2)
				svcs = append(svcs, svc1, svc2)

				svc1.EXPECT().Send(gomock.Any(), "xxxx", []string{"1111"}, []string{"123456"}).
					Return(errors.New("发送不了"))
				svc2.EXPECT().Send(gomock.Any(), "xxxx", []string{"1111"}, []string{"123456"}).
					Return(nil)
				return svcs
			},

			ctx:     context.Background(),
			tpl:     "xxxx",
			args:    []string{"1111"},
			numbers: []string{"123456"},

			wantErr: nil,
		},
		{
			name: "重试最终失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svcs := make([]sms.Service, 0, 2)
				svcs = append(svcs, svc1, svc2)

				svc1.EXPECT().Send(gomock.Any(), "xxxx", []string{"1111"}, []string{"123456"}).
					Return(errors.New("发送不了"))
				svc2.EXPECT().Send(gomock.Any(), "xxxx", []string{"1111"}, []string{"123456"}).
					Return(errors.New("还是失败"))
				return svcs
			},

			ctx:     context.Background(),
			tpl:     "xxxx",
			args:    []string{"1111"},
			numbers: []string{"123456"},
			wantErr: errors.New("发送失败，所有供应商都尝试过了"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewFailOverService(tc.mock(ctrl))

			err := svc.Send(tc.ctx, tc.tpl, tc.args, tc.numbers...)

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
