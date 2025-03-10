package controller

import (
	"errors"
	"kreditplus/config"
	"kreditplus/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(model.User)

	if authUserModel.UserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to use this function"})
		return
	}

	var input struct {
		UserUsername string `json:"user_username"`
		UserPassword string `json:"user_password"`
		UserRole     string `json:"user_role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if input.UserRole == "" {
		input.UserRole = "user"
	}

	var user model.User

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_username = ?", input.UserUsername).First(&user).Error; err == nil {
			return errors.New("username already exists")
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}

		newUser := model.User{
			UserUsername: input.UserUsername,
			UserPassword: string(hashedPassword),
			UserRole:     input.UserRole,
		}
		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func GetUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(model.User)

	if authUserModel.UserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to use this function"})
		return
	}

	var users []model.User
	var userDTOs []model.UserResponseDTO

	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page value"})
		return
	}

	offset := (page - 1) * limit

	config.DB.Limit(limit).Offset(offset).Find(&users)

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
	authUserModel := authUser.(model.User)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if authUserModel.UserRole != "admin" && authUserModel.UserID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view this user"})
		return
	}

	userDTO := user.ToDTO()

	c.JSON(http.StatusOK, userDTO)
}

func UpdateUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(model.User)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if authUserModel.UserRole != "admin" && authUserModel.UserID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only update your own profile"})
		return
	}

	var input struct {
		UserUsername string `json:"user_username,omitempty"`
		UserPassword string `json:"user_password,omitempty"`
		UserRole     string `json:"user_role,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
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
		if authUserModel.UserRole == "admin" {
			user.UserRole = input.UserRole
		} else if authUserModel.UserID == user.UserID {
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
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(model.User)

	if authUserModel.UserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to use this function"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
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
