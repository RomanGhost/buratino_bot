package scheduler

import (
	"time"

	"github.com/go-telegram/bot"
)

type BotSheduler struct {
	timeInterval time.Duration
	b            *bot.Bot
}
