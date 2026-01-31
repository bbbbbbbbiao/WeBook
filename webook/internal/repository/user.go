package repository

import (
	"context"
	"database/sql"
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
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindUserById(ctx context.Context, id int64) (domain.User, error)
	UpdateById(ctx context.Context, u domain.User) error
}

type CachedUserRepository struct {
	ud dao.UserDao
	uc cache.UserCache
}

func NewUserRepository(ud dao.UserDao, uc cache.UserCache) UserRepository {
	return &CachedUserRepository{
		ud: ud,
		uc: uc,
	}
}

func (ur *CachedUserRepository) Create(ctx context.Context, user domain.User) error {
	u := ur.DomainToEntity(user)
	return ur.ud.Insert(ctx, u)
}

func (ur *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.ud.FindByEmail(ctx, email)

	if err != nil {
		return domain.User{}, err
	}

	return ur.EntityToDomain(u), nil
}

func (ur *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.ud.FindByPhone(ctx, phone)

	if err != nil {
		return domain.User{}, err
	}

	return ur.EntityToDomain(u), nil
}

func (ur *CachedUserRepository) FindUserById(ctx context.Context, id int64) (domain.User, error) {

	u, err := ur.uc.Get(ctx, id)

	if err == nil {
		return u, nil
	}
	// err 就有三种错，一个是redis中没有，一个偶然性的没有命中，一种是redis崩掉了（面试时都需要查询，数据库的限流）
	ue, err := ur.ud.FindUserById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = ur.EntityToDomain(ue)

	_ = ur.uc.Set(ctx, u)

	//if err != nil {
	//	// 这里记录一下就行 （可以容忍这里的错误，但是需要记录是否是redis崩了）
	//}

	return u, nil
}

func (ur *CachedUserRepository) UpdateById(ctx context.Context, u domain.User) error {
	return ur.ud.UpdateById(ctx, dao.User{
		Id:           u.Id,
		NickName:     u.NickName,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
	})
}

func (ur *CachedUserRepository) DomainToEntity(user domain.User) dao.User {
	return dao.User{
		Id: user.Id,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Password: user.Password,
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		NickName:     user.NickName,
		Birthday:     user.Birthday,
		Introduction: user.Introduction,
	}
}

func (ur *CachedUserRepository) EntityToDomain(user dao.User) domain.User {
	return domain.User{
		Id:           user.Id,
		Email:        user.Email.String,
		Password:     user.Password,
		Phone:        user.Phone.String,
		NickName:     user.NickName,
		Birthday:     user.Birthday,
		Introduction: user.Introduction,
		Ctime:        user.Ctime,
		Utime:        user.Utime,
	}
}
