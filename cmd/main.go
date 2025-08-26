package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/account"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/handler"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn"
	vpnHandlerBot "github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/bot"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/scheduler"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

/*// TODO
- [+] Задавать время жизни ключа при его создании
- [ ] Добавить возможность изменять длительность ключа
- [ ] Добавить пользователю монет
- [ ] Личный кабинет
- [ ] Продление ключа сразу если он уже все
- [ ] Возможность продлить написав боту - да
*/

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken, exist := os.LookupEnv("BOT_API_TOKEN")
	if !exist {
		panic("Variable not found")
	}

	vpnRepositories := vpn.InitializeRepository()
	vpnServices := vpn.InitService(vpnRepositories)

	accountRepositories := account.InitializeRepository()
	accountServices := account.InitService(accountRepositories)

	accountHandlers := account.InitHandler(accountServices)
	vpnHandlers := vpn.InitHandler(vpnServices, accountServices)

	// initialize bot
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		// key work
		bot.WithCallbackQueryDataHandler(data.CreateKey, bot.MatchTypeExact, vpnHandlers.RegionHandler.GetRegionsInline),
		bot.WithCallbackQueryDataHandler(data.ExtendKey, bot.MatchTypePrefix, vpnHandlers.KeyHandler.ExtendKeyIntline),
		bot.WithCallbackQueryDataHandler(data.RegionChoose, bot.MatchTypePrefix, vpnHandlers.KeyHandler.CreateKeyGetServerInline),
		bot.WithCallbackQueryDataHandler(data.CreateTime, bot.MatchTypePrefix, vpnHandlers.KeyHandler.CreateKeyGetTimeInline),

		bot.WithCallbackQueryDataHandler(data.InfoAboutProject, bot.MatchTypeExact, vpnHandlerBot.InfoAboutInline),
		bot.WithCallbackQueryDataHandler(data.OutlineHelp, bot.MatchTypeExact, vpnHandlerBot.HelpOutlineIntructionInline),

		// time work
		bot.WithCallbackQueryDataHandler(data.TimeAdd, bot.MatchTypePrefix, handler.AddTimeInline),
		bot.WithCallbackQueryDataHandler(data.TimeReduce, bot.MatchTypePrefix, handler.ReduceTimeInline),

		//payment !!!!!
		bot.WithDefaultHandler(accountHandlers.WalletHandler.PaymentHandler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, data.START, bot.MatchTypeExact, accountHandlers.UserHandler.RegisterUser)
	b.RegisterHandler(bot.HandlerTypeMessageText, data.PAY, bot.MatchTypePrefix, accountHandlers.WalletHandler.PayAmount)
	b.RegisterHandler(bot.HandlerTypeMessageText, data.BALANCE, bot.MatchTypeExact, accountHandlers.WalletHandler.GetBalace)

	keyScheduler := scheduler.NewScheduler(time.Minute*5, b, vpnServices.KeyService)
	keyScheduler.Run(ctx)

	b.Start(ctx)
}
