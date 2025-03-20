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

func TestCreateLimit_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO "limits" \("limit_nik","limit_tenor","limit_amount","limit_used_amount","limit_remaining_amount","limit_created_by","limit_created_at","limit_edited_by","limit_edited_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\) RETURNING "limit_id"`).
		WithArgs(
			"1234567890123456",
			12,
			float64(5000000),
			float64(0),
			float64(0),
			0,
			sqlmock.AnyArg(),
			nil,
			nil,
		).
		WillReturnRows(sqlmock.NewRows([]string{"limit_id"}).AddRow(1))

	mock.ExpectCommit()

	limit := &domain.Limit{
		LimitNIK:             "1234567890123456",
		LimitTenor:           12,
		LimitAmount:          5000000,
		LimitUsedAmount:      0,
		LimitRemainingAmount: 0,
		LimitCreatedBy:       0,
		LimitCreatedAt:       time.Now(),
		LimitEditedBy:        nil,
		LimitEditedAt:        nil,
	}

	err = limitRepo.CreateLimit(limit)

	assert.Nil(t, err, "Error should be nil on successful insert")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateLimit_DBError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "limits" \("limit_nik","limit_tenor","limit_amount","limit_used_amount","limit_remaining_amount","limit_created_by","limit_created_at","limit_edited_by","limit_edited_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\) RETURNING "limit_id"`).
		WithArgs(
			"1234567890123456",
			12,
			float64(5000000),
			float64(0),
			float64(0),
			0,
			sqlmock.AnyArg(),
			nil,
			nil,
		).
		WillReturnError(gorm.ErrInvalidTransaction)
	mock.ExpectRollback()

	limit := &domain.Limit{
		LimitNIK:             "1234567890123456",
		LimitTenor:           12,
		LimitAmount:          5000000,
		LimitUsedAmount:      0,
		LimitRemainingAmount: 0,
		LimitCreatedBy:       0,
		LimitCreatedAt:       time.Now(),
		LimitEditedBy:        nil,
		LimitEditedAt:        nil,
	}

	err = limitRepo.CreateLimit(limit)

	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllLimits_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "limits" LIMIT \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"limit_id", "limit_nik", "limit_tenor", "limit_amount"}).
			AddRow(1, "1234567890123456", 12, 5000000).
			AddRow(2, "9876543210987654", 24, 10000000))

	mock.ExpectQuery(`SELECT \* FROM "customers" WHERE "customers"."customer_nik" IN \(\$1,\$2\)`).
		WithArgs("1234567890123456", "9876543210987654").
		WillReturnRows(sqlmock.NewRows([]string{"customer_nik", "customer_name"}).
			AddRow("1234567890123456", "John Doe").
			AddRow("9876543210987654", "Jane Doe"))

	limits, err := limitRepo.GetAllLimits(10, 0)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, limits, "Limits should not be nil")
	if assert.Len(t, limits, 2, "Should return 2 limits") {
		assert.Equal(t, "1234567890123456", limits[0].LimitNIK)
		assert.Equal(t, "9876543210987654", limits[1].LimitNIK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllLimits_DBError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "limits" LIMIT \$1`).
		WithArgs(10).
		WillReturnError(gorm.ErrInvalidTransaction)

	limits, err := limitRepo.GetAllLimits(10, 0)

	assert.Nil(t, limits, "Limits should be nil on DB error")
	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetLimitByIDWithTx_Success(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT \* FROM limits WHERE limit_id = \$1 FOR UPDATE`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"limit_id", "limit_nik", "limit_tenor", "limit_amount"}).
			AddRow(1, "1234567890123456", 12, 5000000.0))

	mock.ExpectCommit()

	tx := gormDB.Begin()

	limit, err := limitRepo.GetLimitByIDWithTx(tx, 1)

	assert.Nil(t, err, "Error should be nil on successful retrieval")
	assert.NotNil(t, limit, "Limit should not be nil")
	assert.Equal(t, "1234567890123456", limit.LimitNIK, "Expected limit NIK to match")
	assert.Equal(t, 5000000.0, limit.LimitAmount, "Expected limit amount to match")

	if err == nil {
		tx.Commit()
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetLimitByIDWithTx_NotFound(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT \* FROM limits WHERE limit_id = \$1 FOR UPDATE`).
		WithArgs(99).
		WillReturnError(gorm.ErrRecordNotFound)

	tx := gormDB.Begin()

	limit, err := limitRepo.GetLimitByIDWithTx(tx, 99)

	assert.Nil(t, limit)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetLimitByIDWithTx_DBError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT \* FROM limits WHERE limit_id = \$1 FOR UPDATE`).
		WithArgs(1).
		WillReturnError(errors.New("database connection failed"))

	tx := gormDB.Begin()

	limit, err := limitRepo.GetLimitByIDWithTx(tx, 1)

	assert.Nil(t, limit)
	assert.Error(t, err)
	assert.Equal(t, "database connection failed", err.Error())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateLimitWithTx_Success(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "limits" SET`).
		WithArgs(
			"",
			0,
			float64(6000000),
			float64(1000000),
			float64(5000000),
			0,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	tx := gormDB.Begin()

	limit := &domain.Limit{
		LimitID:              1,
		LimitNIK:             "",
		LimitTenor:           0,
		LimitAmount:          6000000,
		LimitUsedAmount:      1000000,
		LimitRemainingAmount: 5000000,
		LimitCreatedBy:       0,
		LimitCreatedAt:       time.Time{},
		LimitEditedBy:        new(uint),
		LimitEditedAt:        new(time.Time),
	}

	err = limitRepo.UpdateLimitWithTx(tx, limit)

	if err == nil {
		tx.Commit()
	}

	assert.Nil(t, err, "Error should be nil on successful update")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateLimitWithTx_DBError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "limits" SET`).
		WithArgs(
			"",
			0,
			float64(6000000),
			float64(1000000),
			float64(5000000),
			0,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			1,
		).
		WillReturnError(gorm.ErrInvalidTransaction)

	mock.ExpectRollback()

	tx := gormDB.Begin()

	limit := &domain.Limit{
		LimitID:              1,
		LimitNIK:             "",
		LimitTenor:           0,
		LimitAmount:          6000000,
		LimitUsedAmount:      1000000,
		LimitRemainingAmount: 5000000,
		LimitCreatedBy:       0,
		LimitCreatedAt:       time.Time{},
		LimitEditedBy:        new(uint),
		LimitEditedAt:        new(time.Time),
	}

	err = limitRepo.UpdateLimitWithTx(tx, limit)

	if err != nil {
		tx.Rollback()
	}

	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteLimit_Success(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "limits" WHERE "limits"."limit_id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = limitRepo.DeleteLimit(1)

	assert.Nil(t, err, "Error should be nil on successful delete")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteLimit_NotFound(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "limits" WHERE "limits"."limit_id" = \$1`).
		WithArgs(99).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = limitRepo.DeleteLimit(99)

	assert.Nil(t, err, "Error should be nil even if the limit was not found")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteLimit_DBError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	limitRepo := repository.NewLimitRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "limits" WHERE "limits"."limit_id" = \$1`).
		WithArgs(1).
		WillReturnError(gorm.ErrInvalidTransaction)
	mock.ExpectRollback()

	err = limitRepo.DeleteLimit(1)

	assert.Error(t, err, "Error should not be nil on DB error")
	assert.Equal(t, gorm.ErrInvalidTransaction, err, "Expected gorm.ErrInvalidTransaction error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
