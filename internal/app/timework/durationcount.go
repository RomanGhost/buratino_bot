package timework

import "time"

const DayDuration = 24 * time.Hour
const MonthDuration = 30 * DayDuration

type timeCountDuration struct {
	Minutes int64
	Hours   int64
	Days    int64
	Months  int64
}

func ConcrateDuration(timeDuration time.Duration) *timeCountDuration {
	minutes := (timeDuration % time.Hour) / time.Minute
	hours := (timeDuration % DayDuration) / time.Hour
	days := (timeDuration % MonthDuration) / DayDuration
	months := timeDuration / MonthDuration

	return &timeCountDuration{
		Minutes: int64(minutes),
		Hours:   int64(hours),
		Days:    int64(days),
		Months:  int64(months),
	}
}
