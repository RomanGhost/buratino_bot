package model

import (
	"time"

	"gorm.io/gorm"
)

// Key represents keys table
type Key struct {
	gorm.Model
	KeyID int
	UserID       uint `gorm:"index"`
	ServerID     uint `gorm:"index"`
	DeadlineTime time.Time
	ConnectUrl   string
	KeyName      string
	IsActive     bool `gorm:"default:true"`
	Duration     time.Duration
	SupplierID string `gorm:"size:16;index"`

	// Associations
	User   User   `gorm:"foreignKey:UserID"`
	Server Server `gorm:"foreignKey:ServerID"`	
	Supplier Supplier `gorm:"foreignKey:SupplierID"`
}
