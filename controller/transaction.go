package controller

import (
	"fmt"
	"kreditplus/config"
	"kreditplus/model"
	"kreditplus/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionInput struct {
	TransactionNIK         string  `json:"transaction_nik" validate:"required,len=16,numeric"`
	TransactionOTR         float64 `json:"transaction_otr" validate:"required"`
	TransactionAdminFee    float64 `json:"transaction_admin_fee" validate:"required"`
	TransactionInstallment float64 `json:"transaction_installment" validate:"required"`
	TransactionInterest    float64 `json:"transaction_interest" validate:"required"`
	TransactionAssetName   string  `json:"transaction_asset_name" validate:"required"`
}

func CreateTransaction(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer model.Customer
	if err := config.DB.Where("customer_nik = ?", input.TransactionNIK).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var limit model.Limit
	if err := config.DB.
		Where("limit_nik = ? AND limit_tenor = ?", input.TransactionNIK, input.TransactionInstallment).
		Order("limit_remaining_amount DESC").
		First(&limit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No sufficient limit available for this tenor"})
		return
	}

	totalAmount := input.TransactionOTR + input.TransactionAdminFee + (input.TransactionInterest * input.TransactionOTR * float64(input.TransactionInstallment) / 100)

	if totalAmount > limit.LimitRemainingAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient limit"})
		return
	}

	transaction := model.Transaction{
		TransactionContractNumber: utils.GenerateContractNumber(),
		TransactionNIK:            input.TransactionNIK,
		TransactionLimit:          limit.LimitID,
		TransactionOTR:            input.TransactionOTR,
		TransactionAdminFee:       input.TransactionAdminFee,
		TransactionInstallment:    input.TransactionInstallment,
		TransactionInterest:       input.TransactionInterest,
		TransactionAssetName:      input.TransactionAssetName,
		TransactionDate:           time.Now(),
		TransactionCreatedBy:      authUser.(model.User).UserID,
		TransactionCreatedAt:      time.Now(),
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		limit.LimitUsedAmount += totalAmount
		limit.LimitRemainingAmount -= totalAmount

		if err := tx.Save(&limit).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUser.(model.User).UserID,
			"error":   err.Error(),
		}).Error("Failed to create transaction")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":                     authUser.(model.User).UserID,
		"transaction_contract_number": transaction.TransactionContractNumber,
		"edited_at":                   time.Now(),
	}).Infof("Transaction contract Number %s updated successfully by User %d", transaction.TransactionContractNumber, authUser.(model.User).UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully", "transaction": transaction})
}

func GetTransaction(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page value"})
		return
	}

	offset := (page - 1) * limit

	var transactions []model.Transaction
	err = config.DB.
		Preload("NIKCustomer.CreatedByUser").
		Preload("NIKCustomer.EditedByUser").
		Preload("IDLimit.NIKCustomer").
		Preload("IDLimit.CreatedByUser").
		Preload("IDLimit.EditedByUser").
		Preload("CreatedByUser").
		Preload("EditedByUser").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":         page,
		"limit":        limit,
		"transactions": transactions,
	})
}

func GetTransactionByID(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit ID"})
		return
	}

	var transaction model.Transaction
	err = config.DB.
		Preload("NIKCustomer.CreatedByUser").
		Preload("NIKCustomer.EditedByUser").
		Preload("IDLimit.NIKCustomer").
		Preload("IDLimit.CreatedByUser").
		Preload("IDLimit.EditedByUser").
		Preload("CreatedByUser").
		Preload("EditedByUser").
		First(&transaction, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

func UpdateTransaction(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(model.User)
	timeNow := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var transaction model.Transaction
	if err := config.DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = fmt.Sprintf("failed on '%s' rule", err.Tag())
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "fields": validationErrors})
		return
	}

	var customer model.Customer
	if err := config.DB.Where("customer_nik = ?", input.TransactionNIK).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var limit model.Limit
	if err := config.DB.Where("limit_nik = ? AND limit_tenor = ?", input.TransactionNIK, input.TransactionInstallment).First(&limit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No sufficient limit available for this tenor"})
		return
	}

	totalAmount := input.TransactionOTR + input.TransactionAdminFee + (input.TransactionInterest * input.TransactionOTR * float64(input.TransactionInstallment) / 100)

	originalAmount := transaction.TransactionOTR + transaction.TransactionAdminFee + (transaction.TransactionInterest * transaction.TransactionOTR * float64(transaction.TransactionInstallment) / 100)
	limit.LimitUsedAmount -= originalAmount
	limit.LimitRemainingAmount += originalAmount
	if totalAmount > limit.LimitRemainingAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient limit"})
		return
	}

	transaction.TransactionOTR = input.TransactionOTR
	transaction.TransactionAdminFee = input.TransactionAdminFee
	transaction.TransactionInstallment = input.TransactionInstallment
	transaction.TransactionInterest = input.TransactionInterest
	transaction.TransactionAssetName = input.TransactionAssetName
	transaction.TransactionEditedBy = &authUserModel.UserID
	transaction.TransactionEditedAt = &timeNow

	limit.LimitUsedAmount += totalAmount
	limit.LimitRemainingAmount -= totalAmount

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&transaction).Error; err != nil {
			return err
		}

		if err := tx.Save(&limit).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id":                     authUserModel.UserID,
			"transaction_contract_number": transaction.TransactionContractNumber,
			"error":                       err.Error(),
		}).Error("Failed to update transaction")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":                     authUserModel.UserID,
		"transaction_contract_number": transaction.TransactionContractNumber,
		"edited_at":                   timeNow,
	}).Infof("Transaction Contract Number %s updated successfully by User %d", transaction.TransactionContractNumber, authUserModel.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully", "transaction": transaction})
}

func DeleteTransaction(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(model.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var transaction model.Transaction
	if err := config.DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var limit model.Limit
	if err := config.DB.First(&limit, transaction.TransactionLimit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Associated limit record not found"})
		return
	}

	totalAmount := transaction.TransactionOTR + transaction.TransactionAdminFee +
		(transaction.TransactionInterest * transaction.TransactionOTR * float64(transaction.TransactionInstallment) / 100)

	limit.LimitUsedAmount -= totalAmount
	limit.LimitRemainingAmount += totalAmount

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&transaction).Error; err != nil {
			return err
		}

		if err := tx.Save(&limit).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id":        authUserModel.UserID,
			"transaction_id": transaction.TransactionID,
			"error":          err.Error(),
		}).Error("Failed to delete transaction")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":        authUserModel.UserID,
		"transaction_id": transaction.TransactionID,
	}).Infof("Transaction ID %d deleted successfully by User %d", transaction.TransactionID, authUserModel.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
