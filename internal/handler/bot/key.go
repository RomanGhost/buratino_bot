package handler

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/bot/function"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/outline"
	"github.com/RomanGhost/buratino_bot.git/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type KeyHandler struct {
	outline       *outline.OutlineClient
	keyService    *service.KeyService
	regionService *service.RegionService
	serverService *service.ServerService
}

func NewKeyHandler(outline *outline.OutlineClient, keyService *service.KeyService, regionService *service.RegionService, serverService *service.ServerService) *KeyHandler {
	return &KeyHandler{
		outline:       outline,
		keyService:    keyService,
		regionService: regionService,
		serverService: serverService,
	}
}

// function for get region of server
func (h *KeyHandler) CreateKeyGetRegionInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswer(ctx, b, update.CallbackQuery.ID)

	regions, err := h.regionService.GetRegionsWithServers()
	if err != nil {
		regionsError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	// regions into buttons
	inlineButtons := [][]models.InlineKeyboardButton{}
	line := []models.InlineKeyboardButton{}
	for i, region := range regions {
		button := models.InlineKeyboardButton{Text: region.RegionName, CallbackData: fmt.Sprintf("choosenRegion_%v", region.ShortName)}
		line = append(line, button)

		if (i+1)%3 == 0 {
			inlineButtons = append(inlineButtons, line)
			line = line[0:0]
		}
	}
	if len(line) > 0 {
		inlineButtons = append(inlineButtons, line)
	}

	// send message
	inlineKeyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineButtons,
	}
	messageText := `Выбери регион, из которого нужно принести ключик`
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        messageText,
		ReplyMarkup: inlineKeyboard,
		ParseMode:   "MarkdownV2",
	})

	if err != nil {
		log.Printf("[WARN] Error send region message %v", err)
	}
}

func (h *KeyHandler) CreateKeyGetServerInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswer(ctx, b, update.CallbackQuery.ID)

	// get data from inline
	data := update.CallbackQuery.Data
	shortRegionName := strings.Split(data, "_")[1]

	servers, err := h.serverService.GetServersByRegionShortName(shortRegionName)
	if err != nil || len(servers) == 0 {
		serverError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

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

	if minServer.ID	 == 0 {
		serverError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}
}

func (h *KeyHandler) CreateKeyInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswer(ctx, b, update.CallbackQuery.ID)

	key, err := h.outline.CreateAccessKey()
	if err != nil {
		log.Printf("[WARN] create outline key: %v\n", err)
		missKeyError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	telegramUser := update.CallbackQuery.From
	key.Name = fmt.Sprintf("%v_%v", telegramUser.ID, time.Now().UTC().Unix())
	err = h.outline.RenameAccessKey(key.ID, key.Name)
	if err != nil {
		log.Printf("[WARN] rename outline key: %v\n", err)
		missKeyError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		return
	}

	connectionKey := key.AccessURL + "&prefix=POST%20"

	_, err = h.keyService.CreateKey(telegramUser.ID, 1, connectionKey)
	if err != nil {
		log.Printf("[WARN] write key in db: %v\n", err)
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text: fmt.Sprintf(
			"🔑 *Вот мой волшебный ключик* \\- держи, не потеряй\\! 🪄\n\n`%s`\n\n_Просто нажми — и он скопируется сам собой\\.\\.\\._ ✨",
			bot.EscapeMarkdown(connectionKey),
		),
		ParseMode: "MarkdownV2",
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
