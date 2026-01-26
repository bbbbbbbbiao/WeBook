package ioc

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms/memory"
)

/**
 * @author: biao
 * @date: 2026/1/23 下午9:50
 * @description:
 */

func InitSMSService() sms.Service {
	return memory.NewService()
}
