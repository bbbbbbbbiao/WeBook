package web

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

/**
 * @author: biao
 * @date: 2025/12/18 下午10:03
 * @description: 用户模块的web层
 */

type UserHandler struct {
	svc              *service.UserService
	emailExpr        *regexp.Regexp
	passwordExpr     *regexp.Regexp
	nikeNameExpr     *regexp.Regexp
	birthdayExpr     *regexp.Regexp
	introductionExpr *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
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
		svc:              svc,
		emailExpr:        emailExpr,
		passwordExpr:     passwordExpr,
		nikeNameExpr:     nikeNameExpr,
		birthdayExpr:     birthdayExpr,
		introductionExpr: introductionExpr,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SingUp)
	ug.POST("/login", u.Login)
	ug.POST("/JWTLogin", u.JWTLogin)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

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

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, "注册成功")
}

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

	// 生成jwt结构体
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserId:    user.Id,
		UserAgent: ctx.Request.UserAgent(),
	})

	tokenStr, err := token.SignedString([]byte("3E7QYaUxM5tMhDWwd5HphdYWND7WR2Vx"))

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.Header("x-jwt-token", tokenStr)
	ctx.JSON(http.StatusOK, "登录成功")
}

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

func (u *UserHandler) Profile(ctx *gin.Context) {
	userClaims, _ := ctx.Get("userClaims")
	claims, ok := userClaims.(*UserClaims)
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
	ctx.JSON(http.StatusOK, user)
}

// 声明一个我自己的放到token中数据
type UserClaims struct {
	jwt.RegisteredClaims
	UserId    int64
	UserAgent string
}
