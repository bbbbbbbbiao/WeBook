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

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(cache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}
}

func (repo *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, phone, biz, code)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.Verify(ctx, phone, biz, code)
}

//var ErrCodeSendTooMany = cache.ErrCodeSendTooMany
//
//type CodeRepository struct {
//	cache *cache.CodeCache
//}
//
//func NewCodeRepository(cache *cache.CodeCache) *CodeRepository {
//	return &CodeRepository{
//		cache: cache,
//	}
//}
//
//func (repo *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
//	return repo.cache.Set(ctx, biz, phone, code)
//}
