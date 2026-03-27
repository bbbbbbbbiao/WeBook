package article

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

/**
 * @author: biao
 * @date: 2026/3/21 下午3:56
 * @description:
 */

// 线上表
type ReaderArticle struct {
	Article
}

type ArticleReaderDao interface {
	Upsert(ctx context.Context, article ReaderArticle) (int64, error)
}

type ArticleReaderDaoImpl struct {
	db *gorm.DB
}

func NewArticleReaderDaoImpl(db *gorm.DB) ArticleReaderDao {
	return &ArticleReaderDaoImpl{
		db: db,
	}
}

func (a *ArticleReaderDaoImpl) Upsert(ctx context.Context, article ReaderArticle) (int64, error) {
	// 新建冲突是，则进行修改
	err := a.db.WithContext(ctx).Clauses(clause.OnConflict{
		// mysql只关心这个
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   article.Title,
			"content": article.Content,
			"Utime":   time.Now(),
		}),
	}).Create(&article).Error
	if err != nil {
		return 0, err
	}
	return article.Id, nil
}
