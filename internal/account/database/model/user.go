package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents users table
type User struct {
	gorm.Model
	TelegramID int64 `gorm:"index"`
	IsActive   bool  `gorm:"default:true"`
	BanTime    *time.Time
	Role       string `gorm:"size:16;index"`
	Timezone   time.Location

	// Associations
	UserRole UserRole `gorm:"foreignKey:Role;references:RoleName"`
	Wallet   Wallet   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
