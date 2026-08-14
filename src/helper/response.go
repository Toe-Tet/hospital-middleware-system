package helper

import (
	"net/http"

	apperrors "hospital-middleware-system/src/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{
		Success: status < http.StatusBadRequest,
		Data:    data,
	})
}

func OK(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, data)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Message(c *gin.Context, status int, msg string) {
	c.JSON(status, Response{
		Success: status < http.StatusBadRequest,
		Message: msg,
	})
}

func Error(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		c.JSON(apperrors.StatusCode(appErr), Response{
			Success: false,
			Error: map[string]interface{}{
				"code":    appErr.Code,
				"message": appErr.Message,
				"details": appErr.Details,
			},
		})
		return
	}
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: map[string]interface{}{
			"code":    apperrors.CodeInternalServerError,
			"message": apperrors.MsgInternalServerError,
		},
	})
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

func Paginated(c *gin.Context, items interface{}, total int64, page, perPage int) {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	OK(c, PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}
