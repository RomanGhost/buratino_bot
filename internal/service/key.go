package service

import (
	"fmt"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/database/repository"
)

type KeyService struct {
	keyRepository    *repository.KeyRepository
	userRepository   *repository.UserRepository
	serverRepository *repository.ServerRepository
}

func NewKeyService(keyRepository *repository.KeyRepository, userRepository *repository.UserRepository, serverRepository *repository.ServerRepository) *KeyService {
	return &KeyService{
		keyRepository:    keyRepository,
		userRepository:   userRepository,
		serverRepository: serverRepository,
	}
}

func (s *KeyService) CreateKey(userTelegramID int64, serverID uint, connectURL string) (*model.Key, error) {
	user, err := s.userRepository.GetByTelegramID(userTelegramID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user does not exis: %v", userTelegramID)
	}

	server, err := s.serverRepository.GetByID(serverID)
	if err != nil {
		return nil, err
	}
	if server == nil {
		return nil, fmt.Errorf("server does not exis: %v", serverID)
	}

	newKey := model.Key{
		ServerID:     server.ID,
		UserID:       user.ID,
		DeadlineTime: time.Now().Add(30 * time.Minute),
		ConnectUrl:   connectURL,
	}

	err = s.keyRepository.Create(&newKey)
	if err != nil {
		return nil, fmt.Errorf("error create new key: %v", err)
	}

	return &newKey, nil
}
