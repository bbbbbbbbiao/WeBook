package cache

import (
	"context"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"sync"
	"time"
)

/**
 * @author: biao
 * @date: 2026/3/8 下午7:46
 * @description: 本地内存实现的验证码功能
 */

// LocalCodeCache 技术选型考虑的点
//  1. 功能性：功能是否能够完全覆盖你的需求
//  2. 社区和支持度：社区是否活跃，文档是否齐全
//     以及百度（搜索引擎）能不搜索到你需要的各种信息，有没有帮你踩过坑
//  3. 非功能性：
//     易用性（用户友好度，学习曲线平滑）
//     扩展性（如果开源软件的某些功能需要定制，框架是否支持定制，以及定制的难度高不高）
//     性能（追求性能的公司，往往有能力自研）

type LocalCodeCache struct {
	cache      *lru.Cache[string, CodeConfig]
	expiration time.Duration
	lock       sync.Mutex
	maps       sync.Map
}

func NewLocalCodeCache(cache *lru.Cache[string, CodeConfig], expiration time.Duration) *LocalCodeCache {
	return &LocalCodeCache{
		cache:      cache,
		expiration: expiration,
	}
}

func (l *LocalCodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (l *LocalCodeCache) Set(ctx context.Context, phone string, biz string, code string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	key := l.Key(biz, phone)

	// 优化1，直接对key进行加锁（缺点，如果key很多，那么这个maps就是占据很多内存）
	//lock, _ := l.maps.LoadOrStore(key, &sync.Mutex{})
	//lock.(*sync.Mutex).Lock()
	//defer lock.(*sync.Mutex).Unlock()

	// 优化2
	//lock, _ := l.maps.LoadOrStore(key, &sync.Mutex{})
	//lock.(*sync.Mutex).Lock()
	//defer func() {
	//	l.maps.Delete(key)
	//	lock.(*sync.Mutex).Unlock()
	//}()

	value, ok := l.cache.Get(key)
	if ok {
		if value.Expiration.Sub(time.Now()) > 9*time.Minute {
			// 发送太频繁
			return ErrCodeSendTooMany
		}
		// 小于9分钟（大于1分钟），可以重新发送验证码

	}
	config := CodeConfig{
		Expiration: time.Now().Add(l.expiration),
		Code:       code,
		CntKey:     3,
	}

	l.cache.Add(key, config)
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, phone string, biz string, inputCode string) (bool, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	key := l.Key(biz, phone)
	value, ok := l.cache.Get(key)

	if !ok {
		return false, ErrKeyNotExist
	}

	if value.Expiration.Sub(time.Now()) < 0 {
		return false, errors.New("验证码已过期")
	}

	if value.CntKey <= 0 {
		return false, ErrCodeVerifyTooMany
	}

	value.CntKey--
	return value.Code == inputCode, nil

}

type CodeConfig struct {
	Expiration time.Time
	Code       string
	CntKey     int
}
