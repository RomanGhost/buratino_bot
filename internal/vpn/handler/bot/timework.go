package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/bot/function"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type timeInfoOperations struct {
	NowDataTime uint16
	DeltaTime   uint16
	TimeUnit    *data.TimeUnit
}

func AddTimeInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswer(ctx, b, update)

	callbackData := update.CallbackQuery.Data
	message := update.CallbackQuery.Message.Message
	infoOperations, err := dataFromCallback(message, callbackData[len(data.TimeAdd):])
	if err != nil {
		log.Printf("[WARN] Can't get data date")
		return
	}

	result := infoOperations.NowDataTime + infoOperations.DeltaTime
	if result >= infoOperations.TimeUnit.MaxValue {
		result = 0
	}

	inlineKeyboard, _ := data.UpdateTimeKeyboard(result, infoOperations.TimeUnit, message.ReplyMarkup)

	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      message.Chat.ID,
		MessageID:   message.ID,
		ReplyMarkup: inlineKeyboard,
	})
	if err != nil {
		log.Printf("[WARN] Can't edit message, err:%v\n", err)
	}
}

func ReduceTimeInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswer(ctx, b, update)

	callbackData := update.CallbackQuery.Data
	message := update.CallbackQuery.Message.Message
	infoOperations, err := dataFromCallback(message, callbackData[len(data.TimeReduce):])
	if err != nil {
		log.Printf("[WARN] Can't get data date")
		return
	}

	// Приводим к uint16
	var result uint16
	if infoOperations.NowDataTime < infoOperations.DeltaTime {
		result = infoOperations.TimeUnit.MaxValue - infoOperations.DeltaTime
	} else {
		result = infoOperations.NowDataTime - infoOperations.DeltaTime
	}

	inlineKeyboard, _ := data.UpdateTimeKeyboard(result, infoOperations.TimeUnit, message.ReplyMarkup)

	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      message.Chat.ID,
		MessageID:   message.ID,
		ReplyMarkup: inlineKeyboard,
	})
	if err != nil {
		log.Printf("[WARN] Can't edit message, err:%v\n", err)
	}
}

func dataFromCallback(message *models.Message, cutCallbackData string) (*timeInfoOperations, error) {
	if message == nil {
		log.Println("[ERROR] Message params is nil")
		return nil, fmt.Errorf("message is nil")
	}

	lineLook := -1 // 0 - min, 1 - hours, 2 - days
	var timeUnit data.TimeUnit
	timeUnitStr := cutCallbackData[2:]

	switch timeUnitStr {
	case data.MinutesUnit.CallBackData:
		lineLook = 0
		timeUnit = data.MinutesUnit
	case data.HoursUnit.CallBackData:
		lineLook = 1
		timeUnit = data.HoursUnit
	case data.DaysUnit.CallBackData:
		lineLook = 2
		timeUnit = data.DaysUnit
	default:
		log.Printf("[ERROR] Unknown time unit type: %v\n", timeUnitStr)
		return nil, fmt.Errorf("unknown time unit type")
	}

	buttons := message.ReplyMarkup.InlineKeyboard
	timeButton := buttons[lineLook][len(buttons[lineLook])/2]

	timedataNow, err := strconv.ParseUint(timeButton.CallbackData[:2], 10, 16)
	if err != nil {
		log.Printf("[ERROR] Can't parse timebutton:%v\n", timeButton.CallbackData)
		return nil, fmt.Errorf("can't parse timebutton")
	}

	deltaTimedata := cutCallbackData[:2]
	timedataDelta, err := strconv.ParseUint(deltaTimedata, 10, 16)
	if err != nil {
		log.Printf("[ERROR] Can't parse delta data\n")
		return nil, fmt.Errorf("can't parse deltatime ")
	}

	return &timeInfoOperations{uint16(timedataNow), uint16(timedataDelta), &timeUnit}, nil
}
