package service

import (
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/repository"
)

type ProviderService struct {
	providerRepository *repository.ProviderRepository
}

func NewProviderService(providerRepository *repository.ProviderRepository) *ProviderService {
	return &ProviderService{providerRepository}
}

func (s *ProviderService) GetProviders() ([]model.Provider, error) {
	return s.providerRepository.GetAll()
}
