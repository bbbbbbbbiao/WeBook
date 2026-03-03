package web

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	ws  wechat.Service
	svc service.UserService
	JwtHandler
}

func NewOAuth2WechatHandler(ws wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		ws: ws,
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	oAuth2 := server.Group("/oAuth2/wechat")
	oAuth2.GET("/authUrl", o.AuthUrl)
	oAuth2.Any("/callBack", o.CallBack)
}

func (o *OAuth2WechatHandler) AuthUrl(ctx *gin.Context) {
	url, err := o.ws.AuthUrl(ctx)
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

func (o *OAuth2WechatHandler) CallBack(ctx *gin.Context) {
	code := ctx.Query("code")
	status := ctx.Query("status")
	info, err := o.ws.VerifyCode(ctx, code, status)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  err.Error(),
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
	err = o.SetJWTToken(ctx, u.Id)
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
