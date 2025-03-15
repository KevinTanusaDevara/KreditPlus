package route

import (
	"kreditplus/config"
	"kreditplus/internal/handler"
	"kreditplus/internal/repository"
	"kreditplus/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupTransactionRoutes(protected *gin.RouterGroup) {
	transactionRepo := repository.NewTransactionRepository(config.DB)
	limitRepo := repository.NewLimitRepository(config.DB)
	customerRepo := repository.NewCustomerRepository(config.DB)
	transactionUsecase := usecase.NewTransactionUsecase(customerRepo, limitRepo, transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	transactions := protected.Group("/transactions")
	transactions.GET("/", transactionHandler.GetTransaction)
	transactions.GET("/:id", transactionHandler.GetTransactionByID)
	transactions.POST("/", transactionHandler.CreateTransaction)
	transactions.PUT("/:id", transactionHandler.UpdateTransaction)
	transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
}
