package model

import (
	"time"
)

type Customer struct {
	CustomerID          uint       `gorm:"primaryKey" json:"customer_id"`
	CustomerNIK         string     `gorm:"unique;not null" json:"customer_nik"`
	CustomerFullName    string     `gorm:"not null" json:"customer_full_name"`
	CustomerLegalName   string     `gorm:"not null" json:"customer_legal_name"`
	CustomerBirthPlace  string     `gorm:"not null" json:"customer_birth_place"`
	CustomerBirthDate   time.Time  `gorm:"not null" json:"customer_birth_date"`
	CustomerSalary      float64    `gorm:"not null" json:"customer_salary"`
	CustomerKTPPhoto    string     `gorm:"not null" json:"customer_ktp_photo"`
	CustomerSelfiePhoto string     `gorm:"not null" json:"customer_selfie_photo"`
	CustomerCreatedBy   uint       `gorm:"not null" json:"customer_created_by"`
	CreatedByUser       User       `gorm:"foreignKey:CustomerCreatedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CustomerCreatedAt   time.Time  `gorm:"autoCreateTime" json:"customer_created_at"`
	CustomerEditedBy    *uint      `json:"customer_edited_by"`
	EditedByUser        *User      `gorm:"foreignKey:CustomerEditedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CustomerEditedAt    *time.Time `json:"customer_edited_at"`
}
