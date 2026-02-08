package scheduler

import (
	"time"

	"github.com/go-telegram/bot"
)

type BotScheduler struct {
	timeInterval time.Duration
	b            *bot.Bot
}
