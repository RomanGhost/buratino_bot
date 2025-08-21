package service

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/repository"
)

type UserService struct {
	userRepository     *repository.UserRepository
	userRoleRepository *repository.UserRoleRepository
}

func NewUserService(userRepository *repository.UserRepository, userRoleRepository *repository.UserRoleRepository) *UserService {
	return &UserService{
		userRepository:     userRepository,
		userRoleRepository: userRoleRepository,
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
	return &newUser, nil
}

func (s *UserService) ExistUserByTelegramID(telegramID int64) bool {
	u, err := s.userRepository.FindByTelegramID(telegramID)
	if err == nil && u != nil {
		return true
	}
	return false
}
