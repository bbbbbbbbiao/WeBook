package startup

import (
	logger2 "github.com/bbbbbbbbiao/WeBook/webook/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

/**
 * @author: biao
 * @date: 2026/3/19 上午9:40
 * @description:
 */

func InitDB() *gorm.DB {
	logger, _ := zap.NewDevelopment() // zap日志
	l := logger2.NewZapLogger(logger) // 自定义日志

	viper.SetConfigFile("E:\\Project\\my_project\\WeBook\\webook\\config\\dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	type DBConfig struct {
		DSN string `json:"dsn"`
	}

	var dbConfig DBConfig
	err = viper.UnmarshalKey("db", &dbConfig)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(dbConfig.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Info), glogger.Config{
			SlowThreshold:             time.Millisecond * 50,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
			LogLevel:                  glogger.Info,
		}),
	})

	if err != nil {
		l.Error("获取数据库连接失败", logger2.Error(err))
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger2.Field)

func (g gormLoggerFunc) Printf(msg string, fields ...any) {
	g(msg, logger2.Field{Key: "args", Value: fields})
}
