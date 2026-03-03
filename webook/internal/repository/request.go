package repository

import (
	"context"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/domain"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao"
)

type ReqRepository interface {
	Create(ctx context.Context, req domain.Request) error
	Find(ctx context.Context) ([]domain.Request, error)
}

type SmsReqRepository struct {
	sq dao.SmsReq
}

func NewSmsReqRepository(sq dao.SmsReq) ReqRepository {
	return &SmsReqRepository{
		sq: sq,
	}
}

func (r *SmsReqRepository) Create(ctx context.Context, req domain.Request) error {
	reqEntity := r.DomainToEntity(req)
	return r.sq.Insert(ctx, reqEntity)
}

func (r *SmsReqRepository) Find(ctx context.Context) ([]domain.Request, error) {
	reqs, err := r.sq.FindReq(ctx)
	if err != nil {
		return nil, err
	}

	reqEntities := make([]domain.Request, 0, len(reqs))
	for _, req := range reqs {
		reqEntity := r.EntityToDomain(req)
		reqEntities = append(reqEntities, reqEntity)
	}
	return reqEntities, nil
}

func (r *SmsReqRepository) EntityToDomain(req dao.Request) domain.Request {
	return domain.Request{
		ID:      req.ID,
		Tpl:     req.Tpl,
		Args:    req.Args,
		Numbers: req.Numbers,
		Ctime:   req.Ctime,
		Utime:   req.Utime,
	}
}

func (r *SmsReqRepository) DomainToEntity(req domain.Request) dao.Request {
	return dao.Request{
		ID:      req.ID,
		Tpl:     req.Tpl,
		Args:    req.Args,
		Numbers: req.Numbers,
		Ctime:   req.Ctime,
		Utime:   req.Utime,
	}
}
