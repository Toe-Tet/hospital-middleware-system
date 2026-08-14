package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required,min=1" example:"password" format:"password"`
}
