package model

import "gorm.io/gorm"

type GoodsPrice struct {
	gorm.Model
	SysName string `gorm:"type:varchar(128);unique;not null"`
	Name    string `gorm:"type:varchar(128);unique;not null"`
	Price   int64  `gorm:"not null"`
}
