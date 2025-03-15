package route

import (
	"kreditplus/config"
	"kreditplus/internal/handler"
	"kreditplus/internal/repository"
	"kreditplus/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(protected *gin.RouterGroup) {
	customerRepo := repository.NewCustomerRepository(config.DB)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	customers := protected.Group("/customers")
	customers.GET("/", customerHandler.GetCustomer)
	customers.GET("/:id", customerHandler.GetCustomerByID)
	customers.POST("/", customerHandler.CreateCustomer)
	customers.PUT("/:id", customerHandler.UpdateCustomer)
	customers.DELETE("/:id", customerHandler.DeleteCustomer)
}
