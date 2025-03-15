package domain

import "time"

type Limit struct {
	LimitID              uint       `gorm:"primaryKey" json:"limit_id"`
	LimitNIK             string     `gorm:"not null" json:"limit_nik"`
	NIKCustomer          Customer   `gorm:"foreignKey:LimitNIK;references:CustomerNIK;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LimitTenor           int        `gorm:"not null" json:"limit_tenor"`
	LimitAmount          float64    `gorm:"not null" json:"limit_amount"`
	LimitUsedAmount      float64    `gorm:"not null" json:"limit_used_amount"`
	LimitRemainingAmount float64    `gorm:"not null" json:"limit_remaining_amount"`
	LimitCreatedBy       uint       `gorm:"not null" json:"limit_created_by"`
	CreatedByUser        User       `gorm:"foreignKey:LimitCreatedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LimitCreatedAt       time.Time  `gorm:"autoCreateTime" json:"limit_created_at"`
	LimitEditedBy        *uint      `json:"limit_edited_by"`
	EditedByUser         *User      `gorm:"foreignKey:LimitEditedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LimitEditedAt        *time.Time `json:"limit_edited_at"`
}

type CreateLimitInput struct {
	LimitNIK    string  `json:"limit_nik" validate:"required,len=16,numeric"`
	LimitTenor  int     `json:"limit_tenor" validate:"required"`
	LimitAmount float64 `json:"limit_amount" validate:"required"`
}

type EditLimitInput struct {
	LimitNIK             string  `json:"limit_nik" validate:"required,len=16,numeric"`
	LimitTenor           int     `json:"limit_tenor" validate:"required"`
	LimitAmount          float64 `json:"limit_amount" validate:"required"`
	LimitUsedAmount      *float64 `json:"limit_used_amount" validate:"required"`
	LimitRemainingAmount *float64 `json:"limit_remaining_amount" validate:"required"`
}
