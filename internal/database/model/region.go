package model

// Region represents regions table
type Region struct {
	RegionName string   `gorm:"size:128" json:"region_name"`
	ShortName  string   `gorm:"primaryKey;size:5" json:"short_name"`
	Servers    []Server `gorm:"foreignKey:Region;references:RegionName" json:"servers,omitempty"`
}
