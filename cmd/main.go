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
	"github.com/gin-gonic/gin"
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

	vpnConfigs := vpn.Initialize()
	accountConfigs := account.Initialize()

	// initialize bot
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		// key work
		bot.WithCallbackQueryDataHandler(data.CreateKey, bot.MatchTypeExact, vpnConfigs.Handlers.RegionHandler.GetRegionsInline),
		bot.WithCallbackQueryDataHandler(data.ExtendKey, bot.MatchTypePrefix, vpnConfigs.Handlers.KeyHandler.ExtendKeyIntline),
		bot.WithCallbackQueryDataHandler(data.RegionChoose, bot.MatchTypePrefix, vpnConfigs.Handlers.KeyHandler.CreateKeyGetServerInline),
		bot.WithCallbackQueryDataHandler(data.CreateTime, bot.MatchTypePrefix, vpnConfigs.Handlers.KeyHandler.CreateKeyGetTimeInline),

		bot.WithCallbackQueryDataHandler(data.InfoAboutProject, bot.MatchTypeExact, vpnHandlerBot.InfoAboutInline),
		bot.WithCallbackQueryDataHandler(data.OutlineHelp, bot.MatchTypeExact, vpnHandlerBot.HelpOutlineIntructionInline),

		// time work
		bot.WithCallbackQueryDataHandler(data.TimeAdd, bot.MatchTypePrefix, handler.AddTimeInline),
		bot.WithCallbackQueryDataHandler(data.TimeReduce, bot.MatchTypePrefix, handler.ReduceTimeInline),

		//payment !!!!!
		bot.WithDefaultHandler(accountConfigs.Handlers.WalletHandler.PaymentHandler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, data.START, bot.MatchTypeExact, accountConfigs.Handlers.UserHandler.RegisterUser)
	b.RegisterHandler(bot.HandlerTypeMessageText, data.PAY, bot.MatchTypePrefix, accountConfigs.Handlers.WalletHandler.PayAmount)
	b.RegisterHandler(bot.HandlerTypeMessageText, data.BALANCE, bot.MatchTypeExact, accountConfigs.Handlers.WalletHandler.GetBalace)

	keyScheduler := scheduler.NewScheduler(time.Minute*5, b, vpnConfigs.Services.KeyService)
	keyScheduler.Run(ctx)

	go b.Start(ctx)

	r := gin.Default()
	r.POST("operation/create", accountConfigs.Handlers.OperationHandlerWeb.CreateOperation)

	r.Run(":8080")
}
