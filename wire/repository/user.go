package repository

import "github.com/bbbbbbbbiao/WeBook/wire/repository/dao"

/**
 * @author: biao
 * @date: 2026/1/23 上午10:38
 * @description:
 */

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao: dao}
}
