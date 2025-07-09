package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents users table
type User struct {
	gorm.Model
	TelegramID int64      `gorm:"index" json:"telegram_id"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`
	BanTime    *time.Time `json:"ban_time,omitempty"`
	Role       string     `gorm:"size:16;index" json:"role"`

	// Associations
	UserRole UserRole `gorm:"foreignKey:Role;references:RoleName" json:"user_role,omitempty"`
	Keys     []Key    `gorm:"foreignKey:UserID" json:"keys,omitempty"`
}
