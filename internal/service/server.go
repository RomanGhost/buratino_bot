package service

import (
	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/database/repository"
)

type ServerService struct {
	serverRepository *repository.ServerRepository
}

func NewServerService(serverRepository *repository.ServerRepository) *ServerService {
	return &ServerService{serverRepository}
}

func (s *ServerService) GetServersByRegionShortName(shortRegionName string) ([]model.Server, error) {
	return s.serverRepository.GetByRegion(shortRegionName)
}

func (s *ServerService) GetServerByID(serverID uint) (*model.Server, error) {
	return s.serverRepository.GetByID(serverID)
}
