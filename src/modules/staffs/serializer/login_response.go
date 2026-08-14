package serializer

type LoginResponse struct {
	Token  string           `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0..." format:"jwt"`
	ExpiresAt int              `json:"expires_at" example:"1694502400"`
	Staff        *StaffMeResponse `json:"staff"`
}
