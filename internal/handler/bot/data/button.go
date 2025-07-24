package data

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-telegram/bot/models"
)

const (
	RegionChoose = "chRegion_"
	ExtendKey    = "extKey_"

	CreateKey        = "createKey"
	InfoAboutProject = "infoProject"
	OutlineHelp      = "helpOutline"

	TimeAdd    = "tAdd_"
	TimeReduce = "tReduce_"
	CreateTime = "tCreate_"
)

type TimeUnit struct {
	Name         string
	CallBackData string
	MaxValue     uint16
}

// не больше 99
var (
	MinutesUnit = TimeUnit{"мин", "m", 60}
	HoursUnit   = TimeUnit{"ч", "h", 24}
	DaysUnit    = TimeUnit{"д", "d", 31}
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
	button := models.InlineKeyboardButton{Text: "Узнать об outline", CallbackData: OutlineHelp}

	return button
}

func GetDateFromButton(dateData string) (*TimeDataDuration, error) {
	if len(dateData) < 6 {
		return nil, fmt.Errorf("text is not format: 02d02d02d")
	}

	minutesStr := dateData[0:2]
	minutes, err := strconv.ParseUint(minutesStr, 10, 16)
	if err != nil {
		log.Printf("[WARN] Can't parse minutes err: %v\n", err)
		return nil, fmt.Errorf("can't parse minutes")
	}

	hoursStr := dateData[2:4]
	hours, err := strconv.ParseUint(hoursStr, 10, 16)
	if err != nil {
		log.Printf("[WARN] Can't parse hours err: %v\n", err)
		return nil, fmt.Errorf("can't parse hours")
	}

	daysStr := dateData[4:6]
	days, err := strconv.ParseUint(daysStr, 10, 16)
	if err != nil {
		log.Printf("[WARN] Can't parse days err: %v\n", err)
		return nil, fmt.Errorf("can't parse days")
	}

	return &TimeDataDuration{uint16(minutes), uint16(hours), uint16(days)}, nil
}

func CreateDateButton(keyboardTime *TimeDataDuration) models.InlineKeyboardButton {
	button := models.InlineKeyboardButton{Text: "Создать", CallbackData: fmt.Sprintf("%v%02d%02d%02d", CreateTime, keyboardTime.Minutes, keyboardTime.Hours, keyboardTime.Days)}

	return button
}

func InfoTimeButton(time uint16, timeValue *TimeUnit) models.InlineKeyboardButton {
	button := models.InlineKeyboardButton{Text: fmt.Sprintf("%d%v", time, timeValue.Name), CallbackData: fmt.Sprintf("%02d%v", time, timeValue.CallBackData)}
	return button
}

func TimeMinutesAddButton(minutes uint16) models.InlineKeyboardButton {
	button := timeOperationButton(minutes, &MinutesUnit, true)

	return button
}

func TimeMinutesReduceButton(minutes uint16) models.InlineKeyboardButton {
	button := timeOperationButton(minutes, &MinutesUnit, false)

	return button
}

func TimeHoursAddButton(hours uint16) models.InlineKeyboardButton {
	button := timeOperationButton(hours, &HoursUnit, true)

	return button
}

func TimeHoursReduceButton(hours uint16) models.InlineKeyboardButton {
	button := timeOperationButton(hours, &HoursUnit, false)

	return button
}

func TimeDaysAddButton(days uint16) models.InlineKeyboardButton {
	button := timeOperationButton(days, &DaysUnit, true)

	return button
}

func TimeDaysReduceButton(days uint16) models.InlineKeyboardButton {
	button := timeOperationButton(days, &DaysUnit, false)

	return button
}

func timeOperationButton(time uint16, timeValue *TimeUnit, isAdd bool) models.InlineKeyboardButton {
	if time >= timeValue.MaxValue {
		time = 0
	}

	var sign byte
	var operation string
	if isAdd {
		sign = '+'
		operation = TimeAdd
	} else {
		sign = '-'
		operation = TimeReduce
	}

	button := models.InlineKeyboardButton{Text: fmt.Sprintf("%c%d%v", sign, time, timeValue.Name), CallbackData: fmt.Sprintf("%v%02d%v", operation, time, timeValue.CallBackData)}
	return button
}
