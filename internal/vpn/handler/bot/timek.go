package handler

import (
	"context"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func KeyboardTimeChoose(ctx context.Context, b *bot.Bot, update *models.Update) {
	zeroTimeKeyboard := data.GetCustomTimeKeyboard(&data.TimeDataDuration{Minutes: 30, Hours: 0, Days: 0})

	messageText := `Выбери время\!`
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        messageText,
		ReplyMarkup: zeroTimeKeyboard,
		ParseMode:   "MarkdownV2",
	})

	if err != nil {
		log.Printf("[WARN] Error send notify message %v", err)
	}
}
