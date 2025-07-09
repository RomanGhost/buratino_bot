package model

// Region represents regions table
type Region struct {
	RegionName string   `gorm:"primaryKey;size:128" json:"region_name"`
	ShortName  string   `gorm:"size:5" json:"short_name"`
	Servers    []Server `gorm:"foreignKey:Region;references:RegionName" json:"servers,omitempty"`
}
