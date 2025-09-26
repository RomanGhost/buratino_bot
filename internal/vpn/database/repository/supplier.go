package repository

import (
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"gorm.io/gorm"
)

type ProviderRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{
		db: db,
	}
}

func (r *ProviderRepository) Create(Provider *model.Provider) error {
	return r.db.Create(Provider).Error
}

func (r *ProviderRepository) GetByName(id uint) (*model.Provider, error) {
	var Provider model.Provider
	err := r.db.Where("name=", id, true).First(&Provider).Error
	if err != nil {
		return nil, err
	}
	return &Provider, nil
}

func (r *ProviderRepository) GetAll() ([]model.Provider, error) {
	var providers []model.Provider
	err := r.db.Find(providers).Error
	if err != nil {
		return nil, err
	}
	return providers, nil
}
