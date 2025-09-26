package service

import (
	"math"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/repository"
)

type ServerService struct {
	serverRepository *repository.ServerRepository
	keyService       *KeyService
}

func NewServerService(serverRepository *repository.ServerRepository, keyService *KeyService) *ServerService {
	return &ServerService{
		serverRepository: serverRepository,
		keyService:       keyService,
	}
}

func (s *ServerService) GetNotLoadedByRegionAndProviderServer(shortRegionName, providerName string) (*model.Server, error) {
	servers, err := s.GetServersByRegionShortName(shortRegionName)
	if err != nil || len(servers) == 0 {
		return nil, err
	}

	// chose server with min keys of region
	minCount := math.MaxInt
	var minServer model.Server
	for _, server := range servers {
		val := s.keyService.CountKeysOfServer(server.ID)
		if val == -1 {
			continue
		}
		if minCount > val {
			minCount = val
			minServer = server
		}
	}

	if minServer.ID == 0 {
		return nil, err
	}
	return &minServer, nil
}

func (s *ServerService) GetServersByRegionShortName(shortRegionName string) ([]model.Server, error) {
	return s.serverRepository.GetByRegion(shortRegionName)
}

func (s *ServerService) GetServerByID(serverID uint) (*model.Server, error) {
	return s.serverRepository.GetByID(serverID)
}
