package integration

import (
	"bytes"
	"encoding/json"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/integration/startup"
	aa "github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao/article"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**
 * @author: biao
 * @date: 2026/3/17 下午2:11
 * @description: 文章集成测试
 */

// TDD：Test Driver Develop 测试驱动开发

// 用组合类型的测试套件
// 引入测试套件，避免一些重复操作
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB // 目的是after中验证数据是否正确
}

// SetupSuite 这个会在整个测试之前进行调用初始化
func (a *ArticleTestSuite) SetupSuite() {
	a.server = gin.Default()
	// 模拟登录
	a.server.Use(func(ctx *gin.Context) {
		ctx.Set("userClaims", &jwt.AccessClaims{
			UserId: 123,
		})
	})
	handler := startup.InitArticleHandler()
	handler.RegisterRoutes(a.server)
	a.db = startup.InitDB()
}

// 每个测试之后调用
func (a *ArticleTestSuite) TearDownSuite() {
	a.db.Exec("TRUNCATE TABLE articles")
	a.db.Exec("TRUNCATE TABLE reader_articles")
}

func (a *ArticleTestSuite) TestArticle_Edit() {
	t := a.T()
	testCases := []struct {
		name   string
		before func(t *testing.T) // 准备数据
		after  func(t *testing.T) // 验证数据、清理数据，需要跳过我们写的代码

		wantReq Article

		wantCode int // http状态码
		wantRes  Result[int64]
	}{
		{
			name: "新建文章-保存成功",
			before: func(t *testing.T) {
				// 因为获取的id是自增，所以每次都需要把数据删掉
				//a.db.Exec("TRUNCATE TABLE articles")
			},
			after: func(t *testing.T) {
				// 验证数据：验证数据库数据，所以需要跳过我们写的代码
				var article aa.Article
				err := a.db.Where("author_id = ?", 123).First(&article).Error
				assert.NoError(t, err)
				// 因为验证不了时间等于当前时间，所以验证其大于0即可
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)
				// 验证完后设置为 0
				article.Ctime = 0
				article.Utime = 0
				assert.Equal(t, article, aa.Article{
					Id:       1,
					Title:    "文章标题",
					Content:  "文章内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnPublished.ToUint8(),
				})
			},

			wantReq: Article{
				Title:   "文章标题",
				Content: "文章内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 1,
			},
		},
		{
			name: "修改文章-保存成功",
			before: func(t *testing.T) {
				// 之前需要先新建文章
				err := a.db.Create(&aa.Article{
					Id:       2,
					Title:    "文章标题",
					Content:  "文章内容",
					Status:   domain.ArticleStatusUnPublished.ToUint8(),
					AuthorId: 123,
					Ctime:    123,
					Utime:    456,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据：验证数据库数据，所以需要跳过我们写的代码
				var article aa.Article
				err := a.db.Where("id = ?", 2).First(&article).Error
				assert.NoError(t, err)
				// 因为修改了，所以修改时间看到比现在时间大
				assert.True(t, article.Utime > 456)
				// 验证完后设置为 0
				article.Utime = 0
				assert.Equal(t, article, aa.Article{
					Id:       2,
					Title:    "修改标题",
					Content:  "修改内容",
					Status:   domain.ArticleStatusUnPublished.ToUint8(),
					AuthorId: 123,
					Ctime:    123,
				})
			},

			wantReq: Article{
				Id:      2,
				Title:   "修改标题",
				Content: "修改内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 2,
			},
		},
		{
			name: "修改别人文章-保存失败",
			before: func(t *testing.T) {
				a.db.Create(&aa.Article{
					Id:       5,
					Title:    "文章标题",
					Content:  "文章内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnPublished.ToUint8(),
					Ctime:    123,
					Utime:    123,
				})
			},
			after: func(t *testing.T) {
				// 验证数据：验证数据库数据，所以需要跳过我们写的代码
				var article aa.Article
				err := a.db.Where("id = ?", 5).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Utime == 123)
				// 验证完后设置为 0
				article.Utime = 0
				assert.Equal(t, article, aa.Article{
					Id:       5,
					Title:    "文章标题",
					Content:  "文章内容",
					Status:   domain.ArticleStatusUnPublished.ToUint8(),
					AuthorId: 123,
					Ctime:    123,
				})
			},

			wantReq: Article{
				Id:      6,
				Title:   "修改标题",
				Content: "修改内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "系统错误",
				Code: 5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 构造请求
			tc.before(t)

			wantReq, err := json.Marshal(tc.wantReq)
			assert.NoError(t, err)
			req, err := http.NewRequest("POST", "/articles/edit", bytes.NewBuffer(wantReq))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			// 执行

			a.server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			// 验证响应
			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)

			tc.after(t)
		})
	}
}

