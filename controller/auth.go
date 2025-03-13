package controller

import (
	"kreditplus/config"
	"kreditplus/middleware"
	"kreditplus/model"
	"kreditplus/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_username = ?", input.Username).First(&user).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	accessToken, refreshToken, err := middleware.GenerateTokens(user.UserID, user.UserRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	var user model.User
	var newAccessToken string

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		_, claims, err := middleware.ValidateToken(refreshToken)
		if err != nil {
			return err
		}

		userID := uint(claims["user_id"].(float64))
		userRole := claims["user_role"].(string)

		if err := tx.Where("user_id = ?", userID).First(&user).Error; err != nil {
			return err
		}

		newAccessToken, _, err = middleware.GenerateTokens(userID, userRole)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.SetCookie("access_token", newAccessToken, 15*60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "New access generate successful"})
}

func Logout(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
