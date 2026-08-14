package router

import (
	"database/sql"

	"hospital-middleware-system/src/helper"
	"hospital-middleware-system/src/middleware"
	hospitalModule "hospital-middleware-system/src/modules/hospitals"
	patientModule "hospital-middleware-system/src/modules/patients"
	staffModule "hospital-middleware-system/src/modules/staffs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "hospital-middleware-system/docs"
)

type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"hospital-middleware-system"`
}

func New(db *sql.DB) *gin.Engine {
	r := gin.New()

	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())
	r.Use(gin.Logger())

	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/docs/doc.json")))

	// register all module routes
	hospitalModule.RegisterRoutes(r, db)
	staffModule.RegisterRoutes(r, db)
	patientModule.RegisterRoutes(r, db)

	r.GET("/health", HealthCheck)

	return r
}

// HealthCheck godoc
//
//	@Summary		Liveness / readiness health check
//	@Description	Returns service name and 200 OK when the process is alive. Note: it does not verify DB connectivity.
//	@Tags			System
//	@Produce		json
//	@Success		200	{object}	helper.Response{data=HealthResponse}
//	@Router			/health [get]
func HealthCheck(c *gin.Context) {
	helper.OK(c, HealthResponse{
		Status:  "ok",
		Service: "hospital-middleware-system",
	})
}
