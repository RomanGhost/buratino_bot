package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	accountService "github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/RomanGhost/buratino_bot.git/internal/app/timework"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/function"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type KeyHandler struct {
	userService             *service.UserService
	keyService              *service.KeyService
	serverService           *service.ServerService
	accountOperationService *accountService.OperationService
	walletService           *accountService.WalletService
	keyCreatorInfo          map[int64]keyInfo
}

func NewKeyHandler(
	userService *service.UserService,
	keyService *service.KeyService,
	serverService *service.ServerService,
	accountOperationService *accountService.OperationService,
	walletService *accountService.WalletService,
) *KeyHandler {
	keyCreatorInfo := make(map[int64]keyInfo)
	return &KeyHandler{
		userService:             userService,
		keyService:              keyService,
		serverService:           serverService,
		accountOperationService: accountOperationService,
		keyCreatorInfo:          keyCreatorInfo,
		walletService:           walletService,
	}
}

func (h *KeyHandler) ExtendKeyInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	defer function.InlineAnswerWithDelete(ctx, b, update)

	// key Id get
	callbackData := update.CallbackQuery.Data
	keyIDString := callbackData[len(data.ExtendKey):] //strings.Split(data, "_")[1]

	// number check
	keyID, err := strconv.ParseUint(keyIDString, 10, 64)
	if err != nil {
		errorMissKey(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
	}

	keyIDUint := uint(keyID)

	isActiveKey := h.keyService.IsActiveKey(keyIDUint)
	if !isActiveKey {
		errorExpiredKeys(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	telegramUser := update.CallbackQuery.From
	keyVal, err := h.keyService.GetByID(keyIDUint)
	if err != nil {
		errorMissKey(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	resultDuration := h.makeRequest(telegramUser.ID, keyVal.Duration)
	if resultDuration == 0 {
		// Вернуть ошибку баланса пользователю и не выполнять действий
		function.BalanceOver(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	_, errExtendKey := h.keyService.ExtendKeyByID(keyIDUint)
	if errExtendKey != nil {
		errorExpiredKeys(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	messageText := `Ключик продлен\!`
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		Text:      messageText,
		ParseMode: "MarkdownV2",
	})

	if err != nil {
		log.Printf("[WARN] Error send notify message %v", err)
	}
}

func SendNotifyAboutDeadline(ctx context.Context, b *bot.Bot, chatID int64, keyID uint) {
	inlineKeyboard := data.CreateKeyboard([]models.InlineKeyboardButton{data.ExtendKeyButton(keyID)})

	// notify users
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text: fmt.Sprintf(
			"Ключ №%d  скоро совсем испарится, нажми *продлить*, чтобы продолжать пользоваться",
			keyID,
		),
		ParseMode:   "MarkdownV2",
		ReplyMarkup: inlineKeyboard,
	})

	if err != nil {
		log.Printf("[WARN] Error send key message %v", err)
	}
}

func CreateKeyInlineShutdown(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text: `🔧 *В разработке* 🔮

		Тссс\.\.\. *Буратино колдует над новыми чудесами* 🧙‍♂️✨  
		Скоро здесь появится нечто волшебное, что поможет тебе ещё проще и быстрее получать тайные ключики от свободного интернета 🌍🔑

		*Потерпи немного, добрый странник* \- магия требует времени\! ⏳`,
		ParseMode: models.ParseModeMarkdown,
	})

	if err != nil {
		log.Printf("[WARN] Error send key message %v", err)
	}
}

func formatDuration(timeDuration time.Duration) string {
	cd := timework.ConcrateDuration(timeDuration)

	result := fmt.Sprintf("%02d:%02d", cd.Hours, cd.Minutes)
	if cd.Days > 0 {
		result = fmt.Sprintf("%dд %s", cd.Days, result)
	}
	if cd.Months > 0 {
		result = fmt.Sprintf("%dм %s", cd.Months, result)
	}

	return result
}

func errorForgotUserData(ctx context.Context, b *bot.Bot, chatId int64) {
	log.Printf("[WARN] Error get values from map")
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatId,
		Text:      `Я все забыл, давай по новой\!`,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}

func errorServer(ctx context.Context, b *bot.Bot, chatId int64) {
	log.Printf("[WARN] Error get server")
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatId,
		Text:      `Возникли проблемы со сбором серверов, уже чиним\!`,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}

func errorMissKey(ctx context.Context, b *bot.Bot, chatId int64) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text: `🏃‍♂️💨 Пока я к тебе бежал, *ключик куда\-то выскользнул*\.\.\. 🔑😱  
Но не беда\! *Поиски уже ведутся* \- я задействовал всех сверчков, псов и даже Дуремара с его лягушками 🕵️‍♂️🐸

*Чуток терпения, друг мой* \- скоро всё найдётся, и волшебство продолжится ✨`,
		ParseMode: models.ParseModeMarkdown,
	})

	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}

func errorExpiredKeys(ctx context.Context, b *bot.Bot, chatId int64) {
	inlineKeyboard := data.CreateKeyboard([]models.InlineKeyboardButton{data.CreateKeyButton()})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        `Увы ключ совсем заржавел, придется создать новый`,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: inlineKeyboard,
	})

	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}

func errorTimeChoice(ctx context.Context, b *bot.Bot, chatId int64) {
	inlineKeyboard := data.CreateKeyboard([]models.InlineKeyboardButton{data.CreateKeyButton()})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        `Какая-то проблема с выбором времени, пересоздай ключ`,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: inlineKeyboard,
	})

	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}
