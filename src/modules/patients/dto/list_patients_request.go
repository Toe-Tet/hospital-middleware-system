package dto

type PatientFilters struct {
	NationalID  string `form:"national_id" example:"1100100100123" description:"Exact match on national ID / citizen ID"`
	PassportID  string `form:"passport_id" example:"AB1234567" description:"Exact match on passport number"`
	FirstName   string `form:"first_name" example:"John" description:"Case-insensitive partial match on first name (EN or TH)"`
	MiddleName  string `form:"middle_name" example:"Michael" description:"Case-insensitive partial match on middle name (EN or TH)"`
	LastName    string `form:"last_name" example:"Doe" description:"Case-insensitive partial match on last name (EN or TH)"`
	DateOfBirth string `form:"date_of_birth" example:"1990-01-31" description:"Exact match on birth date (YYYY-MM-DD)"`
	Email       string `form:"email" example:"john@example.com" description:"Case-insensitive partial match on email"`
	PhoneNumber string `form:"phone_number" example:"0812345678" description:"Case-insensitive partial match on phone number"`
	Page        int    `form:"page" binding:"omitempty,gte=1" example:"1" description:"Page number (1-indexed)" minimum:"1"`
	PerPage     int    `form:"per_page" binding:"omitempty,gte=1,lte=100" example:"10" description:"Items per page" minimum:"1" maximum:"100"`
}
