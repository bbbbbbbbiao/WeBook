package main

/**
* @author: biao
* @date: 2025/12/18 下午9:46
* @description:
 */

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	// 初始化数据库
	//db := initDB()
	//re := initRedis()

	//server := initWbeServer()

	//server := initJWTWebServer()
	InitLogger()
	InitViperV1()
	server := InitWebServer()
	////u := initUser(db, re)
	////u.RegisterRoutes(server)
	//
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, I am k8s!!!")
	})
	//
	server.Run(":8080")
}

// 加载配置文件方式1
func InitViper() {
	// 设置默认值，如果配置中没有这个key，则使用默认值，如果用key，不管是啥值还是空值，都用配置
	viper.SetDefault("db.config.dsn", "root:root@tcp(localhost:13316)/webook")
	// viper 使用1
	// 读取的配置文件名
	viper.SetConfigName("dev")
	// 读取的配置文件类型
	viper.SetConfigType("yaml")
	// 读取配置文件的路径（可以是多个路径，会按照顺序读取，直到读取到配置文件）
	// 相对路径：不是当前文件夹下相对的路径，而是working director（工作目录）下的相对路径
	viper.AddConfigPath("./config")

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)

	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// 加载配置文件方式2
func InitViperV1() {
	viper.SetConfigFile("./config/dev.yaml")
	// 监听配置文件
	viper.WatchConfig()
	// 当有操作配置文件时，就会键入这里执行
	// Viper的缺点：
	// 1. 只能监听什么配置文件修改了，但是没有监听配置文件中哪些key修改了
	// 2. 无法知道修改key 修改前 和 修改后 的值
	// 3. viper监听配置文件，是无法做到监听远程配置中心的配置修改
	viper.OnConfigChange(func(in fsnotify.Event) {
		// in.Name: 配置文件的路径
		// in.Op: 配置文件的操作类型（创建/删除/写入/重命名）
		fmt.Print(in.Name, in.Op)
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// 直接读取
func InitViperReader() {
	cfg := `
db.config:
  dsn: "root:root@tcp(localhost:13316)/webook"

redis:
  addr: "localhost:6379"
`
	// 必须设置配置类型，才能解析
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

// 不同环境下加载不同的配置文件
func InitViperV2() {
	// 开发环境/测试环境/线上环境

	// pflag: 会加载我们指令中或goland配置中的配置文件，但是如果没有设置的话，他就会以默认值加载配置文件
	// name: config -> --config=config/config.yaml（命令行或goland配置中路径）
	cfile := pflag.String("config", "./config/config.yaml", "默认配置文件路路径")
	pflag.Parse()
	// 一定要先解析出来，才能加载对应文件（若还未解析出路径，viper就会以默认值加载配置文件）
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// 远程配置中心
func InitViperRemote() {

	err := viper.AddRemoteProvider("etcd3", "localhost:12379",
		// 通过 webook 和 其它使用etcd的区分开来
		"/webook")

	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func InitLogger() {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	// 如果不调用ReplaceGlobals，直接调用zap.L()，什么都打不开
	zap.ReplaceGlobals(logger)
	zap.L().Info("初始化日志成功",
		zap.Error(err),
		zap.String("sss", "ssss"),
		zap.Int32("code", 200))
}
