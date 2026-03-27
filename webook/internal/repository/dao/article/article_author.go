package article

import (
	"context"
	"gorm.io/gorm"
)

/**
 * @author: biao
 * @date: 2026/3/21 下午3:56
 * @description:
 */

type ArticleAuthorDao interface {
	Create(ctx context.Context, article Article) (int64, error)
	Update(ctx context.Context, article Article) (int64, error)
}

type ArticleAuthorDaoImpl struct {
	db *gorm.DB
}

func NewArticleAuthorDaoImpl(db *gorm.DB) ArticleAuthorDao {
	return &ArticleAuthorDaoImpl{
		db: db,
	}
}

func (a *ArticleAuthorDaoImpl) Create(ctx context.Context, article Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *ArticleAuthorDaoImpl) Update(ctx context.Context, article Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}
