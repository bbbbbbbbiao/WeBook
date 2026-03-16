package web

import (
	"fmt"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/oauth2/wechat"
	iJwt "github.com/bbbbbbbbiao/WeBook/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	ws  wechat.Service
	svc service.UserService
	iJwt.Handler
	stateKey []byte
	cfg      WeChatHandlerConfig
}

type WeChatHandlerConfig struct {
	Secure bool
}

func NewOAuth2WechatHandler(ws wechat.Service, svc service.UserService, cfg WeChatHandlerConfig, handler iJwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		ws:       ws,
		svc:      svc,
		stateKey: []byte("3E7QYaUxM5tMhDWwd5dfdshdYWND7WR2Vr"),
		cfg:      cfg,
		Handler:  handler,
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	oAuth2 := server.Group("/oAuth2/wechat")
	oAuth2.GET("/authUrl", o.AuthUrl)
	oAuth2.Any("/callBack", o.CallBack)
}

func (o *OAuth2WechatHandler) AuthUrl(ctx *gin.Context) {
	// 生成state（随机数）
	state := uuid.New().String()
	err := o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}

	url, err := o.ws.AuthUrl(ctx, state)

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "生成微信扫码连接错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Date: url,
	})
}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, WeChatClaim{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenStr, err := token.SignedString(o.stateKey)
	if err != nil {
		return err
	}
	// 将token放入到cookie中
	ctx.SetCookie("jwt-state", tokenStr,
		600,                       // 600秒
		"/oAuth2/wechat/callBack", // 作用域为回调地址（其余没有这个cookie）
		"",
		o.cfg.Secure,
		true)
	return nil
}

func (o *OAuth2WechatHandler) CallBack(ctx *gin.Context) {
	code := ctx.Query("code")

	err := o.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	info, err := o.ws.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	var wechatInfo = domain.WeChatInfo{
		OpenId:  info.OpenId,
		UnionId: info.UnionId,
	}
	u, err := o.svc.FindOrCreateByWechat(ctx, wechatInfo)

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	err = o.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "登陆成功",
	})
}

func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")

	// 先验证state
	tokenStr, err := ctx.Cookie("jwt-state")
	if err != nil || tokenStr == "" {

		return fmt.Errorf("拿不到state的cookie, %w", err)
	}
	var weChatClaim WeChatClaim
	token, err := jwt.ParseWithClaims(tokenStr, &weChatClaim, func(token *jwt.Token) (interface{}, error) {
		return o.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("state的cookie解析错误, %w", err)
	}
	// 验证state是否一致
	if weChatClaim.State != state {
		return fmt.Errorf("state不一致, %w", err)
	}
	return nil
}

type WeChatClaim struct {
	jwt.RegisteredClaims
	State string
}
