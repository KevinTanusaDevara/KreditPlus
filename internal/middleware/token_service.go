package middleware

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	ValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	GenerateTokens(userID uint, userRole string) (string, string, error)
}

type DefaultTokenService struct{}

func (d *DefaultTokenService) ValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, nil, err
	}

	return token, claims, nil
}

func (d *DefaultTokenService) GenerateTokens(userID uint, userRole string) (string, string, error) {
	return GenerateTokens(userID, userRole)
}
