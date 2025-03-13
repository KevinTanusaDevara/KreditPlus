package controller

import (
	"kreditplus/config"
	"kreditplus/model"
	"kreditplus/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LimitInput struct {
	LimitNIK    string  `json:"limit_nik" validate:"required,len=16,numeric"`
	LimitTenor  int     `json:"limit_tenor" validate:"required"`
	LimitAmount float64 `json:"limit_amount" validate:"required"`
}

func CreateLimit(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input LimitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer model.Customer
	if err := config.DB.Where("customer_nik = ?", input.LimitNIK).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	limit := model.Limit{
		LimitNIK:       input.LimitNIK,
		LimitTenor:     input.LimitTenor,
		LimitAmount:    input.LimitAmount,
		LimitCreatedBy: authUser.(model.User).UserID,
		LimitCreatedAt: time.Now(),
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&limit).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUser.(model.User).UserID,
			"error":   err.Error(),
		}).Error("Failed to create limit")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create limit"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":    authUser.(model.User).UserID,
		"limit_nik":  input.LimitNIK,
		"created_at": time.Now(),
	}).Info("Limit created successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Limit created successfully", "limit": limit})
}

func GetLimit(c *gin.Context) {
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

	var limits []model.Limit
	err = config.DB.
		Preload("NIKCustomer").
		Preload("CreatedByUser").
		Preload("EditedByUser").
		Limit(limit).
		Offset(offset).
		Find(&limits).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch limits"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":   page,
		"limit":  limit,
		"limits": limits,
	})
}

func GetLimitByID(c *gin.Context) {
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

	var limit model.Limit
	err = config.DB.
		Preload("NIKCustomer").
		Preload("CreatedByUser").
		Preload("EditedByUser").
		First(&limit, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Limit not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"limit": limit})
}

func UpdateLimit(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(model.User)
	timeNow := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit ID"})
		return
	}

	var limit model.Limit
	if err := config.DB.First(&limit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Limit not found"})
		return
	}

	var input LimitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.LimitNIK != "" {
		limit.LimitNIK = input.LimitNIK
	}

	if input.LimitTenor > 0 {
		limit.LimitTenor = input.LimitTenor
	}

	if input.LimitAmount > 0 {
		limit.LimitAmount = input.LimitAmount
	}

	limit.LimitEditedBy = &authUserModel.UserID
	limit.LimitEditedAt = &timeNow

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&limit).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUserModel.UserID,
			"error":   err.Error(),
		}).Error("Failed to update limit")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update limit"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":    authUserModel.UserID,
		"limit_nik":  input.LimitNIK,
		"created_at": time.Now(),
	}).Infof("Limit ID %d updated successfully by User %d", limit.LimitID, authUserModel.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Limit updated successfully"})
}

func DeleteLimit(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(model.User)
	timeNow := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit ID"})
		return
	}

	var limit model.Limit
	if err := config.DB.First(&limit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Limit not found"})
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&limit).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id":  authUserModel.UserID,
			"limit_id": limit.LimitID,
			"error":    err.Error(),
		}).Error("Failed to delete limit")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete limit"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":    authUserModel.UserID,
		"limit_id":   limit.LimitID,
		"deleted_at": timeNow,
	}).Infof("Limit ID %d deleted successfully by User %d", limit.LimitID, authUserModel.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Limit deleted successfully"})
}
