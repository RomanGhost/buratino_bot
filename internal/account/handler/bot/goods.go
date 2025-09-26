package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type GoodsHandler struct {
	goodsService *service.GoodsService
}

func NewGoodsHandler(goodsService *service.GoodsService) *GoodsHandler {
	return &GoodsHandler{
		goodsService: goodsService,
	}
}

func (h *GoodsHandler) GetPrices(ctx context.Context, b *bot.Bot, update *models.Update) {
	goods, err := h.goodsService.GetAll(1, 100)
	if err != nil {
		log.Printf("[WARN] Error get goods: %s\n", err)
	}

	message := "Актуальные цены:\n"
	for _, g := range goods {
		if g.Price < 0 {
			continue
		}
		message += fmt.Sprintf("+ %s - %.3f\n", g.Name, float64(g.Price)/1000.0)
	}

	_, sendMessageError := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message,
	})
	if sendMessageError != nil {
		log.Printf("[WARN] Error send message %v", sendMessageError)
	}
}
