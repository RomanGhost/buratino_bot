package handler

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func InfoAboutInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	inlineKeyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Создать ключ", CallbackData: "create_key"},
			},
		},
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text: `📜 *Сказ о волшебных ключах* 🗝️
В этой сказочной обители ты встретил *Буратино* \- не просто деревянного мальчишку, а стража потайных троп интернета\! 🌐✨
Он дарует *волшебные VPN\-ключи*, что действуют недолго \- всего около *30 минут*, но дают силу обойти коварных Карабасов и Брандмейстеров\.
Сейчас ключики выдаются через *Outline*, но скоро и *WireGuard* придёт на помощь храбрым странникам\!
Нажми на волшебную кнопку, и путь будет открыт\.\.\. 🧙‍♂️🔑`,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: inlineKeyboard,
	})

	if err != nil {
		log.Printf("[WARN] Error send info message %v", err)
	}

}
