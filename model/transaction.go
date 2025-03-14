package model

import "time"

type Transaction struct {
	TransactionID             uint       `gorm:"primaryKey" json:"transaction_id"`
	TransactionContractNumber string     `gorm:"unique;not null" json:"transaction_contract_number"`
	TransactionNIK            string     `gorm:"not null" json:"transaction_nik"`
	NIKCustomer               Customer   `gorm:"foreignKey:TransactionNIK;references:CustomerNIK;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionLimit          uint       `gorm:"not null" json:"transaction_limit"`
	IDLimit                   Limit      `gorm:"foreignKey:TransactionLimit;references:LimitID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionOTR            float64    `gorm:"not null" json:"transaction_otr"`
	TransactionAdminFee       float64    `gorm:"not null" json:"transaction_admin_fee"`
	TransactionInstallment    float64    `gorm:"not null" json:"transaction_installment"`
	TransactionInterest       float64    `gorm:"not null" json:"transaction_interest"`
	TransactionAssetName      string     `gorm:"not null" json:"transaction_asset_name"`
	TransactionDate           time.Time  `gorm:"not null" json:"transaction_date"`
	TransactionCreatedBy      uint       `gorm:"not null" json:"transaction_created_by"`
	CreatedByUser             User       `gorm:"foreignKey:TransactionCreatedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionCreatedAt      time.Time  `gorm:"autoCreateTime" json:"transaction_created_at"`
	TransactionEditedBy       *uint      `json:"transaction_edited_by"`
	EditedByUser              *User      `gorm:"foreignKey:TransactionEditedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionEditedAt       *time.Time `json:"transaction_edited_at"`
}
