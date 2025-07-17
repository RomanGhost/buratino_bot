package data

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

const (
	RegionChoose     = "chRegion_"
	ExtendKey        = "extKey_"
	CreateKey        = "createKey"
	InfoAboutProject = "infoProject"
	OutlineHelp      = "helpOutline"
)

func CreateKeyButton() models.InlineKeyboardButton {
	button := models.InlineKeyboardButton{Text: "Создать ключ", CallbackData: CreateKey}

	return button
}

func ExtendKeyButton(keyID uint) models.InlineKeyboardButton {
	button := models.InlineKeyboardButton{Text: "Продлить ключ", CallbackData: fmt.Sprintf("%v%d", ExtendKey, keyID)}

	return button
}

func KnowProjectButton() models.InlineKeyboardButton {
	button := models.InlineKeyboardButton{Text: "Узнать о проекте", CallbackData: InfoAboutProject}

	return button
}

func AboutOutlineButton() models.InlineKeyboardButton {
	button := models.InlineKeyboardButton{Text: "Узнать о проекте", CallbackData: InfoAboutProject}

	return button
}
