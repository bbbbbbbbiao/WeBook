package article

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	aa "github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao/article"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/**
 * @author: biao
 * @date: 2026/3/18 下午7:48
 * @description:
 */

var ErrorAuthorIdNotEqual = aa.ErrorAuthorIdNotEqual

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)   // 没有引入事务，采用重试机制
	PublishV1(ctx context.Context, article domain.Article) (int64, error) // 在repository层引入事务（同库不同表）
	Sync(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx *gin.Context, id int64, userId int64, status domain.ArticleStatus) error // 在dao层面引入事务
}

type ArticleRepositoryImpl struct {
	dao       aa.ArticleDao
	authorDao aa.ArticleAuthorDao
	readerDao aa.ArticleReaderDao
	db        *gorm.DB
}

func NewArticleRepositoryImpl(dao aa.ArticleDao, db *gorm.DB) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao: dao,
		db:  db,
	}
}

func (a *ArticleRepositoryImpl) Create(ctx context.Context, article domain.Article) (int64, error) {
	return a.dao.Insert(ctx, a.DomainToEntity(article))
}

func (a *ArticleRepositoryImpl) Update(ctx context.Context, article domain.Article) (int64, error) {
	return a.dao.Update(ctx, a.DomainToEntity(article))
}

func (a *ArticleRepositoryImpl) Publish(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	if article.Id > 0 {
		id, err = a.authorDao.Update(ctx, a.DomainToEntity(article))
	} else {
		id, err = a.authorDao.Create(ctx, a.DomainToEntity(article))
	}
	if err != nil {
		return 0, err
	}
	article.Id = id
	var readerArt = aa.ReaderArticle{
		Article: a.DomainToEntity(article),
	}

	id, err = a.readerDao.Upsert(ctx, readerArt)
	return id, err
}

// repository层引入事务
func (a *ArticleRepositoryImpl) PublishV1(ctx context.Context, article domain.Article) (int64, error) {
	// 一个事务对应一个数据库连接
	// 当事务一直没有结束时，且没有超过数据对应的时间线时，数据库连接会被一直占用
	tx := a.db.Begin()
	defer tx.Rollback() // 若提交了，则这里会抛出异常，不管就行
	authorDaoImpl := aa.NewArticleAuthorDaoImpl(tx)
	readerDaoImpl := aa.NewArticleReaderDaoImpl(tx)
	var (
		id  = article.Id
		err error
	)
	if article.Id > 0 {
		id, err = authorDaoImpl.Update(ctx, a.DomainToEntity(article))
	} else {
		id, err = authorDaoImpl.Create(ctx, a.DomainToEntity(article))
	}
	if err != nil {
		//tx.Rollback() // 若每提交 defer 中会回滚
		return 0, err
	}
	article.Id = id
	var readerArt = aa.ReaderArticle{
		Article: a.DomainToEntity(article),
	}

	id, err = readerDaoImpl.Upsert(ctx, readerArt)
	if err != nil {
		//tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return id, nil
}

// dao层引入事务
func (a *ArticleRepositoryImpl) Sync(ctx context.Context, article domain.Article) (int64, error) {
	return a.dao.Sync(ctx, a.DomainToEntity(article))
}

func (a *ArticleRepositoryImpl) Withdraw(ctx *gin.Context, id int64, userId int64, status domain.ArticleStatus) error {
	return a.dao.Withdraw(ctx, id, userId, status.ToUint8())
}

func (a *ArticleRepositoryImpl) DomainToEntity(article domain.Article) aa.Article {
	return aa.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.AuthorId,
		Status:   article.Status.ToUint8(),
	}
}
