package usecase

import (
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/repository"
	"kreditplus/internal/utils"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionUsecase interface {
	CreateTransaction(input domain.Transaction) error
	CreateTransactionWithLimitUpdate(userID uint, customer *domain.Customer, input domain.TransactionInput) error
	GetAllTransactions(limit, offset int) ([]domain.Transaction, error)
	GetTransactionByID(id uint) (*domain.Transaction, error)
	GetCustomerByNIK(nik string) (*domain.Customer, error)
	UpdateTransactionWithLimitUpdate(userID uint, customer *domain.Customer, transaction *domain.Transaction, input domain.TransactionInput) error
	DeleteTransactionWithLimitUpdate(userID uint, transaction *domain.Transaction) error
}

type transactionUsecase struct {
	customerRepo    repository.CustomerRepository
	limitRepo       repository.LimitRepository
	transactionRepo repository.TransactionRepository
}

func NewTransactionUsecase(customerRepo repository.CustomerRepository, limitRepo repository.LimitRepository, transactionRepo repository.TransactionRepository) TransactionUsecase {
	return &transactionUsecase{customerRepo: customerRepo, limitRepo: limitRepo, transactionRepo: transactionRepo}
}

func (u *transactionUsecase) CreateTransaction(input domain.Transaction) error {
	if input.TransactionNIK == "" || len(input.TransactionNIK) != 16 {
		return errors.New("invalid NIK")
	}
	input.TransactionCreatedAt = time.Now()

	return u.transactionRepo.CreateTransaction(&input)
}

func (u *transactionUsecase) CreateTransactionWithLimitUpdate(userID uint, customer *domain.Customer, input domain.TransactionInput) error {
	return u.transactionRepo.WithTransaction(func(tx *gorm.DB) error {
		limit, err := u.limitRepo.GetLimitByNIKandTenorWithTx(tx, input.TransactionNIK, input.TransactionInstallment)
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"transaction_nik": input.TransactionNIK,
				"error":           err.Error(),
			}).Warn("Failed to retrieve limit")
			return err

		}

		totalAmount := input.TransactionOTR + input.TransactionAdminFee +
			(input.TransactionInterest * input.TransactionOTR * float64(input.TransactionInstallment) / 100)

		if totalAmount > limit.LimitRemainingAmount {
			utils.Logger.Warnf("Insufficient limit for NIK %s", input.TransactionNIK)
			return errors.New("insufficient limit")
		}

		transaction := domain.Transaction{
			TransactionContractNumber: utils.GenerateContractNumber(),
			TransactionNIK:            customer.CustomerNIK,
			TransactionLimit:          limit.LimitID,
			TransactionOTR:            input.TransactionOTR,
			TransactionAdminFee:       input.TransactionAdminFee,
			TransactionInstallment:    input.TransactionInstallment,
			TransactionInterest:       input.TransactionInterest,
			TransactionAssetName:      input.TransactionAssetName,
			TransactionDate:           time.Now(),
			TransactionCreatedBy:      userID,
			TransactionCreatedAt:      time.Now(),
		}

		if transaction.TransactionNIK == "" || len(transaction.TransactionNIK) != 16 {
			utils.Logger.Warn("Validation failed: Invalid NIK")
			return errors.New("invalid NIK")
		}

		if err := u.transactionRepo.CreateTransactionWithTx(tx, &transaction); err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"transaction_nik": transaction.TransactionNIK,
				"error":           err.Error(),
			}).Error("Failed to create transaction")
			return err
		}

		limit.LimitUsedAmount += totalAmount
		limit.LimitRemainingAmount -= totalAmount

		if err := u.limitRepo.UpdateLimitWithTx(tx, limit); err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"transaction_nik": transaction.TransactionNIK,
				"error":           err.Error(),
			}).Error("Failed to update limit")
			return err
		}

		utils.Logger.WithFields(logrus.Fields{
			"user_id":                     userID,
			"transaction_contract_number": transaction.TransactionContractNumber,
			"limit_remaining_amount":      limit.LimitRemainingAmount,
			"updated_at":                  time.Now(),
		}).Info("Transaction successfully created and limit updated")

		return nil
	})
}

func (u *transactionUsecase) GetAllTransactions(limit, offset int) ([]domain.Transaction, error) {
	return u.transactionRepo.GetAllTransactions(limit, offset)
}

func (u *transactionUsecase) GetTransactionByID(id uint) (*domain.Transaction, error) {
	return u.transactionRepo.GetTransactionByID(id)
}

