package dao

import (
	"context"
	"gorm.io/gorm"
)

/**
 * @author: biao
 * @date: 2026/2/6 下午5:38
 * @description: 用于存放异步重试的请求
 */

type SmsRequest struct {
	ID      int64 `gorm:"primaryKey; autoIncrement"`
	Tpl     string
	Args    []string `gorm:"type:json"`
	numbers []string `gorm:"type:json"`

	// 时间，毫秒级别
	Ctime int64
	Utime int64
	Dtime int64
}

func (SmsRequest) TableName() string {
	return "request"
}

type SmsReq interface {
	Insert(ctx context.Context, smsRequest SmsRequest) error // 插入请求
	FindReq(ctx context.Context) error                       // 查询所有未被删除的请求
}

type GORMSmsReq struct {
	db *gorm.DB
}
