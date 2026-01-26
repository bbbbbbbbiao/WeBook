package wire

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
 * @author: biao
 * @date: 2026/1/23 上午10:48
 * @description:
 */

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("dsn"))
	if err != nil {
		panic(err)
	}
	return db
}
