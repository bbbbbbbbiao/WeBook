package sms

import "context"

/**
 * @author: biao
 * @date: 2026/1/18 下午4:13
 * @description:
 */

type Service interface {
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
	//SendV1(ctx context.Context, tpl string, args []NamedArg, numbers ...string) error
}

// NamedArg 调用者类型为 []string 或者 map[string]string
type NamedArg struct {
	Val  string
	Name string
}
