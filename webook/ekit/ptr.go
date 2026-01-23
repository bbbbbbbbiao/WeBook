package ekit

/**
 * @author: biao
 * @date: 2026/1/18 下午4:29
 * @description: 转换为指针泛型工具类
 */

func ToPtr[T any](t T) *T {
	return &t
}
