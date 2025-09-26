package handler

import (
	"context"
	"log"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/function"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type keyInfo struct {
	ShortRegionID    string
	ProviderID       string
	ServerID         uint
	DeadlineDuration time.Duration
}

func (h *KeyHandler) GetRegionSendProvider(providerFunc func(ctx context.Context, b *bot.Bot, update *models.Update)) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer function.InlineAnswerWithDelete(ctx, b, update)
		telegramUser := update.CallbackQuery.From

		callbackData := update.CallbackQuery.Data
		shortRegionName := callbackData[len(data.RegionChoose):]

		val, ok := h.keyCreatorInfo[telegramUser.ID]
		if !ok {
			h.keyCreatorInfo[telegramUser.ID] = keyInfo{ShortRegionID: shortRegionName}
		} else {
			val.ShortRegionID = shortRegionName
			h.keyCreatorInfo[telegramUser.ID] = val
		}

		providerFunc(ctx, b, update)
	}
}

func (h *KeyHandler) GetProviderSendTime(timeFunc func(ctx context.Context, b *bot.Bot, update *models.Update)) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer function.InlineAnswerWithDelete(ctx, b, update)
		telegramUser := update.CallbackQuery.From
		callbackData := update.CallbackQuery.Data

		val, ok := h.keyCreatorInfo[telegramUser.ID]
		// get provider
		providerName := callbackData[len(data.ProviderChoose):]
		if !ok {
			h.keyCreatorInfo[telegramUser.ID] = keyInfo{ProviderID: providerName}
		} else {
			val.ProviderID = providerName
			h.keyCreatorInfo[telegramUser.ID] = val
		}

		//get server
		if val.ShortRegionID == "" || val.ProviderID == "" {
			forgotUserDataError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
			return
		}

		minServer, err := h.serverService.GetNotLoadedByRegionAndProviderServer(val.ShortRegionID, val.ProviderID)
		if err != nil {
			serverError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
			return
		}

		val.ServerID = minServer.ID
		h.keyCreatorInfo[telegramUser.ID] = val

		// send time choose
		timeFunc(ctx, b, update)
	}
}

func (h *KeyHandler) GetTimeToCreateKey(createKeyFunc func(ctx context.Context, b *bot.Bot, update *models.Update)) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer function.InlineAnswerWithDelete(ctx, b, update)
		telegramUser := update.CallbackQuery.From

		val, ok := h.keyCreatorInfo[telegramUser.ID]
		if !ok {
			forgotUserDataError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
			return
		}

		// get time duration from inline
		callbackData := update.CallbackQuery.Data
		timeDurationStr := callbackData[len(data.TimeChoose):]
		timeDataDuration, err := data.GetDateFromButton(timeDurationStr)
		if err != nil {
			log.Printf("[WARN] Can't parse date from callback: %v\n", err)
			errorTimeChoice(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
			return
		}

		duration := time.Duration(timeDataDuration.Days) * 24 * time.Hour
		duration += time.Duration(timeDataDuration.Hours) * time.Hour
		duration += time.Duration(timeDataDuration.Minutes) * time.Minute

		val.DeadlineDuration = duration
		h.keyCreatorInfo[telegramUser.ID] = val
	}
}
