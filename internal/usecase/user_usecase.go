package usecase

import (
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/repository"
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
	return u.repo.UpdateUser(&input)
}

func (u *userUsecase) DeleteUser(id uint) error {
	return u.repo.DeleteUser(id)
}