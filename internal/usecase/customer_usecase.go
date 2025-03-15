package usecase

import (
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/repository"
	"time"
)

type CustomerUsecase interface {
	CreateCustomer(input domain.Customer) error
	GetAllCustomers(limit, offset int) ([]domain.Customer, error)
	GetCustomerByID(id uint) (*domain.Customer, error)
	UpdateCustomer(input domain.Customer) error
	DeleteCustomer(id uint) error
}

type customerUsecase struct {
	repo repository.CustomerRepository
}

func NewCustomerUsecase(repo repository.CustomerRepository) CustomerUsecase {
	return &customerUsecase{repo: repo}
}

func (u *customerUsecase) CreateCustomer(input domain.Customer) error {
	if input.CustomerNIK == "" || len(input.CustomerNIK) != 16 {
		return errors.New("invalid NIK")
	}
	input.CustomerCreatedAt = time.Now()
	return u.repo.CreateCustomer(&input)
}

func (u *customerUsecase) GetAllCustomers(limit, offset int) ([]domain.Customer, error) {
	return u.repo.GetAllCustomers(limit, offset)
}

func (u *customerUsecase) GetCustomerByID(id uint) (*domain.Customer, error) {
	return u.repo.GetCustomerByID(id)
}

func (u *customerUsecase) UpdateCustomer(input domain.Customer) error {
	if input.CustomerNIK != "" && len(input.CustomerNIK) != 16 {
		return errors.New("NIK must be 16 numeric characters")
	}

	input.CustomerEditedAt = new(time.Time)
	*input.CustomerEditedAt = time.Now()
	return u.repo.UpdateCustomer(&input)
}

func (u *customerUsecase) DeleteCustomer(id uint) error {
	return u.repo.DeleteCustomer(id)
}
