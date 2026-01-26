// 编译标签
//go:build wireinject

// 让wire来注入这部分代码
package wire

/**
 * @author: biao
 * @date: 2026/1/23 上午10:39
 * @description:
 */

import (
	"github.com/bbbbbbbbiao/WeBook/wire/repository"
	"github.com/bbbbbbbbiao/WeBook/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	// 这里面传入各个组件的初始化方法(传入的是方法，不是调用)
	wire.Build(repository.NewUserRepository, dao.NewUserDao, InitDB)
	return new(repository.UserRepository)
}
