package service

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/repository"
)

type UserService struct {
	userRepository     *repository.UserRepository
	userRoleRepository *repository.UserRoleRepository
	walletService      *WalletService
}

func NewUserService(userRepository *repository.UserRepository, userRoleRepository *repository.UserRoleRepository, walletService *WalletService) *UserService {
	return &UserService{
		userRepository:     userRepository,
		userRoleRepository: userRoleRepository,
		walletService:      walletService,
	}
}

func (s *UserService) RegisterUser(telegramID int64, username string) (*model.User, error) {
	exist := s.ExistUserByTelegramID(telegramID)
	if exist {
		return nil, fmt.Errorf("error user exist")
	}
	newUser := model.User{
		TelegramID:       telegramID,
		TelegramUsername: username,
		IsActive:         true,
		Role:             "user",
		TimezoneOffset:   3 * 60,
	}
	err := s.userRepository.Create(&newUser)
	if err != nil {
		return nil, fmt.Errorf("error create user: %s", err)
	}

	_, createWalletError := s.walletService.CreateWallet(newUser.ID)
	if createWalletError != nil {
		return nil, fmt.Errorf("error create wallet for user: %d, error: %s", newUser.ID, createWalletError)
	}
	return &newUser, nil
}

func (s *UserService) GetUserByTelegramID(telegramID int64) (*model.User, error) {
	u, err := s.userRepository.FindByTelegramID(telegramID)
	if err != nil {
		return nil, fmt.Errorf("error get user by telegramID: %d", telegramID)
	}
	return u, nil
}

func (s *UserService) ExistUserByTelegramID(telegramID int64) bool {
	u, err := s.userRepository.FindByTelegramID(telegramID)
	if err == nil && u != nil {
		return true
	}
	return false
}
