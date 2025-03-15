package route

import (
	"kreditplus/config"
	"kreditplus/internal/handler"
	"kreditplus/internal/repository"
	"kreditplus/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupLimitRoutes(protected *gin.RouterGroup) {
	limitRepo := repository.NewLimitRepository(config.DB)
	customerRepo := repository.NewCustomerRepository(config.DB)
	limitUsecase := usecase.NewLimitUsecase(limitRepo, customerRepo)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	limits := protected.Group("/limits")
	limits.GET("/", limitHandler.GetLimit)
	limits.GET("/:id", limitHandler.GetLimitByID)
	limits.POST("/", limitHandler.CreateLimit)
	limits.PUT("/:id", limitHandler.UpdateLimit)
	limits.DELETE("/:id", limitHandler.DeleteLimit)
}
