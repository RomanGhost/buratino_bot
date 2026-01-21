package service

import (
	"fmt"
	"log"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/repository"
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

func (s *KeyService) CreateKeyWithDeadline(KeyID int, telegramUserID int64, serverID uint, connectURL string, keyName string, duration time.Duration) (*model.Key, error) {
	user, err := s.userRepository.GetByTelegramID(telegramUserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user does not exis: %v", telegramUserID)
	}

	server, err := s.serverRepository.GetByID(serverID)
	if err != nil {
		return nil, err
	}
	if server == nil {
		return nil, fmt.Errorf("server does not exis: %v", serverID)
	}

	newKey := model.Key{
		KeyID:        KeyID,
		ServerID:     server.ID,
		UserID:       user.ID,
		DeadlineTime: time.Now().UTC().Truncate(time.Minute).Add(duration),
		ConnectUrl:   connectURL,
		KeyName:      keyName,
		Duration:     duration,
	}

	err = s.keyRepository.Create(&newKey)
	if err != nil {
		return nil, fmt.Errorf("error create new key: %v", err)
	}

	return &newKey, nil
}

func (s *KeyService) GetKeysByTelegramUserID(telegramUserID int64) ([]model.Key, error) {
	user, err := s.userRepository.GetByTelegramID(telegramUserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user does not exis: %v", telegramUserID)
	}

	keys, err := s.keyRepository.GetByUserIDIncludeInactive(user.ID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *KeyService) CreateDefaultKey(KeyID int, telegramUserID int64, serverID uint, connectURL string, keyName string) (*model.Key, error) {
	return s.CreateKeyWithDeadline(KeyID, telegramUserID, serverID, connectURL, keyName, time.Duration(30*time.Minute))
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

func (s *KeyService) GetByID(keyID uint) (*model.Key, error) {
	key, err := s.keyRepository.GetByID(keyID)
	if err != nil {
		return nil, fmt.Errorf("error get key by ID: %v", keyID)
	}

	return key, nil
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

func (s *KeyService) ExtendKeyByIDWithUpdate(keyID uint, timeDuration time.Duration) (*model.Key, error) {
	key, err := s.keyRepository.GetByIDIncludeInactive(keyID)
	if err != nil {
		return nil, fmt.Errorf("error get key by ID: %v", keyID)
	}

	newDeadlineDateTime := time.Now().UTC().Truncate(time.Minute).Add(timeDuration)
	// // Если ключ имеет дату конца жизни позже чем новый, то не изменять его
	// if key.DeadlineTime.After(newDeadlineDateTime) {
	// 	return key, nil
	// }

	key.Duration = timeDuration
	key.DeadlineTime = newDeadlineDateTime
	key.IsActive = true

	updateKeyError := s.keyRepository.Update(key)
	if updateKeyError != nil {
		return nil, fmt.Errorf("error update key: %v", updateKeyError)
	}

	return key, nil
}
