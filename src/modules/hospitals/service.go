package hospitals

import (
	"context"

	"hospital-middleware-system/src/errors"
	"hospital-middleware-system/src/modules/hospitals/model"
)

type Service interface {
	List(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error)
}

type service struct {
	repo HospitalRepository
}

func NewService(repo HospitalRepository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
	list, total, err := s.repo.List(ctx, page, perPage)
	if err != nil {
		return nil, 0, errors.NewInternal(err)
	}
	return list, total, nil
}
