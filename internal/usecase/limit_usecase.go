package usecase

import (
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/repository"
	"time"
)

type LimitUsecase interface {
	CreateLimit(input domain.Limit) error
	GetAllLimits(limit, offset int) ([]domain.Limit, error)
	GetLimitByID(id uint) (*domain.Limit, error)
	GetCustomerByNIK(nik string) (*domain.Customer, error)
	UpdateLimit(input domain.Limit) error
	DeleteLimit(id uint) error
}

type limitUsecase struct {
	customerRepo repository.CustomerRepository
	limitRepo    repository.LimitRepository
}

func NewLimitUsecase(limitRepo repository.LimitRepository, customerRepo repository.CustomerRepository) LimitUsecase {
	return &limitUsecase{limitRepo: limitRepo, customerRepo: customerRepo}
}

func (u *limitUsecase) CreateLimit(input domain.Limit) error {
	if input.LimitNIK == "" || len(input.LimitNIK) != 16 {
		return errors.New("invalid NIK")
	}
	input.LimitCreatedAt = time.Now()
	return u.limitRepo.CreateLimit(&input)
}

func (u *limitUsecase) GetAllLimits(limit, offset int) ([]domain.Limit, error) {
	return u.limitRepo.GetAllLimits(limit, offset)
}

func (u *limitUsecase) GetLimitByID(id uint) (*domain.Limit, error) {
	return u.limitRepo.GetLimitByID(id)
}

func (u *limitUsecase) GetCustomerByNIK(nik string) (*domain.Customer, error) {
	return u.customerRepo.GetCustomerByNIK(nik)
}

func (u *limitUsecase) UpdateLimit(input domain.Limit) error {
	if input.LimitNIK != "" && len(input.LimitNIK) != 16 {
		return errors.New("NIK must be 16 numeric characters")
	}

	input.LimitEditedAt = new(time.Time)
	*input.LimitEditedAt = time.Now()
	return u.limitRepo.UpdateLimit(&input)
}

func (u *limitUsecase) DeleteLimit(id uint) error {
	return u.limitRepo.DeleteLimit(id)
}
