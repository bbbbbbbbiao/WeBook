package web

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service"
	iJwt "github.com/bbbbbbbbiao/WeBook/webook/internal/web/jwt"
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @author: biao
 * @date: 2026/3/16 下午9:13
 * @description:
 */

type ArticleHandler struct {
	l   logger.LoggerV2
	svc service.ArticleService
}

func NewArticleHandler(l logger.LoggerV2, svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		l:   l,
		svc: svc,
	}
}

func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)
}

// 新建或修改并保存（制作库）
func (a *ArticleHandler) Edit(ctx *gin.Context) {

	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("参数错误")
		return
	}

	claims, _ := ctx.Get("userClaims")
	accessClaims := claims.(*iJwt.AccessClaims)
	// 新增或修改的保存
	id, err := a.svc.Save(ctx, req.ReqToDomain(accessClaims.UserId))

	if err == service.ErrorAuthorIdNotEqual {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("文章保存失败", logger.Error(err))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("文章保存失败", logger.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg:  "ok",
		Data: id,
	})

}

// 发表
// 1. 新建并发表
// 2. 修改并发表（制作库有，线上库没有）
// 3. 修改并发表（制作可和线上库都有）
func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("参数错误")
		return
	}
	claims, _ := ctx.Get("userClaims")
	accessClaims := claims.(*iJwt.AccessClaims)
	id, err := a.svc.Publish(ctx, req.ReqToDomain(accessClaims.UserId))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("文章发表失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "发表成功",
		Data: id,
	})
}

// 修改发表文章的状态
func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("参数错误", logger.Error(err))
		return
	}
	claims, _ := ctx.Get("userClaims")
	accessClaims := claims.(*iJwt.AccessClaims)
	err := a.svc.Withdraw(ctx, req.Id, accessClaims.UserId)

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("隐藏文章失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  uint8  `json:"status"`
}

func (a ArticleReq) ReqToDomain(authorId int64) domain.Article {
	return domain.Article{
		Id:      a.Id,
		Title:   a.Title,
		Content: a.Content,
		Author: domain.Author{
			AuthorId: authorId,
		},
	}
}
