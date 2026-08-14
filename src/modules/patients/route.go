package patients

import (
	"database/sql"

	"hospital-middleware-system/src/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	ctrl := NewController(db)

	patients := r.Group("/patients")
	{
    protected := patients.Group("")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("", ctrl.List)
			protected.GET("/:id", ctrl.GetByID)
		}
		
	}
}
