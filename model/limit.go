package model

import (
	"time"
)

type Limit struct {
	LimitID        uint       `gorm:"primaryKey" json:"limit_id"`
	LimitNIK       string     `gorm:"not null" json:"limit_nik"`
	NIKCustomer    Customer   `gorm:"foreignKey:LimitNIK;references:CustomerNIK;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LimitTenor     int        `gorm:"not null" json:"limit_tenor"`
	LimitAmount    float64    `gorm:"not null" json:"limit_amount"`
	LimitCreatedBy uint       `gorm:"not null" json:"limit_created_by"`
	CreatedByUser  User       `gorm:"foreignKey:LimitCreatedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LimitCreatedAt time.Time  `gorm:"autoCreateTime" json:"limit_created_at"`
	LimitEditedBy  *uint      `json:"limit_edited_by"`
	EditedByUser   *User      `gorm:"foreignKey:LimitEditedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LimitEditedAt  *time.Time `json:"limit_edited_at"`
}
