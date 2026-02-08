package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/app/timework"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/function"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/provider"
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
			errorForgotUserData(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
			return
		}

		minServer, err := h.serverService.GetNotLoadedByRegionAndProviderServer(val.ShortRegionID, val.ProviderID)
		if err != nil {
			errorServer(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
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
			errorForgotUserData(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
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

		createKeyFunc(ctx, b, update)
	}
}

func (h *KeyHandler) CreateKeyIfNotExists(ctx context.Context, b *bot.Bot, update *models.Update) {
	telegramUser := update.CallbackQuery.From
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	keys, err := h.keyService.GetKeysByTelegramUserID(telegramUser.ID)
	if err != nil {
		errorServer(ctx, b, chatID)
		return
	}

	if len(keys) == 0 {
		h.createOrExtendKey(ctx, b, update, nil)
		return
	}

	val, ok := h.keyCreatorInfo[telegramUser.ID]
	if !ok {
		errorForgotUserData(ctx, b, chatID)
		return
	}

	server, err := h.serverService.GetServerByID(val.ServerID)
	if err != nil {
		errorServer(ctx, b, chatID)
		return
	}

	for _, key := range keys {
		if !key.IsActive && key.ServerID == server.ID {
			h.createOrExtendKey(ctx, b, update, &key)
			return
		}
	}

	h.createOrExtendKey(ctx, b, update, nil)
}

func (h *KeyHandler) createOrExtendKey(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	existingKey *model.Key, // nil → создать новый
) {
	telegramUser := update.CallbackQuery.From
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	defer delete(h.keyCreatorInfo, telegramUser.ID)

	val, ok := h.keyCreatorInfo[telegramUser.ID]
	if !ok {
		errorForgotUserData(ctx, b, chatID)
		return
	}

	// TODO перенести в кеш - нагрузка на БД
	user, err := h.userService.GetUserByTelegramID(telegramUser.ID)
	if err != nil {
		return
	}

	// получение цены
	totalPrice := h.countPrice(val.DeadlineDuration)
	userBalance, err := h.walletService.GetBalance(user.AuthID)
	if err != nil {
		errorServer(ctx, b, chatID)
		return
	}

	if totalPrice > userBalance {
		function.BalanceOverAddInfo(ctx, b, chatID, userBalance, totalPrice)
		return
	}

	// расчет стоимости за время жизни ключа
	resultDuration := h.makeRequest(telegramUser.ID, val.DeadlineDuration)
	if resultDuration == 0 {
		function.BalanceOver(ctx, b, chatID)
		return
	}

	server, err := h.serverService.GetServerByID(val.ServerID)
	if err != nil {
		errorServer(ctx, b, chatID)
		return
	}

	var key *model.Key

	if existingKey == nil {
		// CREATE
		providerClient := provider.NewProvider(server.Access, server.ProviderID)

		connectionKey, err := providerClient.CreateKey(fmt.Sprintf("%d", telegramUser.ID))
		if err != nil {
			log.Printf("[WARN] Can't create key: %v\n", err)
			errorMissKey(ctx, b, chatID)
			return
		}

		key, err = h.keyService.CreateKeyWithDeadline(
			connectionKey.ID,
			telegramUser.ID,
			server.ID,
			connectionKey.ConnectData,
			connectionKey.Name,
			resultDuration,
		)

		if err != nil {
			log.Printf("[WARN] Can't write key in db: %v\n", err)
			errorMissKey(ctx, b, chatID)
			return
		}
	} else {
		// EXTEND
		key, err = h.keyService.ExtendKeyByIDWithUpdate(existingKey.ID, resultDuration)
		if err != nil {
			log.Printf("[WARN] Can't extend key: %v\n", err)
			errorMissKey(ctx, b, chatID)
			return
		}
	}

	switch server.ProviderID {
	case model.Outline.Name:
		sendKeyOutline(ctx, b, update, key)
	case model.Wireguard.Name:
		sendKeyWireguard(ctx, b, update, key)
	default:
		errorServer(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
	}
}

func (h *KeyHandler) makeRequest(telegramID int64, timeDuration time.Duration) time.Duration {
	user, err := h.userService.GetUserByTelegramID(telegramID)
	if err != nil {
		return 0
	}

	concDuration := timework.ConcrateDuration(timeDuration)
	var resDuration time.Duration

	_, minError := h.accountOperationService.CreateOperation(user.AuthID, VPN1Min, uint64(concDuration.Minutes))
	if minError != nil {
		return 0
	}
	resDuration += time.Duration(concDuration.Minutes) * time.Minute

	_, hourError := h.accountOperationService.CreateOperation(user.AuthID, VPN1Hour, uint64(concDuration.Hours))
	if hourError != nil {
		return resDuration
	}
	resDuration += time.Duration(concDuration.Hours) * time.Hour

	_, dayError := h.accountOperationService.CreateOperation(user.AuthID, VPN1Day, uint64(concDuration.Days))
	if dayError != nil {
		return resDuration
	}
	resDuration += time.Duration(concDuration.Days) * timework.DayDuration

	_, monthError := h.accountOperationService.CreateOperation(user.AuthID, VPN1Month, uint64(concDuration.Months))
	if monthError != nil {
		return resDuration
	}
	resDuration += time.Duration(concDuration.Months) * timework.MonthDuration

	return timeDuration
}

func (h *KeyHandler) countPrice(timeDuration time.Duration) int64 {

	var resPrice int64
	cd := timework.ConcrateDuration(timeDuration)

	minPrice, minError := h.accountOperationService.GetPrice(VPN1Min, uint64(cd.Minutes))
	if minError != nil {
		return 0
	}
	resPrice += minPrice

	hourPrice, hourError := h.accountOperationService.GetPrice(VPN1Hour, uint64(cd.Hours))
	if hourError != nil {
		return resPrice
	}
	resPrice += hourPrice

	dayPrice, dayError := h.accountOperationService.GetPrice(VPN1Day, uint64(cd.Days))
	if dayError != nil {
		return resPrice
	}
	resPrice += dayPrice

	monthPrice, monthError := h.accountOperationService.GetPrice(VPN1Month, uint64(cd.Months))
	if monthError != nil {
		return resPrice
	}
	resPrice += monthPrice

	return resPrice
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
