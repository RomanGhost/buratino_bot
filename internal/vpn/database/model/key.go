package model

import (
	"time"

	"gorm.io/gorm"
)

// Key represents keys table
type Key struct {
	gorm.Model
	OutlineKeyId int
	UserID       uint `gorm:"index"`
	ServerID     uint `gorm:"index"`
	DeadlineTime time.Time
	ConnectUrl   string
	KeyName      string
	IsActive     bool `gorm:"default:true"`
	Duration     time.Duration

	// Associations
	User   User   `gorm:"foreignKey:UserID"`
	Server Server `gorm:"foreignKey:ServerID"`
}
