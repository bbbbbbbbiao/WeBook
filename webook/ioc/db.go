package ioc

import (
	"fmt"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao"
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

/**
 * @author: biao
 * @date: 2026/1/23 下午9:35
 * @description:
 */

func InitDB(l logger.LoggerV2) *gorm.DB {

	// 读取某个类型的配置
	//dsn := viper.GetString("db.config.dsn")

	// 自定义配置结构体，反系列化获得配置信息
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var c Config
	// remote 不支持 key 的切割
	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败 %v", err))
	}
	fmt.Println("dsn:", c.DSN)

	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询阈值
			// 一般设置为 50 ~ 100 ms
			// SQL 查询必定命中索引，一次查询最好走一次磁盘IO
			// 一次磁盘 IO 是不到10ms，所以当阈值达到了50ms，则说明要么是SQL写的不好，要么是数据量太大
			SlowThreshold: time.Millisecond * 50,
			// 是否打印"未找到记录"，有些需要打印（查询用户信息），有些不需要打印（订单信息）
			IgnoreRecordNotFoundError: true,
			// 参数能否拼接在日志中
			ParameterizedQueries: true,
			LogLevel:             glogger.Info,
		}),
	})

	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

// 适配器
// 条件：只能是一个抽象方法时，使用适配器
// 使用：gormLoggerFunc(l.Debug) // 相当于类型转换了
// 里面传的参数和l.Debug的参数是一致的
type gormLoggerFunc func(msg string, args ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

// 举例
type DoSomething interface {
	DoABC() string
}

type DoSomethingFunc func() string

func (d DoSomethingFunc) DoABC() string {
	return d()
}
