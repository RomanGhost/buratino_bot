package function

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func InlineAnswerWithDelete(ctx context.Context, b *bot.Bot, update *models.Update) {
	InlineAnswer(ctx, b, update)

	// Удаляем сообщение с inline кнопкой
	_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})
	if err != nil {
		log.Printf("[WARN] can't delete message: %v", err)
	}
}

func InlineAnswer(ctx context.Context, b *bot.Bot, update *models.Update) {
	// inline answer
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		log.Println("[WARN] can't answer on inline request")
	}
}
