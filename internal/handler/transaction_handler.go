package handler

import (
	"kreditplus/internal/domain"
	"kreditplus/internal/usecase"
	"kreditplus/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransactionHandler struct {
	usecase usecase.TransactionUsecase
}

func NewTransactionHandler(usecase usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{usecase: usecase}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt to CreateTransaction")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input domain.TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Logger.Warn("Invalid request format for creating transaction")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.Logger.Warnf("Validation error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.usecase.GetCustomerByNIK(input.TransactionNIK)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"transaction_nik": input.TransactionNIK,
			"error":           err.Error(),
		}).Warn("Customer not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer NIK not found"})
		return
	}

	err = h.usecase.CreateTransactionWithLimitUpdate(authUser.(domain.User).UserID, customer, input)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUser.(domain.User).UserID,
			"error":   err.Error(),
		}).Error("Failed to create transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":         authUser.(domain.User).UserID,
		"transaction_nik": customer.CustomerNIK,
		"created_at":      time.Now(),
	}).Infof("Transaction NIK %s created successfully by User %d", customer.CustomerNIK, authUser.(domain.User).UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully"})
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt to GetTransaction")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		utils.Logger.Warn("Invalid limit value in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page value"})
		return
	}

	offset := (page - 1) * limit

	transactions, err := h.usecase.GetAllTransactions(limit, offset)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"limit":  limit,
			"offset": offset,
			"error":  err.Error(),
		}).Error("Failed to retrieve transactions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"page":         page,
		"limit":        limit,
		"transactions": transactions,
	}).Info("transactions retrieved successfully")

	c.JSON(http.StatusOK, gin.H{
		"page":         page,
		"limit":        limit,
		"transactions": transactions,
	})
}

func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt to GetTransactionByID")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Logger.Warn("Invalid limit transaction in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.usecase.GetTransactionByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"transaction_id": id,
			"error":          err.Error(),
		}).Warn("Transaction not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"transaction_id": transaction.TransactionID,
	}).Info("Transaction retrieved successfully")

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(domain.User).UserRole != "admin" {
		utils.Logger.Warn("Unauthorized access attempt to UpdateTransaction")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Logger.Warn("Invalid user ID provided for update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transasction ID"})
		return
	}

	transaction, err := h.usecase.GetTransactionByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"transaction_id": id,
			"error":          err.Error(),
		}).Warn("Transaction not found for update")
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var input domain.TransactionInput
	if err := c.ShouldBind(&input); err != nil {
		utils.Logger.Warn("Invalid request format for updating transaction")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.Logger.Warnf("Validation error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.usecase.GetCustomerByNIK(input.TransactionNIK)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"transaction_nik": input.TransactionNIK,
			"error":           err.Error(),
		}).Warn("Customer not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer NIK not found"})
		return
	}

	err = h.usecase.UpdateTransactionWithLimitUpdate(authUser.(domain.User).UserID, customer, transaction, input)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUser.(domain.User).UserID,
			"error":   err.Error(),
		}).Error("Failed to update transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":                     authUser.(domain.User).UserID,
		"transaction_contract_number": transaction.TransactionContractNumber,
		"updated_at":                  time.Now(),
	}).Infof("Transaction Contract Number %s edited successfully by User %d", transaction.TransactionContractNumber, authUser.(domain.User).UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction edited successfully"})
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(domain.User).UserRole != "admin" {
		utils.Logger.Warn("Unauthorized access attempt to DeleteTransaction")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Logger.Warn("Invalid user ID provided for delete")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transasction ID"})
		return
	}

	transaction, err := h.usecase.GetTransactionByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"transaction_id": id,
			"error":          err.Error(),
		}).Warn("Transaction not found for delete")
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	err = h.usecase.DeleteTransactionWithLimitUpdate(authUser.(domain.User).UserID, transaction)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id":                     authUser.(domain.User).UserID,
			"transaction_contract_number": transaction.TransactionContractNumber,
			"error":                       err.Error(),
		}).Error("Failed to delete transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":                     authUser.(domain.User).UserID,
		"transaction_contract_number": transaction.TransactionContractNumber,
	}).Infof("Transaction Contract Number %s deleted successfully by User %d", transaction.TransactionContractNumber, authUser.(domain.User).UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
