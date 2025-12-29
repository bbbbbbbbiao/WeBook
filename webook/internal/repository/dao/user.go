package dao

import (
	"context"
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
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

// 数据库表的结构体
type User struct {
	Id           int64  `gorm:"primaryKey, autoIncrement"`
	Email        string `gorm:"unique"`
	Password     string
	NickName     string
	Birthday     string
	Introduction string

	// 时间，时间戳毫秒数
	Ctime int64
	Utime int64
}

// UserDao 是用户的 DAO 层
type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (ud *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := ud.db.WithContext(ctx).Create(&u).Error
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo = 1062
		if mysqlError.Number == uniqueIndexErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (ud *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (ud *UserDao) FindUserById(ctx context.Context, id int64) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

func (ud *UserDao) UpdateById(ctx context.Context, u User) error {
	return ud.db.WithContext(ctx).Where("id = ?", u.Id).Updates(&u).Error
}
