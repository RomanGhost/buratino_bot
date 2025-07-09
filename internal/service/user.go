package service

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/database/repository"
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

func (s *UserService) AddNewUser(telegramID int64) error {
	user, err := s.userRepository.GetByTelegramID(telegramID)
	if user != nil {
		return fmt.Errorf("user exist")
	}
	if err != nil {
		return err
	}

	userRole, err := s.userRoleRepository.GetByRoleName("user")
	if err != nil {
		return err
	}

	newUser := model.User{
		TelegramID: telegramID,
		IsActive:   true,
		UserRole:   *userRole,
	}
	s.userRepository.Create(&newUser)

	return nil
}
