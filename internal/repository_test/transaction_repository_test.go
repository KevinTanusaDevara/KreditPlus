package repository_test

import (
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestWithTransaction_Success(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err = transactionRepo.WithTransaction(func(tx *gorm.DB) error {
		return nil
	})

	assert.Nil(t, err, "Transaction should complete successfully")
}

func TestWithTransaction_DeadlockRetry(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	// Simulate deadlock on the first two attempts
	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errors.New("deadlock detected"))

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errors.New("deadlock detected"))

	// Third attempt should succeed
	mock.ExpectBegin()
	mock.ExpectCommit()

	// Counter to track the number of attempts
	attempt := 0

	// Execute the function under test
	err = transactionRepo.WithTransaction(func(tx *gorm.DB) error {
		attempt++
		if attempt < 3 {
			return errors.New("deadlock detected") // Simulate deadlock for the first two attempts
		}
		return nil // Success on the third attempt
	})

	// Assertions
	assert.Nil(t, err, "Transaction should succeed after retrying deadlock")
	assert.Equal(t, 3, attempt, "Transaction should be retried 3 times")

	// Ensure all mock expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestWithTransaction_Failure(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = transactionRepo.WithTransaction(func(tx *gorm.DB) error {
		return errors.New("failed transaction")
	})

	assert.Error(t, err, "Transaction should fail and rollback")
}

