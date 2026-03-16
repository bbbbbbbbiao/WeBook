package web

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service"
	iJwt "github.com/bbbbbbbbiao/WeBook/webook/internal/web/jwt"
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/logger"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
)

/**
 * @author: biao
 * @date: 2025/12/18 下午10:03
 * @description: 用户模块的web层
 */

const biz = "login"

type UserHandler struct {
	cmd              redis.Cmdable
	svc              service.UserService
	codeSvc          service.CodeService
	emailExpr        *regexp.Regexp
	passwordExpr     *regexp.Regexp
	nikeNameExpr     *regexp.Regexp
	birthdayExpr     *regexp.Regexp
	introductionExpr *regexp.Regexp
	iJwt.Handler
	l logger.LoggerV2
}

func NewUserHandler(svc service.UserService,
	codeSvc service.CodeService,
	cmd redis.Cmdable,
	handler iJwt.Handler,
	l logger.LoggerV2) *UserHandler {
	const (
		EmailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		PasswordRegexPattern = `^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9])(?=.*[!@#$%^&*()_+\-=\[\]{}|;':",./<>?]).{8,}$`
		NikeNamePattern      = `^[\u4e00-\u9fa5a-zA-Z0-9_]{1,20}$`
		BirthdayPattern      = `^(19|20)\d{2}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`
		IntroductionPattern  = `^.{0,200}$`
	)
	emailExpr := regexp.MustCompile(EmailRegexPattern, regexp.None)
	passwordExpr := regexp.MustCompile(PasswordRegexPattern, regexp.None)
	nikeNameExpr := regexp.MustCompile(NikeNamePattern, regexp.None)
	birthdayExpr := regexp.MustCompile(BirthdayPattern, regexp.None)
	introductionExpr := regexp.MustCompile(IntroductionPattern, regexp.None)

	return &UserHandler{
		cmd:              cmd,
		svc:              svc,
		codeSvc:          codeSvc,
		emailExpr:        emailExpr,
		passwordExpr:     passwordExpr,
		nikeNameExpr:     nikeNameExpr,
		birthdayExpr:     birthdayExpr,
		introductionExpr: introductionExpr,
		Handler:          handler,
		l:                l,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SingUp)
	ug.POST("/login", u.Login)
	ug.POST("/JWTLogin", u.JWTLogin)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSms)
	ug.POST("/refresh_token", u.RefreshToken)
	ug.POST("/logout", u.Logout)
}

// 注册
func (u *UserHandler) SingUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析数据到 req 里面
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 验证密码是否一致
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	// 验证邮箱格式是否正确
	ok, err := u.emailExpr.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	// 验证密码格式是否正确
	ok, err = u.passwordExpr.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码格式错误，密码必须包含大写字母、小写字母、数字和特殊字符，长度不能小于8位")
		return
	}

	// 调用 SVC 的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicate {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

// 发送验证码
func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "1111",
		})
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	//if err == service.ErrCodeSendTooMany {
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: 5,
	//		Msg:  "验证码发送太频繁",
	//	})
	//	return
	//}
	//if err != nil {
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: 5,
	//		Msg:  "系统错误",
	//	})
	//	return
	//}
	//ctx.JSON(http.StatusOK, Result{
	//	Msg: "发送成功",
	//})
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

}

// TODO: 测试
// 验证验证码 （注册或登录）
func (u *UserHandler) LoginSms(ctx *gin.Context) {
	type VerifyReq struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req VerifyReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)

	if err == service.ErrCodeVerifyTooMay {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证次数已用完，请重新获取验证码",
		})

		// 当等级不是那么高时，我们可以用warn日志
		// 但是这里可以在告警系统中进行配置
		// 比如说规则，一分钟内出现超过100次WARN，你就告警
		zap.L().Warn("验证次数已用完", zap.Error(err),
			// 这里不能直接打印手机号码，因为手机号码是敏感数据
			// 1. 进行手机号加密处理（能反向解密的，但是打印本就是高频操作，再进行加密的话性能就会很差，加密本就吃CPU和内存的）
			// 2. 脱敏处理：130****5678，没啥用，找不到原始的手机号
			zap.String("phone", req.Phone))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//zap.L().Error("校验验证码出错", zap.Error(err))
		u.l.Info("校验验证码出错", logger.String("err", err))
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误，请重试",
		})
		return
	}

	// 验证通过后，该如何操作呢？
	// 进行注册或登录
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 生成Token
	err = u.SetLoginToken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过",
	})
}

