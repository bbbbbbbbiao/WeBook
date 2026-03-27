package service

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/article"
	"github.com/bbbbbbbbiao/WeBook/webook/pkg/logger"
	"github.com/gin-gonic/gin"
)

/**
 * @author: biao
 * @date: 2026/3/18 下午3:47
 * @description:
 */

var ErrorAuthorIdNotEqual = article.ErrorAuthorIdNotEqual

type ArticleService interface {
	Save(ctx context.Context, artReq domain.Article) (int64, error)
	Publish(ctx *gin.Context, artReq domain.Article) (int64, error)
	Withdraw(ctx *gin.Context, id int64, userId int64) error
}

type ArticleServiceImpl struct {
	repo article.ArticleRepository
	l    logger.LoggerV2
}

func NewArticleServiceImpl(repo article.ArticleRepository) ArticleService {
	return &ArticleServiceImpl{
		repo: repo,
		l:    logger.NewNopLogger(),
	}
}

func (a *ArticleServiceImpl) Save(ctx context.Context, artReq domain.Article) (int64, error) {
	artReq.Status = domain.ArticleStatusUnPublished
	if artReq.Id > 0 {
		// 修改的保存
		return a.repo.Update(ctx, artReq)
	}
	// 新增的保存
	return a.repo.Create(ctx, artReq)
}

func (a *ArticleServiceImpl) Publish(ctx *gin.Context, artReq domain.Article) (int64, error) {
	// 文章发表成功，那么文章状态为"已发表"
	artReq.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, artReq)
}

func (a *ArticleServiceImpl) Withdraw(ctx *gin.Context, id int64, userId int64) error {
	return a.repo.Withdraw(ctx, id, userId, domain.ArticleStatusPrivate)
}

// 采用两个不同的repo实现
type ArticleServiceImplV1 struct {
	authorRepo article.ArticleAuthorRepo
	readerRepo article.ArticleReaderRepo
	l          logger.LoggerV2
}

func (a *ArticleServiceImplV1) Withdraw(ctx *gin.Context, id int64, userId int64) error {
	//TODO implement me
	panic("implement me")
}

func NewArticleServiceImplV1(
	authorRepo article.ArticleAuthorRepo,
	readerRepo article.ArticleReaderRepo) ArticleService {
	return &ArticleServiceImplV1{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
		l:          logger.NewNopLogger(),
	}
}

func (a *ArticleServiceImplV1) Save(ctx context.Context, artReq domain.Article) (int64, error) {
	// 不管是新增还是修改，状态都将变为"未发表状态"
	artReq.Status = domain.ArticleStatusUnPublished
	if artReq.Id > 0 {
		return a.authorRepo.Update(ctx, artReq)
	}
	return a.authorRepo.Create(ctx, artReq)
}

// 保存并发表（新建并发表、修改并发表）
// 我们可以通过id是否大于0，来判断制作库是新建还是更新
// 但是我们无法判断线上库是否已经有对应文章了，所以我们调用upsert（由dao层解决）
func (a *ArticleServiceImplV1) Publish(ctx *gin.Context, artReq domain.Article) (int64, error) {
	var (
		id  = artReq.Id
		err error
	)
	if artReq.Id > 0 {
		id, err = a.authorRepo.Update(ctx, artReq)
	} else {
		id, err = a.authorRepo.Create(ctx, artReq)
	}

	if err != nil {
		a.l.Error("发表到制作库失败", logger.Error(err))
		return 0, err
	}
	// 为了保证制作库和线上库一致性（使得他们的id相同）
	artReq.Id = id

	// 如果保存到制作库成功，保存到线上库失败，怎么办？
	// 高级：分布式事务（因为不知道他们在不在同一张库里面），消息队列，修复好后同步转异步处理，
	// 普通：写项目时一般重试机制（装饰器）
	return a.readerRepo.Upsert(ctx, artReq)
}
