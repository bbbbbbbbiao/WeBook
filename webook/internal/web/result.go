package web

/**
 * @author: biao
 * @date: 2026/1/21 下午10:26
 * @description:
 */

// 测试时，当以any接收值时
// 如果是数字，则以float64接收
// 如果是json，则以map接收
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
