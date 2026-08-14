package serializer

import (
	"hospital-middleware-system/src/modules/staffs/model"
	"time"
)

type StaffMeResponse struct {
	ID         int     `json:"id" example:"1" format:"int64"`
	HospitalID int     `json:"hospital_id" example:"1" format:"int64"`
	Username   string  `json:"username" example:"admin"`
	Name       *string  `json:"name" example:"System Admin"`
	Status     string  `json:"status" example:"active" enum:"active,inactive"`
	CreatedAt  *time.Time `json:"created_at" example:"2026-08-14T10:30:00Z" format:"date-time"`
	UpdatedAt  *time.Time `json:"updated_at" example:"2026-08-14T10:30:00Z" format:"date-time"`
}

func MeResponse(s *model.Staff) *StaffMeResponse {
	if s == nil {
		return nil
	}
	return &StaffMeResponse{
		ID:         s.ID,
		HospitalID: s.HospitalID,
		Username:   s.Username,
		Name:       s.Name,
		Status:     s.Status,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}
