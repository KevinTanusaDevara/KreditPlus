package route

import (
	"kreditplus/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.CSRFMiddleware())

	api := r.Group("/api")

	SetupAuthRoutes(api)

	protected := api.Group("/protected")
	protected.Use(middleware.AuthMiddleware())

	SetupUserRoutes(protected)
	SetupCustomerRoutes(protected)
	SetupLimitRoutes(protected)
	SetupTransactionRoutes(protected)

	return r
}
