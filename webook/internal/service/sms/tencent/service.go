package tencent

import (
	"context"
	"fmt"
	"github.com/bbbbbbbbiao/WeBook/webook/ekit"
	mySms "github.com/bbbbbbbbiao/WeBook/webook/internal/service/sms"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

/**
 * @author: biao
 * @date: 2026/1/18 下午4:15
 * @description:
 */

type Service struct {
	client   *sms.Client
	signName *string
	appId    *string
}

func NewService(client *sms.Client, signName string, appId string) *Service {
	return &Service{
		client:   client,
		signName: ekit.ToPtr[string](signName),
		appId:    ekit.ToPtr[string](appId),
	}
}

func (s *Service) SendV1(ctx context.Context, tpl string, args []mySms.NamedArg, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tpl)
	req.TemplateParamSet = s.ToNamedArgPtrSlice(args)
	req.PhoneNumberSet = s.ToStringPtrSlice(numbers)
	response, err := s.client.SendSms(req)

	if err != nil {
		return err
	}

	// 为啥要遍历
	// 因为会给多个手机发送短信，这些就是查看每个手机的短信是否发送成功
	for _, status := range response.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			// 有一条错误，则返回error
			return fmt.Errorf("发送失败, code: %s, 原因是: %s", *status.Code, *status.Message)
		}
	}

	return nil
}

func (s *Service) ToStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}

func (s *Service) ToNamedArgPtrSlice(src []mySms.NamedArg) []*string {
	return slice.Map[mySms.NamedArg, *string](src, func(idx int, src mySms.NamedArg) *string {
		return &src.Val
	})
}
