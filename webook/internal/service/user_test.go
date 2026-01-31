package service

import (
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository"
	repomocks "github.com/bbbbbbbbiao/WeBook/webook/internal/repository/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

/**
 * @author: biao
 * @date: 2026/1/28 下午10:22
 * @description:
 */

func TestUserService_Login(t *testing.T) {
	now := time.Now().Unix()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		// 输入
		user domain.User
		// 输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomocks.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$ixV5ZH3WMEWgSl58MsFClucp6/SBd/bHZOE80d1chOYK1Pfjz1bVy",
						Phone:    "123456789",
						Ctime:    now,
					}, nil)
				return userRepo
			},
			user: domain.User{
				Email:    "123@qq.com",
				Password: "Hello#World123",
			},
			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$ixV5ZH3WMEWgSl58MsFClucp6/SBd/bHZOE80d1chOYK1Pfjz1bVy",
				Phone:    "123456789",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户没找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomocks.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return userRepo
			},
			user: domain.User{
				Email:    "123@qq.com",
				Password: "Hello#World123",
			},
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "密码不正确",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomocks.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$ixV5ZH3WMEWgSl58MsFClucp6/SBd/bHZOE80d1chOYK1Pfjz1bVy",
						Phone:    "123456789",
						Ctime:    now,
					}, nil)
				return userRepo
			},
			user: domain.User{
				Email:    "123@qq.com",
				Password: "Hello#World124",
			},
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomocks.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("系统错误"))
				return userRepo
			},
			user: domain.User{
				Email:    "123@qq.com",
				Password: "Hello#World124",
			},
			wantUser: domain.User{},
			wantErr:  errors.New("系统错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc := NewUserService(tc.mock(ctrl))

			u, err := userSvc.Login(nil, tc.user)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestUserService_Bcrypt(t *testing.T) {
	password, err := bcrypt.GenerateFromPassword([]byte("Hello#World123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(password))
	}
}
