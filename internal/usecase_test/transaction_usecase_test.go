package usecase_test

import (
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/repository/mocks"
	"kreditplus/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateTransaction_Success(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)

	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	input := domain.TransactionInput{
		TransactionNIK:         "1234567890123456",
		TransactionOTR:         1000000,
		TransactionAdminFee:    50000,
		TransactionInstallment: 12,
		TransactionInterest:    5.0,
	}

	customer := &domain.Customer{
		CustomerNIK: "1234567890123456",
	}

	limit := &domain.Limit{
		LimitID:              1,
		LimitNIK:             "1234567890123456",
		LimitRemainingAmount: 5000000, // ✅
	}

	mockLimitRepo.On("GetLimitByNIKandTenorWithTx", mock.Anything, input.TransactionNIK, input.TransactionInstallment).Return(limit, nil)
	mockTransactionRepo.On("WithTransaction", mock.Anything).Return(nil)

	err := transactionUsecase.CreateTransactionWithLimitUpdate(1, customer, input)

	assert.Nil(t, err, "Transaction creation should be successful")
	mockTransactionRepo.AssertExpectations(t)
}

func TestCreateTransaction_InsufficientLimit(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)

	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	input := domain.TransactionInput{
		TransactionNIK:         "1234567890123456",
		TransactionOTR:         5000000,
		TransactionAdminFee:    50000,
		TransactionInstallment: 12,
		TransactionInterest:    5.0,
	}

	customer := &domain.Customer{
		CustomerNIK: "1234567890123456",
	}

	limit := &domain.Limit{
		LimitID:              1,
		LimitNIK:             "1234567890123456",
		LimitRemainingAmount: 100000,
	}

	mockTransactionRepo.On("WithTransaction", mock.AnythingOfType("func(*gorm.DB) error")).
		Return(func(fn func(*gorm.DB) error) error {
			return fn(nil)
		})

	mockLimitRepo.On("GetLimitByNIKandTenorWithTx", mock.Anything, input.TransactionNIK, input.TransactionInstallment).
		Return(limit, nil)

	err := transactionUsecase.CreateTransactionWithLimitUpdate(1, customer, input)

	assert.Error(t, err, "Transaction should fail due to insufficient limit")
	assert.Equal(t, "insufficient limit", err.Error())
}

func TestCreateTransaction_DBError(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)

	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	input := domain.TransactionInput{
		TransactionNIK:         "1234567890123456",
		TransactionOTR:         1000000,
		TransactionAdminFee:    50000,
		TransactionInstallment: 12,
		TransactionInterest:    5.0,
	}

	customer := &domain.Customer{
		CustomerNIK: "1234567890123456",
	}

	mockTransactionRepo.On("WithTransaction", mock.Anything).Return(errors.New("database error"))

	err := transactionUsecase.CreateTransactionWithLimitUpdate(1, customer, input)

	assert.Error(t, err, "Error should not be nil")
	assert.Equal(t, "database error", err.Error(), "Expected database error")
}

func TestGetAllTransactions_Success(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	mockTransactions := []domain.Transaction{
		{
			TransactionID:             1,
			TransactionContractNumber: "TX12345",
			TransactionNIK:            "1234567890123456",
			TransactionOTR:            10000000,
			TransactionAdminFee:       500000,
			TransactionInstallment:    12,
			TransactionInterest:       5.0,
			TransactionAssetName:      "Motorcycle",
		},
		{
			TransactionID:             2,
			TransactionContractNumber: "TX67890",
			TransactionNIK:            "9876543210987654",
			TransactionOTR:            15000000,
			TransactionAdminFee:       700000,
			TransactionInstallment:    24,
			TransactionInterest:       4.5,
			TransactionAssetName:      "Car",
		},
	}

	mockTransactionRepo.On("GetAllTransactions", 10, 0).Return(mockTransactions, nil)

	transactions, err := transactionUsecase.GetAllTransactions(10, 0)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, transactions, "Transactions should not be nil")
	assert.Len(t, transactions, 2, "Should return 2 transactions")
	assert.Equal(t, "TX12345", transactions[0].TransactionContractNumber)
	assert.Equal(t, "TX67890", transactions[1].TransactionContractNumber)
}

func TestGetAllTransactions_DBError(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	mockTransactionRepo.On("GetAllTransactions", 10, 0).Return(nil, errors.New("database error"))

	transactions, err := transactionUsecase.GetAllTransactions(10, 0)

	assert.Error(t, err, "Error should not be nil")
	assert.Nil(t, transactions, "Transactions should be nil on DB error")
	assert.Equal(t, "database error", err.Error(), "Expected 'database error'")
}

