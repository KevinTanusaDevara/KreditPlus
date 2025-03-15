package handler

import (
	"html"
	"kreditplus/internal/domain"
	"kreditplus/internal/usecase"
	"kreditplus/internal/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CustomerHandler struct {
	usecase usecase.CustomerUsecase
}

func NewCustomerHandler(usecase usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{usecase: usecase}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt to CreateCustomer")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input domain.CustomerInput
	if err := c.ShouldBind(&input); err != nil {
		utils.Logger.Warn("Invalid request format for creating customer")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.Logger.Warnf("Validation error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerBirthDate, err := time.Parse("2006-01-02", input.CustomerBirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birth date format"})
		return
	}

	ktpPhotoPath, err := utils.SaveUploadedFile(c, "customer_ktp_photo", "uploads/ktp")
	if err != nil {
		utils.Logger.Warn("Failed to upload KTP photo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload KTP photo"})
		return
	}

	selfiePhotoPath, err := utils.SaveUploadedFile(c, "customer_selfie_photo", "uploads/selfie")
	if err != nil {
		utils.Logger.Warn("Failed to upload selfie photo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload selfie photo"})
		return
	}

	customer := domain.Customer{
		CustomerNIK:         input.CustomerNIK,
		CustomerFullName:    html.EscapeString(input.CustomerFullName),
		CustomerLegalName:   html.EscapeString(input.CustomerLegalName),
		CustomerBirthPlace:  html.EscapeString(input.CustomerBirthPlace),
		CustomerBirthDate:   customerBirthDate,
		CustomerSalary:      input.CustomerSalary,
		CustomerKTPPhoto:    ktpPhotoPath,
		CustomerSelfiePhoto: selfiePhotoPath,
		CustomerCreatedBy:   authUser.(domain.User).UserID,
		CustomerCreatedAt:   time.Now(),
	}

	err = h.usecase.CreateCustomer(customer)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUser.(domain.User).UserID,
			"error":   err.Error(),
		}).Error("Failed to create customer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":      authUser.(domain.User).UserID,
		"customer_nik": customer.CustomerNIK,
		"created_at":   customer.CustomerCreatedAt,
	}).Infof("Customer NIK %s created successfully by User %d", customer.CustomerNIK, authUser.(domain.User).UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Customer created successfully"})
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt to GetCustomer")
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

	customers, err := h.usecase.GetAllCustomers(limit, offset)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"limit":  limit,
			"offset": offset,
			"error":  err.Error(),
		}).Error("Failed to retrieve customers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"page":      page,
		"limit":     limit,
		"customers": customers,
	}).Info("Customers retrieved successfully")

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"limit":     limit,
		"customers": customers,
	})
}

func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt to GetCustomerByID")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Logger.Warn("Invalid customer ID in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	customer, err := h.usecase.GetCustomerByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"customer_id": id,
			"error":       err.Error(),
		}).Warn("Customer not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"customer_id": id,
	}).Info("Customer retrieved successfully")

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt to UpdateCustomer")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(domain.User)
	timeNow := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Logger.Warn("Invalid customer ID provided for update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	customer, err := h.usecase.GetCustomerByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"customer_id": id,
			"error":       err.Error(),
		}).Warn("Customer not found for update")
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var input domain.CustomerInput
	if err := c.ShouldBind(&input); err != nil {
		utils.Logger.Warn("Invalid request format for updating customer")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.Logger.Warnf("Validation error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.CustomerNIK != "" {
		customer.CustomerNIK = input.CustomerNIK
	}

	if input.CustomerFullName != "" {
		customer.CustomerFullName = input.CustomerFullName
	}

	if input.CustomerLegalName != "" {
		customer.CustomerLegalName = input.CustomerLegalName
	}

	if input.CustomerBirthPlace != "" {
		customer.CustomerBirthPlace = input.CustomerBirthPlace
	}

	if input.CustomerBirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", input.CustomerBirthDate)
		if err != nil {
			utils.Logger.Warn("Invalid birth date format for update")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birth date format"})
			return
		}
		customer.CustomerBirthDate = parsedDate
	}

	if input.CustomerSalary > 0 {
		customer.CustomerSalary = input.CustomerSalary
	}

	if _, err := c.FormFile("customer_ktp_photo"); err == nil {
		if customer.CustomerKTPPhoto != "" {
			if err := os.Remove(customer.CustomerKTPPhoto); err != nil {
				utils.Logger.Warnf("Failed to delete old KTP photo: %s", customer.CustomerKTPPhoto)
			}
		}

		ktpPhotoPath, err := utils.SaveUploadedFile(c, "customer_ktp_photo", "uploads/ktp")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload KTP photo"})
			return
		}
		customer.CustomerKTPPhoto = ktpPhotoPath
	}

	if _, err := c.FormFile("customer_selfie_photo"); err == nil {
		if customer.CustomerSelfiePhoto != "" {
			if err := os.Remove(customer.CustomerSelfiePhoto); err != nil {
				utils.Logger.Warnf("Failed to delete old Selfie photo: %s", customer.CustomerSelfiePhoto)
			}
		}

		selfiePhotoPath, err := utils.SaveUploadedFile(c, "customer_selfie_photo", "uploads/selfie")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload selfie photo"})
			return
		}
		customer.CustomerSelfiePhoto = selfiePhotoPath
	}

	customer.CustomerEditedBy = &authUserModel.UserID
	customer.CustomerEditedAt = &timeNow

	err = h.usecase.UpdateCustomer(*customer)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUserModel.UserID,
			"error":   err.Error(),
		}).Error("Failed to update customer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":      authUserModel.UserID,
		"customer_nik": customer.CustomerNIK,
		"updated_at":   customer.CustomerEditedAt,
	}).Infof("Customer NIK %s updated successfully by User %d", customer.CustomerNIK, authUserModel.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(domain.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	customer, err := h.usecase.GetCustomerByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"customer_id": id,
			"error":       err.Error(),
		}).Warn("Customer not found for update")
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	if customer.CustomerKTPPhoto != "" {
		if err := os.Remove(customer.CustomerKTPPhoto); err != nil {
			utils.Logger.Warnf("Failed to delete KTP photo: %s", customer.CustomerKTPPhoto)
		}
	}

	if customer.CustomerSelfiePhoto != "" {
		if err := os.Remove(customer.CustomerSelfiePhoto); err != nil {
			utils.Logger.Warnf("Failed to delete Selfie photo: %s", customer.CustomerSelfiePhoto)
		}
	}

	err = h.usecase.DeleteCustomer(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id":      authUserModel.UserID,
			"customer_nik": customer.CustomerNIK,
			"error":        err.Error(),
		}).Error("Failed to delete customer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":      authUserModel.UserID,
		"customer_nik": customer.CustomerNIK,
	}).Infof("Customer NIK %s deleted successfully by User %d", customer.CustomerNIK, authUserModel.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
