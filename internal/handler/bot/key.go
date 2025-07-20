package handler

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/bot/data"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/bot/function"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/outline"
	"github.com/RomanGhost/buratino_bot.git/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type keyInfo struct {
	ServerID         uint
	DeadlineDuration time.Duration
}

type KeyHandler struct {
	keyService     *service.KeyService
	serverService  *service.ServerService
	keyCreatorInfo map[int64]keyInfo
}

func NewKeyHandler(keyService *service.KeyService, serverService *service.ServerService) *KeyHandler {
	keyCreatorInfo := make(map[int64]keyInfo)
	return &KeyHandler{
		keyService:     keyService,
		serverService:  serverService,
		keyCreatorInfo: keyCreatorInfo,
	}
}

func (h *KeyHandler) ExtendKeyIntline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswerWithDelete(ctx, b, update)

	// key Id get
	callbackData := update.CallbackQuery.Data
	keyIDString := callbackData[len(data.ExtendKey):] //strings.Split(data, "_")[1]

	// number check
	keyID, err := strconv.ParseUint(keyIDString, 10, 64)
	if err != nil {
		missKeyError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
	}

	keyIDUint := uint(keyID)

	isActiveKey := h.keyService.IsActiveKey(keyIDUint)
	if !isActiveKey {
		errorExpiredKeys(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
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

func (h *KeyHandler) CreateKeyGetServerInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswerWithDelete(ctx, b, update)

	telegramUser := update.CallbackQuery.From

	// get data from inline
	callbackData := update.CallbackQuery.Data
	shortRegionName := callbackData[len(data.RegionChoose):] //strings.Split(data, "_")[1]

	// get servers by region
	servers, err := h.serverService.GetServersByRegionShortName(shortRegionName)
	if err != nil || len(servers) == 0 {
		serverError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	// chose server with min keys of region
	minCount := math.MaxInt
	var minServer model.Server
	for _, server := range servers {
		val := h.keyService.CountKeysOfServer(server.ID)
		if val == -1 {
			continue
		}
		if minCount > val {
			minCount = val
			minServer = server
		}
	}

	if minServer.ID == 0 {
		serverError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	// переписать для пользователя его данные по серверу
	h.keyCreatorInfo[telegramUser.ID] = keyInfo{ServerID: minServer.ID}

	zeroTimeKeyboard := data.GetCustomTimeKeyboard(&data.TimeDataDuration{30, 0, 0})
	messageText := `Выбери время\!`
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        messageText,
		ReplyMarkup: zeroTimeKeyboard,
		ParseMode:   "MarkdownV2",
	})

	if err != nil {
		log.Printf("[WARN] Error send notify message %v", err)
	}
	// h.createKey(ctx, b, update)
}

func (h *KeyHandler) CreateKeyGetTimeInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswerWithDelete(ctx, b, update)

	telegramUser := update.CallbackQuery.From

	// get data from inline
	callbackData := update.CallbackQuery.Data
	timeDurationStr := callbackData[len(data.CreateTime):]
	timeDataDuration, err := data.GetDateFromButton(timeDurationStr)
	if err != nil {
		log.Printf("[WARN] Can't parse date from callback: %v\n", err)
		errorTimeChoice(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	_, exist := h.keyCreatorInfo[telegramUser.ID]
	if !exist {
		log.Printf("[WARN] user %d, can't go to next step", telegramUser.ID)
		errorSkipStep(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	info := h.keyCreatorInfo[telegramUser.ID]

	duration := time.Duration(timeDataDuration.Days) * 24 * time.Hour
	duration += time.Duration(timeDataDuration.Hours) * time.Hour
	duration += time.Duration(timeDataDuration.Minutes) * time.Minute

	info.DeadlineDuration = duration
	h.keyCreatorInfo[telegramUser.ID] = info

	h.createKey(ctx, b, update)
}

func (h *KeyHandler) createKey(ctx context.Context, b *bot.Bot, update *models.Update) {
	telegramUser := update.CallbackQuery.From
	val, ok := h.keyCreatorInfo[telegramUser.ID]
	if !ok {
		// отправить в начало
		serverError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	server, err := h.serverService.GetServerByID(val.ServerID)
	if err != nil {
		serverError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	outlineClient := outline.NewOutlineClient(server.Access)

	// generate new keys with name
	key, err := outlineClient.CreateAccessKey()
	if err != nil {
		log.Printf("[WARN] create outline key: %v\n", err)
		missKeyError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	key.Name = fmt.Sprintf("%v_%v", telegramUser.ID, time.Now().UTC().Unix())
	err = outlineClient.RenameAccessKey(key.ID, key.Name)
	if err != nil {
		log.Printf("[WARN] Can't rename outline key: %v\n", err)
		missKeyError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	connectionKey := key.AccessURL + "&prefix=POST%20"

	keyDB, err := h.keyService.CreateKeyWithDeadline(key.ID, telegramUser.ID, server.ID, connectionKey, key.Name, val.DeadlineDuration)
	if err != nil {
		log.Printf("[WARN] Can't write key in db: %v\n", err)
		return
	}

	// notify users
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text: fmt.Sprintf(
			"🔑 *Вот мой волшебный ключик №%d* \\- держи, не потеряй\\! 🪄\n\n`%s`\n\n_Просто нажми \\- и он скопируется сам собой\\.\\.\\._ ✨",
			keyDB.ID, bot.EscapeMarkdown(connectionKey),
		),
		ParseMode: "MarkdownV2",
	})

	if err != nil {
		log.Printf("[WARN] Error send key message %v", err)
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

func regionsError(ctx context.Context, b *bot.Bot, chatId int64) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatId,
		Text:      `Возникли проблемы с полученим регионов, уже чиним\!`,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}

func serverError(ctx context.Context, b *bot.Bot, chatId int64) {
	log.Printf("[WARN] Error get server of region")
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatId,
		Text:      `Возникли проблемы с полученим серверов выбранного региона, уже чиним\!`,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}

func missKeyError(ctx context.Context, b *bot.Bot, chatId int64) {
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

func errorSkipStep(ctx context.Context, b *bot.Bot, chatId int64) {
	inlineKeyboard := data.CreateKeyboard([]models.InlineKeyboardButton{data.CreateKeyButton()})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        `Был пропущен шаг при выборе ключа, придется начать сначала`,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: inlineKeyboard,
	})

	if err != nil {
		log.Printf("[WARN] Error send info error message %v", err)
	}
}
