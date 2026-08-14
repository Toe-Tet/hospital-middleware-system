package model

import "time"

type Staff struct {
	ID         int        `json:"id"`
	HospitalID int        `json:"hospital_id"`
	Username   string     `json:"username"`
	Name       *string    `json:"name"`
	Password   string     `json:"-"`
	Status     string     `json:"status"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}
