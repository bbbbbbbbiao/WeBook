package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @author: biao
 * @date: 2025/12/18 下午10:03
 * @description: 用户模块的web层
 */

type UserHandler struct {
	emailExpr    *regexp.Regexp
	passwordExpr *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		EmailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		PasswordRegexPattern = `^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9])(?=.*[!@#$%^&*()_+\-=\[\]{}|;':",./<>?]).{8,}$`
	)
	emailExpr := regexp.MustCompile(EmailRegexPattern, regexp.None)
	passwordExpr := regexp.MustCompile(PasswordRegexPattern, regexp.None)

	return &UserHandler{
		emailExpr:    emailExpr,
		passwordExpr: passwordExpr,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SingUp)
	ug.POST("/login", u.Login)
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

	ctx.JSON(http.StatusOK, "注册成功")
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
