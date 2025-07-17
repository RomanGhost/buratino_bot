package handler

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/bot/data"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/bot/function"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/outline"
	"github.com/RomanGhost/buratino_bot.git/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type KeyHandler struct {
	keyService    *service.KeyService
	serverService *service.ServerService
}

func NewKeyHandler(keyService *service.KeyService, serverService *service.ServerService) *KeyHandler {
	return &KeyHandler{
		keyService:    keyService,
		serverService: serverService,
	}
}

func (h *KeyHandler) ExtendKeyIntline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswerWithDelete(ctx, b, update)

	data := update.CallbackQuery.Data
	keyIDString := strings.Split(data, "_")[1]

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

	h.keyService.ExtendKeyByID(keyIDUint)

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

	// get data from inline
	data := update.CallbackQuery.Data
	shortRegionName := strings.Split(data, "_")[1]

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

	outlineClient := outline.NewOutlineClient(minServer.Access)

	h.createKey(ctx, b, update, minServer.ID, outlineClient)
}

func (h *KeyHandler) createKey(ctx context.Context, b *bot.Bot, update *models.Update, serverID uint, outlineClient *outline.OutlineClient) {
	telegramUser := update.CallbackQuery.From

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

	keyDB, err := h.keyService.CreateKey(key.ID, telegramUser.ID, serverID, connectionKey, key.Name)
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
