package model

import (
	"gorm.io/gorm"
)

// User represents users table
type User struct {
	gorm.Model
	AuthID     uint
	TelegramID int64 `gorm:"index"`

	Keys []Key `gorm:"foreignKey:UserID"`
}
