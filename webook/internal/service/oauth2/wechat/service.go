package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// 用于处理Path地址的双引号问题（fmt.Sprintf）
var redirectUri = url.PathEscape("/oAuth2/wechat/calBack")

type Service interface {
	AuthUrl(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (WechatResult, error)
}

type OAuth2WechatService struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewOAuth2WechatService(appId string, appSecret string) Service {
	return &OAuth2WechatService{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient, // 偷懒写法
	}
}

func (o *OAuth2WechatService) AuthUrl(ctx context.Context, state string) (string, error) {
	var urlPattern = " https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprintf(urlPattern, o.appId, redirectUri, state), nil
}

func (o *OAuth2WechatService) VerifyCode(ctx context.Context, code string) (WechatResult, error) {
	var redirectUrlPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	redirectUrl := fmt.Sprintf(redirectUrlPattern, o.appId, o.appSecret, code)
	//resp, err := http.Get(redirectUrl)
	req, err := http.NewRequestWithContext(ctx, "GET", redirectUrl, nil)
	resp, err := o.client.Do(req)

	if err != nil {
		return WechatResult{}, err
	}

	defer resp.Body.Close()
	var wechatInfo WechatResult
	err = json.NewDecoder(resp.Body).Decode(&wechatInfo)
	if err != nil {
		return WechatResult{}, err
	}

	if wechatInfo.ErrCode != 0 {
		return WechatResult{},
			fmt.Errorf("微信返回错误响应，错误码：%d，错误信息：%s", wechatInfo.ErrCode, wechatInfo.ErrMsg)
	}

	return wechatInfo, nil
}

type WechatResult struct {
	AccessToken  string `json:"access_token"`
	ExpireIn     int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenId  string `json:"openid"`
	Scope   string `json:"scope"`
	UnionId string `json:"unionid"`

	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
