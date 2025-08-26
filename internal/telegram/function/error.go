package function

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
)

func UnknownUser(ctx context.Context, b *bot.Bot, chatId int64) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   `Возникла проблема, я тебя забыл, для продолжения выполни команду: /start`,
	})

	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}
