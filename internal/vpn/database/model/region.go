package model

// Region represents regions table
type Region struct {
	RegionName string   `gorm:"size:128" json:"region_name"`
	ShortName  string   `gorm:"primaryKey;size:4" json:"short_name"`
	Servers    []Server `gorm:"foreignKey:Region;references:ShortName" json:"servers,omitempty"`
}
