package dao

import "gorm.io/gorm"

/**
 * @author: biao
 * @date: 2025/12/23 上午10:45
 * @description: 初始化数据库表
 */

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
