package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/database"
	"github.com/RomanGhost/buratino_bot.git/internal/database/repository"
	handlerBot "github.com/RomanGhost/buratino_bot.git/internal/handler/bot"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/outline"
	"github.com/RomanGhost/buratino_bot.git/internal/scheduler"
	"github.com/RomanGhost/buratino_bot.git/internal/service"
	"github.com/go-telegram/bot"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// outline init
	httpUrl := "https://77.233.215.100:3411/g2G6SIZWzAPcXeFVjO_78A"
	outlineClient := outline.NewOutlineClient(httpUrl)

	//initialize database
	dsn := "host=localhost user=main_telegram_user password=jfsdlkfsur3432fd dbname=buratino_vpn port=5434 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.InitDB(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// repository init
	keyRepository := repository.NewKeyRepository(db)
	userRepository := repository.NewUserRepository(db)
	userRoleRepository := repository.NewUserRoleRepository(db)
	serverRepository := repository.NewServerRepository(db)
	regionRepository := repository.NewRegionRepository(db)

	// service init
	keyService := service.NewKeyService(keyRepository, userRepository, serverRepository)
	userService := service.NewUserService(userRepository, userRoleRepository)
	regionService := service.NewRegionService(regionRepository)
	serverService := service.NewServerService(serverRepository)

	// handler init
	RegionHandler := handlerBot.NewRegionHandler(regionService)
	keyHandler := handlerBot.NewKeyHandler(outlineClient, keyService, serverService)
	userHandler := handlerBot.NewUserHandler(userService)

	// initialize bot
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithCallbackQueryDataHandler(handlerBot.RegionChoose, bot.MatchTypePrefix, keyHandler.CreateKeyGetServerInline),
		bot.WithCallbackQueryDataHandler(handlerBot.ExtendKey, bot.MatchTypePrefix, keyHandler.ExtendKeyIntline),
		bot.WithCallbackQueryDataHandler(handlerBot.CreateKey, bot.MatchTypeExact, RegionHandler.GetRegionsInline),

		bot.WithCallbackQueryDataHandler(handlerBot.InfoAboutProject, bot.MatchTypeExact, handlerBot.InfoAboutInline),
		bot.WithCallbackQueryDataHandler(handlerBot.OutlineHelp, bot.MatchTypeExact, handlerBot.HelpOutlineIntructionInline),
	}

	b, err := bot.New("7786090535:AAGg1aj6SkJwc6mURapwQ7AYf4hmRo-ynAE", opts...)
	if err != nil {
		panic(err)
	}

	// scheluder init
	keyScheduler := scheduler.NewScheduler(time.Minute*5, b, ctx, keyService, outlineClient)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, userHandler.RegisterUser)

	keyScheduler.Run()
	b.Start(ctx)
}
