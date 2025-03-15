package handler

import (
	"fmt"
	"kreditplus/internal/domain"
	"kreditplus/internal/usecase"
	"kreditplus/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LimitHandler struct {
	usecase usecase.LimitUsecase
}

func NewLimitHandler(usecase usecase.LimitUsecase) *LimitHandler {
	return &LimitHandler{usecase: usecase}
}

func (h *LimitHandler) CreateLimit(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt to CreateLimit")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input domain.CreateLimitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Logger.Warn("Invalid request format for creating limit")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.Logger.Warnf("Validation error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.usecase.GetCustomerByNIK(input.LimitNIK)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"customer_nik": input.LimitNIK,
			"error":        err.Error(),
		}).Warn("Customer not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer NIK not found"})
		return
	}

	limit := domain.Limit{
		LimitNIK:             customer.CustomerNIK,
		LimitTenor:           input.LimitTenor,
		LimitAmount:          input.LimitAmount,
		LimitUsedAmount:      0,
		LimitRemainingAmount: input.LimitAmount,
		LimitCreatedBy:       authUser.(domain.User).UserID,
		LimitCreatedAt:       time.Now(),
	}

	err = h.usecase.CreateLimit(limit)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUser.(domain.User).UserID,
			"error":   err.Error(),
		}).Error("Failed to create limit")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create limit"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":    authUser.(domain.User).UserID,
		"limit_nik":  limit.LimitNIK,
		"created_at": time.Now(),
	}).Infof("Limit NIK %s created successfully by User %d", limit.LimitNIK, authUser.(domain.User).UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Limit created successfully"})
}

func (h *LimitHandler) GetLimit(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt to GetLimit")
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

	limits, err := h.usecase.GetAllLimits(limit, offset)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"limit":  limit,
			"offset": offset,
			"error":  err.Error(),
		}).Error("Failed to retrieve limits")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve limits"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"page":   page,
		"limit":  limit,
		"limits": limits,
	}).Info("Limits retrieved successfully")

	c.JSON(http.StatusOK, gin.H{
		"page":   page,
		"limit":  limit,
		"limits": limits,
	})
}

func (h *LimitHandler) GetLimitByID(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt to GetLimitByID")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Logger.Warn("Invalid limit ID in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit ID"})
		return
	}

	limit, err := h.usecase.GetLimitByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"limit_id": id,
			"error":    err.Error(),
		}).Warn("Limit not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Limit not found"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"limit_id": limit.LimitID,
	}).Info("Transaction retrieved successfully")

	c.JSON(http.StatusOK, limit)
}

func (h *LimitHandler) UpdateLimit(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(domain.User).UserRole != "admin" {
		utils.Logger.Warn("Unauthorized access attempt to UpdateLimit")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(domain.User)
	timeNow := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Logger.Warn("Invalid user ID provided for update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit ID"})
		return
	}

	limit, err := h.usecase.GetLimitByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"limit_id": id,
			"error":    err.Error(),
		}).Warn("Limit not found for update")
		c.JSON(http.StatusNotFound, gin.H{"error": "Limit not found"})
		return
	}

	var input domain.EditLimitInput
	if err := c.ShouldBind(&input); err != nil {
		utils.Logger.Warn("Invalid request format for updating limit")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	fmt.Println(input)
	if err := utils.Validate.Struct(input); err != nil {
		utils.Logger.Warnf("Validation error: %s", err.Error())
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

	if input.LimitUsedAmount != nil && *input.LimitUsedAmount > 0 {
		limit.LimitUsedAmount = *input.LimitUsedAmount
	}

	if input.LimitRemainingAmount != nil && *input.LimitRemainingAmount > 0 {
		limit.LimitRemainingAmount = *input.LimitRemainingAmount
	}

	limit.LimitEditedBy = &authUserModel.UserID
	limit.LimitEditedAt = &timeNow

	err = h.usecase.UpdateLimit(*limit)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUserModel.UserID,
			"error":   err.Error(),
		}).Error("Failed to update user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":          authUserModel.UserID,
		"update_limit_nik": limit.LimitNIK,
		"updated_at":       limit.LimitEditedAt,
	}).Infof("Limit NIK %s updated successfully by User %d", limit.LimitNIK, authUserModel.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Limit updated successfully"})
}

func (h *LimitHandler) DeleteLimit(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(domain.User)
	timeNow := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit ID"})
		return
	}

	limit, err := h.usecase.GetLimitByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"limit_id": id,
			"error":    err.Error(),
		}).Warn("Limit not found for update")
		c.JSON(http.StatusNotFound, gin.H{"error": "Limit not found"})
		return
	}

	err = h.usecase.DeleteLimit(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id":   authUserModel.UserID,
			"limit_nik": limit.LimitNIK,
			"error":     err.Error(),
		}).Error("Failed to delete customer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete limit"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":    authUserModel.UserID,
		"limit_nik":  limit.LimitNIK,
		"deleted_at": timeNow,
	}).Infof("Limit NIK %s deleted successfully by User %d", limit.LimitNIK, authUserModel.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Limit deleted successfully"})
}
