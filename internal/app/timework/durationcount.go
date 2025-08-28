package timework

import "time"

const DayDuration = 24 * time.Hour
const MonthDuration = 30 * DayDuration

type timeCountDuration struct {
	Minutes time.Duration
	Hours   time.Duration
	Days    time.Duration
	Months  time.Duration
}

func ConcrateDuration(timeDuration time.Duration) *timeCountDuration {
	minutes := (timeDuration % time.Hour) / time.Minute
	hours := (timeDuration % DayDuration) / time.Hour
	days := (timeDuration % MonthDuration) / DayDuration
	months := timeDuration / MonthDuration

	return &timeCountDuration{
		Minutes: minutes,
		Hours:   hours,
		Days:    days,
		Months:  months,
	}
}
