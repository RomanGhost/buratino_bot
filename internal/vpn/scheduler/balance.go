package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/go-telegram/bot"
)

const timeInterval = 12 * time.Hour

type BalanceScheduler struct {
	BotSheduler
	operationService *service.OperationService
	userService      *service.UserService
}

func NewBalanceScheduler(b *bot.Bot, operationService *service.OperationService, userService *service.UserService) *BalanceScheduler {
	return &BalanceScheduler{
		BotSheduler:      BotSheduler{timeInterval, b},
		operationService: operationService,
		userService:      userService,
	}
}

func (s *BalanceScheduler) Run(ctx context.Context) {
	log.Println("[INFO] account scheduler run")

	go func() {
		now := time.Now()

		nextRun := nextNoonOrMidnight(now)
		delay := time.Until(nextRun)

		log.Printf("[INFO] First scheduler run at %v", nextRun)

		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			s.safeRun()
		}

		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				log.Printf("[INFO] Scheduler tick at %v", t)
				s.safeRun()
			}
		}
	}()
}

func (s *BalanceScheduler) safeRun() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] scheduler panic: %v", r)
		}
	}()

	log.Println("[INFO] Check keys into db")
	s.topUpForUsers()
}

func (s *BalanceScheduler) topUpForUsers() {
	users, err := s.userService.GetActiveUser()
	if err != nil {
		log.Println("[ERROR] Can't get users from database, error: ", err)
	}

	price, err := s.operationService.GetPrice(model.VPN1Day.Name, 30)
	if err != nil {
		log.Println("[ERROR] Can't get price from database, error: ", err)
	}

	for _, user := range users {
		s.operationService.TopUpAccount(user.ID, uint64(price/1000), uint64(price%1000))
	}
}

func nextNoonOrMidnight(now time.Time) time.Time {
	loc := now.Location()

	todayMidnight := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, loc,
	)

	todayNoon := todayMidnight.Add(12 * time.Hour)

	switch {
	case now.Before(todayMidnight):
		return todayMidnight
	case now.Before(todayNoon):
		return todayNoon
	default:
		return todayMidnight.Add(24 * time.Hour)
	}
}
