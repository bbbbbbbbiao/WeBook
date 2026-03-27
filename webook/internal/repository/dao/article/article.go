package article

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

/**
 * @author: biao
 * @date: 2026/3/18 下午7:59
 * @description:
 */

var ErrorAuthorIdNotEqual = errors.New("作者id不相同")

// 制作表
type Article struct {
	Id      int64  `gorm:"primary_key,AUTO_INCREMENT"`
	Title   string `gorm:"type=varchar(64)"`
	Content string `gorm:"type=BLOB"`
	// select * from article where AuthorId = 1223 order by Ctime DESC
	// 是否需要在这里给 AuthorId 和 Ctime 建立联合索引呢？
	AuthorId int64 `gorm:"index"`
	Status   uint8
	Ctime    int64
	Utime    int64
}

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	Update(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, article ReaderArticle) (int64, error)
	Sync(ctx context.Context, article Article) (int64, error)
	Withdraw(ctx *gin.Context, id int64, userId int64, status uint8) error // 在dao层面开启事务
}

type ArticleDaoImpl struct {
	db *gorm.DB
}

func NewArticleDaoImpl(db *gorm.DB) ArticleDao {
	return &ArticleDaoImpl{
		db: db,
	}
}

func (a *ArticleDaoImpl) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := a.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

func (a *ArticleDaoImpl) Update(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Utime = now
	// 受影响的条数
	affected := a.db.WithContext(ctx).Model(&article).Where("id = ? and author_id = ?", article.Id, article.AuthorId).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"utime":   article.Utime,
		}).RowsAffected

	if affected == 0 {
		return article.Id, ErrorAuthorIdNotEqual
	}
	return article.Id, nil
}

func (a *ArticleDaoImpl) Upsert(ctx context.Context, article ReaderArticle) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := a.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   article.Title,
			"content": article.Content,
			"utime":   time.Now().UnixMilli(),
			"status":  article.Status,
		}),
	}).Create(&article).Error
	return article.Id, err
}

// 在dao层面开启事务
func (a *ArticleDaoImpl) Sync(ctx context.Context, article Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	// 利用GORM事务的闭包
	err = a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		articleImpl := NewArticleDaoImpl(tx)

		if id > 0 {
			id, err = articleImpl.Update(ctx, article)
		} else {
			id, err = articleImpl.Insert(ctx, article)
		}
		if err != nil {
			return err
		}

		var readerArticle = ReaderArticle{
			Article: article,
		}
		id, err = articleImpl.Upsert(ctx, readerArticle)
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (a *ArticleDaoImpl) Withdraw(ctx *gin.Context, id int64, userId int64, status uint8) error {
	err := a.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		affected := tx.WithContext(ctx).
			Model(&Article{}).
			Where("id = ? and author_id = ?", id, userId).
			Updates(map[string]interface{}{
				"status": status,
				"utime":  now,
			}).RowsAffected
		if affected == 0 {
			return ErrorAuthorIdNotEqual
		}

		err := tx.WithContext(ctx).
			Model(&ReaderArticle{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"status": status,
				"utime":  now,
			}).Error
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
