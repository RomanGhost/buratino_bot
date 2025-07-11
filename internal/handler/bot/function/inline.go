package function

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
)

func InlineAnswer(ctx context.Context, b *bot.Bot, callbackID string) {
	// inline answer
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackID,
		ShowAlert:       false,
	})
	if err != nil {
		log.Println("[WARN] can't answer on inline request")
	}
}

