package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

/**
 * @author: biao
 * @date: 2026/3/6 下午9:47
 * @description:
 */

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
	ExtractToken(ctx *gin.Context) string
	ClearToken(ctx *gin.Context, ssid string) error
	CheckSession(ctx *gin.Context, ssid string) error
}

type AccessClaims struct {
	jwt.RegisteredClaims
	UserId    int64
	SsId      string
	UserAgent string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	SsId   string
	UserId int64
}
