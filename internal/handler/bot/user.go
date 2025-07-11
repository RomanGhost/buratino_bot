package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) RegisterUser(ctx context.Context, b *bot.Bot, update *models.Update) {
	telegramUser := update.Message.From
	log.Printf("[INFO] Registe user: %v, ID: %v", telegramUser.Username, telegramUser.ID)
	if err := h.userService.AddNewUser(telegramUser.ID); err != nil {
		log.Printf("[WARN] user register error: %v", err)
	}

	inlineKeyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Создать ключ", CallbackData: "createKey"},
			}, {
				{Text: "Узнать о проекте", CallbackData: "infoProject"},
			},
		},
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf(
			"🎭 *Здравствуй, %v*\\!\n\nЯ \\- _Буратино_, не простой деревянный мальчишка, а хранитель волшебных ключей от потайных дверей интернета\\! 🌍✨\n\nВ этом сказочном чате ты сможешь получить *волшебный VPN\\-ключ*, который укроет тебя от злых Карабасов и злобных Брандмейстеров 🔒\n\n🔑 *Нажми на волшебную кнопку*, и я подарю тебе ключик… но помни:\n\\- он действует лишь ограниченное время ⏳\n\\- добро должно быть с кулаками… но под шифрованием\\! 🛡️\n\n🧙‍♂️ Поехали в страну свободного интернета\\!",
			bot.EscapeMarkdown(telegramUser.Username),
		),
		ReplyMarkup: inlineKeyboard,
		ParseMode:   models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("[WARN] Error send message %v", err)
	}
}
