package model

import "gorm.io/gorm"

type Operation struct {
	gorm.Model
	WalletID uint       `gorm:"not null;index"`
	Wallet   Wallet     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GoodsID  uint       `gorm:"not null;index"` // связь с товаром
	Goods    GoodsPrice `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Count    uint64
}
