package model

import (
	"time"

	"gorm.io/gorm"
)

// Key represents keys table
type Key struct {
	gorm.Model
	OutlineKeyId int       `json:"outline_key_id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	ServerID     uint      `gorm:"index" json:"server_id"`
	DeadlineTime time.Time `json:"deadline_time"`
	ConnectUrl   string    `json:"connect_url"`
	KeyName      string    `json:"key_name"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`

	// Associations
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Server Server `gorm:"foreignKey:ServerID" json:"server,omitempty"`
}
