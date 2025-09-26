package model

type Supplier struct {
	Name string   `gorm:"size:16;primaryKey"`
}