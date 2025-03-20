package usecase_test

import (
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/repository/mocks"
	"kreditplus/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateCustomer_Success(t *testing.T) {
	customerRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)

	input := domain.Customer{
		CustomerNIK:        "1234567890123456",
		CustomerFullName:   "Test User",
		CustomerLegalName:  "Test Legal",
		CustomerBirthPlace: "Jakarta",
		CustomerSalary:     5000000,
	}

	customerRepo.On("CreateCustomer", mock.Anything).Return(nil)

	err := customerUsecase.CreateCustomer(input)

	assert.Nil(t, err)
	customerRepo.AssertExpectations(t)
}

func TestCreateCustomer_InvalidNIK(t *testing.T) {
	customerRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)

	input := domain.Customer{
		CustomerNIK:        "123",
		CustomerFullName:   "Test User",
		CustomerLegalName:  "Test Legal",
		CustomerBirthPlace: "Jakarta",
		CustomerSalary:     5000000,
	}

	err := customerUsecase.CreateCustomer(input)

	assert.NotNil(t, err)
	assert.Equal(t, "invalid NIK", err.Error())
}

func TestCreateCustomer_DBError(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	customer := domain.Customer{
		CustomerNIK:      "1234567890123456",
		CustomerFullName: "John Doe",
		CustomerSalary:   5000000,
	}

	mockRepo.On("CreateCustomer", mock.Anything).Return(errors.New("database error"))

	err := customerUsecase.CreateCustomer(customer)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error(), "Should return database error")
}

func TestGetAllCustomers_Success(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	mockCustomers := []domain.Customer{
		{CustomerNIK: "1234567890123456", CustomerFullName: "John Doe"},
		{CustomerNIK: "6543210987654321", CustomerFullName: "Jane Doe"},
	}

	mockRepo.On("GetAllCustomers", 10, 0).Return(mockCustomers, nil)

	customers, err := customerUsecase.GetAllCustomers(10, 0)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, customers, "Customers should not be nil")
	assert.Len(t, customers, 2, "Should return 2 customers")
	assert.Equal(t, "John Doe", customers[0].CustomerFullName)
	assert.Equal(t, "Jane Doe", customers[1].CustomerFullName)

	mockRepo.AssertExpectations(t)
}

func TestGetAllCustomers_DBError(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	mockRepo.On("GetAllCustomers", 10, 0).Return(nil, errors.New("database error"))

	customers, err := customerUsecase.GetAllCustomers(10, 0)

	assert.Nil(t, customers, "Customers should be nil on DB error")
	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, "database error", err.Error(), "Expected database error")

	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByID_Success(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	mockCustomer := &domain.Customer{
		CustomerID:       1,
		CustomerNIK:      "1234567890123456",
		CustomerFullName: "John Doe",
	}

	mockRepo.On("GetCustomerByID", uint(1)).Return(mockCustomer, nil)

	customer, err := customerUsecase.GetCustomerByID(1)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, customer, "Customer should not be nil")
	assert.Equal(t, uint(1), customer.CustomerID)
	assert.Equal(t, "1234567890123456", customer.CustomerNIK)
	assert.Equal(t, "John Doe", customer.CustomerFullName)

	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	mockRepo.On("GetCustomerByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	customer, err := customerUsecase.GetCustomerByID(99)

	assert.Nil(t, customer, "Customer should be nil when not found")
	assert.Error(t, err, "Error should be returned when customer is not found")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Expected gorm.ErrRecordNotFound error")

	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByID_DBError(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	mockRepo.On("GetCustomerByID", uint(1)).Return(nil, errors.New("database error"))

	customer, err := customerUsecase.GetCustomerByID(1)

	assert.Nil(t, customer, "Customer should be nil on DB error")
	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, "database error", err.Error(), "Expected 'database error' message")

	mockRepo.AssertExpectations(t)
}

func TestUpdateCustomer_Success(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	updatedCustomer := domain.Customer{
		CustomerNIK:      "1234567890123456",
		CustomerFullName: "John Doe Updated",
		CustomerSalary:   6000000,
	}

	mockRepo.On("UpdateCustomer", mock.Anything).Return(nil)

	err := customerUsecase.UpdateCustomer(updatedCustomer)

	assert.Nil(t, err, "Customer update should be successful")
	mockRepo.AssertExpectations(t)
}

func TestUpdateCustomer_InvalidNIK(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	customer := domain.Customer{
		CustomerNIK: "12345",
	}

	err := customerUsecase.UpdateCustomer(customer)

	assert.Error(t, err, "Error should occur due to invalid NIK")
	assert.Equal(t, "NIK must be 16 numeric characters", err.Error(), "Should return NIK validation error")
}

func TestUpdateCustomer_DBError(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	customer := domain.Customer{
		CustomerNIK:      "1234567890123456",
		CustomerFullName: "John Doe Updated",
	}

	mockRepo.On("UpdateCustomer", mock.Anything).Return(errors.New("database error"))

	err := customerUsecase.UpdateCustomer(customer)

	assert.Error(t, err, "Error should occur due to database failure")
	assert.Equal(t, "database error", err.Error(), "Should return database error")
}

func TestDeleteCustomer_Success(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	mockRepo.On("DeleteCustomer", uint(1)).Return(nil)

	err := customerUsecase.DeleteCustomer(1)

	assert.Nil(t, err, "Customer deletion should be successful")
	mockRepo.AssertExpectations(t)
}

func TestDeleteCustomer_NotFound(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	mockRepo.On("DeleteCustomer", uint(99)).Return(gorm.ErrRecordNotFound)

	err := customerUsecase.DeleteCustomer(99)

	assert.Error(t, err, "Expected an error for non-existing customer")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Should return record not found error")
	mockRepo.AssertExpectations(t)
}

func TestDeleteCustomer_DBError(t *testing.T) {
	mockRepo := new(mocks.CustomerRepository)
	customerUsecase := usecase.NewCustomerUsecase(mockRepo)

	mockRepo.On("DeleteCustomer", uint(1)).Return(errors.New("database error"))

	err := customerUsecase.DeleteCustomer(1)

	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, "database error", err.Error(), "Expected database error")
	mockRepo.AssertExpectations(t)
}
