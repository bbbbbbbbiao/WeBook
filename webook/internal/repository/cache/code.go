package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

/**
 * @author: biao
 * @date: 2026/1/18 下午7:11
 * @description:
 */

var (
	ErrCodeSendTooMany   = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooMany = errors.New("验证次数太多")
	ErrUnknownForCode    = errors.New("我也不清楚什么错误")
)

// 编译器会在编译的时候，把 set_code.lua 的代码放进来这个 luaSetCode 变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var verifyCode string

type CodeCache struct {
	cmd redis.Cmdable
}

func NewCodeCache(cmd redis.Cmdable) *CodeCache {
	return &CodeCache{
		cmd: cmd,
	}
}

func (c *CodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *CodeCache) Set(ctx context.Context, phone string, biz string, code string) error {
	// .Int() 表示你Lua脚本写的返回是啥，就用啥接收
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case 0:
		// 没有问题
		return nil
	case -1:
		// 发送太频繁
		return ErrCodeSendTooMany
	default:
		// 系统错误
		return errors.New("系统错误")
	}
}

func (c *CodeCache) Verify(ctx context.Context, phone string, biz string, inputCode string) (bool, error) {
	res, err := c.cmd.Eval(ctx, verifyCode, []string{c.Key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return false, nil
	}
}

//var ErrCodeSendTooMany = errors.New("发送验证码太频繁")
//
////go:embed lua/set_code.lua
//var luaSetCode string
//
//type CodeCache struct {
//	cmd redis.Cmdable
//}
//
//func NewCodeCache(cmd redis.Cmdable) *CodeCache {
//	return &CodeCache{
//		cmd: cmd,
//	}
//}
//
//func (cache *CodeCache) Key(biz, phone string) string {
//	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
//}
//
//func (cache *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
//	res, err := cache.cmd.Eval(ctx, luaSetCode, []string{cache.Key(biz, phone)}, code).Int()
//	if err != nil {
//		return err
//	}
//
//	switch res {
//	case 0:
//		return nil
//	case -1:
//		return ErrCodeSendTooMany
//	default:
//		return err
//	}
//}
