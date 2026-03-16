package logger

/**
 * @author: biao
 * @date: 2026/3/14 下午6:51
 * @description:
 */

// 当不再使用zap作为日志框架时，或者需要灵活替换时，则可以使用接口进行封装
// 兼容性最好：第一种
// 参数需要有名字：第二种
// 有严格的代码评审流程：第三种

// 参数采字符串形式拼接，拼接方式由外部传进来
type LoggerV1 interface {
	Debug(Msg string, args ...any)
	Info(Msg string, args ...any)
	Warn(Msg string, args ...any)
	Error(Msg string, args ...any)
}

// 参数采用键值对形式拼接
type LoggerV2 interface {
	Debug(Msg string, args ...Field)
	Info(Msg string, args ...Field)
	Warn(Msg string, args ...Field)
	Error(Msg string, args ...Field)
}

type Field struct {
	Key   string
	Value any
}

// 参数必须是偶数个，组成键值对
type LoggerV3 interface {
	Debug(Msg string, args ...any)
	Info(Msg string, args ...any)
	Warn(Msg string, args ...any)
	Error(Msg string, args ...any)
}

func LoggerExample() {
	var l1 LoggerV1
	phone := "13800000000"
	l1.Info("手机号为 %s", phone)

	var l2 LoggerV2
	l2.Info("手机号错误", Field{
		Key:   "key",
		Value: phone,
	})

	var l3 LoggerV3
	l3.Info("手机号错误", "key", phone)

}
