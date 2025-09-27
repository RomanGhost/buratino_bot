package handler

import (
	"context"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/telegram/function"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/service"
	vpnTelegram "github.com/RomanGhost/buratino_bot.git/internal/vpn/telegram/function"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type RegionHandler struct {
	regionService *service.RegionService
}

func NewRegionHandler(regionService *service.RegionService) *RegionHandler {
	return &RegionHandler{
		regionService: regionService,
	}
}

// function for get region of server
func (h *RegionHandler) GetRegionsInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswerWithDelete(ctx, b, update)

	regions, err := h.regionService.GetRegionsWithServers()
	if err != nil {
		regionError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	inlineKeyboard := vpnTelegram.GetRegionsInlineKeyboard(regions)

	messageText := `Выбери регион, из которого нужно принести ключик`
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        messageText,
		ReplyMarkup: inlineKeyboard,
		ParseMode:   "MarkdownV2",
	})

	if err != nil {
		log.Printf("[WARN] Error send region message %v", err)
	}
}

func regionError(ctx context.Context, b *bot.Bot, chatId int64) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatId,
		Text:      `Возникли проблемы с полученим регионов, уже чиним\!`,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}
