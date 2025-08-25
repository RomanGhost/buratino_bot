package service

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/repository"
	"github.com/RomanGhost/buratino_bot.git/internal/pagination"
)

type GoodsService struct {
	goodsRepository *repository.GoodsRepository
}

func NewGoodsService(goodsRepository *repository.GoodsRepository) *GoodsService {
	return &GoodsService{
		goodsRepository: goodsRepository,
	}
}

func (s *GoodsService) GetByName(goodsName string) (*model.GoodsPrice, error) {
	g, err := s.goodsRepository.FindByName(goodsName)
	if err != nil {
		return nil, fmt.Errorf("error get goods: %s", err)
	}
	if g == nil {
		return nil, fmt.Errorf("can't find goods: %s", goodsName)
	}

	return g, nil
}

func (s *GoodsService) GetAll(pageNum int, limit int) ([]model.GoodsPrice, error) {
	pages := pagination.Pagination{
		Page:  pageNum,
		Limit: limit,
	}
	return s.goodsRepository.PaginationAll(&pages)
}
