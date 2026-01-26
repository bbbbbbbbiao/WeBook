package repository

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/cache"
)

/**
 * @author: biao
 * @date: 2026/1/20 下午1:27
 * @description: 验证码缓存
 */

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCacheCodeRepository(cache cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: cache,
	}
}

func (repo *CacheCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, phone, biz, code)
}

func (repo *CacheCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.Verify(ctx, phone, biz, code)
}

//var ErrCodeSendTooMany = cache.ErrCodeSendTooMany
//
//type CacheCodeRepository struct {
//	cache *cache.RedisCodeCache
//}
//
//func NewCacheCodeRepository(cache *cache.RedisCodeCache) *CacheCodeRepository {
//	return &CacheCodeRepository{
//		cache: cache,
//	}
//}
//
//func (repo *CacheCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
//	return repo.cache.Set(ctx, biz, phone, code)
//}
