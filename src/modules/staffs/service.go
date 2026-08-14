package staffs

import (
	"context"
	"fmt"

	apperrors "hospital-middleware-system/src/errors"
	"hospital-middleware-system/src/helper"
	"hospital-middleware-system/src/modules/staffs/dto"
	"hospital-middleware-system/src/modules/staffs/model"
	"hospital-middleware-system/src/modules/staffs/serializer"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error)
	GetByID(ctx context.Context, id int) (*model.Staff, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error)
	Me(ctx context.Context, staffID int) (interface{}, error)
}

type service struct {
	repo StaffRepository
}

func NewService(repo StaffRepository) Service {
	return &service{repo: repo}
}

func hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(hash, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}

func (s *service) Create(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error) {
	pwdHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	staff := &model.Staff{
		HospitalID: req.HospitalID,
		Username:   req.Username,
		Password:   pwdHash,
		Status:     "active",
	}
  fmt.Println("init staff...", staff)

	created, err := s.repo.Create(ctx, staff)
  fmt.Println("created...", created)
	if err != nil {
		switch err.Error() {
		case "staff with this username already exists":
			return nil, apperrors.NewConflict(err.Error())
		case "staff with this email already exists":
			return nil, apperrors.NewConflict(err.Error())
		case "hospital does not exist":
			return nil, apperrors.NewBadRequest(err.Error())
		default:
			return nil, apperrors.NewInternal(err)
		}
	}
	return created, nil
}

func (s *service) GetByID(ctx context.Context, id int) (*model.Staff, error) {
	staff, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if staff == nil {
		return nil, apperrors.NewNotFound("Staff not found")
	}
	return staff, nil
}

func (s *service) Login(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error) {
	staff, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if staff == nil {
		return nil, apperrors.NewUnauthorized("Invalid username or password")
	}
	if staff.Status != "active" {
		return nil, apperrors.NewUnauthorized("Account is inactive")
	}
  fmt.Println("staff...", staff)
	if !VerifyPassword(staff.Password, req.Password) {
		return nil, apperrors.NewUnauthorized("Invalid username or password")
	}
	tokenResp, err := helper.GenerateToken(staff.ID, staff.HospitalID)
	if err != nil {
		return nil, apperrors.NewInternal(err, "Failed to generate token")
	}
  fmt.Println("token...", tokenResp)
  fmt.Println("staff...", staff)
	return &serializer.LoginResponse{
		Token:      tokenResp.Token,
		ExpiresAt:  tokenResp.ExpiresAt,
		Staff:        serializer.MeResponse(staff),
	}, nil
}

func (s *service) Me(ctx context.Context, staffID int) (interface{}, error) {
	staff, err := s.repo.GetByID(ctx, staffID)
	if err != nil {
		return nil, err
	}
	return serializer.MeResponse(staff), nil
}
