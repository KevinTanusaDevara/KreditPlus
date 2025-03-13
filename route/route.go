package route

import (
	"kreditplus/controller"
	"kreditplus/middleware"
	"kreditplus/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	api := r.Group("/api")

	setupPublicRoutes(api)
	setupProtectedRoutes(api)

	return r
}

func setupPublicRoutes(api *gin.RouterGroup) {
	api.POST("/login", controller.Login)
	api.POST("/refresh-token", controller.RefreshToken)
	api.POST("/logout", controller.Logout)
}

func setupProtectedRoutes(api *gin.RouterGroup) {
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

	protected.PUT("/profile", func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		authUser := user.(model.User)
		c.Params = append(c.Params, gin.Param{Key: "id", Value: strconv.Itoa(int(authUser.UserID))})
		controller.UpdateUser(c)
	})

	setupUserRoutes(protected)
	setupCustomerRoutes(protected)
	setupLimitRoutes(protected)
}

// func setupAdminRoutes(api *gin.RouterGroup) {
// 	admin := api.Group("/admin")
// 	admin.Use(middleware.AuthMiddleware())
// 	admin.Use(middleware.RoleMiddleware("admin"))

// 	setupUserRoutes(admin)
// }

func setupUserRoutes(admin *gin.RouterGroup) {
	users := admin.Group("/users")

	users.GET("/", controller.GetUser)
	users.GET("/:id", controller.GetUserByID)
	users.POST("/", controller.CreateUser)
	users.PUT("/:id", controller.UpdateUser)
	users.DELETE("/:id", controller.DeleteUser)
}

func setupCustomerRoutes(admin *gin.RouterGroup) {
	customers := admin.Group("/customers")

	customers.GET("/", controller.GetCustomer)
	customers.GET("/:id", controller.GetCustomerByID)
	customers.POST("/", controller.CreateCustomer)
	customers.PUT("/:id", controller.UpdateCustomer)
	customers.DELETE("/:id", controller.DeleteCustomer)
}

func setupLimitRoutes(admin *gin.RouterGroup) {
	limits := admin.Group("/limits")

	limits.GET("/", controller.GetLimit)
	limits.GET("/:id", controller.GetLimitByID)
	limits.POST("/", controller.CreateLimit)
	limits.PUT("/:id", controller.UpdateLimit)
	limits.DELETE("/:id", controller.DeleteLimit)
}
