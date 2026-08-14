package serializer

import (
	"hospital-middleware-system/src/modules/staffs/model"
	"time"
)

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	cp := s
	return &cp
}

type CreateStaffResponse struct {
	ID         int     `json:"id" example:"7" format:"int64"`
	HospitalID int     `json:"hospital_id" example:"1" format:"int64"`
	Username   string  `json:"username" example:"alice.carter"`
	Name       *string  `json:"name" example:"Dr. Alice Carter"`
	Status     string  `json:"status" example:"active" enum:"active,inactive"`
	CreatedAt  *time.Time `json:"created_at" example:"2026-08-14T10:30:00Z" format:"date-time"`
	UpdatedAt  *time.Time `json:"updated_at" example:"2026-08-14T10:30:00Z" format:"date-time"`
}

func SerializeCreateStaff(s *model.Staff) *CreateStaffResponse {
	if s == nil {
		return nil
	}
	return &CreateStaffResponse{
		ID:         s.ID,
		HospitalID: s.HospitalID,
		Username:   s.Username,
		Name:       s.Name,
		Status:     s.Status,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}
