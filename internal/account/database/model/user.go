package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramID       int64 `gorm:"index"`
	TelegramUsername string
	IsActive         bool `gorm:"default:true"`
	BanTime          *time.Time
	Role             string `gorm:"size:16;index"`
	TimezoneOffset   int    // смещение от UTC в минутах

	// Associations
	UserRole UserRole `gorm:"foreignKey:Role;references:RoleName"`
	Wallet   Wallet   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