func (u *transactionUsecase) GetCustomerByNIK(nik string) (*domain.Customer, error) {
	return u.customerRepo.GetCustomerByNIK(nik)
}

func (u *transactionUsecase) UpdateTransactionWithLimitUpdate(userID uint, customer *domain.Customer, transaction *domain.Transaction, input domain.TransactionInput) error {
	return u.transactionRepo.WithTransaction(func(tx *gorm.DB) error {
		timeNow := time.Now()

		limit, err := u.limitRepo.GetLimitByNIKandTenorWithTx(tx, input.TransactionNIK, input.TransactionInstallment)
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"transaction_nik": input.TransactionNIK,
				"error":           err.Error(),
			}).Warn("Failed to retrieve limit")
			return err
		}

		originalAmount := transaction.TransactionOTR + transaction.TransactionAdminFee +
			(transaction.TransactionInterest * transaction.TransactionOTR * float64(transaction.TransactionInstallment) / 100)

		newAmount := input.TransactionOTR + input.TransactionAdminFee +
			(input.TransactionInterest * input.TransactionOTR * float64(input.TransactionInstallment) / 100)

		limit.LimitUsedAmount -= originalAmount
		limit.LimitRemainingAmount += originalAmount

		if newAmount > limit.LimitRemainingAmount {
			utils.Logger.Warnf("Insufficient limit for NIK %s", input.TransactionNIK)
			return errors.New("insufficient limit")
		}

		if transaction.TransactionOTR > 0 {
			transaction.TransactionOTR = input.TransactionOTR
		}

		if transaction.TransactionAdminFee > 0 {
			transaction.TransactionAdminFee = input.TransactionAdminFee
		}

		if transaction.TransactionInstallment > 0 {
			transaction.TransactionInstallment = input.TransactionInstallment
		}

		if transaction.TransactionInterest > 0 {
			transaction.TransactionInterest = input.TransactionInterest
		}

		if transaction.TransactionAssetName != "" {
			transaction.TransactionAssetName = input.TransactionAssetName
		}

		transaction.TransactionEditedBy = &userID
		transaction.TransactionEditedAt = &timeNow

		if err := u.transactionRepo.UpdateTransactionWithTx(tx, transaction); err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"transaction_nik": transaction.TransactionNIK,
				"error":           err.Error(),
			}).Error("Failed to create transaction")
			return err
		}

		limit.LimitUsedAmount += newAmount
		limit.LimitRemainingAmount -= newAmount

		if err := u.limitRepo.UpdateLimitWithTx(tx, limit); err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"transaction_nik": transaction.TransactionNIK,
				"error":           err.Error(),
			}).Error("Failed to update limit")
			return err
		}

		utils.Logger.WithFields(logrus.Fields{
			"user_id":                     userID,
			"transaction_contract_number": transaction.TransactionContractNumber,
			"limit_remaining_amount":      limit.LimitRemainingAmount,
			"updated_at":                  time.Now(),
		}).Info("Transaction successfully updated and limit updated")

		return nil
	})
}

func (u *transactionUsecase) DeleteTransactionWithLimitUpdate(userID uint, transaction *domain.Transaction) error {
	return u.transactionRepo.WithTransaction(func(tx *gorm.DB) error {
		limit, err := u.limitRepo.GetLimitByID(transaction.TransactionLimit)
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"transaction_limit_id": transaction.TransactionLimit,
				"error":                err.Error(),
			}).Warn("Failed to retrieve associated limit record")
			return err
		}

		originalAmount := transaction.TransactionOTR + transaction.TransactionAdminFee +
			(transaction.TransactionInterest * transaction.TransactionOTR * float64(transaction.TransactionInstallment) / 100)

		limit.LimitUsedAmount -= originalAmount
		limit.LimitRemainingAmount += originalAmount

		if err := u.transactionRepo.DeleteTransactionWithTx(tx, transaction); err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"transaction_contract_number": transaction.TransactionContractNumber,
				"error":                       err.Error(),
			}).Error("Failed to create transaction")
			return err
		}

		if err := u.limitRepo.UpdateLimitWithTx(tx, limit); err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"transaction_contract_number": transaction.TransactionContractNumber,
				"error":                       err.Error(),
			}).Error("Failed to update limit")
			return err
		}

		utils.Logger.WithFields(logrus.Fields{
			"user_id":                     userID,
			"transaction_contract_number": transaction.TransactionContractNumber,
			"limit_remaining_amount":      limit.LimitRemainingAmount,
			"updated_at":                  time.Now(),
		}).Info("Transaction successfully updated and limit updated")

		return nil
	})
}
