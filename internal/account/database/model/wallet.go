package model

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserID     uint        `gorm:"not null;index"` // индекс для поиска по пользователю
	MoneyCount int         `gorm:"not null"`
	History    []Operation `gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // каскад при удалении
}
