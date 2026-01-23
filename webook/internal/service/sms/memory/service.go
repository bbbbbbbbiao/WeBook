package memory

import (
	"context"
	"fmt"
)

/**
 * @author: biao
 * @date: 2026/1/21 上午10:17
 * @description:
 */

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}
