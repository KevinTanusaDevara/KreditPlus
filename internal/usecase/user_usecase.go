package usecase

import (
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/repository"
	"kreditplus/internal/utils"
)

type UserUsecase interface {
	CreateUser(input domain.User) error
	GetAllUsers(limit, offset int) ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	UpdateUser(input domain.User) error
	DeleteUser(id uint) error
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) CreateUser(input domain.User) error {
	input.UserUsername = utils.SanitizeString(input.UserUsername)
	input.UserPassword = utils.SanitizeString(input.UserPassword)
	input.UserRole = utils.SanitizeString(input.UserRole)

	existingUser, _ := u.repo.GetUserByUsername(input.UserUsername)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	return u.repo.CreateUser(&input)
}

func (u *userUsecase) GetAllUsers(limit, offset int) ([]domain.User, error) {
	return u.repo.GetAllUsers(limit, offset)
}

func (u *userUsecase) GetUserByID(id uint) (*domain.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *userUsecase) UpdateUser(input domain.User) error {
	input.UserUsername = utils.SanitizeString(input.UserUsername)
	input.UserPassword = utils.SanitizeString(input.UserPassword)
	input.UserRole = utils.SanitizeString(input.UserRole)

	existingUser, err := u.repo.GetUserByID(input.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if input.UserUsername != existingUser.UserUsername {
		userWithSameUsername, _ := u.repo.GetUserByUsername(input.UserUsername)
		if userWithSameUsername != nil {
			return errors.New("username already exists")
		}
	}

	return u.repo.UpdateUser(&input)
}

func (u *userUsecase) DeleteUser(id uint) error {
	return u.repo.DeleteUser(id)
}
