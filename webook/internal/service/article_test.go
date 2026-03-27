package service

import (
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/article"
	repomocks "github.com/bbbbbbbbiao/WeBook/webook/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
 * @author: biao
 * @date: 2026/3/21 下午3:54
 * @description:
 */

func TestArticleService_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (article.ArticleAuthorRepo, article.ArticleReaderRepo)
		// 输入
		wantReq domain.Article
		// 输出
		wantId  int64
		wantErr error
	}{
		{
			name: "新建文章-发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepo, article.ArticleReaderRepo) {
				authorRepo := repomocks.NewMockArticleAuthorRepo(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepo(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "新建标题",
					Content: "新建内容",
					Author: domain.Author{
						AuthorId: 123,
					},
				}).Return(int64(2), nil)
				readerRepo.EXPECT().Upsert(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "新建标题",
					Content: "新建内容",
					Author: domain.Author{
						AuthorId: 123,
					},
				}).Return(int64(2), nil)
				return authorRepo, readerRepo
			},
			wantReq: domain.Article{
				Title:   "新建标题",
				Content: "新建内容",
				Author: domain.Author{
					AuthorId: 123,
				},
			},
			wantId:  2,
			wantErr: nil,
		},
		{
			name: "修改文章-发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepo, article.ArticleReaderRepo) {
				authorRepo := repomocks.NewMockArticleAuthorRepo(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepo(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      33,
					Title:   "修改标题",
					Content: "修改内容",
					Author: domain.Author{
						AuthorId: 123,
					},
				}).Return(int64(33), nil)
				readerRepo.EXPECT().Upsert(gomock.Any(), domain.Article{
					Id:      33,
					Title:   "修改标题",
					Content: "修改内容",
					Author: domain.Author{
						AuthorId: 123,
					},
				}).Return(int64(33), nil)
				return authorRepo, readerRepo
			},
			wantReq: domain.Article{
				Id:      33,
				Title:   "修改标题",
				Content: "修改内容",
				Author: domain.Author{
					AuthorId: 123,
				},
			},
			wantId:  33,
			wantErr: nil,
		},
		{
			name: "修改文章-保存到制作库失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepo, article.ArticleReaderRepo) {
				authorRepo := repomocks.NewMockArticleAuthorRepo(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepo(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      33,
					Title:   "修改标题",
					Content: "修改内容",
					Author: domain.Author{
						AuthorId: 123,
					},
				}).Return(int64(0), errors.New("保存到制作库失败"))
				return authorRepo, readerRepo
			},
			wantReq: domain.Article{
				Id:      33,
				Title:   "修改标题",
				Content: "修改内容",
				Author: domain.Author{
					AuthorId: 123,
				},
			},
			wantId:  0,
			wantErr: errors.New("保存到制作库失败"),
		},
		{
			name: "修改文章-保存到线上库失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepo, article.ArticleReaderRepo) {
				authorRepo := repomocks.NewMockArticleAuthorRepo(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepo(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      33,
					Title:   "修改标题",
					Content: "修改内容",
					Author: domain.Author{
						AuthorId: 123,
					},
				}).Return(int64(33), nil)
				readerRepo.EXPECT().Upsert(gomock.Any(), domain.Article{
					Id:      33,
					Title:   "修改标题",
					Content: "修改内容",
					Author: domain.Author{
						AuthorId: 123,
					},
				}).Return(int64(0), errors.New("保存到制作库失败"))
				return authorRepo, readerRepo
			},
			wantReq: domain.Article{
				Id:      33,
				Title:   "修改标题",
				Content: "修改内容",
				Author: domain.Author{
					AuthorId: 123,
				},
			},
			wantId:  0,
			wantErr: errors.New("保存到制作库失败"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			artServiceImpl := NewArticleServiceImplV1(tc.mock(ctrl))
			id, err := artServiceImpl.Publish(nil, tc.wantReq)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)
		})
	}
}
