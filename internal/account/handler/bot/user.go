package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	vpnService "github.com/RomanGhost/buratino_bot.git/internal/vpn/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserHandler struct {
	userService    *service.UserService
	userVPNService *vpnService.UserService
}

func NewUserHandler(userService *service.UserService, userVPNService *vpnService.UserService) *UserHandler {
	return &UserHandler{
		userService:    userService,
		userVPNService: userVPNService,
	}
}

func (h *UserHandler) MiddleWareLookup(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		h.lookupUserChange(update)
		next(ctx, b, update)
	}
}

// check userChange
func (h *UserHandler) lookupUserChange(update *models.Update) {
	if update != nil && update.Message != nil && update.Message.From != nil {
		telegramUser := update.Message.From
		log.Printf("[INFO] Regist user: %v, ID: %v", telegramUser.Username, telegramUser.ID)

		user, err := h.userService.GetOrCreateUser(telegramUser.ID, telegramUser.Username)
		if err != nil {
			log.Printf("[WARN] user register error: %v", err)
		} else {
			_ = h.userVPNService.AddNewUser(telegramUser.ID, user.ID)
		}
	}
}

// stsrt message
func (h *UserHandler) RegisterUser(ctx context.Context, b *bot.Bot, update *models.Update) {
	telegramUser := update.Message.From
	h.lookupUserChange(update)

	inlineKeyboard := data.CreateKeyboard(
		[]models.InlineKeyboardButton{data.KnowProjectButton()},
		[]models.InlineKeyboardButton{data.AboutOutlineButton(), data.AboutWireguardButton()},
		[]models.InlineKeyboardButton{data.CreateKeyButton()},
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
		log.Printf("[WARN] Error send message %v", sendMessageError)
	}
}
