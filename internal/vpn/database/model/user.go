package model

import (
	"gorm.io/gorm"
)

// User represents users table
type User struct {
	gorm.Model
	TelegramID int64 `gorm:"index" json:"telegram_id"`

	Keys []Key `gorm:"foreignKey:UserID" json:"keys,omitempty"`
}
