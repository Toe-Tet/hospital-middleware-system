package patients

import (
	"database/sql"

	apperrors "hospital-middleware-system/src/errors"
	"hospital-middleware-system/src/helper"
	"hospital-middleware-system/src/middleware"
	"hospital-middleware-system/src/modules/patients/dto"
	"hospital-middleware-system/src/modules/patients/serializer"

	"github.com/gin-gonic/gin"
)

var _ = apperrors.CodeBadRequest

type Controller struct {
	service Service
}

func NewController(db *sql.DB) *Controller {
	repo := NewRepository(db)
	svc := NewService(repo)
	return &Controller{service: svc}
}

// GetPatient godoc
//
//	@Summary		Get a patient by identifier
//	@Description	Retrieve a single patient using any identifier: integer patient ID, national ID, or passport number. Results are automatically scoped to the caller's hospital from JWT claims.
//	@Tags			Patients
//	@Produce		json
//	@Param			id	path		string	true	"Patient identifier: national_id, or passport_id" example:"1100100100123"
//	@Success		200	{object}	helper.Response{data=serializer.GetPatientResponse}
//	@Router			/patients/{id} [get]
//	@Security		BearerAuth
func (ctrl *Controller) GetByID(c *gin.Context) {
	hospitalID := middleware.GetHospitalID(c)
	identifier := c.Param("id")
	patient, err := ctrl.service.GetByID(c.Request.Context(), hospitalID, identifier)
	if err != nil {
		helper.Error(c, err)
		return
	}
	helper.OK(c, serializer.SerializeGetPatient(patient))
}

// ListPatients godoc
//
//	@Summary		List patients
//	@Description	Return a paginated list of patients scoped to the caller's hospital (from JWT claims). Supports optional filtering by identifiers, name components, date of birth, email, and phone number.
//	@Tags			Patients
//	@Produce		json
//	@Param			national_id		query		string	false	"Exact match on national ID / citizen ID number"
//	@Param			passport_id		query		string	false	"Exact match on passport number"
//	@Param			first_name		query		string	false	"Case-insensitive partial match on first name (English or Thai)"
//	@Param			middle_name		query		string	false	"Case-insensitive partial match on middle name (English or Thai)"
//	@Param			last_name		query		string	false	"Case-insensitive partial match on last name (English or Thai)"
//	@Param			date_of_birth	query		string	false	"Exact match on date of birth in YYYY-MM-DD format"
//	@Param			email			query		string	false	"Case-insensitive partial match on email address"
//	@Param			phone_number	query		string	false	"Case-insensitive partial match on phone number"
//	@Param			page			query		int		false	"Page number (1-indexed)"						minimum(1)	default(1)
//	@Param			per_page		query		int		false	"Items per page (1..100)"						minimum(1)	maximum(100)	default(10)
//	@Success		200				{object}	helper.Response{data=helper.PaginatedResponse{items=serializer.ListPatientsResponseItems}}
//	@Router			/patients [get]
//	@Security		BearerAuth
func (ctrl *Controller) List(c *gin.Context) {
	hospitalID := middleware.GetHospitalID(c)

	filters := &dto.PatientFilters{}
	if err := c.ShouldBindQuery(filters); err != nil {
		helper.Error(c, apperrors.NewBadRequest("Invalid query parameters"))
		return
	}
	if err := helper.ValidateStruct(filters); err != nil {
		helper.Error(c, err)
		return
	}

	page := filters.Page
	if page < 1 {
		page = 1
	}
	perPage := filters.PerPage
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	list, total, err := ctrl.service.List(c.Request.Context(), hospitalID, filters, page, perPage)
	if err != nil {
		helper.Error(c, err)
		return
	}
	helper.Paginated(c, serializer.SerializeListPatients(list), total, page, perPage)
}
