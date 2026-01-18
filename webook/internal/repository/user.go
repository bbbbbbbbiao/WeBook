package repository

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/cache"
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
	uc *cache.UserCache
}

func NewUserRepository(ud *dao.UserDao, uc *cache.UserCache) *UserRepository {
	return &UserRepository{
		ud: ud,
		uc: uc,
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

	u, err := ur.uc.Get(ctx, id)

	if err == nil {
		return u, nil
	}
	// err 就有三种错，一个是redis中没有，一个偶然性的没有命中，一种是redis崩掉了（面试时都需要查询，数据库的限流）
	ue, err := ur.ud.FindUserById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = domain.User{
		Id:           ue.Id,
		Email:        ue.Email,
		NickName:     ue.NickName,
		Birthday:     ue.Birthday,
		Introduction: ue.Introduction,
	}

	err = ur.uc.Set(ctx, u)

	if err != nil {
		// 这里记录一下就行 （可以容忍这里的错误，但是需要记录是否是redis崩了）
	}

	return u, nil
}

func (ur *UserRepository) UpdateById(ctx context.Context, u domain.User) error {
	return ur.ud.UpdateById(ctx, dao.User{
		Id:           u.Id,
		NickName:     u.NickName,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
	})
}
