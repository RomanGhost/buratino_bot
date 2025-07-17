package data

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"github.com/go-telegram/bot/models"
)

func GetRegionsInlineKeyboard(regions []model.Region) *models.InlineKeyboardMarkup {
	// regions into buttons
	inlineButtons := [][]models.InlineKeyboardButton{}
	line := []models.InlineKeyboardButton{}
	for i, region := range regions {
		if len(region.Servers) == 0 {
			continue
		}
		button := models.InlineKeyboardButton{Text: region.RegionName, CallbackData: fmt.Sprintf("%v%v", RegionChoose, region.ShortName)}
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

	return inlineKeyboard
}

func CreateKeyboard(buttons ...[]models.InlineKeyboardButton) *models.InlineKeyboardMarkup {
	var keyboardButtons [][]models.InlineKeyboardButton
	keyboardButtons = append(keyboardButtons, buttons...)

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: keyboardButtons,
	}
	return &keyboard
}
