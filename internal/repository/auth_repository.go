package repository

import (
	"kreditplus/internal/domain"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByUsername(username string) (*domain.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) FindUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("user_username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
