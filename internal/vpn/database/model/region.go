package model

// Region represents regions table
type Region struct {
	RegionName string   `gorm:"size:128"`
	ShortName  string   `gorm:"primaryKey;size:4"`
	Servers    []Server `gorm:"foreignKey:Region;references:ShortName"`
}
