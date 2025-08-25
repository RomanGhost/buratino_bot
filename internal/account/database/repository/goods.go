package repository

import (
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/pagination"
	"gorm.io/gorm"
)

type GoodsRepository struct {
	db *gorm.DB
}

func NewGoodsRepository(db *gorm.DB) *GoodsRepository {
	return &GoodsRepository{db: db}
}

func (r *GoodsRepository) Create(goods *model.GoodsPrice) error {
	return r.db.Create(goods).Error
}

func (r *GoodsRepository) FindByID(id uint) (*model.GoodsPrice, error) {
	var goods model.GoodsPrice
	err := r.db.First(&goods, id).Error
	return &goods, err
}

func (r *GoodsRepository) FindByName(name string) (*model.GoodsPrice, error) {
	var goods model.GoodsPrice
	err := r.db.First(&goods, "name = ?", name).Error
	return &goods, err
}

func (r *GoodsRepository) Update(goods *model.GoodsPrice) error {
	return r.db.Save(goods).Error
}

func (r *GoodsRepository) All() ([]model.GoodsPrice, error) {
	var goods []model.GoodsPrice
	err := r.db.Find(&goods).Error
	return goods, err
}

func (r *GoodsRepository) PaginationAll(p *pagination.Pagination) ([]model.GoodsPrice, error) {
	var goods []model.GoodsPrice

	// выбираем данные с лимитом и оффсетом
	err := r.db.
		Limit(p.Limit).
		Offset(p.GetOffset()).
		Find(&goods).Error

	if err != nil {
		return nil, err
	}

	return goods, nil
}
