package handler

import (
	"context"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/service"
	vpnTelegram "github.com/RomanGhost/buratino_bot.git/internal/vpn/telegram/function"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ProviderHandler struct {
	providerService *service.ProviderService
}

func NewProviderHandler(providerService *service.ProviderService) *ProviderHandler {
	return &ProviderHandler{
		providerService: providerService,
	}
}

// function for get Provider of server
func (h *ProviderHandler) GetProvidersInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	// function.InlineAnswerWithDelete(ctx, b, update)

	providers, err := h.providerService.GetProviders()
	if err != nil {
		providerError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	inlineKeyboard := vpnTelegram.GetProvidersInlineKeyboard(providers)

	messageText := `Выбери материал(протокол), из которого нужно изготовить ключик`
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        messageText,
		ReplyMarkup: inlineKeyboard,
	})

	if err != nil {
		log.Printf("[WARN] Error send Provider message %v", err)
	}
}

func providerError(ctx context.Context, b *bot.Bot, chatId int64) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatId,
		Text:      `Возникли проблемы с полученим протоколов, уже чиним\!`,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}
