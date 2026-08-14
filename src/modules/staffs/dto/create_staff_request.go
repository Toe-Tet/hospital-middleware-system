package dto

type CreateStaffRequest struct {
	HospitalID int    `json:"hospital_id" binding:"required,gte=1" example:"1" minimum:"1"`
	Username   string `json:"username" binding:"required,min=3,max=255" example:"alice.carter" minLength:"3" maxLength:"255"`
	Password   string `json:"password" binding:"required,min=6,max=100" example:"SecurePass123!" minLength:"6" maxLength:"100" format:"password"`
}
