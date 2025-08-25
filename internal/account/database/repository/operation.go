package repository

import (
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"gorm.io/gorm"
)

type OperationRepository struct {
	db *gorm.DB
}

func NewOperationRepository(db *gorm.DB) *OperationRepository {
	return &OperationRepository{db: db}
}

func (r *OperationRepository) Create(op *model.Operation) error {
	return r.db.Create(op).Error
}

func (r *OperationRepository) FindByID(id uint) (*model.Operation, error) {
	var op model.Operation
	err := r.db.Preload("Goods").Preload("Wallet").First(&op, id).Error
	return &op, err
}

func (r *OperationRepository) FindByWalletID(walletID uint) ([]model.Operation, error) {
	var ops []model.Operation
	err := r.db.Preload("Wallet").Find(&ops, "wallet_id = ?", walletID).Error
	return ops, err
}

func (r *OperationRepository) FindByGoodsID(goodID uint) ([]model.Operation, error) {
	var ops []model.Operation
	err := r.db.Preload("Goods").Find(&ops, "goods_id = ?", goodID).Error
	return ops, err
}
