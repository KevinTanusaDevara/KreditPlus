package route

import (
	"kreditplus/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the main router
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Global Middleware
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	api := r.Group("/api")

	// Public Routes
	SetupAuthRoutes(api)

	// Protected Routes
	protected := api.Group("/protected")
	protected.Use(middleware.AuthMiddleware())

	SetupUserRoutes(protected)
	SetupCustomerRoutes(protected)
	SetupLimitRoutes(protected)
	SetupTransactionRoutes(protected)

	return r
}
