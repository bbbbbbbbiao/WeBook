package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

/**
 * @author: biao
 * @date: 2026/3/6 下午9:52
 * @description:
 */

var (
	AccessTokenKey  = []byte("3E7QYaUxM5tMhDWwd5HphdYWND7WR2Vx")
	RefreshTokenKey = []byte("3E7QYaUxM5sdsafsfsf5HphdYWNdfdfddx")
)

type RedisJwtHandler struct {
	cmd redis.Cmdable
}

func NewRedisJwtHandler(cmd redis.Cmdable) Handler {
	return &RedisJwtHandler{
		cmd: cmd,
	}
}

func (r *RedisJwtHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.SetJWTToken(ctx, uid, ssid)
	err = r.SetRefreshToken(ctx, uid, ssid)
	return err
}

func (r *RedisJwtHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	// 生成jwt结构体
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserId:    uid,
		SsId:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	})

	accessTokenStr, err := accessToken.SignedString(AccessTokenKey)

	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", accessTokenStr)
	return nil
}
func (r *RedisJwtHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims{
		UserId: uid,
		SsId:   ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	})
	refreshTokenStr, err := refreshToken.SignedString(RefreshTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)

	return nil
}

func (r *RedisJwtHandler) ExtractToken(ctx *gin.Context) string {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return ""
	}
	segs := strings.Split(header, " ")
	if len(segs) != 2 {
		return ""
	}

	tokenStr := segs[1]
	return tokenStr
}

func (r *RedisJwtHandler) ClearToken(ctx *gin.Context, ssid string) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	_, err := r.cmd.Set(ctx, fmt.Sprintf("user:ssid:%s", ssid), ssid, time.Hour*24*7).Result()

	return err

}

func (r *RedisJwtHandler) CheckSession(ctx *gin.Context, ssid string) error {
	res, err := r.cmd.Exists(ctx, fmt.Sprintf("user:ssid:%s", ssid)).Result()
	//if err != nil {
	//	return err
	//}
	//if res > 0 {
	//	return errors.New("token 已过期")
	//}
	//return nil

	switch err {
	case redis.Nil:
		return nil
	case nil:
		if res == 0 {
			return nil
		}
		return errors.New("token 已过期")
	default:
		return err
	}
}
