package service

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/repository"
)

type OperationService struct {
	operationRepository *repository.OperationRepository
	walletService       *WalletService
	goodsService        *GoodsService
}

func NewOperationService(operationRepository *repository.OperationRepository) *OperationService {
	return &OperationService{
		operationRepository: operationRepository,
	}
}

func (s *OperationService) CreateOperation(goodsName string, userID uint) (*model.Operation, error) {
	g, err := s.goodsService.GetByName(goodsName)
	if err != nil {
		return nil, fmt.Errorf("error get goods: %s", err)
	}

	wallet, err := s.walletService.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("can't find wallet")
	}

	addWalletError := s.walletService.Sub(wallet.ID, uint(g.Price))
	if addWalletError != nil {
		return nil, fmt.Errorf("error sub of wallet: %s", addWalletError)
	}

	newOperation := model.Operation{
		WalletID: wallet.ID,
		GoodsID:  g.ID,
	}
	createOperationError := s.operationRepository.Create(&newOperation)
	if createOperationError != nil {
		return nil, fmt.Errorf("error create operation: %s", createOperationError)
	}

	return &newOperation, nil
}
