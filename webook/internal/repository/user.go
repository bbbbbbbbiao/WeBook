package repository

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao"
)

/**
 * @author: biao
 * @date: 2025/12/22 下午9:35
 * @description:
 */

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	ud *dao.UserDao
}

func NewUserRepository(ud *dao.UserDao) *UserRepository {
	return &UserRepository{
		ud: ud,
	}
}

func (ur *UserRepository) Create(ctx context.Context, user domain.User) error {
	u := dao.User{
		Email:    user.Email,
		Password: user.Password,
	}

	return ur.ud.Insert(ctx, u)
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.ud.FindByEmail(ctx, email)

	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (ur *UserRepository) FindUserById(ctx context.Context, id int64) (domain.User, error) {
	u, err := ur.ud.FindUserById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:           u.Id,
		Email:        u.Email,
		NickName:     u.NickName,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
	}, err
}

func (ur *UserRepository) UpdateById(ctx context.Context, u domain.User) error {
	return ur.ud.UpdateById(ctx, dao.User{
		Id:           u.Id,
		NickName:     u.NickName,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
	})
}
