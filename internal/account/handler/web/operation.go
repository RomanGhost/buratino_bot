package web

import (
	"log"
	"net/http"

	"github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/gin-gonic/gin"
)

type Data struct {
	UserID     uint        `json:"user_id"`
	Operations []Operation `json:"operations"`
}

type Operation struct {
	OperationName string `json:"operation_name"`
	Count         uint64 `json:"count"`
}

type OperationHandler struct {
	operationService *service.OperationService
}

func NewOperationHandler(operationService *service.OperationService) *OperationHandler {
	return &OperationHandler{
		operationService: operationService,
	}
}

func (h *OperationHandler) CreateOperation(c *gin.Context) {
	var input Data

	// Привязываем JSON к структуре Data
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Логика обработки (например, печать операций)
	var resultAmount int64
	for _, op := range input.Operations {
		operation, err := h.operationService.CreateOperation(input.UserID, op.OperationName, op.Count)
		if err != nil {
			log.Printf("[WARN] Can't create operation for user: %d, err: %s", input.UserID, err)
		}
		resultAmount += int64(operation.Count) * operation.Goods.Price

	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Данные получены успешно",
		"user":    input.UserID,
		"result":  resultAmount,
	})
}
