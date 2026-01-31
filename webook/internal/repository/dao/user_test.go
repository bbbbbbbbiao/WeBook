package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	my "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

/**
 * @author: biao
 * @date: 2026/1/30 下午3:03
 * @description:
 */

func TestGORMUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name string
		// 这里为啥不用ctrl
		// 因为这里不是gomock，而是sqlmock
		mock func(t *testing.T) *sql.DB
		ctx  context.Context
		u    User

		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			u: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
				Phone: sql.NullString{
					String: "123456789",
					Valid:  true,
				},
			},
			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(&my.MySQLError{
						Number: 1062,
					})
				require.NoError(t, err)
				return mockDB
			},
			u: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
				Phone: sql.NullString{
					String: "123456789",
					Valid:  true,
				},
			},
			wantErr: ErrUserDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(errors.New("数据库错误"))
				require.NoError(t, err)
				return mockDB
			},
			u: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
				Phone: sql.NullString{
					String: "123456789",
					Valid:  true,
				},
			},
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// 使用sqlmock 创建DB给GORM使用
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true, // 是否跳过初始化
				// 如果这里为false，就会发起一个初始化的调用，但是不是我预期的调用就会报错
			}), &gorm.Config{
				DisableAutomaticPing: true, // 跳过ping

				SkipDefaultTransaction: true, // 是否跳过事务
				// gorm默认会在执行sql语句之前，会先开启一个事务（不是我们测试时预期的）
			})
			dao := NewUserDao(db)
			err = dao.Insert(tc.ctx, tc.u)

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
