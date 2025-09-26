package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/account"
	accountHandlerBot "github.com/RomanGhost/buratino_bot.git/internal/account/handler/bot"
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

func initHandlerVPN(s *vpn.Services, as *account.Services) *vpn.Handlers {
	regionHandler := vpnHandlerBot.NewRegionHandler(s.RegionService)
	keyHandler := vpnHandlerBot.NewKeyHandler(s.UserService, s.KeyService, s.ServerService, as.OperationService)
	provider := vpnHandlerBot.NewProviderHandler(s.ProviderService)

	return &vpn.Handlers{
		RegionHandler:   regionHandler,
		KeyHandler:      keyHandler,
		ProviderHandler: provider,
	}
}

func initHandlerAccount(s *account.Services, vpnS *vpn.Services) *account.Handlers {
	userHandler := accountHandlerBot.NewUserHandler(s.UserService, vpnS.UserService)
	walletHandler := accountHandlerBot.NewWalletHandler(s.WalletService, s.OperationService, s.UserService)
	goodsHandler := accountHandlerBot.NewGoodsHandler(s.GoodsService)

	return &account.Handlers{
		UserHandler:   userHandler,
		WalletHandler: walletHandler,
		GoodsHandler:  goodsHandler,
	}
}

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

	accountHandlers := initHandlerAccount(accountServices, vpnServices)
	vpnHandlers := initHandlerVPN(vpnServices, accountServices)

	// initialize bot
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		// key work
		bot.WithCallbackQueryDataHandler(data.CreateKeyRequest, bot.MatchTypeExact, vpnHandlers.RegionHandler.GetRegionsInline),                                                // first - get request for create key -> send to get the region
		bot.WithCallbackQueryDataHandler(data.RegionChoose, bot.MatchTypePrefix, vpnHandlers.KeyHandler.GetRegionSendProvider(vpnHandlers.ProviderHandler.GetProvidersInline)), // get region, send a request of the Provider
		bot.WithCallbackQueryDataHandler(data.ProviderChoose, bot.MatchTypePrefix, vpnHandlers.KeyHandler.GetProviderSendTime(vpnHandlerBot.KeyboardTimeChoose)),               // get provider, send a request of the time
		bot.WithCallbackQueryDataHandler(data.TimeChoose, bot.MatchTypePrefix, vpnHandlers.KeyHandler.GetTimeToCreateKey(vpnHandlers.KeyHandler.CreateKey)),

		bot.WithCallbackQueryDataHandler(data.ExtendKey, bot.MatchTypePrefix, vpnHandlers.KeyHandler.ExtendKeyIntline),

		bot.WithCallbackQueryDataHandler(data.InfoAboutProject, bot.MatchTypeExact, vpnHandlerBot.InfoAboutInline),
		bot.WithCallbackQueryDataHandler(data.OutlineHelp, bot.MatchTypeExact, vpnHandlerBot.HelpOutlineIntructionInline),

		// time work
		bot.WithCallbackQueryDataHandler(data.TimeAdd, bot.MatchTypePrefix, handler.AddTimeInline),
		bot.WithCallbackQueryDataHandler(data.TimeReduce, bot.MatchTypePrefix, handler.ReduceTimeInline),

		//payment !!!!!
		bot.WithDefaultHandler(accountHandlers.WalletHandler.PaymentHandler),

		// lookup
		bot.WithMiddlewares(accountHandlers.UserHandler.MiddleWareLookup),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, data.START, bot.MatchTypeExact, accountHandlers.UserHandler.RegisterUser)
	b.RegisterHandler(bot.HandlerTypeMessageText, data.PAY, bot.MatchTypePrefix, accountHandlers.WalletHandler.PayAmount)
	b.RegisterHandler(bot.HandlerTypeMessageText, data.BALANCE, bot.MatchTypeExact, accountHandlers.WalletHandler.GetBalace)
	b.RegisterHandler(bot.HandlerTypeMessageText, data.PRICES, bot.MatchTypeExact, accountHandlers.GoodsHandler.GetPrices)

	keyScheduler := scheduler.NewScheduler(time.Minute*5, b, vpnServices.KeyService)
	keyScheduler.Run(ctx)

	b.Start(ctx)
}
