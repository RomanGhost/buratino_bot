package model

import "gorm.io/gorm"

type GoodsPrice struct {
	gorm.Model
	Name   string      `gorm:"type:varchar(128);unique;not null"`
	Price  int         `gorm:"not null"`
	Mount  int         `gorm:"not null"`
	Orders []Operation `gorm:"foreignKey:GoodsID"` // связь: какой товар покупался в операциях
}
