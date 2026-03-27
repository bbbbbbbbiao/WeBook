package article

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
)

/**
 * @author: biao
 * @date: 2026/3/21 下午3:57
 * @description:
 */

type ArticleReaderRepo interface {
	Upsert(ctx context.Context, article domain.Article) (int64, error)
}
