package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/handler/outline"
	"github.com/RomanGhost/buratino_bot.git/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type KeyHandler struct {
	outline    *outline.OutlineClient
	keyService *service.KeyService
}

func NewKeyHandler(outline *outline.OutlineClient, keyService *service.KeyService) *KeyHandler {
	return &KeyHandler{outline: outline, keyService: keyService}
}

func (h *KeyHandler) CreateKeyInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

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
	//
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