func (a *ArticleTestSuite) TestArticle_Publish() {
	t := a.T()
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		wantReq  Article
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "新建文章-发表成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var article aa.Article
				err := a.db.Where("author_id = ?", 123).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)
				article.Ctime = 0
				article.Utime = 0
				assert.Equal(t, article, aa.Article{
					Id:       1,
					Title:    "新建文章标题",
					Content:  "新建文章内容",
					Status:   domain.ArticleStatusPublished.ToUint8(),
					AuthorId: 123,
				})
				var readArticle aa.ReaderArticle
				err = a.db.Where("id = ?", 1).First(&readArticle).Error
				assert.NoError(t, err)
				assert.True(t, readArticle.Ctime > 0)
				assert.True(t, readArticle.Utime > 0)
				readArticle.Ctime = 0
				readArticle.Utime = 0
				assert.Equal(t, readArticle, aa.ReaderArticle{
					Article: aa.Article{
						Id:       1,
						Title:    "新建文章标题",
						Content:  "新建文章内容",
						Status:   domain.ArticleStatusPublished.ToUint8(),
						AuthorId: 123,
					},
				})
			},
			wantReq: Article{
				Title:   "新建文章标题",
				Content: "新建文章内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "发表成功",
				Data: 1,
			},
		},
		{
			name: "修改文章-发表成功（制作库，线上库都有）",
			before: func(t *testing.T) {
				var article = aa.Article{
					Id:       3,
					Title:    "新建文章标题",
					Content:  "新建文章内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
					Ctime:    111,
					Utime:    222,
				}
				var readArticle = aa.ReaderArticle{
					Article: article,
				}
				a.db.Create(&article)
				a.db.Create(&readArticle)
			},
			after: func(t *testing.T) {
				var article aa.Article
				err := a.db.Where("id = ?", 3).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Utime > 222)
				article.Utime = 0
				assert.Equal(t, article, aa.Article{
					Id:       3,
					Title:    "修改文章标题",
					Content:  "修改文章内容",
					Status:   domain.ArticleStatusPublished.ToUint8(),
					AuthorId: 123,
					Ctime:    111,
				})
				var readArticle aa.ReaderArticle
				err = a.db.Where("id = ?", 3).First(&readArticle).Error
				assert.NoError(t, err)
				assert.True(t, readArticle.Utime > 222)
				readArticle.Utime = 0
				assert.Equal(t, readArticle, aa.ReaderArticle{
					Article: aa.Article{
						Id:       3,
						Title:    "修改文章标题",
						Content:  "修改文章内容",
						Status:   domain.ArticleStatusPublished.ToUint8(),
						AuthorId: 123,
						Ctime:    111,
					},
				})
			},
			wantReq: Article{
				Id:      3,
				Title:   "修改文章标题",
				Content: "修改文章内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "发表成功",
				Data: 3,
			},
		},
		{
			name: "修改文章-发表成功（只有制作库有）",
			before: func(t *testing.T) {
				var article = aa.Article{
					Id:       3,
					Title:    "新建文章标题",
					Content:  "新建文章内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
					Ctime:    111,
					Utime:    222,
				}
				a.db.Create(&article)
			},
			after: func(t *testing.T) {
				var article aa.Article
				err := a.db.Where("id = ?", 3).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Utime > 222)
				article.Utime = 0
				assert.Equal(t, article, aa.Article{
					Id:       3,
					Title:    "修改文章标题",
					Content:  "修改文章内容",
					Status:   domain.ArticleStatusPublished.ToUint8(),
					AuthorId: 123,
					Ctime:    111,
				})
				var readArticle aa.ReaderArticle
				err = a.db.Where("id = ?", 3).First(&readArticle).Error
				assert.NoError(t, err)
				assert.True(t, readArticle.Ctime > 0)
				assert.True(t, readArticle.Utime > 0)
				readArticle.Ctime = 0
				readArticle.Utime = 0
				assert.Equal(t, readArticle, aa.ReaderArticle{
					Article: aa.Article{
						Id:       3,
						Title:    "修改文章标题",
						Content:  "修改文章内容",
						Status:   domain.ArticleStatusPublished.ToUint8(),
						AuthorId: 123,
					},
				})
			},
			wantReq: Article{
				Id:      3,
				Title:   "修改文章标题",
				Content: "修改文章内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "发表成功",
				Data: 3,
			},
		},
		{
			name: "发表文章-失败",
			before: func(t *testing.T) {
				var article = aa.Article{
					Id:       888,
					Title:    "新建文章标题",
					Content:  "新建文章内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
					Ctime:    111,
					Utime:    222,
				}
				var readArticle = aa.ReaderArticle{
					Article: article,
				}
				a.db.Create(&article)
				a.db.Create(&readArticle)
			},
			after: func(t *testing.T) {
				var article aa.Article
				err := a.db.Where("id = ?", 888).First(&article).Error
				assert.NoError(t, err)
				assert.Equal(t, article, aa.Article{
					Id:       888,
					Title:    "新建文章标题",
					Content:  "新建文章内容",
					Status:   domain.ArticleStatusPublished.ToUint8(),
					AuthorId: 123,
					Ctime:    111,
					Utime:    222,
				})
				var readArticle aa.ReaderArticle
				err = a.db.Where("id = ?", 888).First(&readArticle).Error
				assert.NoError(t, err)
				assert.Equal(t, readArticle, aa.ReaderArticle{
					Article: aa.Article{
						Id:       888,
						Title:    "新建文章标题",
						Content:  "新建文章内容",
						Status:   domain.ArticleStatusPublished.ToUint8(),
						AuthorId: 123,
						Ctime:    111,
						Utime:    222,
					},
				})
			},
			wantReq: Article{
				Id:      5,
				Title:   "修改文章标题",
				Content: "修改文章内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "系统错误",
				Code: 5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			wantReq, err := json.Marshal(tc.wantReq)
			assert.NoError(t, err)
			req, err := http.NewRequest("POST", "/articles/publish", bytes.NewBuffer(wantReq))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			resp := httptest.NewRecorder()
			a.server.ServeHTTP(resp, req)

			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantRes, res)
			tc.after(t)
		})
	}
}

