package ioc

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

/**
 * @author: biao
 * @date: 2026/3/18 下午8:41
 * @description:
 */

func InitViper() {
	viper.SetConfigFile("E:\\Project\\my_project\\WeBook\\webook\\config\\dev.yaml")
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
