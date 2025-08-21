package repository

import (
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *model.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) FindByID(id uint) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Preload("User").Preload("History").First(&wallet, id).Error
	return &wallet, err
}

func (r *WalletRepository) FindByUserID(userID uint) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Preload("History").First(&wallet, "user_id = ?", userID).Error
	return &wallet, err
}

func (r *WalletRepository) Update(wallet *model.Wallet) error {
	return r.db.Save(wallet).Error
}
