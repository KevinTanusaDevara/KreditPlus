package handler

import (
	"kreditplus/internal/domain"
	"kreditplus/internal/middleware"
	"kreditplus/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(usecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input domain.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := h.usecase.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	csrfToken, err := middleware.GenerateCSRFToken()
	if err != nil {
		logrus.Error("Failed to generate CSRF token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
		return
	}

	c.SetCookie("csrf_token", csrfToken, 3600, "/", "", false, true)
	c.SetCookie("access_token", response.AccessToken, 15*60, "/", "", false, true)
	c.SetCookie("refresh_token", response.RefreshToken, 7*24*60*60, "/", "", false, true)

	c.Header("X-CSRF-Token", csrfToken)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	response, err := h.usecase.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	csrfToken, err := middleware.GenerateCSRFToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
		return
	}

	c.SetCookie("access_token", response.AccessToken, 15*60, "/", "", false, true)
	c.SetCookie("csrf_token", csrfToken, 3600, "/", "", false, true)
	c.Header("X-CSRF-Token", csrfToken)

	c.JSON(http.StatusOK, gin.H{"message": "Refresh Token successful"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
