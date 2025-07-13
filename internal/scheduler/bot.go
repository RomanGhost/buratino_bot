package scheduler

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
)

type BotSheduler struct {
	timeInterval time.Duration
	b            *bot.Bot
	ctx          context.Context
}
