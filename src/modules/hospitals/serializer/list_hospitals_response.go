package serializer

import (
	"hospital-middleware-system/src/modules/hospitals/model"
)

type ListHospitalItem struct {
	ID        int    `json:"id" example:"42" description:"Unique hospital ID" format:"int64" required:"true"`
	Name      string `json:"name" example:"St. Mary's General Hospital" description:"Full legal name of the hospital" required:"true"`
	Address   string `json:"address" example:"123 Medical Lane, Springfield, IL 62704" description:"Physical street address"`
	Phone     string `json:"phone" example:"+12175550123" description:"Main contact phone in E.164 format"`
	Status    string `json:"status" example:"active" description:"Operational status" enum:"active,inactive" required:"true"`
	CreatedAt string `json:"created_at" example:"2026-08-14T10:30:00Z" description:"ISO-8601 timestamp of creation" format:"date-time" required:"true"`
	UpdatedAt string `json:"updated_at" example:"2026-08-14T10:30:00Z" description:"ISO-8601 timestamp of last update" format:"date-time" required:"true"`
}

type ListHospitalsResponseItems = []*ListHospitalItem

func SerializeListHospital(h *model.Hospital) *ListHospitalItem {
	if h == nil {
		return nil
	}
	return &ListHospitalItem{
		ID:        h.ID,
		Name:      h.Name,
		Address:   h.Address,
		Phone:     h.Phone,
		Status:    h.Status,
		CreatedAt: h.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: h.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func SerializeListHospitals(list []*model.Hospital) ListHospitalsResponseItems {
	result := make([]*ListHospitalItem, 0, len(list))
	for _, h := range list {
		result = append(result, SerializeListHospital(h))
	}
	return result
}
