package repository

import (
	"kreditplus/internal/domain"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	WithTransaction(fn func(tx *gorm.DB) error) error
	CreateTransactionWithTx(tx *gorm.DB, transaction *domain.Transaction) error
	CreateTransaction(transaction *domain.Transaction) error
	GetAllTransactions(transaction, offset int) ([]domain.Transaction, error)
	GetTransactionByID(id uint) (*domain.Transaction, error)
	UpdateTransactionWithTx(tx *gorm.DB, transaction *domain.Transaction) error
	DeleteTransactionWithTx(tx *gorm.DB, transaction *domain.Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *transactionRepository) CreateTransactionWithTx(tx *gorm.DB, transaction *domain.Transaction) error {
	return tx.Create(transaction).Error
}

func (r *transactionRepository) CreateTransaction(transaction *domain.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) GetAllTransactions(limit, offset int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.Preload("NIKCustomer.CreatedByUser").
		Preload("NIKCustomer.EditedByUser").
		Preload("IDLimit.NIKCustomer").
		Preload("IDLimit.CreatedByUser").
		Preload("IDLimit.EditedByUser").
		Preload("CreatedByUser").
		Preload("EditedByUser").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetTransactionByID(id uint) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.Preload("NIKCustomer.CreatedByUser").
		Preload("NIKCustomer.EditedByUser").
		Preload("IDLimit.NIKCustomer").
		Preload("IDLimit.CreatedByUser").
		Preload("IDLimit.EditedByUser").
		Preload("CreatedByUser").
		Preload("EditedByUser").
		First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) UpdateTransactionWithTx(tx *gorm.DB, transaction *domain.Transaction) error {
	return tx.Save(transaction).Error
}

func (r *transactionRepository) DeleteTransactionWithTx(tx *gorm.DB, transaction *domain.Transaction) error {
	return tx.Delete(&transaction).Error
}
