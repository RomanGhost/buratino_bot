package service

import (
	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/database/repository"
)

type RegionService struct {
	regionRepository *repository.RegionRepository
}

func NewRegionService(regionRepository *repository.RegionRepository) *RegionService {
	return &RegionService{regionRepository}
}

func (s *RegionService) GetRegionsWithServers() ([]model.Region, error) {
	return s.regionRepository.GetAllWithServers()
}
