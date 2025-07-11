package model

import (
	"gorm.io/gorm"
)

// Server represents servers table
type Server struct {
	gorm.Model
	Region string `gorm:"size:5;index" json:"region"`
	Access string `gorm:"size:512,type:char" json:"acess"`

	// Associations
	RegionInfo Region `gorm:"foreignKey:Region;references:ShortName" json:"region_info,omitempty"`
	Keys       []Key  `gorm:"foreignKey:ServerID" json:"keys,omitempty"`
}
