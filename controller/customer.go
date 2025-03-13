package controller

import (
	"html"
	"kreditplus/config"
	"kreditplus/model"
	"kreditplus/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerInput struct {
	CustomerNIK        string  `form:"customer_nik" validate:"required,len=16,numeric"`
	CustomerFullName   string  `form:"customer_full_name" validate:"required"`
	CustomerLegalName  string  `form:"customer_legal_name" validate:"required"`
	CustomerBirthPlace string  `form:"customer_birth_place" validate:"required"`
	CustomerBirthDate  string  `form:"customer_birth_date" validate:"required,datetime=2006-01-02"`
	CustomerSalary     float64 `form:"customer_salary" validate:"required,gte=1000000,lte=100000000"`
}

func CreateCustomer(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input CustomerInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	if err := utils.Validate.Struct(input); err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload KTP photo"})
		return
	}

	selfiePhotoPath, err := utils.SaveUploadedFile(c, "customer_selfie_photo", "uploads/selfie")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload selfie photo"})
		return
	}

	customer := model.Customer{
		CustomerNIK:         input.CustomerNIK,
		CustomerFullName:    html.EscapeString(input.CustomerFullName),
		CustomerLegalName:   html.EscapeString(input.CustomerLegalName),
		CustomerBirthPlace:  html.EscapeString(input.CustomerBirthPlace),
		CustomerBirthDate:   customerBirthDate,
		CustomerSalary:      input.CustomerSalary,
		CustomerKTPPhoto:    ktpPhotoPath,
		CustomerSelfiePhoto: selfiePhotoPath,
		CustomerCreatedBy:   authUser.(model.User).UserID,
		CustomerCreatedAt:   time.Now(),
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&customer).Error; err != nil {
			if err := os.Remove(customer.CustomerKTPPhoto); err != nil {
				utils.Logger.Warnf("Failed to delete KTP photo: %s", customer.CustomerKTPPhoto)
			}

			if err := os.Remove(customer.CustomerSelfiePhoto); err != nil {
				utils.Logger.Warnf("Failed to delete Selfie photo: %s", customer.CustomerSelfiePhoto)
			}
			return err
		}
		return nil
	})

	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUser.(model.User).UserID,
			"error":   err.Error(),
		}).Error("Failed to create customer")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":      authUser.(model.User).UserID,
		"customer_nik": input.CustomerNIK,
		"created_at":   time.Now(),
	}).Info("Customer created successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Customer created successfully", "customer": customer})
}

func GetCustomer(c *gin.Context) {
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

	var customers []model.Customer
	err = config.DB.
		Preload("CreatedByUser").
		Preload("EditedByUser").
		Limit(limit).
		Offset(offset).
		Find(&customers).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"limit":     limit,
		"customers": customers,
	})
}

func GetCustomerByID(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var customer model.Customer
	err = config.DB.
		Preload("CreatedByUser").
		Preload("EditedByUser").
		First(&customer, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"customer": customer})
}

func UpdateCustomer(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(model.User)
	timeNow := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var customer model.Customer
	if err := config.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var input CustomerInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
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
			err := os.Remove(customer.CustomerKTPPhoto)
			if err != nil {
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
			err := os.Remove(customer.CustomerSelfiePhoto)
			if err != nil {
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

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&customer).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUserModel.UserID,
			"error":   err.Error(),
		}).Error("Failed to update customer")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":      authUserModel.UserID,
		"customer_nik": input.CustomerNIK,
		"created_at":   time.Now(),
	}).Infof("Customer ID %d updated successfully by User %d", customer.CustomerID, authUserModel.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully", "customer": customer})
}

func DeleteCustomer(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(model.User)
	timeNow := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var customer model.Customer
	if err := config.DB.First(&customer, id).Error; err != nil {
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

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&customer).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id":     authUserModel.UserID,
			"customer_id": customer.CustomerID,
			"error":       err.Error(),
		}).Error("Failed to delete customer")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":     authUserModel.UserID,
		"customer_id": customer.CustomerID,
		"deleted_at":  timeNow,
	}).Infof("Customer ID %d deleted successfully by User %d", customer.CustomerID, authUserModel.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
