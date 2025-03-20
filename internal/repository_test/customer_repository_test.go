package repository_test

import (
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/repository"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateCustomer_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO "customers" \("customer_nik","customer_full_name","customer_legal_name","customer_birth_place","customer_birth_date","customer_salary","customer_ktp_photo","customer_selfie_photo","customer_created_by","customer_created_at","customer_edited_by","customer_edited_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11,\$12\) RETURNING "customer_id"`).
		WithArgs(
			"1234567890123456",
			"John Doe",
			"John D",
			"Jakarta",
			sqlmock.AnyArg(),
			float64(0),
			"",
			"",
			0,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"customer_id"}).AddRow(1))

	mock.ExpectCommit()

	customer := &domain.Customer{
		CustomerNIK:         "1234567890123456",
		CustomerFullName:    "John Doe",
		CustomerLegalName:   "John D",
		CustomerBirthPlace:  "Jakarta",
		CustomerBirthDate:   time.Now(),
		CustomerSalary:      0,
		CustomerKTPPhoto:    "",
		CustomerSelfiePhoto: "",
		CustomerCreatedBy:   0,
		CustomerCreatedAt:   time.Now(),
		CustomerEditedBy:    nil,
		CustomerEditedAt:    nil,
	}

	err = customerRepo.CreateCustomer(customer)

	assert.Nil(t, err, "Error should be nil on successful insert")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateCustomer_DBError(t *testing.T) {
	// Initialize mock SQL database
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	// Open GORM DB with sqlmock
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	// Initialize CustomerRepository with mock DB
	customerRepo := repository.NewCustomerRepository(gormDB)

	// Mock expected SQL queries
	mock.ExpectBegin() // GORM starts a transaction automatically

	// Mock the INSERT query with RETURNING clause
	mock.ExpectQuery(`INSERT INTO "customers" .* RETURNING "customer_id"`).
		WithArgs(
			"1234567890123456", // customer_nik
			"John Doe",         // customer_full_name
			"John D",           // customer_legal_name
			"Jakarta",          // customer_birth_place
			sqlmock.AnyArg(),   // customer_birth_date (dynamic timestamp)
			float64(0),         // customer_salary (must be float64)
			"",                 // customer_ktp_photo (empty string)
			"",                 // customer_selfie_photo (empty string)
			0,                  // customer_created_by
			sqlmock.AnyArg(),   // customer_created_at (dynamic timestamp)
			sqlmock.AnyArg(),   // customer_edited_by (can be nil)
			sqlmock.AnyArg(),   // customer_edited_at (can be nil)
		).
		WillReturnError(errors.New("invalid transaction")) // Simulate DB failure

	mock.ExpectRollback() // Rollback transaction on error

	// Create sample customer
	customer := &domain.Customer{
		CustomerNIK:         "1234567890123456",
		CustomerFullName:    "John Doe",
		CustomerLegalName:   "John D",
		CustomerBirthPlace:  "Jakarta",
		CustomerBirthDate:   time.Now(),
		CustomerSalary:      0,
		CustomerKTPPhoto:    "",
		CustomerSelfiePhoto: "",
		CustomerCreatedBy:   0,
		CustomerCreatedAt:   time.Now(),
		CustomerEditedBy:    nil,
		CustomerEditedAt:    nil,
	}

	// Call the function under test
	err = customerRepo.CreateCustomer(customer)

	// Assertions
	assert.Error(t, err, "Expected an error, but got nil")
	assert.True(t, strings.Contains(err.Error(), "invalid transaction"), "Expected error message to contain 'invalid transaction'")

	// Ensure all mock expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllCustomers_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "customers" LIMIT \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"customer_id", "customer_nik", "customer_full_name"}).
			AddRow(1, "1234567890123456", "John Doe").
			AddRow(2, "9876543210987654", "Jane Doe"))

	customers, err := customerRepo.GetAllCustomers(10, 0)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, customers, "Customers list should not be nil")
	assert.Len(t, customers, 2, "Should return 2 customers")
	assert.Equal(t, "John Doe", customers[0].CustomerFullName)
	assert.Equal(t, "Jane Doe", customers[1].CustomerFullName)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllCustomers_DBError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "customers" LIMIT \$1`).
		WithArgs(10).
		WillReturnError(gorm.ErrInvalidTransaction)

	customers, err := customerRepo.GetAllCustomers(10, 0)

	assert.Nil(t, customers, "Customers list should be nil on DB error")
	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetCustomerByID_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "customers" WHERE "customers"."customer_id" = \$1 ORDER BY "customers"."customer_id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"customer_id", "customer_nik", "customer_full_name"}).
			AddRow(1, "1234567890123456", "John Doe"))

	customer, err := customerRepo.GetCustomerByID(1)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, customer, "Customer should not be nil")
	assert.Equal(t, uint(1), customer.CustomerID)
	assert.Equal(t, "1234567890123456", customer.CustomerNIK)
	assert.Equal(t, "John Doe", customer.CustomerFullName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetCustomerByID_NotFound(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "customers" WHERE "customers"."customer_id" = \$1 ORDER BY "customers"."customer_id" LIMIT \$2`).
		WithArgs(99, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	customer, err := customerRepo.GetCustomerByID(99)

	assert.Nil(t, customer, "Customer should be nil when not found")
	assert.Error(t, err, "Error should not be nil")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Expected gorm.ErrRecordNotFound error")
}

func TestGetCustomerByID_DBError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "customers" WHERE "customers"."customer_id" = \$1 ORDER BY "customers"."customer_id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(gorm.ErrInvalidTransaction)

	customer, err := customerRepo.GetCustomerByID(1)

	assert.Nil(t, customer, "Customer should be nil on DB error")
	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")
}

func TestGetCustomerByNIK_Success(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "customers" WHERE customer_nik = \$1 ORDER BY "customers"."customer_id" LIMIT \$2`).
		WithArgs("1234567890123456", 1).
		WillReturnRows(sqlmock.NewRows([]string{"customer_id", "customer_nik", "customer_full_name"}).
			AddRow(1, "1234567890123456", "John Doe"))

	customer, err := customerRepo.GetCustomerByNIK("1234567890123456")

	assert.Nil(t, err, "Error should be nil when customer is found")
	assert.NotNil(t, customer, "Customer should not be nil")
	assert.Equal(t, "1234567890123456", customer.CustomerNIK)
	assert.Equal(t, "John Doe", customer.CustomerFullName)
}

