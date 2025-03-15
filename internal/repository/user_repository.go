package repository

import (
	"kreditplus/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetAllUsers(limit, offset int) ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetAllUsers(limit, offset int) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Limit(limit).
		Offset(offset).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("user_username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}
