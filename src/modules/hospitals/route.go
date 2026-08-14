package hospitals

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	ctrl := NewController(db)

	hospitals := r.Group("/hospitals")
	{
		hospitals.GET("", ctrl.List)
	}
}