// 账号JWT登录
func (u *UserHandler) JWTLogin(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq

	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 生成Token
	err = u.SetLoginToken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, "登录成功")
}

// 账号Session 登录
func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq

	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	session := sessions.Default(ctx)
	session.Set("userId", user.Id)
	// sessions 中设置参数
	session.Options(sessions.Options{
		MaxAge: 60,
	})
	err = session.Save()
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, "登录成功")
}

// 编辑
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		NickName     string `json:"nikeName"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ok, err := u.nikeNameExpr.MatchString(req.NickName)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "昵称格式错误，昵称只能包含中文、字母、数字和下划线，长度不能超过20个字符")
	}

	ok, err = u.birthdayExpr.MatchString(req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "生日格式错误，生日必须符合yyyy-MM-dd格式")
		return
	}

	ok, err = u.introductionExpr.MatchString(req.Introduction)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "简介格式错误，简介长度不能超过200个字符")
		return
	}
	err = u.svc.Edit(ctx, domain.User{
		Id:           ctx.GetInt64("userId"),
		NickName:     req.NickName,
		Birthday:     req.Birthday,
		Introduction: req.Introduction,
	})

	if err == service.ErrUserNotFound {
		ctx.String(http.StatusOK, "用户不存在")
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, "编辑成功")
}

// 获取用户信息
func (u *UserHandler) Profile(ctx *gin.Context) {
	userClaims, _ := ctx.Get("userClaims")
	claims, ok := userClaims.(*iJwt.AccessClaims)
	if !ok || claims.UserId == 0 {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	user, err := u.svc.Profile(ctx, claims.UserId)
	if err == service.ErrUserNotFound {
		ctx.String(http.StatusOK, "用户不存在")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// TODO: 这里不能直接将domain暴露出去
	// 首先不能让别人知道你的domain，同时你的密码也在里面
	ctx.JSON(http.StatusOK, user)
}

// 刷新Token
func (u *UserHandler) RefreshToken(ctx *gin.Context) {

	// 当调用该接口时，Header中携带的便是RefreshToken
	tokenStr := u.ExtractToken(ctx)
	if tokenStr == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	var refreshClaims iJwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return iJwt.RefreshTokenKey, nil
	})

	if err != nil || token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 查询一下是否已经退出登录
	err = u.CheckSession(ctx, refreshClaims.SsId)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = u.SetJWTToken(ctx, refreshClaims.UserId, refreshClaims.SsId)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 用来标记哪个位置打印的日志
		//zap.L().Error("设置JWTToken失败",
		//	zap.Error(err), zap.String("method", "UserHandler.RefreshToken"))
		// 正常来说Msg就应该能包含足够的定位信息
		zap.L().Error("UserHandler:RefreshToken 设置JWTToken失败", zap.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "刷新成功",
	})
}

// 登出
func (u *UserHandler) Logout(ctx *gin.Context) {
	tokenStr := u.ExtractToken(ctx)
	var accessClaims iJwt.AccessClaims

	token, err := jwt.ParseWithClaims(tokenStr, &accessClaims, func(token *jwt.Token) (interface{}, error) {
		return iJwt.AccessTokenKey, nil
	})

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	if token == nil || !token.Valid {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	ssid := accessClaims.SsId
	// 将 ssid 加入到Redis中，表示该Token已失效
	err = u.ClearToken(ctx, ssid)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "登出成功",
	})
}
