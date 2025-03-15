package route

// import (
// 	"kreditplus/config"
// 	"kreditplus/controller"
// 	"kreditplus/internal/domain"
// 	"kreditplus/internal/handler"
// 	"kreditplus/internal/repository"
// 	"kreditplus/internal/usecase"
// 	"kreditplus/middleware"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// func SetupRouter() *gin.Engine {
// 	r := gin.Default()

// 	r.Use(middleware.SecurityHeadersMiddleware())
// 	r.Use(middleware.LoggerMiddleware())
// 	r.Use(middleware.RateLimitMiddleware())

// 	api := r.Group("/api")

// 	setupPublicRoutes(api)
// 	setupProtectedRoutes(api)

// 	return r
// }

// func setupPublicRoutes(api *gin.RouterGroup) {
// 	api.POST("/login", controller.Login)
// 	api.POST("/refresh-token", controller.RefreshToken)
// 	api.POST("/logout", controller.Logout)
// }

// func setupProtectedRoutes(api *gin.RouterGroup) {
// 	protected := api.Group("/protected")
// 	protected.Use(middleware.AuthMiddleware())

// 	userRepo := repository.NewUserRepository(config.DB)
// 	userUsecase := usecase.NewUserUsecase(userRepo)
// 	userHandler := handler.NewUserHandler(userUsecase)

// 	protected.GET("/profile", func(c *gin.Context) {
// 		user, exists := c.Get("user")
// 		if !exists {
// 			c.JSON(401, gin.H{"error": "Unauthorized"})
// 			return
// 		}

// 		userModel := user.(domain.User)
// 		userDTO := userModel.ToDTO()

// 		c.JSON(200, gin.H{"user": userDTO})
// 	})

// 	protected.PUT("/profile", func(c *gin.Context) {
// 		user, exists := c.Get("user")
// 		if !exists {
// 			c.JSON(401, gin.H{"error": "Unauthorized"})
// 			return
// 		}
// 		authUser := user.(domain.User)
// 		c.Params = append(c.Params, gin.Param{Key: "id", Value: strconv.Itoa(int(authUser.UserID))})
// 		userHandler.UpdateUser(c)
// 	})

// 	setupUserRoutes(protected)
// 	setupCustomerRoutes(protected)
// 	setupLimitRoutes(protected)
// 	setupTransactionRoutes(protected)
// }

// // func setupAdminRoutes(api *gin.RouterGroup) {
// // 	admin := api.Group("/admin")
// // 	admin.Use(middleware.AuthMiddleware())
// // 	admin.Use(middleware.RoleMiddleware("admin"))

// // 	setupUserRoutes(admin)
// // }

// func setupUserRoutes(admin *gin.RouterGroup) {
// 	userRepo := repository.NewUserRepository(config.DB)
// 	userUsecase := usecase.NewUserUsecase(userRepo)
// 	userHandler := handler.NewUserHandler(userUsecase)

// 	users := admin.Group("/users")
// 	users.GET("/", userHandler.GetUser)
// 	users.GET("/:id", userHandler.GetUserByID)
// 	users.POST("/", userHandler.CreateUser)
// 	users.PUT("/:id", userHandler.UpdateUser)
// 	users.DELETE("/:id", userHandler.DeleteUser)
// }

// func setupCustomerRoutes(admin *gin.RouterGroup) {
// 	customerRepo := repository.NewCustomerRepository(config.DB)
// 	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
// 	customerHandler := handler.NewCustomerHandler(customerUsecase)

// 	customers := admin.Group("/customers")
// 	customers.GET("/", customerHandler.GetCustomer)
// 	customers.GET("/:id", customerHandler.GetCustomerByID)
// 	customers.POST("/", customerHandler.CreateCustomer)
// 	customers.PUT("/:id", customerHandler.UpdateCustomer)
// 	customers.DELETE("/:id", customerHandler.DeleteCustomer)
// }

// func setupLimitRoutes(admin *gin.RouterGroup) {
// 	limitRepo := repository.NewLimitRepository(config.DB)
// 	customerRepo := repository.NewCustomerRepository(config.DB)
// 	limitUsecase := usecase.NewLimitUsecase(limitRepo, customerRepo)
// 	limitHandler := handler.NewLimitHandler(limitUsecase)

// 	limits := admin.Group("/limits")
// 	limits.GET("/", limitHandler.GetLimit)
// 	limits.GET("/:id", limitHandler.GetLimitByID)
// 	limits.POST("/", limitHandler.CreateLimit)
// 	limits.PUT("/:id", limitHandler.UpdateLimit)
// 	limits.DELETE("/:id", limitHandler.DeleteLimit)
// }

// func setupTransactionRoutes(admin *gin.RouterGroup) {
// 	transactionRepo := repository.NewTransactionRepository(config.DB)
// 	limitRepo := repository.NewLimitRepository(config.DB)
// 	customerRepo := repository.NewCustomerRepository(config.DB)
// 	transactionUsecase := usecase.NewTransactionUsecase(customerRepo, limitRepo, transactionRepo)
// 	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

// 	transactions := admin.Group("/transactions")
// 	transactions.GET("/", transactionHandler.GetTransaction)
// 	transactions.GET("/:id", transactionHandler.GetTransactionByID)
// 	transactions.POST("/", transactionHandler.CreateTransaction)
// 	transactions.PUT("/:id", transactionHandler.UpdateTransaction)
// 	transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
// }