func TestGetTransactionByID_Success(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)

	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	mockTransaction := &domain.Transaction{
		TransactionID:             1,
		TransactionContractNumber: "TX12345",
		TransactionNIK:            "1234567890123456",
		TransactionOTR:            15000000,
		TransactionAdminFee:       500000,
		TransactionInstallment:    12,
		TransactionInterest:       5.0,
		TransactionAssetName:      "Car",
		TransactionDate:           time.Now(),
	}

	mockTransactionRepo.On("GetTransactionByID", uint(1)).Return(mockTransaction, nil)

	transaction, err := transactionUsecase.GetTransactionByID(1)

	assert.Nil(t, err, "Error should be nil when transaction exists")
	assert.NotNil(t, transaction, "Transaction should not be nil")
	assert.Equal(t, uint(1), transaction.TransactionID, "Transaction ID should match expected value")
	assert.Equal(t, "TX12345", transaction.TransactionContractNumber, "Transaction contract number should match expected value")
}

func TestGetTransactionByID_NotFound(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)

	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	mockTransactionRepo.On("GetTransactionByID", uint(999)).Return(nil, errors.New("transaction not found"))

	transaction, err := transactionUsecase.GetTransactionByID(999)

	assert.Nil(t, transaction, "Transaction should be nil when not found")
	assert.NotNil(t, err, "Error should not be nil for a non-existing transaction")
	assert.Equal(t, "transaction not found", err.Error(), "Error message should indicate transaction not found")
}

