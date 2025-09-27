package vpnfunction

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
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
		button := models.InlineKeyboardButton{Text: region.RegionName, CallbackData: fmt.Sprintf("%v%v", data.RegionChoose, region.ShortName)}
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

func GetProvidersInlineKeyboard(providers []model.Provider) *models.InlineKeyboardMarkup {
	// providers into buttons
	inlineButtons := [][]models.InlineKeyboardButton{}
	line := []models.InlineKeyboardButton{}
	for i, Provider := range providers {
		button := models.InlineKeyboardButton{Text: Provider.Name, CallbackData: fmt.Sprintf("%v%v", data.ProviderChoose, Provider.Name)}
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
