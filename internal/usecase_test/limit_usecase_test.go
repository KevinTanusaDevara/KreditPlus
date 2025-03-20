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

func TestCreateLimit_Success(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	limit := domain.Limit{
		LimitNIK:    "1234567890123456",
		LimitTenor:  12,
		LimitAmount: 5000000,
	}

	mockLimitRepo.On("CreateLimit", mock.Anything).Return(nil)

	err := limitUsecase.CreateLimit(limit)

	assert.Nil(t, err, "Limit creation should be successful")
	mockLimitRepo.AssertExpectations(t)
}

func TestCreateLimit_InvalidNIK(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	limit := domain.Limit{
		LimitNIK:    "12345",
		LimitTenor:  12,
		LimitAmount: 5000000,
	}

	err := limitUsecase.CreateLimit(limit)

	assert.Error(t, err, "Limit creation should fail due to invalid NIK")
	assert.Equal(t, "invalid NIK", err.Error(), "Expected 'invalid NIK' error message")

	mockLimitRepo.AssertNotCalled(t, "CreateLimit", mock.Anything)
}

func TestCreateLimit_DBError(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	limit := domain.Limit{
		LimitNIK:    "1234567890123456",
		LimitTenor:  12,
		LimitAmount: 5000000,
	}

	mockLimitRepo.On("CreateLimit", mock.Anything).Return(errors.New("database error"))

	err := limitUsecase.CreateLimit(limit)

	assert.Error(t, err, "Error should occur due to database failure")
	assert.Equal(t, "database error", err.Error(), "Expected database error message")
}

func TestGetAllLimits_Success(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockLimits := []domain.Limit{
		{
			LimitNIK:    "1234567890123456",
			LimitTenor:  12,
			LimitAmount: 5000000,
		},
		{
			LimitNIK:    "9876543210987654",
			LimitTenor:  24,
			LimitAmount: 10000000,
		},
	}

	mockLimitRepo.On("GetAllLimits", 10, 0).Return(mockLimits, nil)

	limits, err := limitUsecase.GetAllLimits(10, 0)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, limits, "Limits should not be nil")
	assert.Len(t, limits, 2, "Should return 2 limits")
	assert.Equal(t, "1234567890123456", limits[0].LimitNIK)
	assert.Equal(t, "9876543210987654", limits[1].LimitNIK)

	mockLimitRepo.AssertExpectations(t)
}

func TestGetAllLimits_DBError(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockLimitRepo.On("GetAllLimits", 10, 0).Return(nil, errors.New("database error"))

	limits, err := limitUsecase.GetAllLimits(10, 0)

	assert.Error(t, err, "Error should not be nil")
	assert.Nil(t, limits, "Limits should be nil on DB error")
	assert.Equal(t, "database error", err.Error(), "Expected database error")

	mockLimitRepo.AssertExpectations(t)
}

func TestGetLimitByID_Success(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockLimit := &domain.Limit{
		LimitID:     1,
		LimitNIK:    "1234567890123456",
		LimitTenor:  12,
		LimitAmount: 5000000,
	}

	mockLimitRepo.On("GetLimitByID", uint(1)).Return(mockLimit, nil)

	limit, err := limitUsecase.GetLimitByID(1)

	assert.Nil(t, err)
	assert.NotNil(t, limit)
	assert.Equal(t, uint(1), limit.LimitID)
	assert.Equal(t, "1234567890123456", limit.LimitNIK)
	mockLimitRepo.AssertExpectations(t)
}

