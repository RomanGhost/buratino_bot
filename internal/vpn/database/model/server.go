package model

import (
	"gorm.io/gorm"
)

// Server represents servers table
type Server struct {
	gorm.Model
	Region string `gorm:"size:5;index"`
	Access string `gorm:"size:512,type:char"`

	// Associations
	RegionInfo Region `gorm:"foreignKey:Region;references:ShortName"`
	Keys       []Key  `gorm:"foreignKey:ServerID"`
}
