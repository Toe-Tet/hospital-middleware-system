package dto

type ListHospitalsRequest struct {
	Page    int `form:"page" binding:"omitempty,gte=1" example:"1" description:"Page number (1-indexed)" minimum:"1"`
	PerPage int `form:"per_page" binding:"omitempty,gte=1,lte=100" example:"10" description:"Items returned per page" minimum:"1" maximum:"100"`
}
