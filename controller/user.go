package controller

import (
	"kreditplus/config"
	"kreditplus/model"
	"kreditplus/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserInput struct {
	UserUsername string `json:"user_username" validate:"required,min=3"`
	UserPassword string `json:"user_password,omitempty" validate:"omitempty,min=6"`
	UserRole     string `json:"user_role,omitempty" validate:"omitempty,oneof=admin user"`
}

func CreateUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(model.User).UserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to use this function"})
		return
	}

	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil || utils.Validate.Struct(input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := model.User{
		UserUsername: input.UserUsername,
		UserPassword: string(hashedPassword),
		UserRole:     input.UserRole,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Logger.WithError(err).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func GetUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(model.User).UserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to use this function"})
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

	var users []model.User
	if err := config.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	var userDTOs []model.UserResponseDTO
	for _, user := range users {
		userDTOs = append(userDTOs, user.ToDTO())
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"users": userDTOs,
	})
}

func GetUserByID(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if authUser.(model.User).UserRole != "admin" && authUser.(model.User).UserID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.ToDTO()})
}

func UpdateUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if authUser.(model.User).UserRole != "admin" && authUser.(model.User).UserID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only update your own profile"})
		return
	}

	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil || utils.Validate.Struct(input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.UserUsername != "" {
		user.UserUsername = input.UserUsername
	}

	if input.UserPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.UserPassword = string(hashedPassword)
	}

	if input.UserRole != "" {
		if authUser.(model.User).UserRole == "admin" {
			user.UserRole = input.UserRole
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You cannot change your own role"})
			return
		}
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(model.User).UserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to use this function"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
