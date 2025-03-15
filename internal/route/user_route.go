package route

import (
	"kreditplus/config"
	"kreditplus/internal/handler"
	"kreditplus/internal/repository"
	"kreditplus/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(protected *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	users := protected.Group("/users")
	users.GET("/", userHandler.GetUser)
	users.GET("/:id", userHandler.GetUserByID)
	users.POST("/", userHandler.CreateUser)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)
}
