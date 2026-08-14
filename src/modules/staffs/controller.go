package staffs

import (
	"database/sql"

	apperrors "hospital-middleware-system/src/errors"
	"hospital-middleware-system/src/helper"
	"hospital-middleware-system/src/middleware"
	"hospital-middleware-system/src/modules/staffs/dto"
	"hospital-middleware-system/src/modules/staffs/serializer"

	"fmt"

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

// CreateStaff godoc
//
//	@Summary		Create a new staff account
//	@Description	Create a staff user bound to a hospital.
//	@Tags			Staff
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.CreateStaffRequest	true	"Staff details"
//	@Success		200		{object}	helper.Response{data=serializer.CreateStaffResponse}
//	@Router			/staffs [post]
func (ctrl *Controller) Create(c *gin.Context) {
	var req dto.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, helper.ValidateStruct(&req))
		return
	}
	if err := helper.ValidateStruct(&req); err != nil {
		helper.Error(c, err)
		return
	}
	staff, err := ctrl.service.Create(c.Request.Context(), &req)
	if err != nil {
		helper.Error(c, err)
		return
	}
  fmt.Println(staff)
	helper.OK(c, serializer.SerializeCreateStaff(staff))
}

// LoginStaff godoc
//
//	@Summary		Authenticate staff and issue JWT
//	@Description	Verify staff email + password (bcrypt compare) and return a short-lived Bearer JWT plus staff profile.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.LoginRequest	true	"Staff credentials"
//	@Success		200		{object}	helper.Response{data=serializer.LoginResponse}
//	@Router			/staffs/login [post]
func (ctrl *Controller) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, helper.ValidateStruct(&req))
		return
	}
	if err := helper.ValidateStruct(&req); err != nil {
		helper.Error(c, err)
		return
	}
	resp, err := ctrl.service.Login(c.Request.Context(), &req)
	if err != nil {
		helper.Error(c, err)
		return
	}
	helper.OK(c, resp)
}

// StaffMe godoc
//
// @Summary      Get the current authenticated staff profile
// @Description  Return the staff profile of the caller, resolved from the staff_id JWT claim.
// @Tags         Staff
// @Produce      json
// @Success      200 {object} helper.Response{data=serializer.StaffMeResponse}
// @Router       /staffs/me [get]
// @Security     BearerAuth
func (ctrl *Controller) Me(c *gin.Context) {
	staffID := middleware.GetStaffID(c)
	result, err := ctrl.service.Me(c.Request.Context(), staffID)
	if err != nil {
		helper.Error(c, err)
		return
	}
	helper.OK(c, result)
}
