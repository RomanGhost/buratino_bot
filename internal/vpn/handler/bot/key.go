package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	accountService "github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/RomanGhost/buratino_bot.git/internal/app/timework"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/function"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/provider"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type KeyHandler struct {
	userService             *service.UserService
	keyService              *service.KeyService
	serverService           *service.ServerService
	accountOperationService *accountService.OperationService
	keyCreatorInfo          map[int64]keyInfo
}

func NewKeyHandler(userService *service.UserService, keyService *service.KeyService, serverService *service.ServerService, accountOperationService *accountService.OperationService) *KeyHandler {
	keyCreatorInfo := make(map[int64]keyInfo)
	return &KeyHandler{
		userService:             userService,
		keyService:              keyService,
		serverService:           serverService,
		accountOperationService: accountOperationService,
		keyCreatorInfo:          keyCreatorInfo,
	}
}

func (h *KeyHandler) ExtendKeyIntline(ctx context.Context, b *bot.Bot, update *models.Update) {
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

func (h *KeyHandler) CreateKey(ctx context.Context, b *bot.Bot, update *models.Update) {
	telegramUser := update.CallbackQuery.From
	defer delete(h.keyCreatorInfo, telegramUser.ID)

	val, ok := h.keyCreatorInfo[telegramUser.ID]
	if !ok {
		// отправить в начало
		errorForgotUserData(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	resultDuration := h.makeRequest(telegramUser.ID, val.DeadlineDuration)
	if resultDuration == 0 {
		// Вернуть ошибку баланса пользователю и не выполнять действий
		function.BalanceOver(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	server, err := h.serverService.GetServerByID(val.ServerID)
	if err != nil {
		errorServer(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	providerClient := provider.NewProvider(server.Access, server.ProviderID)

	newKeyName := fmt.Sprintf("%d", telegramUser.ID)
	connectionKey, err := providerClient.CreateKey(newKeyName)
	log.Println("[DEBUG] created key", connectionKey)
	if err != nil {
		log.Printf("[WARN] Can't create key: %v\n", err)
		errorMissKey(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	keyDB, err := h.keyService.CreateKeyWithDeadline(connectionKey.ID, telegramUser.ID, server.ID, connectionKey.ConnectData, connectionKey.Name, resultDuration)
	if err != nil {
		log.Printf("[WARN] Can't write key in db: %v\n", err)
		return
	}

	switch server.ProviderID {
	case model.Outline.Name:
		sendKeyOutline(ctx, b, update, keyDB)
	case model.Wireguard.Name:
		sendKeyWireguard(ctx, b, update, keyDB)
	default:
		errorServer(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
	}

}

func sendKeyOutline(ctx context.Context, b *bot.Bot, update *models.Update, keyData *model.Key) {
	// notify users
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text: fmt.Sprintf(
			"🔑 *Вот мой волшебный **Outline** ключик №%d* \\- держи, не потеряй\\! 🪄\n`%s`\n⌚ Время жизни: %s\n_Просто нажми \\- и он скопируется сам собой\\.\\.\\._ ✨",
			keyData.ID, bot.EscapeMarkdown(keyData.ConnectUrl), bot.EscapeMarkdown(formatDuration(keyData.Duration)),
		),
		ParseMode: "MarkdownV2",
	})

	if err != nil {
		log.Printf("[WARN] Error send key message %v", err)
	}
}

func sendKeyWireguard(ctx context.Context, b *bot.Bot, update *models.Update, keyData *model.Key) {
	fileName := fmt.Sprintf("key%d.conf", keyData.ID)
	tempFile, err := os.CreateTemp("./cache", fmt.Sprintf("*-%s", fileName))
	if err != nil {
		log.Printf("[WARN] error create temp file: %v", err)
		errorMissKey(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}
	// defer os.Remove(tmpFile.Name()) // удаляем после использования
	defer tempFile.Close()

	_, err = tempFile.WriteString(keyData.ConnectUrl)
	if err != nil {
		log.Printf("[WARN] error write to temp file: %v", err)
		errorMissKey(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	_, err = tempFile.Seek(0, 0) // переместить курсор в начало
	if err != nil {
		log.Printf("[WARN] error seek temp file: %v", err)
		errorMissKey(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	textMessage := fmt.Sprintf(
		"🔑 *Вот мой волшебный Wireguard ключик №%d* \\- держи, не потеряй\\! 🪄\n"+
			"⌚ Время жизни: %s ✨",
		keyData.ID,
		bot.EscapeMarkdown(formatDuration(keyData.Duration)),
	)

	_, err = b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Document: &models.InputFileUpload{
			Filename: fileName, // имя файла, которое увидит пользователь
			Data:     tempFile, // сам файл
		},
		Caption:   textMessage,
		ParseMode: "MarkdownV2",
	})
	if err != nil {
		log.Printf("[ERROR] send document: %v", err)
	}
}

func (h *KeyHandler) makeRequest(telegramID int64, timeDuration time.Duration) time.Duration {
	user, err := h.userService.GetUserByTelegramID(telegramID)
	if err != nil {
		return 0
	}

	cd := timework.ConcrateDuration(timeDuration)
	var res time.Duration

	_, minError := h.accountOperationService.CreateOperation(user.AuthID, "1m vpn", uint64(cd.Minutes))
	if minError != nil {
		return 0
	}
	res += cd.Minutes * time.Minute

	_, hourError := h.accountOperationService.CreateOperation(user.AuthID, "1h vpn", uint64(cd.Hours))
	if hourError != nil {
		return res
	}
	res += cd.Hours * time.Hour

	_, dayError := h.accountOperationService.CreateOperation(user.AuthID, "1d vpn", uint64(cd.Days))
	if dayError != nil {
		return res
	}
	res += cd.Days * timework.DayDuration

	_, monthError := h.accountOperationService.CreateOperation(user.AuthID, "1month vpn", uint64(cd.Months))
	if monthError != nil {
		return res
	}
	res += cd.Months * timework.MonthDuration

	return timeDuration
}

func (h *KeyHandler) countPrice(timeDuration time.Duration) int64 {

	var resPrice int64
	cd := timework.ConcrateDuration(timeDuration)

	minPrice, minError := h.accountOperationService.GetPrice("1m vpn", uint64(cd.Minutes))
	if minError != nil {
		return 0
	}
	resPrice += minPrice

	hourPrice, hourError := h.accountOperationService.GetPrice("1h vpn", uint64(cd.Hours))
	if hourError != nil {
		return resPrice
	}
	resPrice += hourPrice

	dayPrice, dayError := h.accountOperationService.GetPrice("1d vpn", uint64(cd.Days))
	if dayError != nil {
		return resPrice
	}
	resPrice += dayPrice

	monthPrice, monthError := h.accountOperationService.GetPrice("1month vpn", uint64(cd.Months))
	if monthError != nil {
		return resPrice
	}
	resPrice += monthPrice

	return resPrice
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
