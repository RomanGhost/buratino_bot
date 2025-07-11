package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/database/repository"
	handlerBot "github.com/RomanGhost/buratino_bot.git/internal/handler/bot"
	"github.com/go-telegram/bot"
)

type KeyScheduler struct {
	BotSheduler
	keyRepository *repository.KeyRepository
}

func NewScheduler(keyRepository *repository.KeyRepository, intervalSeconds time.Duration, b *bot.Bot, ctx context.Context) {

}

func (s *KeyScheduler) Run() {
	go func() {
		ticker := time.NewTicker(time.Duration(s.timeInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case t := <-ticker.C:
				log.Printf("[INFO] Check keys into db: %v", t)

				// notify
				keysExpiringSoon, err := s.keyRepository.GetExpiringSoon(s.timeInterval)
				if err != nil || len(keysExpiringSoon) == 0 {
					log.Printf("[WARN] Can't get keys: %v\n", err)
				}

				// close keys
				keysExpired, err := s.keyRepository.GetExpiredKeys()
				if err != nil || len(keysExpired) == 0 {
					log.Printf("[WARN] Can't get keys: %v\n", err)
				}

				var wg sync.WaitGroup
				wg.Add(2)

				go func() {
					defer wg.Done()
					for _, key := range keysExpiringSoon {
						select {
						case <-s.ctx.Done():
							return
						default:
							chatId := key.User.TelegramID

							handlerBot.SendNotifyAboutDeadline(s.ctx, s.b, chatId)
						}
					}
				}()

				go func() {
					defer wg.Done()
					for _, key := range keysExpired {
						select {
						case <-s.ctx.Done():
							return
						default:
							log.Printf("[INFO] diactivate key #%v", key.ID)
							s.keyRepository.DeactivateKey(key.ID)
						}
					}
				}()

				wg.Wait()

			}
		}
	}()
}