func TestGetCustomerByNIK_NotFound(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "customers" WHERE customer_nik = \$1 ORDER BY "customers"."customer_id" LIMIT \$2`).
		WithArgs("9999999999999999", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	customer, err := customerRepo.GetCustomerByNIK("9999999999999999")

	assert.Nil(t, customer, "Customer should be nil when not found")
	assert.Error(t, err, "Error should not be nil when customer is not found")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Expected gorm.ErrRecordNotFound error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetCustomerByNIK_DBError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "customers" WHERE customer_nik = \$1 ORDER BY "customers"."customer_id" LIMIT \$2`).
		WithArgs("1234567890123456", 1).
		WillReturnError(errors.New("database connection failed"))

	customer, err := customerRepo.GetCustomerByNIK("1234567890123456")

	assert.Nil(t, customer, "Customer should be nil when DB error occurs")
	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, "database connection failed", err.Error(), "Expected database error message")
}

func TestUpdateCustomer_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "customers" SET .*`).
		WithArgs(
			"1234567890123456",
			"John Doe Updated",
			"John D Updated",
			"",
			sqlmock.AnyArg(),
			float64(0),
			"",
			"",
			0,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	customer := &domain.Customer{
		CustomerID:          1,
		CustomerNIK:         "1234567890123456",
		CustomerFullName:    "John Doe Updated",
		CustomerLegalName:   "John D Updated",
		CustomerBirthPlace:  "",
		CustomerBirthDate:   time.Time{},
		CustomerSalary:      0,
		CustomerKTPPhoto:    "",
		CustomerSelfiePhoto: "",
		CustomerCreatedBy:   0,
		CustomerCreatedAt:   time.Time{},
		CustomerEditedBy:    nil,
		CustomerEditedAt:    nil,
	}

	err = customerRepo.UpdateCustomer(customer)

	assert.Nil(t, err, "Error should be nil on successful update")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateCustomer_DBError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "customers" SET .*`).
		WithArgs(
			"1234567890123456",
			"John Doe Updated",
			"John D Updated",
			"",
			sqlmock.AnyArg(),
			float64(0),
			"",
			"",
			0,
			sqlmock.AnyArg(),
			nil,
			nil,
			1,
		).
		WillReturnError(gorm.ErrInvalidTransaction)

	mock.ExpectRollback()

	customer := &domain.Customer{
		CustomerID:          1,
		CustomerNIK:         "1234567890123456",
		CustomerFullName:    "John Doe Updated",
		CustomerLegalName:   "John D Updated",
		CustomerBirthPlace:  "",
		CustomerBirthDate:   time.Time{},
		CustomerSalary:      0,
		CustomerKTPPhoto:    "",
		CustomerSelfiePhoto: "",
		CustomerCreatedBy:   0,
		CustomerCreatedAt:   time.Time{},
		CustomerEditedBy:    nil,
		CustomerEditedAt:    nil,
	}

	err = customerRepo.UpdateCustomer(customer)

	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteCustomer_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`DELETE FROM "customers" WHERE "customers"."customer_id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = customerRepo.DeleteCustomer(1)

	assert.Nil(t, err, "Error should be nil on successful delete")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteCustomer_NotFound(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`DELETE FROM "customers" WHERE "customers"."customer_id" = \$1`).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectCommit()

	err = customerRepo.DeleteCustomer(999)

	assert.Nil(t, err, "Error should be nil when deleting a non-existing customer")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteCustomer_DBError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`DELETE FROM "customers" WHERE "customers"."customer_id" = \$1`).
		WithArgs(1).
		WillReturnError(gorm.ErrInvalidTransaction)

	mock.ExpectRollback()

	err = customerRepo.DeleteCustomer(1)

	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
