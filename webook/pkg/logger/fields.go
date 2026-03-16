package logger

/**
 * @author: biao
 * @date: 2026/3/14 下午7:29
 * @description:
 */

func String(key string, value any) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
