package logger

import "go.uber.org/zap"

/**
 * @author: biao
 * @date: 2026/3/14 下午7:18
 * @description:
 */

// 适配器模式
type ZapLogger struct {
	zl *zap.Logger
}

func NewZapLogger(zl *zap.Logger) LoggerV2 {
	return &ZapLogger{
		zl: zl,
	}
}

func (z *ZapLogger) Debug(Msg string, args ...Field) {
	z.zl.Debug(Msg, z.ToZapFields(args)...)
}

func (z *ZapLogger) Info(Msg string, args ...Field) {
	z.zl.Info(Msg, z.ToZapFields(args)...)
}

func (z *ZapLogger) Warn(Msg string, args ...Field) {
	z.zl.Warn(Msg, z.ToZapFields(args)...)
}

func (z *ZapLogger) Error(Msg string, args ...Field) {
	z.zl.Error(Msg, z.ToZapFields(args)...)
}

func (z *ZapLogger) ToZapFields(args []Field) []zap.Field {
	fields := make([]zap.Field, 0, len(args))

	for _, arg := range args {
		fields = append(fields, zap.Any(arg.Key, arg.Value))
	}

	return fields
}
