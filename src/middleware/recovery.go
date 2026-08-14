package middleware

import (
	apperrors "hospital-middleware-system/src/errors"
	"hospital-middleware-system/src/helper"
	"hospital-middleware-system/src/logger"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error().
					Interface("panic", err).
					Str("stack", string(debug.Stack())).
					Str("path", c.Request.URL.Path).
					Msg("Panic recovered")

				helper.Error(c, apperrors.NewInternal(nil, "Internal server error"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
