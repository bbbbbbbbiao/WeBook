package service

import (
	"context"
	"errors"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/**
 * @author: biao
 * @date: 2025/12/22 下午9:18
 * @description:
 */

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("用户名或密码不正确")
	ErrUserNotFound          = errors.New("用户不存在")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, u.Email)

	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return user, nil
}

func (svc *UserService) Edit(ctx *gin.Context, user domain.User) error {
	_, err := svc.repo.FindUserById(ctx, user.Id)

	if err == repository.ErrUserNotFound {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	err = svc.repo.UpdateById(ctx, user)
	return err
}

func (svc *UserService) Profile(ctx *gin.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindUserById(ctx, id)

	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrUserNotFound
	}

	return u, err
}

func (svc *UserService) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)

	// 快路径
	if err != repository.ErrUserNotFound {
		// err = nil 会进来
		// err != ErrUserNotFound 也会进来
		return u, err
	}

	// 在系统资源不足，触发降级之后，不执行慢路径了
	// 慢路径
	// 明确是没有这个用户的
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return u, err
	}

	// 这里会有问题，可能会出现主从延迟问题
	u, err = svc.repo.FindByPhone(ctx, phone)
	return u, err
}
