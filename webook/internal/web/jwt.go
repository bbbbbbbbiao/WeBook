package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtHandler struct {
}

func (u *JwtHandler) SetJWTToken(ctx *gin.Context, uid int64) error {
	// 生成jwt结构体
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserId:    uid,
		UserAgent: ctx.Request.UserAgent(),
	})

	tokenStr, err := token.SignedString([]byte("3E7QYaUxM5tMhDWwd5HphdYWND7WR2Vx"))

	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

// 声明一个我自己的放到token中数据
type UserClaims struct {
	jwt.RegisteredClaims
	UserId    int64
	UserAgent string
}
