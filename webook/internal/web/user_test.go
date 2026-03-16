package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service"
	svcmocks "github.com/bbbbbbbbiao/WeBook/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**
 * @author: biao
 * @date: 2026/1/27 下午9:51
 * @description:
 */

// HTTP 相关的测试
func TestUserHandler_signUp(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) service.UserService
		// 输入
		reqBody string
		// 输出
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				// 使用mock模拟UserService接口的实现
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Hello#World123",
				}).Return(nil)
				return userSvc
			},
			reqBody:  `{"email":"123@qq.com","password":"Hello#World123","confirmPassword":"Hello#World123"}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对，bind异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				// 使用mock模拟UserService接口的实现
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody:  `{"email":"123@qq.com","password":"Hello#World123","confirmPassword":"Hello#World123"`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				// 使用mock模拟UserService接口的实现
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody:  `{"email":"123@qq.com","password":"Hello#World123","confirmPassword":"Hello#World124"}`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				// 使用mock模拟UserService接口的实现
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody:  `{"email":"@qq.com","password":"Hello#World123","confirmPassword":"Hello#World123"}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱格式错误",
		},
		{
			name: "密码格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				// 使用mock模拟UserService接口的实现
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody:  `{"email":"123@qq.com","password":"ello#orld123","confirmPassword":"ello#orld123"}`,
			wantCode: http.StatusOK,
			wantBody: "密码格式错误，密码必须包含大写字母、小写字母、数字和特殊字符，长度不能小于8位",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				// 使用mock模拟UserService接口的实现
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Hello#World123",
				}).Return(service.ErrUserDuplicate)
				return userSvc
			},
			reqBody:  `{"email":"123@qq.com","password":"Hello#World123","confirmPassword":"Hello#World123"}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				// 使用mock模拟UserService接口的实现
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Hello#World123",
				}).Return(errors.New("随便一个错误"))
				return userSvc
			},
			reqBody:  `{"email":"123@qq.com","password":"Hello#World123","confirmPassword":"Hello#World123"}`,
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 注册go mock的Controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 注册gin
			server := gin.Default()
			// 注册UserHandler
			h := NewUserHandler(tc.mock(ctrl), nil)
			// 注册路由
			h.RegisterRoutes(server)

			// 构造请求体
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			// 这里的错误只有几种情况-方法类型错误、请求路径错误、构造的数据错误
			require.NoError(t, err)

			// 构造响应体
			resp := httptest.NewRecorder() // 记住就好
			t.Log(resp)

			// 将请求传入gin，由gin处理对应请求
			// 将响应回写与resp
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestUserHandler_LoginSms(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		reqBody string

		wantCode int
		wantBody Result
	}{
		{
			name: "验证码校验通过",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(
					gomock.Any(),
					biz,
					"123456",
					"360520",
				).Return(true, nil)

				userSvc.EXPECT().FindOrCreate(
					gomock.Any(),
					"123456",
				).Return(domain.User{
					Id:    123,
					Phone: "123456",
				}, nil)
				return userSvc, codeSvc
			},
			reqBody:  `{"phone":"123456", "code":"360520"}`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Msg: "验证码校验通过",
			},
		},
		{
			name: "参数不对，bind失败",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBody:  `{"phone":"123456", "code":"360520"`,
			wantCode: http.StatusBadRequest, // 参数异常，专用http.StatusBadRequest
			wantBody: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "验证次数已用完，请重新获取验证码",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(
					gomock.Any(),
					biz,
					"123456",
					"360520",
				).Return(false, service.ErrCodeVerifyTooMay)

				return userSvc, codeSvc
			},
			reqBody:  `{"phone":"123456", "code":"360520"}`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 5,
				Msg:  "验证次数已用完，请重新获取验证码",
			},
		},
		{
			name: "校验后的未知错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(
					gomock.Any(),
					biz,
					"123456",
					"360520",
				).Return(false, errors.New("校验后的未知错误"))
				return userSvc, codeSvc
			},
			reqBody:  `{"phone":"123456", "code":"360520"}`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "验证码错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(
					gomock.Any(),
					biz,
					"123456",
					"360520",
				).Return(false, nil)
				return userSvc, codeSvc
			},
			reqBody:  `{"phone":"123456", "code":"360520"}`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 4,
				Msg:  "验证码错误，请重试",
			},
		},
		{
			name: "登录或注册错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(
					gomock.Any(),
					biz,
					"123456",
					"360520",
				).Return(true, nil)

				userSvc.EXPECT().FindOrCreate(
					gomock.Any(),
					"123456",
				).Return(domain.User{}, errors.New("登录或注册错误"))
				return userSvc, codeSvc
			},
			reqBody:  `{"phone":"123456", "code":"360520"}`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()

			userSvc, codeSvc := tc.mock(ctrl)
			h := NewUserHandler(userSvc, codeSvc)

			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			resp := httptest.NewRecorder()
			t.Log(resp)

			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)

			var respRes Result
			_ = json.NewDecoder(resp.Body).Decode(&respRes)
			assert.Equal(t, tc.wantBody, respRes)
		})
	}
}