func TestGetCustomerByNIKForTransaction_Success(t *testing.T) {
	mockCustomerRepo := new(mocks.CustomerRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockTransactionRepo := new(mocks.TransactionRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	expectedCustomer := &domain.Customer{
		CustomerNIK:      "1234567890123456",
		CustomerFullName: "John Doe",
	}

	mockCustomerRepo.On("GetCustomerByNIK", "1234567890123456").Return(expectedCustomer, nil)

	customer, err := transactionUsecase.GetCustomerByNIK("1234567890123456")

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, customer, "Customer should not be nil")
	assert.Equal(t, "1234567890123456", customer.CustomerNIK)
	assert.Equal(t, "John Doe", customer.CustomerFullName)
}

func TestGetCustomerByNIKForTransaction_NotFound(t *testing.T) {
	mockCustomerRepo := new(mocks.CustomerRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockTransactionRepo := new(mocks.TransactionRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	mockCustomerRepo.On("GetCustomerByNIK", "0000000000000000").Return(nil, gorm.ErrRecordNotFound)

	customer, err := transactionUsecase.GetCustomerByNIK("0000000000000000")

	assert.Nil(t, customer, "Customer should be nil")
	assert.Error(t, err, "Error should not be nil")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Should return record not found error")
}

func TestGetCustomerByNIKForTransaction_DBError(t *testing.T) {
	mockCustomerRepo := new(mocks.CustomerRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockTransactionRepo := new(mocks.TransactionRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	mockCustomerRepo.On("GetCustomerByNIK", "1234567890123456").Return(nil, errors.New("database error"))

	customer, err := transactionUsecase.GetCustomerByNIK("1234567890123456")

	assert.Nil(t, customer, "Customer should be nil")
	assert.Error(t, err, "Error should not be nil")
	assert.Equal(t, "database error", err.Error(), "Should return database error")
}

func TestUpdateTransaction_Success(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	userID := uint(1)

	mockTransaction := &domain.Transaction{
		TransactionContractNumber: "TX12345",
		TransactionNIK:            "1234567890123456",
		TransactionOTR:            10000000,
		TransactionInstallment:    12,
		TransactionLimit:          1,
	}

	mockCustomer := &domain.Customer{
		CustomerNIK: "1234567890123456",
	}

	mockLimit := &domain.Limit{
		LimitID:              1,
		LimitNIK:             "1234567890123456",
		LimitTenor:           12,
		LimitAmount:          20000000,
		LimitRemainingAmount: 10000000,
	}

	input := domain.TransactionInput{
		TransactionNIK:         "1234567890123456",
		TransactionOTR:         8000000,
		TransactionInstallment: 12,
	}

	mockTransactionRepo.On("WithTransaction", mock.Anything).Return(nil)
	mockLimitRepo.On("GetLimitByNIKandTenorWithTx", mock.Anything, input.TransactionNIK, input.TransactionInstallment).Return(mockLimit, nil)
	mockTransactionRepo.On("UpdateTransactionWithTx", mock.Anything, mockTransaction).Return(nil)

	err := transactionUsecase.UpdateTransactionWithLimitUpdate(userID, mockCustomer, mockTransaction, input)

	assert.Nil(t, err, "Transaction should be updated successfully")
}

func TestUpdateTransaction_InsufficientLimit(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	userID := uint(1)

	mockTransaction := &domain.Transaction{
		TransactionContractNumber: "TX12345",
		TransactionNIK:            "1234567890123456",
		TransactionOTR:            10000000,
		TransactionInstallment:    12,
		TransactionLimit:          1,
	}

	mockCustomer := &domain.Customer{
		CustomerNIK: "1234567890123456",
	}

	mockLimit := &domain.Limit{
		LimitID:              1,
		LimitNIK:             "1234567890123456",
		LimitTenor:           12,
		LimitAmount:          20000000,
		LimitRemainingAmount: 5000000,
	}

	input := domain.TransactionInput{
		TransactionNIK:         "1234567890123456",
		TransactionOTR:         15000000,
		TransactionInstallment: 12,
	}

	mockTransactionRepo.On("WithTransaction", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(tx *gorm.DB) error)
		_ = fn(nil)
	}).Return(errors.New("insufficient limit"))

	mockLimitRepo.On("GetLimitByNIKandTenorWithTx", mock.Anything, input.TransactionNIK, input.TransactionInstallment).Return(mockLimit, nil)

	mockLimitRepo.On("UpdateLimitWithTx", mock.Anything, mock.Anything).Return(nil).Maybe()

	mockTransactionRepo.On("UpdateTransactionWithTx", mock.Anything, mock.Anything).Return(nil).Maybe()

	err := transactionUsecase.UpdateTransactionWithLimitUpdate(userID, mockCustomer, mockTransaction, input)

	assert.Error(t, err)
	assert.Equal(t, "insufficient limit", err.Error(), "Should return insufficient limit error")

	mockTransactionRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
}

func TestUpdateTransaction_DBError(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	userID := uint(1)

	mockTransaction := &domain.Transaction{
		TransactionContractNumber: "TX12345",
		TransactionNIK:            "1234567890123456",
		TransactionOTR:            10000000,
		TransactionInstallment:    12,
		TransactionLimit:          1,
	}

	mockCustomer := &domain.Customer{
		CustomerNIK: "1234567890123456",
	}

	mockLimit := &domain.Limit{
		LimitID:              1,
		LimitNIK:             "1234567890123456",
		LimitTenor:           12,
		LimitAmount:          20000000,
		LimitRemainingAmount: 15000000,
	}

	input := domain.TransactionInput{
		TransactionNIK:         "1234567890123456",
		TransactionOTR:         12000000,
		TransactionInstallment: 12,
	}

	mockTransactionRepo.On("WithTransaction", mock.AnythingOfType("func(*gorm.DB) error")).
		Return(func(fn func(*gorm.DB) error) error {
			return fn(nil)
		})

	mockLimitRepo.On("GetLimitByNIKandTenorWithTx", mock.Anything, input.TransactionNIK, input.TransactionInstallment).
		Return(mockLimit, nil)

	mockTransactionRepo.On("UpdateTransactionWithTx", mock.Anything, mockTransaction).
		Return(errors.New("database error"))

	err := transactionUsecase.UpdateTransactionWithLimitUpdate(userID, mockCustomer, mockTransaction, input)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error(), "Should return database error")

	mockTransactionRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
}

func TestDeleteTransaction_Success(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	mockTransaction := &domain.Transaction{
		TransactionID:             1,
		TransactionContractNumber: "TX12345",
		TransactionLimit:          1,
		TransactionOTR:            1000000,
		TransactionAdminFee:       50000,
		TransactionInstallment:    12,
		TransactionInterest:       5.0,
	}

	mockLimit := &domain.Limit{
		LimitID:              1,
		LimitNIK:             "1234567890123456",
		LimitTenor:           12,
		LimitAmount:          20000000,
		LimitUsedAmount:      1000000,
		LimitRemainingAmount: 19000000,
	}

	mockTransactionRepo.On("WithTransaction", mock.Anything).Return(nil)
	mockLimitRepo.On("GetLimitByIDWithTx", mock.Anything, mockTransaction.TransactionLimit).Return(mockLimit, nil)
	mockTransactionRepo.On("DeleteTransactionWithTx", mock.Anything, mockTransaction).Return(nil)
	mockLimitRepo.On("UpdateLimitWithTx", mock.Anything, mockLimit).Return(nil)

	err := transactionUsecase.DeleteTransactionWithLimitUpdate(1, mockTransaction)

	assert.Nil(t, err, "Transaction should be deleted successfully")
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	mockTransactionRepo.On("WithTransaction", mock.Anything).Return(errors.New("transaction not found"))

	err := transactionUsecase.DeleteTransactionWithLimitUpdate(1, &domain.Transaction{TransactionID: 999})

	assert.Error(t, err)
	assert.Equal(t, "transaction not found", err.Error(), "Should return transaction not found error")
}

func TestDeleteTransaction_DBError(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLimitRepo := new(mocks.LimitRepository)
	mockCustomerRepo := new(mocks.CustomerRepository)
	transactionUsecase := usecase.NewTransactionUsecase(mockCustomerRepo, mockLimitRepo, mockTransactionRepo)

	mockTransaction := &domain.Transaction{
		TransactionID:             1,
		TransactionContractNumber: "TX12345",
		TransactionLimit:          1,
		TransactionOTR:            1000000,
		TransactionAdminFee:       50000,
		TransactionInstallment:    12,
		TransactionInterest:       5.0,
	}

	mockTransactionRepo.On("WithTransaction", mock.Anything).Return(errors.New("database error"))

	err := transactionUsecase.DeleteTransactionWithLimitUpdate(1, mockTransaction)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error(), "Should return database error")
}
