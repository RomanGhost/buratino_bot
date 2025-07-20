package service

import (
	"fmt"
	"log"
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

func (s *KeyService) CreateKeyWithDeadline(outlineKeyId int, userTelegramID int64, serverID uint, connectURL string, keyName string, deadline time.Duration) (*model.Key, error) {
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
		OutlineKeyId: outlineKeyId,
		ServerID:     server.ID,
		UserID:       user.ID,
		DeadlineTime: time.Now().UTC().Truncate(time.Minute).Add(deadline),
		ConnectUrl:   connectURL,
		KeyName:      keyName,
	}

	err = s.keyRepository.Create(&newKey)
	if err != nil {
		return nil, fmt.Errorf("error create new key: %v", err)
	}

	return &newKey, nil
}

func (s *KeyService) CreateDefaultKey(outlineKeyId int, userTelegramID int64, serverID uint, connectURL string, keyName string) (*model.Key, error) {
	return s.CreateKeyWithDeadline(outlineKeyId, userTelegramID, serverID, connectURL, keyName, time.Duration(30*time.Minute))
}

func (s *KeyService) CountKeysOfServer(serverID uint) int {
	key, err := s.keyRepository.GetByServerID(serverID)
	if err != nil {
		return -1
	}
	return len(key)
}

func (s *KeyService) Delete(keyID uint) {
	s.keyRepository.Delete(keyID)
}

func (s *KeyService) DeactivateKey(keyID uint) error {
	return s.keyRepository.DeactivateKey(keyID)
}

func (s *KeyService) GetExpiringSoon(timeDuration time.Duration) ([]model.Key, error) {
	timeStart := time.Now().UTC().Truncate(time.Minute)
	timeEnd := timeStart.Add(timeDuration)

	return s.keyRepository.GetExpiringSoon(timeStart, timeEnd)
}

func (s *KeyService) GetExpiredKeys() ([]model.Key, error) {
	timeDeadline := time.Now().UTC().Truncate(time.Minute)

	return s.keyRepository.GetExpiredActiveKeys(timeDeadline)
}

func (s *KeyService) IsActiveKey(keyID uint) bool {
	key, err := s.keyRepository.GetByID(keyID)
	if err != nil {
		log.Println("[INFO] Can't get key")
		return false
	}
	if key == nil {
		return false
	}
	if !key.IsActive {
		return false
	}
	return true
}

func (s *KeyService) ExtendKeyByID(keyID uint) (*model.Key, error) {
	key, err := s.keyRepository.GetByID(keyID)
	if err != nil {
		return nil, fmt.Errorf("error get key by ID: %v", keyID)
	}
	newKeyDeadlineTime := key.DeadlineTime.Add(key.Duration)
	key.DeadlineTime = newKeyDeadlineTime

	err = s.keyRepository.ExtendKey(keyID, newKeyDeadlineTime)
	if err != nil {
		return nil, fmt.Errorf("error expire key: %v", err)
	}
	return key, nil
}
