package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/cache"
	cachemocks "github.com/bbbbbbbbiao/WeBook/webook/internal/repository/cache/mocks"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao"
	daomocks "github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

/**
 * @author: biao
 * @date: 2026/1/29 下午3:09
 * @description:
 */

func TestUserRepository_FindUserById(t *testing.T) {
	now := time.Now().Unix()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		// 输入
		ctx context.Context
		id  int64
		// 输出
		wantErr  error
		wantUser domain.User
	}{
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)

				ud := daomocks.NewMockUserDao(ctrl)
				ud.EXPECT().FindUserById(gomock.Any(), int64(123)).
					Return(dao.User{
						Id: 123,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "baibiao",
						Phone: sql.NullString{
							String: "123456789",
							Valid:  true,
						},
						NickName:     "111",
						Birthday:     "2002-06-21",
						Introduction: "哥只是个传说",
						Ctime:        now,
					}, nil)

				uc.EXPECT().Set(gomock.Any(), domain.User{
					Id:           123,
					Email:        "123@qq.com",
					Password:     "baibiao",
					Phone:        "123456789",
					NickName:     "111",
					Birthday:     "2002-06-21",
					Introduction: "哥只是个传说",
					Ctime:        now,
				})
				return ud, uc
			},
			ctx:     context.Background(),
			id:      123,
			wantErr: nil,
			wantUser: domain.User{
				Id:           123,
				Email:        "123@qq.com",
				Password:     "baibiao",
				Phone:        "123456789",
				NickName:     "111",
				Birthday:     "2002-06-21",
				Introduction: "哥只是个传说",
				Ctime:        now,
			},
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:           123,
						Email:        "123@qq.com",
						Password:     "baibiao",
						Phone:        "123456789",
						NickName:     "111",
						Birthday:     "2002-06-21",
						Introduction: "哥只是个传说",
						Ctime:        now,
					}, nil)

				ud := daomocks.NewMockUserDao(ctrl)
				return ud, uc
			},
			ctx:     context.Background(),
			id:      123,
			wantErr: nil,
			wantUser: domain.User{
				Id:           123,
				Email:        "123@qq.com",
				Password:     "baibiao",
				Phone:        "123456789",
				NickName:     "111",
				Birthday:     "2002-06-21",
				Introduction: "哥只是个传说",
				Ctime:        now,
			},
		},
		{
			name: "缓存未命中，查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)

				ud := daomocks.NewMockUserDao(ctrl)
				ud.EXPECT().FindUserById(gomock.Any(), int64(123)).
					Return(dao.User{}, errors.New("gomock db error"))
				return ud, uc
			},
			ctx:      context.Background(),
			id:       123,
			wantErr:  errors.New("gomock db error"),
			wantUser: domain.User{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)

			u, err := repo.FindUserById(tc.ctx, tc.id)
			//require.NoError(t, err)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}
