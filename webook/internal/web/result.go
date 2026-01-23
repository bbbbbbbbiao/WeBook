package web

/**
 * @author: biao
 * @date: 2026/1/21 下午10:26
 * @description:
 */

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Date any    `json:"date"`
}
