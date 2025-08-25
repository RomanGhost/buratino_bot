package service

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/repository"
)

type WalletService struct {
	walletRepository *repository.WalletRepository
}

func NewWalletService(walletRepository *repository.WalletRepository) *WalletService {
	return &WalletService{walletRepository}
}

func (s *WalletService) GetByUserID(userID uint) (*model.Wallet, error) {
	w, err := s.walletRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *WalletService) Add(walletID uint, amount uint) error {
	w, err := s.walletRepository.FindByID(walletID)
	if err != nil {
		return fmt.Errorf("error to get wallet with id: %d, error: %s", walletID, err)
	}
	w.MoneyCount += int(amount)

	return nil
}

func (s *WalletService) Sub(walletID uint, amount uint) error {
	w, err := s.walletRepository.FindByID(walletID)
	if err != nil {
		return fmt.Errorf("error to get wallet with id: %d, error: %s", walletID, err)
	}
	if w.MoneyCount < int(amount) {
		return fmt.Errorf("wallet balance less that amount: %d", amount)
	}
	w.MoneyCount -= int(amount)

	return nil
}
