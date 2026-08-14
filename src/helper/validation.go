package helper

import (
	"strconv"

	apperrors "hospital-middleware-system/src/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return apperrors.NewInternal(err)
	}

	validationErrors := err.(validator.ValidationErrors)
	var fieldErrors []ValidationFieldError

	for _, fe := range validationErrors {
		fieldErrors = append(fieldErrors, ValidationFieldError{
			Field:   fe.Field(),
			Message: buildFieldErrorMsg(fe),
		})
	}

	return apperrors.NewValidationError("", fieldErrors)
}

func buildFieldErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Minimum length is " + fe.Param()
	case "max":
		return "Maximum length is " + fe.Param()
	case "gte":
		return "Must be greater than or equal to " + fe.Param()
	case "lte":
		return "Must be less than or equal to " + fe.Param()
	case "len":
		return "Length must be exactly " + fe.Param()
	case "oneof":
		return "Must be one of: " + fe.Param()
	default:
		return "Invalid value"
	}
}

func ParseIDParam(c *gin.Context, param string) (int, error) {
	idStr := c.Param(param)
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return 0, apperrors.NewBadRequest("Invalid ID parameter")
	}
	return id, nil
}

func ParsePageParams(c *gin.Context) (page, perPage int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	perPage, _ = strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}
	return
}
