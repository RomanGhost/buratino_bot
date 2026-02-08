package function

import (
	"context"
	"fmt"
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

func BalanceOver(ctx context.Context, b *bot.Bot, chatId int64) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   `Возникла проблема, недостаточен баланс, для продолжения выполни команду: /balance`,
	})

	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}

func BalanceOverAddInfo(ctx context.Context, b *bot.Bot, chatId int64, available, needMoney int64) {
	BalanceOver(ctx, b, chatId)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		// TODO вынести в отдельную функцию расчет баланса(деление на 1000)
		Text: fmt.Sprintf("Сумма на балансе: %6.2f🪙\nНеобходимо для оплаты: %6.2f🪙", float64(available)/1000.0, float64(needMoney)/1000.0),
	})

	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}
