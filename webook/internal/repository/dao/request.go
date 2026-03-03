package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

/**
 * @author: biao
 * @date: 2026/2/6 下午5:38
 * @description: 用于存放异步重试的请求
 */

type Request struct {
	ID      int64 `gorm:"primaryKey; autoIncrement"`
	Tpl     string
	Args    []string `gorm:"type:json"`
	Numbers []string `gorm:"type:json"`

	// 时间，毫秒级别
	Ctime int64
	Utime int64
	Dtime int64
}

func (Request) TableName() string {
	return "request"
}

type SmsReq interface {
	Insert(ctx context.Context, smsRequest Request) error // 插入请求
	FindReq(ctx context.Context) ([]Request, error)       // 查询所有未被删除的请求
}

type GORMSmsReq struct {
	db *gorm.DB
}

func NewGORMSmsReq(db *gorm.DB) *GORMSmsReq {
	return &GORMSmsReq{
		db: db,
	}
}

func (r *GORMSmsReq) Insert(ctx context.Context, smsRequest Request) error {
	now := time.Now().UnixMilli()
	smsRequest.Ctime = now
	smsRequest.Utime = now
	return r.db.WithContext(ctx).Create(&smsRequest).Error
}

func (r *GORMSmsReq) FindReq(ctx context.Context) ([]Request, error) {
	var smsReqs []Request
	// 加了一个行锁，避免并发问题
	err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("dtime is null").Find(&smsReqs).Error

	if err != nil {
		return nil, err
	}

	if len(smsReqs) == 0 {
		return smsReqs, nil
	}

	ids := make([]int64, 0, len(smsReqs))
	for _, smsReq := range smsReqs {
		ids = append(ids, smsReq.ID)
	}

	err = r.db.WithContext(ctx).Model(&Request{}).Where("id in (?)", ids).Update("dtime", time.Now().UnixMilli()).Error
	return smsReqs, err
}
