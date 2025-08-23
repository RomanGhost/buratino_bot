package service

import "github.com/RomanGhost/buratino_bot.git/internal/account/database/repository"

type OperationService struct {
	operationRepository *repository.OperationRepository
}
