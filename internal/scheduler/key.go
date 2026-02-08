package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	handlerBot "github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/bot"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/provider"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/service"
	"github.com/go-telegram/bot"
)

type KeyScheduler struct {
	BotScheduler
	keyService *service.KeyService
}

func NewKeyScheduler(duration time.Duration, b *bot.Bot, keyService *service.KeyService) *KeyScheduler {
	return &KeyScheduler{
		BotScheduler: BotScheduler{duration, b},
		keyService:   keyService,
	}
}

func (s *KeyScheduler) Run(ctx context.Context) {
	log.Println("[INFO] Key scheduler run")
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
					s.deactivateExpiredKeys(ctx)
				}()

				wg.Wait()
			}
		}
	}()
}

func (s *KeyScheduler) notifyExpired(ctx context.Context) {
	// Получить истекут в промежуток шедулера
	keysExpiringSoon, err := s.keyService.GetExpiringSoon(s.timeInterval * 2)
	if err != nil || len(keysExpiringSoon) == 0 {
		log.Printf("[WARN] Can't get keys: %v\n", err)
	}

	nowTime := time.Now().UTC()
	deadlineKeyData := make(map[time.Duration][]Notify)

	for _, key := range keysExpiringSoon {
		select {
		case <-ctx.Done():
			return
		default:

			delta := key.DeadlineTime.Sub(nowTime)
			deadlinePercent := time.Duration(float64(key.Duration) * 0.1)
			if delta < deadlinePercent {
				continue
			}

			val, ok := deadlineKeyData[delta-s.timeInterval]
			if !ok {
				deadlineKeyData[delta-s.timeInterval] = []Notify{}
			}

			chatID := key.User.TelegramID
			dataNotify := Notify{ChatID: chatID, KeyID: key.ID}
			deadlineKeyData[delta-s.timeInterval] = append(val, dataNotify)
		}
	}
	s.notify(ctx, deadlineKeyData)
}

func (s *KeyScheduler) notify(ctx context.Context, timersData map[time.Duration][]Notify) {
	var wg sync.WaitGroup
	for durationTime, notifyData := range timersData {
		if ctx.Err() != nil {
			return
		}

		wg.Add(1)
		go func(d time.Duration, data []Notify) {
			defer wg.Done()

			select {
			case <-time.After(d):
				// Продолжаем работу после ожидания
			case <-ctx.Done():
				return
			}

			// Последовательно отправляем уведомления с проверкой контекста
			for _, nd := range data {
				if ctx.Err() != nil {
					return
				}

				handlerBot.SendNotifyAboutDeadline(ctx, s.b, nd.ChatID, nd.KeyID)
			}
		}(durationTime, notifyData)
	}
	wg.Wait()
}

func (s *KeyScheduler) deactivateExpiredKeys(ctx context.Context) {
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
			providerClient := provider.NewProvider(key.Server.Access, key.Server.ProviderID)

			// errOutline := providerClient.DeleteAccessKey
			// deactivate into vpn provider
			errOutline := providerClient.DisableKey(key.KeyID)
			if errOutline != nil {
				log.Printf("[ERROR] Can't deactivate key #%v, err: %v", key.ID, errOutline)
				_ = s.keyService.Delete(key.ID)
				continue
			}

			// deactivate in db
			err := s.keyService.DeactivateKey(key.ID)
			if err != nil {
				log.Printf("[ERROR] Can't deactivate key #%v", key.ID)
				_ = s.keyService.Delete(key.ID)
				continue
			}
		}
	}

}
