package route

import (
	"kreditplus/config"
	"kreditplus/internal/handler"
	"kreditplus/internal/repository"
	"kreditplus/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(api *gin.RouterGroup) {
	authRepo := repository.NewAuthRepository(config.DB)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh-token", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)
}
