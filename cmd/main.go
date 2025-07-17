package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/database"
	"github.com/RomanGhost/buratino_bot.git/internal/database/repository"
	handlerBot "github.com/RomanGhost/buratino_bot.git/internal/handler/bot"
	"github.com/RomanGhost/buratino_bot.git/internal/handler/bot/data"
	"github.com/RomanGhost/buratino_bot.git/internal/scheduler"
	"github.com/RomanGhost/buratino_bot.git/internal/service"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
	keyHandler := handlerBot.NewKeyHandler(keyService, serverService)
	userHandler := handlerBot.NewUserHandler(userService)

	// initialize bot
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithCallbackQueryDataHandler(data.RegionChoose, bot.MatchTypePrefix, keyHandler.CreateKeyGetServerInline),
		bot.WithCallbackQueryDataHandler(data.ExtendKey, bot.MatchTypePrefix, keyHandler.ExtendKeyIntline),
		bot.WithCallbackQueryDataHandler(data.CreateKey, bot.MatchTypeExact, RegionHandler.GetRegionsInline),

		bot.WithCallbackQueryDataHandler(data.InfoAboutProject, bot.MatchTypeExact, handlerBot.InfoAboutInline),
		bot.WithCallbackQueryDataHandler(data.OutlineHelp, bot.MatchTypeExact, handlerBot.HelpOutlineIntructionInline),
	}

	botToken, exist := os.LookupEnv("BOT_API_TOKEN")
	if !exist {
		panic("Variable not found")
	}
	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	// scheluder init
	keyScheduler := scheduler.NewScheduler(time.Minute*5, b, keyService)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, userHandler.RegisterUser)

	keyScheduler.Run(ctx)
	b.Start(ctx)
}
