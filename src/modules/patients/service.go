package patients

import (
	"context"
	"time"

	apperrors "hospital-middleware-system/src/errors"
	"hospital-middleware-system/src/modules/patients/dto"
	"hospital-middleware-system/src/modules/patients/model"
)

type Service interface {
	GetByID(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error)
	List(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error)
}

type service struct {
	repo PatientRepository
}

func NewService(repo PatientRepository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error) {
	if identifier == "" {
		return nil, apperrors.NewBadRequest("Patient identifier is required")
	}
	p, err := s.repo.FindByAnyID(ctx, hospitalID, identifier)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if p == nil {
		return nil, apperrors.NewNotFound("Patient not found")
	}
	return p, nil
}

func (s *service) List(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
	if filters.DateOfBirth != "" {
		if _, err := parseDate(filters.DateOfBirth); err != nil {
			return nil, 0, err
		}
	}
	list, total, err := s.repo.List(ctx, hospitalID, filters, page, perPage)
	if err != nil {
		return nil, 0, apperrors.NewInternal(err)
	}
	return list, total, nil
}

func parseDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, apperrors.NewBadRequest("Invalid date format, use YYYY-MM-DD")
	}
	return &t, nil
}
