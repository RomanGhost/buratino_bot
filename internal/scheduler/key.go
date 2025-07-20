package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	handlerBot "github.com/RomanGhost/buratino_bot.git/internal/handler/bot"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/outline"
	"github.com/RomanGhost/buratino_bot.git/internal/service"
	"github.com/go-telegram/bot"
)

type KeyScheduler struct {
	BotSheduler
	keyService *service.KeyService
}

func NewScheduler(duration time.Duration, b *bot.Bot, keyService *service.KeyService) *KeyScheduler {
	return &KeyScheduler{
		BotSheduler: BotSheduler{duration, b},
		keyService:  keyService,
	}
}

func (s *KeyScheduler) Run(ctx context.Context) {
	log.Println("[INFO] scheduler run")
	go func() {
		ticker := time.NewTicker(s.timeInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				log.Printf("[INFO] Check keys into db: %v", t)

				var wg sync.WaitGroup
				wg.Add(2)

				go func() {
					defer wg.Done()
					s.notifyExpired(ctx)
				}()

				go func() {
					defer wg.Done()
					s.diactivateExpiredKeys(ctx)
				}()

				wg.Wait()
			}
		}
	}()
}

func (s *KeyScheduler) notifyExpired(ctx context.Context) {
	keysExpiringSoon, err := s.keyService.GetExpiringSoon(s.timeInterval * 2)
	if err != nil || len(keysExpiringSoon) == 0 {
		log.Printf("[WARN] Can't get keys: %v\n", err)
	}

	for _, key := range keysExpiringSoon {
		select {
		case <-ctx.Done():
			return
		default:
			chatId := key.User.TelegramID

			handlerBot.SendNotifyAboutDeadline(ctx, s.b, chatId, key.ID)
		}
	}
}

func (s *KeyScheduler) diactivateExpiredKeys(ctx context.Context) {
	keysExpired, err := s.keyService.GetExpiredKeys()
	if err != nil || len(keysExpired) == 0 {
		log.Printf("[WARN] Can't get keys: %v\n", err)
	}

	for _, key := range keysExpired {
		select {
		case <-ctx.Done():
			return
		default:
			log.Printf("[INFO] diactivate key #%v", key.ID)

			// TODO edit to change url
			outlineClient := outline.NewOutlineClient(key.Server.Access)

			errOutline := outlineClient.SetDataLimit(key.OutlineKeyId, 0)
			if errOutline != nil {
				log.Printf("[ERROR] Can't change datalimit key #%v, err: %v", key.ID, errOutline)
				continue
			}

			err := s.keyService.DeactivateKey(key.ID)
			if err != nil {
				log.Printf("[ERROR] Can't diactivate key #%v", key.ID)
				continue
			}
			s.keyService.Delete(key.ID)
		}
	}

}
