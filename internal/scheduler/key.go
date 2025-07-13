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

func NewScheduler(intervalSeconds time.Duration, b *bot.Bot, ctx context.Context, keyService *service.KeyService) *KeyScheduler {
	return &KeyScheduler{
		BotSheduler: BotSheduler{intervalSeconds, b, ctx},
		keyService:  keyService,
	}
}

func (s *KeyScheduler) Run() {
	log.Println("[INFO] scheduler run")
	go func() {
		ticker := time.NewTicker(s.timeInterval)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case t := <-ticker.C:
				log.Printf("[INFO] Check keys into db: %v", t)

				var wg sync.WaitGroup
				wg.Add(2)

				go func() {
					defer wg.Done()
					s.notifyExpired()
				}()

				go func() {
					defer wg.Done()
					s.diactivateExpiredKeys()
				}()

				wg.Wait()
			}
		}
	}()
}

func (s *KeyScheduler) notifyExpired() {
	keysExpiringSoon, err := s.keyService.GetExpiringSoon(s.timeInterval)
	if err != nil || len(keysExpiringSoon) == 0 {
		log.Printf("[WARN] Can't get keys: %v\n", err)
	}

	for _, key := range keysExpiringSoon {
		select {
		case <-s.ctx.Done():
			return
		default:
			chatId := key.User.TelegramID

			handlerBot.SendNotifyAboutDeadline(s.ctx, s.b, chatId, key.ID)
		}
	}
}

func (s *KeyScheduler) diactivateExpiredKeys() {
	keysExpired, err := s.keyService.GetExpiredKeys()
	if err != nil || len(keysExpired) == 0 {
		log.Printf("[WARN] Can't get keys: %v\n", err)
	}

	for _, key := range keysExpired {
		select {
		case <-s.ctx.Done():
			return
		default:
			log.Printf("[INFO] diactivate key #%v", key.ID)

			// TODO edit to change url
			outlineClient := outline.NewOutlineClient(key.Server.Access)

			errOutline := outlineClient.SetDataLimit(key.OutlineKeyId, 0)
			if errOutline != nil {
				log.Printf("[ERROR] Can't change datalimit key #%v", key.ID)
				s.keyService.Delete(key.ID)
				continue
			}

			err := s.keyService.DeactivateKey(key.ID)
			if err != nil {
				log.Printf("[ERROR] Can't diactivate key #%v", key.ID)
				continue
			}
		}
	}

}
