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

func (u *transactionUsecase) CreateTransactionWithLimitUpdate(userID uint, customer *domain.Customer, input domain.TransactionInput) error {
	input.TransactionNIK = utils.SanitizeString(input.TransactionNIK)
	input.TransactionAssetName = utils.SanitizeString(input.TransactionAssetName)

	input.TransactionOTR = utils.SanitizeNumberFloat64(input.TransactionOTR)
	input.TransactionAdminFee = utils.SanitizeNumberFloat64(input.TransactionAdminFee)
	input.TransactionInstallment = utils.SanitizeNumberFloat64(input.TransactionInstallment)
	input.TransactionInterest = utils.SanitizeNumberFloat64(input.TransactionInterest)

	const maxRetries = 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := u.transactionRepo.WithTransaction(func(tx *gorm.DB) error {
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

		if err != nil {
			if err.Error() == "database is locked" || err.Error() == "deadlock detected" {
				utils.Logger.Warnf("Deadlock detected. Retrying transaction %d/%d", attempt+1, maxRetries)
				time.Sleep(time.Millisecond * 100)
				continue
			}
			return err
		}
		return nil
	}
	return errors.New("failed to process transaction after multiple retries")
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
	input.TransactionNIK = utils.SanitizeString(input.TransactionNIK)
	input.TransactionAssetName = utils.SanitizeString(input.TransactionAssetName)

	input.TransactionOTR = utils.SanitizeNumberFloat64(input.TransactionOTR)
	input.TransactionAdminFee = utils.SanitizeNumberFloat64(input.TransactionAdminFee)
	input.TransactionInstallment = utils.SanitizeNumberFloat64(input.TransactionInstallment)
	input.TransactionInterest = utils.SanitizeNumberFloat64(input.TransactionInterest)

	const maxRetries = 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := u.transactionRepo.WithTransaction(func(tx *gorm.DB) error {
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

			if input.TransactionOTR > 0 {
				transaction.TransactionOTR = input.TransactionOTR
			}
			if input.TransactionAdminFee > 0 {
				transaction.TransactionAdminFee = input.TransactionAdminFee
			}
			if input.TransactionInstallment > 0 {
				transaction.TransactionInstallment = input.TransactionInstallment
			}
			if input.TransactionInterest > 0 {
				transaction.TransactionInterest = input.TransactionInterest
			}
			if input.TransactionAssetName != "" {
				transaction.TransactionAssetName = input.TransactionAssetName
			}

			transaction.TransactionEditedBy = &userID
			transaction.TransactionEditedAt = &timeNow

			if err := u.transactionRepo.UpdateTransactionWithTx(tx, transaction); err != nil {
				utils.Logger.WithFields(logrus.Fields{
					"transaction_nik": transaction.TransactionNIK,
					"error":           err.Error(),
				}).Error("Failed to update transaction")
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

		if err != nil && (err.Error() == "database is locked" || err.Error() == "deadlock detected") {
			utils.Logger.Warnf("Deadlock detected. Retrying update transaction %d/%d", attempt+1, maxRetries)
			time.Sleep(time.Millisecond * 100)
			continue
		}
		return err
	}
	return errors.New("failed to update transaction after multiple retries")
}

func (u *transactionUsecase) DeleteTransactionWithLimitUpdate(userID uint, transaction *domain.Transaction) error {
	const maxRetries = 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := u.transactionRepo.WithTransaction(func(tx *gorm.DB) error {
			limit, err := u.limitRepo.GetLimitByIDWithTx(tx, transaction.TransactionLimit)
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
				}).Error("Failed to delete transaction")
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
			}).Info("Transaction successfully deleted and limit updated")

			return nil
		})

		if err != nil && (err.Error() == "database is locked" || err.Error() == "deadlock detected") {
			utils.Logger.Warnf("Deadlock detected. Retrying delete transaction %d/%d", attempt+1, maxRetries)
			time.Sleep(time.Millisecond * 100)
			continue
		}
		return err
	}
	return errors.New("failed to delete transaction after multiple retries")
}
