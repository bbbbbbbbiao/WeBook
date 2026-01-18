package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

/**
 * @author: biao
 * @date: 2026/1/11 下午12:33
 * @description: 用户缓存
 */

var ErrKeyNotExist = redis.Nil

// 面向接口编程
type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

// 依赖注入
func NewUserCache(cmd redis.Cmdable, expiration time.Duration) *UserCache {
	return &UserCache{
		cmd:        cmd,
		expiration: expiration,
	}
}

func (cache *UserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (cache *UserCache) Set(ctx context.Context, user domain.User) error {

	key := cache.Key(user.Id)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return cache.cmd.Set(ctx, key, data, cache.expiration).Err()
}

// 一般使用error做数据信息的传递
// error 为空 可以获取
// error 不为空
//   - ErrKeyNotExist Key不存在
//   - nil 可能没有命中/redis崩掉了
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {

	key := cache.Key(id)
	data, err := cache.cmd.Get(ctx, key).Result()

	if err != nil {
		return domain.User{}, err
	}

	var user domain.User
	err = json.Unmarshal([]byte(data), &user)
	return user, err
}
