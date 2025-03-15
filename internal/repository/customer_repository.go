package repository

import (
	"kreditplus/internal/domain"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	CreateCustomer(customer *domain.Customer) error
	GetAllCustomers(limit, offset int) ([]domain.Customer, error)
	GetCustomerByID(id uint) (*domain.Customer, error)
	GetCustomerByNIK(nik string) (*domain.Customer, error)
	UpdateCustomer(customer *domain.Customer) error
	DeleteCustomer(id uint) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) CreateCustomer(customer *domain.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) GetAllCustomers(limit, offset int) ([]domain.Customer, error) {
	var customers []domain.Customer
	err := r.db.Preload("CreatedByUser").
		Preload("EditedByUser").
		Limit(limit).
		Offset(offset).
		Find(&customers).Error
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *customerRepository) GetCustomerByID(id uint) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.Preload("CreatedByUser").
		Preload("EditedByUser").
		First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetCustomerByNIK(nik string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.Where("customer_nik = ?", nik).
		First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) UpdateCustomer(customer *domain.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) DeleteCustomer(id uint) error {
	return r.db.Delete(&domain.Customer{}, id).Error
}
