package middleware

import (
	"strings"

	apperrors "hospital-middleware-system/src/errors"
	"hospital-middleware-system/src/helper"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	StaffIDKey    contextKey = "staff_id"
	HospitalIDKey contextKey = "hospital_id"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			helper.Error(c, apperrors.NewUnauthorized("Authorization header is required"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			helper.Error(c, apperrors.NewUnauthorized("Invalid authorization format, use Bearer <token>"))
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(parts[1])
		if err != nil {
			helper.Error(c, apperrors.NewUnauthorized("Invalid or expired token"))
			c.Abort()
			return
		}

		c.Set(string(StaffIDKey), claims.StaffID)
		c.Set(string(HospitalIDKey), claims.HospitalID)

		c.Next()
	}
}

func GetStaffID(c *gin.Context) int {
	id, _ := c.Get(string(StaffIDKey))
	return id.(int)
}

func GetHospitalID(c *gin.Context) int {
	id, _ := c.Get(string(HospitalIDKey))
	return id.(int)
}
