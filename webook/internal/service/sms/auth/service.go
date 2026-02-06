package auth

import (
	"context"
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"github.com/golang-jwt/jwt/v5"
)

/**
 * @author: biao
 * @date: 2026/2/5 下午4:25
 * @description: 装饰器-权限控制，仅允许得到申请的部门使用
 */

type AuthService struct {
	smsSvc sms.Service
	key    string // JWT 中，使用加密的密钥
}

func NewAuthService(smsSvc sms.Service) sms.Service {
	return &AuthService{
		smsSvc: smsSvc,
	}
}

// Send
// Auth 中的 tpl 就是 所需的Token(静态)，这里的Token由申请获得，调用时携带这个token
// 其他的 tpl 就是模板
func (a *AuthService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	var aClaims authClaims
	token, err := jwt.ParseWithClaims(tpl, &aClaims, func(token *jwt.Token) (interface{}, error) {
		return a.key, nil
	})

	if err != nil {
		return err
	}

	if token == nil || !token.Valid {
		return errors.New("无效的Token")
	}

	return a.smsSvc.Send(ctx, aClaims.actualTpl, args, numbers...)
}

type authClaims struct {
	jwt.RegisteredClaims
	actualTpl string
}
