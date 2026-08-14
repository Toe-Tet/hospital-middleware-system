package hospitals

import (
	"database/sql"

	apperrors "hospital-middleware-system/src/errors"
	"hospital-middleware-system/src/helper"
	"hospital-middleware-system/src/modules/hospitals/serializer"

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

// ListHospitals godoc
//
//	@Summary		List hospitals
//	@Description	Return a paginated list of all hospitals sorted by most recent first.
//	@Tags			Hospitals
//	@Produce		json
//	@Param			page		query		int	false	"Page number (1-indexed)"					minimum(1)	default(1)
//	@Param			per_page	query		int	false	"Items returned per page (1..100)"			minimum(1)	maximum(100)	default(10)
//	@Success		200			{object}	helper.Response{data=helper.PaginatedResponse{items=serializer.ListHospitalsResponseItems}}
//	@Router			/hospitals [get]
func (ctrl *Controller) List(c *gin.Context) {
	page, perPage := helper.ParsePageParams(c)
	list, total, err := ctrl.service.List(c.Request.Context(), page, perPage)
	if err != nil {
		helper.Error(c, err)
		return
	}
	helper.Paginated(c, serializer.SerializeListHospitals(list), total, page, perPage)
}
