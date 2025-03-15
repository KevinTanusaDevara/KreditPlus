package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("arise")

const AccessTokenExpiry = 15 * time.Minute
const RefreshTokenExpiry = 7 * 24 * time.Hour

func GenerateTokens(userID uint, userRole string) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"user_id":   userID,
		"user_role": userRole,
		"exp":       time.Now().Add(AccessTokenExpiry).Unix(),
	}
	refreshTokenClaims := jwt.MapClaims{
		"user_id":   userID,
		"user_role": userRole,
		"exp":       time.Now().Add(RefreshTokenExpiry).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	accessString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, errors.New("invalid token")
	}

	return token, claims, nil
}
