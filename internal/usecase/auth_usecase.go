package usecase

import (
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/middleware"
	"kreditplus/internal/repository"
	"kreditplus/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(input domain.LoginInput) (*domain.AuthResponse, error)
	RefreshToken(refreshToken string) (*domain.AuthResponse, error)
}

type authUsecase struct {
	repo repository.AuthRepository
}

func NewAuthUsecase(repo repository.AuthRepository) AuthUsecase {
	return &authUsecase{repo}
}

func (u *authUsecase) Login(input domain.LoginInput) (*domain.AuthResponse, error) {
	input.Username = utils.SanitizeString(input.Username)
	input.Password = utils.SanitizeString(input.Password)

	user, err := u.repo.FindUserByUsername(input.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	accessToken, refreshToken, err := middleware.GenerateTokens(user.UserID, user.UserRole)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Login successful",
	}, nil
}

func (u *authUsecase) RefreshToken(refreshToken string) (*domain.AuthResponse, error) {
	_, claims, err := middleware.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID := uint(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	accessToken, refreshToken, err := middleware.GenerateTokens(userID, userRole)
	if err != nil {
		return nil, errors.New("failed to generate new token")
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "New access token generated",
	}, nil
}

