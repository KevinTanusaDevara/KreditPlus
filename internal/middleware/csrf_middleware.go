package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GenerateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/auth/login") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/auth/refresh-token") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/auth/logout") {
			c.Next()
			return
		}

		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodDelete {
			clientToken := c.GetHeader("X-CSRF-Token")
			cookieToken, err := c.Cookie("csrf_token")

			if err != nil || clientToken == "" || clientToken != cookieToken {
				logrus.Warn("CSRF validation failed: Invalid token")
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
				c.Abort()
				return
			}

			logrus.Info("CSRF Token Validated Successfully")
		}

		c.Next()
	}
}
