package route

import (
	"kreditplus/controller"
	"kreditplus/middleware"
	"kreditplus/model"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.LoggerMiddleware())

	api := r.Group("/api")

	api.Use(middleware.RateLimitMiddleware())

	api.POST("/login", controller.Login)
	api.POST("/refresh-token", controller.RefreshToken)
	api.POST("/logout", controller.Logout)

	protected := api.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		userModel := user.(model.User)
		userDTO := userModel.ToDTO()

		c.JSON(200, gin.H{"user": userDTO})
	})

	admin := protected.Group("/admin")
	admin.Use(middleware.RoleMiddleware("admin"))
	admin.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome Admin!"})
	})

	return r
}
