package model

import (
	"gorm.io/gorm"
)

// Server represents servers table
type Server struct {
	gorm.Model
	Region string `gorm:"size:128;index" json:"region"`
	IPv4   string `gorm:"type:cidr" json:"ipv4"`
	IPv6   string `gorm:"type:cidr" json:"ipv6"`

	// Associations
	RegionInfo Region `gorm:"foreignKey:Region;references:RegionName" json:"region_info,omitempty"`
	Keys       []Key  `gorm:"foreignKey:ServerID" json:"keys,omitempty"`
}
