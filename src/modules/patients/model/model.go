package model

import "time"

type Patient struct {
	ID           int        `json:"id"`
	HospitalID   int        `json:"hospital_id"`
	PatientHN    string     `json:"patient_hn"`
	FirstNameTh  *string    `json:"first_name_th"`
	MiddleNameTh *string    `json:"middle_name_th"`
	LastNameTh   *string    `json:"last_name_th"`
	FirstNameEn  string     `json:"first_name_en"`
	MiddleNameEn *string    `json:"middle_name_en"`
	LastNameEn   string     `json:"last_name_en"`
	NationalID   *string    `json:"national_id"`
	PassportID   *string    `json:"passport_id"`
	DateOfBirth  time.Time  `json:"date_of_birth"`
	Gender       string     `json:"gender"`
	Email       string     `json:"email"`
	PhoneNumber  string     `json:"phone_number"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
