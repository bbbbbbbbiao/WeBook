package ioc

import (
	"github.com/bbbbbbbbiao/WeBook/webook/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
 * @author: biao
 * @date: 2026/1/23 下午9:35
 * @description:
 */

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	return db
}
