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

func (s *WalletService) CreateWallet(userID uint) (*model.Wallet, error) {
	newWallet := model.Wallet{
		UserID:     userID,
		MoneyCount: 0,
	}

	err := s.walletRepository.Create(&newWallet)
	if err != nil {
		return nil, fmt.Errorf("error create wallet")
	}
	return &newWallet, nil
}

func (s *WalletService) GetByUserID(userID uint) (*model.Wallet, error) {
	w, err := s.walletRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *WalletService) Add(walletID uint, amount uint64) error {
	w, err := s.walletRepository.FindByID(walletID)
	if err != nil {
		return fmt.Errorf("error to get wallet with id: %d, error: %s", walletID, err)
	}
	w.MoneyCount += int64(amount)

	return nil
}

func (s *WalletService) Sub(walletID uint, amount int64) error {
	w, err := s.walletRepository.FindByID(walletID)
	if err != nil {
		return fmt.Errorf("error to get wallet with id: %d, error: %s", walletID, err)
	}

	if w.MoneyCount < amount {
		return fmt.Errorf("wallet balance less that amount: %d", amount)
	}

	w.MoneyCount -= amount

	updateWalletError := s.walletRepository.Update(w)
	if updateWalletError != nil {
		return fmt.Errorf("wallet update error: %s", updateWalletError)
	}

	return nil
}

func (s *WalletService) GetBalance(userID uint) (int64, error) {
	w, err := s.walletRepository.FindByUserID(userID)
	if err != nil {
		return 0, err
	}

	return w.MoneyCount, nil
}
