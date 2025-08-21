package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// first message
func (h *UserHandler) RegisterUser(ctx context.Context, b *bot.Bot, update *models.Update) {
	telegramUser := update.Message.From
	log.Printf("[INFO] Regist user: %v, ID: %v", telegramUser.Username, telegramUser.ID)

	_, err := h.userService.RegisterUser(telegramUser.ID, telegramUser.Username)
	if err != nil {
		log.Printf("[WARN] user register error: %v", err)
	}

	inlineKeyboard := data.CreateKeyboard(
		[]models.InlineKeyboardButton{data.CreateKeyButton()},
		[]models.InlineKeyboardButton{data.AboutOutlineButton(), data.KnowProjectButton()},
	)

	_, sendMessageError := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf(
			"🎭 *Здравствуй, %v*\\!\n\nЯ \\- _Буратино_, не простой деревянный мальчишка, а хранитель волшебных ключей от потайных дверей интернета\\! 🌍✨\n\nВ этом сказочном чате ты сможешь получить *волшебный VPN\\-ключ*, который укроет тебя от злых Карабасов и злобных Брандмейстеров 🔒\n\n🔑 *Нажми на волшебную кнопку*, и я подарю тебе ключик… но помни:\n\\- он действует лишь ограниченное время ⏳\n\\- добро должно быть с кулаками… но под шифрованием\\! 🛡️\n\n🧙‍♂️ Поехали в страну свободного интернета\\!",
			bot.EscapeMarkdown(telegramUser.Username),
		),
		ReplyMarkup: inlineKeyboard,
		ParseMode:   models.ParseModeMarkdown,
	})
	if sendMessageError != nil {
		log.Printf("[WARN] Error send message %v", err)
	}
}