func TestGetLimitByID_NotFound(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockLimitRepo.On("GetLimitByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	limit, err := limitUsecase.GetLimitByID(99)

	assert.Nil(t, limit)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	mockLimitRepo.AssertExpectations(t)
}

func TestGetLimitByID_DBError(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockLimitRepo.On("GetLimitByID", uint(1)).Return(nil, errors.New("database error"))

	limit, err := limitUsecase.GetLimitByID(1)

	assert.Nil(t, limit)
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockLimitRepo.AssertExpectations(t)
}

func TestGetCustomerByNIKForLimit_Success(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockCustomer := &domain.Customer{
		CustomerNIK:      "1234567890123456",
		CustomerFullName: "John Doe",
		CustomerSalary:   5000000,
	}

	mockCustomerRepo.On("GetCustomerByNIK", "1234567890123456").Return(mockCustomer, nil)

	customer, err := limitUsecase.GetCustomerByNIK("1234567890123456")

	assert.Nil(t, err)
	assert.NotNil(t, customer, "Customer should not be nil")
	assert.Equal(t, "1234567890123456", customer.CustomerNIK)
	assert.Equal(t, "John Doe", customer.CustomerFullName)
	mockCustomerRepo.AssertExpectations(t)
}

func TestGetCustomerByNIKForLimit_NotFound(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockCustomerRepo.On("GetCustomerByNIK", "9999999999999999").Return(nil, gorm.ErrRecordNotFound)

	customer, err := limitUsecase.GetCustomerByNIK("9999999999999999")

	assert.Nil(t, customer, "Customer should be nil when not found")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Error should be gorm.ErrRecordNotFound")
	mockCustomerRepo.AssertExpectations(t)
}

func TestGetCustomerByNIKForLimit_DBError(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockCustomerRepo.On("GetCustomerByNIK", "1234567890123456").Return(nil, errors.New("database error"))

	customer, err := limitUsecase.GetCustomerByNIK("1234567890123456")

	assert.Nil(t, customer, "Customer should be nil on DB error")
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error(), "Error should match expected DB error")
	mockCustomerRepo.AssertExpectations(t)
}

func TestUpdateLimit_Success(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	limit := domain.Limit{
		LimitNIK:             "1234567890123456",
		LimitTenor:           12,
		LimitAmount:          6000000,
		LimitUsedAmount:      1000000,
		LimitRemainingAmount: 5000000,
	}

	mockLimitRepo.On("UpdateLimit", mock.Anything).Return(nil)

	err := limitUsecase.UpdateLimit(limit)

	assert.Nil(t, err, "Limit update should be successful")
	mockLimitRepo.AssertExpectations(t)
}

func TestUpdateLimit_InvalidNIK(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	limit := domain.Limit{
		LimitNIK: "12345", // Invalid NIK
		LimitTenor: 12,
		LimitAmount: 6000000,
	}

	err := limitUsecase.UpdateLimit(limit)

	assert.Error(t, err)
	assert.Equal(t, "NIK must be 16 numeric characters", err.Error(), "Should return NIK validation error")
}

func TestUpdateLimit_DBError(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	limit := domain.Limit{
		LimitNIK:             "1234567890123456",
		LimitTenor:           12,
		LimitAmount:          6000000,
		LimitUsedAmount:      1000000,
		LimitRemainingAmount: 5000000,
	}

	mockLimitRepo.On("UpdateLimit", mock.Anything).Return(errors.New("database error"))

	err := limitUsecase.UpdateLimit(limit)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error(), "Should return database error")
}

func TestDeleteLimit_Success(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockLimitRepo.On("DeleteLimit", uint(1)).Return(nil)

	err := limitUsecase.DeleteLimit(1)

	assert.Nil(t, err, "Limit deletion should be successful")
	mockLimitRepo.AssertExpectations(t)
}

func TestDeleteLimit_NotFound(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockLimitRepo.On("DeleteLimit", uint(99)).Return(gorm.ErrRecordNotFound)

	err := limitUsecase.DeleteLimit(99)

	assert.Error(t, err, "Error should not be nil when deleting non-existent record")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Should return a record not found error")
}

func TestDeleteLimit_DBError(t *testing.T) {
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	limitUsecase := usecase.NewLimitUsecase(mockLimitRepo, mockCustomerRepo)

	mockLimitRepo.On("DeleteLimit", uint(1)).Return(errors.New("database error"))

	err := limitUsecase.DeleteLimit(1)

	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, "database error", err.Error(), "Expected database error")
}
