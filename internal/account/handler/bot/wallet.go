package bot

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type WalletHandler struct {
	walletService    *service.WalletService
	operationService *service.OperationService
	userService      *service.UserService
}

func NewWalletHandler(walletService *service.WalletService, operationService *service.OperationService, userService *service.UserService) *WalletHandler {
	return &WalletHandler{
		operationService: operationService,
		userService:      userService,
		walletService:    walletService,
	}
}

func (h *WalletHandler) GetBalace(ctx context.Context, b *bot.Bot, update *models.Update) {
	telegramUser := update.Message.From
	user, err := h.userService.GetUserByTelegramID(telegramUser.ID)
	if err != nil {
		log.Printf("[WARN] Unknown user: %v", telegramUser.Username)
		return // TODO register
	}

	wallet, err := h.walletService.GetByUserID(user.ID)
	if err != nil {
		log.Printf("Error get wallet by userID: %d", user.ID)
	}

	_, sendMessageError := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf(
			"Баланс:%6.2f🪙", float64(wallet.MoneyCount)/1000.0,
		),
	})
	if sendMessageError != nil {
		log.Printf("[WARN] Error send message %v", err)
	}
}

// /pay <amount>
func (h *WalletHandler) PayAmount(ctx context.Context, b *bot.Bot, update *models.Update) {
	telegramUser := update.Message.From
	user, err := h.userService.GetUserByTelegramID(telegramUser.ID)
	if err != nil {
		log.Printf("[WARN] Unknown user: %v", telegramUser.Username)
		return // TODO register
	}

	messageText := update.Message.Text

	re := regexp.MustCompile(`(\d+)([,.](\d{1,2}))?`)
	matches := re.FindStringSubmatch(messageText)
	if matches == nil {
		log.Printf("[WARN] Incorect amount: %v", messageText)
		return // TODO ERROR
	}

	integerPartStr := matches[1]
	integerPart, _ := strconv.ParseUint(integerPartStr, 10, 64)
	fractionalPartStr := matches[3]
	fractionalPart, _ := strconv.ParseUint(fractionalPartStr, 10, 64)

	// TODO Get payment from telegram

	operation, err := h.operationService.TopUpAccount(user.ID, integerPart, fractionalPart)
	if err != nil {
		log.Printf("[ERROR] Can't top up account for user: %d, Error: %s", telegramUser.ID, err)
		return
	}

	_, sendMessageError := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf(
			"Счет успешно пополнен %6.2f🪙!", float64(operation.Count)/1000.0,
		),
	})
	if sendMessageError != nil {
		log.Printf("[WARN] Error send message %v", err)
	}

}
