package ioc

import (
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/logger"
	"go.uber.org/zap"
)

/**
 * @author: biao
 * @date: 2026/3/14 下午7:33
 * @description:
 */

func InitLogger() logger.LoggerV2 {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
