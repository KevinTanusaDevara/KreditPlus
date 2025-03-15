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
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(domain.User).UserRole != "admin" {
		utils.Logger.Warn("Unauthorized access attempt to CreateUser")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input domain.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Logger.Warn("Invalid request format for creating user")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.Logger.Warnf("Validation error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Warnf("Failed to hash password: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := domain.User{
		UserUsername: input.UserUsername,
		UserPassword: string(hashedPassword),
		UserRole:     input.UserRole,
	}

	err = h.usecase.CreateUser(user)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUser.(domain.User).UserID,
			"error":   err.Error(),
		}).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":       authUser.(domain.User).UserID,
		"user_username": user.UserUsername,
		"created_at":    time.Now(),
	}).Infof("User ID %d created successfully by User %d", user.UserID, authUser.(domain.User).UserID)
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(domain.User).UserRole != "admin" {
		utils.Logger.Warn("Unauthorized access attempt to GetUser")
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

	users, err := h.usecase.GetAllUsers(limit, offset)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"limit":  limit,
			"offset": offset,
			"error":  err.Error(),
		}).Error("Failed to retrieve users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"page":  page,
		"limit": limit,
		"users": users,
	}).Info("Users retrieved successfully")

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"users": users,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(domain.User).UserRole != "admin" {
		utils.Logger.Warn("Unauthorized access attempt to GetUserByID")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Logger.Warn("Invalid user ID in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.usecase.GetUserByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": id,
			"error":   err.Error(),
		}).Warn("User not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if authUser.(domain.User).UserRole != "admin" && authUser.(domain.User).UserID != user.UserID {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUser.(domain.User).UserID,
		}).Info("You do not have permission to view this user")
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view this user"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User retrieved successfully")

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(domain.User).UserRole != "admin" {
		utils.Logger.Warn("Unauthorized access attempt to UpdateUser")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(domain.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Logger.Warn("Invalid user ID provided for update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.usecase.GetUserByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": id,
			"error":   err.Error(),
		}).Warn("User not found for update")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if authUser.(domain.User).UserRole != "admin" && authUser.(domain.User).UserID != user.UserID {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUser.(domain.User).UserID,
		}).Info("You do not have permission to view this user")
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view this user"})
		return
	}

	var input domain.UserInput
	if err := c.ShouldBind(&input); err != nil {
		utils.Logger.Warn("Invalid request format for updating user")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.Logger.Warnf("Validation error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.UserUsername != "" {
		user.UserUsername = input.UserUsername
	}

	if input.UserPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			utils.Logger.Warn("Failed to hash password")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.UserPassword = string(hashedPassword)
	}

	if input.UserRole != "" {
		if authUser.(domain.User).UserRole == "admin" {
			user.UserRole = input.UserRole
		} else {
			utils.Logger.Warn("You cannot change your own role")
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You cannot change your own role"})
			return
		}
	}

	err = h.usecase.UpdateUser(*user)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUserModel.UserID,
			"error":   err.Error(),
		}).Error("Failed to update user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":        authUserModel.UserID,
		"update_user_id": user.UserID,
	}).Infof("User ID %d updated successfully by User %d", user.UserID, authUserModel.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists || authUser.(domain.User).UserRole != "admin" {
		utils.Logger.Warn("Unauthorized access attempt to DeleteUser")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserModel := authUser.(domain.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.usecase.GetUserByID(uint(id))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": id,
			"error":   err.Error(),
		}).Warn("User not found for update")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = h.usecase.DeleteUser(user.UserID)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"user_id": authUserModel.UserID,
			"error":   err.Error(),
		}).Error("Failed to delete user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":        authUserModel.UserID,
		"delete_user_id": user.UserID,
	}).Infof("User ID %d deleted successfully by User %d", user.UserID, authUserModel.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
