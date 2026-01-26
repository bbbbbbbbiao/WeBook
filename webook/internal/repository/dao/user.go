package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

/**
 * @author: biao
 * @date: 2025/12/22 下午9:35
 * @description:
 */

var (
	ErrUserDuplicate = errors.New("邮箱或手机号冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

// 数据库表的结构体
type User struct {
	Id       int64          `gorm:"primaryKey, autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string

	// 唯一索引，可以允许有多个空值
	// 不允许有多个空字符串
	Phone sql.NullString `gorm:"unique"`

	NickName     string
	Birthday     string
	Introduction string

	// 时间，时间戳毫秒数
	Ctime int64
	Utime int64
}

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindUserById(ctx context.Context, id int64) (User, error)
	UpdateById(ctx context.Context, u User) error
}

// UserDao 是用户的 DAO 层
type GORMUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GORMUserDao{
		db: db,
	}
}

func (ud *GORMUserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := ud.db.WithContext(ctx).Create(&u).Error
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo = 1062
		if mysqlError.Number == uniqueIndexErrNo {
			// 邮箱或手机号冲突
			return ErrUserDuplicate
		}
	}
	return err
}

func (ud *GORMUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (ud *GORMUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (ud *GORMUserDao) FindUserById(ctx context.Context, id int64) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

func (ud *GORMUserDao) UpdateById(ctx context.Context, u User) error {
	return ud.db.WithContext(ctx).Where("id = ?", u.Id).Updates(&u).Error
}
