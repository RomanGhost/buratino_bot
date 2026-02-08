package service

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUserByTelegramID(telegramID int64) (*model.User, error) {
	return s.userRepository.GetByTelegramID(telegramID)
}

func (s *UserService) AddNewUser(telegramID int64, authUserID uint) error {
	user, err := s.userRepository.GetByTelegramID(telegramID)
	if user != nil {
		return fmt.Errorf("user exist")
	}
	if err != nil {
		return err
	}

	newUser := model.User{
		TelegramID: telegramID,
		AuthID:     authUserID,
	}
	_ = s.userRepository.Create(&newUser)

	return nil
}
