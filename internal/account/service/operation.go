package service

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/repository"
	apperror "github.com/RomanGhost/buratino_bot.git/internal/app/error"
)

type OperationService struct {
	operationRepository *repository.OperationRepository
	walletService       *WalletService
	goodsService        *GoodsService
}

func NewOperationService(operationRepository *repository.OperationRepository, walletService *WalletService, goodsService *GoodsService) *OperationService {
	return &OperationService{
		operationRepository: operationRepository,
		walletService:       walletService,
		goodsService:        goodsService,
	}
}

func (s *OperationService) TopUpAccount(userID uint, integerPart, fractionalPart uint64) (*model.Operation, error) {
	count := integerPart*1000 + fractionalPart
	return s.CreateOperation(userID, "Пополнение", count)
}

func (s *OperationService) CreateOperation(userID uint, goodsName string, count uint64) (*model.Operation, error) {
	g, err := s.goodsService.GetByName(goodsName)
	if err != nil {
		return nil, apperror.NotFound("Goods not found", err)
	}

	wallet, err := s.walletService.GetByUserID(userID)
	if err != nil {
		return nil, apperror.NotFound("Can't find wallet", err)
	}

	addWalletError := s.walletService.Sub(wallet.ID, g.Price*int64(count))
	if addWalletError != nil {
		return nil, apperror.BadRequest("Account is over", err)
	}

	newOperation := model.Operation{
		WalletID: wallet.ID,
		GoodsID:  g.ID,
		Count:    count,
	}
	createOperationError := s.operationRepository.Create(&newOperation)
	if createOperationError != nil {
		return nil, fmt.Errorf("error create operation: %s", createOperationError)
	}

	return &newOperation, nil
}
