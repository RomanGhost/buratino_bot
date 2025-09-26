package model

import (
	"gorm.io/gorm"
)

// Server represents servers table
type Server struct {
	gorm.Model
	Region     string `gorm:"size:5;index"`
	Access     string `gorm:"size:512,type:char"`
	ProviderID string `gorm:"size:16;index"`

	// Associations
	RegionInfo Region   `gorm:"foreignKey:Region;references:ShortName"`
	Keys       []Key    `gorm:"foreignKey:ServerID"`
	Provider   Provider `gorm:"foreignKey:ProviderID"`
}
