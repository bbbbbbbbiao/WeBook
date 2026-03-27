package domain

/**
 * @author: biao
 * @date: 2026/3/18 下午3:48
 * @description:
 */

type Article struct {
	Id      int64         `json:"id"`
	Title   string        `json:"title"`
	Content string        `json:"content"`
	Author  Author        `json:"author"`
	Status  ArticleStatus `json:"status"`

	Ctime int64 `json:"ctime"`
	Utime int64 `json:"utime"`
}

const (
	ArticleStatusUnKnow ArticleStatus = iota
	ArticleStatusPublished
	ArticleStatusUnPublished
	ArticleStatusPrivate
)

type ArticleStatus uint8

func (a ArticleStatus) ToUint8() uint8 {
	return uint8(a)
}

func (a ArticleStatus) GetName() string {
	switch a {
	case ArticleStatusPublished:
		return "published"
	case ArticleStatusUnPublished:
		return "unpublished"
	case ArticleStatusPrivate:
		return "private"
	default:
		return "unknown"
	}
}

type Author struct {
	AuthorId int64 `json:"author_id"`
}
