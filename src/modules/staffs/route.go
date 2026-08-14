package staffs

import (
	"database/sql"

	"hospital-middleware-system/src/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	ctrl := NewController(db)

	staffs := r.Group("/staffs")
	{
		staffs.POST("", ctrl.Create)
		staffs.POST("/create", ctrl.Create)
		staffs.POST("/login", ctrl.Login)

		protected := staffs.Group("")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("/me", ctrl.Me)
		}
	}
}
