package repository

import (
	"kreditplus/internal/domain"

	"gorm.io/gorm"
)

type LimitRepository interface {
	CreateLimit(limit *domain.Limit) error
	GetAllLimits(limit, offset int) ([]domain.Limit, error)
	GetLimitByID(id uint) (*domain.Limit, error)
	GetLimitByIDWithTx(tx *gorm.DB, id uint) (*domain.Limit, error)
	GetLimitByNIKandTenorWithTx(tx *gorm.DB, nik string, tenor float64) (*domain.Limit, error)
	UpdateLimitWithTx(tx *gorm.DB, limit *domain.Limit) error
	UpdateLimit(limit *domain.Limit) error
	DeleteLimit(id uint) error
}

type limitRepository struct {
	db *gorm.DB
}

func NewLimitRepository(db *gorm.DB) LimitRepository {
	return &limitRepository{db: db}
}

func (r *limitRepository) CreateLimit(limit *domain.Limit) error {
	return r.db.Create(limit).Error
}

func (r *limitRepository) GetAllLimits(limit, offset int) ([]domain.Limit, error) {
	var limits []domain.Limit
	err := r.db.Preload("NIKCustomer").
		Preload("CreatedByUser").
		Preload("EditedByUser").
		Limit(limit).
		Offset(offset).
		Find(&limits).Error
	if err != nil {
		return nil, err
	}
	return limits, nil
}

func (r *limitRepository) GetLimitByID(id uint) (*domain.Limit, error) {
	var limit domain.Limit
	err := r.db.Preload("NIKCustomer").
		Preload("CreatedByUser").
		Preload("EditedByUser").
		First(&limit, id).Error
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (r *limitRepository) GetLimitByIDWithTx(tx *gorm.DB, id uint) (*domain.Limit, error) {
	var limit domain.Limit
	err := tx.Raw(`SELECT * FROM limits WHERE limit_id = ? FOR UPDATE`, id).Scan(&limit).Error
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (r *limitRepository) GetLimitByNIKandTenorWithTx(tx *gorm.DB, nik string, tenor float64) (*domain.Limit, error) {
	var limit domain.Limit
	err := tx.Raw(`SELECT * FROM limits WHERE limit_nik = ? AND limit_tenor = ? FOR UPDATE`, nik, tenor).Scan(&limit).Error
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (r *limitRepository) UpdateLimitWithTx(tx *gorm.DB, limit *domain.Limit) error {
	return tx.Save(limit).Error
}

func (r *limitRepository) UpdateLimit(limit *domain.Limit) error {
	return r.db.Save(limit).Error
}

func (r *limitRepository) DeleteLimit(id uint) error {
	return r.db.Delete(&domain.Limit{}, id).Error
}
