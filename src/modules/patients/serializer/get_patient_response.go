package serializer

import (
	"hospital-middleware-system/src/modules/patients/model"
)

type GetPatientResponse struct {
	ID           int     `json:"id" example:"101" format:"int64" required:"true"`
	HospitalID   int     `json:"hospital_id" example:"1" format:"int64" required:"true"`
	PatientHN    string  `json:"patient_hn" example:"HN-0001" required:"true"`
	FirstNameTh  *string `json:"first_name_th,omitempty" example:"สมชาย"`
	MiddleNameTh *string `json:"middle_name_th,omitempty"`
	LastNameTh   *string `json:"last_name_th,omitempty" example:"ไทย"`
	FirstNameEn  string  `json:"first_name_en" example:"John" required:"true"`
	MiddleNameEn *string `json:"middle_name_en,omitempty" example:"Michael"`
	LastNameEn   string  `json:"last_name_en" example:"Doe" required:"true"`
	NationalID   *string `json:"national_id,omitempty" example:"1100100100123"`
	PassportID   *string `json:"passport_id,omitempty" example:"AB1234567"`
	DateOfBirth  string  `json:"date_of_birth" example:"1985-03-15" format:"date" required:"true"`
	Gender       string  `json:"gender" example:"M" enum:"M,F" required:"true"`
	Email        string  `json:"email" example:"john@example.com"`
	PhoneNumber  string  `json:"phone_number" example:"0812345678" required:"true"`
	Status       string  `json:"status" example:"active" enum:"active,inactive" required:"true"`
	CreatedAt    string  `json:"created_at" example:"2026-08-14T10:30:00Z" format:"date-time" required:"true"`
	UpdatedAt    *string `json:"updated_at" example:"2026-08-14T10:30:00Z" format:"date-time"`
}

func SerializeGetPatient(p *model.Patient) *GetPatientResponse {
	if p == nil {
		return nil
	}
	return &GetPatientResponse{
		ID:           p.ID,
		HospitalID:   p.HospitalID,
		PatientHN:    p.PatientHN,
		FirstNameTh:  p.FirstNameTh,
		MiddleNameTh: p.MiddleNameTh,
		LastNameTh:   p.LastNameTh,
		FirstNameEn:  p.FirstNameEn,
		MiddleNameEn: p.MiddleNameEn,
		LastNameEn:   p.LastNameEn,
		NationalID:   p.NationalID,
		PassportID:   p.PassportID,
		DateOfBirth:  p.DateOfBirth.Format("2006-01-02"),
		Gender:       p.Gender,
		Email:        p.Email,
		PhoneNumber:  p.PhoneNumber,
		Status:       p.Status,
		CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    formatNullableTime(p.UpdatedAt),
	}
}