func TestCreateTransactionWithTx_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create SQL mock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO "transactions"`).
		WithArgs(
			"TRX-001",
			"1234567890123456",
			0,
			1000000.0,
			50000.0,
			200000.0,
			10000.0,
			"Car",
			sqlmock.AnyArg(),
			0,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"transaction_id"}).AddRow(1))

	mock.ExpectCommit()

	transaction := &domain.Transaction{
		TransactionContractNumber: "TRX-001",
		TransactionNIK:            "1234567890123456",
		TransactionOTR:            1000000.0,
		TransactionAdminFee:       50000.0,
		TransactionInstallment:    200000.0,
		TransactionInterest:       10000.0,
		TransactionAssetName:      "Car",
		TransactionDate:           time.Now(),
	}

	tx := gormDB.Begin()

	err = transactionRepo.CreateTransactionWithTx(tx, transaction)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to create transaction: %v", err)
	}

	err = tx.Commit().Error
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	assert.Nil(t, err, "Error should be nil on successful transaction insert")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestCreateTransactionWithTx_DBError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create SQL mock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO "transactions"`).
		WithArgs(
			"TRX-001",
			"1234567890123456",
			0,
			1000000.0,
			50000.0,
			200000.0,
			10000.0,
			"Car",
			sqlmock.AnyArg(),
			0,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(gorm.ErrInvalidTransaction)

	mock.ExpectRollback()

	transaction := &domain.Transaction{
		TransactionContractNumber: "TRX-001",
		TransactionNIK:            "1234567890123456",
		TransactionOTR:            1000000.0,
		TransactionAdminFee:       50000.0,
		TransactionInstallment:    200000.0,
		TransactionInterest:       10000.0,
		TransactionAssetName:      "Car",
		TransactionDate:           time.Now(),
	}

	tx := gormDB.Begin()

	err = transactionRepo.CreateTransactionWithTx(tx, transaction)

	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetAllTransactions_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "transactions" LIMIT \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "transaction_nik", "transaction_amount"}).
			AddRow(1, "1234567890123456", 5000000).
			AddRow(2, "9876543210987654", 3000000))

	mock.ExpectQuery(`SELECT \* FROM "customers" WHERE "customers"."customer_nik" IN \(\$1,\$2\)`).
		WithArgs("1234567890123456", "9876543210987654").
		WillReturnRows(sqlmock.NewRows([]string{"customer_nik", "customer_name"}).
			AddRow("1234567890123456", "John Doe").
			AddRow("9876543210987654", "Jane Doe"))

	transactions, err := transactionRepo.GetAllTransactions(10, 0)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, transactions, "Transactions should not be nil")
	assert.Len(t, transactions, 2, "Should return 2 transactions")
	assert.Equal(t, "1234567890123456", transactions[0].TransactionNIK)
	assert.Equal(t, "9876543210987654", transactions[1].TransactionNIK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetAllTransactions_DBError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "transactions" LIMIT \$1`).
		WithArgs(10).
		WillReturnError(gorm.ErrInvalidTransaction)

	transactions, err := transactionRepo.GetAllTransactions(10, 0)

	assert.Nil(t, transactions, "Transactions should be nil on DB error")
	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionByID_Success(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "transactions" WHERE "transactions"."transaction_id" = \$1 ORDER BY "transactions"."transaction_id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"transaction_id", "transaction_nik", "transaction_otr", "transaction_admin_fee", "transaction_installment",
		}).
			AddRow(1, "1234567890123456", 5000000.0, 100000, 250000))

	mock.ExpectQuery(`SELECT \* FROM "customers" WHERE "customers"."customer_nik" = \$1`).
		WithArgs("1234567890123456").
		WillReturnRows(sqlmock.NewRows([]string{"customer_nik", "customer_full_name"}).
			AddRow("1234567890123456", "John Doe"))

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."user_id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "user_username"}).
			AddRow(1, "admin"))

	transaction, err := transactionRepo.GetTransactionByID(1)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, transaction, "Transaction should not be nil")
	assert.Equal(t, "1234567890123456", transaction.TransactionNIK)
	assert.Equal(t, 5000000.0, transaction.TransactionOTR)
}

func TestGetTransactionByID_NotFound(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "transactions" WHERE "transactions"."transaction_id" = \$1 ORDER BY "transactions"."transaction_id" LIMIT \$2`).
		WithArgs(99, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	transaction, err := transactionRepo.GetTransactionByID(99)

	assert.Nil(t, transaction, "Transaction should be nil when not found")
	assert.Error(t, err, "Error should be returned for non-existing transaction")
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionByID_DBError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "transactions" WHERE "transactions"."transaction_id" = \$1 ORDER BY "transactions"."transaction_id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(errors.New("database connection failed"))

	transaction, err := transactionRepo.GetTransactionByID(1)

	assert.Nil(t, transaction, "Transaction should be nil on DB error")
	assert.Error(t, err, "Error should be returned on DB failure")
	assert.Equal(t, "database connection failed", err.Error())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestUpdateTransactionWithTx_Success(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "transactions" SET`).
		WithArgs(
			"",
			"1234567890123456",
			int64(0),
			5000000.0,
			100000.0,
			500000.0,
			5.0,
			"Motorcycle",
			sqlmock.AnyArg(),
			int64(0),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			int64(1),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	tx := gormDB.Begin()

	transaction := &domain.Transaction{
		TransactionID:             1,
		TransactionNIK:            "1234567890123456",
		TransactionLimit:          0,
		TransactionOTR:            5000000.0,
		TransactionAdminFee:       100000.0,
		TransactionInstallment:    500000.0,
		TransactionInterest:       5.0,
		TransactionAssetName:      "Motorcycle",
		TransactionContractNumber: "",
		TransactionDate:           time.Now(),
		TransactionCreatedBy:      0,
		TransactionCreatedAt:      time.Now(),
		TransactionEditedBy:       nil,
		TransactionEditedAt:       nil,
	}

	err = transactionRepo.UpdateTransactionWithTx(tx, transaction)

	assert.Nil(t, err, "Error should be nil on successful update")
}

func TestUpdateTransactionWithTx_DBError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "transactions" SET`).
		WithArgs(
			"",
			"1234567890123456",
			int64(0),
			5000000.0,
			100000.0,
			500000.0,
			5.0,
			"Motorcycle",
			sqlmock.AnyArg(),
			int64(0),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			int64(1),
		).
		WillReturnError(gorm.ErrInvalidTransaction)

	mock.ExpectRollback()

	tx := gormDB.Begin()

	transaction := &domain.Transaction{
		TransactionID:             1,
		TransactionNIK:            "1234567890123456",
		TransactionLimit:          0,
		TransactionOTR:            5000000.0,
		TransactionAdminFee:       100000.0,
		TransactionInstallment:    500000.0,
		TransactionInterest:       5.0,
		TransactionAssetName:      "Motorcycle",
		TransactionContractNumber: "",
		TransactionDate:           time.Now(),
		TransactionCreatedBy:      0,
		TransactionCreatedAt:      time.Now(),
		TransactionEditedBy:       nil,
		TransactionEditedAt:       nil,
	}

	err = transactionRepo.UpdateTransactionWithTx(tx, transaction)

	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")
}

func TestDeleteTransactionWithTx_Success(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`DELETE FROM "transactions" WHERE "transactions"."transaction_id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	tx := gormDB.Begin()
	transaction := &domain.Transaction{TransactionID: 1}

	err = transactionRepo.DeleteTransactionWithTx(tx, transaction)

	tx.Commit()

	assert.Nil(t, err, "Error should be nil on successful delete")

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "All expected database operations should be met")
}

func TestDeleteTransactionWithTx_DBError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`DELETE FROM "transactions" WHERE "transactions"."transaction_id" = \$1`).
		WithArgs(1).
		WillReturnError(gorm.ErrInvalidTransaction)

	mock.ExpectRollback()

	tx := gormDB.Begin()
	transaction := &domain.Transaction{TransactionID: 1}

	err = transactionRepo.DeleteTransactionWithTx(tx, transaction)

	tx.Rollback()

	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "All expected database operations should be met")
}
