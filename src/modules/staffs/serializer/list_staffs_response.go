package serializer

import (
	"hospital-middleware-system/src/modules/staffs/model"
)

type ListStaffItem struct {
	ID         int     `json:"id" example:"7" format:"int64"`
	HospitalID int     `json:"hospital_id" example:"1" format:"int64"`
	Username   string  `json:"username" example:"alice.carter"`
	Name       *string  `json:"name" example:"Dr. Alice Carter"`
	Email      string  `json:"email" example:"alice.carter@hospital.com" format:"email"`
	Role       string  `json:"role" example:"doctor" enum:"staff,admin,doctor,nurse,receptionist"`
	Status     string  `json:"status" example:"active" enum:"active,inactive"`
	CreatedAt  string  `json:"created_at" example:"2026-08-14T10:30:00Z" format:"date-time"`
	UpdatedAt  string  `json:"updated_at" example:"2026-08-14T10:30:00Z" format:"date-time"`
}

type ListStaffsResponseItems = []*ListStaffItem

func SerializeListStaff(s *model.Staff) *ListStaffItem {
	if s == nil {
		return nil
	}
	return &ListStaffItem{
		ID:         s.ID,
		HospitalID: s.HospitalID,
		Username:   s.Username,
		Name:       s.Name,
		Status:     s.Status,
		CreatedAt:  s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func SerializeListStaffs(list []*model.Staff) ListStaffsResponseItems {
	result := make([]*ListStaffItem, 0, len(list))
	for _, s := range list {
		result = append(result, SerializeListStaff(s))
	}
	return result
}
