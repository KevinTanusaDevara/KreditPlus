package middleware

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type MockMiddlewareService struct{}

func (m *MockMiddlewareService) GenerateCSRFToken() (string, error) {
	return "mock_csrf_token", nil
}

func (m *MockMiddlewareService) ValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	if tokenString == "valid_refresh_token" {
		return nil, jwt.MapClaims{
			"user_id":   float64(1),
			"user_role": "user",
		}, nil
	}
	return nil, nil, errors.New("invalid refresh token")
}

func (m *MockMiddlewareService) GenerateTokens(userID uint, userRole string) (string, string, error) {
	if userID == 1 {
		return "new_access_token", "new_refresh_token", nil
	}
	return "", "", errors.New("failed to generate new token")
}
