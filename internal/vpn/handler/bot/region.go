package handler

import (
	"context"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/bot/data"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/bot/function"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/service"
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
		regionsError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	inlineKeyboard := data.GetRegionsInlineKeyboard(regions)

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
