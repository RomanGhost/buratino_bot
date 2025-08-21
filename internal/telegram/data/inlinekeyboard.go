package data

import (
	"fmt"
	"log"
	"strconv"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"github.com/go-telegram/bot/models"
)

type TimeDataDuration struct {
	Minutes uint16
	Hours   uint16
	Days    uint16
}

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

func GetZeroTimeKeyboard() *models.InlineKeyboardMarkup {
	keyboardDataParam := TimeDataDuration{0, 0, 0}

	return timeChooseInlineKeyboard(&keyboardDataParam)
}

func GetCustomTimeKeyboard(timeDuration *TimeDataDuration) *models.InlineKeyboardMarkup {
	if timeDuration.Minutes >= MinutesUnit.MaxValue {
		timeDuration.Minutes = 0
	}
	if timeDuration.Hours >= HoursUnit.MaxValue {
		timeDuration.Hours = 0
	}
	if timeDuration.Days >= DaysUnit.MaxValue {
		timeDuration.Days = 0
	}

	return timeChooseInlineKeyboard(timeDuration)
}

func UpdateTimeKeyboard(updateValue uint16, updateTimeUnit *TimeUnit, keyboard *models.InlineKeyboardMarkup) (*models.InlineKeyboardMarkup, error) {
	keyboardDataParam, err := getDataTimeFromKeyboard(keyboard)
	if err != nil {
		return nil, fmt.Errorf("error get data from Inline keyboard")
	}

	switch updateTimeUnit.CallBackData {
	case MinutesUnit.CallBackData:
		keyboardDataParam.Minutes = updateValue
	case HoursUnit.CallBackData:
		keyboardDataParam.Hours = updateValue
	case DaysUnit.CallBackData:
		keyboardDataParam.Days = updateValue
	}

	return timeChooseInlineKeyboard(keyboardDataParam), nil

}

func getDataTimeFromKeyboard(keyboard *models.InlineKeyboardMarkup) (*TimeDataDuration, error) {
	minutesStr := keyboard.InlineKeyboard[0][len(keyboard.InlineKeyboard[0])/2].CallbackData[:2]
	minutes, err := strconv.ParseUint(minutesStr, 10, 16)
	if err != nil {
		log.Printf("[WARN] Can't parse minutes err: %v\n", err)
		return nil, fmt.Errorf("can't parse minutes")
	}

	hoursStr := keyboard.InlineKeyboard[1][len(keyboard.InlineKeyboard[1])/2].CallbackData[:2]
	hours, err := strconv.ParseUint(hoursStr, 10, 16)
	if err != nil {
		log.Printf("[WARN] Can't parse hours err: %v\n", err)
		return nil, fmt.Errorf("can't parse hours")
	}

	daysStr := keyboard.InlineKeyboard[2][len(keyboard.InlineKeyboard[2])/2].CallbackData[:2]
	days, err := strconv.ParseUint(daysStr, 10, 16)
	if err != nil {
		log.Printf("[WARN] Can't parse days err: %v\n", err)
		return nil, fmt.Errorf("can't parse days")
	}

	return &TimeDataDuration{uint16(minutes), uint16(hours), uint16(days)}, nil

}

func timeChooseInlineKeyboard(keyboardTime *TimeDataDuration) *models.InlineKeyboardMarkup {
	/*/example format/
	[-5m][time_m][+5m]
	[-1h][time_h][+1h]
	[-1d][time_d][+1d]
	[     create     ]
	//createTime_mmhhdd
	*/
	lineMin := []models.InlineKeyboardButton{TimeMinutesReduceButton(10), TimeMinutesReduceButton(5), InfoTimeButton(keyboardTime.Minutes, &MinutesUnit), TimeMinutesAddButton(5), TimeMinutesAddButton(10)}
	lineHour := []models.InlineKeyboardButton{TimeHoursReduceButton(3), TimeHoursReduceButton(1), InfoTimeButton(keyboardTime.Hours, &HoursUnit), TimeHoursAddButton(1), TimeHoursAddButton(3)}
	lineDays := []models.InlineKeyboardButton{TimeDaysReduceButton(1), InfoTimeButton(keyboardTime.Days, &DaysUnit), TimeDaysAddButton(1)}
	createButtonLine := []models.InlineKeyboardButton{CreateDateButton(keyboardTime)}

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			lineMin,
			lineHour,
			lineDays,
			createButtonLine,
		},
	}

	return &keyboard
}

func CreateKeyboard(buttons ...[]models.InlineKeyboardButton) *models.InlineKeyboardMarkup {
	var keyboardButtons [][]models.InlineKeyboardButton
	keyboardButtons = append(keyboardButtons, buttons...)

	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: keyboardButtons,
	}
	return &keyboard
}
