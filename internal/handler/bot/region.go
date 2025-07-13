package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/handler/bot/function"
	"github.com/RomanGhost/buratino_bot.git/internal/service"
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

	// regions into buttons
	inlineButtons := [][]models.InlineKeyboardButton{}
	line := []models.InlineKeyboardButton{}
	for i, region := range regions {
		if len(region.Servers) == 0 {
			continue
		}
		button := models.InlineKeyboardButton{Text: region.RegionName, CallbackData: fmt.Sprintf("%v%v", RegionChoose, region.ShortName)}
		line = append(line, button)

		if (i+1)%3 == 0 {
			inlineButtons = append(inlineButtons, line)
			line = line[0:0]
		}
	}

	if len(line) > 0 {
		inlineButtons = append(inlineButtons, line)
	}

	// send message
	inlineKeyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineButtons,
	}
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
