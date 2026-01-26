package dao

import "gorm.io/gorm"

/**
 * @author: biao
 * @date: 2026/1/23 上午10:36
 * @description:
 */

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}
