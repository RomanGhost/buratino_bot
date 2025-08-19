package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database"
	handlerBot "github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/bot"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/bot/data"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/scheduler"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*// TODO
- [+] Задавать время жизни ключа при его создании
- [ ] Добавить возможность изменять длительность ключа
- [ ] Добавить пользователю монет
- [ ] Личный кабинет
- [ ] Продление ключа сразу если он уже все
- [ ] Возможность продлить написав боту - да
*/

func buildDSN() string {
	host := os.Getenv("DATABASE_ADDR")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//initialize database
	dsn := buildDSN()
	botToken, exist := os.LookupEnv("BOT_API_TOKEN")
	if !exist {
		panic("Variable not found")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.InitDB(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	vpnConfigs := vpn.Initialize(db)

	// initialize bot
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		// key work
		bot.WithCallbackQueryDataHandler(data.CreateKey, bot.MatchTypeExact, vpnConfigs.Handlers.RegionHandler.GetRegionsInline),
		bot.WithCallbackQueryDataHandler(data.ExtendKey, bot.MatchTypePrefix, vpnConfigs.Handlers.KeyHandler.ExtendKeyIntline),
		bot.WithCallbackQueryDataHandler(data.RegionChoose, bot.MatchTypePrefix, vpnConfigs.Handlers.KeyHandler.CreateKeyGetServerInline),
		bot.WithCallbackQueryDataHandler(data.CreateTime, bot.MatchTypePrefix, vpnConfigs.Handlers.KeyHandler.CreateKeyGetTimeInline),

		bot.WithCallbackQueryDataHandler(data.InfoAboutProject, bot.MatchTypeExact, handlerBot.InfoAboutInline),
		bot.WithCallbackQueryDataHandler(data.OutlineHelp, bot.MatchTypeExact, handlerBot.HelpOutlineIntructionInline),

		// time work
		bot.WithCallbackQueryDataHandler(data.TimeAdd, bot.MatchTypePrefix, handlerBot.AddTimeInline),
		bot.WithCallbackQueryDataHandler(data.TimeReduce, bot.MatchTypePrefix, handlerBot.ReduceTimeInline),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, vpnConfigs.Handlers.UserHandler.RegisterUser)

	keyScheduler := scheduler.NewScheduler(time.Minute*5, b, vpnConfigs.Services.KeyService)
	keyScheduler.Run(ctx)

	b.Start(ctx)
}