func (a *ArticleTestSuite) TestArticle_Withdraw() {
	t := a.T()
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		wantReq  Article
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "修改文章为不可见-成功",
			before: func(t *testing.T) {
				article := aa.Article{
					Id:       100,
					Title:    "新建文章标题",
					Content:  "新建文章内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
					Ctime:    111,
					Utime:    222,
				}
				var readArticle = aa.ReaderArticle{
					Article: article,
				}
				a.db.Create(&article)
				a.db.Create(&readArticle)
			},
			after: func(t *testing.T) {
				var article aa.Article
				err := a.db.Where("id = ? and author_id = ?", 100, 123).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Utime > 222)
				article.Utime = 0
				assert.Equal(t, article, aa.Article{
					Id:       100,
					Title:    "新建文章标题",
					Content:  "新建文章内容",
					Status:   domain.ArticleStatusPrivate.ToUint8(),
					AuthorId: 123,
					Ctime:    111,
				})
				var readArticle aa.ReaderArticle
				err = a.db.Where("id = ? and author_id = ?", 100, 123).First(&readArticle).Error
				assert.NoError(t, err)
				assert.True(t, readArticle.Utime > 222)
				readArticle.Utime = 0
				assert.Equal(t, readArticle, aa.ReaderArticle{
					Article: aa.Article{
						Id:       100,
						Title:    "新建文章标题",
						Content:  "新建文章内容",
						Status:   domain.ArticleStatusPrivate.ToUint8(),
						AuthorId: 123,
						Ctime:    111,
					},
				})
			},
			wantReq: Article{
				Id: 100,
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg: "ok",
			},
		},
		{
			name: "修改文章为不可见-失败",
			before: func(t *testing.T) {
				article := aa.Article{
					Id:       200,
					Title:    "新建文章标题",
					Content:  "新建文章内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
					Ctime:    111,
					Utime:    222,
				}
				var readArticle = aa.ReaderArticle{
					Article: article,
				}
				a.db.Create(&article)
				a.db.Create(&readArticle)
			},
			after: func(t *testing.T) {
				var article aa.Article
				err := a.db.Where("id = ? and author_id = ?", 200, 123).First(&article).Error
				assert.NoError(t, err)
				assert.Equal(t, article, aa.Article{
					Id:       200,
					Title:    "新建文章标题",
					Content:  "新建文章内容",
					Status:   domain.ArticleStatusPublished.ToUint8(),
					AuthorId: 123,
					Ctime:    111,
					Utime:    222,
				})
				var readArticle aa.ReaderArticle
				err = a.db.Where("id = ? and author_id = ?", 200, 123).First(&readArticle).Error
				assert.NoError(t, err)
				assert.Equal(t, readArticle, aa.ReaderArticle{
					Article: aa.Article{
						Id:       200,
						Title:    "新建文章标题",
						Content:  "新建文章内容",
						Status:   domain.ArticleStatusPublished.ToUint8(),
						AuthorId: 123,
						Ctime:    111,
						Utime:    222,
					},
				})
			},
			wantReq: Article{
				Id: 500,
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "系统错误",
				Code: 5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			wantReq, err := json.Marshal(tc.wantReq)
			assert.NoError(t, err)
			req, err := http.NewRequest("POST", "/articles/withdraw", bytes.NewBuffer(wantReq))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			resp := httptest.NewRecorder()
			a.server.ServeHTTP(resp, req)

			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantRes, res)
			tc.after(t)
		})
	}
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func TestArticleTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))

}
