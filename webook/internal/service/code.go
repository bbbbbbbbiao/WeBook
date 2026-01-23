package service

import (
	"context"
	"fmt"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"math/rand"
)

/**
 * @author: biao
 * @date: 2026/1/18 下午7:41
 * @description:
 */

const codeTplId = "1877556"

var (
	ErrCodeSendTooMany  = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooMay = repository.ErrCodeVerifyTooMany
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// 发送验证码，发送给谁，以及区别业务场景
func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	// 三个步骤，生成一个验证码，放到redis中，发送验证码
	code := svc.GenerateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	if err != nil {
		// 这里有问题了
		// 这意味着，redis中有验证码了，但是不知道发出去了没（因为可能是因为超时问题）
	}

	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) GenerateCode() string {
	// 取值为 0 - 999999
	num := rand.Intn(10000000)
	return fmt.Sprintf("%06d", num)
}

//const codeTplId = "1877556"
//
//var ErrCodeSendTooMany = repository.ErrCodeSendTooMany
//
//type CodeService struct {
//	repo   *repository.CodeRepository
//	smsSvc sms.Service
//}
//
//func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
//	return &CodeService{
//		repo:   repo,
//		smsSvc: smsSvc,
//	}
//}
//
//func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
//	// 三件事：生成验证码、 存入Redis、 发送验证码
//	code := svc.GenerateCode()
//
//	err := svc.repo.Store(ctx, biz, phone, code)
//	if err != nil {
//		return err
//	}
//
//	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
//	// 这里有错误，存在了redis中，但是发送有问题，那么该删除redis中的值吗
//	// 不能，可能是因为发送超时了
//	return err
//}
//
//func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
//
//}
//
//func (svc *CodeService) GenerateCode() string {
//	intn := rand.Intn(1000000)
//	return fmt.Sprintf("%06d", intn)
//}
