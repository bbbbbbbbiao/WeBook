package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service"
	svcmocks "github.com/bbbbbbbbiao/WeBook/webook/internal/service/mocks"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web/jwt"
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**
 * @author: biao
 * @date: 2026/3/21 上午11:18
 * @description:
 */

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) service.ArticleService
		wantReq articleReq

		wantCode int
		wantResp Result
	}{
		{
			name: "新建文章-发表成功",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleSvc := svcmocks.NewMockArticleService(ctrl)
				articleSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "新建文章",
					Content: "新建内容",
					Author: domain.Author{
						AuthorId: 123,
					},
				}).Return(int64(1), nil)
				return articleSvc
			},
			wantReq: articleReq{
				Title:   "新建文章",
				Content: "新建内容",
			},
			wantCode: http.StatusOK,
			wantResp: Result{
				Data: float64(1),
				Msg:  "发表成功",
			},
		},
		{
			name: "修改文章-发表成功",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleSvc := svcmocks.NewMockArticleService(ctrl)
				articleSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      3,
					Title:   "修改文章",
					Content: "修改内容",
					Author: domain.Author{
						AuthorId: 123,
					},
				}).Return(int64(3), nil)
				return articleSvc
			},
			wantReq: articleReq{
				Id:      3,
				Title:   "修改文章",
				Content: "修改内容",
			},
			wantCode: http.StatusOK,
			wantResp: Result{
				Data: float64(3),
				Msg:  "发表成功",
			},
		},
		{
			name: "发表失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleSvc := svcmocks.NewMockArticleService(ctrl)
				articleSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      3,
					Title:   "修改文章",
					Content: "修改内容",
					Author: domain.Author{
						AuthorId: 123,
					},
				}).Return(int64(0), errors.New("发表失败"))
				return articleSvc
			},
			wantReq: articleReq{
				Id:      3,
				Title:   "修改文章",
				Content: "修改内容",
			},
			wantCode: http.StatusOK,
			wantResp: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := NewArticleHandler(logger.NewNopLogger(), tc.mock(ctrl))
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("userClaims", &jwt.AccessClaims{
					UserId: 123,
				})
			})
			handler.RegisterRoutes(server)

			// 构建请求
			reqBody, err := json.Marshal(tc.wantReq)
			assert.NoError(t, err)
			req, err := http.NewRequest("POST", "/articles/publish", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// 执行
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			// 验证响应
			assert.Equal(t, tc.wantCode, resp.Code)
			var res Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResp, res)
		})
	}
}

type articleReq struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorId int64  `json:"author_id"`
}
